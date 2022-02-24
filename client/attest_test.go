package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-tpm-tools/internal/test"
)

// Returns an x509 Certificate with the provided issuingURL and signed with the provided parent certificate and key.
// If parentCert and parentKey are nil, the certificate will be self-signed.
func getTestCert(t *testing.T, issuingURL []string, parentCert *x509.Certificate, parentKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey) {
	t.Helper()

	certKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
		IssuingCertificateURL: issuingURL,
	}

	if parentCert == nil && parentKey == nil {
		parentCert = template
		parentKey = certKey
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, parentCert, certKey.Public(), parentKey)
	if err != nil {
		t.Fatalf("Unable to create test certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		t.Fatalf("Unable to parse test certificate: %v", err)
	}

	return cert, certKey
}

func TestFetchIssuingCertificateSucceeds(t *testing.T) {
	testCA, caKey := getTestCert(t, nil, nil, nil)

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write(testCA.Raw)
	}))
	defer ts.Close()

	leafCert, _ := getTestCert(t, []string{"invalid.URL", ts.URL}, testCA, caKey)

	_, err := fetchIssuingCertificate(leafCert)
	if err != nil {
		t.Errorf("fetchIssuingCertificate() returned error: %v", err)
	}
}

func TestFetchIssuingCertificateReturnsErrorIfNoValidCertificateFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("these are some random bytes"))
	}))
	defer ts.Close()

	testCA, caKey := getTestCert(t, nil, nil, nil)
	leafCert, _ := getTestCert(t, []string{ts.URL}, testCA, caKey)

	_, err := fetchIssuingCertificate(leafCert)
	if err == nil {
		t.Error("fetchIssuingCertificate returned successfully, but expected error.")
	}
}

func TestGetCertificateChainSucceeds(t *testing.T) {
	// Create CA and corresponding server.
	testCA, caKey := getTestCert(t, nil, nil, nil)

	caServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write(testCA.Raw)
	}))

	defer caServer.Close()

	// Create intermediate cert and corresponding server.
	intermediateCert, intermediateKey := getTestCert(t, []string{caServer.URL}, testCA, caKey)

	intermediateServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write(intermediateCert.Raw)
	}))
	defer intermediateServer.Close()

	// Create leaf cert.
	leafCert, _ := getTestCert(t, []string{intermediateServer.URL}, intermediateCert, intermediateKey)

	key := &Key{cert: leafCert}

	certChain, err := key.getCertificateChain()
	if err != nil {
		t.Fatalf("getCertificateChain returned error: %v", err)
	}

	if len(certChain) != 2 {
		t.Fatalf("getCertificateChain did not return the expected number of certificates: got %v, want 2", len(certChain))
	}
}

func TestGetCertificateChainFailsIfCertHasTooManyIssuingURLs(t *testing.T) {
	// Create CA and corresponding server so the only failure point is the IssuingCertificateURL length check.
	testCA, caKey := getTestCert(t, nil, nil, nil)

	caServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write(testCA.Raw)
	}))
	defer caServer.Close()

	leafCert, _ := getTestCert(t, []string{caServer.URL, "url2", "url3", "url4"}, testCA, caKey)

	key := &Key{cert: leafCert}

	_, err := key.getCertificateChain()
	if err == nil {
		t.Error("getCertificateChain was successful, expected error")
	}
}

func TestKeyAttestSucceedsWithCertChainRetrieval(t *testing.T) {
	testCA, caKey := getTestCert(t, nil, nil, nil)

	caServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write(testCA.Raw)
	}))

	defer caServer.Close()

	leafCert, _ := getTestCert(t, []string{caServer.URL}, testCA, caKey)

	rwc := test.GetTPM(t)
	defer CheckedClose(t, rwc)

	ak, err := AttestationKeyRSA(rwc)
	if err != nil {
		t.Fatalf("Failed to generate test AK: %v", err)
	}

	ak.cert = leafCert

	attestation, err := ak.Attest(AttestOpts{Nonce: []byte("some nonce"), FetchCertChain: true})
	if err != nil {
		t.Fatalf("Attest returned with error: %v", err)
	}

	// Expect one cert retrieved.
	if len(attestation.IntermediateCerts) != 1 {
		t.Fatalf("Got %v intermediate certs, want 1.", len(attestation.IntermediateCerts))
	}

	if !bytes.Equal(attestation.IntermediateCerts[0], testCA.Raw) {
		t.Errorf("Attestation does not contain the expected intermediate cert: got %v, want %v", attestation.IntermediateCerts[0], testCA.Raw)
	}
}

func TestKeyAttestFailsForFetchCertChainWithNoAKCert(t *testing.T) {
	rwc := test.GetTPM(t)
	defer CheckedClose(t, rwc)

	ak, err := AttestationKeyRSA(rwc)
	if err != nil {
		t.Fatalf("Failed to generate test AK: %v", err)
	}

	_, err = ak.Attest(AttestOpts{Nonce: []byte("some nonce"), FetchCertChain: true})
	if err == nil {
		t.Errorf("Attest returned successfully, expected error due to nil AK cert.")
	}

}
