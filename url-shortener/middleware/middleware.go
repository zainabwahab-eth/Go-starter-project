package lMiddleware

import (
	"net"
	"net/http"
	"sync"
	"url-shortener/utils"

	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, exists := limiters[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(5), 10)
	limiters[ip] = limiter
	return limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if !getLimiter(ip).Allow() {
			utils.WriteJSON(w, http.StatusTooManyRequests, &utils.Response{Message: "rate limit exceeded"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
