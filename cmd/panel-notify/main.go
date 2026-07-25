package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ayeama/panel/internal/types"
)

type Ntfy struct {
	host  string
	topic string
	token string
}

type WebhookHandler struct {
	ntfy *Ntfy
}

func NewWebhookHandler(ntfy *Ntfy) WebhookHandler {
	return WebhookHandler{ntfy}
}

func (h *WebhookHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /", h.handleWebhook)
}

func (h *WebhookHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	var event types.WebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Fatal(err)
	}

	switch event.Type {
	case types.WebhookEventInstanceCreated:
		var eventData types.WebhookEventDataInstanceCreated
		if err := json.Unmarshal(event.Data, &eventData); err != nil {
			log.Fatal(err)
		}

		msg := fmt.Sprintf("instance %s created", eventData.Name)
		url := fmt.Sprintf("%s/%s", h.ntfy.host, h.ntfy.topic)
		req, _ := http.NewRequest("POST", url, strings.NewReader(msg))
		req.Header.Set("Authorization", "Bearer "+h.ntfy.token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("WARNING:", err)
			return
		}
		defer resp.Body.Close()

		log.Println("handled notification")
	case types.WebhookEventInstanceDeleted:
		var eventData types.WebhookEventDataInstanceDeleted
		if err := json.Unmarshal(event.Data, &eventData); err != nil {
			log.Fatal(err)
		}

		msg := fmt.Sprintf("instance %s deleted", eventData.Name)
		url := fmt.Sprintf("%s/%s", h.ntfy.host, h.ntfy.topic)
		req, _ := http.NewRequest("POST", url, strings.NewReader(msg))
		req.Header.Set("Authorization", "Bearer "+h.ntfy.token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("WARNING:", err)
			return
		}
		defer resp.Body.Close()

		log.Println("handled notification")
	default:
		log.Print("ERROR: unknown webhook event type")
		return
	}
}

func main() {
	host := os.Getenv("PANEL_NOTIFY_NTFY_HOST")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_NOTIFY_NTFY_HOST' environment variable"))
	}

	topic := os.Getenv("PANEL_NOTIFY_NTFY_TOPIC")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_NOTIFY_NTFY_TOPIC' environment variable"))
	}

	token := os.Getenv("PANEL_NOTIFY_NTFY_TOKEN")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_NOTIFY_NTFY_TOKEN' environment variable"))
	}

	ntfy := Ntfy{
		host:  host,
		topic: topic,
		token: token,
	}

	mux := http.NewServeMux()

	webhookHandler := NewWebhookHandler(&ntfy)
	webhookHandler.RegisterHandlers(mux)

	log.Println("starting")
	if err := http.ListenAndServe("0.0.0.0:8003", mux); err != nil {
		log.Fatal(err)
	}
}
