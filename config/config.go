package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr        string
	RedisPassword    string
	RateLimitIP      int
	RateLimitToken   int
	BlockTimeSeconds int
	AllowedTokens    map[string]int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	rateLimitIP, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	rateLimitToken, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))
	blockTimeSeconds, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_SECONDS"))
	allowedTokens, _ := os.Open("tokens.json")
	defer allowedTokens.Close()
	tokens := make(map[string]int)
	json.NewDecoder(allowedTokens).Decode(&tokens)

	return &Config{
		RedisAddr:        os.Getenv("REDIS_ADDR"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RateLimitIP:      rateLimitIP,
		RateLimitToken:   rateLimitToken,
		BlockTimeSeconds: blockTimeSeconds,
		AllowedTokens:    tokens,
	}
}
