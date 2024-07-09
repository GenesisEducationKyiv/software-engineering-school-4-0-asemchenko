package api

import "github.com/google/uuid"

const (
	MessageTypeRateUpdate         = "rate_update"
	MessageTypeAddSubscription    = "add_subscription"
	MessageTypeRemoveSubscription = "remove_subscription"
)

type Message struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type RateUpdateMessage struct {
	Currencies map[string]float64 `json:"currencies"`
}

type AddSubscriptionMessage struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

type RemoveSubscriptionMessage struct {
	UserID uuid.UUID `json:"user_id"`
}
