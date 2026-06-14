package interfaces

import (
	"time"

	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/entities"
)

type IRateLimitStrategy interface {
	Allow(ctx *entities.RequestContext) (bool, time.Duration)
}
