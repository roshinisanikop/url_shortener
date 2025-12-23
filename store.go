package main

import (
	"errors"
	"sync"
	"time"
)

// URLMapping represents a shortened URL mapping
type URLMapping struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	Clicks      int       `json:"clicks"`
}

// URLStore manages URL mappings
type URLStore struct {
	mu       sync.RWMutex
	urls     map[string]*URLMapping
	reverse  map[string]string // original URL -> short code for deduplication
}

// NewURLStore creates a new URL store
func NewURLStore() *URLStore {
	return &URLStore{
		urls:    make(map[string]*URLMapping),
		reverse: make(map[string]string),
	}
}

// Save stores a new URL mapping
func (s *URLStore) Save(shortCode, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.urls[shortCode]; exists {
		return errors.New("short code already exists")
	}

	mapping := &URLMapping{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		Clicks:      0,
	}

	s.urls[shortCode] = mapping
	s.reverse[originalURL] = shortCode

	return nil
}

// Get retrieves the original URL for a short code
func (s *URLStore) Get(shortCode string) (*URLMapping, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mapping, exists := s.urls[shortCode]
	if !exists {
		return nil, errors.New("short code not found")
	}

	return mapping, nil
}

// IncrementClicks increments the click counter for a short code
func (s *URLStore) IncrementClicks(shortCode string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if mapping, exists := s.urls[shortCode]; exists {
		mapping.Clicks++
	}
}

// GetByOriginalURL retrieves the short code for an original URL
func (s *URLStore) GetByOriginalURL(originalURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	shortCode, exists := s.reverse[originalURL]
	return shortCode, exists
}

// GetAll returns all URL mappings
func (s *URLStore) GetAll() []*URLMapping {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mappings := make([]*URLMapping, 0, len(s.urls))
	for _, mapping := range s.urls {
		mappings = append(mappings, mapping)
	}

	return mappings
}

// Exists checks if a short code exists
func (s *URLStore) Exists(shortCode string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.urls[shortCode]
	return exists
}
