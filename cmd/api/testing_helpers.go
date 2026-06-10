package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/auth"
	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func newTestAPIConfig(t *testing.T) *apiConfig {
	t.Helper()
	/*
		err := godotenv.Load(".env")
		if err != nil {
			t.Fatalf("failed to load .env: %v", err)
		}
	*/
	db, err := sql.Open("postgres", os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	if !strings.Contains(os.Getenv("TEST_DB_URL"), "test") {
		t.Fatal("TEST_DB_URL does not appear to be a test database")
	}
	t.Cleanup(func() { db.Close() })
	return &apiConfig{db: database.New(db)}
}

func createTestUser(t *testing.T, cfg *apiConfig) database.User {
	t.Helper()
	hash := createTestHashPassword(t, "testpassword123")
	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:           uuid.New(),
		Email:        fmt.Sprintf("%v@test.com", time.Now().UTC()),
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	t.Cleanup(func() {
		err := cfg.db.DeleteUserByID(context.Background(), user.ID)
		if err != nil {
			t.Errorf("failed to cleanup user: %v", err)
		}
	})
	return user
}

func createTestHashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return hash
}

func createTestWorkoutSession(t *testing.T, cfg *apiConfig, userID uuid.UUID) database.WorkoutSession {
	t.Helper()
	session, err := cfg.db.CreateWorkoutSessions(context.Background(), database.CreateWorkoutSessionsParams{
		UserID:      userID,
		WorkoutDate: time.Now().UTC(),
		Description: sql.NullString{
			String: "testing",
			Valid:  true,
		},
		Notes: sql.NullString{
			String: "test",
			Valid:  true,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to create test workout session: %v", err)
	}
	t.Cleanup(func() {
		err := cfg.db.DeleteWorkoutSession(context.Background(), database.DeleteWorkoutSessionParams{
			ID:     session.ID,
			UserID: session.UserID})
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})
	return session
}

func createTestWorkoutExercise(t *testing.T, cfg *apiConfig, session database.WorkoutSession, exerciseID uuid.UUID) database.WorkoutExercise {
	t.Helper()
	woExercise, err := cfg.db.CreateWorkoutExercises(context.Background(), database.CreateWorkoutExercisesParams{
		ID:               uuid.New(),
		WorkoutSessionID: session.ID,
		ExerciseID:       exerciseID,
		OrderIndex:       1,
		Notes:            sql.NullString{},
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to create test workout exercise: %v", err)
	}
	t.Cleanup(func() {
		err := cfg.db.DeleteWorkoutExercises(context.Background(), database.DeleteWorkoutExercisesParams{
			ID:               woExercise.ID,
			WorkoutSessionID: woExercise.WorkoutSessionID,
		})
		if err != nil {
			t.Errorf("cleanup failed: %v", err)
		}
	})
	return woExercise
}

func createTestExercise(t *testing.T, cfg *apiConfig, userID uuid.UUID) database.Exercise {
	t.Helper()
	exercise, err := cfg.db.CreateExercises(context.Background(), database.CreateExercisesParams{
		ID: uuid.New(),
		UserID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
		ExerciseName:  fmt.Sprintf("test-exercise-%s", uuid.New()),
		TargetMuscles: []string{"chest"},
		Equipment: sql.NullString{
			String: "barbell",
			Valid:  true,
		},
		DifficultyLevel: sql.NullString{
			String: "beginner",
			Valid:  true,
		},
		Description: sql.NullString{
			String: "test description",
			Valid:  true,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to create test exercise: %v", err)
	}

	t.Cleanup(func() {
		err := cfg.db.DeleteExerciseByID(
			context.Background(),
			database.DeleteExerciseByIDParams{
				ID:     exercise.ID,
				UserID: exercise.UserID,
			},
		)
		if err != nil {
			t.Errorf("failed to cleanup exercise: %v", err)
		}
	})

	return exercise
}
