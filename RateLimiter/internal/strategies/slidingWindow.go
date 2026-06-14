package strategies

import (
	"sync"
	"time"

	"github.com/jagjeet-singh-23/rate-limiter/internal/config"
	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/entities"
	"github.com/jagjeet-singh-23/rate-limiter/internal/utils/pq"
)

type SlidingWindow struct {
	pq     *pq.Heap[entities.RequestContext]
	config config.RateLimitConfig
	mu     sync.Mutex
}

func comparator(a, b entities.RequestContext) bool {
	return a.Timestamp.Before(b.Timestamp)
}

func NewSlidingWindow(config config.RateLimitConfig) *SlidingWindow {
	pq := pq.New(comparator)
	return &SlidingWindow{
		pq:     pq,
		config: config,
	}
}

func (sw *SlidingWindow) Allow(ctx *entities.RequestContext) (bool, time.Duration) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.clearExpiredRequests()

	if sw.pq.Len() >= int(sw.config.MaxCapacity) {
		oldest := sw.pq.Peek()

		retryAfter := time.Until(
			oldest.Timestamp.Add(sw.config.Window),
		)

		retryAfter = max(0, retryAfter)

		return false, retryAfter
	}

	sw.pq.Push(ctx)

	return true, 0
}

func (sw *SlidingWindow) clearExpiredRequests() {
	currentTimestamp := time.Now()

	for !sw.pq.Empty() {
		earliestRequest := sw.pq.Peek()

		if earliestRequest.Timestamp.Before(currentTimestamp.Add(-sw.config.Window)) {
			sw.pq.Pop()
		} else {
			break
		}

	}
}
