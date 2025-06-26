// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/letscloud-community/terraform-provider-letscloud/internal/provider/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SSHKeysDataSource{}

func NewSSHKeysDataSource() datasource.DataSource {
	return &SSHKeysDataSource{}
}

// SSHKeysDataSource defines the data source implementation.
type SSHKeysDataSource struct {
	client client.LetsCloudClient
}

// SSHKeysDataSourceModel describes the data source data model for multiple SSH keys.
type SSHKeysDataSourceModel struct {
	SSHKeys []SSHKeyDataSourceModel `tfsdk:"ssh_keys"`
}

func (d *SSHKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_keys"
}

func (d *SSHKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about all SSH keys in the account.",

		Attributes: map[string]schema.Attribute{
			"ssh_keys": schema.ListNestedAttribute{
				MarkdownDescription: "List of SSH keys in the account.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier for the SSH key.",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "The label of the SSH key.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SSHKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SSHKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSHKeysDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch all SSH keys
	sshKeys, err := d.client.SSHKeys()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list SSH keys, got error: %s", err))
		return
	}

	// Map response to model
	data.SSHKeys = make([]SSHKeyDataSourceModel, len(sshKeys))
	for i, key := range sshKeys {
		data.SSHKeys[i] = SSHKeyDataSourceModel{
			Id:    types.StringValue(key.Slug),
			Label: types.StringValue(key.Title),
		}
	}

	tflog.Info(ctx, "SSH keys data source read successfully", map[string]interface{}{
		"count": len(sshKeys),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
