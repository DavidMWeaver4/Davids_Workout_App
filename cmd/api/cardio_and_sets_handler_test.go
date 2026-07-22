package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestHandlerCreateCardioSet_Success(t *testing.T) {
	x := testingCardioSetHandlerSetup(t)

	type requestBody struct {
		WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
		SetNumber          int32     `json:"set_number"`
		Distance           float64   `json:"distance,omitempty"`
		IsKilometers       bool      `json:"is_kilometers"`
		DurationSeconds    int32     `json:"duration_seconds,omitempty"`
	}
	body := requestBody{
		WorkoutExercisesID: x.workoutExercise.ID,
		SetNumber:          2,
		Distance:           12.1,
		IsKilometers:       true,
		DurationSeconds:    7200,
	}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPost, "/api/v1/cardio_sets", JSONBody, x.token)
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerCreateCardioSet), req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[cardioAndSets](t, rr)
	if response.Distance == nil {
		t.Fatalf("expected distance %.2f, got nil", body.Distance)
	}
	if *response.Distance != body.Distance {
		t.Fatalf("expected %.2f, got %.2f", body.Distance, *response.Distance)
	}
	if response.DurationSeconds == nil {
		t.Fatalf("expected DurationSeconds %d, got nil", body.DurationSeconds)
	}
	if *response.DurationSeconds != body.DurationSeconds {
		t.Fatalf("expected %d, got %d", body.DurationSeconds, *response.DurationSeconds)
	}
	if response.WorkoutExercisesID != x.workoutExercise.ID {
		t.Fatalf("got wrong session id, expected %v, got %v", x.workoutExercise.ID, response.WorkoutExercisesID)
	}
	if response.SetNumber != body.SetNumber {
		t.Fatalf("expected set number %d got %d", body.SetNumber, response.SetNumber)
	}
	dbCardioSet, err := x.cfg.db.GetCardioSetFromID(context.Background(), response.ID)
	t.Cleanup(func() {
		x.cfg.db.DeleteCardioSet(context.Background(), database.DeleteCardioSetParams{
			ID:                 response.ID,
			WorkoutExercisesID: x.workoutExercise.ID,
		})
	})
	if err != nil {
		t.Fatal("cardio set was not inserted into database")
	}
	if dbCardioSet.Distance.Float64 != body.Distance {
		t.Fatal("database values does not match entry data")
	}
}

/*
	mux.HandleFunc("GET /api/v1/workout_exercises/{id}/cardio_sets", apiCfg.handlerGetAllCardioFromExercise)
	mux.HandleFunc("DELETE /api/v1/cardio_sets/{id}", apiCfg.handlerDeleteCardioSet)
	mux.HandleFunc("PUT /api/v1/cardio_sets/{id}", apiCfg.handlerUpdateCardioSet)
	mux.HandleFunc("GET /api/v1/cardio_sets/{id}/distance", apiCfg.handlerGetSetDistance)
	mux.HandleFunc("GET /api/v1/cardio_sets/{id}/duration", apiCfg.handlerGetCardioSetDuration)
	mux.HandleFunc("GET /api/v1/workout_exercises/{id}/distance", apiCfg.handlerGetExerciseDistance)
	mux.HandleFunc("GET /api/v1/workout_exercises/{id}/duration", apiCfg.handlerGetCardioExerciseDuration)
	mux.HandleFunc("GET /api/v1/workout_sessions/{id}/cardio_sets", apiCfg.handlerGetAllCardioSessionSets)
	mux.HandleFunc("GET /api/v1/workout_sessions/{id}/distance", apiCfg.handlerGetAllSessionDistance)
	mux.HandleFunc("GET /api/v1/workout_sessions/{id}/duration", apiCfg.handlerGetAllCardioSessionDuration)

*/
//
//helpers
//
type CardioSetHandlerFixture struct {
	ctx             context.Context
	cfg             *apiConfig
	user            database.User
	token           string
	workoutSession  database.WorkoutSession
	exercise        database.Exercise
	workoutExercise database.WorkoutExercise
}

func testingCardioSetHandlerSetup(t *testing.T) CardioSetHandlerFixture {
	t.Helper()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	workoutSession := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	workoutExercise := createTestWorkoutExercise(t, cfg, workoutSession, exercise.ID)
	return CardioSetHandlerFixture{
		ctx:             context.Background(),
		cfg:             cfg,
		user:            user,
		token:           token,
		workoutSession:  workoutSession,
		exercise:        exercise,
		workoutExercise: workoutExercise,
	}
}
