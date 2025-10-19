package handler

import (
	"encoding/json"
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

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MetricsResponse struct {
	TopDomains []service.DomainMetric `json:"top_domains"`
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics, err := h.service.GetTopDomains(3)
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