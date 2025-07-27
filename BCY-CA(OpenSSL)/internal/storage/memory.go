package storage

import "sync"

type User struct {
	ID     string
	Secret string
}

type MemoryStore struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{users: make(map[string]User)}
}

func (s *MemoryStore) Register(id, secret string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[id] = User{ID: id, Secret: secret}
}

func (s *MemoryStore) Verify(id, secret string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	return ok && u.Secret == secret
}
