package auth

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type TokenStore struct {
	mu     sync.Mutex
	tokens map[string]tokenEntry // token -> email
}

type tokenEntry struct {
	email     string
	expiresAt time.Time
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]tokenEntry),
	}
}

func GenerateToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *TokenStore) Add(token, email string, ttl time.Duration) { // ttl = "time to live"
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = tokenEntry{
		email:     email,
		expiresAt: time.Now().Add(ttl),
	}
}

func (s *TokenStore) Use(token string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.tokens[token]
	if !ok {
		return "", false
	}

	if time.Now().After(entry.expiresAt) {
		delete(s.tokens, token)
		return "", false
	}

	delete(s.tokens, token)
	return entry.email, true
}

func (s *TokenStore) CleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, entry := range s.tokens {
		if now.After(entry.expiresAt) {
			delete(s.tokens, token)
		}
	}
}

func (s *TokenStore) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			s.CleanupExpired()
		}
	}()
}
