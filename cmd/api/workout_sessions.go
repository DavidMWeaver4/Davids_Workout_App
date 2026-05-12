package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
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
	respondWithJSON(w, http.StatusCreated, workoutSessions{
		ID:          workoutSession.ID,
		UserID:      workoutSession.UserID,
		WorkoutDate: workoutSession.WorkoutDate,
		Description: workoutSession.Description,
		Notes:       workoutSession.Notes,
		CreatedAt:   workoutSession.CreatedAt,
		UpdatedAt:   workoutSession.UpdatedAt,
	})

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
		response = append(response, workoutSessions{
			ID:          ws.ID,
			UserID:      ws.UserID,
			WorkoutDate: ws.WorkoutDate,
			Description: ws.Description,
			Notes:       ws.Notes,
			CreatedAt:   ws.CreatedAt,
			UpdatedAt:   ws.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, response)

}
func (cfg *apiConfig) handlerGetWorkoutSessionById(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.URL.Query().Get("id")
	var ID uuid.UUID
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	ID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your session", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, workoutSessions{
		ID:          workSess.ID,
		UserID:      workSess.UserID,
		WorkoutDate: workSess.WorkoutDate,
		Description: workSess.Description,
		Notes:       workSess.Notes,
		CreatedAt:   workSess.CreatedAt,
		UpdatedAt:   workSess.UpdatedAt,
	})
}
func (cfg *apiConfig) handlerUpdateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		ID          uuid.UUID `json:"id"`
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
	checkSession, err := cfg.db.GetWorkoutSessionByID(r.Context(), params.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't check database", err)
		return
	}
	if checkSession.UserID != userID {
		respondWithError(w, http.StatusForbidden, "This sessions does not belong to you", nil)
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
		ID: params.ID,
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
	idString := r.URL.Query().Get("id")
	var iD uuid.UUID
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	iD, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
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
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No sessions found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, workoutSessions{
		ID:          workSess.ID,
		UserID:      workSess.UserID,
		WorkoutDate: workSess.WorkoutDate,
		Description: workSess.Description,
		Notes:       workSess.Notes,
		CreatedAt:   workSess.CreatedAt,
		UpdatedAt:   workSess.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerGetMyXNumberLastSessions(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	lastXString := r.URL.Query().Get("lastx")
	var lastX int
	lastX, err = strconv.Atoi(lastXString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid number", err)
		return
	}
	if lastX <= 0 {
		lastX = 1
	}
	workSess, err := cfg.db.GetLastNWorkoutSessions(r.Context(), database.GetLastNWorkoutSessionsParams{
		UserID: userID,
		Limit:  int32(lastX),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}
	response := make([]workoutSessions, 0)
	for _, ws := range workSess {
		response = append(response, workoutSessions{
			ID:          ws.ID,
			UserID:      ws.UserID,
			WorkoutDate: ws.WorkoutDate,
			Description: ws.Description,
			Notes:       ws.Notes,
			CreatedAt:   ws.CreatedAt,
			UpdatedAt:   ws.UpdatedAt,
		})
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
	response := fmt.Sprintf("Your total workout sessions is: %d", count)
	respondWithJSON(w, http.StatusOK, msgResponse{response})
}
func (cfg *apiConfig) handlerSearchWOSByDescription(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	query := r.URL.Query().Get("query")
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
		response = append(response, workoutSessions{
			ID:          ws.ID,
			UserID:      ws.UserID,
			WorkoutDate: ws.WorkoutDate,
			Description: ws.Description,
			Notes:       ws.Notes,
			CreatedAt:   ws.CreatedAt,
			UpdatedAt:   ws.UpdatedAt,
		})
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
		response = append(response, workoutSessions{
			ID:          ws.ID,
			UserID:      ws.UserID,
			WorkoutDate: ws.WorkoutDate,
			Description: ws.Description,
			Notes:       ws.Notes,
			CreatedAt:   ws.CreatedAt,
			UpdatedAt:   ws.UpdatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, response)

}
