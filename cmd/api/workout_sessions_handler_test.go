package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func TestHandlerCreateWorkoutSession_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	type requestBody struct {
		WorkoutDate time.Time `json:"workout_date"`
		Description string    `json:"description"`
		Notes       string    `json:"notes"`
	}
	body := requestBody{
		WorkoutDate: time.Now().UTC(),
		Description: "this is a test workout session",
		Notes:       "test note",
	}
	jsonBody := testingMarshalJSON(t, body)

	req := testingCreateAuthenticatedJSONRequest(http.MethodPost, "/api/v1/workout_sessions", jsonBody, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerCreateWorkoutSession), req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[workoutSessions](t, rr)

	if response.Description == nil || *response.Description != body.Description {
		t.Fatalf("expected %s, got %v", body.Description, response.Description)
	}
	if response.Notes == nil || *response.Notes != body.Notes {
		t.Fatalf("expected %s, got %v", body.Notes, response.Notes)
	}
	if response.UserID != user.ID {
		t.Fatalf("Got wrong user ID, expected %v, got %v", user.ID, response.UserID)
	}
	dbWorkoutSession, err := cfg.db.GetWorkoutSessionByID(context.Background(), response.ID)
	t.Cleanup(func() {
		cfg.db.DeleteWorkoutSession(context.Background(), database.DeleteWorkoutSessionParams{
			ID:     response.ID,
			UserID: dbWorkoutSession.UserID,
		})
	})
	if err != nil {
		t.Fatal("workout session was not inserted into database")
	}
	if dbWorkoutSession.Description.String != body.Description {
		t.Fatal("database value does not match entry data")
	}
}

func TestHandlerDeleteWorkoutSession_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	workoutSessionToDelete := createTestWorkoutSession(t, cfg, user.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodDelete, "/api/v1/workout_sessions/"+workoutSessionToDelete.ID.String(), nil, token)
	req.SetPathValue("id", workoutSessionToDelete.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerDeleteWorkoutSession), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	_, err := cfg.db.GetWorkoutSessionByID(context.Background(), workoutSessionToDelete.ID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHandlerGetAllMyWorkoutSessions_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	workoutSessions1 := createTestWorkoutSession(t, cfg, user.ID)
	workoutSessions2 := createTestWorkoutSession(t, cfg, user.ID)
	workoutSessions3 := createTestWorkoutSession(t, cfg, user.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/me", nil, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetAllMyWorkoutSessions), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[[]workoutSessions](t, rr)
	expectedIDs := map[uuid.UUID]bool{
		workoutSessions1.ID: false,
		workoutSessions2.ID: false,
		workoutSessions3.ID: false,
	}

	for _, rep := range response {
		if _, ok := expectedIDs[rep.ID]; ok {
			expectedIDs[rep.ID] = true
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Fatalf("expected session %v in response, but it was not found", id)
		}
	}

}

func TestHandlerGetWorkoutSessionsByID_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	sessionToFind := createTestWorkoutSession(t, cfg, user.ID)

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/"+sessionToFind.ID.String(), nil, token)
	req.SetPathValue("id", sessionToFind.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetWorkoutSessionById), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[workoutSessions](t, rr)

	if response.ID != sessionToFind.ID {
		t.Fatalf("expected session ID %v, to match %v", response.ID, sessionToFind.ID)
	}
}

func TestHandlerUpdateWorkoutSession_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	sessionToUpdate := createTestWorkoutSession(t, cfg, user.ID)
	type requestBody struct {
		WorkoutDate time.Time `json:"workout_date"`
		Description string    `json:"description"`
		Notes       string    `json:"notes"`
	}
	body := requestBody{
		WorkoutDate: time.Now().UTC(),
		Description: "this is a new test workout session",
		Notes:       "new test note",
	}
	jsonBody := testingMarshalJSON(t, body)

	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/workout_sessions/"+sessionToUpdate.ID.String(), jsonBody, token)
	req.SetPathValue("id", sessionToUpdate.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerUpdateWorkoutSession), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	dbSession, err := cfg.db.GetWorkoutSessionByID(context.Background(), sessionToUpdate.ID)
	if err != nil {
		t.Fatal("session was not updated in database")
	}
	if dbSession.Description.String != body.Description ||
		dbSession.Notes.String != body.Notes {
		t.Fatal("database value does not match entry data")
	}

}

func TestHandlerGetMyLastSession_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	sessionToNotGet := createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-24*time.Hour))
	lastSession := createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC())

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/last/me", nil, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetMyLastSession), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[workoutSessions](t, rr)
	if response.ID == sessionToNotGet.ID {
		t.Fatal("expected last session, got first session")
	}
	if response.ID != lastSession.ID {
		t.Fatalf("expected session id %v, got %v", lastSession.ID, response.ID)
	}
}
func TestHandlerGetMyXNumberLastSessions_Success(t *testing.T) {
	t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/last/count/{lastX}", apiCfg.handlerGetMyXNumberLastSessions)
}
func TestHandlerCountSessions_Success(t *testing.T) {
	t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/count/me", apiCfg.handlerCountSessions)
}
func TestHandlerSearchWOSByDescription_Success(t *testing.T) {
	t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/search/desc", apiCfg.handlerSearchWOSByDescription)
}
func TestHandlerSearchWOSByDateRange_Sucess(t *testing.T) {
	t.Skip()
	//mux.HandleFunc("GET /api/v1/workout_sessions/search/date", apiCfg.handlerSearchWOSByDateRange)
}
