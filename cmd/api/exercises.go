package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

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
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}
	authorizedExercise, err := cfg.authorizeExercise(r.Context(), exerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = cfg.db.DeleteExerciseByID(r.Context(), database.DeleteExerciseByIDParams{
		ID:     authorizedExercise.ID,
		UserID: authorizedExercise.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete exercise", err)
		return
	}
	respondWithJSON(w, http.StatusOK, msgResponse{"Deleted successful"})

}
func (cfg *apiConfig) handlerGetExerciseFromID(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}
	exercise, err := cfg.db.GetExerciseFromID(r.Context(), exerciseID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Exercise not found", nil)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get exercise", err)
		return
	}
	respondWithJSON(w, http.StatusOK, dbExerciseLibraryToResponse(exercise))
}
func (cfg *apiConfig) handlerSearchExercises(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	nullUserID := uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	}
	muscle := r.URL.Query().Get("muscle")
	equipment := r.URL.Query().Get("equipment")
	difficulty := r.URL.Query().Get("difficulty")
	var muscleParam []string

	if muscle != "" {
		muscleParam = []string{muscle}
	}

	exercises, err := cfg.db.SearchExercises(r.Context(), database.SearchExercisesParams{
		UserID:  nullUserID,
		Column2: difficulty,
		Column3: equipment,
		Column4: muscleParam,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve exercises", err)
		return
	}

	response := make([]exercise, 0)
	for _, el := range exercises {
		response = append(response, dbExerciseLibraryToResponse(el))
	}
	respondWithJSON(w, http.StatusOK, response)
}
func (cfg *apiConfig) handlerUpdateExercise(w http.ResponseWriter, r *http.Request) {
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
	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	authorizedExercise, err := cfg.authorizeExercise(r.Context(), exerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	updatedExercise, err := cfg.db.UpdateExercise(r.Context(), database.UpdateExerciseParams{
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
		ID:     authorizedExercise.ID,
		UserID: authorizedExercise.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, dbExerciseLibraryToResponse(updatedExercise))
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
