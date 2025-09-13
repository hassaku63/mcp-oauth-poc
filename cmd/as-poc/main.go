package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"
	"time"

	oh "mcp-oauth-poc/internal/oauth"
	wh "mcp-oauth-poc/internal/wellknown"
)

func main() {
	// Initialize signing key (RSA 2048) and issuer config
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	cfg := oh.Config{
		Issuer:     "http://localhost:8080",
		TokenTTL:   1 * time.Hour,
		CodeTTL:    2 * time.Minute,
		KeyID:      "poc-key-1",
		PrivateKey: priv,
	}

	store := oh.NewStore(cfg)
	issuer := oh.NewIssuer(cfg)
	validator := oh.NewValidator(cfg)

	mux := http.NewServeMux()

	// Well-known endpoints
	mux.HandleFunc("/.well-known/oauth-authorization-server", wh.AuthorizationServerMetadataHandler(cfg))
	mux.HandleFunc("/.well-known/jwks.json", wh.JWKSHandler(cfg))

	// OAuth endpoints
	mux.HandleFunc("/oauth2/authorize", oh.AuthorizeHandler(store, validator))
	mux.HandleFunc("/oauth2/token", oh.TokenHandler(store, issuer, validator))

	srv := &http.Server{Addr: ":8080", Handler: logging(mux)}
	log.Println("AS PoC listening on", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
