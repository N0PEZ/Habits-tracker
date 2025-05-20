package main

import (
	"fmt"

	"huibitica/internal/config"
	"huibitica/internal/logger"
	"huibitica/internal/postgresql"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	log.Logger = logger.SetupLogger(cfg)

	db, err := postgresql.CreateDB()
	fmt.Println(err)
	db, err = postgresql.InitDB("postgresql://postgres:8968@localhost:5432/huibitica")
	fmt.Println(err)
	_ = db
}
