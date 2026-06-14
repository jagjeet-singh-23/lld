package notification

import (
	"fmt"

	notificationInterfaces "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification/interfaces"
	sharedInterfaces "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
)

type NotificationObserver struct {
	service *NotificationService
}

func NewNotificationObserver(service *NotificationService) *NotificationObserver {
	return &NotificationObserver{
		service: service,
	}
}

func (o *NotificationObserver) Update(event sharedInterfaces.Event) error {
	orderID, _ := event.Metadata["orderID"].(string)
	userID, ok := event.Metadata["userID"].(string)
	if !ok {
		// Fallback/dummy user ID if none is present in event metadata
		userID = "user_1"
	}

	ctx := notificationInterfaces.MessageContext{
		Message:   fmt.Sprintf("Order %s status updated: %s", orderID, event.EventType),
		UserID:    userID,
		CreatedAt: event.Timestamp,
	}

	o.service.Handle(ctx)
	return nil
}
