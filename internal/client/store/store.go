package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Session struct {
	ID            string    `json:"session_id"`
	ClientID      string    `json:"client_id"`
	RedirectURI   string    `json:"redirect_uri"`
	Resource      string    `json:"resource"`
	Scope         string    `json:"scope"`
	CodeVerifier  string    `json:"code_verifier"`
	CodeChallenge string    `json:"code_challenge"`
	State         string    `json:"state"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresInSec  int       `json:"expires_in"`
}

func dir() (string, error) {
	d := ".auth-cli-sessions"
	if err := os.MkdirAll(d, 0o700); err != nil {
		return "", err
	}
	return d, nil
}

func pathFor(id string) (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, id+".json"), nil
}

func Save(s Session) error {
	p, err := pathFor(s.ID)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

func Load(id string) (Session, error) {
	p, err := pathFor(id)
	if err != nil {
		return Session{}, err
	}
	f, err := os.Open(p)
	if err != nil {
		return Session{}, err
	}
	defer f.Close()
	var s Session
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return Session{}, err
	}
	// TTL check (best-effort)
	if s.ExpiresInSec > 0 && time.Since(s.CreatedAt) > time.Duration(s.ExpiresInSec)*time.Second {
		return Session{}, errors.New("session expired")
	}
	return s, nil
}

func Delete(id string) error {
	p, err := pathFor(id)
	if err != nil {
		return err
	}
	_ = os.Remove(p)
	return nil
}
