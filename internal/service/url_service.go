package service

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"urlshortener/internal/storage"
)

type URLService struct {
	store storage.Store
}

func NewURLService(store storage.Store) *URLService {
	return &URLService{store: store}
}

func (s *URLService) ShortenURL(originalURL string) (string, error) {
	if !strings.Contains(originalURL, "://") {
        originalURL = "http://" + originalURL
    }
	// Validate URL
	if err := validateURL(originalURL); err != nil {
		return "", err
	}

	// Check if URL already exists
	if shortCode, exists := s.store.URLExists(originalURL); exists {
		return shortCode, nil
	}

	// Generate new short code
	shortCode := s.generateShortCode(originalURL)

	// Save to store
	if err := s.store.SaveURL(originalURL, shortCode); err != nil {
		return "", err
	}

	return shortCode, nil
}

func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	originalURL, exists := s.store.GetURL(shortCode)
	if !exists {
		return "", fmt.Errorf("short code not found")
	}
	return originalURL, nil
}

func (s *URLService) GetTopDomains(limit int) ([]DomainMetric, error) {
	mappings := s.store.GetAllMappings()
	domainCounts := make(map[string]int)

	for _, mapping := range mappings {
		domain, err := extractDomain(mapping.OriginalURL)
		if err == nil {
			domainCounts[domain]++
		}
	}

	metrics := make([]DomainMetric, 0)
	for domain, count := range domainCounts {
		metrics = append(metrics, DomainMetric{
			Domain: domain,
			Count:  count,
		})
	}

	// Sort by count descending
	sortMetrics(metrics)

	// Limit results
	if len(metrics) > limit {
		metrics = metrics[:limit]
	}

	return metrics, nil
}

func (s *URLService) generateShortCode(originalURL string) string {
	hash := md5.Sum([]byte(originalURL))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	// Return first 7 characters for a reasonable short code
	return strings.TrimRight(encoded[:7], "=")
}

func validateURL(urlString string) error {
	
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("URL must have scheme and host")
	}

	return nil
}

func extractDomain(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host
	// Remove www. prefix if present
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}

	// Extract just the domain without port
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	return host, nil
}

func sortMetrics(metrics []DomainMetric) {
	sort.Slice(metrics, func(i,j int) bool{
		return metrics[i].Count > metrics[j].Count
	})
}

type DomainMetric struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}