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
		t.Errorf("expected distance %.2f, got nil", body.Distance)
	}
	if *response.Distance != body.Distance {
		t.Errorf("expected %.2f, got %.2f", body.Distance, *response.Distance)
	}
	if response.DurationSeconds == nil {
		t.Errorf("expected DurationSeconds %d, got nil", body.DurationSeconds)
	}
	if *response.DurationSeconds != body.DurationSeconds {
		t.Errorf("expected %d, got %d", body.DurationSeconds, *response.DurationSeconds)
	}
	if response.WorkoutExercisesID != x.workoutExercise.ID {
		t.Errorf("got wrong session id, expected %v, got %v", x.workoutExercise.ID, response.WorkoutExercisesID)
	}
	if response.SetNumber != body.SetNumber {
		t.Errorf("expected set number %d got %d", body.SetNumber, response.SetNumber)
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
	if !dbCardioSet.Distance.Valid || dbCardioSet.Distance.Float64 != body.Distance {
		t.Errorf("database values does not match entry data")
	}
}

func TestHandlerGetAllCardioFromExercise_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_exercises/{id}/cardio_sets", apiCfg.handlerGetAllCardioFromExercise)
	x := testingCardioSetHandlerSetup(t)
	set1 := createTestCardioSet(t, x.cfg, x.workoutExercise)
	set2 := createTestCardioSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_exercises/"+x.workoutExercise.ID.String()+"/cardio_sets", nil, x.token)
	req.SetPathValue("id", x.workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetAllCardioFromExercise), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[[]cardioAndSets](t, rr)
	if len(response) != 2 {
		t.Errorf("expected length 2, got %d", len(response))
	}
	expected := map[uuid.UUID]bool{
		set1.ID: false,
		set2.ID: false,
	}

	for _, s := range response {
		if _, ok := expected[s.ID]; ok {
			expected[s.ID] = true
		}
	}

	if !expected[set1.ID] || !expected[set2.ID] {
		t.Errorf("did not receive expected cardio sets")
	}

}
func TestHandlerDeleteCardioSet_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("DELETE /api/v1/cardio_sets/{id}", apiCfg.handlerDeleteCardioSet)
	x := testingCardioSetHandlerSetup(t)
	set := createTestCardioSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodDelete, "/api/v1/cardio_sets/"+set.ID.String(), nil, x.token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerDeleteCardioSet), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	_, err := x.cfg.db.GetCardioSetFromID(context.Background(), set.ID)
	if err == nil {
		t.Fatalf("expected cardio set to be deleted")
	}
}
func TestHandlerUpdateCardioSet_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("PUT /api/v1/cardio_sets/{id}", apiCfg.handlerUpdateCardioSet)
	x := testingCardioSetHandlerSetup(t)
	set := createTestCardioSet(t, x.cfg, x.workoutExercise)

	type requestBody struct {
		SetNumber       int32   `json:"set_number"`
		Distance        float64 `json:"distance,omitempty"`
		IsKilometers    bool    `json:"is_kilometers"`
		DurationSeconds int32   `json:"duration_seconds,omitempty"`
	}
	body := requestBody{
		SetNumber:       6,
		Distance:        42.195,
		IsKilometers:    true,
		DurationSeconds: 7235,
	}
	JSONBody := testingMarshalJSON(t, body)

	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/cardio_sets/"+set.ID.String(), JSONBody, x.token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateCardioSet), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[cardioAndSets](t, rr)
	if response.Distance == nil || *response.Distance != body.Distance {
		t.Errorf("expected distance %.2f, got %v", body.Distance, response.Distance)
	}
	dbSet, err := x.cfg.db.GetCardioSetFromID(context.Background(), set.ID)
	if err != nil {
		t.Fatalf("failed retrieving updated set")
	}

	if !dbSet.Distance.Valid || dbSet.Distance.Float64 != body.Distance {
		t.Errorf("database distance does not match request body")
	}
	if dbSet.SetNumber != body.SetNumber {
		t.Errorf("database SetNUmber does not match request body")
	}
	if dbSet.IsKilometers != body.IsKilometers {
		t.Errorf("database IsKilometers does not match request body")
	}
	if !dbSet.DurationSeconds.Valid || dbSet.DurationSeconds.Int32 != body.DurationSeconds {
		t.Errorf("database DurationSeconds does not match request body")
	}
}
func TestHandlerGetSetDistance_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("GET /api/v1/cardio_sets/{id}/distance", apiCfg.handlerGetSetDistance)
	x := testingCardioSetHandlerSetup(t)
	cardioSet := createTestCardioSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/cardio_sets/"+cardioSet.ID.String(), nil, x.token)
	req.SetPathValue("id", cardioSet.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetSetDistance), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[distanceResponse](t, rr)
	expectedDistance := cardioSet.Distance.Float64
	if response.Distance != expectedDistance {
		t.Errorf("expected %.2f, got %.2f", expectedDistance, response.Distance)
	}
}
func TestHandlerGetCardioSetDuration_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("GET /api/v1/cardio_sets/{id}/duration", apiCfg.handlerGetCardioSetDuration)
	x := testingCardioSetHandlerSetup(t)
	set := createTestCardioSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/cardio_sets/"+set.ID.String()+"/duration", nil, x.token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetCardioSetDuration), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[cardioDurationResponse](t, rr)
	expected := set.DurationSeconds.Int32

	if response.TotalSeconds != int64(expected) {
		t.Errorf("expected duration %d got %d", expected, response.TotalSeconds)
	}
}
func TestHandlerGetExerciseDistance_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_exercises/{id}/distance", apiCfg.handlerGetExerciseDistance)
	x := testingCardioSetHandlerSetup(t)

	set1 := createTestCardioSet(t, x.cfg, x.workoutExercise)
	set2 := createTestCardioSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_exercises/"+x.workoutExercise.ID.String()+"/distance", nil, x.token)
	req.SetPathValue("id", x.workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetExerciseDistance), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[distanceResponse](t, rr)
	expected := set1.Distance.Float64 + set2.Distance.Float64
	if response.Distance != expected {
		t.Errorf("expected distance %.2f, got %.2f", expected, response.Distance)
	}
}
func TestHandlerGetCardioExerciseDuration_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_exercises/{id}/duration", apiCfg.handlerGetCardioExerciseDuration)
	x := testingCardioSetHandlerSetup(t)
	set1 := createTestCardioSet(t, x.cfg, x.workoutExercise)
	set2 := createTestCardioSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_exercises/"+x.workoutExercise.ID.String()+"/duration", nil, x.token)
	req.SetPathValue("id", x.workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetCardioExerciseDuration), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[cardioDurationResponse](t, rr)
	expected := set1.DurationSeconds.Int32 + set2.DurationSeconds.Int32
	if response.TotalSeconds != int64(expected) {
		t.Errorf("expected duration %d, got %d", expected, response.TotalSeconds)
	}

}
func TestHandlerGetAllCardioSessionSets_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/{id}/cardio_sets", apiCfg.handlerGetAllCardioSessionSets)
	x := testingCardioSetHandlerSetup(t)
	exercise2 := createTestExercise(t, x.cfg, x.user.ID)
	workoutExercise2 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise2.ID)

	set1 := createTestCardioSet(t, x.cfg, x.workoutExercise)
	set2 := createTestCardioSet(t, x.cfg, workoutExercise2)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/"+x.workoutSession.ID.String()+"/cardio_sets", nil, x.token)
	req.SetPathValue("id", x.workoutSession.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetAllCardioSessionSets), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[[]cardioAndSets](t, rr)
	if len(response) != 2 {
		t.Errorf("expected 2 sets got %d", len(response))
	}
	expectedID := map[uuid.UUID]bool{set1.ID: false, set2.ID: false}
	for _, item := range response {
		if _, ok := expectedID[item.ID]; ok {
			expectedID[item.ID] = true
		}
	}
	if !expectedID[set1.ID] || !expectedID[set2.ID] {
		t.Errorf("did not find correct sets")
	}
}
func TestHandlerGetAllSessionDistance_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/{id}/distance", apiCfg.handlerGetAllSessionDistance)
	x := testingCardioSetHandlerSetup(t)
	exercise1 := createTestExercise(t, x.cfg, x.user.ID)
	exercise2 := createTestExercise(t, x.cfg, x.user.ID)
	workoutExercise1 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise1.ID)
	workoutExercise2 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise2.ID)
	set1 := createTestCardioSet(t, x.cfg, workoutExercise1)
	set2 := createTestCardioSet(t, x.cfg, workoutExercise2)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/"+x.workoutSession.ID.String()+"/distance", nil, x.token)
	req.SetPathValue("id", x.workoutSession.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetAllSessionDistance), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[distanceResponse](t, rr)
	expected := set1.Distance.Float64 + set2.Distance.Float64
	if response.Distance != expected {
		t.Errorf("expected distance %.2f, got %.2f", expected, response.Distance)
	}
}
func TestHandlerGetAllCardioSessionDuration_Success(t *testing.T) {
	//t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/{id}/duration", apiCfg.handlerGetAllCardioSessionDuration)
	x := testingCardioSetHandlerSetup(t)
	exercise1 := createTestExercise(t, x.cfg, x.user.ID)
	exercise2 := createTestExercise(t, x.cfg, x.user.ID)
	workoutExercise1 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise1.ID)
	workoutExercise2 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise2.ID)
	set1 := createTestCardioSet(t, x.cfg, workoutExercise1)
	set2 := createTestCardioSet(t, x.cfg, workoutExercise2)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/"+x.workoutSession.ID.String()+"/duration", nil, x.token)
	req.SetPathValue("id", x.workoutSession.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetAllCardioSessionDuration), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[cardioDurationResponse](t, rr)
	expected := set1.DurationSeconds.Int32 + set2.DurationSeconds.Int32
	if response.TotalSeconds != int64(expected) {
		t.Errorf("expected %d, got %d", expected, response.TotalSeconds)
	}
}

// helpers
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
