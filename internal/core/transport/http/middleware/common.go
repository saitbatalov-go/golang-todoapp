package core_http_middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			if requestID == "" {
				requestID = uuid.New().String()
			}
			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID:= r.Header.Get(requestIDHeader)

			l:= log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String())
			)

			ctx:=context.WithValue(r.Context(), "logger", l)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
