package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/auth"
	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// testing helpers
func testingMarshalJSON(t *testing.T, v any) []byte {
	t.Helper()

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	return data
}
func testingExecuteRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}
func testingDecodeJSONResponse[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()

	var result T

	err := json.NewDecoder(rr.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func testingCreateJWT(t *testing.T, cfg *apiConfig, userID uuid.UUID) string {
	t.Helper()

	token, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	return token
}
func testingCreateAuthenticatedJSONRequest(method string, path string, body []byte, token string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	return req
}

// auth and setup handlers
func newTestAPIConfig(t *testing.T) *apiConfig {
	t.Helper()
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file was found!")
	}

	testDBURL := os.Getenv("TEST_DB_URL")
	if testDBURL == "" {
		// Intentionally hardcoded so the project runs immediately after cloning.
		// #nosec G101 -- development fallback
		testDBURL = "postgres://postgres:12345@localhost:5432/workout_tracker_test?sslmode=disable"
	}
	if !strings.Contains(testDBURL, "test") {
		t.Fatal("TEST_DB_URL does not appear to be a test database")
	}

	db, err := sql.Open("postgres", testDBURL)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	})
	return &apiConfig{
		db:        database.New(db),
		jwtSecret: "test-secret",
	}
}

func createTestUser(t *testing.T, cfg *apiConfig) database.User {
	t.Helper()
	hash := createTestHashPassword(t, "testpassword123")
	user, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:           uuid.New(),
		Email:        fmt.Sprintf("%v@test.com", uuid.New().String()),
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	t.Cleanup(func() {
		err := cfg.db.DeleteUserByID(context.Background(), user.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
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

// create test database entries
func createTestWorkoutSession(t *testing.T, cfg *apiConfig, userID uuid.UUID) database.WorkoutSession {
	t.Helper()
	session, err := cfg.db.CreateWorkoutSessions(context.Background(), database.CreateWorkoutSessionsParams{
		ID:          uuid.New(),
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
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("cleanup failed: %v", err)
		}
	})
	return session
}
func createTestWorkoutSessionWithDate(t *testing.T, cfg *apiConfig, userID uuid.UUID, date time.Time) database.WorkoutSession {
	t.Helper()
	session, err := cfg.db.CreateWorkoutSessions(context.Background(), database.CreateWorkoutSessionsParams{
		ID:          uuid.New(),
		UserID:      userID,
		WorkoutDate: date,
		Description: sql.NullString{
			String: "testing",
			Valid:  true,
		},
		Notes: sql.NullString{
			String: "test",
			Valid:  true,
		},
		CreatedAt: date,
		UpdatedAt: date,
	})
	if err != nil {
		t.Fatalf("failed to create test workout session: %v", err)
	}
	t.Cleanup(func() {
		err := cfg.db.DeleteWorkoutSession(context.Background(), database.DeleteWorkoutSessionParams{
			ID:     session.ID,
			UserID: session.UserID})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("cleanup failed: %v", err)
		}
	})
	return session
}
func createTestWorkoutSessionWithDescription(t *testing.T, cfg *apiConfig, userID uuid.UUID, desc string) database.WorkoutSession {
	t.Helper()
	session, err := cfg.db.CreateWorkoutSessions(context.Background(), database.CreateWorkoutSessionsParams{
		ID:          uuid.New(),
		UserID:      userID,
		WorkoutDate: time.Now().UTC(),
		Description: sql.NullString{
			String: desc,
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
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
		err := cfg.db.DeleteExerciseByID(context.Background(), database.DeleteExerciseByIDParams{
			ID:     exercise.ID,
			UserID: exercise.UserID,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("failed to cleanup exercise: %v", err)
		}
	})

	return exercise
}

func createTestWeightSet(t *testing.T, cfg *apiConfig, exercise database.WorkoutExercise) database.WeightsAndSet {
	t.Helper()
	weightSet, err := cfg.db.CreateWeightAndSet(context.Background(), database.CreateWeightAndSetParams{
		ID:                 uuid.New(),
		WorkoutExercisesID: exercise.ID,
		Weight:             25,
		IsKilograms:        true,
		SetNumber:          1,
		RepsTarget:         10,
		RepsActual:         10,
		DurationSeconds: sql.NullInt32{
			Int32: 30,
			Valid: true,
		},
		RestTimeSeconds: sql.NullInt32{
			Int32: 90,
			Valid: true,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to create test weight set: %v", err)
	}
	t.Cleanup(func() {
		err := cfg.db.DeleteWeightAndSet(context.Background(), database.DeleteWeightAndSetParams{
			ID:                 weightSet.ID,
			WorkoutExercisesID: weightSet.WorkoutExercisesID,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("failed to cleanup weight set: %v", err)
		}
	})

	return weightSet
}

func createTestCardioSet(t *testing.T, cfg *apiConfig, exercise database.WorkoutExercise) database.CardioAndSet {
	t.Helper()
	cardioSet, err := cfg.db.CreateCardioSet(context.Background(), database.CreateCardioSetParams{
		ID:                 uuid.New(),
		WorkoutExercisesID: exercise.ID,
		SetNumber:          5,
		Distance: sql.NullFloat64{
			Float64: 20,
			Valid:   true,
		},
		IsKilometers: true,
		DurationSeconds: sql.NullInt32{
			Int32: 30,
			Valid: true,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to create test cardio set: %v", err)
	}
	t.Cleanup(func() {
		err := cfg.db.DeleteCardioSet(context.Background(), database.DeleteCardioSetParams{
			ID:                 cardioSet.ID,
			WorkoutExercisesID: cardioSet.WorkoutExercisesID,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("failed to cleanup cardio set: %v", err)
		}
	})
	return cardioSet
}
