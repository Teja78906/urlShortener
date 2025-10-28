package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"urlshortener/internal/service"
)

type Handler struct {
	service *service.URLService
}

func NewHandler(service *service.URLService) *Handler {
	return &Handler{service: service}
}

type ShortenRequest struct {
	URL string `json:"url"`
}
type BatchRequest struct {
	URLs []string `json:"urls"`
}
type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
}
type BatchResponse struct {
	Results []ShortenResponse `json:"results"`
	Errors  []string          `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
type MetricRequest struct {
	Limit int `json:"limit"`
}
type MetricsResponse struct {
	TopDomains []service.DomainMetric `json:"top_domains"`
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("request received on /shorten.....")
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if strings.TrimSpace(req.URL) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "URL cannot be empty"})
		return
	}
	if !strings.Contains(req.URL, "://") {
		req.URL = "http://" + req.URL
	}
	shortCode, err := h.service.ShortenURL(req.URL)
	if err != nil {
		log.Printf("Error shortening URL: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ShortenResponse{
		ShortCode: shortCode,
		URL:       req.URL,
	})
}

func (h *Handler) ShortenMultipleURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("request received on /shortenMultipleURLs.....")
	var req BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}
	var responses []ShortenResponse
	var errors []string
	for _, url := range req.URLs {
		if strings.TrimSpace(url) == "" {
			continue
		}
		if !strings.Contains(url, "://") {
			url = "http://" + url
		}
		shortCode, err := h.service.ShortenURL(url)
		if err == nil {
			responses = append(responses, ShortenResponse{
				ShortCode: shortCode,
				URL:       url,
			})
		} else {
			errors = append(errors, fmt.Sprintf("Error shortening URL %s: %v", url, err))

		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BatchResponse{
		Results: responses,
		Errors:  errors,
	})
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/redirect/")

	if strings.TrimSpace(shortCode) == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.GetOriginalURL(shortCode)
	if err != nil {
		log.Printf("Short code not found: %s", shortCode)
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (h *Handler) GetTopDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MetricRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	metrics, err := h.service.GetTopDomains(req.Limit)

	if err != nil {
		log.Printf("Error getting metrics: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MetricsResponse{TopDomains: metrics})
}
