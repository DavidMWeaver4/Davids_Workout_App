package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	_, err = cfg.authorizeWorkoutSession(r.Context(), params.WorkoutSessionID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	exerciseName := strings.TrimSpace(params.ExerciseName)
	if exerciseName == "" {
		respondWithError(w, http.StatusBadRequest, "Please enter an exercise name", nil)
		return
	}
	exercise, err := cfg.db.GetExerciseFromName(r.Context(), exerciseName)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Exercise not found", nil)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot find exercise", err)
		return
	}
	if params.OrderIndex < 0 {
		respondWithError(w, http.StatusBadRequest,
			"order_index must be non-negative", nil)
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
	respondWithJSON(w, http.StatusCreated, dbExerciseToResponse(workExer))
}

func (cfg *apiConfig) handlerGetWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	exerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}
	exercise, err := cfg.authorizeWorkoutExercise(r.Context(), exerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, dbExerciseToResponse(exercise))

}
func (cfg *apiConfig) handlerGetWorkoutExercisesInSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	exerciseSessionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID", err)
		return
	}
	_, err = cfg.authorizeWorkoutSession(r.Context(), exerciseSessionID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	var sessionExercises []database.WorkoutExercise
	sessionExercises, err = cfg.db.GetWorkoutsInSession(r.Context(), exerciseSessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout exercises", err)
		return
	}
	response := make([]workoutExercise, 0)
	for _, ws := range sessionExercises {
		response = append(response, dbExerciseToResponse(ws))
	}
	respondWithJSON(w, http.StatusOK, response)

}

func (cfg *apiConfig) handlerDeleteWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	workoutExerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}
	exercise, err := cfg.authorizeWorkoutExercise(r.Context(), workoutExerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}

	err = cfg.db.DeleteWorkoutExercises(r.Context(), database.DeleteWorkoutExercisesParams{
		ID:               workoutExerciseID,
		WorkoutSessionID: exercise.WorkoutSessionID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete from database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, msgResponse{"successfully deleted"})
}

func (cfg *apiConfig) handlerGetNumOfWorkoutsInSession(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	workoutSessionID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
		return
	}
	_, err = cfg.authorizeWorkoutSession(r.Context(), workoutSessionID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}
	count, err := cfg.db.GetNumberOfExercisesInWorkout(r.Context(), workoutSessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve from database", err)
		return
	}
	type countResponse struct {
		Count int64 `json:"count"`
	}
	respondWithJSON(w, http.StatusOK, countResponse{Count: count})
}

/*
//This func has been depreciated due to its function replaced by the more roboust handlerUpdateWorkoutExercise

	func (cfg *apiConfig) handlerUpdateWorkoutExerciseOrder(w http.ResponseWriter, r *http.Request) {
		userID, err := cfg.getUserIDFromToken(w, r)
		if err != nil {
			return
		}

		exerciseID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid exercise ID", err)
			return
		}
		workExercise, err := cfg.authorizeWorkoutExercise(r.Context(), exerciseID, userID)
		if err != nil {
			respondWithAuthError(w, err)
			return
		}
		type parameters struct {
			OrderIndex int32 `json:"order_index"`
		}

		var params parameters

		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}
		if params.OrderIndex < 0 {
			respondWithError(w, http.StatusBadRequest, "order_index must be non-negative", nil)
			return
		}
		updatedExercise, err := cfg.db.UpdateWorkoutExerciseOrder(r.Context(), database.UpdateWorkoutExerciseOrderParams{
			OrderIndex:       params.OrderIndex,
			ID:               workExercise.ID,
			WorkoutSessionID: workExercise.WorkoutSessionID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update exercise", err)
			return
		}

		respondWithJSON(w, http.StatusOK, dbExerciseToResponse(updatedExercise))
	}
*/
func (cfg *apiConfig) handlerUpdateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}

	workoutExerciseID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid workout exercise ID", err)
		return
	}

	workoutExercise, err := cfg.authorizeWorkoutExercise(r.Context(), workoutExerciseID, userID)
	if err != nil {
		respondWithAuthError(w, err)
		return
	}

	type parameters struct {
		ExerciseName *string `json:"exercise_name,omitempty"`
		OrderIndex   *int32  `json:"order_index,omitempty"`
		Notes        *string `json:"notes,omitempty"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	exerciseID := workoutExercise.ExerciseID
	orderIndex := workoutExercise.OrderIndex
	notes := workoutExercise.Notes
	changed := false

	if params.ExerciseName != nil {
		exerciseName := strings.TrimSpace(*params.ExerciseName)
		changed = true
		if exerciseName == "" {
			respondWithError(w, http.StatusBadRequest, "Please enter an exercise name", nil)
			return
		}

		exercise, err := cfg.db.GetExerciseFromName(r.Context(), exerciseName)
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Exercise not found", nil)
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve exercise", err)
			return
		}

		exerciseID = exercise.ID
	}

	if params.OrderIndex != nil {
		changed = true
		if *params.OrderIndex < 0 {
			respondWithError(w, http.StatusBadRequest, "order_index must be non-negative", nil)
			return
		}
		orderIndex = *params.OrderIndex
	}

	if params.Notes != nil {
		changed = true
		notes = sql.NullString{
			String: *params.Notes,
			Valid:  *params.Notes != "",
		}
	}
	if !changed {
		respondWithError(w, http.StatusBadRequest, "No fields provided to update", nil)
		return
	}
	updatedExercise, err := cfg.db.UpdateWorkoutExercise(r.Context(), database.UpdateWorkoutExerciseParams{
		ExerciseID:       exerciseID,
		OrderIndex:       orderIndex,
		Notes:            notes,
		ID:               workoutExercise.ID,
		WorkoutSessionID: workoutExercise.WorkoutSessionID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update workout exercise", err)
		return
	}

	respondWithJSON(w, http.StatusOK, dbExerciseToResponse(updatedExercise))
}

/*
 * helpers
 *
 */
func dbExerciseToResponse(we database.WorkoutExercise) workoutExercise {
	return workoutExercise{
		ID:               we.ID,
		WorkoutSessionID: we.WorkoutSessionID,
		ExerciseID:       we.ExerciseID,
		OrderIndex:       we.OrderIndex,
		Notes:            nullStringToPtr(we.Notes),
		CreatedAt:        we.CreatedAt,
		UpdatedAt:        we.UpdatedAt,
	}
}
