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
	Instance panel.Instance `json:"instance" jsonschema:"the created instance"`
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

type InstanceDeleteInput struct {
	ID string `json:"id" jsonschema:"the instance ID"`
}

type InstanceDeleteOutput struct{}

type InstanceStartInput struct {
	ID string `json:"id" jsonschema:"the instance ID"`
}

type InstanceStartOutput struct{}

type InstanceStopInput struct {
	ID string `json:"id" jsonschema:"the instance ID"`
}

type InstanceStopOutput struct{}

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

func (h *Handler) InstanceDelete(ctx context.Context, req *mcp.CallToolRequest, input InstanceDeleteInput) (
	*mcp.CallToolResult,
	InstanceDeleteOutput,
	error,
) {
	if err := h.client.Instance.Delete(ctx, input.ID); err != nil {

		return nil, InstanceDeleteOutput{}, err
	}

	return nil, InstanceDeleteOutput{}, nil
}

func (h *Handler) InstanceStart(ctx context.Context, req *mcp.CallToolRequest, input InstanceStartInput) (
	*mcp.CallToolResult,
	InstanceStartOutput,
	error,
) {
	err := h.client.Instance.Start(ctx, input.ID)
	if err != nil {
		return nil, InstanceStartOutput{}, err
	}

	return nil, InstanceStartOutput{}, nil
}

func (h *Handler) InstanceStop(ctx context.Context, req *mcp.CallToolRequest, input InstanceStopInput) (
	*mcp.CallToolResult,
	InstanceStopOutput,
	error,
) {
	err := h.client.Instance.Stop(ctx, input.ID)
	if err != nil {
		return nil, InstanceStopOutput{}, err
	}

	return nil, InstanceStopOutput{}, nil
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

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "instance_delete",
			Description: "delete instance",
			Annotations: annotationWrite(),
		},
		handler.InstanceDelete,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "instance_start",
			Description: "start instance",
			Annotations: annotationWrite(),
		},
		handler.InstanceStart,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "instance_stop",
			Description: "stop instance",
			Annotations: annotationWrite(),
		},
		handler.InstanceStop,
	)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
