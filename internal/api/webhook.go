package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/ayeama/panel/internal/runtime"
	"github.com/ayeama/panel/internal/types"
	"github.com/ayeama/panel/pkg/api"
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

		for e := range events {
			switch e.Type {
			case types.EventTypeInstanceCreated:
				eventData := e.Data.(types.EventInstanceCreated)

				id := eventData.ID
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

					webhookData, err := json.Marshal(api.WebhookEventDataInstanceCreated{
						Instance: api.Instance{
							ID:     instance.ID,
							Name:   instance.Name,
							Image:  instance.Image,
							Status: instance.Status,
							Ports:  instance.Ports,
							Resources: api.InstanceResources{
								CPU:    instance.Resources.CPU,
								Memory: instance.Resources.Memory,
								Disk:   instance.Resources.Disk,
							},
							Webhooks: instance.Webhooks,
						},
					})
					if err != nil {
						log.Println("WARNING:", err)
						continue
					}

					webhookRequest := api.WebhookEvent{
						ID:   uuid.NewString(),
						Type: api.WebhookEventInstanceCreated,
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

			case types.EventTypeInstanceDeleted:
				eventData := e.Data.(types.EventInstanceDeleted)

				id := eventData.ID
				if id == "" {
					continue
				}

				name := eventData.Name

				webhookData, err := json.Marshal(api.WebhookEventDataInstanceDeleted{
					ID:   id,
					Name: name,
				})
				if err != nil {
					log.Fatal(err)
				}

				webhookRequest := api.WebhookEvent{
					ID:   uuid.NewString(),
					Type: api.WebhookEventInstanceDeleted,
					Data: webhookData,
				}

				body, err := json.Marshal(webhookRequest)
				if err != nil {
					log.Fatal(err)
				}

				webhooks := eventData.Webhooks

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
