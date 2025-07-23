package handlers

import (
	"fmt"
	"net/http"
	"pull-up-mania/internal/models"
	"time"
)

// HomeHandler displays today's workout or a message if no workout is scheduled
func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	today := getDayOfWeek()

	// Get today's workout plan
	workoutPlan, err := h.storage.GetWorkoutPlanForDay(today)
	if err != nil {
		// No workout scheduled for today
		data := struct {
			NoWorkout bool
			DayName   string
		}{
			NoWorkout: true,
			DayName:   time.Now().Format("Monday"),
		}
		h.templates.ExecuteTemplate(w, "home.html", data)
		return
	}

	// Get or create today's workout session
	session, err := h.storage.GetTodaysWorkoutSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session == nil {
		// Create a new workout session for today
		session = &models.WorkoutSession{
			ID:            generateID(),
			WorkoutPlanID: workoutPlan.ID,
			Date:          time.Now(),
			Completed:     false,
			ExerciseLogs:  make([]models.ExerciseLog, 0),
			Notes:         "",
		}

		// Initialize exercise logs based on the workout plan
		for _, exercise := range workoutPlan.Exercises {
			exerciseLog := models.ExerciseLog{
				ExerciseID:   exercise.ID,
				ExerciseName: exercise.Name,
				Sets:         make([]models.SetLog, exercise.Sets),
				Notes:        "",
			}

			// Initialize sets
			for i := 0; i < exercise.Sets; i++ {
				exerciseLog.Sets[i] = models.SetLog{
					SetNumber: i + 1,
					Weight:    0,
					Reps:      exercise.Reps, // Default to planned reps
					Completed: false,
				}
			}

			session.ExerciseLogs = append(session.ExerciseLogs, exerciseLog)
		}

		if err := h.storage.SaveWorkoutSession(session); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Get settings for weight unit
	settings, err := h.storage.GetSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		WorkoutPlan *models.WorkoutPlan
		Session     *models.WorkoutSession
		WeightUnit  string
		DayName     string
	}{
		WorkoutPlan: workoutPlan,
		Session:     session,
		WeightUnit:  settings.WeightUnit,
		DayName:     time.Now().Format("Monday"),
	}

	h.templates.ExecuteTemplate(w, "home.html", data)
}

// SaveWorkoutHandler handles saving workout progress via HTMX
func (h *Handlers) SaveWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.FormValue("session_id")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	// Get the current session
	session, err := h.storage.GetWorkoutSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse form data and update session
	r.ParseForm()

	// Update exercise logs
	for i := range session.ExerciseLogs {
		exerciseLog := &session.ExerciseLogs[i]

		for j := range exerciseLog.Sets {
			setLog := &exerciseLog.Sets[j]

			weightKey := fmt.Sprintf("exercise_%d_set_%d_weight", i, j)
			repsKey := fmt.Sprintf("exercise_%d_set_%d_reps", i, j)
			completedKey := fmt.Sprintf("exercise_%d_set_%d_completed", i, j)

			if weight := r.FormValue(weightKey); weight != "" {
				setLog.Weight = parseFormFloat(weight)
			}

			if reps := r.FormValue(repsKey); reps != "" {
				setLog.Reps = parseFormInt(reps)
			}

			setLog.Completed = r.FormValue(completedKey) == "on"
		}

		// Update exercise notes
		notesKey := fmt.Sprintf("exercise_%d_notes", i)
		if notes := r.FormValue(notesKey); notes != "" {
			exerciseLog.Notes = notes
		}
	}

	// Update session notes
	if notes := r.FormValue("session_notes"); notes != "" {
		session.Notes = notes
	}

	// Check if workout is completed (all sets marked as completed)
	allCompleted := true
	for _, exerciseLog := range session.ExerciseLogs {
		for _, setLog := range exerciseLog.Sets {
			if !setLog.Completed {
				allCompleted = false
				break
			}
		}
		if !allCompleted {
			break
		}
	}
	session.Completed = allCompleted

	// Save the session
	if err := h.storage.SaveWorkoutSession(session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response for HTMX
	w.Header().Set("Content-Type", "text/html")
	if session.Completed {
		w.Write([]byte(`<div class="alert alert-success">Workout completed! 🎉</div>`))
	} else {
		w.Write([]byte(`<div class="alert alert-info">Progress saved ✓</div>`))
	}
}

// HistoryHandler displays workout history
func (h *Handlers) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get sessions from the last 30 days
	end := time.Now()
	start := end.AddDate(0, 0, -30)

	sessions, err := h.storage.GetWorkoutSessionsByDateRange(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Sessions []models.WorkoutSession
	}{
		Sessions: sessions,
	}

	h.templates.ExecuteTemplate(w, "history.html", data)
}
