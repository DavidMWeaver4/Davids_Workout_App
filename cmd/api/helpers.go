package main

import (
	"context"
	"database/sql"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) authorizeWorkoutExercise(ctx context.Context, workoutExerciseID uuid.UUID, userID uuid.UUID) (database.WorkoutExercise, error) {
	workExercise, err := cfg.db.GetWorkoutExerciseFromID(ctx, workoutExerciseID)
	if err == sql.ErrNoRows {
		return database.WorkoutExercise{}, ErrNotFound
	}
	if err != nil {
		return database.WorkoutExercise{}, err
	}

	workSess, err := cfg.db.GetWorkoutSessionByID(ctx, workExercise.WorkoutSessionID)
	if err == sql.ErrNoRows {
		return database.WorkoutExercise{}, ErrNotFound
	}
	if err != nil {
		return database.WorkoutExercise{}, err
	}

	if workSess.UserID != userID {
		return database.WorkoutExercise{}, ErrForbidden
	}

	return workExercise, nil
}
func (cfg *apiConfig) authorizeWorkoutSession(ctx context.Context, workoutSessionID uuid.UUID, userID uuid.UUID) (database.WorkoutSession, error) {

	workSess, err := cfg.db.GetWorkoutSessionByID(
		ctx,
		workoutSessionID,
	)

	if err == sql.ErrNoRows {
		return database.WorkoutSession{}, ErrNotFound
	}

	if err != nil {
		return database.WorkoutSession{}, err
	}

	if workSess.UserID != userID {
		return database.WorkoutSession{}, ErrForbidden
	}

	return workSess, nil
}
func (cfg *apiConfig) authorizeWeightSet(ctx context.Context, weightSetID uuid.UUID, userID uuid.UUID) (database.WeightsAndSet, error) {

	weightSet, err := cfg.db.GetWeightAndSetFromID(ctx, weightSetID)
	if err == sql.ErrNoRows {
		return database.WeightsAndSet{}, ErrNotFound
	}
	if err != nil {
		return database.WeightsAndSet{}, err
	}

	_, err = cfg.authorizeWorkoutExercise(ctx, weightSet.WorkoutExercisesID, userID)
	if err != nil {
		return database.WeightsAndSet{}, err
	}

	return weightSet, nil
}
func (cfg *apiConfig) authorizeCardioSet(ctx context.Context, cardioSetID uuid.UUID, userID uuid.UUID) (database.CardioAndSet, error) {

	cardioSet, err := cfg.db.GetCardioAndSetFromID(ctx, cardioSetID)
	if err == sql.ErrNoRows {
		return database.CardioAndSet{}, ErrNotFound
	}
	if err != nil {
		return database.CardioAndSet{}, err
	}

	_, err = cfg.authorizeWorkoutExercise(ctx, cardioSet.WorkoutExercisesID, userID)
	if err != nil {
		return database.CardioAndSet{}, err
	}

	return cardioSet, nil
}

func (cfg *apiConfig) authorizeExercise(ctx context.Context, exerciseID uuid.UUID, userID uuid.UUID) (database.Exercise, error) {
	exercise, err := cfg.db.GetExerciseFromID(ctx, exerciseID)
	if err == sql.ErrNoRows {
		return database.Exercise{}, ErrNotFound
	}
	if err != nil {
		return database.Exercise{}, err
	}
	if !exercise.UserID.Valid || exercise.UserID.UUID != userID {
		return database.Exercise{}, ErrForbidden
	}
	return exercise, nil
}
