package storage

import (
	"pull-up-mania/internal/models"
	"time"
)

// Storage defines the interface for all storage operations
type Storage interface {
	// Settings
	GetSettings() (*models.Settings, error)
	SaveSettings(settings *models.Settings) error

	// Workout Plans
	GetWorkoutPlans() ([]models.WorkoutPlan, error)
	GetWorkoutPlan(id string) (*models.WorkoutPlan, error)
	SaveWorkoutPlan(plan *models.WorkoutPlan) error
	DeleteWorkoutPlan(id string) error
	GetWorkoutPlanForDay(dayOfWeek int) (*models.WorkoutPlan, error)

	// Workout Sessions
	GetWorkoutSessions() ([]models.WorkoutSession, error)
	GetWorkoutSession(id string) (*models.WorkoutSession, error)
	GetWorkoutSessionsByDateRange(start, end time.Time) ([]models.WorkoutSession, error)
	SaveWorkoutSession(session *models.WorkoutSession) error
	DeleteWorkoutSession(id string) error
	GetTodaysWorkoutSession() (*models.WorkoutSession, error)

	// Initialization
	Init() error
	Close() error
}
