package token

import (
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

type ExchangeRequest struct {
    Code         string
    RedirectURI  string
    ClientID     string
    CodeVerifier string
    Resource     string
}

type ExchangeResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
    Scope       string `json:"scope"`
}

func Exchange(tokenEndpoint string, in ExchangeRequest) (ExchangeResponse, error) {
    form := url.Values{}
    form.Set("grant_type", "authorization_code")
    form.Set("code", in.Code)
    form.Set("redirect_uri", in.RedirectURI)
    form.Set("client_id", in.ClientID)
    form.Set("code_verifier", in.CodeVerifier)
    if in.Resource != "" { form.Set("resource", in.Resource) }

    req, _ := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(form.Encode()))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Accept", "application/json")

    hc := &http.Client{ Timeout: 10 * time.Second }
    resp, err := hc.Do(req)
    if err != nil {
        return ExchangeResponse{}, err
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    if resp.StatusCode != http.StatusOK {
        return ExchangeResponse{}, errors.New(resp.Status + ": " + string(body))
    }
    var out ExchangeResponse
    if err := json.Unmarshal(body, &out); err != nil {
        return ExchangeResponse{}, err
    }
    return out, nil
}

