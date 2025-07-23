package handlers

import (
	"fmt"
	"log"
	"net/http"
	"pull-up-mania/internal/models"
	"strconv"
)

// AdminLoginHandler displays the admin login page
func (h *Handlers) AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		settings, err := h.storage.GetSettings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if password == settings.AdminPassword {
			// Set a simple session cookie (in production, use proper session management)
			http.SetCookie(w, &http.Cookie{
				Name:  "admin_session",
				Value: "authenticated",
				Path:  "/",
			})
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		} else {
			data := struct {
				Error string
			}{
				Error: "Invalid password",
			}
			h.templates.ExecuteTemplate(w, "admin_login.html", data)
			return
		}
	}

	h.templates.ExecuteTemplate(w, "admin_login.html", nil)
}

// AdminDashboardHandler displays the admin dashboard
func (h *Handlers) AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAdminAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	workoutPlans, err := h.storage.GetWorkoutPlans()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	settings, err := h.storage.GetSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		WorkoutPlans []models.WorkoutPlan
		Settings     *models.Settings
		DayNames     []string
	}{
		WorkoutPlans: workoutPlans,
		Settings:     settings,
		DayNames:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	}

	err = h.templates.ExecuteTemplate(w, "admin_dashboard.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AdminSettingsHandler handles updating app settings
func (h *Handlers) AdminSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAdminAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		settings, err := h.storage.GetSettings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update settings from form
		if freq := r.FormValue("workout_frequency"); freq != "" {
			if f, err := strconv.Atoi(freq); err == nil {
				settings.WorkoutFrequency = f
			}
		}

		if password := r.FormValue("admin_password"); password != "" {
			settings.AdminPassword = password
		}

		if unit := r.FormValue("weight_unit"); unit != "" {
			settings.WeightUnit = unit
		}

		if err := h.storage.SaveSettings(settings); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return success response for HTMX
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="alert alert-success">Settings updated successfully!</div>`))
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// AdminWorkoutPlanHandler handles creating/editing workout plans
func (h *Handlers) AdminWorkoutPlanHandler(w http.ResponseWriter, r *http.Request) {
	if !h.isAdminAuthenticated(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		action := r.FormValue("action")

		switch action {
		case "create", "edit":
			h.saveWorkoutPlan(w, r)
		case "delete":
			h.deleteWorkoutPlan(w, r)
		default:
			http.Error(w, "Invalid action", http.StatusBadRequest)
		}
		return
	}

	// GET request - show form
	planID := r.URL.Query().Get("id")
	if planID != "" {
		// Edit existing plan
		plan, err := h.storage.GetWorkoutPlan(planID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Plan     *models.WorkoutPlan
			DayNames []string
			IsEdit   bool
		}{
			Plan:     plan,
			DayNames: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			IsEdit:   true,
		}

		h.templates.ExecuteTemplate(w, "admin_workout_plan.html", data)
	} else {
		// Create new plan
		data := struct {
			Plan     *models.WorkoutPlan
			DayNames []string
			IsEdit   bool
		}{
			Plan: &models.WorkoutPlan{
				ID:        generateID(),
				Exercises: []models.Exercise{},
				Active:    true,
			},
			DayNames: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			IsEdit:   false,
		}

		h.templates.ExecuteTemplate(w, "admin_workout_plan.html", data)
	}
}

func (h *Handlers) saveWorkoutPlan(w http.ResponseWriter, r *http.Request) {
	planID := r.FormValue("plan_id")
	if planID == "" {
		planID = generateID()
	}

	name := r.FormValue("name")
	dayOfWeek, _ := strconv.Atoi(r.FormValue("day_of_week"))
	active := r.FormValue("active") == "on"

	// Parse exercises from form
	var exercises []models.Exercise

	// Get exercise count
	exerciseCount := 0
	for {
		if r.FormValue(fmt.Sprintf("exercise_%d_name", exerciseCount)) == "" {
			break
		}
		exerciseCount++
	}

	for i := 0; i < exerciseCount; i++ {
		exerciseName := r.FormValue(fmt.Sprintf("exercise_%d_name", i))
		if exerciseName == "" {
			continue
		}

		sets, _ := strconv.Atoi(r.FormValue(fmt.Sprintf("exercise_%d_sets", i)))
		reps, _ := strconv.Atoi(r.FormValue(fmt.Sprintf("exercise_%d_reps", i)))
		description := r.FormValue(fmt.Sprintf("exercise_%d_description", i))

		exercise := models.Exercise{
			ID:          generateID(),
			Name:        exerciseName,
			Sets:        sets,
			Reps:        reps,
			Description: description,
		}

		exercises = append(exercises, exercise)
	}

	plan := &models.WorkoutPlan{
		ID:        planID,
		Name:      name,
		DayOfWeek: dayOfWeek,
		Exercises: exercises,
		Active:    active,
	}

	if err := h.storage.SaveWorkoutPlan(plan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to admin dashboard
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handlers) deleteWorkoutPlan(w http.ResponseWriter, r *http.Request) {
	planID := r.FormValue("plan_id")
	if planID == "" {
		http.Error(w, "Plan ID required", http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteWorkoutPlan(planID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response for HTMX
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div class="alert alert-success">Workout plan deleted successfully!</div>`))
}

// AdminLogoutHandler handles admin logout
func (h *Handlers) AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "admin_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// isAdminAuthenticated checks if the user is authenticated as admin
func (h *Handlers) isAdminAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return false
	}
	return cookie.Value == "authenticated"
}
