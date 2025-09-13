package wellknown

import (
	"encoding/json"
	"net/http"

	oh "mcp-oauth-poc/internal/oauth"
)

func AuthorizationServerMetadataHandler(cfg oh.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		issuer := cfg.Issuer
		resp := map[string]any{
			"issuer":                                issuer,
			"authorization_endpoint":                issuer + "/oauth2/authorize",
			"token_endpoint":                        issuer + "/oauth2/token",
			"jwks_uri":                              issuer + "/.well-known/jwks.json",
			"grant_types_supported":                 []string{"authorization_code"},
			"code_challenge_methods_supported":      []string{"S256"},
			"token_endpoint_auth_methods_supported": []string{"none"},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func JWKSHandler(cfg oh.Config) http.HandlerFunc {
	iss := oh.NewIssuer(cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(iss.JWKS())
	}
}
