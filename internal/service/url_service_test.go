// internal/service/url_service_test.go
package service

import (
	"testing"

	"urlshortener/internal/storage"
)

func TestShortenURL_ValidURL(t *testing.T) {
	store := storage.NewInMemoryStore()
	svc := NewURLService(store)

	url := "https://www.example.com/path"
	shortCode, err := svc.ShortenURL(url)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if shortCode == "" {
		t.Fatalf("Expected non-empty short code")
	}

	if len(shortCode) != 7 {
		t.Fatalf("Expected short code length 7, got %d", len(shortCode))
	}
}

func TestShortenURL_InvalidURL(t *testing.T) {
	store := storage.NewInMemoryStore()
	svc := NewURLService(store)

	testCases := []struct {
		url string
		name string
	}{
		{"invalid-url", "missing scheme"},
		{"http://", "missing host"},
		{"", "empty URL"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.ShortenURL(tc.url)
			if err == nil {
				t.Fatalf("Expected error for %s", tc.name)
			}
		})
	}
}

func TestShortenURL_DuplicateURL(t *testing.T) {
	store := storage.NewInMemoryStore()
	svc := NewURLService(store)

	url := "https://www.example.com"
	shortCode1, err := svc.ShortenURL(url)
	if err != nil {
		t.Fatalf("First shorten failed: %v", err)
	}

	shortCode2, err := svc.ShortenURL(url)
	if err != nil {
		t.Fatalf("Second shorten failed: %v", err)
	}

	if shortCode1 != shortCode2 {
		t.Fatalf("Expected same short code for duplicate URL, got %s and %s", shortCode1, shortCode2)
	}
}

func TestGetOriginalURL(t *testing.T) {
	store := storage.NewInMemoryStore()
	svc := NewURLService(store)

	url := "https://www.youtube.com/watch?v=123"
	shortCode, _ := svc.ShortenURL(url)

	retrieved, err := svc.GetOriginalURL(shortCode)
	if err != nil {
		t.Fatalf("Error retrieving URL: %v", err)
	}

	if retrieved != url {
		t.Fatalf("Expected %s, got %s", url, retrieved)
	}
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	store := storage.NewInMemoryStore()
	svc := NewURLService(store)

	_, err := svc.GetOriginalURL("nonexistent")
	if err == nil {
		t.Fatalf("Expected error for non-existent short code")
	}
}

func TestGetTopDomains(t *testing.T) {
	store := storage.NewInMemoryStore()
	svc := NewURLService(store)

	urls := []string{
		"https://www.youtube.com/video1",
		"https://www.youtube.com/video2",
		"https://www.youtube.com/video3",
		"https://www.youtube.com/video4",
		"https://stackoverflow.com/question1",
		"https://wikipedia.org/page1",
		"https://wikipedia.org/page2",
		"https://udemy.com/course1",
		"https://udemy.com/course2",
		"https://udemy.com/course3",
		"https://udemy.com/course4",
		"https://udemy.com/course5",
		"https://udemy.com/course6",
	}

	for _, url := range urls {
		svc.ShortenURL(url)
	}

	metrics, err := svc.GetTopDomains(3)
	if err != nil {
		t.Fatalf("Error getting metrics: %v", err)
	}

	expectedOrder := []struct {
		domain string
		count  int
	}{
		{"udemy.com", 6},
		{"youtube.com", 4},
		{"wikipedia.org", 2},
	}

	if len(metrics) != 3 {
		t.Fatalf("Expected 3 metrics, got %d", len(metrics))
	}

	for i, expected := range expectedOrder {
		if metrics[i].Domain != expected.domain || metrics[i].Count != expected.count {
			t.Fatalf("At index %d: expected %s with count %d, got %s with count %d",
				i, expected.domain, expected.count, metrics[i].Domain, metrics[i].Count)
		}
	}
}

func TestExtractDomain(t *testing.T) {
	testCases := []struct {
		url      string
		expected string
		name     string
	}{
		{"https://www.example.com/path", "example.com", "with www"},
		{"https://github.com/repo", "github.com", "without www"},
		{"https://localhost:8080/path", "localhost", "with port"},
		{"https://sub.example.com", "sub.example.com", "subdomain"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			domain, err := extractDomain(tc.url)
			if err != nil {
				t.Fatalf("Error extracting domain: %v", err)
			}
			if domain != tc.expected {
				t.Fatalf("Expected %s, got %s", tc.expected, domain)
			}
		})
	}
}
