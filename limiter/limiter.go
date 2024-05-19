package limiter

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/walterlicinio/ratelimiter/db"
)

type RateLimiter struct {
	Store          db.Store
	RateLimitIP    int
	RateLimitToken map[string]int
	BlockTime      time.Duration
	mu             sync.Mutex
}

func NewRateLimiter(store db.Store, rateLimitIP int, blockTimeSeconds int, rateLimitToken map[string]int) *RateLimiter {
	return &RateLimiter{
		Store:          store,
		RateLimitIP:    rateLimitIP,
		RateLimitToken: rateLimitToken,
		BlockTime:      time.Duration(blockTimeSeconds) * time.Second,
	}
}

func (rl *RateLimiter) AllowRequest(ctx context.Context, identifier string, isToken bool) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limit := rl.RateLimitIP
	if isToken {
		if l, exists := rl.RateLimitToken[identifier]; exists {
			limit = l
		} else {
			limit = rl.RateLimitIP
		}
	}

	if rl.IsBlocked(ctx, identifier) {
		return false
	}

	count, err := rl.Store.Incr(ctx, identifier)
	if err != nil {
		return false
	}

	if count == 1 {
		rl.Store.Expire(ctx, identifier, time.Second)
	}

	if count > int64(limit) {
		rl.Store.Set(ctx, fmt.Sprintf("block:%s", identifier), "blocked", rl.BlockTime)
		return false
	}

	return true
}

func (rl *RateLimiter) IsBlocked(ctx context.Context, identifier string) bool {
	val, err := rl.Store.Get(ctx, fmt.Sprintf("block:%s", identifier))
	return err == nil && val == "blocked"
}

func (rl *RateLimiter) Reset(ctx context.Context, identifier string) {
	rl.Store.Set(ctx, identifier, "0", 0)
	rl.Store.Set(ctx, fmt.Sprintf("block:%s", identifier), "", 0)
}

func (rl *RateLimiter) Count(ctx context.Context, identifier string) int {
	count, err := rl.Store.Get(ctx, identifier)
	if err != nil {
		return -1
	}

	val, err := strconv.Atoi(count)
	if err != nil {
		return -1
	}

	return val
}
