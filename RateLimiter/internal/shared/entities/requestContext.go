package entities

import (
	"time"

	"github.com/jagjeet-singh-23/rate-limiter/internal/utils"
)

type RequestContext struct {
	RequestID string
	UserID    string
	Timestamp time.Time
}

func NewRequestContext(userID string) (*RequestContext, error) {
	requestID, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	return &RequestContext{
		RequestID: requestID,
		UserID:    userID,
		Timestamp: time.Now(),
	}, nil
}
