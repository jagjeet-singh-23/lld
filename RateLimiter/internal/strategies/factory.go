package strategies

import (
	"fmt"

	"github.com/jagjeet-singh-23/rate-limiter/internal/config"
	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/interfaces"
)

func NewRateLimiterFactory(strategy string, config config.RateLimitConfig) (interfaces.IRateLimitStrategy, error) {
	switch strategy {
	case "token_bucket":
		return NewTokenBucket(config), nil
	default:
		return nil, fmt.Errorf("no such strategy found: %s", strategy)
	}
}
