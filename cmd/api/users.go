package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/auth"
	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("userId")
	var userID uuid.UUID

	if idString == "" {
		var err error
		userID, err = cfg.getUserIDFromToken(w, r)
		if err != nil {
			return
		}
	} else {
		var err error
		userID, err = uuid.Parse(idString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user_id", err)
			return
		}
	}
	user, err := cfg.db.GetUserFromID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, userResponse{
		UserID:    user.ID,
		Email:     sanitizeEmail(user.Email),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerGetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	user, err := cfg.db.GetUserFromID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, userResponse{
		UserID:    user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerChangeEmail(w http.ResponseWriter, r *http.Request) {

	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		Email string `json:"email"`
	}
	var params parameters
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	_, err = checkEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid email", err)
		return
	}
	err = cfg.db.UpdateUserEmail(r.Context(), database.UpdateUserEmailParams{
		Email: params.Email,
		ID:    userID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad Request", err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "email updated"})
}
func (cfg *apiConfig) handlerChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		Password string `json:"password"`
	}
	var params parameters

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash new password", nil)
		return
	}
	err = cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		PasswordHash: hashed,
		ID:           userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to store to database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
}

func (cfg *apiConfig) handlerDeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		Password string `json:"password"`
	}
	var params parameters

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	user, err := cfg.db.GetUserFromID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user from database", err)
		return
	}
	isValid, err := auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error checking password", nil)
		return
	}
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password entered", nil)
		return
	}
	err = cfg.db.DeleteUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user from database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
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

func (cfg *apiConfig) getUserIDFromToken(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing token", err)
		return uuid.Nil, err
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return uuid.Nil, err
	}
	return userID, nil
}
