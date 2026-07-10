package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ayeama/panel/internal/types"
	"github.com/google/uuid"
	mobyEvents "github.com/moby/moby/api/types/events"
	"go.podman.io/podman/v6/pkg/bindings/system"
	podmanTypes "go.podman.io/podman/v6/pkg/domain/entities/types"
)

func webhook(runtime *context.Context) {
	// TODO better lifecycle management
	for {
		// TODO use podman events or our own
		events := make(chan podmanTypes.Event)
		cancel := make(chan bool)
		if err := system.Events(*runtime, events, cancel, nil); err != nil {
			log.Fatal(err)
		}

		for event := range events {
			switch event.Type {
			case mobyEvents.ContainerEventType:
				switch event.Action {
				case mobyEvents.ActionCreate:
					webhookData, err := json.Marshal(types.WebhookEventDataInstanceCreated{
						InstanceID:   event.Actor.Attributes["com.github.ayeama.panel.instance.id"],
						InstanceName: event.Actor.Attributes["name"],
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

					req, err := http.NewRequest(http.MethodPost, "http://localhost:8001/webhook", bytes.NewReader(body))
					if err != nil {
						log.Fatal(err)
					}

					req.Header.Set("Content-Type", "application/json")

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Fatal(err)
					}
					defer resp.Body.Close()
				default:
					break
				}
			default:
				break
			}

			// fmt.Println("event:", event.Actor.ID, event.Type, event.Action, event.Actor.Attributes)
		}
	}
}
