package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"strings"
)

type Validator struct{ cfg Config }

func NewValidator(cfg Config) *Validator { return &Validator{cfg: cfg} }

func (v *Validator) ValidateRedirectURI(client *Client, ru string) bool {
	u, err := url.Parse(ru)
	if err != nil {
		return false
	}
	// Allow loopback with any port and fixed path /callback
	if (u.Hostname() == "127.0.0.1" || u.Hostname() == "[::1]" || u.Hostname() == "::1") && u.Path == "/callback" && (u.Scheme == "http" || u.Scheme == "https") {
		return true
	}
	// Also allow any exact match from AllowedRedirects
	for _, ar := range client.AllowedRedirects {
		if ru == ar {
			return true
		}
	}
	return false
}

func (v *Validator) ValidatePKCE(codeVerifier, codeChallenge, method string) bool {
	if method != "S256" {
		return false
	}
	sum := sha256.Sum256([]byte(codeVerifier))
	want := base64.RawURLEncoding.EncodeToString(sum[:])
	return subtleConstEq(want, codeChallenge)
}

func subtleConstEq(a, b string) bool {
	// constant-time-ish comparison for equal-length strings
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}
	return v == 0
}

func (v *Validator) SplitScopes(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Fields(s)
	return parts
}
