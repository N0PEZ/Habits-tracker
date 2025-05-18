package logger

import (
	"huibitica/internal/config"
	"os"

	"github.com/rs/zerolog"
)

func SetupLogger(cfg *config.Config) zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Dev: только красивый вывод в консоль
	if cfg.Env == "dev" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else if cfg.Env == "prod" { // Prod: логи в файл + консоль
		logFile, err := os.OpenFile(
			"app.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to open log file")
		}

		// MultiWriter: пишем и в консоль, и в файл
		multiWriter := zerolog.MultiLevelWriter(
			zerolog.ConsoleWriter{Out: os.Stderr}, // Консоль
			logFile,                               // Файл
		)
		logger = logger.Output(multiWriter)
	}

	return logger
}
