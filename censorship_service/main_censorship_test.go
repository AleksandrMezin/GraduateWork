// main_censorship_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCensorHandler(t *testing.T) {
	testCases := []struct {
		name       string
		comment    Comment
		wantStatus int
	}{
		{
			name:       "allowed comment",
			comment:    Comment{Text: "hello world"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "comment with forbidden word",
			comment:    Comment{Text: "это слово пизда запрещено"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.comment)
			req, _ := http.NewRequest(http.MethodPost, "/censor", bytes.NewBuffer(reqBody))
			rec := httptest.NewRecorder()

			censorHandler(rec, req)

			res := rec.Result()
			if res.StatusCode != tc.wantStatus {
				t.Errorf("unexpected status code: got %v, want %v", res.StatusCode, tc.wantStatus)
			}
		})
	}
}
