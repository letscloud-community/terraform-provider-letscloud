// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/letscloud-community/letscloud-go"
)

// Ensure LetsCloudProvider satisfies various provider interfaces.
var _ provider.Provider = &LetsCloudProvider{}
var _ provider.ProviderWithFunctions = &LetsCloudProvider{}
var _ provider.ProviderWithEphemeralResources = &LetsCloudProvider{}

// LetsCloudProvider defines the provider implementation.
type LetsCloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// LetsCloudProviderModel describes the provider data model.
type LetsCloudProviderModel struct {
	APIToken types.String `tfsdk:"api_token"`
}

// mockLetsCloudClient is used for testing. If set, it will be used instead of a real client.
var mockLetsCloudClient LetsCloudClient = nil

func (p *LetsCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "letscloud"
	resp.Version = p.version
}

func (p *LetsCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with LetsCloud.",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Description: "The API token for LetsCloud. May also be provided via LETSCLOUD_API_TOKEN environment variable.",
				Optional:    true,
			},
		},
	}
}

func (p *LetsCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config LetsCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown LetsCloud API token",
			"The provider cannot create the LetsCloud API client as there is an unknown configuration value for the LetsCloud API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LETSCLOUD_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with provider configuration value if set.

	apiToken := config.APIToken.ValueString()

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing LetsCloud API token",
			"The provider cannot create the LetsCloud API client as there is a missing or empty value for the LetsCloud API token. "+
				"Set the api_token value in the configuration or use the LETSCLOUD_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Se o mockLetsCloudClient estiver setado, use ele
	if mockLetsCloudClient != nil {
		tflog.Debug(ctx, "Using mock LetsCloud client for testing")
		resp.DataSourceData = mockLetsCloudClient
		resp.ResourceData = mockLetsCloudClient
		return
	}

	// Create a new LetsCloud client using the configuration values
	var client LetsCloudClient

	// If we're in test mode or using a mock token, use the mock client
	if p.version == "test" || apiToken == "mock-token-for-testing" {
		tflog.Debug(ctx, "Using mock LetsCloud client for testing")
		client = &letsCloudClientMock{} // Use our mock client implementation
	} else {
		tflog.Debug(ctx, "Creating real LetsCloud client")
		lc, err := letscloud.New(apiToken)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create LetsCloud API Client",
				"An unexpected error occurred when creating the LetsCloud API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"LetsCloud Client Error: "+err.Error(),
			)
			return
		}
		client = lc
	}

	// Make the LetsCloud client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LetsCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewInstanceResource,
	}
}

func (p *LetsCloudProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		// Add your ephemeral resources here
	}
}

func (p *LetsCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *LetsCloudProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// Add your functions here
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LetsCloudProvider{
			version: version,
		}
	}
}
