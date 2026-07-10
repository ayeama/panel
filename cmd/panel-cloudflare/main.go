package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/ayeama/panel/internal/types"
	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
)

type Cloudflare struct {
	client *cloudflare.Client
	zoneID string

	host    string
	comment string
}

func (cf *Cloudflare) subdomainName(instanceName string) string {
	return instanceName + "." + cf.host
}

type WebhookHandler struct {
	cf *Cloudflare
}

func NewWebhookHandler(cf *Cloudflare) WebhookHandler {
	return WebhookHandler{cf}
}

func (h *WebhookHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/webhook", h.handleWebhook)
}

func (h *WebhookHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

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

		// TODO make async
		name := h.cf.subdomainName(eventData.InstanceName)
		ipaddresses, err := net.LookupIP(h.cf.host)
		if err != nil {
			log.Fatal(err)
		}
		content := ipaddresses[0].String()

		_, err = (*h.cf.client).DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(h.cf.zoneID),
			Body: dns.ARecordParam{
				Type:    cloudflare.F(dns.ARecordTypeA),
				Name:    cloudflare.F(name),
				Content: cloudflare.F(content),
				TTL:     cloudflare.F(dns.TTL(60)),
				Comment: cloudflare.F(h.cf.comment),
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = (*h.cf.client).DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(h.cf.zoneID),
			Body: dns.SRVRecordParam{
				Type: cloudflare.F(dns.SRVRecordTypeSRV),
				Name: cloudflare.F("_minecraft._tcp." + name),
				Data: cloudflare.F(dns.SRVRecordDataParam{
					Priority: cloudflare.F(float64(0)),
					Weight:   cloudflare.F(float64(5)),
					Port:     cloudflare.F(float64(25565)),
					Target:   cloudflare.F(name),
				}),
				TTL:     cloudflare.F(dns.TTL(60)),
				Comment: cloudflare.F(h.cf.comment),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown webhook event type")
	}
}

func main() {
	host := os.Getenv("PANEL_CLOUDFLARE_HOST")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_CLOUDFLARE_HOST' environment variable"))
	}

	comment := os.Getenv("PANEL_CLOUDFLARE_COMMENT")
	if comment == "" {
		comment = "managed by panel"
	}

	zoneID := os.Getenv("PANEL_CLOUDFLARE_ZONE_ID")
	if zoneID == "" {
		log.Fatal(errors.New("missing 'PANEL_CLOUDFLARE_ZONE_ID' environment variable"))
	}

	client := cloudflare.NewClient()

	cf := Cloudflare{
		client:  client,
		zoneID:  zoneID,
		host:    host,
		comment: comment,
	}

	mux := http.NewServeMux()

	webhookHandler := NewWebhookHandler(&cf)
	webhookHandler.RegisterHandlers(mux)

	log.Println("starting")
	if err := http.ListenAndServe("0.0.0.0:8001", mux); err != nil {
		log.Fatal(err)
	}
}
