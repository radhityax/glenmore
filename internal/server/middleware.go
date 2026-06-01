package server

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

func ActivityPubContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request ){
		if strings.HasPrefix(r.URL.Path, "/actor/") ||
		strings.HasPrefix(r.URL.Path, "/.well-known/") {
			accept := r.Header.Get("Accept")
			if accept == "" || strings.Contains(accept, "*/*") {
			}
		}
		next.ServeHTTP(w,r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w,r)
		slog.Info("request",
		"method", r.Method,
		"path", r.URL.Path,
		"duration", time.Since(start).String(),
	)
})
}

type RateLimiter struct {
	mu sync.Mutex
	visitors map[string]*rateLimitEntry
	rate int
	burst int
}


type rateLimitEntry struct {
	tokens int
	lastCheck time.Time
}

func NewRateLimiter(rate, burst int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rateLimitEntry),
		rate: rate,
		burst: burst,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, ok := rl.visitors[ip]
	if !ok {
		entry = &rateLimitEntry{tokens: rl.burst, lastCheck: time.Now()}
		rl.visitors[ip] = entry
	}

	now := time.Now()
	elapsed := now.Sub(entry.lastCheck).Seconds()
	entry.tokens += int(elapsed * float64(rl.rate))
	if entry.tokens > rl.burst {
		entry.tokens = rl.burst
	}
	entry.lastCheck = now

	if entry.tokens > 0 {
		entry.tokens--
		return true
	}
	return false
}

func RateLimit(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if idx := strings.LastIndex(ip, ":"); idx != -1 {
				ip = ip[:idx]
			}
			if !rl.Allow(ip) {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
