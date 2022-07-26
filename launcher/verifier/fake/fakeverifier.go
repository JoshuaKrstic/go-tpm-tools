// Package fakeverifier is a fake implementation of the AttestationVerifier for testing.
package fakeverifier

import (
	"context"
	"crypto"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-tpm-tools/launcher/verifier"
)

type fakeClient struct {
	//pbClient servpb.AttestationVerifierClient
	signer crypto.Signer
}

func NewClient(signer crypto.Signer) verifier.Client {
	return &fakeClient{signer}
}

// CreateChallenge returns a hard coded, basic challenge.
//
// If you have found this method is insufficient for your tests, this class must be updated to
// allow for better testing.
func (fc *fakeClient) CreateChallenge(ctx context.Context) (*verifier.Challenge, error) {
	return &verifier.Challenge{
		Name:  "FakeName",
		Nonce: []byte{0x0},
	}, nil
}

// VerifyAttestation does basic checks and returns a hard coded attestation response.
//
// If you have found this method is insufficient for your tests, this class must be updated to
// allow for better testing.
func (fc *fakeClient) VerifyAttestation(ctx context.Context, request verifier.VerifyAttestationRequest) (*verifier.VerifyAttestationResponse, error) {
	// Determine signing algorithm.
	signingMethod := jwt.SigningMethodRS256
	now := jwt.TimeFunc().Unix()
	claims := jwt.StandardClaims{
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: now + 360,
		Audience:  "TestingAudience",
		Issuer:    "TestingIssuer",
		Subject:   "TestingSubject",
	}

	token := jwt.NewWithClaims(signingMethod, claims)

	// Instead of a private key, provide the signer.
	signed, _ := token.SignedString(fc.signer)

	response := verifier.VerifyAttestationResponse{
		ClaimsToken: []byte(signed),
	}

	return &response, nil
}
