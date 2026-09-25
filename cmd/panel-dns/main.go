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

	"github.com/ayeama/panel/pkg/api"
	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
)

const (
	// TODO rename?
	imageLabelDnsSrvName string = "com.github.ayeama.panel-dns.srv.name"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

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
	mux.HandleFunc("POST /", h.handleWebhook)
}

func (h *WebhookHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var event api.WebhookEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println("WARNING failed to deserialise event envelope")
		handleError(w, err)
		return
	}

	switch event.Type {
	case api.WebhookEventInstanceCreated:
		var eventData api.WebhookEventDataInstanceCreated
		if err := json.Unmarshal(event.Data, &eventData); err != nil {
			log.Println("WARNING failed to deserialise event data")
			handleError(w, err)
			return
		}

		comment := strings.ReplaceAll(eventData.ID, "-", "")

		name := h.cf.subdomainName(eventData.Name)
		ipaddresses, err := net.LookupIP(h.cf.host)
		if err != nil {
			log.Println("WARNING failed to lookup IP")
			handleError(w, err)
			return
		}
		content := ipaddresses[0].String() // TODO

		// TODO hardcoded
		var firstPort string
		for k := range eventData.Ports {
			firstPort = k
			break
		}
		port, err := strconv.ParseFloat(firstPort, 10)
		if err != nil {
			log.Println("WARNING failed to parse port")
			handleError(w, err)
			return
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
			log.Println("WARNING failed to create DNS A record")
			handleError(w, err)
			return
		}
		log.Println("created", dns.ARecordTypeA, name)

		if eventData.Labels[imageLabelDnsSrvName] != "" {
			srvName := "_" + eventData.Labels[imageLabelDnsSrvName] + "._tcp." + name // TODO tcp is still hardcoded
			_, err = (*h.cf.client).DNS.Records.New(ctx, dns.RecordNewParams{
				ZoneID: cloudflare.F(h.cf.zoneID),
				Body: dns.SRVRecordParam{
					Type: cloudflare.F(dns.SRVRecordTypeSRV),
					Name: cloudflare.F(srvName),
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
				log.Println("WARNING failed to create DNS SRV record")
				handleError(w, err)
				return
			}
			log.Println("created", dns.SRVRecordTypeSRV, srvName)
		}
	case api.WebhookEventInstanceDeleted:
		var eventData api.WebhookEventDataInstanceDeleted
		if err := json.Unmarshal(event.Data, &eventData); err != nil {
			log.Println("WARNING failed to deserialise event data")
			handleError(w, err)
			return
		}

		comment := strings.ReplaceAll(eventData.ID, "-", "")

		records, err := (*h.cf.client).DNS.Records.List(ctx, dns.RecordListParams{
			ZoneID:  cloudflare.F(h.cf.zoneID),
			Comment: cloudflare.F(dns.RecordListParamsComment{Exact: cloudflare.F(comment)}),
		})
		if err != nil {
			log.Println("WARNING failed to list DNS records")
			handleError(w, err)
			return
		}

		for _, record := range records.Result {
			_, err = (*h.cf.client).DNS.Records.Delete(ctx, record.ID, dns.RecordDeleteParams{
				ZoneID: cloudflare.F(h.cf.zoneID),
			})
			log.Println("deleted", record.Type, record.Name)
		}
	default:
		log.Println("WARNING unknown webhook event type")
		handleError(w, nil)
		return
	}
}

func main() {
	host := os.Getenv("PANEL_DNS_HOST")
	if host == "" {
		log.Fatal(errors.New("missing 'PANEL_DNS_HOST' environment variable"))
	}

	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	if zoneID == "" {
		log.Fatal(errors.New("missing 'CLOUDFLARE_ZONE_ID' environment variable"))
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
