package order

import (
	"time"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"
)

type TransitionAudit struct {
	OrderID   string
	FromState enums.State
	ToState   enums.State
	Timestamp time.Time
}

type TransitionAuditRepo struct {
	history map[string][]TransitionAudit
}

func NewTransitionAuditRepo() *TransitionAuditRepo {
	return &TransitionAuditRepo{
		history: make(map[string][]TransitionAudit),
	}
}

func (r *TransitionAuditRepo) Record(orderID string, fromState, toState enums.State) {
	r.history[orderID] = append(r.history[orderID], TransitionAudit{
		OrderID:   orderID,
		FromState: fromState,
		ToState:   toState,
		Timestamp: time.Now(),
	})
}

func (r *TransitionAuditRepo) History(orderID string) []TransitionAudit {
	history := r.history[orderID]
	result := make([]TransitionAudit, len(history))
	copy(result, history)
	return result
}
