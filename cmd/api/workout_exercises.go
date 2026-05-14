package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		ExerciseName     string    `json:"exercise_name"`
		WorkoutSessionID uuid.UUID `json:"workout_session_id"`
		OrderIndex       int32     `json:"order_index"`
		Notes            string    `json:"notes"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), params.WorkoutSessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your session", nil)
		return
	}
	if params.ExerciseName == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter an exercise name", nil)
		return
	}
	exercise, err := cfg.db.GetExerciseFromName(r.Context(), params.ExerciseName)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot find exercise", err)
		return
	}
	workExer, err := cfg.db.CreateWorkoutExercises(r.Context(), database.CreateWorkoutExercisesParams{
		WorkoutSessionID: params.WorkoutSessionID,
		ExerciseID:       exercise.ID,
		OrderIndex:       params.OrderIndex,
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
	respondWithJSON(w, http.StatusCreated, workoutExercise{
		ID:               workExer.ID,
		WorkoutSessionID: workExer.WorkoutSessionID,
		ExerciseID:       workExer.ExerciseID,
		OrderIndex:       workExer.OrderIndex,
		Notes:            nullStringToPtr(workExer.Notes),
		CreatedAt:        workExer.CreatedAt,
		UpdatedAt:        workExer.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerGetWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.PathValue("id")
	var execID uuid.UUID
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	execID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	exercise, err := cfg.db.GetWorkoutExerciseFromID(r.Context(), execID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), exercise.WorkoutSessionID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your exercise", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, workoutExercise{
		ID:               exercise.ID,
		WorkoutSessionID: exercise.WorkoutSessionID,
		ExerciseID:       exercise.ExerciseID,
		OrderIndex:       exercise.OrderIndex,
		Notes:            nullStringToPtr(exercise.Notes),
		CreatedAt:        exercise.CreatedAt,
		UpdatedAt:        exercise.UpdatedAt,
	})

}
func (cfg *apiConfig) handlerGetWorkoutExercisesInSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.PathValue("id")
	var exeSessID uuid.UUID
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	exeSessID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), exeSessID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your exercise session", nil)
		return
	}
	var sessionExercises []database.WorkoutExercise
	sessionExercises, err = cfg.db.GetWorkoutsInSession(r.Context(), exeSessID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	response := make([]workoutExercise, 0)
	for _, ws := range sessionExercises {
		response = append(response, workoutExercise{
			ID:               ws.ID,
			WorkoutSessionID: ws.WorkoutSessionID,
			ExerciseID:       ws.ExerciseID,
			OrderIndex:       ws.OrderIndex,
			Notes:            nullStringToPtr(ws.Notes),
			CreatedAt:        ws.CreatedAt,
			UpdatedAt:        ws.UpdatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, response)

}

func (cfg *apiConfig) handlerDeleteWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.PathValue("id")
	var workoutExercise uuid.UUID
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	workoutExercise, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}
	exercise, err := cfg.db.GetWorkoutExerciseFromID(r.Context(), workoutExercise)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}
	workoutSession := exercise.WorkoutSessionID
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), workoutSession)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your exercise session", nil)
		return
	}
	err = cfg.db.DeleteWorkoutExercises(r.Context(), database.DeleteWorkoutExercisesParams{
		ID:               workoutExercise,
		WorkoutSessionID: workoutSession,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete from database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, msgResponse{"successfully deleted"})
}

func (cfg *apiConfig) handlerGetWorkoutsInSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.PathValue("id")
	var workoutSession uuid.UUID
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	workoutSession, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), workoutSession)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your exercise session", nil)
		return
	}
	var workExer []database.WorkoutExercise
	workExer, err = cfg.db.GetWorkoutsInSession(r.Context(), workoutSession)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve from database", err)
		return
	}
	response := make([]workoutExercise, 0)
	for _, ws := range workExer {
		response = append(response, workoutExercise{
			ID:               ws.ID,
			WorkoutSessionID: ws.WorkoutSessionID,
			ExerciseID:       ws.ExerciseID,
			OrderIndex:       ws.OrderIndex,
			Notes:            nullStringToPtr(ws.Notes),
			CreatedAt:        ws.CreatedAt,
			UpdatedAt:        ws.UpdatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, response)

}

/*
 * updateeworkoutexerciseorder
 * getworkoutexercisecount
 *
 */
