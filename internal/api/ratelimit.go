package api

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rateBucket struct {
	tokens float64
	last   time.Time
}

// ipRateLimiter is a simple per-IP token bucket (in-process; fine for a single Cloud Run instance).
type ipRateLimiter struct {
	mu      sync.Mutex
	rate    float64 // tokens per second
	burst   float64
	buckets map[string]*rateBucket
}

func newIPRateLimiter(ratePerSec float64, burst int) *ipRateLimiter {
	return &ipRateLimiter{
		rate:    ratePerSec,
		burst:   float64(burst),
		buckets: make(map[string]*rateBucket),
	}
}

func (l *ipRateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[key]
	if !ok {
		l.buckets[key] = &rateBucket{tokens: l.burst - 1, last: now}
		return true
	}
	elapsed := now.Sub(b.last).Seconds()
	b.tokens += elapsed * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (l *ipRateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.allow(clientIP(r)) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"Too many requests"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
