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

type weightAndSets struct {
	ID                 uuid.UUID `json:"id"`
	WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
	Weight             float64   `json:"weight"`
	IsKilograms        bool      `json:"is_kilogram"`
	SetNumber          int32     `json:"set_number"`
	RepsTarget         int32     `json:"reps_target"`
	RepsActual         int32     `json:"reps_actual"`
	DurationSeconds    *int32    `json:"duration_seconds,omitempty"`
	RestTimeSeconds    *int32    `json:"rest_time_seconds,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type cardioAndSets struct {
	ID                 uuid.UUID `json:"id"`
	WorkoutExercisesID uuid.UUID `json:"workout_exercises_id"`
	SetNumber          int32     `json:"set_number"`
	Distance           *float64  `json:"distance,omitempty"`
	IsKilometers       bool      `json:"is_kilometers"`
	DurationSeconds    *int32    `json:"duration_seconds,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type exercise struct {
	ID              uuid.UUID  `json:"id"`
	UserID          *uuid.UUID `json:"user_id,omitempty"`
	ExerciseName    string     `json:"exercise_name"`
	TargetMuscles   []string   `json:"target_miscles"`
	Equipment       *string    `json:"equipment,omitempty"`
	DifficultyLevel *string    `json:"difficulty_level,omitempty"`
	Description     *string    `json:"description,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
