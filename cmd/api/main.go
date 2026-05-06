package main

import(
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
	"github.com/DavidMWeaver4/Davids_Workout_App/internal/handlers"
	"github.com/DavidMWeaver4/Davids_Workout_App/internal/middleware"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main(){
	if err := godotenv.Load(); err != nil{
		log.Println("No .env file was found!")
	}
	envPlatform := os.Getenv("PLATFORM")
	dbUrl := os.Getenv("DB_URL")
	secretKey := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil{
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
    	log.Fatal("Cannot reach database: ", err)
	}
	data := database.New(db)

	const port = "8080"
	const filepathRoot = "."
	apiCfg := apiConfig{
		db: data,
		platform: envPlatform,
		jwtSecret: secretKey,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middelwareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	//mux.HandleFunc("POST /api/v1/auth/register", apiCFG.handlerRegister)

	myServer := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(myServer.ListenAndServe())
}
