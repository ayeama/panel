package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	panel "github.com/ayeama/panel/pkg/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Handler struct {
	client *panel.Client
}

func newHandler() *Handler {
	url := "http://localhost:8000"
	client := panel.NewClient(url)

	handler := &Handler{
		client: client,
	}

	return handler
}

type ImageReadManyInput struct{}

type ImageReadManyOutput struct {
	Images []panel.Image `json:"images" jsonschema:"a list of images"`
}

type InstanceCreateInput struct {
	// TODO types
	Image    string   `json:"image" jsonschema:"the image ID to to create the instance with"`
	Webhooks []string `json:"webhooks" jsonschema:"the webhooks to create the instance with"`

	panel.InstanceResources
}

type InstanceCreateOutput struct {
	InstanceID string `json:"instance_id" jsonschema:"the created instance's ID"`
}

type InstanceReadManyInput struct{}

type InstanceReadManyOutput struct {
	Instances []panel.Instance `json:"instances" jsonschema:"a list of instances"`
}

type InstanceReadInput struct {
	ID string `json:"id" jsonschema:"the instance ID"`
}

type InstanceReadOutput struct {
	Instance panel.Instance `json:"instance" jsonschema:"the instance"`
}

func (h *Handler) ImageReadMany(ctx context.Context, req *mcp.CallToolRequest, input ImageReadManyInput) (
	*mcp.CallToolResult,
	ImageReadManyOutput,
	error,
) {
	images, err := h.client.Image.Read(ctx)
	if err != nil {
		return nil, ImageReadManyOutput{}, err
	}

	return nil, ImageReadManyOutput{Images: images}, nil
}

func (h *Handler) InstanceCreate(ctx context.Context, req *mcp.CallToolRequest, input InstanceCreateInput) (
	*mcp.CallToolResult,
	InstanceCreateOutput,
	error,
) {
	var instance InstanceCreateOutput

	body, err := json.Marshal(input)
	if err != nil {
		return nil, instance, err
	}

	url := "http://localhost:8000/instances"
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, instance, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, instance, err
	}

	return nil, instance, nil
}

func (h *Handler) InstanceReadMany(ctx context.Context, req *mcp.CallToolRequest, input InstanceReadManyInput) (
	*mcp.CallToolResult,
	InstanceReadManyOutput,
	error,
) {
	instances, err := h.client.Instance.Read(ctx)
	if err != nil {
		return nil, InstanceReadManyOutput{}, err
	}

	return nil, InstanceReadManyOutput{Instances: instances}, nil
}

func (h *Handler) InstanceRead(ctx context.Context, req *mcp.CallToolRequest, input InstanceReadInput) (
	*mcp.CallToolResult,
	InstanceReadOutput,
	error,
) {
	instance, err := h.client.Instance.ReadOne(ctx, input.ID)
	if err != nil {
		return nil, InstanceReadOutput{}, err
	}

	return nil, InstanceReadOutput{Instance: instance}, nil
}

func annotationRead() *mcp.ToolAnnotations {
	destructive := false
	idempotent := false
	openworld := true
	readonly := true
	return &mcp.ToolAnnotations{DestructiveHint: &destructive, IdempotentHint: idempotent, OpenWorldHint: &openworld, ReadOnlyHint: readonly}
}

func annotationWrite() *mcp.ToolAnnotations {
	destructive := true
	idempotent := false
	openworld := true
	readonly := false
	return &mcp.ToolAnnotations{DestructiveHint: &destructive, IdempotentHint: idempotent, OpenWorldHint: &openworld, ReadOnlyHint: readonly}
}

func main() {
	handler := newHandler()

	server := mcp.NewServer(&mcp.Implementation{Name: "panel", Title: "Panel", Version: "v0.0.1"}, nil)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "image_read_many",
			Description: "read many images",
			Annotations: annotationRead(),
		},
		handler.ImageReadMany,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "instance_create",
			Description: "create instance",
			Annotations: annotationWrite(),
		},
		handler.InstanceCreate,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "instance_read_many",
			Description: "read many instances",
			Annotations: annotationRead(),
		},
		handler.InstanceReadMany,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "instance_read",
			Description: "read one instance",
			Annotations: annotationRead(),
		},
		handler.InstanceRead,
	)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
