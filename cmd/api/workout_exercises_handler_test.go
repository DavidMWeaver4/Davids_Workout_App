package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestHandlerCreateWorkoutExercise_Success(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	type requestBody struct {
		ExerciseName     string    `json:"exercise_name"`
		WorkoutSessionID uuid.UUID `json:"workout_session_id"`
		OrderIndex       int32     `json:"order_index"`
		Notes            string    `json:"notes"`
	}

	body := requestBody{
		ExerciseName:     x.exercise.ExerciseName,
		WorkoutSessionID: x.workoutSession.ID,
		OrderIndex:       0,
		Notes:            "testing note",
	}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPost, "/api/v1/workout_exercises", JSONBody, x.token)
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerCreateWorkoutExercise), req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[workoutExercise](t, rr)
	if response.Notes == nil || *response.Notes != body.Notes {
		t.Fatalf("expected %s, got %v", body.Notes, response.Notes)
	}
	if response.WorkoutSessionID != x.workoutSession.ID {
		t.Fatalf("got wrong session ID, expected %v, got %v", x.workoutSession.ID, response.WorkoutSessionID)
	}
	if response.OrderIndex != 0 {
		t.Fatalf("expect order 0, got %v", response.OrderIndex)
	}
	dbWorkoutExercise, err := x.cfg.db.GetWorkoutExerciseFromID(context.Background(), response.ID)
	t.Cleanup(func() {
		x.cfg.db.DeleteWorkoutExercises(context.Background(), database.DeleteWorkoutExercisesParams{
			ID:               response.ID,
			WorkoutSessionID: x.workoutSession.ID,
		})
	})
	if err != nil {
		t.Fatal("workout exercise was not inserted into database")
	}
	if dbWorkoutExercise.Notes.String != body.Notes {
		t.Fatal("database value does not match entry data")
	}
}
func TestHandlerGetWorkoutExercise_Success(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	exerciseToFind := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_exercises/"+exerciseToFind.ID.String(), nil, x.token)
	req.SetPathValue("id", exerciseToFind.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetWorkoutExercises), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[workoutExercise](t, rr)

	if response.ID != exerciseToFind.ID {
		t.Fatalf("expected session ID %v, to match %v", response.ID, exerciseToFind.ID)
	}
}

func TestHandlerGetWorkoutExercisesInSession(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	exercise2 := createTestExercise(t, x.cfg, x.user.ID)
	workoutExercise1 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	workoutExercise2 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise2.ID)
	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/"+x.workoutSession.ID.String()+"/exercises", nil, x.token)
	req.SetPathValue("id", x.workoutSession.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetWorkoutExercisesInSession), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[[]workoutExercise](t, rr)
	if len(response) != 2 {
		t.Fatalf("expected 2 workout exercises, got %d", len(response))
	}
	expectedIDs := map[uuid.UUID]bool{workoutExercise1.ID: false, workoutExercise2.ID: false}
	for _, item := range response {
		if _, ok := expectedIDs[item.ID]; ok {
			expectedIDs[item.ID] = true
		}
	}
	if expectedIDs[workoutExercise1.ID] != true || expectedIDs[workoutExercise2.ID] != true {
		t.Fatal("did not find correct workout exercises")
	}
}

func TestHandlerDeleteWorkoutExercises_Success(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	req := testingCreateAuthenticatedJSONRequest(http.MethodDelete, "/api/v1/workout_exercises/"+workoutExercise.ID.String(), nil, x.token)
	req.SetPathValue("id", workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerDeleteWorkoutExercises), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	_, err := x.cfg.db.GetWorkoutExerciseFromID(context.Background(), workoutExercise.ID)
	if err == nil {
		t.Fatal("expected workout exercise to be deleted")
	}
}

func TestHandlerGetNumOfWorkoutsInSession_Success(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	_ = createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	secondEx := createTestExercise(t, x.cfg, x.user.ID)
	_ = createTestWorkoutExercise(t, x.cfg, x.workoutSession, secondEx.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/"+x.workoutSession.ID.String()+"/exercises/count", nil, x.token)
	req.SetPathValue("id", x.workoutSession.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetNumOfWorkoutsInSession), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[countResponse](t, rr)
	if response.Count != 2 {
		t.Fatalf("got unexpected length: %d", response.Count)
	}
}

func TestHandlerUpdateWorkoutExercise_Success(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	newExercise := createTestExercise(t, x.cfg, x.user.ID)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	type requestBody struct {
		ExerciseName string  `json:"exercise_name"`
		OrderIndex   int32   `json:"order_index"`
		Notes        *string `json:"notes,omitempty"`
	}
	newString := "updated test note"
	body := requestBody{
		ExerciseName: newExercise.ExerciseName,
		OrderIndex:   2,
		Notes:        &newString,
	}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/workout_exercises/"+workoutExercise.ID.String(), JSONBody, x.token)
	req.SetPathValue("id", workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateWorkoutExercise), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	dbWOExercise, err := x.cfg.db.GetWorkoutExerciseFromID(context.Background(), workoutExercise.ID)
	if err != nil {
		t.Fatal("session was not in database")
	}
	if dbWOExercise.Notes.String != *body.Notes ||
		dbWOExercise.OrderIndex != body.OrderIndex ||
		dbWOExercise.ExerciseID != newExercise.ID {
		t.Fatal("database value does not match entry data")
	}
}

/*
 *
 *
 * helpers
 *
 *
 */
func testingWorkoutExerciseHandlerSetup(t *testing.T) WorkoutExerciseHandlerFixture {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	workoutSession := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	return WorkoutExerciseHandlerFixture{
		ctx:            context.Background(),
		cfg:            cfg,
		user:           user,
		token:          token,
		workoutSession: workoutSession,
		exercise:       exercise,
	}
}

type WorkoutExerciseHandlerFixture struct {
	ctx            context.Context
	cfg            *apiConfig
	user           database.User
	token          string
	workoutSession database.WorkoutSession
	exercise       database.Exercise
}
