package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestServer_HandleGetReviews(t *testing.T) {
	t.Parallel()

	data := &ReviewsData{
		Colleges: []College{
			{Name: "Alpha University", URL: "alpha"},
			{Name: "Beta College", URL: "beta-college"},
		},
		Reviews: map[string][]string{
			"alpha":        {"Great school.", "Would recommend."},
			"beta-college": {"Decent program."},
		},
	}
	srv := NewServer(data)

	type reviewsResponse struct {
		College College  `json:"college"`
		Reviews []string `json:"reviews"`
	}

	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantBody   reviewsResponse
		wantErrSub string
	}{
		{
			name:       "returns college and reviews for known url",
			target:     "/reviews?url=alpha",
			wantStatus: http.StatusOK,
			wantBody: reviewsResponse{
				College: College{Name: "Alpha University", URL: "alpha"},
				Reviews: []string{"Great school.", "Would recommend."},
			},
		},
		{
			name:       "single-review url",
			target:     "/reviews?url=beta-college",
			wantStatus: http.StatusOK,
			wantBody: reviewsResponse{
				College: College{Name: "Beta College", URL: "beta-college"},
				Reviews: []string{"Decent program."},
			},
		},
		{
			name:       "missing url param returns 400",
			target:     "/reviews",
			wantStatus: http.StatusBadRequest,
			wantErrSub: "missing url query param",
		},
		{
			name:       "empty url param returns 400",
			target:     "/reviews?url=",
			wantStatus: http.StatusBadRequest,
			wantErrSub: "missing url query param",
		},
		{
			name:       "unknown url returns 404",
			target:     "/reviews?url=unknown",
			wantStatus: http.StatusNotFound,
			wantErrSub: "college not found",
		},
		{
			name:       "url with surrounding whitespace is trimmed",
			target:     "/reviews?url=%20alpha%20",
			wantStatus: http.StatusOK,
			wantBody: reviewsResponse{
				College: College{Name: "Alpha University", URL: "alpha"},
				Reviews: []string{"Great school.", "Would recommend."},
			},
		},
		{
			name:       "repeated dashes in url are collapsed",
			target:     "/reviews?url=beta----college",
			wantStatus: http.StatusOK,
			wantBody: reviewsResponse{
				College: College{Name: "Beta College", URL: "beta-college"},
				Reviews: []string{"Decent program."},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, tc.target, nil)
			rr := httptest.NewRecorder()
			srv.Router.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body = %q", rr.Code, tc.wantStatus, rr.Body.String())
			}

			if got := rr.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "*" {
				t.Errorf("Access-Control-Allow-Origin = %q, want *", got)
			}

			if tc.wantStatus != http.StatusOK {
				if !strings.Contains(rr.Body.String(), tc.wantErrSub) {
					t.Errorf("body = %q, want substring %q", rr.Body.String(), tc.wantErrSub)
				}
				return
			}

			var got reviewsResponse
			if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v; body = %q", err, rr.Body.String())
			}
			if !reflect.DeepEqual(got, tc.wantBody) {
				t.Errorf("body mismatch\n got:  %+v\n want: %+v", got, tc.wantBody)
			}
		})
	}
}

func TestServer_HandleGetReviews_RejectsNonGET(t *testing.T) {
	t.Parallel()

	data := &ReviewsData{
		Colleges: []College{{Name: "Alpha", URL: "alpha"}},
		Reviews:  map[string][]string{"alpha": {"r"}},
	}
	srv := NewServer(data)

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(method, "/reviews?url=alpha", nil)
			rr := httptest.NewRecorder()
			srv.Router.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}
