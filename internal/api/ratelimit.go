package api

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	rateBucketTTL   = 10 * time.Minute
	rateBucketMax   = 10_000
	ratePruneEvery  = 100 // allow() calls between opportunistic prunes
)

type rateBucket struct {
	tokens float64
	last   time.Time
}

// ipRateLimiter is a simple per-IP token bucket (in-process; fine for a single Cloud Run instance).
type ipRateLimiter struct {
	mu       sync.Mutex
	rate     float64 // tokens per second
	burst    float64
	buckets  map[string]*rateBucket
	ops      uint64
	ttl      time.Duration
	maxKeys  int
}

func newIPRateLimiter(ratePerSec float64, burst int) *ipRateLimiter {
	return &ipRateLimiter{
		rate:    ratePerSec,
		burst:   float64(burst),
		buckets: make(map[string]*rateBucket),
		ttl:     rateBucketTTL,
		maxKeys: rateBucketMax,
	}
}

func (l *ipRateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.ops++
	if l.ops%ratePruneEvery == 0 {
		l.pruneLocked(now)
	}

	b, ok := l.buckets[key]
	if !ok {
		if len(l.buckets) >= l.maxKeys {
			l.pruneLocked(now)
			if len(l.buckets) >= l.maxKeys {
				l.evictOldestLocked()
			}
		}
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

func (l *ipRateLimiter) pruneLocked(now time.Time) {
	for k, b := range l.buckets {
		if now.Sub(b.last) > l.ttl {
			delete(l.buckets, k)
		}
	}
}

func (l *ipRateLimiter) evictOldestLocked() {
	var oldestKey string
	var oldestTime time.Time
	first := true
	for k, b := range l.buckets {
		if first || b.last.Before(oldestTime) {
			oldestKey = k
			oldestTime = b.last
			first = false
		}
	}
	if oldestKey != "" {
		delete(l.buckets, oldestKey)
	}
}

// clientIP returns a rate-limit key for the request.
// Cloud Run / GFE append the connecting client IP to X-Forwarded-For, so the
// rightmost entry is the platform-trusted hop. Leftmost is client-spoofable.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		for i := len(parts) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(parts[i])
			if ip != "" {
				return ip
			}
		}
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
