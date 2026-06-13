package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type msgResponse struct {
	Message string `json:"message"`
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	respondWithJSON(w, code, map[string]string{"error": msg})
}
func respondWithAuthError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		respondWithError(w, http.StatusNotFound, "Resource not found", nil)
		return
	}
	if errors.Is(err, ErrForbidden) {
		respondWithError(w, http.StatusForbidden, "Access denied", nil)
		return
	}
	respondWithError(w, http.StatusInternalServerError, "Authorization failed", err)
}
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("response write failed: %v", err)
	}
}

func sanitizeEmail(email string) string {
	emailParts, err := checkEmail(email)
	if err != nil {
		return "invalid email"
	}
	username := emailParts[0]
	domain := emailParts[1]
	if len(username) > 2 {
		return username[:2] + "****@" + domain
	}
	return username + "****@" + domain
}

func checkEmail(email string) ([]string, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid email format: missing @ symbol")
	}
	return parts, nil
}

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullInt32ToPtr(v sql.NullInt32) *int32 {
	if v.Valid {
		return &v.Int32
	}
	return nil
}

func nullFloat64ToPtr(v sql.NullFloat64) *float64 {
	if v.Valid {
		return &v.Float64
	}
	return nil
}

func nullUUIDToUUID(v uuid.NullUUID) *uuid.UUID {
	if v.Valid {
		return &v.UUID
	}
	return nil
}
