package oauth

import (
    "crypto"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "errors"
    "strings"
    "time"
)

type TokenClaims struct {
    Iss string `json:"iss"`
    Sub string `json:"sub"`
    Aud any    `json:"aud"` // string or []string
    Iat int64  `json:"iat"`
    Exp int64  `json:"exp"`
    Scope string `json:"scope"`
}

func parseAud(a any) []string {
    switch v := a.(type) {
    case string:
        if v == "" { return nil }
        return []string{v}
    case []any:
        out := make([]string, 0, len(v))
        for _, e := range v {
            if s, ok := e.(string); ok { out = append(out, s) }
        }
        return out
    default:
        return nil
    }
}

// VerifyAndExtract verifies a JWT (RS256) and returns parsed claims.
func VerifyAndExtract(cfg Config, token string) (TokenClaims, error) {
    parts := strings.Split(token, ".")
    if len(parts) != 3 { return TokenClaims{}, errors.New("invalid token format") }
    headerB, err := base64.RawURLEncoding.DecodeString(parts[0])
    if err != nil { return TokenClaims{}, err }
    _ = headerB // we could check alg=RS256, kid=...
    payloadB, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil { return TokenClaims{}, err }
    sigB, err := base64.RawURLEncoding.DecodeString(parts[2])
    if err != nil { return TokenClaims{}, err }

    // Verify signature
    sum := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
    pub := cfg.PrivateKey.Public().(*rsa.PublicKey)
    if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, sum[:], sigB); err != nil {
        return TokenClaims{}, errors.New("invalid signature")
    }

    var claims TokenClaims
    if err := json.Unmarshal(payloadB, &claims); err != nil {
        return TokenClaims{}, err
    }
    // Validate iss/exp/iat minimally here; audience/scope per handler
    if claims.Iss != cfg.Issuer {
        return TokenClaims{}, errors.New("issuer mismatch")
    }
    now := time.Now().Unix()
    if claims.Exp <= now {
        return TokenClaims{}, errors.New("token expired")
    }
    // iat: allow skew (not strictly enforced here)
    return claims, nil
}

func AudienceContains(claims TokenClaims, want string) bool {
    for _, a := range parseAud(claims.Aud) {
        if a == want { return true }
    }
    return false
}

func ScopeContains(claims TokenClaims, scope string) bool {
    if claims.Scope == "" { return false }
    parts := strings.Fields(claims.Scope)
    for _, s := range parts {
        if s == scope { return true }
    }
    return false
}
