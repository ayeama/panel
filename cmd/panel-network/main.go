package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ayeama/panel/internal/types"
)

type Unify struct {
	// client *http.Client
	host        string
	apiKey      string
	siteID      string
	forwardHost string
}

type WebhookHandler struct {
	unify *Unify
}

func NewWebhookHandler(unify *Unify) WebhookHandler {
	return WebhookHandler{unify}
}

func (h *WebhookHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /", h.handleWebhook)
}

func (h *WebhookHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	var event types.WebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Fatal(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	switch event.Type {
	case types.WebhookEventInstanceCreated:
		var eventData types.WebhookEventDataInstanceCreated
		if err := json.Unmarshal(event.Data, &eventData); err != nil {
			log.Fatal(err)
		}

		type unifyRequestPortforward struct {
			Enabled               bool     `json:"enabled"`
			Name                  string   `json:"name"`
			PortforwardInterface  string   `json:"pfwd_interface"`
			DestinationIP         string   `json:"destination_ip"`
			DestinationIPs        []string `json:"destination_ips"`
			DestinationPort       string   `json:"dst_port"`
			ForwardIP             string   `json:"fwd"`
			ForwardPort           string   `json:"fwd_port"`
			Protocol              string   `json:"proto"`
			SourceLimitingEnabled bool     `json:"src_limiting_enabled"`
			Log                   bool     `json:"log"`
		}

		for _, port := range eventData.Ports {
			url := fmt.Sprintf(
				"https://%s/proxy/network/api/s/%s/rest/portforward",
				h.unify.host,
				h.unify.siteID,
			)

			body, err := json.Marshal(unifyRequestPortforward{
				Enabled:               true,
				Name:                  fmt.Sprintf("panel %s %s %s", eventData.Name, port, eventData.ID),
				PortforwardInterface:  "wan",
				DestinationIP:         "any",
				DestinationIPs:        make([]string, 0),
				DestinationPort:       port,
				ForwardIP:             h.unify.forwardHost,
				ForwardPort:           port,
				Protocol:              "tcp_udp",
				SourceLimitingEnabled: false,
				Log:                   false,
			})
			if err != nil {
				log.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Set("X-API-KEY", h.unify.apiKey)

			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			log.Println("opened port", port, eventData.Name, eventData.ID)
		}

	case types.WebhookEventInstanceDeleted:
		var eventData types.WebhookEventDataInstanceDeleted
		if err := json.Unmarshal(event.Data, &eventData); err != nil {
			log.Fatal(err)
		}

		type unifyResponsePortforward struct {
			Enabled               bool     `json:"enabled"`
			Name                  string   `json:"name"`
			PortforwardInterface  string   `json:"pfwd_interface"`
			DestinationIP         string   `json:"destination_ip"`
			DestinationIPs        []string `json:"destination_ips"`
			DestinationPort       string   `json:"dst_port"`
			ForwardIP             string   `json:"fwd"`
			ForwardPort           string   `json:"fwd_port"`
			Protocol              string   `json:"proto"`
			SourceLimitingEnabled bool     `json:"src_limiting_enabled"`
			Log                   bool     `json:"log"`

			ID     string `json:"_id"`
			SiteID string `json:"site_id"`
		}

		type unifyResponse struct {
			Data []unifyResponsePortforward `json:"data"`
		}

		url := fmt.Sprintf(
			"https://%s/proxy/network/api/s/%s/rest/portforward",
			h.unify.host,
			h.unify.siteID,
		)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("X-API-KEY", h.unify.apiKey)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var data unifyResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatal(err)
		}

		for _, rule := range data.Data {
			if strings.Contains(rule.Name, eventData.ID) {
				url = fmt.Sprintf(
					"https://%s/proxy/network/api/s/%s/rest/portforward/%s",
					h.unify.host,
					h.unify.siteID,
					rule.ID,
				)

				req, err := http.NewRequest(http.MethodDelete, url, nil)
				if err != nil {
					log.Fatal(err)
				}

				req.Header.Set("X-API-KEY", h.unify.apiKey)

				resp, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				log.Println("closed port", rule.ForwardPort, eventData.Name, eventData.ID)
			}
		}

	default:
		log.Println("WARNING unknown webhook event type")
	}
}

func main() {
	host := os.Getenv("PANEL_NETWORK_UNIFY_HOST")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_NETWORK_UNIFY_HOST' environment variable"))
	}

	apiKey := os.Getenv("PANEL_NETWORK_UNIFY_API_KEY")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_NETWORK_UNIFY_API_KEY' environment variable"))
	}

	siteID := os.Getenv("PANEL_NETWORK_UNIFY_SITE_ID")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_NETWORK_UNIFY_SITE_ID' environment variable"))
	}

	forwardHost := os.Getenv("PANEL_NETWORK_UNIFY_FORWARD_HOST")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_NETWORK_UNIFY_FORWARD_HOST' environment variable"))
	}

	unify := Unify{
		host:        host,
		apiKey:      apiKey,
		siteID:      siteID,
		forwardHost: forwardHost,
	}

	mux := http.NewServeMux()

	webhookHandler := NewWebhookHandler(&unify)
	webhookHandler.RegisterHandlers(mux)

	log.Println("starting")
	if err := http.ListenAndServe("0.0.0.0:8002", mux); err != nil {
		log.Fatal(err)
	}
}
