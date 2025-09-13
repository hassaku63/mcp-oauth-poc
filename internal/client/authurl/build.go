package authurl

import (
	"net/url"

	"mcp-oauth-poc/internal/client/store"
)

func Build(authorizeEndpoint string, s store.Session) (string, error) {
	u, err := url.Parse(authorizeEndpoint)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", s.ClientID)
	q.Set("redirect_uri", s.RedirectURI)
	if s.Scope != "" {
		q.Set("scope", s.Scope)
	}
	q.Set("state", s.State)
	q.Set("code_challenge", s.CodeChallenge)
	q.Set("code_challenge_method", "S256")
	if s.Resource != "" {
		q.Set("resource", s.Resource)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
