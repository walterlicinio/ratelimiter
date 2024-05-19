package tests

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/walterlicinio/ratelimiter/db"
	"github.com/walterlicinio/ratelimiter/limiter"
)

func TestRateLimiter(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer s.Close()

	store := db.NewRedisStore(s.Addr(), "")
	allowedTokens := map[string]int{
		"token1": 3,
		"token2": 6,
	}
	rl := limiter.NewRateLimiter(store, 5, 1, allowedTokens)

	ctx := context.Background()
	identifierIP := "testIP"
	identifierToken := "token1"

	for i := 0; i < 5; i++ {
		if !rl.AllowRequest(ctx, identifierIP, false) {
			t.Fatalf("Request %d should have been allowed (IP)", i+1)
		}
	}

	if rl.AllowRequest(ctx, identifierIP, false) {
		t.Fatalf("6th request should have been blocked (IP)")
	}

	if !rl.IsBlocked(ctx, identifierIP) {
		t.Fatalf("IP should be blocked")
	}

	t.Logf("IP Request count before sleep: %d", rl.Count(ctx, identifierIP))

	time.Sleep(2 * time.Second)

	t.Logf("IP Request count after sleep: %d", rl.Count(ctx, identifierIP))

	rl.Reset(ctx, identifierIP)

	if !rl.AllowRequest(ctx, identifierIP, false) {
		t.Fatalf("Request after expiry should have been allowed (IP)")
	}

	if rl.IsBlocked(ctx, identifierIP) {
		t.Fatalf("IP should not be blocked after expiry")
	}

	for i := 0; i < 3; i++ {
		if !rl.AllowRequest(ctx, identifierToken, true) {
			t.Fatalf("Request %d should have been allowed (Token)", i+1)
		}
	}

	if rl.AllowRequest(ctx, identifierToken, true) {
		t.Fatalf("4th request should have been blocked (Token)")
	}

	if !rl.IsBlocked(ctx, identifierToken) {
		t.Fatalf("Token should be blocked")
	}

	t.Logf("Token Request count before sleep: %d", rl.Count(ctx, identifierToken))

	time.Sleep(2 * time.Second)

	t.Logf("Token Request count after sleep: %d", rl.Count(ctx, identifierToken))

	rl.Reset(ctx, identifierToken)

	if !rl.AllowRequest(ctx, identifierToken, true) {
		t.Fatalf("Request after expiry should have been allowed (Token)")
	}

	if rl.IsBlocked(ctx, identifierToken) {
		t.Fatalf("Token should not be blocked after expiry")
	}
}
