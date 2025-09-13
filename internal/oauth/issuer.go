package oauth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"
)

type Issuer struct{ cfg Config }

func NewIssuer(cfg Config) *Issuer { return &Issuer{cfg: cfg} }

func (i *Issuer) MintAccessToken(sub, aud, scope string) (string, error) {
	now := time.Now().Unix()
	exp := time.Now().Add(i.cfg.TokenTTL).Unix()
	header := map[string]any{
		"alg": "RS256",
		"kid": i.cfg.KeyID,
		"typ": "JWT",
	}
	claims := map[string]any{
		"iss":   i.cfg.Issuer,
		"sub":   sub,
		"aud":   aud,
		"iat":   now,
		"exp":   exp,
		"scope": scope,
	}
	return signJWT(i.cfg.PrivateKey, header, claims)
}

func signJWT(priv *rsa.PrivateKey, header, claims map[string]any) (string, error) {
	h, _ := json.Marshal(header)
	c, _ := json.Marshal(claims)
	enc := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
	unsigned := enc(h) + "." + enc(c)
	sum := sha256.Sum256([]byte(unsigned))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, sum[:])
	if err != nil {
		return "", err
	}
	token := unsigned + "." + enc(sig)
	return token, nil
}

// JWKS produces a minimal JWK Set for the current RSA key.
func (i *Issuer) JWKS() map[string]any {
	pub := i.cfg.PrivateKey.Public().(*rsa.PublicKey)
	// n and e as base64url
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(bigEndianBytes(pub.E))
	jwk := map[string]any{
		"kty": "RSA",
		"use": "sig",
		"alg": "RS256",
		"kid": i.cfg.KeyID,
		"n":   n,
		"e":   e,
	}
	return map[string]any{"keys": []any{jwk}}
}

func bigEndianBytes(e int) []byte {
	// minimal bytes for public exponent
	if e == 0 {
		return []byte{0}
	}
	var b []byte
	for v := e; v > 0; v >>= 8 {
		b = append([]byte{byte(v & 0xff)}, b...)
	}
	// strip leading zeros is not needed as we assembled minimal
	return b
}
