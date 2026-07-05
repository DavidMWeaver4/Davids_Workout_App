package main

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
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
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session1 := createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-24*time.Hour))
	session2 := createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-12*time.Hour))
	session3 := createTestWorkoutSession(t, cfg, user.ID)
	sessionToNotGet := createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-128*time.Hour))
	numberToGet := 3

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/last/count/3?lastX="+strconv.Itoa(numberToGet), nil, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerGetMyXNumberLastSessions), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	response := testingDecodeJSONResponse[[]workoutSessions](t, rr)

	expectedIDs := map[uuid.UUID]bool{
		session1.ID: false,
		session2.ID: false,
		session3.ID: false,
	}

	for _, rep := range response {
		if _, ok := expectedIDs[rep.ID]; ok {
			expectedIDs[rep.ID] = true
		}
		if rep.ID == sessionToNotGet.ID {
			t.Fatalf("got wrong sessions")
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Fatalf("expected session %v in response, but it was not found", id)
		}
	}
}

func TestHandlerCountSessions_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	_ = createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-24*time.Hour))
	_ = createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-12*time.Hour))
	_ = createTestWorkoutSession(t, cfg, user.ID)
	_ = createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-128*time.Hour))

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/count/me", nil, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerCountSessions), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	type countResponse struct {
		Count int64 `json:"count"`
	}
	response := testingDecodeJSONResponse[countResponse](t, rr)

	if response.Count != 4 {
		t.Fatalf("expected 4, got %d", response.Count)
	}

}
func TestHandlerSearchWOSByDescription_Success(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	session1 := createTestWorkoutSessionWithDescription(t, cfg, user.ID, "first sessions descrip")
	session2 := createTestWorkoutSessionWithDescription(t, cfg, user.ID, "second sessions descrip")
	session3 := createTestWorkoutSession(t, cfg, user.ID)
	sessionToGet := createTestWorkoutSessionWithDescription(t, cfg, user.ID, "Let's match here")

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, "/api/v1/workout_sessions/search/desc?query=Let%27s+match", nil, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerSearchWOSByDescription), req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[[]workoutSessions](t, rr)

	for _, rep := range response {
		if rep.ID == session1.ID || rep.ID == session2.ID || rep.ID == session3.ID {
			t.Fatalf("got unexpected session in results: %v", rep.ID)
		}
	}
	found := false
	for _, rep := range response {
		if rep.ID == sessionToGet.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected session was not returned")
	}
}
func TestHandlerSearchWOSByDateRange_Sucess(t *testing.T) {
	cfg := newTestAPIConfig(t)
	user := createTestUser(t, cfg)
	token := testingCreateJWT(t, cfg, user.ID)

	_ = createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-24*time.Hour))
	_ = createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-12*time.Hour))
	_ = createTestWorkoutSession(t, cfg, user.ID)
	sessionOutOfRange := createTestWorkoutSessionWithDate(t, cfg, user.ID, time.Now().UTC().Add(-128*time.Hour))

	endTime := time.Now().UTC()
	startTime := time.Now().UTC().Add(-72 * time.Hour)
	baseURL := "/api/v1/workout_sessions/search/date"
	params := url.Values{}
	params.Add("start", startTime.Format(time.RFC3339))
	params.Add("end", endTime.Format(time.RFC3339))
	fullURL := baseURL + "?" + params.Encode()

	req := testingCreateAuthenticatedJSONRequest(http.MethodGet, fullURL, nil, token)
	rr := testingExecuteRequest(http.HandlerFunc(cfg.handlerSearchWOSByDateRange), req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
	response := testingDecodeJSONResponse[[]workoutSessions](t, rr)

	for _, rep := range response {
		if rep.WorkoutDate.Before(startTime) || rep.WorkoutDate.After(endTime) {
			t.Fatalf("session %v has date %v outside expected range", rep.ID, rep.WorkoutDate)
		}
		if rep.ID == sessionOutOfRange.ID {
			t.Fatal("session outside date range was returned")
		}
	}

}
