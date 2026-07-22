package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestHandlerCreateWorkoutExercise_Success(t *testing.T) {
	//POST /api/v1/workout_exercises", apiCfg.handlerCreateWorkoutExercise
	x := testingWorkoutExerciseHandlerSetup(t)
	type requestBody struct {
		ExerciseName     string    `json:"exercise_name"`
		WorkoutSessionID uuid.UUID `json:"workout_session_id"`
		OrderIndex       int32     `json:"order_index"`
		Notes            string    `json:"notes"`
	}
	exercise := createTestExercise(t, x.cfg, x.user.ID)

	body := requestBody{
		ExerciseName:     exercise.ExerciseName,
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

	//GET /api/v1/workout_exercises/{id}", apiCfg.handlerGetWorkoutExercises)
	x := testingWorkoutExerciseHandlerSetup(t)
	exercise := createTestExercise(t, x.cfg, x.user.ID)

	exerciseToFind := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise.ID)

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
	t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/{session_id}/exercises", apiCfg.handlerGetWorkoutExercisesInSession)
}
func TestHandlerDeleteWorkoutExercises_Success(t *testing.T) {
	t.Skip()
	//mux.HandleFunc("DELETE /api/v1/workout_exercises/{id}", apiCfg.handlerDeleteWorkoutExercises)

}
func TestHandlerGetNumOfWorkoutsInSession_Success(t *testing.T) {
	t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/{session_id}/exercises/count", apiCfg.handlerGetNumOfWorkoutsInSession)
}
func TestHandlerUpdateWorkoutExerciseOrder_Success(t *testing.T) {
	t.Skip()
	//mux.HandleFunc("PUT /api/v1/workout_exercises/{id}", apiCfg.handlerUpdateWorkoutExerciseOrder)
}

func testingWorkoutExerciseHandlerSetup(t *testing.T) WorkoutExerciseHandlerFixture {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	workoutSession := createTestWorkoutSession(t, cfg, user.ID)
	return WorkoutExerciseHandlerFixture{
		ctx:            context.Background(),
		cfg:            cfg,
		user:           user,
		token:          token,
		workoutSession: workoutSession,
	}
}

type WorkoutExerciseHandlerFixture struct {
	ctx            context.Context
	cfg            *apiConfig
	user           database.User
	token          string
	workoutSession database.WorkoutSession
}
