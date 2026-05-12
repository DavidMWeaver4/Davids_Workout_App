package main
import(
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/google/uuid"
)

func(cfg *apiConfig) handlerCreateExercise(w http.ResponseWriter, r *http.Request){
	userID, err := cfg.getUserIDFromToken(w, r)
	if err != nil {
		return
	}
	type parameters struct{
		WorkoutSessionID uuid.UUID `json:"workout_session_id"`
		OrderIndex int32 `json:"order_index"`
		Notes 	string `json:"notes"`
	}
	decoder := json.NewDecoder(r.Body)
	var params parameters
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}
	workSess, err := cfg.db.GetWorkoutSessionByID(r.Context(), params.WorkoutSessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve workout session", err)
		return
	}
	if workSess.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not your session", nil)
		return
	}
	workExer, err := cfg.db.CreateWorkoutExercises(r.Context(), database.CreateWorkoutExercisesParams{
		WorkoutSessionID: params.WorkoutSessionID,
		OrderIndex: params.OrderIndex,
		Notes: sql.NullString{
			String: params.Notes,
			Valid:  params.Notes != "",
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Error storing to database", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, WorkoutExercise{
		ID: workExer.ID,
		WorkoutSessionID: workExer.WorkoutSessionID,
		ExerciseID: workExer.ExerciseID,
		OrderIndex: workExer.OrderIndex,
		Notes: workExer.Notes,
		CreatedAt: workExer.CreatedAt,
		UpdatedAt: workExer.UpdatedAt,
	})
}
