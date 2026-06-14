package internal

import (
	"time"

	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/entities"
	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/interfaces"
)

type RateLimiter struct {
	strategy interfaces.IRateLimitStrategy
}

func NewRateLimiter(strategy interfaces.IRateLimitStrategy) *RateLimiter {
	return &RateLimiter{
		strategy: strategy,
	}
}

func (rl *RateLimiter) Allow(ctx *entities.RequestContext) (bool, time.Duration) {
	allowed, retryAfter := rl.strategy.Allow(ctx)
	return allowed, retryAfter
}
