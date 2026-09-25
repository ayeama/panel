package podman

import (
	"fmt"
	"strings"

	"github.com/ayeama/panel/internal/runtime"
	"github.com/ayeama/panel/internal/types"
	mobyEvents "github.com/moby/moby/api/types/events"
	"go.podman.io/podman/v6/pkg/bindings/system"
	entitiesTypes "go.podman.io/podman/v6/pkg/domain/entities/types"
)

func (r *Runtime) Events(events chan types.Event, cancel chan bool) error {
	podmanEvents := make(chan entitiesTypes.Event)
	if err := system.Events(*r.ctx, podmanEvents, cancel, nil); err != nil {
		return &runtime.Error{Op: "read", Resource: "events", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	for podmanEvent := range podmanEvents {
		switch podmanEvent.Type {
		case mobyEvents.ContainerEventType:
			switch podmanEvent.Action {
			case mobyEvents.ActionCreate:
				events <- types.Event{
					Type: types.EventTypeInstanceCreated,
					Data: types.EventInstanceCreated{
						ID: podmanEvent.Actor.Attributes[instanceLabelID],
					},
				}
			case mobyEvents.ActionRemove:
				var webhooks []string
				if podmanEvent.Actor.Attributes[instanceLabelWebhooks] != "" {
					webhooks = strings.Split(podmanEvent.Actor.Attributes[instanceLabelWebhooks], ",")
				} else {
					webhooks = make([]string, 0)
				}

				events <- types.Event{
					Type: types.EventTypeInstanceDeleted,
					Data: types.EventInstanceDeleted{
						ID:       podmanEvent.Actor.Attributes[instanceLabelID],
						Name:     podmanEvent.Actor.Attributes["name"],
						Webhooks: webhooks,
					},
				}
			}
		}
	}

	return nil
}
