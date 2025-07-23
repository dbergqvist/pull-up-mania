package storage

import (
	"database/sql"
	"encoding/json"
	"pull-up-mania/internal/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	storage := &SQLiteStorage{db: db}
	if err := storage.Init(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *SQLiteStorage) Init() error {
	// Create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS settings (
			id TEXT PRIMARY KEY,
			workout_frequency INTEGER,
			admin_password TEXT,
			weight_unit TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS workout_plans (
			id TEXT PRIMARY KEY,
			name TEXT,
			day_of_week INTEGER,
			exercises TEXT,
			active BOOLEAN
		)`,
		`CREATE TABLE IF NOT EXISTS workout_sessions (
			id TEXT PRIMARY KEY,
			workout_plan_id TEXT,
			date TEXT,
			completed BOOLEAN,
			exercise_logs TEXT,
			notes TEXT
		)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return err
		}
	}

	// Insert default settings if they don't exist
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		defaultSettings := &models.Settings{
			ID:               "default",
			WorkoutFrequency: 3,
			AdminPassword:    "admin123",
			WeightUnit:       "kg",
		}
		if err := s.SaveSettings(defaultSettings); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// Settings
func (s *SQLiteStorage) GetSettings() (*models.Settings, error) {
	var settings models.Settings
	row := s.db.QueryRow("SELECT id, workout_frequency, admin_password, weight_unit FROM settings WHERE id = 'default'")
	err := row.Scan(&settings.ID, &settings.WorkoutFrequency, &settings.AdminPassword, &settings.WeightUnit)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *SQLiteStorage) SaveSettings(settings *models.Settings) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO settings (id, workout_frequency, admin_password, weight_unit) 
		VALUES (?, ?, ?, ?)`,
		settings.ID, settings.WorkoutFrequency, settings.AdminPassword, settings.WeightUnit)
	return err
}

// Workout Plans
func (s *SQLiteStorage) GetWorkoutPlans() ([]models.WorkoutPlan, error) {
	rows, err := s.db.Query("SELECT id, name, day_of_week, exercises, active FROM workout_plans")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.WorkoutPlan
	for rows.Next() {
		var plan models.WorkoutPlan
		var exercisesJSON string

		err := rows.Scan(&plan.ID, &plan.Name, &plan.DayOfWeek, &exercisesJSON, &plan.Active)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(exercisesJSON), &plan.Exercises); err != nil {
			return nil, err
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

func (s *SQLiteStorage) GetWorkoutPlan(id string) (*models.WorkoutPlan, error) {
	var plan models.WorkoutPlan
	var exercisesJSON string

	row := s.db.QueryRow("SELECT id, name, day_of_week, exercises, active FROM workout_plans WHERE id = ?", id)
	err := row.Scan(&plan.ID, &plan.Name, &plan.DayOfWeek, &exercisesJSON, &plan.Active)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(exercisesJSON), &plan.Exercises); err != nil {
		return nil, err
	}

	return &plan, nil
}

func (s *SQLiteStorage) SaveWorkoutPlan(plan *models.WorkoutPlan) error {
	exercisesJSON, err := json.Marshal(plan.Exercises)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`INSERT OR REPLACE INTO workout_plans (id, name, day_of_week, exercises, active) 
		VALUES (?, ?, ?, ?, ?)`,
		plan.ID, plan.Name, plan.DayOfWeek, string(exercisesJSON), plan.Active)
	return err
}

func (s *SQLiteStorage) DeleteWorkoutPlan(id string) error {
	_, err := s.db.Exec("DELETE FROM workout_plans WHERE id = ?", id)
	return err
}

func (s *SQLiteStorage) GetWorkoutPlanForDay(dayOfWeek int) (*models.WorkoutPlan, error) {
	var plan models.WorkoutPlan
	var exercisesJSON string

	row := s.db.QueryRow("SELECT id, name, day_of_week, exercises, active FROM workout_plans WHERE day_of_week = ? AND active = 1", dayOfWeek)
	err := row.Scan(&plan.ID, &plan.Name, &plan.DayOfWeek, &exercisesJSON, &plan.Active)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(exercisesJSON), &plan.Exercises); err != nil {
		return nil, err
	}

	return &plan, nil
}

// Workout Sessions
func (s *SQLiteStorage) GetWorkoutSessions() ([]models.WorkoutSession, error) {
	rows, err := s.db.Query("SELECT id, workout_plan_id, date, completed, exercise_logs, notes FROM workout_sessions ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.WorkoutSession
	for rows.Next() {
		var session models.WorkoutSession
		var dateStr, logsJSON string

		err := rows.Scan(&session.ID, &session.WorkoutPlanID, &dateStr, &session.Completed, &logsJSON, &session.Notes)
		if err != nil {
			return nil, err
		}

		if session.Date, err = time.Parse(time.RFC3339, dateStr); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(logsJSON), &session.ExerciseLogs); err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *SQLiteStorage) GetWorkoutSession(id string) (*models.WorkoutSession, error) {
	var session models.WorkoutSession
	var dateStr, logsJSON string

	row := s.db.QueryRow("SELECT id, workout_plan_id, date, completed, exercise_logs, notes FROM workout_sessions WHERE id = ?", id)
	err := row.Scan(&session.ID, &session.WorkoutPlanID, &dateStr, &session.Completed, &logsJSON, &session.Notes)
	if err != nil {
		return nil, err
	}

	if session.Date, err = time.Parse(time.RFC3339, dateStr); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(logsJSON), &session.ExerciseLogs); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SQLiteStorage) GetWorkoutSessionsByDateRange(start, end time.Time) ([]models.WorkoutSession, error) {
	startStr := start.Format(time.RFC3339)
	endStr := end.Format(time.RFC3339)

	rows, err := s.db.Query("SELECT id, workout_plan_id, date, completed, exercise_logs, notes FROM workout_sessions WHERE date BETWEEN ? AND ? ORDER BY date DESC", startStr, endStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.WorkoutSession
	for rows.Next() {
		var session models.WorkoutSession
		var dateStr, logsJSON string

		err := rows.Scan(&session.ID, &session.WorkoutPlanID, &dateStr, &session.Completed, &logsJSON, &session.Notes)
		if err != nil {
			return nil, err
		}

		if session.Date, err = time.Parse(time.RFC3339, dateStr); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(logsJSON), &session.ExerciseLogs); err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *SQLiteStorage) SaveWorkoutSession(session *models.WorkoutSession) error {
	logsJSON, err := json.Marshal(session.ExerciseLogs)
	if err != nil {
		return err
	}

	dateStr := session.Date.Format(time.RFC3339)

	_, err = s.db.Exec(`INSERT OR REPLACE INTO workout_sessions (id, workout_plan_id, date, completed, exercise_logs, notes) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		session.ID, session.WorkoutPlanID, dateStr, session.Completed, string(logsJSON), session.Notes)
	return err
}

func (s *SQLiteStorage) DeleteWorkoutSession(id string) error {
	_, err := s.db.Exec("DELETE FROM workout_sessions WHERE id = ?", id)
	return err
}

func (s *SQLiteStorage) GetTodaysWorkoutSession() (*models.WorkoutSession, error) {
	today := time.Now().Format("2006-01-02")

	var session models.WorkoutSession
	var dateStr, logsJSON string

	row := s.db.QueryRow("SELECT id, workout_plan_id, date, completed, exercise_logs, notes FROM workout_sessions WHERE date(date) = ?", today)
	err := row.Scan(&session.ID, &session.WorkoutPlanID, &dateStr, &session.Completed, &logsJSON, &session.Notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No session found for today
		}
		return nil, err
	}

	if session.Date, err = time.Parse(time.RFC3339, dateStr); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(logsJSON), &session.ExerciseLogs); err != nil {
		return nil, err
	}

	return &session, nil
}
