// Package middleware - logging.go
package middleware

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"
	"time"
)

type contextKey string

var requestIDKey = contextKey("requestID")

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// LoggingMiddleware является обработчиком middleware, ведущим журнал запросов.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// ошибка получения IP-адреса клиента
			clientIP = "unknown"
		}

		log.Printf("Request: %s %s (ID: %s, IP: %s, Time: %s)",
			r.Method, r.URL.Path, requestID, clientIP, time.Now().Format(time.RFC3339))

		sw := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r.WithContext(ctx))

		log.Printf("Response status: %d (Request ID: %s)", sw.status, requestID)
	})
}

// RequestIDMiddleware является обработчиком middleware, устанавливающим уникальный идентификатор запроса.
func RequestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)

		next(w, r.WithContext(ctx))
	})
}
