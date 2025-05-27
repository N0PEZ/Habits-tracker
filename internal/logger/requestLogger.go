package logger

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

func RequestLogger(requestID string, r *http.Request, log zerolog.Logger) ([]byte, error) {
	if r.Body == nil || r.Body == http.NoBody {
		return nil, nil
	}

	// Читаем тело
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Восстанавливаем тело для обработчика
	r.Body = io.NopCloser(bytes.NewReader(body))

	// Логируем (если тело не пустое)
	if len(body) > 0 {
		log.Debug().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Bytes("body", body).
			Msg("HTTP request body")
	}

	return body, nil
}
