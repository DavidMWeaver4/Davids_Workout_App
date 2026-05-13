package main

import (
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
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	WorkoutDate time.Time `json:"workout_date"`
	Description *string   `json:"description,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type workoutExercise struct {
	ID               uuid.UUID `json:"id"`
	WorkoutSessionID uuid.UUID `json:"workout_session_id"`
	ExerciseID       uuid.UUID `json:"exercise_id"`
	OrderIndex       int32     `json:"order_index"`
	Notes            *string   `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
