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

func TestHandlerGetWorkoutExercisesInSession_Success(t *testing.T) {
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
	if dbWOExercise.Notes.String != *body.Notes {
		t.Errorf("expected notes %q, got %q", *body.Notes, dbWOExercise.Notes.String)
	}
	if dbWOExercise.OrderIndex != body.OrderIndex {
		t.Errorf("expected order index %d, got %d", body.OrderIndex, dbWOExercise.OrderIndex)
	}
	if dbWOExercise.ExerciseID != newExercise.ID {
		t.Errorf("expected exercise id %v, got %v", newExercise.ID, dbWOExercise.ExerciseID)
	}
}
func TestHandlerUpdateWorkoutExercise_Unauthorized(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_Forbidden(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_NotFound(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_InvalidWorkoutExerciseID(t *testing.T) {
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
	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/workout_exercises/"+"notarealuuid", JSONBody, x.token)
	req.SetPathValue("id", "notarealuuid")
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateWorkoutExercise), req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	dbWOExercise, err := x.cfg.db.GetWorkoutExerciseFromID(context.Background(), workoutExercise.ID)
	if err != nil {
		t.Fatal("workout exercise should still exist in database")
	}
	if dbWOExercise.ExerciseID != x.exercise.ID {
		t.Fatal("exercise ID should not have changed")
	}
}
func TestHandlerUpdateWorkoutExercise_InvalidJSON(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	JSONBody := []byte(`{"exercise_name":`)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/workout_exercises/"+workoutExercise.ID.String(), JSONBody, x.token)
	req.SetPathValue("id", workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateWorkoutExercise), req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	dbWOExercise, err := x.cfg.db.GetWorkoutExerciseFromID(context.Background(), workoutExercise.ID)
	if err != nil {
		t.Fatal("workout exercise should still exist")
	}

	if dbWOExercise.ExerciseID != workoutExercise.ExerciseID {
		t.Fatal("exercise should not have been modified")
	}

	if dbWOExercise.OrderIndex != workoutExercise.OrderIndex {
		t.Fatal("order index should not have been modified")
	}

	if dbWOExercise.Notes.String != workoutExercise.Notes.String {
		t.Fatal("notes should not have been modified")
	}
}
func TestHandlerUpdateWorkoutExercise_NoFieldsProvided(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	type requestBody struct {
		ExerciseName *string `json:"exercise_name,omitempty"`
		OrderIndex   *int32  `json:"order_index,omitempty"`
		Notes        *string `json:"notes,omitempty"`
	}
	body := requestBody{}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/workout_exercises/"+workoutExercise.ID.String(), JSONBody, x.token)
	req.SetPathValue("id", workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateWorkoutExercise), req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	dbWOExercise, err := x.cfg.db.GetWorkoutExerciseFromID(context.Background(), workoutExercise.ID)
	if err != nil {
		t.Fatal("workout exercise should still exist")
	}

	if dbWOExercise.ExerciseID != workoutExercise.ExerciseID {
		t.Fatal("exercise should not have been modified")
	}

	if dbWOExercise.OrderIndex != workoutExercise.OrderIndex {
		t.Fatal("order index should not have been modified")
	}

	if dbWOExercise.Notes.String != workoutExercise.Notes.String {
		t.Fatal("notes should not have been modified")
	}
}
func TestHandlerUpdateWorkoutExercise_ExerciseNameEmpty(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	type requestBody struct {
		ExerciseName string  `json:"exercise_name"`
		OrderIndex   int32   `json:"order_index"`
		Notes        *string `json:"notes,omitempty"`
	}
	newString := "updated test note"
	body := requestBody{
		ExerciseName: "",
		OrderIndex:   2,
		Notes:        &newString,
	}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/workout_exercises/"+workoutExercise.ID.String(), JSONBody, x.token)
	req.SetPathValue("id", workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateWorkoutExercise), req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
	dbWOExercise, err := x.cfg.db.GetWorkoutExerciseFromID(context.Background(), workoutExercise.ID)
	if err != nil {
		t.Fatal("workout exercise should still exist in database")
	}
	if dbWOExercise.ExerciseID != x.exercise.ID {
		t.Fatal("exercise ID should not have changed")
	}
}
func TestHandlerUpdateWorkoutExercise_ExerciseNameNotFound(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	type requestBody struct {
		ExerciseName string  `json:"exercise_name"`
		OrderIndex   int32   `json:"order_index"`
		Notes        *string `json:"notes,omitempty"`
	}
	newString := "updated test note"
	body := requestBody{
		ExerciseName: "invalid exercise name",
		OrderIndex:   2,
		Notes:        &newString,
	}
	JSONBody := testingMarshalJSON(t, body)
	req := testingCreateAuthenticatedJSONRequest(http.MethodPut, "/api/v1/workout_exercises/"+workoutExercise.ID.String(), JSONBody, x.token)
	req.SetPathValue("id", workoutExercise.ID.String())
	rr := testingExecuteRequest(http.HandlerFunc(x.cfg.handlerUpdateWorkoutExercise), req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d got %d: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
	dbWOExercise, err := x.cfg.db.GetWorkoutExerciseFromID(context.Background(), workoutExercise.ID)
	if err != nil {
		t.Fatal("workout exercise should still exist in database")
	}
	if dbWOExercise.ExerciseID != x.exercise.ID {
		t.Fatal("exercise ID should not have changed")
	}
}
func TestHandlerUpdateWorkoutExercise_ExerciseNameOnly_PreservesOtherFields(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	newExercise := createTestExercise(t, x.cfg, x.user.ID)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	type requestBody struct {
		ExerciseName string `json:"exercise_name"`
	}
	body := requestBody{
		ExerciseName: newExercise.ExerciseName,
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
	if dbWOExercise.Notes.String != workoutExercise.Notes.String {
		t.Fatalf("expected notes %q, got %q", workoutExercise.Notes.String, dbWOExercise.Notes.String)
	}

	if dbWOExercise.OrderIndex != workoutExercise.OrderIndex {
		t.Fatalf("expected order index %d, got %d", workoutExercise.OrderIndex, dbWOExercise.OrderIndex)
	}

	if dbWOExercise.ExerciseID != newExercise.ID {
		t.Fatalf("expected exercise ID %v, got %v", newExercise.ID, dbWOExercise.ExerciseID)
	}
}
func TestHandlerUpdateWorkoutExercise_OrderIndexNegative(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_OrderIndexOnly_PreservesOtherFields(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	type requestBody struct {
		OrderIndex int32 `json:"order_index"`
	}
	body := requestBody{
		OrderIndex: 2,
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
		t.Fatal("workout exercise was not in database")
	}
	if dbWOExercise.Notes.String != workoutExercise.Notes.String {
		t.Errorf("expected notes %q, got %q", workoutExercise.Notes.String, dbWOExercise.Notes.String)
	}
	if dbWOExercise.OrderIndex != body.OrderIndex {
		t.Errorf("expected order index %d, got %d", body.OrderIndex, dbWOExercise.OrderIndex)
	}
	if dbWOExercise.ExerciseID != x.exercise.ID {
		t.Errorf("expected exercise id %v, got %v", x.exercise.ID, dbWOExercise.ExerciseID)
	}
	if dbWOExercise.OrderIndex == workoutExercise.OrderIndex {
		t.Fatal("database value did not update")
	}
}
func TestHandlerUpdateWorkoutExercise_ClearNotes(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_UpdateNotesFromEmpty(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_NotesOnly_SetNonEmpty_PreservesOtherFields(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_NotesOmitted_PreservesExistingNotes(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_UpdateMultipleFields(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_ResponseBodyMatchesUpdateRecord(t *testing.T) {
	t.Skip()
}
func TestHandlerUpdateWorkoutExercise_ExerciseNameTrimmed(t *testing.T) {
	x := testingWorkoutExerciseHandlerSetup(t)
	newExercise := createTestExercise(t, x.cfg, x.user.ID)
	workoutExercise := createTestWorkoutExercise(t, x.cfg, x.workoutSession, x.exercise.ID)
	originalExerciseID := workoutExercise.ExerciseID
	type requestBody struct {
		ExerciseName string `json:"exercise_name"`
	}
	body := requestBody{
		ExerciseName: "      " + newExercise.ExerciseName + "        ",
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
		t.Fatal("workout exercise was not in database")
	}
	if dbWOExercise.ExerciseID == originalExerciseID {
		t.Fatal("expected exercise to change after trimming the exercise name")
	}
	if dbWOExercise.ExerciseID != newExercise.ID {
		t.Fatalf("expected exercise ID %v, got %v", newExercise.ID, dbWOExercise.ExerciseID)
	}
}
func TestHandlerUpdateWorkoutExercise_ExerciseNameInternalWhitespace(t *testing.T) {
	t.Skip()
}

//TODO Fully test update
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
