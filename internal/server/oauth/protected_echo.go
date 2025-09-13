package oauth

import (
    "encoding/json"
    "net/http"

    httpmw "mcp-oauth-poc/internal/server/http"
)

// ProtectedEchoHandler serves GET /resource/echo, requiring mcp.read scope and proper audience.
func ProtectedEchoHandler(cfg Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        token, ok := httpmw.ExtractBearer(r)
        if !ok || token == "" {
            httpmw.WriteWWWAuth(w, http.StatusUnauthorized, map[string]string{
                "realm": "resource",
                "error": "invalid_token",
            })
            return
        }
        claims, err := VerifyAndExtract(cfg, token)
        if err != nil {
            httpmw.WriteWWWAuth(w, http.StatusUnauthorized, map[string]string{
                "realm":              "resource",
                "error":              "invalid_token",
                "error_description":  err.Error(),
            })
            return
        }
        if !AudienceContains(claims, cfg.ResourceID) {
            httpmw.WriteWWWAuth(w, http.StatusUnauthorized, map[string]string{
                "realm":             "resource",
                "error":             "invalid_token",
                "error_description": "audience mismatch",
            })
            return
        }
        if !ScopeContains(claims, "mcp.read") {
            httpmw.WriteWWWAuth(w, http.StatusForbidden, map[string]string{
                "realm": "resource",
                "error": "insufficient_scope",
                "scope": "mcp.read",
            })
            return
        }
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(map[string]any{
            "sub":   claims.Sub,
            "aud":   claims.Aud,
            "scope": claims.Scope,
            "exp":   claims.Exp,
        })
    }
}

