package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

//TODO

type volumeResponse struct {
	Volume float64 `json:"volume"`
}

type durationResponse struct {
	TotalSeconds int32 `json:"total_seconds"`
}

func (cfg *apiConfig) handlerCreateWeightsAndSets(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	_, err = cfg.authorizeWorkoutExercise(r.Context(), params.WorkoutExercisesID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}
	weightSet, err := cfg.db.CreateWeightsAndSets(r.Context(), database.CreateWeightsAndSetsParams{
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
	_, err = cfg.authorizeWorkoutExercise(r.Context(), workoutExerciseID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}
	setsInSession, err := cfg.db.GetAllSetsFromSession(r.Context(), workoutExerciseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}
	response := make([]weightAndSets, 0)
	for _, sis := range setsInSession {
		response = append(response, weightSetResponse(sis))
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
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
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

func (cfg *apiConfig) handlerUpdateWeightAndSets(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
		Weight             float64   `json:"weight"`
		IsKilograms        bool      `json:"is_kilogram"`
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
	idString := r.PathValue("id")
	setID, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), setID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}
	updatedSet, err := cfg.db.UpdateWeightsAndSets(r.Context(), database.UpdateWeightsAndSetsParams{
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
	idString := r.PathValue("id")
	var weightSetID uuid.UUID

	weightSetID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
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
func (cfg *apiConfig) handlerGetTotalVolumeFromAllSet(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.PathValue("id")
	var weightSetID uuid.UUID

	weightSetID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}
	setVolume, err := cfg.db.GetTotalVolumeFromAllSets(r.Context(), weightSet.WorkoutExercisesID)
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
	idString := r.PathValue("id")
	var weightSetID uuid.UUID

	weightSetID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
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

func (cfg *apiConfig) handlerGetTotalDurationFromAllSets(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	idString := r.PathValue("id")
	var weightSetID uuid.UUID

	weightSetID, err = uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	weightSet, err := cfg.authorizeWeightSet(r.Context(), weightSetID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}
	totalDuration, err := cfg.db.GetTotalDurationForAllSets(r.Context(), weightSet.WorkoutExercisesID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get duration", err)
		return
	}
	respondWithJSON(w, http.StatusOK, durationResponse{TotalSeconds: totalDuration})
}

/*
*
// Helper functions
*
*/
func (cfg *apiConfig) authorizeWorkoutExercise(ctx context.Context, workoutExerciseID uuid.UUID, userID uuid.UUID) (database.WorkoutExercise, error) {

	workExercise, err := cfg.db.GetWorkoutExerciseFromID(ctx, workoutExerciseID)
	if err != nil {
		return database.WorkoutExercise{}, err
	}

	workSess, err := cfg.db.GetWorkoutSessionByID(ctx, workExercise.WorkoutSessionID)
	if err != nil {
		return database.WorkoutExercise{}, err
	}

	if workSess.UserID != userID {
		return database.WorkoutExercise{}, errors.New("forbidden")
	}

	return workExercise, nil
}

func (cfg *apiConfig) authorizeWeightSet(ctx context.Context, weightSetID uuid.UUID, userID uuid.UUID) (database.WeightsAndSet, error) {

	weightSet, err := cfg.db.GetWeightAndSetFromID(ctx, weightSetID)
	if err != nil {
		return database.WeightsAndSet{}, err
	}

	_, err = cfg.authorizeWorkoutExercise(ctx, weightSet.WorkoutExercisesID, userID)
	if err != nil {
		return database.WeightsAndSet{}, err
	}

	return weightSet, nil
}

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
