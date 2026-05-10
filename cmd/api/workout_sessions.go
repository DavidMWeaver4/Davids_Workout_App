package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		WorkoutDate time.Time `json:"workout_date"`
		Description string    `json:"description"`
		Notes       string    `json:"notes"`
	}
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
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
	respondWithJSON(w, http.StatusOK, workoutSessions{
		UserID:      workoutSession.UserID,
		WorkoutDate: workoutSession.WorkoutDate,
		Description: workoutSession.Description,
		Notes:       workoutSession.Notes,
		CreatedAt:   workoutSession.CreatedAt,
		UpdatedAt:   workoutSession.UpdatedAt,
	})

}
func (cfg *apiConfig) handlerGetMyWorkoutSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
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
func (cfg *apiConfig) handlerGetWorkoutSessionById(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("ID")
	var ID uuid.UUID
	if idString == "" {
		cfg.handlerGetMyWorkoutSession(w, r)
		return
	}
	ID, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
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
	//TODO
}
func (cfg *apiConfig) handlerDeleteWorkoutSession(w http.ResponseWriter, r *http.Request) {
	//TODO
}

/*
  POST   /api/v1/workout-sessions
 GET    /api/v1/workout-sessions
 GET    /api/v1/workout-sessions/{id}
 PUT    /api/v1/workout-sessions/{id}
 DELETE /api/v1/workout-sessions/{id}


 //after
 GET    /api/v1/workout-sessions/last
 GET    /api/v1/workout-sessions/count
 GET    /api/v1/workout-sessions/search
*/
