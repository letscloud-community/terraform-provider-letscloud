// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"context"
	"fmt"

	"github.com/letscloud-community/letscloud-go/domains"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/letscloud-community/terraform-provider-letscloud/internal/provider/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SSHKeyDataSource{}

func NewSSHKeyDataSource() datasource.DataSource {
	return &SSHKeyDataSource{}
}

// SSHKeyDataSource defines the data source implementation.
type SSHKeyDataSource struct {
	client client.LetsCloudClient
}

// SSHKeyDataSourceModel describes the data source data model.
type SSHKeyDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
}

func (d *SSHKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key_lookup"
}

func (d *SSHKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about an SSH key by its ID or label.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the SSH key. Either 'id' or 'label' must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "The label of the SSH key. Either 'id' or 'label' must be specified.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (d *SSHKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.LetsCloudClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected client.LetsCloudClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *SSHKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSHKeyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either id or label is provided
	if data.Id.IsNull() && data.Label.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'label' must be specified to identify the SSH key.",
		)
		return
	}

	// If both are provided, prefer id
	if !data.Id.IsNull() && !data.Label.IsNull() {
		tflog.Warn(ctx, "Both 'id' and 'label' provided, using 'id' to identify SSH key")
	}

	var sshKey *domains.SSHKey
	var err error

	if !data.Id.IsNull() {
		// Fetch by ID
		sshKey, err = d.client.SSHKey(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSH key by ID, got error: %s", err))
			return
		}
	} else {
		// Fetch by label - need to list all and find by label
		sshKeys, err := d.client.SSHKeys()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list SSH keys, got error: %s", err))
			return
		}

		labelToFind := data.Label.ValueString()
		for _, key := range sshKeys {
			if key.Title == labelToFind {
				sshKey = &key
				break
			}
		}

		if sshKey == nil {
			resp.Diagnostics.AddError(
				"SSH Key Not Found",
				fmt.Sprintf("No SSH key found with label '%s'", labelToFind),
			)
			return
		}
	}

	// Map response body to model
	data.Id = types.StringValue(sshKey.Slug)
	data.Label = types.StringValue(sshKey.Title)

	tflog.Info(ctx, "SSH key data source read successfully", map[string]interface{}{
		"id":    sshKey.Slug,
		"label": sshKey.Title,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
