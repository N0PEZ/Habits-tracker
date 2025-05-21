package logger

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func RequestLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем RequestID из контекста
			requestID := middleware.GetReqID(r.Context())

			// Читаем и восстанавливаем тело запроса
			var requestBody []byte
			if r.Body != nil {
				requestBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			// Подготовка для захвата ответа
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			var responseBuf bytes.Buffer
			ww.Tee(&responseBuf)

			// Замер времени выполнения
			start := time.Now()
			defer func() {
				// Формируем полную запись лога
				logger.Info().
					Str("request_id", requestID).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("query", r.URL.RawQuery).
					Str("ip", r.RemoteAddr).
					Str("user_agent", r.UserAgent()).
					Bytes("request_body", requestBody).
					Int("status", ww.Status()).
					Bytes("response_body", responseBuf.Bytes()).
					Dur("duration", time.Since(start)).
					Msg("request")
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
