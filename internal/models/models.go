package models

import (
	"time"
)

// Exercise represents a single exercise with its configuration
type Exercise struct {
	ID          string `json:"id" firestore:"id"`
	Name        string `json:"name" firestore:"name"`
	Sets        int    `json:"sets" firestore:"sets"`
	Reps        int    `json:"reps" firestore:"reps"`
	Description string `json:"description" firestore:"description"`
}

// WorkoutPlan represents a workout plan for a specific day
type WorkoutPlan struct {
	ID        string     `json:"id" firestore:"id"`
	Name      string     `json:"name" firestore:"name"`
	DayOfWeek int        `json:"day_of_week" firestore:"day_of_week"` // 0=Sunday, 1=Monday, etc.
	Exercises []Exercise `json:"exercises" firestore:"exercises"`
	Active    bool       `json:"active" firestore:"active"`
}

// WorkoutSession represents a completed or in-progress workout
type WorkoutSession struct {
	ID            string        `json:"id" firestore:"id"`
	WorkoutPlanID string        `json:"workout_plan_id" firestore:"workout_plan_id"`
	Date          time.Time     `json:"date" firestore:"date"`
	Completed     bool          `json:"completed" firestore:"completed"`
	ExerciseLogs  []ExerciseLog `json:"exercise_logs" firestore:"exercise_logs"`
	Notes         string        `json:"notes" firestore:"notes"`
}

// ExerciseLog represents the logged performance for a specific exercise
type ExerciseLog struct {
	ExerciseID   string   `json:"exercise_id" firestore:"exercise_id"`
	ExerciseName string   `json:"exercise_name" firestore:"exercise_name"`
	Sets         []SetLog `json:"sets" firestore:"sets"`
	Notes        string   `json:"notes" firestore:"notes"`
}

// SetLog represents the performance for a single set
type SetLog struct {
	SetNumber int     `json:"set_number" firestore:"set_number"`
	Weight    float64 `json:"weight" firestore:"weight"`
	Reps      int     `json:"reps" firestore:"reps"`
	Completed bool    `json:"completed" firestore:"completed"`
}

// Settings represents global app settings
type Settings struct {
	ID               string `json:"id" firestore:"id"`
	WorkoutFrequency int    `json:"workout_frequency" firestore:"workout_frequency"` // workouts per week
	AdminPassword    string `json:"admin_password" firestore:"admin_password"`
	WeightUnit       string `json:"weight_unit" firestore:"weight_unit"` // "kg" or "lbs"
}
