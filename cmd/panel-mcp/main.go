package main

import (
	"context"
	"log"

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

type Image struct {
	ID   string `json:"id" jsonschema:"title=ID, description=The image ID"`
	Name string `json:"name" jsonschema:"title=Name, description=The image name"`
}

type Instance struct {
	ID        string            `json:"id" jsonschema:"title=ID, description=The instance ID"`
	Name      string            `json:"name" jsonschema:"title=Name, description=The instance Name"`
	Image     string            `json:"image" jsonschema:"title=Image, description=The instance image name"`
	Status    string            `json:"status" jsonschema:"title=Status, description=The instance status"`
	Ports     map[string]string `json:"ports" jsonschema:"title=Ports, description=The instance ports"`
	Resources InstanceResources `json:"resources" jsonschema:"title=Resources, description=The instance resources"`
	Webhooks  []string          `json:"webhooks" jsonschema:"title=Webhooks, description=The instance webhooks"`
}

type InstanceCreate struct {
	Image     string            `json:"image" jsonschema:"title=Image, description=The instance image name, required=true"`
	Resources InstanceResources `json:"resources" jsonschema:"title=Resources, description=The instance resources, required=true"`
	Webhooks  []string          `json:"webhooks" jsonschema:"title=Webhooks, description=The instance webhooks, required=true, uniqueItems=true"`
}

type InstanceResources struct {
	CPU    float64 `json:"cpu" jsonschema:"title=CPU, description=The instance resource CPU, required=true, default=1, minimum=0.1, maximum=12, multipleof=0.1"`
	Memory float64 `json:"memory" jsonschema:"title=Memory description=The instance resource memory, required=true, default=1, minimum=0.1, maximum=32, multipleof=0.1"`
	Disk   float64 `json:"disk" jsonschema:"title=Disk, description=The instance resource disk, required=true, default=0, minimum=0.1, maximum=120, multipleof=0.1"`
}

type ImageReadManyInput struct{}

type ImageReadManyOutput struct {
	Images []Image `json:"images" jsonschema:"a list of images"`
}

type InstanceCreateInput struct {
	InstanceCreate
}

type InstanceCreateOutput struct {
	Instance Instance `json:"instance" jsonschema:"the created instance"`
}

type InstanceReadManyInput struct{}

type InstanceReadManyOutput struct {
	Instances []Instance `json:"instances" jsonschema:"a list of instances"`
}

type InstanceReadInput struct {
	ID string `json:"id" jsonschema:"the instance ID"`
}

type InstanceReadOutput struct {
	Instance Instance `json:"instance" jsonschema:"the instance"`
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
	var images ImageReadManyOutput

	respImages, err := h.client.Image.Read(ctx)
	if err != nil {
		return nil, images, err
	}

	images = ImageReadManyOutput{Images: make([]Image, len(respImages))}
	for i, respImage := range respImages {
		images.Images[i] = Image{
			ID:   respImage.ID,
			Name: respImage.Name,
		}
	}

	return nil, images, nil
}

func (h *Handler) InstanceCreate(ctx context.Context, req *mcp.CallToolRequest, input InstanceCreateInput) (
	*mcp.CallToolResult,
	InstanceCreateOutput,
	error,
) {
	var instance InstanceCreateOutput

	respInstanceOptions := panel.InstanceCreate{
		Image: input.Image,
		Resources: panel.InstanceResources{
			CPU:    input.Resources.CPU,
			Memory: input.Resources.Memory,
			Disk:   input.Resources.Disk,
		},
		Webhooks: input.Webhooks,
	}
	respInstance, err := h.client.Instance.Create(ctx, respInstanceOptions)
	if err != nil {
		return nil, instance, err
	}

	instance.Instance = Instance{
		ID:     respInstance.ID,
		Name:   respInstance.Name,
		Image:  respInstance.Image,
		Status: respInstance.Status,
		Ports:  respInstance.Ports,
		Resources: InstanceResources{
			CPU:    respInstance.Resources.CPU,
			Memory: respInstance.Resources.Memory,
			Disk:   respInstance.Resources.Disk,
		},
		Webhooks: respInstance.Webhooks,
	}

	return nil, instance, nil
}

func (h *Handler) InstanceReadMany(ctx context.Context, req *mcp.CallToolRequest, input InstanceReadManyInput) (
	*mcp.CallToolResult,
	InstanceReadManyOutput,
	error,
) {
	var instances InstanceReadManyOutput

	respInstances, err := h.client.Instance.Read(ctx)
	if err != nil {
		return nil, instances, err
	}

	instances = InstanceReadManyOutput{Instances: make([]Instance, len(respInstances))}
	for i, respInstance := range respInstances {
		instances.Instances[i] = Instance{
			ID:     respInstance.ID,
			Name:   respInstance.Name,
			Image:  respInstance.Image,
			Status: respInstance.Status,
			Ports:  respInstance.Ports,
			Resources: InstanceResources{
				CPU:    respInstance.Resources.CPU,
				Memory: respInstance.Resources.Memory,
				Disk:   respInstance.Resources.Disk,
			},
			Webhooks: respInstance.Webhooks,
		}
	}

	return nil, instances, nil
}

func (h *Handler) InstanceRead(ctx context.Context, req *mcp.CallToolRequest, input InstanceReadInput) (
	*mcp.CallToolResult,
	InstanceReadOutput,
	error,
) {
	var instance InstanceReadOutput

	respInstance, err := h.client.Instance.ReadOne(ctx, input.ID)
	if err != nil {
		return nil, instance, err
	}

	instance = InstanceReadOutput{
		Instance: Instance{
			ID:     respInstance.ID,
			Name:   respInstance.Name,
			Image:  respInstance.Image,
			Status: respInstance.Status,
			Ports:  respInstance.Ports,
			Resources: InstanceResources{
				CPU:    respInstance.Resources.CPU,
				Memory: respInstance.Resources.Memory,
				Disk:   respInstance.Resources.Disk,
			},
			Webhooks: respInstance.Webhooks,
		},
	}

	return nil, instance, nil
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
	if err := h.client.Instance.Start(ctx, input.ID); err != nil {
		return nil, InstanceStartOutput{}, err
	}

	return nil, InstanceStartOutput{}, nil
}

func (h *Handler) InstanceStop(ctx context.Context, req *mcp.CallToolRequest, input InstanceStopInput) (
	*mcp.CallToolResult,
	InstanceStopOutput,
	error,
) {
	if err := h.client.Instance.Stop(ctx, input.ID); err != nil {
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
