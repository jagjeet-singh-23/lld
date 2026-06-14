package strategies

import (
	"fmt"
	"sync"
	"time"

	"github.com/jagjeet-singh-23/rate-limiter/internal/config"
	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/entities"
)

type LeakyBucket struct {
	requestQueue  []*entities.RequestContext
	config        config.RateLimitConfig
	mu            sync.RWMutex
	lastCleanupAt time.Time
}

func NewLeakyBucket(config config.RateLimitConfig) *LeakyBucket {
	return &LeakyBucket{
		requestQueue: make([]*entities.RequestContext, int(config.MaxCapacity)),
		config:       config,
	}
}

func (b *LeakyBucket) Allow(ctx *entities.RequestContext) (bool, time.Duration) {
	b.cleanupExpiredRequests()
	allowed, retryAfter := b.canAcceptRequest()
	if !allowed {
		return allowed, retryAfter
	}

	err := b.enqueue(ctx)
	if err != nil {
		fmt.Printf("unable to enqueue the request: %s", ctx.RequestID)
		return false, 0
	}

	return true, 0
}

func (b *LeakyBucket) cleanupExpiredRequests() {
	now := time.Now()
	b.mu.Lock()
	defer b.mu.Unlock()

	cutoffTime := now.Add(-time.Duration(1.0/b.config.RefillRate) * time.Second)

	newQueue := []*entities.RequestContext{}
	for _, req := range b.requestQueue {
		if req.Timestamp.After(cutoffTime) {
			newQueue = append(newQueue, req)
		}
	}

	b.requestQueue = newQueue
}

func (b *LeakyBucket) canAcceptRequest() (bool, time.Duration) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	queueSize := len(b.requestQueue)
	if queueSize >= int(b.config.MaxCapacity) {
		leakRate := b.config.RefillRate
		neededSpace := float32(b.config.MaxCapacity - float32(queueSize))
		retryAfter := neededSpace / leakRate
		return false, time.Duration(retryAfter * float32(time.Second))
	}

	return true, 0
}

func (b *LeakyBucket) enqueue(ctx *entities.RequestContext) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.requestQueue = append(b.requestQueue, ctx)
	return nil
}
