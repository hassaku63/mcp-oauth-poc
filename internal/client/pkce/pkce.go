package pkce

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "time"

    "mcp-oauth-poc/internal/client/store"
)

func randomB64URL(n int) (string, error) {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}

func s256(verifier string) string {
    sum := sha256.Sum256([]byte(verifier))
    return base64.RawURLEncoding.EncodeToString(sum[:])
}

// NewSession creates a new session with PKCE and state.
func NewSession(redirectURI, resource, scope, clientID string) (store.Session, error) {
    // 64 random bytes -> 86 chars base64url, within 43-128
    verifier, err := randomB64URL(64)
    if err != nil {
        return store.Session{}, err
    }
    state, err := randomB64URL(32)
    if err != nil {
        return store.Session{}, err
    }
    sid, err := randomB64URL(16)
    if err != nil {
        return store.Session{}, err
    }
    return store.Session{
        ID:            "sess_" + sid,
        ClientID:      clientID,
        RedirectURI:   redirectURI,
        Resource:      resource,
        Scope:         scope,
        CodeVerifier:  verifier,
        CodeChallenge: s256(verifier),
        State:         state,
        CreatedAt:     time.Now().UTC(),
        ExpiresInSec:  600,
    }, nil
}

