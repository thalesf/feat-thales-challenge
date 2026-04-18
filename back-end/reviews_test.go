package main

import (
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
			name:    "empty url",
			input:   validHeader + "u1,Alpha,,\"text\"\n",
			wantErr: "missing college",
		},
		{
			name:    "empty file",
			input:   "",
			wantErr: "read header",
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
