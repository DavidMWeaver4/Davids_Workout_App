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
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	WOSession := createTestWorkoutSession(t, cfg, user.ID)
	type requestBody struct {
		ExerciseName     string    `json:"exercise_name"`
		WorkoutSessionID uuid.UUID `json:"workout_session_id"`
		OrderIndex       int32     `json:"order_index"`
		Notes            string    `json:"notes"`
	}
	exercise := createTestExercise(t, cfg, user.ID)

	body := requestBody{
		ExerciseName:     exercise.ExerciseName,
		WorkoutSessionID: WOSession.ID,
		OrderIndex:       0,
		Notes:            "testing note",
	}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPost, "/api/v1/workout_exercises", JSONBody, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerCreateWorkoutExercise), req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[workoutExercise](t, rr)
	if response.Notes == nil || *response.Notes != body.Notes {
		t.Fatalf("expected %s, got %v", body.Notes, response.Notes)
	}
	if response.WorkoutSessionID != WOSession.ID {
		t.Fatalf("got wrong session ID, expected %v, got %v", WOSession.ID, response.WorkoutSessionID)
	}
	if response.OrderIndex != 0 {
		t.Fatalf("expect order 0, got %v", response.OrderIndex)
	}
	dbWorkoutExercise, err := cfg.db.GetWorkoutExerciseFromID(context.Background(), response.ID)
	t.Cleanup(func() {
		cfg.db.DeleteWorkoutExercises(context.Background(), database.DeleteWorkoutExercisesParams{
			ID:               response.ID,
			WorkoutSessionID: WOSession.ID,
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
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	WOSession := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)

	exerciseToFind := createTestWorkoutExercise(t, cfg, WOSession, exercise.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_exercises/"+exerciseToFind.ID.String(), nil, token)
	req.SetPathValue("id", exerciseToFind.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetWorkoutExercises), req)
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
