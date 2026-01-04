package auth

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

type TokenStore struct {
	mu     sync.Mutex
	tokens map[string]string // token -> email
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]string),
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

func (s *TokenStore) Add(token, email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = email
}

func (s *TokenStore) Use(token string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	email, ok := s.tokens[token]
	if ok {
		delete(s.tokens, token)
	}
	return email, ok
}
