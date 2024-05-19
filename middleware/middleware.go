package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/walterlicinio/ratelimiter/limiter"
)

func RateLimiterMiddleware(rl *limiter.RateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ip := strings.Split(r.RemoteAddr, ":")[0]
			token := r.Header.Get("API_KEY")

			if token != "" {
				if _, exists := rl.RateLimitToken[token]; exists {
					if rl.IsBlocked(ctx, token) || !rl.AllowRequest(ctx, token, true) {
						http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame (token)", http.StatusTooManyRequests)
						return
					}
				}
			} else {
				if rl.IsBlocked(ctx, ip) || !rl.AllowRequest(ctx, ip, false) {
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame (IP)", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
