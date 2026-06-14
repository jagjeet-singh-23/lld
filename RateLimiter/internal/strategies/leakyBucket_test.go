package strategies

import (
	"sync"
	"testing"
	"time"

	"github.com/jagjeet-singh-23/rate-limiter/internal/config"
	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/entities"
)

func TestLeakyBucket_Allow(t *testing.T) {
	cfg := config.RateLimitConfig{
		MaxCapacity: 3,
		RefillRate:  2, // 2 requests per second -> leak interval is 500ms
	}

	lb := NewLeakyBucket(cfg)

	// Create requests
	req1 := &entities.RequestContext{RequestID: "req1", UserID: "user1", Timestamp: time.Now()}
	req2 := &entities.RequestContext{RequestID: "req2", UserID: "user1", Timestamp: time.Now()}
	req3 := &entities.RequestContext{RequestID: "req3", UserID: "user1", Timestamp: time.Now()}
	req4 := &entities.RequestContext{RequestID: "req4", UserID: "user1", Timestamp: time.Now()}

	// 1. Fill the bucket
	if allowed, _ := lb.Allow(req1); !allowed {
		t.Error("Expected req1 to be allowed")
	}
	if allowed, _ := lb.Allow(req2); !allowed {
		t.Error("Expected req2 to be allowed")
	}
	if allowed, _ := lb.Allow(req3); !allowed {
		t.Error("Expected req3 to be allowed")
	}

	// 2. The bucket is now full. The 4th request should be rejected.
	allowed, retryAfter := lb.Allow(req4)
	if allowed {
		t.Error("Expected req4 to be rejected because the bucket is full")
	}
	if retryAfter <= 0 {
		t.Errorf("Expected positive retryAfter, got %v", retryAfter)
	}

	// 3. Wait for the oldest request to leak (at least 500ms).
	time.Sleep(510 * time.Millisecond)

	// 4. Request 4 should now be allowed
	req4.Timestamp = time.Now()
	if allowed, _ = lb.Allow(req4); !allowed {
		t.Error("Expected req4 to be allowed after oldest request leaked")
	}
}

func TestLeakyBucket_Concurrent(t *testing.T) {
	cfg := config.RateLimitConfig{
		MaxCapacity: 100,
		RefillRate:  1000, // Very high leak rate
	}
	lb := NewLeakyBucket(cfg)

	var wg sync.WaitGroup
	numGoroutines := 10
	requestsPerGoroutine := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(gID int) {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				req := &entities.RequestContext{
					RequestID: "req",
					UserID:    "user1",
					Timestamp: time.Now(),
				}
				lb.Allow(req)
			}
		}(i)
	}
	wg.Wait()
}
