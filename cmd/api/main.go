package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file was found!")
	}
	envPlatform := os.Getenv("PLATFORM")
	dbUrl := os.Getenv("DB_URL")
	secretKey := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot reach database: ", err)
	}
	data := database.New(db)

	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{
		db:        data,
		platform:  envPlatform,
		jwtSecret: secretKey,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	//auth handlers
	mux.HandleFunc("POST /api/v1/register", apiCfg.handlerRegister)
	mux.HandleFunc("POST /api/v1/login", apiCfg.handlerLogin)
	//user handlers
	mux.HandleFunc("GET /api/v1/users", apiCfg.handlerUsers)
	mux.HandleFunc("GET /api/v1/users/me", apiCfg.handlerGetMe)
	mux.HandleFunc("PUT /api/v1/users/me/email", apiCfg.handlerChangeEmail)
	mux.HandleFunc("PUT /api/v1/users/me/password", apiCfg.handlerChangePassword)
	mux.HandleFunc("DELETE /api/v1/users/me", apiCfg.handlerDeleteMe)
	//workout_session handlers
	mux.HandleFunc("POST /api/v1/workout_sessions", apiCfg.handlerCreateWorkoutSession)
	mux.HandleFunc("GET /api/v1/workout_sessions/me", apiCfg.handlerGetAllMyWorkoutSessions)
	mux.HandleFunc("GET /api/v1/workout_sessions/{id}", apiCfg.handlerGetWorkoutSessionById)
	mux.HandleFunc("PUT /api/v1/workout_sessions/{id}", apiCfg.handlerUpdateWorkoutSession)
	mux.HandleFunc("DELETE /api/v1/workout_sessions/{id}", apiCfg.handlerDeleteWorkoutSession)
	mux.HandleFunc("GET /api/v1/workout_sessions/last/me", apiCfg.handlerGetMyLastSession)
	mux.HandleFunc("GET /api/v1/workout_sessions/last/count/{lastX}", apiCfg.handlerGetMyXNumberLastSessions)
	mux.HandleFunc("GET /api/v1/workout_sessions/count/me", apiCfg.handlerCountSessions)
	mux.HandleFunc("GET /api/v1/workout_sessions/search/desc", apiCfg.handlerSearchWOSByDescription)
	mux.HandleFunc("GET /api/v1/workout_sessions/search/date", apiCfg.handlerSearchWOSByDateRange)
	//workout_exercise handlers
	mux.HandleFunc("POST /api/v1/workout_exercises", apiCfg.handlerCreateWorkoutExercise)
	mux.HandleFunc("GET /api/v1/workout_exercises/{id}", apiCfg.handlerGetWorkoutExercises)
	mux.HandleFunc("GET /api/v1/workout_sessions/{session_id}/exercises", apiCfg.handlerGetWorkoutExercisesInSession)
	mux.HandleFunc("DELETE /api/v1/workout_exercises/{id}", apiCfg.handlerDeleteWorkoutExercises)
	mux.HandleFunc("GET /api/v1/workout_sessions/{session_id}/exercises/count", apiCfg.handlerGetNumOfWorkoutsInSession)
	//weights and sets handlers
	mux.HandleFunc("POST /api/v1/weights_and_sets", apiCfg.handlerCreateWeightsAndSets)
	mux.HandleFunc("GET /api/v1/workout_exercises/{id}/sets", apiCfg.handlerGetAllSetsFromSession)
	mux.HandleFunc("DELETE /api/v1/weights_and_sets/{id}", apiCfg.handlerDeleteWeightandSet)
	mux.HandleFunc("PUT /api/v1/weights_and_sets/{id}", apiCfg.handlerUpdateWeightAndSets)
	mux.HandleFunc("GET /api/v1/weights_and_sets/{id}/volume", apiCfg.handlerGetVolumeSet)
	mux.HandleFunc("GET /api/v1/weights_and_sets/{id}/volume/all", apiCfg.handlerGetTotalVolumeFromAllSet)
	mux.HandleFunc("GET /api/v1/weights_and_sets/{id}/duration", apiCfg.handlerGetTotalDuration)
	mux.HandleFunc("GET /api/v1/weights_and_sets/{id}/duration/all", apiCfg.handlerGetTotalDurationFromAllSets)
	//cardio and sets handlers
	mux.HandleFunc("POST /api/v1/cardio_and_sets", apiCfg.handlerCreateCardioAndSets)
	mux.HandleFunc("GET /api/v1/workout_exercises/{id}/cardio", apiCfg.handlerGetAllCardioFromSession)
	mux.HandleFunc("DELETE /api/v1/cardio_and_sets/{id}", apiCfg.handlerDeleteCardioAndSets)
	mux.HandleFunc("PUT /api/v1/cardio_and_sets/{id}", apiCfg.handlerUpdateCardioAndSets)
	mux.HandleFunc("GET /api/v1/cardio_and_sets/{id}/distance", apiCfg.handlerGetSetDistance)
	mux.HandleFunc("GET /api/v1/cardio_and_sets/{id}/distance/all", apiCfg.handlerGetAllSetsDistance)
	mux.HandleFunc("GET /api/v1/cardio_and_sets/{id}/duration", apiCfg.handlerGetSetDuration)
	mux.HandleFunc("GET /api/v1/cardio_and_sets/{id}/duration/all", apiCfg.handlerGetAllSetsDuration)
	//exercise handlers
	mux.HandleFunc("POST /api/v1/exercises", apiCfg.handlerCreateExercises)
	mux.HandleFunc("DELETE /api/v1/exercises/{id}", apiCfg.handlerDeleteExercisesByID)
	mux.HandleFunc("GET /api/v1/exercises/search", apiCfg.handlerSearchExercises)
	mux.HandleFunc("GET /api/v1/exercises/{id}", apiCfg.handlerGetExerciseFromID)
	mux.HandleFunc("PUT /api/v1/exercises/{id}", apiCfg.handlerUpdateExercise)

	myServer := http.Server{
		Addr:    ":" + port,
		Handler: middlewareLogger(mux),
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(myServer.ListenAndServe())
}
