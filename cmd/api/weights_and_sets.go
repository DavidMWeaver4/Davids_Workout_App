package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

type volumeResponse struct {
	Volume float64 `json:"volume"`
}

type durationResponse struct {
	TotalSeconds int32 `json:"total_seconds"`
}

func (cfg *apiConfig) handlerCreateWeightAndSet(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
		Weight             float64   `json:"weight"`
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
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	_, err = cfg.authorizeWorkoutExercise(r.Context(), params.WorkoutExercisesID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = validateWeightSet(params.Weight, params.SetNumber, params.RepsTarget, params.RepsActual, params.DurationSeconds, params.RestTimeSeconds)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	weightSet, err := cfg.db.CreateWeightAndSet(r.Context(), database.CreateWeightAndSetParams{
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
	respondWithJSON(w, http.StatusCreated, weightSetResponse(weightSet))
}

func (cfg *apiConfig) handlerGetAllSetsFromExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	workoutExerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workout exerciseID", err)
		return
	}
	_, err = cfg.authorizeWorkoutExercise(r.Context(), workoutExerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	setsInSession, err := cfg.db.GetAllSetsFromExercise(r.Context(), workoutExerciseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout exercise", err)
		return
	}
	response := make([]weightAndSets, 0)
	for _, sis := range setsInSession {
		response = append(response, weightSetResponse(sis))
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerDeleteWeightAndSet(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	weightSetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}

	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = cfg.db.DeleteWeightAndSet(r.Context(), database.DeleteWeightAndSetParams{
		ID:                 weightSet.ID,
		WorkoutExercisesID: weightSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting from database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, msgResponse{"Deletion successful"})
}

func (cfg *apiConfig) handlerUpdateWeightAndSet(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		Weight          float64 `json:"weight"`
		IsKilograms     bool    `json:"is_kilogram"`
		SetNumber       int32   `json:"set_number"`
		RepsTarget      int32   `json:"reps_target"`
		RepsActual      int32   `json:"reps_actual"`
		DurationSeconds int32   `json:"duration_seconds"`
		RestTimeSeconds int32   `json:"rest_time_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	setID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}

	weightSet, err := cfg.authorizeWeightSet(r.Context(), setID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = validateWeightSet(params.Weight, params.SetNumber, params.RepsTarget, params.RepsActual, params.DurationSeconds, params.RestTimeSeconds)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	updatedSet, err := cfg.db.UpdateWeightAndSet(r.Context(), database.UpdateWeightAndSetParams{
		Weight:      params.Weight,
		IsKilograms: params.IsKilograms,
		SetNumber:   params.SetNumber,
		RepsTarget:  params.RepsTarget,
		RepsActual:  params.RepsActual,
		DurationSeconds: sql.NullInt32{
			Int32: params.DurationSeconds,
			Valid: params.DurationSeconds != 0,
		},
		RestTimeSeconds: sql.NullInt32{
			Int32: params.RestTimeSeconds,
			Valid: params.RestTimeSeconds != 0,
		},
		ID:                 setID,
		WorkoutExercisesID: weightSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, weightSetResponse(updatedSet))
}

func (cfg *apiConfig) handlerGetVolumeSet(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	weightSetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	setVolume, err := cfg.db.GetSetVolume(r.Context(), database.GetSetVolumeParams{
		ID:                 weightSet.ID,
		WorkoutExercisesID: weightSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get volume", err)
		return
	}
	respondWithJSON(w, http.StatusOK, volumeResponse{Volume: float64(setVolume)})
}

func (cfg *apiConfig) handlerGetTotalVolumeFromExerciseSets(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	weightSetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	setVolume, err := cfg.db.GetTotalVolumeFromExerciseSets(r.Context(), weightSet.WorkoutExercisesID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get volume", err)
		return
	}

	respondWithJSON(w, http.StatusOK, volumeResponse{Volume: setVolume})
}
func (cfg *apiConfig) handlerGetTotalDuration(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	weightSetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	totalDuration, err := cfg.db.GetTotalDuration(r.Context(), database.GetTotalDurationParams{
		ID:                 weightSetID,
		WorkoutExercisesID: weightSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get duration", err)
		return
	}
	respondWithJSON(w, http.StatusOK, durationResponse{TotalSeconds: totalDuration})
}

func (cfg *apiConfig) handlerGetTotalDurationForExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	weightSetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	totalDuration, err := cfg.db.GetTotalDurationForExercise(r.Context(), weightSet.WorkoutExercisesID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get duration", err)
		return
	}
	respondWithJSON(w, http.StatusOK, durationResponse{TotalSeconds: totalDuration})
}
func (cfg *apiConfig) handlerGetAllSetsFromSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}

	_, err = cfg.authorizeWorkoutSession(r.Context(), sessionID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}

	sets, err := cfg.db.GetAllWeightSetsFromSession(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve sets", err)
		return
	}

	response := make([]weightAndSets, 0, len(sets))
	for _, s := range sets {
		response = append(response, weightSetResponse(s))
	}

	respondWithJSON(w, http.StatusOK, response)
}
func (cfg *apiConfig) handlerGetTotalSessionVolume(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}

	_, err = cfg.authorizeWorkoutSession(r.Context(), sessionID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}

	volume, err := cfg.db.GetTotalSessionVolume(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get session volume", err)
		return
	}

	respondWithJSON(w, http.StatusOK, volumeResponse{
		Volume: volume,
	})
}

/*
*
// Helper functions
*
*/

func weightSetResponse(ws database.WeightsAndSet) weightAndSets {
	return weightAndSets{
		ID:                 ws.ID,
		WorkoutExercisesID: ws.WorkoutExercisesID,
		Weight:             ws.Weight,
		IsKilograms:        ws.IsKilograms,
		SetNumber:          ws.SetNumber,
		RepsTarget:         ws.RepsTarget,
		RepsActual:         ws.RepsActual,
		DurationSeconds:    nullInt32ToPtr(ws.DurationSeconds),
		RestTimeSeconds:    nullInt32ToPtr(ws.RestTimeSeconds),
		CreatedAt:          ws.CreatedAt,
		UpdatedAt:          ws.UpdatedAt,
	}
}
func validateWeightSet(weight float64, setNumber int32, repsTarget int32, repsActual int32, durationSeconds int32, restTimeSeconds int32) error {

	if weight < 0 {
		return errors.New("weight cannot be negative")
	}

	if setNumber < 1 {
		return errors.New("set_number must be greater than 0")
	}

	if repsTarget < 0 {
		return errors.New("reps_target cannot be negative")
	}

	if repsActual < 0 {
		return errors.New("reps_actual cannot be negative")
	}

	if durationSeconds < 0 {
		return errors.New("duration_seconds cannot be negative")
	}

	if restTimeSeconds < 0 {
		return errors.New("rest_time_seconds cannot be negative")
	}

	return nil
}
