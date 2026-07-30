// Handler Testing
package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestHandlerCreateExercise_Success(t *testing.T) {
	x := testingExerciseHandlerSetup(t)
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

	JSONBody := testingMarshalJSON(t, body)

	req := testingCreateAuthenticatedJSONRequest(http.MethodPost, "/api/v1/exercises", JSONBody, x.token)
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerCreateExercises), req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[exercise](t, rr)

	if response.ExerciseName != body.ExerciseName {
		t.Errorf("expected %s, got %s", body.ExerciseName, response.ExerciseName)
	}

	if len(response.TargetMuscles) != 1 {
		t.Errorf("expected 1 target muscle")
	}

	dbExercise, err := x.cfg.db.GetExerciseFromID(context.Background(), response.ID)
	t.Cleanup(func() {
		x.cfg.db.DeleteExerciseByID(context.Background(), database.DeleteExerciseByIDParams{
			ID: response.ID,
			UserID: uuid.NullUUID{
				UUID:  x.user.ID,
				Valid: true,
			},
		},
		)
	})
	if err != nil {
		t.Fatalf("GetExerciseFromID failed: %v", err)
	}

	if dbExercise.ExerciseName != body.ExerciseName {
		t.Errorf("database value does not match entry data")
	}
}

func TestHandlerDeleteExercise_Success(t *testing.T) {
	x := testingExerciseHandlerSetup(t)

	exerciseToDelete := createTestExercise(t, x.cfg, x.user.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodDelete, "/api/v1/exercises/"+exerciseToDelete.ID.String(), nil, x.token)
	req.SetPathValue("id", exerciseToDelete.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerDeleteExercisesByID), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	_, err := x.cfg.db.GetExerciseFromID(context.Background(), exerciseToDelete.ID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHandlerGetExerciseFromID_Success(t *testing.T) {
	x := testingExerciseHandlerSetup(t)

	exerciseToFind := createTestExercise(t, x.cfg, x.user.ID)
	exerciseID := exerciseToFind.ID

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/exercises/"+exerciseID.String(), nil, x.token)
	req.SetPathValue("id", exerciseID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetExerciseFromID), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[exercise](t, rr)
	if response.ID != exerciseToFind.ID {
		t.Errorf("expected exercise ID %v, to match %v", response.ID, exerciseToFind.ID)
	}
}

func TestHandlerSearchExercises_Success(t *testing.T) {
	x := testingExerciseHandlerSetup(t)

	testExercise := createTestExercise(t, x.cfg, x.user.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/exercises/search?muscle=chest&equipment=barbell&difficulty=beginner", nil, x.token)
	req.Header.Set("Authorization", "Bearer "+x.token)
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerSearchExercises), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[[]exercise](t, rr)

	if len(response) == 0 {
		t.Errorf("expected at least one exercise")
	}

	found := false
	for _, ex := range response {
		if ex.ID == testExercise.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("created exercise was not returned")
	}
}
func TestHandlerUpdateExercise_Success(t *testing.T) {
	x := testingExerciseHandlerSetup(t)

	testExercise := createTestExercise(t, x.cfg, x.user.ID)
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

	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/exercises/"+testExercise.ID.String(), jsonBody, x.token)
	req.SetPathValue("id", testExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateExercise), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[exercise](t, rr)

	dbExercise, err := x.cfg.db.GetExerciseFromID(context.Background(), response.ID)

	if err != nil {
		t.Fatal("exercise was not updated in database")
	}

	if dbExercise.ExerciseName != body.ExerciseName {
		t.Errorf("database exercise name does not match entry data")
	}
	if dbExercise.Equipment.String != body.Equipment {
		t.Errorf("database equipment does not match entry data")
	}
	if dbExercise.Description.String != body.Description {
		t.Errorf("database description does not match entry data")
	}
	if dbExercise.DifficultyLevel.String != body.DifficultyLevel {
		t.Errorf("database difficulty level does not match entry data")
	}
}

type ExerciseHandlerFixture struct {
	ctx   context.Context
	cfg   *apiConfig
	user  database.User
	token string
}

func testingExerciseHandlerSetup(t *testing.T) ExerciseHandlerFixture {
	t.Helper()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	return ExerciseHandlerFixture{
		ctx:   context.Background(),
		cfg:   cfg,
		user:  user,
		token: token,
	}
}
