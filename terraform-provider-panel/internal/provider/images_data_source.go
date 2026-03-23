// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"panel/pkg/api"
	"terraform-provider-panel/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &ImagesDataSource{}

func NewImagesDataSource() datasource.DataSource {
	return &ImagesDataSource{}
}

type ImagesDataSource struct {
	client *client.Client
}

type ImagesDataSourceModel struct {
	Filter types.Set  `tfsdk:"filter"`
	Images types.List `tfsdk:"images"`
}

type ImagesDataSourceFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

type ImagesDataSourceImageModel struct {
	ID         types.String `tfsdk:"id"`
	Repository types.String `tfsdk:"repository"`
	Tag        types.String `tfsdk:"tag"`
	Env        types.Map    `tfsdk:"env"`
}

func (d *ImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

func (d *ImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Images data source",

		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"repository": schema.StringAttribute{
							Computed: true,
						},
						"tag": schema.StringAttribute{
							Computed: true,
						},
						"env": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},

		Blocks: map[string]schema.Block{
			"filter": schema.SetNestedBlock{
				Description: "Filter images",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ImagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ImagesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpReq, err := d.client.NewRequest(ctx, http.MethodGet, "/images", nil)
	if err != nil {
		resp.Diagnostics.AddError("Request error", err.Error())
		return
	}

	httpResp, err := d.client.Client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client error", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"API error",
			fmt.Sprintf("Unexpected status code when reading images: %d", httpResp.StatusCode))
		return
	}

	var imagesResp api.ImageListResponse
	decode := json.NewDecoder(httpResp.Body)
	err = decode.Decode(&imagesResp)
	if err != nil {
		resp.Diagnostics.AddError("Parsing error", err.Error())
		return
	}

	var imageModels []ImagesDataSourceImageModel
	for _, image := range imagesResp.Items {
		envValue, diags := types.MapValueFrom(ctx, types.StringType, image.Env)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		imageModels = append(imageModels, ImagesDataSourceImageModel{
			ID:         types.StringValue(image.Id),
			Repository: types.StringValue(image.Repository),
			Tag:        types.StringValue(image.Tag),
			Env:        envValue,
		})
	}

	data.Images, resp.Diagnostics = types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":         types.StringType,
			"repository": types.StringType,
			"tag":        types.StringType,
			"env":        types.MapType{ElemType: types.StringType},
		},
	}, imageModels)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read images data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
