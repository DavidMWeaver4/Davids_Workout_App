package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type user struct {
	UserID       uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

type userResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type workoutSessions struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	WorkoutDate time.Time
	Description sql.NullString
	Notes       sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WorkoutExercise struct {
	ID               uuid.UUID
	WorkoutSessionID uuid.UUID
	ExerciseID       uuid.UUID
	OrderIndex       int32
	Notes            sql.NullString
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
