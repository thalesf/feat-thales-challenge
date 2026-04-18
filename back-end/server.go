package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const defaultAutocompleteLimit = 20

// Server represents the HTTP server for our application
type Server struct {
	Router      *http.ServeMux
	ReviewsData *ReviewsData
}

// NewServer creates a new HTTP server with the given reviews data
func NewServer(reviewsData *ReviewsData) *Server {
	server := &Server{
		Router:      http.NewServeMux(),
		ReviewsData: reviewsData,
	}

	server.Router.HandleFunc("GET /autocomplete", server.handleAutocomplete)
	server.Router.HandleFunc("GET /reviews", server.handleGetReviews)

	return server
}

func (s *Server) handleAutocomplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := r.URL.Query().Get("q")

	limit := defaultAutocompleteLimit
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"limit must be a positive integer"}`))
			return
		}
		limit = n
	}

	results := s.ReviewsData.Autocomplete(query, limit)
	if results == nil {
		results = []College{}
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		log.Printf("autocomplete: encode response: %v", err)
	}
}

func (s *Server) handleGetReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"missing url query param"}`))
		return
	}

	reviews, college, ok := s.ReviewsData.ReviewsForURL(url)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"college not found"}`))
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]any{
		"college": college,
		"reviews": reviews,
	}); err != nil {
		log.Printf("reviews: encode response: %v", err)
	}
}
