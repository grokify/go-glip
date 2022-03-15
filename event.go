package glip

import "time"

type GenericEvent struct {
	UUID           string    `json:"uuid,omitempty"`
	Event          string    `json:"event,omitempty"`
	Timestamp      time.Time `json:"timestamp,omitempty"`
	SubscriptionID string    `json:"subscriptionId,omitempty"`
	OwnerID        string    `json:"ownerId,omitempty"`
}

type GlipEvent struct {
	UUID           string          `json:"uuid,omitempty"`
	Event          string          `json:"event,omitempty"`
	Timestamp      time.Time       `json:"timestamp,omitempty"`
	SubscriptionID string          `json:"subscriptionId,omitempty"`
	OwnerID        string          `json:"ownerId,omitempty"`
	Body           TextMessageBody `json:"body,omitempty"`
}

type TextMessageBody struct {
	ID               string             `json:"id,omitempty"`
	GroupID          string             `json:"groupId,omitempty"`
	Type             string             `json:"type,omitempty"`
	Text             string             `json:"text,omitempty"`
	CreatorID        string             `json:"creatorId,omitempty"`
	CreationTime     time.Time          `json:"creationTime,omitempty"`
	LastModifiedTime time.Time          `json:"lastModifiedTime,omitempty"`
	Mentions         []GlipEventMention `json:"mentions,omitempty"`
	EventType        string             `json:"eventType,omitempty"`
}

type GlipEventMention struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}
