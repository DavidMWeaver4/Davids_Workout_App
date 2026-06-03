package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

// TODO
func (cfg *apiConfig) handlerCreateExercises(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		ExerciseName    string   `json:"exercise_name"`
		TargetMuscles   []string `json:"target_muscles"`
		Equipment       string   `json:"equipment,omitempty"`
		DifficultyLevel string   `json:"difficulty_level,omitempty"`
		Description     string   `json:"description,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if params.ExerciseName == "" {
		respondWithError(w, http.StatusBadRequest, "exercise name is required", nil)
		return
	}
	exercise, err := cfg.db.CreateExercises(r.Context(), database.CreateExercisesParams{
		ID: uuid.New(),
		UserID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
		ExerciseName:  params.ExerciseName,
		TargetMuscles: params.TargetMuscles,
		Equipment: sql.NullString{
			String: params.Equipment,
			Valid:  params.Equipment != "",
		},
		DifficultyLevel: sql.NullString{
			String: params.DifficultyLevel,
			Valid:  params.DifficultyLevel != "",
		},
		Description: sql.NullString{
			String: params.Description,
			Valid:  params.Description != "",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, dbExerciseLibraryToResponse(exercise))
}

func (cfg *apiConfig) handlerDeleteExercisesByID(w http.ResponseWriter, r *http.Request) {
	//TODO
	// ownership check
}
func (cfg *apiConfig) handlerGetExerciseFromID(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func (cfg *apiConfig) handlerListAvaiableExercises(w http.ResponseWriter, r *http.Request) {
	//TODO
}
func (cfg *apiConfig) handlerUpdateExercise(w http.ResponseWriter, r *http.Request) {
	//TODO
	//ownership check
}

/*
 *
 *
 */
//Helper functions
func dbExerciseLibraryToResponse(ws database.Exercise) exercise {
	return exercise{
		ID:              ws.ID,
		UserID:          nullUUIDToUUID(ws.UserID),
		ExerciseName:    ws.ExerciseName,
		TargetMuscles:   ws.TargetMuscles,
		Equipment:       nullStringToPtr(ws.Equipment),
		DifficultyLevel: nullStringToPtr(ws.DifficultyLevel),
		Description:     nullStringToPtr(ws.Description),
		CreatedAt:       ws.CreatedAt,
		UpdatedAt:       ws.UpdatedAt,
	}
}
