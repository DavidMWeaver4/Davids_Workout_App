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
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerCreateWeightsAndSets), req)
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
		cfg.db.DeleteWeightAndSets(context.Background(), database.DeleteWeightAndSetsParams{
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
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetAllSetsFromSession), req)
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
func TestHandlerDeleteWeightAndSet_Success(t *testing.T) {
	t.Skip()
	//mux.HandleFunc("DELETE /api/v1/weights_and_sets/{id}", apiCfg.handlerDeleteWeightandSet)
}
func TestHandlerUpdateWeightAndSet_Success(t *testing.T) {
	t.Skip()
	//ux.HandleFunc("PUT /api/v1/weights_and_sets/{id}", apiCfg.handlerUpdateWeightAndSets)
}
func TestHandlerGetVolumeSet_Success(t *testing.T) {
	t.Skip()
	//x.HandleFunc("GET /api/v1/weights_and_sets/{id}/volume", apiCfg.handlerGetVolumeSet)
}
func TestHandlerGetTotalVolumeFromAllSets_Success(t *testing.T) {
	t.Skip()
	//x.HandleFunc("GET /api/v1/weights_and_sets/{id}/volume/all", apiCfg.handlerGetTotalVolumeFromAllSet)
}
func TestHandlerGetTotalDuration_Success(t *testing.T) {
	t.Skip()
	//x.HandleFunc("GET /api/v1/weights_and_sets/{id}/duration", apiCfg.handlerGetTotalDuration)
}
func TestHandlerGetTotalDurationFromAllSets_Success(t *testing.T) {
	t.Skip()
	//x.HandleFunc("GET /api/v1/weights_and_sets/{id}/duration/all", apiCfg.handlerGetTotalDurationFromAllSets)
}
