package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// AuthorizeHandler handles GET /oauth2/authorize for code + PKCE.
func AuthorizeHandler(store *Store, v *Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		q := r.URL.Query()
		responseType := q.Get("response_type")
		clientID := q.Get("client_id")
		redirectURI := q.Get("redirect_uri")
		state := q.Get("state")
		codeChallenge := q.Get("code_challenge")
		codeChallengeMethod := q.Get("code_challenge_method")
		scope := q.Get("scope")
		resource := q.Get("resource")

		if responseType != "code" || clientID == "" || redirectURI == "" || state == "" || codeChallenge == "" || codeChallengeMethod == "" {
			writeOAuthErrorRedirect(w, redirectURI, "invalid_request", "missing required parameter", state)
			return
		}

		client := FindClient(clientID)
		if client == nil {
			writeOAuthErrorRedirect(w, redirectURI, "unauthorized_client", "unknown client", state)
			return
		}
		if !v.ValidateRedirectURI(client, redirectURI) {
			writeOAuthErrorRedirect(w, redirectURI, "invalid_request", "invalid redirect_uri", state)
			return
		}
		// Auto-login/consent for PoC; in real, render pages
		code, err := store.IssueCode(CodeRecord{
			ClientID:            clientID,
			RedirectURI:         redirectURI,
			Scope:               scope,
			Resource:            resource,
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: codeChallengeMethod,
		})
		if err != nil {
			writeOAuthErrorRedirect(w, redirectURI, "server_error", "failed to issue code", state)
			return
		}
		// 302 redirect
		u, _ := url.Parse(redirectURI)
		q2 := u.Query()
		q2.Set("code", code)
		q2.Set("state", state)
		u.RawQuery = q2.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
	}
}

func writeOAuthErrorRedirect(w http.ResponseWriter, redirectURI, errCode, desc, state string) {
	if redirectURI == "" {
		http.Error(w, fmt.Sprintf("%s: %s", errCode, desc), http.StatusBadRequest)
		return
	}
	u, pErr := url.Parse(redirectURI)
	if pErr != nil {
		http.Error(w, fmt.Sprintf("%s: %s", errCode, desc), http.StatusBadRequest)
		return
	}
	q := u.Query()
	q.Set("error", errCode)
	if desc != "" {
		q.Set("error_description", desc)
	}
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	http.Redirect(w, &http.Request{}, u.String(), http.StatusFound)
}

// TokenHandler handles POST /oauth2/token for authorization_code.
func TokenHandler(store *Store, issuer *Issuer, v *Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			writeOAuthError(w, http.StatusBadRequest, "invalid_request", "content-type must be application/x-www-form-urlencoded")
			return
		}
		if err := r.ParseForm(); err != nil {
			writeOAuthError(w, http.StatusBadRequest, "invalid_request", "invalid form")
			return
		}
		grantType := r.PostForm.Get("grant_type")
		switch grantType {
		case "authorization_code":
			handleAuthCodeGrant(w, r, store, issuer, v)
		default:
			writeOAuthError(w, http.StatusBadRequest, "unsupported_grant_type", "only authorization_code is supported")
		}
	}
}

func handleAuthCodeGrant(w http.ResponseWriter, r *http.Request, store *Store, issuer *Issuer, v *Validator) {
	code := r.PostForm.Get("code")
	redirectURI := r.PostForm.Get("redirect_uri")
	clientID := r.PostForm.Get("client_id")
	codeVerifier := r.PostForm.Get("code_verifier")
	resource := r.PostForm.Get("resource")

	if code == "" || redirectURI == "" || clientID == "" || codeVerifier == "" {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "missing required parameter")
		return
	}
	client := FindClient(clientID)
	if client == nil {
		writeOAuthError(w, http.StatusBadRequest, "unauthorized_client", "unknown client")
		return
	}
	if !v.ValidateRedirectURI(client, redirectURI) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "invalid redirect_uri")
		return
	}
	rec, ok := store.ConsumeCode(code)
	if !ok {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "invalid or expired code")
		return
	}
	// Bindings check
	if rec.ClientID != clientID || rec.RedirectURI != redirectURI {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "mismatched binding")
		return
	}
	if !v.ValidatePKCE(codeVerifier, rec.CodeChallenge, rec.CodeChallengeMethod) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "pkce verification failed")
		return
	}
    aud := rec.Resource
    if aud == "" {
        aud = resource
    }
    if aud == "" {
        aud = store.cfg.ResourceID
    }
	tok, err := issuer.MintAccessToken("user-123", aud, rec.Scope)
	if err != nil {
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "failed to mint token")
		return
	}
	resp := map[string]any{
		"access_token": tok,
		"token_type":   "Bearer",
		"expires_in":   int(store.cfg.TokenTTL.Seconds()),
		"scope":        rec.Scope,
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeOAuthError(w http.ResponseWriter, status int, code, desc string) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":             code,
		"error_description": desc,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
