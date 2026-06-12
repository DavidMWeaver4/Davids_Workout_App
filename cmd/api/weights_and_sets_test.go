package main

import (
	"context"
	"errors"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestAuthorizeWeightAndSets_Success(t *testing.T) {
	x := newWeightSetFixture(t)
	gotWeightSet, err := x.cfg.authorizeWeightSet(x.ctx, x.weightSet.ID, x.user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotWeightSet.ID != x.weightSet.ID {
		t.Fatalf("expected exercise ID %v, got %v", x.weightSet.ID, gotWeightSet.ID)
	}
	if gotWeightSet.WorkoutExercisesID != x.weightSet.WorkoutExercisesID {
		t.Fatalf("expected exercise ID %v, got %v: ", x.weightSet.WorkoutExercisesID, gotWeightSet.WorkoutExercisesID)
	}
}

func TestAuthorizeWeightAndSets_Forbidden(t *testing.T) {
	x := newWeightSetFixture(t)
	otherUser := createTestUser(t, x.cfg)
	_, err := x.cfg.authorizeWeightSet(x.ctx, x.weightSet.ID, otherUser.ID)
	if err == nil {
		t.Fatal("expected ErrForbidden, got nil")
	}

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthorizeWeightAndSets_NotFound(t *testing.T) {
	x := newWeightSetFixture(t)
	_, err := x.cfg.authorizeWeightSet(x.ctx, uuid.New(), x.user.ID)
	if err == nil {
		t.Fatal("expected ErrNotFound, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

type weightSetFixture struct {
	ctx        context.Context
	cfg        *apiConfig
	user       database.User
	session    database.WorkoutSession
	exercise   database.Exercise
	woExercise database.WorkoutExercise
	weightSet  database.WeightsAndSet
}

func newWeightSetFixture(t *testing.T) weightSetFixture {
	t.Helper()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)
	weightSet := createTestWeightSet(t, cfg, woExercise)
	return weightSetFixture{
		ctx:        context.Background(),
		cfg:        cfg,
		user:       user,
		session:    session,
		exercise:   exercise,
		woExercise: woExercise,
		weightSet:  weightSet,
	}
}
