package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithError(w, http.StatusBadRequest, "something went wrong", nil)

	result := w.Result()
	if result.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", result.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(result.Body).Decode(&body)

	if body["error"] != "something went wrong" {
		t.Fatalf("expected error message 'something went wrong', got %q", body["error"])
	}
}

func TestRespondWithAuthError(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithAuthError(w, ErrNotFound)

	result := w.Result()
	if result.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", result.StatusCode)
	}
	var body map[string]string
	json.NewDecoder(result.Body).Decode(&body)

	if body["error"] != "Resource not found" {
		t.Fatalf("expected error message 'Resource not found', got %q", body["error"])
	}
}
