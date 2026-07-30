package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestHandlerCreateWeightAndSets_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)
	type requestBody struct {
		WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
		Weight             float64   `json:"weight"`
		IsKilogram         bool      `json:"is_kilogram"`
		SetNumber          int32     `json:"set_number"`
		RepsTarget         int32     `json:"reps_target"`
		RepsActual         int32     `json:"reps_actual"`
		DurationSeconds    int32     `json:"duration_seconds"`
		RestTimeSeconds    int32     `json:"rest_time_seconds"`
	}
	body := requestBody{
		WorkoutExercisesID: x.workoutExercise.ID,
		Weight:             7.5,
		IsKilogram:         true,
		SetNumber:          3,
		RepsTarget:         12,
		RepsActual:         11,
		DurationSeconds:    45,
		RestTimeSeconds:    90,
	}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPost, "/api/v1/weights_and_sets", JSONBody, x.token)
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerCreateWeightAndSet), req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[weightAndSets](t, rr)
	if response.Weight != body.Weight {
		t.Errorf("expected %.2f, got %.2f", body.Weight, response.Weight)
	}
	if response.DurationSeconds == nil {
		t.Fatalf("expected DurationSeconds %d, got nil", body.DurationSeconds)
	}
	if *response.DurationSeconds != body.DurationSeconds {
		t.Errorf("expected %d, got %d", body.DurationSeconds, *response.DurationSeconds)
	}
	if response.WorkoutExercisesID != x.workoutExercise.ID {
		t.Errorf("got wrong session id, expected %v, got %v", x.workoutExercise.ID, response.WorkoutExercisesID)
	}
	dbWeightSet, err := x.cfg.db.GetWeightAndSetFromID(context.Background(), response.ID)
	t.Cleanup(func() {
		x.cfg.db.DeleteWeightAndSet(context.Background(), database.DeleteWeightAndSetParams{
			ID:                 response.ID,
			WorkoutExercisesID: x.workoutExercise.ID,
		})
	})
	if err != nil {
		t.Fatal("weight set was not inserted into database")
	}
	if dbWeightSet.RepsActual != body.RepsActual {
		t.Errorf("database value does not match entry data")
	}
}

func TestHandlerGetAllSetsFromSession_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)
	exercise2 := createTestExercise(t, x.cfg, x.user.ID)
	workoutExercise2 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise2.ID)

	set1 := createTestWeightSet(t, x.cfg, x.workoutExercise)
	set2 := createTestWeightSet(t, x.cfg, workoutExercise2)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_session/"+x.workoutSession.ID.String()+"/sets", nil, x.token)
	req.SetPathValue("id", x.workoutSession.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetAllSetsFromSession), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[[]weightAndSets](t, rr)
	if len(response) != 2 {
		t.Errorf("expected 2 sets got %d", len(response))
	}
	expectedIDs := map[uuid.UUID]bool{set1.ID: false, set2.ID: false}
	for _, item := range response {
		if _, ok := expectedIDs[item.ID]; ok {
			expectedIDs[item.ID] = true
		}
	}
	if expectedIDs[set1.ID] != true || expectedIDs[set2.ID] != true {
		t.Errorf("did not find correct sets")
	}
}
func TestHandlerGetTotalDuration_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)

	set := createTestWeightSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/weights_and_sets/"+set.ID.String()+"/duration", nil, x.token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetTotalDuration), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[durationResponse](t, rr)
	expected := set.DurationSeconds.Int32 + set.RestTimeSeconds.Int32

	if response.TotalSeconds != expected {
		t.Errorf("expected duration %d got %d", expected, response.TotalSeconds)
	}
}
func TestHandlerGetVolumeSet_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)
	weightSet := createTestWeightSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/weights_and_sets/"+weightSet.ID.String()+"/volume", nil, x.token)
	req.SetPathValue("id", weightSet.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetVolumeSet), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[volumeResponse](t, rr)

	expectedVolume := weightSet.Weight * float64(weightSet.RepsActual)

	if response.Volume != expectedVolume {
		t.Errorf("expected %.2f got %.2f", expectedVolume, response.Volume)
	}
}
func TestHandlerGetTotalVolumeFromExerciseSets_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)

	set1 := createTestWeightSet(t, x.cfg, x.workoutExercise)
	set2 := createTestWeightSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/weights_and_sets/"+set1.ID.String()+"/volume/all", nil, x.token)
	req.SetPathValue("id", set1.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetTotalVolumeFromExerciseSets), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[volumeResponse](t, rr)

	expected := (set1.Weight * float64(set1.RepsActual)) + (set2.Weight * float64(set2.RepsActual))

	if response.Volume != expected {
		t.Errorf("expected volume %.2f got %.2f", expected, response.Volume)
	}
}
func TestHandlerGetTotalSessionVolume_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)

	exercise1 := createTestExercise(t, x.cfg, x.user.ID)
	exercise2 := createTestExercise(t, x.cfg, x.user.ID)

	workoutExercise1 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise1.ID)
	workoutExercise2 := createTestWorkoutExercise(t, x.cfg, x.workoutSession, exercise2.ID)

	set1 := createTestWeightSet(t, x.cfg, workoutExercise1)
	set2 := createTestWeightSet(t, x.cfg, workoutExercise2)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/"+x.workoutSession.ID.String()+"/volume", nil, x.token)
	req.SetPathValue("id", x.workoutSession.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetTotalSessionVolume), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[volumeResponse](t, rr)

	expected := (set1.Weight * float64(set1.RepsActual)) + (set2.Weight * float64(set2.RepsActual))

	if response.Volume != expected {
		t.Errorf("expected volume %.2f got %.2f", expected, response.Volume)
	}
}
func TestHandlerGetTotalDurationForExercise_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)

	set1 := createTestWeightSet(t, x.cfg, x.workoutExercise)
	set2 := createTestWeightSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/weights_and_sets/"+set1.ID.String()+"/duration/all", nil, x.token)
	req.SetPathValue("id", set1.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerGetTotalDurationForExercise), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[durationResponse](t, rr)

	expected := set1.DurationSeconds.Int32 + set1.RestTimeSeconds.Int32 + set2.DurationSeconds.Int32 + set2.RestTimeSeconds.Int32

	if response.TotalSeconds != expected {
		t.Errorf("expected duration %d got %d", expected, response.TotalSeconds)
	}
}
func TestHandlerUpdateWeightAndSet_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)

	set := createTestWeightSet(t, x.cfg, x.workoutExercise)

	type requestBody struct {
		Weight          float64 `json:"weight"`
		IsKilograms     bool    `json:"is_kilogram"`
		SetNumber       int32   `json:"set_number"`
		RepsTarget      int32   `json:"reps_target"`
		RepsActual      int32   `json:"reps_actual"`
		DurationSeconds int32   `json:"duration_seconds"`
		RestTimeSeconds int32   `json:"rest_time_seconds"`
	}

	body := requestBody{
		Weight:          100,
		IsKilograms:     true,
		SetNumber:       5,
		RepsTarget:      8,
		RepsActual:      8,
		DurationSeconds: 60,
		RestTimeSeconds: 120,
	}

	JSONBody := testingMarshalJSON(t, body)

	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/weights_and_sets/"+set.ID.String(), JSONBody, x.token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateWeightAndSet), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[weightAndSets](t, rr)

	if response.Weight != body.Weight {
		t.Errorf("expected weight %.2f got %.2f", body.Weight, response.Weight)
	}

	dbSet, err := x.cfg.db.GetWeightAndSetFromID(context.Background(), set.ID)

	if err != nil {
		t.Fatal("failed retrieving updated set")
	}

	if dbSet.Weight != body.Weight {
		t.Errorf("database Weight does not match update body")
	}
	if dbSet.RepsActual != body.RepsActual {
		t.Errorf("database RepsActual does not match update body")
	}
	if dbSet.RepsTarget != body.RepsTarget {
		t.Errorf("database RepsTarget does not match update body")
	}
	if dbSet.SetNumber != body.SetNumber {
		t.Errorf("database SetNumber does not match update body")
	}
}
func TestHandlerDeleteWeightAndSet_Success(t *testing.T) {
	x := testingWeightHandlerSetup(t)

	set := createTestWeightSet(t, x.cfg, x.workoutExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodDelete, "/api/v1/weights_and_sets/"+set.ID.String(), nil, x.token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerDeleteWeightAndSet), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	_, err := x.cfg.db.GetWeightAndSetFromID(context.Background(), set.ID)

	if err == nil {
		t.Fatal("expected weight set to be deleted")
	}
}

// helpers
type WeightSetHandlerFixture struct {
	ctx             context.Context
	cfg             *apiConfig
	user            database.User
	token           string
	workoutSession  database.WorkoutSession
	exercise        database.Exercise
	workoutExercise database.WorkoutExercise
}

func testingWeightHandlerSetup(t *testing.T) WeightSetHandlerFixture {
	t.Helper()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	workoutSession := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	workoutExercise := createTestWorkoutExercise(t, cfg, workoutSession, exercise.ID)
	return WeightSetHandlerFixture{
		ctx:             context.Background(),
		cfg:             cfg,
		user:            user,
		token:           token,
		workoutSession:  workoutSession,
		exercise:        exercise,
		workoutExercise: workoutExercise,
	}
}
