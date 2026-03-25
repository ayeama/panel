// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"panel/pkg/api"
	"terraform-provider-panel/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &ImageDataSource{}

func NewImageDataSource() datasource.DataSource {
	return &ImageDataSource{}
}

type ImageDataSource struct {
	client *client.Client
}

type ImageDataSourceModel struct {
	Filter     types.Set    `tfsdk:"filter"`
	Sort       types.String `tfsdk:"sort"`
	MostRecent types.Bool   `tfsdk:"most_recent"`

	// ID         types.String `tfsdk:"id"`
	Repository types.String `tfsdk:"repository"`
	Tag        types.String `tfsdk:"tag"`
	// Env        types.Map    `tfsdk:"env"`
}

type ImageDataSourceFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

func (d *ImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (d *ImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Image data source",

		Attributes: map[string]schema.Attribute{
			"sort": schema.StringAttribute{
				Optional: true,
			},
			"most_recent": schema.BoolAttribute{
				Optional: true,
			},
			// "id": schema.StringAttribute{
			// 	Computed: true,
			// },
			"repository": schema.StringAttribute{
				Computed: true,
			},
			"tag": schema.StringAttribute{
				Computed: true,
			},
			// "env": schema.MapAttribute{
			// 	ElementType: types.StringType,
			// 	Computed:    true,
			// },
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

func (d *ImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ImageDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := "/images"

	values := url.Values{}
	if !data.Filter.IsNull() && !data.Filter.IsUnknown() {
		var filters []ImageDataSourceFilterModel
		resp.Diagnostics.Append(data.Filter.ElementsAs(ctx, &filters, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, filter := range filters {
			name := filter.Name.ValueString()
			if name == "" {
				continue
			}

			var filterValues []string
			resp.Diagnostics.Append(filter.Values.ElementsAs(ctx, &filterValues, false)...)
			if resp.Diagnostics.HasError() {
				return
			}

			for _, v := range filterValues {
				if v == "" {
					continue
				}
				values.Add(name, v)
			}
		}
	}

	if len(values) > 0 {
		path += "?" + values.Encode()
	}

	httpReq, err := d.client.NewRequest(ctx, http.MethodGet, path, nil)
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

	if len(imagesResp.Items) == 0 {
		resp.Diagnostics.AddError(
			"No matching image",
			"No images matched the supplied filters.",
		)
		return
	}

	if len(imagesResp.Items) > 1 {
		resp.Diagnostics.AddError(
			"Multiple matching images",
			fmt.Sprintf("Expected exactly 1 image, found %d. Add more specific filters.", len(imagesResp.Items)),
		)
		return
	}

	image := imagesResp.Items[0]

	// envValue, diags := types.MapValueFrom(ctx, types.StringType, image.Env)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// data.ID = types.StringValue(image.Id)
	data.Repository = types.StringValue(image.Repository)
	data.Tag = types.StringValue(image.Tag)
	// data.Env = envValue

	tflog.Trace(ctx, "read images data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
