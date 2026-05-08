package main

import (
	"sync/atomic"

	"github.com/DavidMWeaver4/Davids_Workout_App/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}
