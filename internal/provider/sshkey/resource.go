// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"context"
	"fmt"
	"strings"

	"github.com/letscloud-community/letscloud-go/domains"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/letscloud-community/terraform-provider-letscloud/internal/provider/client"
)

// SSHKeyResource is the resource implementation.
type SSHKeyResource struct {
	client client.LetsCloudClient
}

// NewSSHKeyResource is a helper function to simplify the provider implementation.
func NewSSHKeyResource() resource.Resource {
	return &SSHKeyResource{}
}

// Metadata returns the resource type name.
func (r *SSHKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

// Schema defines the schema for the resource.
func (r *SSHKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an SSH key in LetsCloud. Note: Updates to existing SSH keys are not supported by the LetsCloud API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the SSH key.",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Description: "The label for the SSH key.",
				Required:    true,
			},
			"key": schema.StringAttribute{
				Description: "The public SSH key.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *SSHKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *SSHKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SSHKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate and clean SSH key
	key := data.Key.ValueString()
	if !strings.HasPrefix(key, "ssh-rsa ") && !strings.HasPrefix(key, "ssh-ed25519 ") {
		resp.Diagnostics.AddError("Validation Error", "SSH key must start with 'ssh-rsa ' or 'ssh-ed25519 '")
		return
	}

	// Remove any comments from the key
	keyParts := strings.Split(key, " ")
	if len(keyParts) < 2 {
		resp.Diagnostics.AddError("Validation Error", "Invalid SSH key format")
		return
	}
	key = strings.Join(keyParts[:2], " ")

	createRequest := &domains.SSHKeyCreateRequest{
		Title: data.Label.ValueString(),
		Key:   key,
	}

	// Validate required fields
	if createRequest.Title == "" {
		resp.Diagnostics.AddError("Validation Error", "Label is required")
		return
	}
	if createRequest.Key == "" {
		resp.Diagnostics.AddError("Validation Error", "Key is required")
		return
	}

	// Check if label already exists
	existingKeys, listErr := r.client.SSHKeys()
	if listErr != nil {
		tflog.Error(ctx, "Error checking for existing SSH keys", map[string]interface{}{"error": listErr.Error()})
		resp.Diagnostics.AddError("Client Error", "Error checking for existing SSH keys: "+listErr.Error())
		return
	}

	for _, key := range existingKeys {
		if key.Title == createRequest.Title {
			tflog.Error(ctx, "Label already exists", map[string]interface{}{
				"label": createRequest.Title,
				"id":    key.Slug,
			})
			resp.Diagnostics.AddError("Validation Error", fmt.Sprintf("Label '%s' already exists. Please choose a different label.", createRequest.Title))
			return
		}
	}

	sshKey, err := r.client.CreateSSHKey(createRequest)
	if err != nil {
		tflog.Error(ctx, "Failed to create SSH key", map[string]interface{}{
			"error":   err.Error(),
			"request": fmt.Sprintf("%+v", createRequest),
		})
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create SSH key, got error: %s", err))
		return
	}

	data.Id = types.StringValue(sshKey.Slug)

	tflog.Info(ctx, "SSH key created successfully", map[string]interface{}{
		"id":    sshKey.Slug,
		"label": sshKey.Title,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *SSHKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SSHKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshKey, err := r.client.SSHKey(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSH key, got error: %s", err))
		return
	}

	data.Label = types.StringValue(sshKey.Title)
	// We don't update the key field as it's sensitive and not returned by the API

	tflog.Info(ctx, "SSH key read successfully", map[string]interface{}{
		"id":    sshKey.Slug,
		"label": sshKey.Title,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is not supported by the LetsCloud API.
func (r *SSHKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Updates to SSH keys are not supported by the LetsCloud API. Please delete and recreate the SSH key instead.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SSHKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SSHKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSSHKey(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete SSH key, got error: %s", err))
		return
	}

	tflog.Info(ctx, "SSH key deleted successfully", map[string]interface{}{
		"id": data.Id.ValueString(),
	})
}

// ImportState imports an existing resource into Terraform.
func (r *SSHKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
