package main

import (
	"log"
	"net/http"

	"urlshortener/internal/handler"
	"urlshortener/internal/service"
	"urlshortener/internal/storage"

)

func main() {
	// Initialize storage
	store := storage.NewInMemoryStore()

	// Initialize service
	svc := service.NewURLService(store)

	// Initialize handler
	h := handler.NewHandler(svc)

	// Setup routes
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", h.ShortenURL)
	mux.HandleFunc("/redirect/", h.RedirectURL)
	mux.HandleFunc("/metrics/top-domains", h.GetTopDomains)
	mux.HandleFunc("/shortenMultipleURLs", h.ShortenMultipleURLs)
	
	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}