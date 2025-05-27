package handlers

import (
	"encoding/json"
	"huibitica/internal/logger"
	"huibitica/internal/models"
	"huibitica/internal/postgresql"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewHandler(db *pgxpool.Pool, log zerolog.Logger) *Handler {
	return &Handler{db: db, log: log}
}

func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterUserRequest

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		h.log.Debug().Str("request_id", requestID).Bytes("request_body", body).Msg("Request body content")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Валидация
	if err := json.Unmarshal(body, &user); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Запись в БД
	if err := postgresql.RegisterUser(user, h.db); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Registration failed")

		// Обработка специфичных ошибок БД
		switch err.Error() {
		case "username already exists":
			log.Warn().Str("request_id", requestID).Msg("username already exists")
			http.Error(w, err.Error(), http.StatusConflict)
		case "email already exists":
			log.Warn().Str("request_id", requestID).Msg("email already exists")
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	log.Debug().Str("request_id", requestID).Msg("OK")
	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User registered",
	})
}

func (h *Handler) NewHabit(w http.ResponseWriter, r *http.Request) {
	var habit models.Habit

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		h.log.Debug().Str("request_id", requestID).Bytes("request_body", body).Msg("Request body content")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Валидация
	if err := json.Unmarshal(body, &habit); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Логирование попытки создания
	h.log.Info().Msg("Attempting to create new habit")

	// Запись в БД
	if err := postgresql.AddHabit(habit, h.db); err != nil {
		h.log.Error().Err(err).Str("request_id", requestID).Msg("Failed to create habit")

		http.Error(w, "Failed to create habit", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status":  "success",
		"message": "Habit created successfully",
		"habit":   habit,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to encode response")
	}
}

func (h *Handler) NewDaily(w http.ResponseWriter, r *http.Request) {
	var daily models.Daily

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		h.log.Debug().Str("request_id", requestID).Bytes("request_body", body).Msg("Request body content")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &daily); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().Str("request_id", requestID).Msg("Attempting to create new habit")

	if err := postgresql.AddDaily(daily, h.db); err != nil {
		h.log.Error().Err(err).Str("request_id", requestID).Msg("Failed to create habit")
		http.Error(w, "Failed to create habit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status":  "success",
		"message": "Habit created successfully",
		"daily":   daily,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to encode response")
	}
}

func (h *Handler) NewTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		h.log.Debug().Str("request_id", requestID).Bytes("request_body", body).Msg("Request body content")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &task); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().Str("request_id", requestID).Msg("Attempting to create new task")
	if err := postgresql.AddTask(task, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to create task")
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"status":  "success",
		"message": "Task created successfully",
		"task":    task,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := req.UserID

	h.log.Info().Str("request_id", requestID).Int("user_id", userID).Msg("Fetching user data")

	user, err := postgresql.GetUserByID(userID, h.db)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to fetch user data")
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetHabits(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := req.UserID

	h.log.Info().Str("request_id", requestID).Int("user_id", userID).Msg("Fetching habits")

	habits, err := postgresql.GetHabits(userID, h.db)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to fetch habits")
		http.Error(w, "Failed to fetch habits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(habits); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetDailies(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := req.UserID

	h.log.Info().Str("request_id", requestID).Int("user_id", userID).Msg("Fetching dailies")

	dailies, err := postgresql.GetDailies(userID, h.db)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to fetch dailies")
		http.Error(w, "Failed to fetch dailies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dailies); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := req.UserID

	h.log.Info().Str("request_id", requestID).Int("user_id", userID).Msg("Fetching tasks")

	tasks, err := postgresql.GetTasks(userID, h.db)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to fetch tasks")
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) EditHabit(w http.ResponseWriter, r *http.Request) {
	var habit models.Habit

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &habit); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().Str("request_id", requestID).Msg("Attempting to edit habit")

	if err := postgresql.EditHabit(habit, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to edit habit")
		http.Error(w, "Failed to edit habit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "success",
		"message": "Habit edited successfully",
		"habit":   habit,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditDaily(w http.ResponseWriter, r *http.Request) {
	var daily models.Daily

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &daily); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().
		Str("request_id", requestID).Msg("Attempting to edit habit")

	if err := postgresql.EditDaily(daily, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to edit habit")

		http.Error(w, "Failed to edit habit", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &task); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().Str("request_id", requestID).Msg("Attempting to edit task")

	if err := postgresql.EditTask(task, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to edit task")
		http.Error(w, "Failed to edit task", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditUserUsername(w http.ResponseWriter, r *http.Request) {
	var user models.EditUserData

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().Str("request_id", requestID).Msg("Attempting to edit username")

	if err := postgresql.EditUserUsername(user, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to edit username")
		http.Error(w, "Failed to edit username", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditUserEmail(w http.ResponseWriter, r *http.Request) {
	var user models.EditUserData

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().Str("request_id", requestID).Msg("Attempting to edit email")

	if err := postgresql.EditUserEmail(user, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to edit email")
		http.Error(w, "Failed to edit email", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditUserPhone(w http.ResponseWriter, r *http.Request) {
	var user models.EditUserData

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().Str("request_id", requestID).Msg("Attempting to edit phone")

	if err := postgresql.EditUserPhone(user, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to edit phone")
		http.Error(w, "Failed to edit phone", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditPassword(w http.ResponseWriter, r *http.Request) {
	var user models.Password

	requestID := middleware.GetReqID(r.Context())

	body, err := logger.RequestLogger(requestID, r, h.log)
	if err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &user); err != nil {
		h.log.Warn().Str("request_id", requestID).Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Info().Str("request_id", requestID).Msg("Attempting to edit password")

	if err := postgresql.EditPassword(user, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to edit password")
		http.Error(w, "Failed to edit password", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	var req struct {
		HabitID int `json:"habit_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	habitID := req.HabitID

	h.log.Info().Str("request_id", requestID).Int("habit_id", habitID).Msg("Attempting to delete habit")

	if err := postgresql.DeleteHabit(habitID, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to delete habit")
		http.Error(w, "Failed to delete habit", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteDaily(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	var req struct {
		DailyID int `json:"daily_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	dailyID := req.DailyID

	h.log.Info().Str("request_id", requestID).Int("daily_id", dailyID).Msg("Attempting to delete daily")

	if err := postgresql.DeleteDaily(dailyID, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to delete daily")
		http.Error(w, "Failed to delete daily", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	var req struct {
		TaskID int `json:"task_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	taskID := req.TaskID

	h.log.Info().Str("request_id", requestID).Int("task_id", taskID).Msg("Attempting to delete task")

	if err := postgresql.DeleteTask(taskID, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to delete task")
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := req.UserID

	h.log.Info().Str("request_id", requestID).Int("user_id", userID).Msg("Attempting to delete user")

	if err := postgresql.DeleteUser(userID, h.db); err != nil {
		h.log.Error().Str("request_id", requestID).Err(err).Msg("Failed to delete user")
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
