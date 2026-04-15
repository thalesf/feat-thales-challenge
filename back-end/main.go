package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Load the reviews data into memory
	fmt.Println("Starting to load reviews...")
	reviewsData, err := loadReviews()
	if err != nil {
		log.Fatalf("Failed to load reviews: %v", err)
	}
	fmt.Println("Successfully loaded reviews structure")

	// Start the HTTP server
	fmt.Println("Starting server on port 8080...")
	server := NewServer(reviewsData)
	log.Fatal(http.ListenAndServe(":8080", server.Router))
}
