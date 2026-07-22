package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Unify struct {
	// client *http.Client
	host   string
	apiKey string
	siteID string
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
	type unifyPortforwardResponse struct {
		ID            string `json:"_id"`
		SiteID        string `json:"site_id"`
		Enabled       bool   `json:"enabled"`
		Name          string `json:"name"`
		PfwdInterface string `json:"pfwd_interface"`
		DstPort       string `json:"dst_port"`
		Fwd           string `json:"fwd"`
		FwdPort       string `json:"fwd_port"`
		Proto         string `json:"proto"`
	}

	type unifyResponse struct {
		Data []unifyPortforwardResponse `json:"data"`
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

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

	log.Println("response:", data)
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

	unify := Unify{
		host:   host,
		apiKey: apiKey,
		siteID: siteID,
	}

	mux := http.NewServeMux()

	webhookHandler := NewWebhookHandler(&unify)
	webhookHandler.RegisterHandlers(mux)

	log.Println("starting")
	if err := http.ListenAndServe("0.0.0.0:8002", mux); err != nil {
		log.Fatal(err)
	}
}
