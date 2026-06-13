package notification

import (
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification/channels"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification/enums"
	"github.com/jagjeet-singh-23/food-delivery-pipeline/internal/domain/notification/interfaces"
	shared "github.com/jagjeet-singh-23/food-delivery-pipeline/internal/shared/interfaces"
)

type NotificationService struct {
	registry        *channels.ChannelRegistry
	userPreferences map[string][]enums.Channel
}

func NewNotificationService(bus *shared.EventBus) *NotificationService {
	return &NotificationService{
		registry: channels.NewChannelRegistry(),
	}
}

func (s *NotificationService) Handle(ctx interfaces.MessageContext) {
	userPreferedChannels := s.getUserPrefferedChannels(ctx.UserID)
	for _, channel := range userPreferedChannels {
		s.registry.Notify(channel, ctx)
	}
}

func (s *NotificationService) getUserPrefferedChannels(userID string) []enums.Channel {
	userPreferedChannels, ok := s.userPreferences[userID]
	if !ok {
		allChannels := s.registry.GetAllChannels()
		s.userPreferences[userID] = allChannels
		userPreferedChannels = allChannels
	}

	return userPreferedChannels
}
