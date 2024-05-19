package main

import (
	"log"
	"net/http"

	"github.com/walterlicinio/ratelimiter/config"
	"github.com/walterlicinio/ratelimiter/db"
	"github.com/walterlicinio/ratelimiter/handlers"
	"github.com/walterlicinio/ratelimiter/limiter"
	"github.com/walterlicinio/ratelimiter/middleware"
)

func main() {
	cfg := config.LoadConfig()

	store := db.NewRedisStore(cfg.RedisAddr, cfg.RedisPassword)
	rl := limiter.NewRateLimiter(store, cfg.RateLimitIP, cfg.RateLimitToken, cfg.AllowedTokens)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomeHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", middleware.RateLimiterMiddleware(rl)(mux)); err != nil {
		log.Fatal(err)
	}
}
