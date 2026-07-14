package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ayeama/panel/internal/types"
	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
)

type Cloudflare struct {
	host   string
	zoneID string
	client *cloudflare.Client
}

func (cf *Cloudflare) subdomainName(name string) string {
	return name + "." + cf.host
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

		comment := strings.ReplaceAll(eventData.ID, "-", "")

		// TODO make async
		name := h.cf.subdomainName(eventData.Name)
		ipaddresses, err := net.LookupIP(h.cf.host)
		if err != nil {
			log.Fatal(err)
		}
		content := ipaddresses[0].String()

		port, err := strconv.ParseFloat(eventData.Ports["25565"], 10) // TODO hardcoded
		if err != nil {
			log.Fatal(err)
		}

		_, err = (*h.cf.client).DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(h.cf.zoneID),
			Body: dns.ARecordParam{
				Type:    cloudflare.F(dns.ARecordTypeA),
				Name:    cloudflare.F(name),
				Content: cloudflare.F(content),
				TTL:     cloudflare.F(dns.TTL(60)),
				Comment: cloudflare.F(comment),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("created", dns.ARecordTypeA, name)

		_, err = (*h.cf.client).DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(h.cf.zoneID),
			Body: dns.SRVRecordParam{
				Type: cloudflare.F(dns.SRVRecordTypeSRV),
				Name: cloudflare.F("_minecraft._tcp." + name),
				Data: cloudflare.F(dns.SRVRecordDataParam{
					Priority: cloudflare.F(float64(0)),
					Weight:   cloudflare.F(float64(5)),
					Port:     cloudflare.F(port),
					Target:   cloudflare.F(name),
				}),
				TTL:     cloudflare.F(dns.TTL(60)),
				Comment: cloudflare.F(comment),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("created", dns.SRVRecordTypeSRV, "_minecraft._tcp."+name)
	case types.WebhookEventInstanceDeleted:
		var eventData types.WebhookEventDataInstanceDeleted
		if err := json.Unmarshal(event.Data, &eventData); err != nil {
			log.Fatal(err)
		}

		comment := strings.ReplaceAll(eventData.ID, "-", "")

		// TODO make async
		records, err := (*h.cf.client).DNS.Records.List(ctx, dns.RecordListParams{
			ZoneID:  cloudflare.F(h.cf.zoneID),
			Comment: cloudflare.F(dns.RecordListParamsComment{Exact: cloudflare.F(comment)}),
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, record := range records.Result {
			_, err = (*h.cf.client).DNS.Records.Delete(ctx, record.ID, dns.RecordDeleteParams{
				ZoneID: cloudflare.F(h.cf.zoneID),
			})
			log.Println("deleted", record.Type, record.Name)
		}
	default:
		log.Fatal("unknown webhook event type")
	}
}

func main() {
	host := os.Getenv("PANEL_DNS_HOST")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_DNS_HOST' environment variable"))
	}

	zoneID := os.Getenv("PANEL_DNS_ZONE_ID")
	if zoneID == "" {
		log.Fatal(errors.New("missing 'PANEL_DNS_ZONE_ID' environment variable"))
	}

	client := cloudflare.NewClient()

	cf := Cloudflare{
		client: client,
		zoneID: zoneID,
		host:   host,
	}

	mux := http.NewServeMux()

	webhookHandler := NewWebhookHandler(&cf)
	webhookHandler.RegisterHandlers(mux)

	log.Println("starting")
	if err := http.ListenAndServe("0.0.0.0:8001", mux); err != nil {
		log.Fatal(err)
	}
}
