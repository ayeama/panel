package internal

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ayeama/panel/internal/runtime"
	"github.com/ayeama/panel/internal/types"
	"github.com/google/uuid"
)

func webhook(runtime runtime.Runtime) {
	// TODO better lifecycle management
	for {
		events := make(chan types.Event)
		cancel := make(chan bool)

		// TODO missing cleanup?
		go func() {
			if err := runtime.Events(events, cancel); err != nil {
				log.Println("about to fail in webhook events")
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

					for _, webhook := range instance.Webhooks {
						// TODO move validation?
						_, err = url.ParseRequestURI(webhook)
						if err != nil {
							log.Println(err)
							continue
						}

						webhookData, err := json.Marshal(types.WebhookEventDataInstanceCreated{
							Instance: instance,
						})
						if err != nil {
							log.Println("WARNING:", err)
							continue
						}

						webhookRequest := types.WebhookEvent{
							ID:   uuid.NewString(),
							Type: types.WebhookEventInstanceCreated,
							Data: webhookData,
						}

						body, err := json.Marshal(webhookRequest)
						if err != nil {
							log.Println("WARNING:", err)
							continue
						}

						req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(body))
						if err != nil {
							log.Println("WARNING:", err)
							continue
						}

						req.Header.Set("Content-Type", "application/json")

						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							log.Println("WARNING", err.Error())
							continue
						}
						resp.Body.Close()

						log.Println("sent webhook")
					}

				case types.EventActionDelete:
					id := event.Actor.Attributes[types.InstanceLabelID]
					if id == "" {
						continue
					}

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

					var webhooks []string
					if event.Actor.Attributes[types.InstanceLabelWebhooks] != "" {
						webhooks = strings.Split(event.Actor.Attributes[types.InstanceLabelWebhooks], ",")
					} else {
						webhooks = make([]string, 0)
					}

					for _, webhook := range webhooks {
						req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(body))
						if err != nil {
							log.Println("WARNING:", err)
							continue
						}

						req.Header.Set("Content-Type", "application/json")

						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							log.Println("WARNING", err.Error())
							continue
						}
						resp.Body.Close()

						log.Println("sent webhook")
					}
				}
			}
		}
	}
}
