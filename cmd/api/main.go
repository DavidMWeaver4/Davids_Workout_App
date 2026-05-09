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

	myServer := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(myServer.ListenAndServe())
}
