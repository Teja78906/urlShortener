
// internal/storage/store_test.go
package storage

import "testing"

func TestSaveAndGetURL(t *testing.T) {
	store := NewInMemoryStore()

	originalURL := "https://example.com"
	shortCode := "abc123"

	err := store.SaveURL(originalURL, shortCode)
	if err != nil {
		t.Fatalf("SaveURL failed: %v", err)
	}

	url, exists := store.GetURL(shortCode)
	if !exists {
		t.Fatalf("URL not found after saving")
	}

	if url != originalURL {
		t.Fatalf("Expected %s, got %s", originalURL, url)
	}
}

func TestURLExists(t *testing.T) {
	store := NewInMemoryStore()

	originalURL := "https://example.com"
	shortCode := "abc123"

	store.SaveURL(originalURL, shortCode)

	retrievedCode, exists := store.URLExists(originalURL)
	if !exists {
		t.Fatalf("URLExists returned false")
	}

	if retrievedCode != shortCode {
		t.Fatalf("Expected %s, got %s", shortCode, retrievedCode)
	}
}

func TestURLNotExists(t *testing.T) {
	store := NewInMemoryStore()

	_, exists := store.URLExists("https://nonexistent.com")
	if exists {
		t.Fatalf("Expected URLExists to return false")
	}
}

func TestGetAllMappings(t *testing.T) {
	store := NewInMemoryStore()

	mappings := []struct {
		original  string
		shortCode string
	}{
		{"https://example1.com", "abc"},
		{"https://example2.com", "def"},
		{"https://example3.com", "ghi"},
	}

	for _, m := range mappings {
		store.SaveURL(m.original, m.shortCode)
	}

	allMappings := store.GetAllMappings()
	if len(allMappings) != 3 {
		t.Fatalf("Expected 3 mappings, got %d", len(allMappings))
	}
}