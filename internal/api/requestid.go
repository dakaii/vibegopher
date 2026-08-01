package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ctxKey string

const requestIDKey ctxKey = "requestID"

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", id)

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))

		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration_ms=%d",
			id, r.Method, r.URL.Path, rec.status, time.Since(start).Milliseconds(),
		)
	})
}

func requestIDFrom(r *http.Request) string {
	if v, ok := r.Context().Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
