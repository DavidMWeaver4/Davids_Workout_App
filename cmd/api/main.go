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
	mux.HandleFunc("POST /api/v1/register", apiCfg.handlerRegister)
	mux.HandleFunc("POST /api/v1/login", apiCfg.handlerLogin)
	mux.HandleFunc("GET /api/v1/users", apiCfg.handlerUsers)
	mux.HandleFunc("GET /api/v1/users/me", apiCfg.handlerGetMe)
	mux.HandleFunc("PUT /api/v1/users/me/email", apiCfg.handlerChangeEmail)
	mux.HandleFunc("PUT /api/v1/users/me/password", apiCfg.handlerChangePassword)
	mux.HandleFunc("DELETE /api/v1/users/me", apiCfg.handlerDeleteMe)
	mux.HandleFunc("POST /api/v1/workout_sessions", apiCfg.handlerCreateWorkoutSession)
	mux.HandleFunc("GET /api/v1/workout_sessions/me", apiCfg.handlerGetAllMyWorkoutSessions)
	mux.HandleFunc("GET /api/v1/workout_sessions/{id}", apiCfg.handlerGetWorkoutSessionById)
	mux.HandleFunc("PUT /api/v1/workout_sessions/{id}", apiCfg.handlerUpdateWorkoutSession)
	mux.HandleFunc("DELETE /api/v1/workout_sessions/{id}", apiCfg.handlerDeleteWorkoutSession)
	mux.HandleFunc("GET /api/v1/workout_sessions/last/me", apiCfg.handlerGetMyLastSession)
	mux.HandleFunc("GET /api/v1/workout_sessions/getX/me", apiCfg.handlerGetMyXNumberLastSessions)
	mux.HandleFunc("GET /api/v1/workout_sessions/count", apiCfg.handlerCountSessions)
	mux.HandleFunc("GET /api/v1/workout_sessions/search/desc", apiCfg.handlerSearchWOSByDescription)
	mux.HandleFunc("GET /api/v1/workout_sessions/search/date", apiCfg.handlerSearchWOSByDateRange)
	mux.HandleFunc("POST /api/v1/workout_exercises", apiCfg.handlerCreateWorkoutExercise)
	mux.HandleFunc("GET /api/v1/workout_exercises/{id}", apiCfg.handlerGetWorkoutExercises)
	mux.HandleFunc("GET /api/v1/workout_sessions/{session_id}/exercises", apiCfg.handlerGetWorkoutExercisesInSession)
	mux.HandleFunc("DELETE /api/v1/workout_exercises/{id}", apiCfg.handlerDeleteWorkoutExercises)
	mux.HandleFunc("GET /api/v1/workout_sessions/{session_id}/exercises/count", apiCfg.handlerGetNumOfWorkoutsInSession)

	myServer := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(myServer.ListenAndServe())
}
