package oauth

import (
	"crypto/rsa"
	"time"
)

// Config holds server-wide settings and signing material.
type Config struct {
    Issuer     string
    TokenTTL   time.Duration
    CodeTTL    time.Duration
    KeyID      string
    PrivateKey *rsa.PrivateKey
    // ResourceID is the audience identifier that access tokens must target
    // to be accepted by the protected resource in this PoC.
    ResourceID string
}

// CodeRecord stores authorization code binding info for PKCE verification and token minting.
type CodeRecord struct {
	ClientID            string
	RedirectURI         string
	Scope               string
	Resource            string
	CodeChallenge       string
	CodeChallengeMethod string // expect "S256"
	ExpiresAt           time.Time
	Used                bool
}

// Client represents a statically-registered OAuth client.
type Client struct {
	ID               string
	AllowedRedirects []string // patterns or exact URIs
	Public           bool     // true: no secret
}
