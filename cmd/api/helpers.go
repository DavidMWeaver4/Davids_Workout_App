package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

// add respondwitherror to this helper function
func (cfg *apiConfig) authorizeWorkoutExercise(ctx context.Context, workoutExerciseID uuid.UUID, userID uuid.UUID) (database.WorkoutExercise, error) {

	workExercise, err := cfg.db.GetWorkoutExerciseFromID(ctx, workoutExerciseID)
	if err == sql.ErrNoRows {
		return database.WorkoutExercise{}, err
	}
	if err != nil {
		return database.WorkoutExercise{}, err
	}

	workSess, err := cfg.db.GetWorkoutSessionByID(ctx, workExercise.WorkoutSessionID)
	if err == sql.ErrNoRows {
		return database.WorkoutExercise{}, err
	}
	if err != nil {
		return database.WorkoutExercise{}, err
	}

	if workSess.UserID != userID {
		return database.WorkoutExercise{}, errors.New("Forbidden")
	}

	return workExercise, nil
}
