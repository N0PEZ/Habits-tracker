package main

import (
	"huibitica/internal/config"
	"huibitica/internal/handlers"
	"huibitica/internal/logger"
	"huibitica/internal/postgresql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg)

	log.Info().Msg("Logger initialized")

	db, err := postgresql.InitDB(cfg.PostgreAddress, cfg.DBName)
	if err != nil {
		log.Fatal().AnErr("DB INIT ERROR", err)
	}
	log.Info().Msg("Database initialized")
	defer db.Close()

	handler := handlers.NewHandler(db, log)

	_ = handler

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Post("/api/users", handler.NewUser)
	r.Post("/api/habits", handler.NewHabit)
	r.Post("/api/dailies", handler.NewDaily)
	r.Post("/api/tasks", handler.NewTask)

	r.Get("/api/users", handler.GetUser)
	r.Get("/api/habits", handler.GetHabits)
	r.Get("/api/dailies", handler.GetDailies)
	r.Get("/api/tasks", handler.GetTasks)

	r.Put("/api/habits", handler.EditHabit)
	r.Put("/api/dailies", handler.EditDaily)
	r.Put("/api/tasks", handler.EditTask)
	r.Put("/api/users/username", handler.EditUserUsername)
	r.Put("/api/users/email", handler.EditUserEmail)
	r.Put("/api/users/phone", handler.EditUserPhone)
	r.Put("/api/users/password", handler.EditPassword)

	r.Delete("/api/habits", handler.DeleteHabit)
	r.Delete("/api/dailies", handler.DeleteDaily)
	r.Delete("/api/tasks", handler.DeleteTask)
	r.Delete("/api/users", handler.DeleteUser)

	http.ListenAndServe(":8080", r)
}
