package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const validHeader = "COLLEGE_UUID,COLLEGE_NAME,COLLEGE_URL,REVIEW_TEXT\n"

func writeCSV(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "reviews.csv")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func TestLoadReviewsFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr string // substring match; empty means no error expected
		want    *ReviewsData
	}{
		{
			name: "happy path with dedup",
			input: validHeader +
				`u1,Alpha University,alpha,"Great school."` + "\n" +
				`u2,Beta College,beta,"Loved it."` + "\n" +
				`u2,Beta College,beta,"Would recommend."` + "\n",
			want: &ReviewsData{
				Colleges: []College{
					{Name: "Alpha University", URL: "alpha"},
					{Name: "Beta College", URL: "beta"},
				},
				Reviews: map[string][]string{
					"alpha": {"Great school."},
					"beta":  {"Loved it.", "Would recommend."},
				},
			},
		},
		{
			name: "embedded comma and escaped quote",
			input: validHeader +
				`u1,Alpha,alpha,"Great, but ""expensive"" classes."` + "\n",
			want: &ReviewsData{
				Colleges: []College{{Name: "Alpha", URL: "alpha"}},
				Reviews:  map[string][]string{"alpha": {`Great, but "expensive" classes.`}},
			},
		},
		{
			name:  "multi-line quoted review",
			input: validHeader + "u1,Alpha,alpha,\"line one\nline two\"\n",
			want: &ReviewsData{
				Colleges: []College{{Name: "Alpha", URL: "alpha"}},
				Reviews:  map[string][]string{"alpha": {"line one\nline two"}},
			},
		},
		{
			name:  "valid header with no rows",
			input: validHeader,
			want: &ReviewsData{
				Colleges: nil,
				Reviews:  map[string][]string{},
			},
		},
		{
			name:    "wrong header",
			input:   "A,B,C,D\n",
			wantErr: "invalid header",
		},
		{
			name:    "wrong column count in row",
			input:   validHeader + "u1,Alpha,alpha\n",
			wantErr: "read row",
		},
		{
			name: "rows with missing name or url are skipped",
			input: validHeader +
				`u1,Alpha,alpha,"keep"` + "\n" +
				`u2,NoURL,,"drop"` + "\n" +
				`u3,,no-name,"drop"` + "\n",
			want: &ReviewsData{
				Colleges: []College{{Name: "Alpha", URL: "alpha"}},
				Reviews:  map[string][]string{"alpha": {"keep"}},
			},
		},
		{
			name: "row with both name and url empty is skipped",
			input: validHeader +
				`u1,Alpha,alpha,"keep"` + "\n" +
				`u2,,,"drop"` + "\n",
			want: &ReviewsData{
				Colleges: []College{{Name: "Alpha", URL: "alpha"}},
				Reviews:  map[string][]string{"alpha": {"keep"}},
			},
		},
		{
			name: "rows with whitespace-only name or url are skipped",
			input: validHeader +
				`u1,Alpha,alpha,"keep"` + "\n" +
				`u2,"   ",beta,"drop"` + "\n" +
				`u3,Gamma,"   ","drop"` + "\n",
			want: &ReviewsData{
				Colleges: []College{{Name: "Alpha", URL: "alpha"}},
				Reviews:  map[string][]string{"alpha": {"keep"}},
			},
		},
		{
			name: "skipped row does not pollute reviews or colleges of subsequent valid row",
			input: validHeader +
				`u1,,,"drop"` + "\n" +
				`u2,Alpha,alpha,"keep"` + "\n",
			want: &ReviewsData{
				Colleges: []College{{Name: "Alpha", URL: "alpha"}},
				Reviews:  map[string][]string{"alpha": {"keep"}},
			},
		},
		{
			name:    "empty file",
			input:   "",
			wantErr: "read header",
		},
		{
			name: "urls with repeated dashes are normalized and deduped",
			input: validHeader +
				`u1,Academy of Art University,academy-of-art-university,"first"` + "\n" +
				`u2,Academy of Art University,academy----of-art-university,"second"` + "\n" +
				`u3,Academy of Art University,academy--of-art-university,"third"` + "\n",
			want: &ReviewsData{
				Colleges: []College{{Name: "Academy of Art University", URL: "academy-of-art-university"}},
				Reviews: map[string][]string{
					"academy-of-art-university": {"first", "second", "third"},
				},
			},
		},
		{
			name: "colleges sorted case-insensitive",
			input: validHeader +
				`u1,zeta,zeta,"r"` + "\n" +
				`u2,Alpha,alpha,"r"` + "\n" +
				`u3,beta,beta,"r"` + "\n",
			want: &ReviewsData{
				Colleges: []College{
					{Name: "Alpha", URL: "alpha"},
					{Name: "beta", URL: "beta"},
					{Name: "zeta", URL: "zeta"},
				},
				Reviews: map[string][]string{
					"alpha": {"r"},
					"beta":  {"r"},
					"zeta":  {"r"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := writeCSV(t, tc.input)
			got, err := loadReviewsFile(path)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("want error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want substring %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got.Colleges, tc.want.Colleges) {
				t.Errorf("Colleges mismatch\n got:  %+v\n want: %+v", got.Colleges, tc.want.Colleges)
			}
			if !reflect.DeepEqual(got.Reviews, tc.want.Reviews) {
				t.Errorf("Reviews mismatch\n got:  %+v\n want: %+v", got.Reviews, tc.want.Reviews)
			}
		})
	}
}

func TestLoadReviewsFile_FileMissing(t *testing.T) {
	t.Parallel()

	missing := filepath.Join(t.TempDir(), "nope.csv")
	if _, err := loadReviewsFile(missing); err == nil {
		t.Fatal("want error, got nil")
	}
}

func TestServer_HandleAutocomplete(t *testing.T) {
	t.Parallel()

	data := &ReviewsData{Colleges: []College{
		{Name: "Alpha University", URL: "alpha"},
		{Name: "alpha State", URL: "alpha-state"},
		{Name: "Beta College", URL: "beta"},
		{Name: "Boston Tech", URL: "boston"},
		{Name: "Zeta Institute", URL: "zeta"},
	}}
	srv := NewServer(data)

	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantBody   []College
		wantErrSub string
	}{
		{
			name:       "matches returned as json array",
			target:     "/autocomplete?q=alp",
			wantStatus: http.StatusOK,
			wantBody: []College{
				{Name: "Alpha University", URL: "alpha"},
				{Name: "alpha State", URL: "alpha-state"},
			},
		},
		{
			name:       "case-insensitive match",
			target:     "/autocomplete?q=ALP",
			wantStatus: http.StatusOK,
			wantBody: []College{
				{Name: "Alpha University", URL: "alpha"},
				{Name: "alpha State", URL: "alpha-state"},
			},
		},
		{
			name:       "empty q returns empty array not null",
			target:     "/autocomplete?q=",
			wantStatus: http.StatusOK,
			wantBody:   []College{},
		},
		{
			name:       "missing q param returns empty array",
			target:     "/autocomplete",
			wantStatus: http.StatusOK,
			wantBody:   []College{},
		},
		{
			name:       "no match returns empty array",
			target:     "/autocomplete?q=xyz",
			wantStatus: http.StatusOK,
			wantBody:   []College{},
		},
		{
			name:       "limit caps results",
			target:     "/autocomplete?q=a&limit=1",
			wantStatus: http.StatusOK,
			wantBody: []College{
				{Name: "Alpha University", URL: "alpha"},
			},
		},
		{
			name:       "limit larger than matches returns all",
			target:     "/autocomplete?q=a&limit=99",
			wantStatus: http.StatusOK,
			wantBody: []College{
				{Name: "Alpha University", URL: "alpha"},
				{Name: "alpha State", URL: "alpha-state"},
			},
		},
		{
			name:       "empty limit value falls back to default",
			target:     "/autocomplete?q=a&limit=",
			wantStatus: http.StatusOK,
			wantBody: []College{
				{Name: "Alpha University", URL: "alpha"},
				{Name: "alpha State", URL: "alpha-state"},
			},
		},
		{
			name:       "non-numeric limit is rejected",
			target:     "/autocomplete?q=a&limit=abc",
			wantStatus: http.StatusBadRequest,
			wantErrSub: "limit must be a positive integer",
		},
		{
			name:       "zero limit is rejected",
			target:     "/autocomplete?q=a&limit=0",
			wantStatus: http.StatusBadRequest,
			wantErrSub: "limit must be a positive integer",
		},
		{
			name:       "negative limit is rejected",
			target:     "/autocomplete?q=a&limit=-5",
			wantStatus: http.StatusBadRequest,
			wantErrSub: "limit must be a positive integer",
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

			if tc.wantStatus != http.StatusOK {
				if !strings.Contains(rr.Body.String(), tc.wantErrSub) {
					t.Errorf("body = %q, want substring %q", rr.Body.String(), tc.wantErrSub)
				}
				return
			}

			var got []College
			if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v; body = %q", err, rr.Body.String())
			}
			if got == nil {
				got = []College{}
			}
			if !reflect.DeepEqual(got, tc.wantBody) {
				t.Errorf("body mismatch\n got:  %+v\n want: %+v", got, tc.wantBody)
			}
		})
	}
}
