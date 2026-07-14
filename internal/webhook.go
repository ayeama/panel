package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/ayeama/panel/internal/runtime"
	"github.com/ayeama/panel/internal/types"
	"github.com/google/uuid"
)

func webhook(runtime runtime.Runtime) {
	// TODO better lifecycle management
	for {
		// TODO use podman events or our own
		events := make(chan types.Event)
		cancel := make(chan bool)

		// TODO missing cleanup?
		go func() {
			if err := runtime.Events(events, cancel); err != nil {
				log.Fatal(err)
			}
		}()

		// TODO standardise webhooks, support comma seperated
		for event := range events {
			switch event.Type {
			case types.EventTypeInstance:
				switch event.Action {
				case types.EventActionCreate:
					id := event.Actor.Attributes[types.InstanceLabelID]
					if id == "" {
						continue
					}

					instance, err := runtime.InstanceRead(id)
					if err != nil {
						log.Fatal(err)
					}

					// TODO move validation?
					_, err = url.ParseRequestURI(instance.Webhook)
					if err != nil {
						log.Println(err)
						continue
					}

					webhookData, err := json.Marshal(types.WebhookEventDataInstanceCreated{
						Instance: instance,
					})
					if err != nil {
						log.Fatal(err)
					}

					webhookRequest := types.WebhookEvent{
						ID:   uuid.NewString(),
						Type: types.WebhookEventInstanceCreated,
						Data: webhookData,
					}

					body, err := json.Marshal(webhookRequest)
					if err != nil {
						log.Fatal(err)
					}

					req, err := http.NewRequest(http.MethodPost, instance.Webhook, bytes.NewReader(body))
					if err != nil {
						log.Fatal(err)
					}

					req.Header.Set("Content-Type", "application/json")

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Println("WARNING", err.Error())
						break
					}
					defer resp.Body.Close()

					log.Println("sent webhook")
				case types.EventActionDelete:
					id := event.Actor.Attributes[types.InstanceLabelID]
					if id == "" {
						continue
					}

					webhook := event.Actor.Attributes[types.InstanceLabelWebhook]

					webhookData, err := json.Marshal(types.WebhookEventDataInstanceDeleted{
						ID: id,
					})
					if err != nil {
						log.Fatal(err)
					}

					webhookRequest := types.WebhookEvent{
						ID:   uuid.NewString(),
						Type: types.WebhookEventInstanceDeleted,
						Data: webhookData,
					}

					body, err := json.Marshal(webhookRequest)
					if err != nil {
						log.Fatal(err)
					}

					req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(body))
					if err != nil {
						log.Fatal(err)
					}

					req.Header.Set("Content-Type", "application/json")

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Println("WARNING", err.Error())
						break
					}
					defer resp.Body.Close()

					log.Println("sent webhook")
				default:
					break
				}
			default:
				break
			}
		}
	}
}
