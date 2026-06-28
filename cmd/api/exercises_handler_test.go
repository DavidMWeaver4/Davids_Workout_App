// Handler Testing
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/auth"
	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestHandlerCreateExercise_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	type requestBody struct {
		ExerciseName    string   `json:"exercise_name"`
		TargetMuscles   []string `json:"target_muscles"`
		Equipment       string   `json:"equipment"`
		DifficultyLevel string   `json:"difficulty_level"`
		Description     string   `json:"description"`
	}

	body := requestBody{
		ExerciseName:    "Bench Press",
		TargetMuscles:   []string{"chest"},
		Equipment:       "barbell",
		DifficultyLevel: "beginner",
		Description:     "compound movement",
	}

	jsonBody := testingMarshalJSON(t, body)

	req := testingCreateAuthenticatedJSONRequest(http.MethodPost, "/api/v1/exercises", jsonBody, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerCreateExercises), req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[exercise](t, rr)

	if response.ExerciseName != body.ExerciseName {
		t.Fatalf("expected %s, got %s", body.ExerciseName, response.ExerciseName)
	}

	if len(response.TargetMuscles) != 1 {
		t.Fatal("expected 1 target muscle")
	}

	dbExercise, err := cfg.db.GetExerciseFromID(context.Background(), response.ID)
	t.Cleanup(func() {
		cfg.db.DeleteExerciseByID(context.Background(), database.DeleteExerciseByIDParams{
			ID:     response.ID,
			UserID: dbExercise.UserID,
		})
	})

	if err != nil {
		t.Fatal("exercise was not inserted into database")
	}

	if dbExercise.ExerciseName != body.ExerciseName {
		t.Fatal("database value does not match")
	}
}

func TestHandlerDeleteExercise_Success(t *testing.T) {

	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	exerciseToDelete := createTestExercise(t, cfg, user.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodDelete, "/api/v1/exercises/"+exerciseToDelete.ID.String(), nil, token)
	req.SetPathValue("id", exerciseToDelete.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerDeleteExercisesByID), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	_, err := cfg.db.GetExerciseFromID(context.Background(), exerciseToDelete.ID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHandlerGetExerciseFromID_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	exerciseToFind := createTestExercise(t, cfg, user.ID)
	exerciseID := exerciseToFind.ID

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/exercises/"+exerciseID.String(), nil, token)
	req.SetPathValue("id", exerciseID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetExerciseFromID), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[exercise](t, rr)
	if response.ID != exerciseToFind.ID {
		t.Fatalf("expected exercise ID %v, to match %v", response.ID, exerciseToFind.ID)
	}
}

func TestHandlerSearchExercises_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	testExercise := createTestExercise(t, cfg, user.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/exercises/search?muscle=chest&equipment=barbell&difficulty=beginner", nil, token)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerSearchExercises), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[[]exercise](t, rr)

	if len(response) == 0 {
		t.Fatal("expected at least one exercise")
	}

	found := false
	for _, ex := range response {
		if ex.ID == testExercise.ID {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("created exercise was not returned")
	}
}
func TestHandlerUpdateExercise_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	testExercise := createTestExercise(t, cfg, user.ID)
	type requestBody struct {
		ExerciseName    string   `json:"exercise_name"`
		TargetMuscles   []string `json:"target_muscles"`
		Equipment       string   `json:"equipment"`
		DifficultyLevel string   `json:"difficulty_level"`
		Description     string   `json:"description"`
	}

	body := requestBody{
		ExerciseName:    "Bench Press",
		TargetMuscles:   []string{"chest"},
		Equipment:       "barbell",
		DifficultyLevel: "beginner",
		Description:     "compound movement",
	}

	jsonBody := testingMarshalJSON(t, body)

	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/exercises/"+testExercise.ID.String(), jsonBody, token)
	req.SetPathValue("id", testExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerUpdateExercise), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[exercise](t, rr)

	dbExercise, err := cfg.db.GetExerciseFromID(context.Background(), response.ID)

	if err != nil {
		t.Fatal("exercise was not updated in database")
	}

	if dbExercise.ExerciseName != body.ExerciseName ||
		dbExercise.Equipment.String != body.Equipment ||
		dbExercise.Description.String != body.Description ||
		dbExercise.DifficultyLevel.String != body.DifficultyLevel {
		t.Fatal("database value does not match")
	}
}

// helper funcs
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
