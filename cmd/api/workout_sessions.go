package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		WorkoutDate time.Time `json:"workout_date"`
		Description string    `json:"description"`
		Notes       string    `json:"notes"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	workoutSession, err := cfg.db.CreateWorkoutSessions(r.Context(), database.CreateWorkoutSessionsParams{
		UserID:      userID,
		WorkoutDate: params.WorkoutDate,
		Description: sql.NullString{
			String: params.Description,
			Valid:  params.Description != "",
		},

		Notes: sql.NullString{
			String: params.Notes,
			Valid:  params.Notes != "",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error storing to database", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, dbSessionToResponse(workoutSession))

}
func (cfg *apiConfig) handlerGetAllMyWorkoutSessions(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	var workSess []database.WorkoutSession
	workSess, err = cfg.db.GetAllWorkoutSessionsSorted(r.Context(), userID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}

	response := make([]workoutSessions, 0)

	for _, ws := range workSess {
		response = append(response, dbSessionToResponse(ws))
	}

	respondWithJSON(w, http.StatusOK, response)

}
func (cfg *apiConfig) handlerGetWorkoutSessionById(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	ID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	workSess, err := cfg.authorizeWorkoutSession(r.Context(), ID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, dbSessionToResponse(workSess))
}
func (cfg *apiConfig) handlerUpdateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		WorkoutDate time.Time `json:"workout_date"`
		Description string    `json:"description"`
		Notes       string    `json:"notes"`
	}
	var params parameters

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	ID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}

	_, err = cfg.authorizeWorkoutSession(r.Context(), ID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = cfg.db.UpdateWorkoutSession(r.Context(), database.UpdateWorkoutSessionParams{
		WorkoutDate: params.WorkoutDate,
		Description: sql.NullString{
			String: params.Description,
			Valid:  params.Description != "",
		},

		Notes: sql.NullString{
			String: params.Notes,
			Valid:  params.Notes != "",
		},
		ID: ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, msgResponse{"Updated session"})

}
func (cfg *apiConfig) handlerDeleteWorkoutSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	iD, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	_, err = cfg.authorizeWorkoutSession(r.Context(), iD, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = cfg.db.DeleteWorkoutSession(r.Context(), database.DeleteWorkoutSessionParams{
		ID:     iD,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete session", err)
		return
	}
	respondWithJSON(w, http.StatusOK, msgResponse{"Deleted session"})
}

func (cfg *apiConfig) handlerGetMyLastSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	workSess, err := cfg.db.GetLastWorkoutSession(r.Context(), userID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Not found", nil)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve records", err)
		return
	}
	respondWithJSON(w, http.StatusOK, dbSessionToResponse(workSess))
}

func (cfg *apiConfig) handlerGetMyXNumberLastSessions(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	lastXString := r.URL.Query().Get("lastX")
	if lastXString == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required query param: lastx", nil)
		return
	}

	parsedLastX, err := strconv.ParseInt(lastXString, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid number", err)
		return
	}

	if parsedLastX <= 0 {
		respondWithError(w, http.StatusBadRequest, "lastx must be greater than 0", nil)
		return
	}
	workSess, err := cfg.db.GetLastNWorkoutSessions(r.Context(), database.GetLastNWorkoutSessionsParams{
		UserID: userID,
		Limit:  int32(parsedLastX),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}
	response := make([]workoutSessions, 0)
	for _, ws := range workSess {
		response = append(response, dbSessionToResponse(ws))
	}
	respondWithJSON(w, http.StatusOK, response)
}
func (cfg *apiConfig) handlerCountSessions(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	count, err := cfg.db.GetWorkoutSessionsCount(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get count", err)
		return
	}
	type countResponse struct {
		Count int64 `json:"count"`
	}
	respondWithJSON(w, http.StatusOK, countResponse{Count: count})
}
func (cfg *apiConfig) handlerSearchWOSByDescription(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "Missing search query", nil)
		return
	}
	var workSess []database.WorkoutSession
	workSess, err = cfg.db.SearchWorkoutSessionsByDescription(r.Context(), database.SearchWorkoutSessionsByDescriptionParams{
		UserID: userID,
		Column2: sql.NullString{
			String: query,
			Valid:  true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search database", err)
		return
	}
	response := make([]workoutSessions, 0)
	for _, ws := range workSess {
		response = append(response, dbSessionToResponse(ws))
	}
	respondWithJSON(w, http.StatusOK, response)

}

func (cfg *apiConfig) handlerSearchWOSByDateRange(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	if startStr == "" || endStr == "" {
		respondWithError(w, http.StatusBadRequest, "2 dates are required to search", nil)
		return
	}
	startDate, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format, expected: 2006-01-02T15:04:05Z", err)
		return
	}
	endDate, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date format, expected: 2006-01-02T15:04:05Z", err)
		return
	}
	if startDate.After(endDate) {
		respondWithError(w, http.StatusBadRequest, "start date must be before end date", nil)
		return
	}
	var workSess []database.WorkoutSession
	workSess, err = cfg.db.SearchWorkoutSessionsByDateRange(r.Context(), database.SearchWorkoutSessionsByDateRangeParams{
		UserID:        userID,
		WorkoutDate:   startDate,
		WorkoutDate_2: endDate,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search database", err)
		return
	}
	response := make([]workoutSessions, 0)
	for _, ws := range workSess {
		response = append(response, dbSessionToResponse(ws))
	}
	respondWithJSON(w, http.StatusOK, response)

}

// response helper
func dbSessionToResponse(ws database.WorkoutSession) workoutSessions {
	return workoutSessions{
		ID:          ws.ID,
		UserID:      ws.UserID,
		WorkoutDate: ws.WorkoutDate,
		Description: nullStringToPtr(ws.Description),
		Notes:       nullStringToPtr(ws.Notes),
		CreatedAt:   ws.CreatedAt,
		UpdatedAt:   ws.UpdatedAt,
	}
}
