package handlers

import(
	"context"
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"encoding/hex"
	"strings"
	"crypto/rand"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

func (cfg *apiConfig)handlerRegister(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Password string `json:"password"`
		Email	string	`json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	hashed_pw, err := HashPassword(params.Password)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}
	dbparams := database.CreateUserParams{
		ID: uuid.New(),
		Email: params.Email,
		PasswordHash: hashed_pw,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	user, err := cfg.db.CreateUser(r.Context(), dbparams)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error storing user in database", err )
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Email	string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	user, err := cfg.db.GetUserFromEmail(r.Context(), params.Email)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user", err)
		return
	}
	isValid, err := CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Internal Error", nil)
		return
	}
	if !isValid{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}
	token, err := MakeRefreshToken(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}
	refreshedTokenString := RotateRefreshToken()
	refresh_token, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		HashedToken: refreshedTokenString,
		UserID:	user.ID,
		ExpiresAt:	time.Now().UTC().Add(15*24*time.Hour),
		CreatedAt:	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
	})
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		UserID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		HashedToken: refresh_token.HashedToken,
	})

}
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
/* TOKEN ROTATION LOGIC
func RotateRefreshToken()(string, err){
hashedToken := hashToken(rawToken)
token, err := database.GetRefreshToken(ctx, hashedToken)
if err != nil { // not found
    return unauthorized
}
if token.IsRevoked || token.ExpiresAt.Before(time.Now()) {
    return unauthorized
}
db.RevokeRefreshToken(ctx, hashedToken)

newJWT := generateJWT(token.UserID)
newRawToken := generateRandomToken()
newHashedToken := hashToken(newRawToken)
database.CreateRefreshToken(ctx, newHashedToken, token.UserID, ...)
//return
}
*/
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing authorization token")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("authorization header must be in format 'Bearer {token}'")
	}

	return strings.TrimSpace(parts[1]), nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil{
		return "", err
	}
	string_token := hex.EncodeToString(token)

	return string_token, nil
}

func GetAPIKey (headers http.Header)(string, error){
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Missing authorization token")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "ApiKey") {
		return "", fmt.Errorf("authorization header must be in format 'ApiKey {token}'")
	}

	return strings.TrimSpace(parts[1]), nil
}
