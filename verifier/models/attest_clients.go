// Package models contains models needed in client and server
package models

// AttestClients contains clients for supported verifier services that can be used to
// get attestation tokens.
type AttestClients struct {
	GCA Client
	ITA Client
}

// HasThirdPartyClient returns true if AttestClients contains an initialzied
// third-party verifier client.
func (ac *AttestClients) HasThirdPartyClient() bool {
	return ac.ITA != nil
}
