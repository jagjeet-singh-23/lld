package channels

import (
	"fmt"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification/interfaces"
)

type SmsChannel struct{}

func NewSmsChannel() *SmsChannel {
	return &SmsChannel{}
}

func (c *SmsChannel) Notify(ctx interfaces.MessageContext) {
	fmt.Printf("%s[INFO] User%s, update: %s", ctx.CreatedAt.Format("2006-01-02 15:04:02"), ctx.UserID, ctx.Message)
}
