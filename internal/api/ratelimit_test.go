package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
