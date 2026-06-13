package interfaces

import "time"

type MessageContext struct {
	Message   string
	UserID    string
	CreatedAt time.Time
}
