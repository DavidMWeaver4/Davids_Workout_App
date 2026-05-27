package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

//TODO

func (cfg *apiConfig) handlerCreateCardioAndSets(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
		SetNumber          int32     `json:"set_number"`
		Distance           float64   `json:"distance,omitempty"`
		IsKilometers       bool      `json:"is_kilometers"`
		DurationSeconds    int32     `json:"duration_seconds,omitempty"`
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

	cardioSet, err := cfg.db.CreateCardioAndSets(r.Context(), database.CreateCardioAndSetsParams{
		WorkoutExercisesID: params.WorkoutExercisesID,
		SetNumber:          params.SetNumber,
		Distance: sql.NullFloat64{
			Float64: params.Distance,
			Valid:   params.Distance != 0,
		},
		IsKilometers: params.IsKilometers,
		DurationSeconds: sql.NullInt32{
			Int32: params.DurationSeconds,
			Valid: params.DurationSeconds != 0,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing to database", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, cardioSetResponse(cardioSet))
}

func (cfg *apiConfig) handlerGetAllCardioFromSession(w http.ResponseWriter, r *http.Request) {
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
	setsInSession, err := cfg.db.GetAllCardioFromSession(r.Context(), workoutExerciseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}
	response := make([]cardioAndSets, 0)
	for _, sis := range setsInSession {
		response = append(response, cardioSetResponse(sis))
	}
	respondWithJSON(w, http.StatusOK, response)
}
func (cfg *apiConfig) handlerDeleteCardioAndSets(w http.ResponseWriter, r *http.Request) {
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
	err = cfg.db.DeleteCardioAndSets(r.Context(), database.DeleteCardioAndSetsParams{
		ID:                 weightSet.ID,
		WorkoutExercisesID: weightSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting from database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, msgResponse{"Deletion successful"})
}

func (cfg *apiConfig) handlerUpdateCardioAndSets(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
		SetNumber          int32     `json:"set_number"`
		Distance           float64   `json:"distance,omitempty"`
		IsKilometers       bool      `json:"is_kilometers"`
		DurationSeconds    int32     `json:"duration_seconds,omitempty"`
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
	updatedSet, err := cfg.db.UpdateCardioAndSets(r.Context(), database.UpdateCardioAndSetsParams{
		SetNumber: params.SetNumber,
		Distance: sql.NullFloat64{
			Float64: params.Distance,
			Valid:   params.Distance != 0,
		},
		IsKilometers: params.IsKilometers,
		DurationSeconds: sql.NullInt32{
			Int32: params.DurationSeconds,
			Valid: params.DurationSeconds != 0,
		},
		ID:                 setID,
		WorkoutExercisesID: weightSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, cardioSetResponse(updatedSet))
}

//TODO
// handlerGetAllSetsDistance
// handlerGetAllSetsDuration
// handlerGetCardioAndSetFromID
// handlerGetSetDistance
// handlerGetSetDuration
/*
*
// Helper functions
*
*/
func cardioSetResponse(ws database.CardioAndSet) cardioAndSets {
	return cardioAndSets{
		ID:                 ws.ID,
		WorkoutExercisesID: ws.WorkoutExercisesID,
		SetNumber:          ws.SetNumber,
		Distance:           nullFloat64ToPtr(ws.Distance),
		IsKilometers:       ws.IsKilometers,
		DurationSeconds:    nullInt32ToPtr(ws.DurationSeconds),
		CreatedAt:          ws.CreatedAt,
		UpdatedAt:          ws.UpdatedAt,
	}
}
