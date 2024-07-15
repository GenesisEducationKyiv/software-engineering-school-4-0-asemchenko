package messaging

import (
	"currency-notifier/internal/shared/api"
	"encoding/json"
	"github.com/emicklei/go-restful/v3/log"
	"github.com/google/uuid"
	"sync"
)

type SubscriptionStatusListener interface {
	Consume(message string)
	SubscribeOnUser(userId uuid.UUID) *StatusChannel
	CloseSubscription(channel *StatusChannel)
}

type subscriptionStatusListener struct {
	userIdToChannels map[uuid.UUID]map[uuid.UUID]*StatusChannel
	rwMutex          *sync.RWMutex
}

type StatusChannel struct {
	UserId    uuid.UUID
	ChannelId uuid.UUID
	Channel   chan api.SubscriptionStatusMessage
}

func NewSubscriptionStatusListener() SubscriptionStatusListener {
	return &subscriptionStatusListener{
		userIdToChannels: make(map[uuid.UUID]map[uuid.UUID]*StatusChannel),
		rwMutex:          &sync.RWMutex{},
	}
}

func (c *subscriptionStatusListener) Consume(message string) {
	var statusMessage api.SubscriptionStatusMessage
	err := json.Unmarshal([]byte(message), &statusMessage)
	if err != nil {
		return
	}

	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	if channels, ok := c.userIdToChannels[statusMessage.UserID]; ok {
		for _, channel := range channels {
			channel.Channel <- statusMessage
		}
	}
	log.Printf("Notified all consumer on user %v about status %v", statusMessage.UserID, statusMessage.Status)
}

func (c *subscriptionStatusListener) SubscribeOnUser(userId uuid.UUID) *StatusChannel {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	channel := &StatusChannel{
		UserId:    userId,
		ChannelId: uuid.New(),
		Channel:   make(chan api.SubscriptionStatusMessage),
	}

	if _, ok := c.userIdToChannels[userId]; !ok {
		c.userIdToChannels[userId] = make(map[uuid.UUID]*StatusChannel)
	}
	c.userIdToChannels[userId][channel.ChannelId] = channel

	return channel
}

func (c *subscriptionStatusListener) CloseSubscription(channel *StatusChannel) {
	userIdChannels := c.userIdToChannels[channel.UserId]
	delete(userIdChannels, channel.ChannelId)
}
