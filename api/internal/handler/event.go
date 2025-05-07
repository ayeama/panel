package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

type EventHandler struct {
	serverService *service.ServerService
}

func NewEventHandler(serverService *service.ServerService) *EventHandler {
	return &EventHandler{serverService: serverService}
}

func (h *EventHandler) Events(w http.ResponseWriter, r *http.Request) {
	c := Upgrade(w, r)
	ctx, cancel := context.WithCancel(r.Context())

	ticker := time.NewTicker(time.Second * 5)

	go func() {
		for {
			var buf []byte
			_, err := c.Read(buf)
			if errors.Is(err, io.EOF) {
				cancel()
				break
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				p := domain.Pagination{Limit: 100, Offset: 0}
				paginatedServers := h.serverService.Read(p)

				for _, server := range paginatedServers.Items {
					event := types.EventServerStatus{Id: server.Id, Status: server.Status}
					data, err := json.Marshal(event)
					if err != nil {
						panic(err)
					}
					c.Write(data)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()
}

func (h *EventHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("GET /events", h.Events)
}
