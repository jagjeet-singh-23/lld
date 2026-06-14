package middleware

import (
	"fmt"
	"net/http"

	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/entities"
	"github.com/jagjeet-singh-23/rate-limiter/internal/shared/interfaces"
)

type RateLimitMiddleware struct {
	next    http.Handler
	limiter interfaces.IRateLimiter
}

func NewRateLimitMiddleware(next http.Handler, limiter interfaces.IRateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		next:    next,
		limiter: limiter,
	}
}

func (m *RateLimitMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the user/client identifier (e.g., from a header or RemoteAddr)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = r.RemoteAddr
	}

	// 2. Initialize the RequestContext
	ctx, err := entities.NewRequestContext(userID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 3. Call Allow on the RateLimiter passing the context
	allowed, retryAfter := m.limiter.Allow(ctx)
	if !allowed {
		w.Header().Set("Retry-After", fmt.Sprintf("%.2f", retryAfter.Seconds()))
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	// 4. Call the next handler in the chain
	m.next.ServeHTTP(w, r)
}
