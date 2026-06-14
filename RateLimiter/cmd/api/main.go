package main

import (
	"log"
	"net/http"

	"github.com/jagjeet-singh-23/rate-limiter/internal"
	"github.com/jagjeet-singh-23/rate-limiter/internal/config"
	"github.com/jagjeet-singh-23/rate-limiter/internal/middleware"
	"github.com/jagjeet-singh-23/rate-limiter/internal/strategies"
)

func main() {
	// 1. Define the rate limiter configuration
	cfg := config.RateLimitConfig{
		MaxCapacity: 10.0, // Maximum capacity of tokens per user
		RefillRate:  1.0,  // Refill rate (e.g., tokens per nanosecond or second, depending on implementation)
	}

	// 2. Initialize the desired rate limiting strategy via the factory
	strategy, err := strategies.NewRateLimiterFactory("token_bucket", cfg)
	if err != nil {
		log.Fatalf("failed to initialize strategy: %v", err)
	}

	// 3. Initialize the RateLimiter service
	limiter := internal.NewRateLimiter(strategy)

	// 4. Define your API handler/handler function
	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access granted to the API resource!\n"))
	})

	// 5. Wrap your API handler with the RateLimitMiddleware
	rateLimitedHandler := middleware.NewRateLimitMiddleware(apiHandler, limiter)

	// 6. Register routes and start the HTTP server
	http.Handle("/api", rateLimitedHandler)

	log.Println("Server is running on port :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
