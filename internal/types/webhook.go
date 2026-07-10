package types

import "encoding/json"

type WebhookEventType string

const (
	WebhookEventInstanceCreated WebhookEventType = "instance.created"
	WebhookEventInstanceDeleted WebhookEventType = "instance.deleted"
)

type WebhookEvent struct {
	ID   string           `json:"uuid"`
	Type WebhookEventType `json:"type"`
	Data json.RawMessage  `json:"data"`
}

type WebhookEventDataInstanceCreated struct {
	Instance
}

type WebhookEventDataInstanceDeleted struct {
	ID string `json:"id"`
}
