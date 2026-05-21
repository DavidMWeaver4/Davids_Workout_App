package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	database "command-line-arguments/Users/dmweaver/workspace/bootdotdev/workout_app/internal/database/weights_and_sets.sql.go"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

//TODO

func (cfg *apiConfig) handlerCreateWeightsAndSets(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
		Weight             string    `json:"weight"`
		IsKilogram         bool      `json:"is_kilogram"`
		SetNumber          int32     `json:"set_number"`
		RepsTarget         int32     `json:"reps_target"`
		RepsActual         int32     `json:"reps_actual"`
		DurationSeconds    int32     `json:"duration_seconds"`
		RestTimeSeconds    int32     `json:"rest_time_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	work_exercise, err := cfg.db.GetWorkoutExerciseFromID(r.Context(), params.WorkoutExercisesID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), work_exercise.WorkoutSessionID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your session", nil)
		return
	}
	weight_set, err := cfg.db.CreateWeightsAndSets(r.Context(), database.CreateWeightsAndSetsParams{
		WorkoutExercisesID: params.WorkoutExercisesID,
		Weight:             params.Weight,
		IsKilograms:        params.IsKilogram,
		SetNumber:          params.SetNumber,
		RepsTarget:         params.RepsTarget,
		RepsActual:         params.RepsActual,
		DurationSeconds: sql.NullInt32{
			Int32: params.DurationSeconds,
			Valid: params.DurationSeconds != 0,
		},
		RestTimeSeconds: sql.NullInt32{
			Int32: params.RestTimeSeconds,
			Valid: params.RestTimeSeconds != 0,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing to database", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, weightAndSets{
		ID:                 weight_set.ID,
		WorkoutExercisesID: weight_set.WorkoutExercisesID,
		Weight:             weight_set.Weight,
		IsKilograms:        weight_set.IsKilograms,
		SetNumber:          weight_set.SetNumber,
		RepsTarget:         weight_set.RepsTarget,
		RepsActual:         weight_set.RepsActual,
		DurationSeconds:    nullInt32ToPtr(weight_set.DurationSeconds),
		RestTimeSeconds:    nullInt32ToPtr(weight_set.RestTimeSeconds),
		CreatedAt:          weight_set.CreatedAt,
		UpdatedAt:          weight_set.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerGetAllSetsFromSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.PathValue("id")
	var workoutExerciseID uuid.UUID
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	workoutExerciseID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	work_exercise, err := cfg.db.GetWorkoutExerciseFromID(r.Context(), workoutExerciseID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), work_exercise.WorkoutSessionID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your session", nil)
		return
	}
	setsInSession, err := cfg.db.GetAllSetsFromSession(r.Context(), workoutExerciseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}
	response := make([]weightAndSets, 0)
	for _, sis := range setsInSession {
		response = append(response, weightAndSets{
			ID:                 sis.ID,
			WorkoutExercisesID: sis.WorkoutExercisesID,
			Weight:             sis.Weight,
			IsKilograms:        sis.IsKilograms,
			SetNumber:          sis.SetNumber,
			RepsTarget:         sis.RepsTarget,
			RepsActual:         sis.RepsActual,
			DurationSeconds:    nullInt32ToPtr(sis.DurationSeconds),
			RestTimeSeconds:    nullInt32ToPtr(sis.RestTimeSeconds),
			CreatedAt:          sis.CreatedAt,
			UpdatedAt:          sis.UpdatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerDeleteWeightandSet(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.PathValue("id")
	var weightSetID uuid.UUID
	if idString == "" {
		respondWithError(w, http.StatusBadRequest, "No ID provided", nil)
		return
	}
	weightSetID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	weightSet, err := cfg.db.GetWeightSetByID(r.Context(), weightSetID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Weight set not found", err)
		return
	}
	work_exercise, err := cfg.db.GetWorkoutExerciseFromID(r.Context(), weightSet.WorkoutExercisesID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), work_exercise.WorkoutSessionID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your session", nil)
		return
	}
	err = cfg.db.DeleteWeightAndSets(r.Context(), database.DeleteWeightAndSetsParams{
		ID:                 weightSet.ID,
		WorkoutExercisesID: weightSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting from database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, msgResponse{"Deletion successful"})
}
