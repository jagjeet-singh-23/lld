package channels

import (
	"fmt"

	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification/enums"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification/interfaces"
)

type ChannelRegistry struct {
	channels map[enums.Channel]interfaces.Channel
}

func NewChannelRegistry() *ChannelRegistry {
	registeredChannels := map[enums.Channel]interfaces.Channel{
		enums.SMS:  NewSmsChannel(),
		enums.Push: NewPushNotificationChannel(),
	}
	return &ChannelRegistry{
		channels: registeredChannels,
	}
}

func (r *ChannelRegistry) GetAllChannels() []enums.Channel {
	var allChannels []enums.Channel
	for k := range r.channels {
		allChannels = append(allChannels, k)
	}
	return allChannels
}

func (r *ChannelRegistry) Notify(channel enums.Channel, ctx interfaces.MessageContext) error {
	if _, ok := r.channels[channel]; !ok {
		return fmt.Errorf("unable to publish a notification, err: no such channel exists: %s", channel)
	}

	notificationChannel := r.channels[channel]
	notificationChannel.Notify(ctx)
	return nil
}
