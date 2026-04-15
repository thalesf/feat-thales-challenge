package main

import (
	"net/http"
)

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

	// Register routes
	server.Router.HandleFunc("/autocomplete", server.handleAutocomplete)
	server.Router.HandleFunc("/reviews", server.handleGetReviews)

	return server
}

// handleAutocomplete handles the autocomplete endpoint
func (s *Server) handleAutocomplete(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// TODO: Implement the autocomplete functionality
	// This should search through the colleges in ReviewsData
	// and return matches based on the query string

	// Write a simple 200 OK with empty JSON response for now
	w.Write([]byte("{}"))
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
