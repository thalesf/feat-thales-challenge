package main

import (
	"encoding/json"
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
	query := r.URL.Query().Get("q")

	limit := defaultAutocompleteLimit
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			http.Error(w, `{"error":"limit must be a positive integer"}`, http.StatusBadRequest)
			return
		}
		limit = n
	}

	results := s.ReviewsData.Autocomplete(query, limit)
	if results == nil {
		results = []College{}
	}

	json.NewEncoder(w).Encode(results)
}

// handleGetReviews handles the reviews endpoint
func (s *Server) handleGetReviews(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// TODO: Implement the reviews endpoint
	// This should retrieve reviews for a specific college
	// and return them in the response

	// Write a simple 200 OK with empty JSON response for now
	w.Write([]byte("{}"))
}
