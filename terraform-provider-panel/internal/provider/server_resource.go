// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"panel/pkg/api"
	"terraform-provider-panel/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ServerResource{}
var _ resource.ResourceWithImportState = &ServerResource{}

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

type ServerResource struct {
	client *client.Client
}

type ServerResourceModel struct {
	// ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Image  types.String `tfsdk:"image"`
	Status types.String `tfsdk:"status"`
	Ports  types.List   `tfsdk:"ports"`
	// Env    types.MapType  `tfsdk:"env"`
}

func (r *ServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Server resource",

		Attributes: map[string]schema.Attribute{
			// "id": schema.StringAttribute{
			// 	MarkdownDescription: "Server identifier",
			// 	Computed:            true,
			// 	PlanModifiers: []planmodifier.String{
			// 		stringplanmodifier.UseStateForUnknown(),
			// 	},
			// },
			"name": schema.StringAttribute{
				MarkdownDescription: "Server name",
				Computed:            true,
			},
			"image": schema.StringAttribute{
				MarkdownDescription: "Server image",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Server status",
				Computed:            true,
			},
			"ports": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			// "env": schema.MapAttribute{
			// 	ElementType: types.StringType,
			// 	Optional:    true,
			// },
		},
	}
}

func (r *ServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestServer := api.ServerCreateRequest{
		Image: data.Image.ValueString(),
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(requestServer); err != nil {
		resp.Diagnostics.AddError("Encode error", err.Error())
		return
	}

	httpReq, err := r.client.NewRequest(ctx, http.MethodPost, "/servers", buf)
	if err != nil {
		resp.Diagnostics.AddError("Client error", err.Error())
		return
	}

	httpResp, err := r.client.Client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client error", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError(
			"API error",
			fmt.Sprintf("Unexpected status code when creating server: %d", httpResp.StatusCode))
		return
	}

	var responseServer api.ServerResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&responseServer); err != nil {
		resp.Diagnostics.AddError("Parsing error", err.Error())
		return
	}

	// data.ID = types.StringValue(responseServer.Id)
	data.Name = types.StringValue(responseServer.Name)
	data.Image = types.StringValue(responseServer.Image)
	data.Status = types.StringValue(responseServer.Status)

	ports, diags := types.ListValueFrom(ctx, types.StringType, responseServer.Ports)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Ports = ports

	tflog.Trace(ctx, "created a server resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()

	httpReq, err := r.client.NewRequest(ctx, http.MethodGet, "/servers/"+url.PathEscape(name), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client error", err.Error())
		return
	}

	httpResp, err := r.client.Client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client error", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"API error",
			fmt.Sprintf("Unexpected status code when reading server: %d", httpResp.StatusCode),
		)
		return
	}

	var responseServer api.ServerResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&responseServer); err != nil {
		resp.Diagnostics.AddError("Parsing error", err.Error())
		return
	}

	// data.ID = types.StringValue(responseServer.Id)
	data.Name = types.StringValue(responseServer.Name)
	data.Image = types.StringValue(responseServer.Image)
	data.Status = types.StringValue(responseServer.Status)

	ports, diags := types.ListValueFrom(ctx, types.StringType, responseServer.Ports)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Ports = ports

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()

	httpReq, err := r.client.NewRequest(ctx, http.MethodDelete, "/servers/"+name, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client error", err.Error())
		return
	}

	httpResp, err := r.client.Client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client error", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return
	}

	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"API error",
			fmt.Sprintf("Unexpected status code when deleting server: %d", httpResp.StatusCode),
		)
		return
	}

	tflog.Trace(ctx, "deleted server resource", map[string]any{"name": name})
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
