// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"terraform-provider-panel/client"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &PanelProvider{}
var _ provider.ProviderWithFunctions = &PanelProvider{}
var _ provider.ProviderWithEphemeralResources = &PanelProvider{}
var _ provider.ProviderWithActions = &PanelProvider{}

// PanelProvider defines the provider implementation.
type PanelProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PanelProviderModel describes the provider data model.
type PanelProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Verify   types.Bool   `tfsdk:"verify"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *PanelProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "panel"
	resp.Version = p.version
}

func (p *PanelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "API endpoint",
				Optional:            true,
			},
			"verify": schema.BoolAttribute{
				MarkdownDescription: "API verify self-signed certificates",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "API basic authentication username",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "API basic authentication password",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *PanelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PanelProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("PANEL_ENDPOINT")

	verify := true
	if value := os.Getenv("PANEL_VERIFY"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid PANEL_VERIFY value",
				fmt.Sprintf("Unable to parse PANEL_VERIFY=%q as a boolean: %s", value, err),
			)
			return
		}
		verify = parsed
	}

	username := os.Getenv("PANEL_USERNAME")
	password := os.Getenv("PANEL_PASSWORD")

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	if !data.Verify.IsNull() {
		verify = data.Verify.ValueBool()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	client := &client.Client{
		Endpoint: endpoint,
		Username: username,
		Password: password,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: !verify,
				},
			},
		},
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PanelProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServerResource,
	}
}

func (p *PanelProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *PanelProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewImageDataSource,
		NewImagesDataSource,
		NewServerDataSource,
	}
}

func (p *PanelProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *PanelProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PanelProvider{
			version: version,
		}
	}
}
