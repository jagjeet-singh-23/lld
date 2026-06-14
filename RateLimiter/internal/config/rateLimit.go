package config

import "time"

type RateLimitConfig struct {
	MaxCapacity float32
	RefillRate  float32
	Window      time.Duration
}
