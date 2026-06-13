package order

import (
	"fmt"
	"time"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/order/enums"
)

type Command interface {
	Execute()
}

type TransitionCommand struct {
	orderID   string
	fromState enums.State
	toState   enums.State
	timestamp string
	auditRepo *TransitionAuditRepo
}

func NewTransitionCommand(orderID string, fromState, toState enums.State, auditRepo *TransitionAuditRepo) *TransitionCommand {
	return &TransitionCommand{
		orderID:   orderID,
		fromState: fromState,
		toState:   toState,
		timestamp: time.Now().Format("2006-01-02 15:04:05"),
		auditRepo: auditRepo,
	}
}

func (c *TransitionCommand) Execute() {
	c.auditRepo.Record(c.orderID, c.fromState, c.toState)
	fmt.Printf("%s [INFO] Order: %s transitioned from :%s to state: %s", c.timestamp, c.orderID, c.fromState, c.toState)
}
