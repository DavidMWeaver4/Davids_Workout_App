package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestHandlerCreateWeightAndSets_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	WOSession := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	WOExecercise := createTestWorkoutExercise(t, cfg, WOSession, exercise.ID)

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
		WorkoutExercisesID: WOExecercise.ID,
		Weight:             7.5,
		IsKilogram:         true,
		SetNumber:          3,
		RepsTarget:         12,
		RepsActual:         11,
		DurationSeconds:    45,
		RestTimeSeconds:    90,
	}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPost, "/api/v1/weights_and_sets", JSONBody, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerCreateWeightAndSet), req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[weightAndSets](t, rr)
	if response.Weight != body.Weight {
		t.Fatalf("expected %.2f, got %.2f", body.Weight, response.Weight)
	}
	if response.DurationSeconds == nil {
		t.Fatalf("expected DurationSeconds %d, got nil", body.DurationSeconds)
	}
	if *response.DurationSeconds != body.DurationSeconds {
		t.Fatalf("expected %d, got %d", body.DurationSeconds, *response.DurationSeconds)
	}
	if response.WorkoutExercisesID != WOExecercise.ID {
		t.Fatalf("got wrong session id, expected %v, got %v", WOExecercise.ID, response.WorkoutExercisesID)
	}
	dbWeightSet, err := cfg.db.GetWeightAndSetFromID(context.Background(), response.ID)
	t.Cleanup(func() {
		cfg.db.DeleteWeightAndSet(context.Background(), database.DeleteWeightAndSetParams{
			ID:                 response.ID,
			WorkoutExercisesID: WOExecercise.ID,
		})
	})
	if err != nil {
		t.Fatal("workout exercise was not inserted into database")
	}
	if dbWeightSet.RepsActual != body.RepsActual {
		t.Fatal("database value does not match entry data")
	}
}

func TestHandlerGetAllSetsFromSession_Success(t *testing.T) {
	t.Skip()
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)
	WOSession := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	exercise2 := createTestExercise(t, cfg, user.ID)
	WOExecercise := createTestWorkoutExercise(t, cfg, WOSession, exercise.ID)
	WOExecercise2 := createTestWorkoutExercise(t, cfg, WOSession, exercise2.ID)

	set1 := createTestWeightSet(t, cfg, WOExecercise)
	set2 := createTestWeightSet(t, cfg, WOExecercise2)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_exercises/"+WOSession.ID.String()+"/sets", nil, token)
	req.SetPathValue("id", WOSession.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetAllSetsFromExercise), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[[]weightAndSets](t, rr)
	if len(response) != 2 {
		t.Fatalf("response was not correct size")
	}
	expectedIDs := map[uuid.UUID]bool{set1.ID: false, set2.ID: false}
	for _, item := range response {
		if _, ok := expectedIDs[item.ID]; ok {
			expectedIDs[item.ID] = true
		}
	}
	if expectedIDs[set1.ID] != true || expectedIDs[set2.ID] != true {
		t.Fatal("did not find correct sets")
	}
}
func TestHandlerGetTotalDuration_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)

	set := createTestWeightSet(t, cfg, woExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/weights_and_sets/"+set.ID.String()+"/duration", nil, token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetTotalDuration), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[durationResponse](t, rr)
	expected := set.DurationSeconds.Int32 + set.RestTimeSeconds.Int32

	if response.TotalSeconds != expected {
		t.Fatalf("expected duration %d got %d", expected, response.TotalSeconds)
	}
}
func TestHandlerGetVolumeSet_Success(t *testing.T) {

	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)
	weightSet := createTestWeightSet(t, cfg, woExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/weights_and_sets/"+weightSet.ID.String()+"/volume", nil, token)
	req.SetPathValue("id", weightSet.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetVolumeSet), req)
	t.Log(rr.Body.String())
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[volumeResponse](t, rr)

	expectedVolume := weightSet.Weight * float64(weightSet.RepsActual)

	if response.Volume != expectedVolume {
		t.Fatalf("expected %.2f got %.2f", expectedVolume, response.Volume)
	}
}
func TestHandlerGetTotalVolumeFromExerciseSets_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)

	set1 := createTestWeightSet(t, cfg, woExercise)
	set2 := createTestWeightSet(t, cfg, woExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/weights_and_sets/"+set1.ID.String()+"/volume/all", nil, token)
	req.SetPathValue("id", set1.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetTotalVolumeFromExerciseSets), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[volumeResponse](t, rr)

	expected := (set1.Weight * float64(set1.RepsActual)) + (set2.Weight * float64(set2.RepsActual))

	if response.Volume != expected {
		t.Fatalf("expected volume %.2f got %.2f", expected, response.Volume)
	}
}
func TestHandlerGetTotalSessionVolume_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session := createTestWorkoutSession(t, cfg, user.ID)

	exercise1 := createTestExercise(t, cfg, user.ID)
	exercise2 := createTestExercise(t, cfg, user.ID)

	woExercise1 := createTestWorkoutExercise(t, cfg, session, exercise1.ID)
	woExercise2 := createTestWorkoutExercise(t, cfg, session, exercise2.ID)

	set1 := createTestWeightSet(t, cfg, woExercise1)
	set2 := createTestWeightSet(t, cfg, woExercise2)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/"+session.ID.String()+"/volume", nil, token)
	req.SetPathValue("id", session.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetTotalSessionVolume), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[volumeResponse](t, rr)

	expected := (set1.Weight * float64(set1.RepsActual)) + (set2.Weight * float64(set2.RepsActual))

	if response.Volume != expected {
		t.Fatalf("expected volume %.2f got %.2f", expected, response.Volume)
	}
}
func TestHandlerGetTotalDurationForExercise_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)

	set1 := createTestWeightSet(t, cfg, woExercise)
	set2 := createTestWeightSet(t, cfg, woExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/weights_and_sets/"+set1.ID.String()+"/duration/all", nil, token)
	req.SetPathValue("id", set1.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetTotalDurationForExercise), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[durationResponse](t, rr)

	expected := set1.DurationSeconds.Int32 + set1.RestTimeSeconds.Int32 + set2.DurationSeconds.Int32 + set2.RestTimeSeconds.Int32

	if response.TotalSeconds != expected {
		t.Fatalf("expected duration %d got %d", expected, response.TotalSeconds)
	}
}
func TestHandlerUpdateWeightAndSet_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)

	set := createTestWeightSet(t, cfg, woExercise)

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

	jsonBody := testingMarshalJSON(t, body)

	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/weights_and_sets/"+set.ID.String(), jsonBody, token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerUpdateWeightAndSet), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[weightAndSets](t, rr)

	if response.Weight != body.Weight {
		t.Fatalf("expected weight %.2f got %.2f", body.Weight, response.Weight)
	}

	dbSet, err := cfg.db.GetWeightAndSetFromID(context.Background(), set.ID)

	if err != nil {
		t.Fatal("failed retrieving updated set")
	}

	if dbSet.RepsActual != body.RepsActual {
		t.Fatal("database was not updated")
	}
}
func TestHandlerDeleteWeightAndSet_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)

	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session := createTestWorkoutSession(t, cfg, user.ID)
	exercise := createTestExercise(t, cfg, user.ID)
	woExercise := createTestWorkoutExercise(t, cfg, session, exercise.ID)

	set := createTestWeightSet(t, cfg, woExercise)

	req := testingCreateAuthenticatedJSONRequest(http.MethodDelete, "/api/v1/weights_and_sets/"+set.ID.String(), nil, token)
	req.SetPathValue("id", set.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerDeleteWeightAndSet), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	_, err := cfg.db.GetWeightAndSetFromID(context.Background(), set.ID)

	if err == nil {
		t.Fatal("expected weight set to be deleted")
	}
}
