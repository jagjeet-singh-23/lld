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
	case "leaky_bucket":
		return NewLeakyBucket(config), nil
	case "sliding_window":
		return NewSlidingWindow(config), nil
	default:
		return nil, fmt.Errorf("no such strategy found: %s", strategy)
	}
}
