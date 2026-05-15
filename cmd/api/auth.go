package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/auth"
	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRegister(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	hashed_pw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}
	dbparams := database.CreateUserParams{
		ID:           uuid.New(),
		Email:        params.Email,
		PasswordHash: hashed_pw,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	use, err := cfg.db.CreateUser(r.Context(), dbparams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error storing user in database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user{
		UserID:    use.ID,
		Email:     use.Email,
		CreatedAt: use.CreatedAt,
		UpdatedAt: use.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	use, err := cfg.db.GetUserFromEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user", err)
		return
	}
	isValid, err := auth.CheckPasswordHash(params.Password, use.PasswordHash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Error", nil)
		return
	}
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}
	rawToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}
	jwtString, err := auth.MakeJWT(use.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT string", err)
		return
	}
	hashedToken := auth.HashToken(rawToken)
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		HashedToken: hashedToken,
		UserID:      use.ID,
		ExpiresAt:   time.Now().UTC().Add(15 * 24 * time.Hour),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, user{
		UserID:       use.ID,
		CreatedAt:    use.CreatedAt,
		UpdatedAt:    use.UpdatedAt,
		Email:        use.Email,
		Token:        jwtString,
		RefreshToken: rawToken,
	})

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
