package main

import (
	"huibitica/internal/config"
	"huibitica/internal/logger"
	"huibitica/internal/postgresql"

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

	_ = db

	r := chi.NewRouter()

	r.Use(logger.RequestLogger(log))
	r.Use(middleware.Recoverer)
}
