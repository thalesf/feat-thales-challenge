package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type College struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ReviewsData struct {
	Reviews  map[string][]string
	Colleges []College
}

const reviewsCSVPath = "data/niche_reviews.csv"

var expectedHeaders = []string{
	"COLLEGE_UUID",
	"COLLEGE_NAME",
	"COLLEGE_URL",
	"REVIEW_TEXT",
}

func loadReviews() (*ReviewsData, error) {
	return loadReviewsFile(reviewsCSVPath)
}

func loadReviewsFile(path string) (data *ReviewsData, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close reviews csv: %w", closeErr))
		}
	}()

	r := csv.NewReader(file)

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	for i, expected := range expectedHeaders {
		if i >= len(header) || strings.TrimSpace(header[i]) != expected {
			return nil, fmt.Errorf("invalid header: column %d want %q", i, expected)
		}
	}

	data = &ReviewsData{Reviews: map[string][]string{}}
	seen := map[string]bool{}

	for {
		row, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}

		name := strings.TrimSpace(row[1])
		url := strings.TrimSpace(row[2])
		if name == "" || url == "" {
			log.Printf("reviews: skipping row with missing college name or url (uuid=%q)", row[0])
			continue
		}

		data.Reviews[url] = append(data.Reviews[url], row[3])
		if !seen[url] {
			seen[url] = true
			data.Colleges = append(data.Colleges, College{Name: name, URL: url})
		}
	}

	sort.Slice(data.Colleges, func(i, j int) bool {
		return strings.ToLower(data.Colleges[i].Name) < strings.ToLower(data.Colleges[j].Name)
	})

	return data, nil
}
