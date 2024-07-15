package api

import "github.com/google/uuid"

const (
	MessageTypeRateUpdate         = "rate_update"
	MessageTypeAddSubscription    = "add_subscription"
	MessageTypeRemoveSubscription = "remove_subscription"
	MessageTypeSubscriptionStatus = "subscription_status"
)

type Message struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type RateUpdateMessage struct {
	Currencies map[string]float64 `json:"currencies"`
}

type AddSubscriptionMessage struct {
	Email  string    `json:"email"`
	UserID uuid.UUID `json:"user_id"`
}

type RemoveSubscriptionMessage struct {
	UserID uuid.UUID `json:"user_id"`
}

// subscription status enum
const (
	SubscriptionStatusSuccessfullyCreated = "successfully_created"
	SubscriptionStatusFailedToCreate      = "failed_to_create"
)

type SubscriptionStatusMessage struct {
	UserID  uuid.UUID `json:"user_id"`
	Status  string    `json:"status"`
	Details string    `json:"details"`
}

func (m *SubscriptionStatusMessage) IsSuccessfullyCreated() bool {
	return m.Status == SubscriptionStatusSuccessfullyCreated
}
