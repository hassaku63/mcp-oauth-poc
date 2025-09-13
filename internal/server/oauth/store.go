package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type Store struct {
	cfg   Config
	mu    sync.Mutex
	codes map[string]CodeRecord
}

func NewStore(cfg Config) *Store {
	s := &Store{cfg: cfg, codes: make(map[string]CodeRecord)}
	go s.gc()
	return s
}

func (s *Store) gc() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, v := range s.codes {
			if v.ExpiresAt.Before(now) || v.Used {
				delete(s.codes, k)
			}
		}
		s.mu.Unlock()
	}
}

func randString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *Store) IssueCode(rec CodeRecord) (string, error) {
	code, err := randString(32)
	if err != nil {
		return "", err
	}
	rec.ExpiresAt = time.Now().Add(s.cfg.CodeTTL)
	s.mu.Lock()
	s.codes[code] = rec
	s.mu.Unlock()
	return code, nil
}

func (s *Store) ConsumeCode(code string) (CodeRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.codes[code]
	if !ok || rec.Used || rec.ExpiresAt.Before(time.Now()) {
		return CodeRecord{}, false
	}
	rec.Used = true
	s.codes[code] = rec
	return rec, true
}
