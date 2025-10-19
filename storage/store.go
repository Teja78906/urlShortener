package storage

import "sync"

type URLMapping struct {
	OriginalURL string
	ShortCode   string
}

type Store interface {
	SaveURL(originalURL, shortCode string) error
	GetURL(shortCode string) (string, bool)
	URLExists(originalURL string) (string, bool)
	GetAllMappings() []URLMapping
}

type InMemoryStore struct {
	urlToShort map[string]string // original URL -> short code
	shortToURL map[string]string // short code -> original URL
	mu         sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		urlToShort: make(map[string]string),
		shortToURL: make(map[string]string),
	}
}

func (s *InMemoryStore) SaveURL(originalURL, shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urlToShort[originalURL] = shortCode
	s.shortToURL[shortCode] = originalURL
	return nil
}

func (s *InMemoryStore) GetURL(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.shortToURL[shortCode]
	return url, exists
}

func (s *InMemoryStore) URLExists(originalURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	shortCode, exists := s.urlToShort[originalURL]
	return shortCode, exists
}

func (s *InMemoryStore) GetAllMappings() []URLMapping {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mappings := make([]URLMapping, 0, len(s.urlToShort))
	for originalURL, shortCode := range s.urlToShort {
		mappings = append(mappings, URLMapping{
			OriginalURL: originalURL,
			ShortCode:   shortCode,
		})
	}
	return mappings
}