package httpmw

import (
    "net/http"
    "strings"
)

func ExtractBearer(r *http.Request) (string, bool) {
    h := r.Header.Get("Authorization")
    if h == "" { return "", false }
    if !strings.HasPrefix(h, "Bearer ") && !strings.HasPrefix(h, "bearer ") { return "", false }
    return strings.TrimSpace(h[len("Bearer "):]), true
}

func WriteWWWAuth(w http.ResponseWriter, status int, params map[string]string) {
    // Build: Bearer realm="resource", error="invalid_token", error_description="...", scope="mcp.read"
    parts := []string{"Bearer"}
    for k, v := range params {
        if v == "" { continue }
        parts = append(parts, k+"=\""+v+"\"")
    }
    w.Header().Set("WWW-Authenticate", strings.Join(parts, " "))
    w.WriteHeader(status)
}

