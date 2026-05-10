package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*
	type msgResponse struct {
		Message string `json:"message"`
	}
*/
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	respondWithJSON(w, code, map[string]string{"error": msg})
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
	w.Write(dat)
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
