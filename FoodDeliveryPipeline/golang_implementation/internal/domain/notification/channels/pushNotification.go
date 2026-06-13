package channels

import (
	"fmt"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification/interfaces"
)

type PushNotificationChannel struct{}

func NewPushNotificationChannel() *PushNotificationChannel {
	return &PushNotificationChannel{}
}

func (c *PushNotificationChannel) Notify(ctx interfaces.MessageContext) {
	fmt.Printf("%s[INFO] User%s, update: %s", ctx.CreatedAt.Format("2006-01-02 15:04:02"), ctx.UserID, ctx.Message)
}
