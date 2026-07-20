package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

type distanceResponse struct {
	Distance float64 `json:"distance"`
}

type cardioDurationResponse struct {
	TotalSeconds int64 `json:"total_seconds"`
}

func (cfg *apiConfig) handlerCreateCardioSet(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	_, err = cfg.authorizeWorkoutExercise(r.Context(), params.WorkoutExercisesID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = validateCardioSet(params.SetNumber, params.Distance, params.DurationSeconds)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	cardioSet, err := cfg.db.CreateCardioSet(r.Context(), database.CreateCardioSetParams{
		ID:                 uuid.New(),
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
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error writing to database", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, cardioSetResponse(cardioSet))
}

func (cfg *apiConfig) handlerGetAllCardioFromExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	workoutExerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	_, err = cfg.authorizeWorkoutExercise(r.Context(), workoutExerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	setsInSession, err := cfg.db.GetAllCardioFromExercise(r.Context(), workoutExerciseID)
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
func (cfg *apiConfig) handlerDeleteCardioSet(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	cardioSetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	cardioSet, err := cfg.authorizeCardioSet(r.Context(), cardioSetID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = cfg.db.DeleteCardioSet(r.Context(), database.DeleteCardioSetParams{
		ID:                 cardioSet.ID,
		WorkoutExercisesID: cardioSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting from database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, msgResponse{"Deletion successful"})
}

func (cfg *apiConfig) handlerUpdateCardioSet(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct {
		SetNumber       int32   `json:"set_number"`
		Distance        float64 `json:"distance,omitempty"`
		IsKilometers    bool    `json:"is_kilometers"`
		DurationSeconds int32   `json:"duration_seconds,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	setID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	cardioSet, err := cfg.authorizeCardioSet(r.Context(), setID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	err = validateCardioSet(params.SetNumber, params.Distance, params.DurationSeconds)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	updatedSet, err := cfg.db.UpdateCardioSet(r.Context(), database.UpdateCardioSetParams{
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
		WorkoutExercisesID: cardioSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, cardioSetResponse(updatedSet))
}

func (cfg *apiConfig) handlerGetSetDistance(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	cardioSetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}
	cardioSet, err := cfg.authorizeCardioSet(r.Context(), cardioSetID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	setDistance, err := cfg.db.GetSetDistance(r.Context(), database.GetSetDistanceParams{
		ID:                 cardioSet.ID,
		WorkoutExercisesID: cardioSet.WorkoutExercisesID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get distance", err)
		return
	}
	respondWithJSON(w, http.StatusOK, distanceResponse{Distance: setDistance.Float64})

}
func (cfg *apiConfig) handlerGetExerciseDistance(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}

	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}

	_, err = cfg.authorizeWorkoutExercise(r.Context(), exerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}

	distance, err := cfg.db.GetAllSetsDistance(r.Context(), exerciseID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get distance", err)
		return
	}

	respondWithJSON(w, http.StatusOK, distanceResponse{
		Distance: float64(distance),
	})
}
func (cfg *apiConfig) handlerGetCardioExerciseDuration(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}

	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}

	_, err = cfg.authorizeWorkoutExercise(r.Context(), exerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}

	duration, err := cfg.db.GetAllCardioSetsDuration(r.Context(), exerciseID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get duration", err)
		return
	}

	respondWithJSON(w, http.StatusOK, cardioDurationResponse{
		TotalSeconds: duration,
	})
}

func (cfg *apiConfig) handlerGetAllCardioSessionSets(w http.ResponseWriter, r *http.Request) {
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

	sets, err := cfg.db.GetAllCardioSetsFromSession(r.Context(), sessionID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve sets", err)
		return
	}

	response := make([]cardioAndSets, 0, len(sets))

	for _, s := range sets {
		response = append(response, cardioSetResponse(s))
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetAllCardioSessionDuration(w http.ResponseWriter, r *http.Request) {
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

	duration, err := cfg.db.GetTotalSessionCardioDuration(r.Context(), sessionID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get duration", err)
		return
	}

	respondWithJSON(w, http.StatusOK, cardioDurationResponse{
		TotalSeconds: duration,
	})
}
func (cfg *apiConfig) handlerGetAllSessionDistance(w http.ResponseWriter, r *http.Request) {
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

	distance, err := cfg.db.GetTotalSessionDistance(r.Context(), sessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get distance", err)
		return
	}

	respondWithJSON(w, http.StatusOK, distanceResponse{
		Distance: distance,
	})
}
func (cfg *apiConfig) handlerGetCardioSetDuration(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}

	cardioSetID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid set ID", err)
		return
	}

	cardioSet, err := cfg.authorizeCardioSet(r.Context(), cardioSetID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}

	duration, err := cfg.db.GetCardioSetDuration(r.Context(), database.GetCardioSetDurationParams{
		ID:                 cardioSet.ID,
		WorkoutExercisesID: cardioSet.WorkoutExercisesID,
	},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get duration", err)
		return
	}

	respondWithJSON(w, http.StatusOK, cardioDurationResponse{
		TotalSeconds: int64(duration.Int32),
	})
}

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

func validateCardioSet(setNumber int32, distance float64, durationSeconds int32) error {

	if setNumber < 1 {
		return errors.New("set_number must be greater than 0")
	}

	if distance < 0 {
		return errors.New("distance cannot be negative")
	}

	if durationSeconds < 0 {
		return errors.New("duration cannot be negative")
	}

	return nil
}
