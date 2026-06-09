package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errorResponse struct {
	Error string `json:"error"`
}

type payloadResponse struct {
	Message string `json:"message"`
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithError(w, http.StatusBadRequest, "something went wrong", nil)

	result := w.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", result.StatusCode)
	}

	var body errorResponse

	decodeResponseBody(t, result, &body)

	if body.Error != "something went wrong" {
		t.Fatalf("expected error message 'something went wrong', got %q", body.Error)
	}
}

func TestRespondWithAuthError_NotFound(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithAuthError(w, ErrNotFound)

	result := w.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", result.StatusCode)
	}
	var body errorResponse
	decodeResponseBody(t, result, &body)

	if body.Error != "Resource not found" {
		t.Fatalf("expected error message 'Resource not found', got %q", body.Error)
	}
}

func TestRespondWithAuthError_Forbidden(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithAuthError(w, ErrForbidden)

	result := w.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", result.StatusCode)
	}
	var body errorResponse

	decodeResponseBody(t, result, &body)

	if body.Error != "Access denied" {
		t.Fatalf("expected error message 'Access denied', got %q", body.Error)
	}
}
func TestRespondWithAuthError_Unknown(t *testing.T) {
	w := httptest.NewRecorder()
	respondWithAuthError(w, fmt.Errorf("some unexpected error"))

	result := w.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", result.StatusCode)
	}
	var body errorResponse

	decodeResponseBody(t, result, &body)

	if body.Error != "Authorization failed" {
		t.Fatalf("expected error message 'Authorization failed', got %q", body.Error)
	}

}

func TestRespondWithJSON(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithJSON(w, http.StatusOK, msgResponse{"working"})

	result := w.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", result.StatusCode)
	}
	var body payloadResponse

	decodeResponseBody(t, result, &body)

	if body.Message != "working" {
		t.Fatalf("expected message 'working', got %q", body.Message)
	}
}
func TestRespondWithJSON_ContentType(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithJSON(w, http.StatusOK, msgResponse{"working"})

	result := w.Result()
	defer result.Body.Close()

	contentType := result.Header.Get("Content-Type")

	if contentType != "application/json" {
		t.Fatalf("expected content type application/json, got %q", contentType)
	}
}

/*
 *
 * Helpers
 *
 */
func decodeResponseBody(t *testing.T, result *http.Response, v any) {
	t.Helper()

	if err := json.NewDecoder(result.Body).Decode(v); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}
