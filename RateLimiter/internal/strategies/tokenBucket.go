package strategies

import (
	"fmt"
	"sync"
	"time"

	"github.com/jagjeet-singh-23/rate-limiter/internal/config"
	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/entities"
)

type UserRequestCtx struct {
	LastRefillAt time.Time
	Tokens       float32
}

type TokenBucket struct {
	mu     sync.RWMutex
	config config.RateLimitConfig
	redis  map[string]*UserRequestCtx
}

func NewTokenBucket(config config.RateLimitConfig) *TokenBucket {
	return &TokenBucket{
		config: config,
		redis:  make(map[string]*UserRequestCtx),
	}
}

func (b *TokenBucket) Allow(ctx *entities.RequestContext) (bool, time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()

	userCtx, ok := b.redis[ctx.UserID]
	if !ok {
		fmt.Printf("unauthorized: no user found: %s", ctx.UserID)
		return false, 0
	}

	userCurrentTokens := userCtx.Tokens
	lastRefillAt := userCtx.LastRefillAt

	timeSinceLastReq := float32(time.Since(lastRefillAt).Nanoseconds())
	tokensToAdd := timeSinceLastReq * b.config.RefillRate

	userCtx.Tokens = max(b.config.MaxCapacity, userCtx.Tokens+tokensToAdd)
	userCtx.LastRefillAt = time.Now()

	b.redis[ctx.UserID] = userCtx

	if userCurrentTokens >= 1 {
		userCtx.Tokens--
		b.redis[ctx.UserID] = userCtx
		return true, 0
	}

	tokensNeeded := 1.0 - userCtx.Tokens
	timeToWait := tokensNeeded / float32(b.config.RefillRate)

	return false, time.Duration(timeToWait)
}
