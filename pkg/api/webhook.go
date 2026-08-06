package api

import (
	"encoding/json"
)

type WebhookEventType string

const (
	WebhookEventInstanceCreated WebhookEventType = "instance.created"
	WebhookEventInstanceDeleted WebhookEventType = "instance.deleted"
)

// TODO update api json structure
type WebhookEvent struct {
	ID   string           `json:"uuid"`
	Type WebhookEventType `json:"type"`
	Data json.RawMessage  `json:"data"`
}

type WebhookEventDataInstanceCreated struct {
	Instance
}

type WebhookEventDataInstanceDeleted struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
