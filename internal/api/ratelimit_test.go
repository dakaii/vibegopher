package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIPRateLimiterBlocksBurst(t *testing.T) {
	lim := newIPRateLimiter(100, 2)
	ok := 0
	blocked := 0
	for i := 0; i < 5; i++ {
		if lim.allow("1.2.3.4") {
			ok++
		} else {
			blocked++
		}
	}
	if ok != 2 || blocked != 3 {
		t.Fatalf("got ok=%d blocked=%d, want ok=2 blocked=3", ok, blocked)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	lim := newIPRateLimiter(100, 1)
	h := lim.middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/posts", nil)
	req.RemoteAddr = "10.0.0.1:1234"

	rr1 := httptest.NewRecorder()
	h.ServeHTTP(rr1, req)
	if rr1.Code != http.StatusNoContent {
		t.Fatalf("first request: %d", rr1.Code)
	}

	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: %d", rr2.Code)
	}
}

func TestClientIPUsesRightmostXForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	req.Header.Set("X-Forwarded-For", "1.1.1.1, 2.2.2.2, 9.9.9.9")
	if got := clientIP(req); got != "9.9.9.9" {
		t.Fatalf("clientIP = %q, want rightmost 9.9.9.9", got)
	}
}

func TestClientIPFallsBackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	if got := clientIP(req); got != "10.0.0.1" {
		t.Fatalf("clientIP = %q, want 10.0.0.1", got)
	}
}

func TestIPRateLimiterPrunesStaleBuckets(t *testing.T) {
	lim := newIPRateLimiter(100, 1)
	lim.ttl = time.Millisecond
	lim.allow("stale")
	time.Sleep(5 * time.Millisecond)
	lim.mu.Lock()
	lim.pruneLocked(time.Now())
	n := len(lim.buckets)
	lim.mu.Unlock()
	if n != 0 {
		t.Fatalf("expected stale bucket pruned, got %d keys", n)
	}
}

func TestIPRateLimiterEvictsWhenAtMaxKeys(t *testing.T) {
	lim := newIPRateLimiter(100, 1)
	lim.maxKeys = 2
	lim.allow("a")
	lim.allow("b")
	lim.allow("c") // should evict oldest so map stays bounded
	lim.mu.Lock()
	n := len(lim.buckets)
	lim.mu.Unlock()
	if n > 2 {
		t.Fatalf("expected <=2 buckets after eviction, got %d", n)
	}
}
