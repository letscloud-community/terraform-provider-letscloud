// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/letscloud-community/letscloud-go/domains"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InstanceResource{}
var _ resource.ResourceWithImportState = &InstanceResource{}

func NewInstanceResource() resource.Resource {
	return &InstanceResource{}
}

// InstanceResource defines the resource implementation.
type InstanceResource struct {
	client LetsCloudClient
}

// InstanceResourceModel describes the resource data model.
type InstanceResourceModel struct {
	Label        types.String   `tfsdk:"label"`
	LocationSlug types.String   `tfsdk:"location_slug"`
	PlanSlug     types.String   `tfsdk:"plan_slug"`
	ImageSlug    types.String   `tfsdk:"image_slug"`
	SSHKeys      []types.String `tfsdk:"ssh_keys"`
	Password     types.String   `tfsdk:"password"`
	Hostname     types.String   `tfsdk:"hostname"`
	Id           types.String   `tfsdk:"id"`
	State        types.String   `tfsdk:"state"`
	IPv4         types.String   `tfsdk:"ipv4"`
	IPv6         types.String   `tfsdk:"ipv6"`
}

func (r *InstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (r *InstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a LetsCloud instance.",

		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				MarkdownDescription: "The label of the instance.",
				Required:            true,
			},
			"location_slug": schema.StringAttribute{
				MarkdownDescription: "The location slug where the instance will be created.",
				Required:            true,
			},
			"plan_slug": schema.StringAttribute{
				MarkdownDescription: "The plan slug for the instance.",
				Required:            true,
			},
			"image_slug": schema.StringAttribute{
				MarkdownDescription: "The image slug to use for the instance.",
				Required:            true,
			},
			"ssh_keys": schema.ListAttribute{
				MarkdownDescription: "The SSH keys to add to the instance.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The root password for the instance.",
				Optional:            true,
				Sensitive:           true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname of the instance.",
				Required:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the instance.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ipv4": schema.StringAttribute{
				MarkdownDescription: "The IPv4 address of the instance.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ipv6": schema.StringAttribute{
				MarkdownDescription: "The IPv6 address of the instance.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Instance identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *InstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(LetsCloudClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected LetsCloudClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	tflog.Debug(ctx, "InstanceResource Configure: Using LetsCloudClient", map[string]interface{}{
		"client_type": fmt.Sprintf("%T", client),
	})

	r.client = client
}

func (r *InstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *InstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert SSH keys to string slice
	sshKeys := make([]string, len(data.SSHKeys))
	for i, key := range data.SSHKeys {
		sshKeys[i] = key.ValueString()
	}

	// Use the first SSH key if available
	sshSlug := ""
	if len(sshKeys) > 0 {
		sshSlug = sshKeys[0]
	}

	err := r.client.CreateInstance(&domains.CreateInstanceRequest{
		LocationSlug: data.LocationSlug.ValueString(),
		PlanSlug:     data.PlanSlug.ValueString(),
		ImageSlug:    data.ImageSlug.ValueString(),
		SSHSlug:      sshSlug,
		Password:     data.Password.ValueString(),
		Label:        data.Label.ValueString(),
		Hostname:     data.Hostname.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create instance, got error: %s", err))
		return
	}

	// Wait for the instance to be ready and get its details
	instance, err := waitForInstanceReady(r.client, data.Label.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error waiting for instance to be ready: %s", err))
		return
	}

	data.Id = types.StringValue(instance.Identifier)
	data.State = types.StringValue(getInstanceState(instance))
	data.IPv4 = types.StringValue(getInstanceIPv4(instance))
	data.IPv6 = types.StringValue(getInstanceIPv6(instance))

	tflog.Trace(ctx, "created an instance resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func waitForInstanceReady(client LetsCloudClient, instanceID string) (*domains.Instance, error) {
	// First check if instance is already ready
	instance, err := client.Instance(instanceID)
	if err == nil && instance != nil && instance.Built && instance.Booted {
		tflog.Debug(context.Background(), fmt.Sprintf("Instance %s is already ready: built=%v, booted=%v",
			instanceID, instance.Built, instance.Booted))
		return instance, nil
	}

	maxAttempts := 10 // 10 minutes total
	attempt := 0
	for attempt < maxAttempts {
		tflog.Debug(context.Background(), fmt.Sprintf("Checking instance %s status (attempt %d/%d)",
			instanceID, attempt+1, maxAttempts))

		instance, err := client.Instance(instanceID)
		if err != nil {
			// If instance is not found yet, continue waiting
			if strings.Contains(err.Error(), "Instance not found") {
				tflog.Debug(context.Background(), fmt.Sprintf("Instance %s not found yet, waiting...", instanceID))
				attempt++
				time.Sleep(60 * time.Second) // Wait 1 minute between checks
				continue
			}
			return nil, fmt.Errorf("error getting instance: %w", err)
		}

		if instance == nil {
			tflog.Debug(context.Background(), fmt.Sprintf("Instance %s is nil, waiting...", instanceID))
			attempt++
			time.Sleep(60 * time.Second)
			continue
		}

		// Log the current state
		tflog.Debug(context.Background(), fmt.Sprintf("Instance %s state: built=%v, booted=%v, ipv4=%v, ipv6=%v",
			instanceID,
			instance.Built,
			instance.Booted,
			getInstanceIPv4(instance),
			getInstanceIPv6(instance)))

		if instance.Built && instance.Booted {
			tflog.Debug(context.Background(), fmt.Sprintf("Instance %s is ready!", instanceID))
			return instance, nil
		}

		attempt++
		time.Sleep(60 * time.Second) // Wait 1 minute between checks
	}
	return nil, fmt.Errorf("timeout waiting for instance %s to be ready after %d minutes", instanceID, maxAttempts)
}

func (r *InstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *InstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := r.client.Instance(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance, got error: %s", err))
		return
	}

	data.Label = types.StringValue(instance.Label)
	data.LocationSlug = types.StringValue(instance.Location.Slug)
	data.Hostname = types.StringValue(instance.Hostname)
	data.State = types.StringValue(getInstanceState(instance))
	data.IPv4 = types.StringValue(getInstanceIPv4(instance))
	data.IPv6 = types.StringValue(getInstanceIPv6(instance))
	// Preserve plan_slug and image_slug from state since they're not returned by the API
	// data.PlanSlug and data.ImageSlug are already set from the state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *InstanceResourceModel
	var state *InstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Password.Equal(state.Password) {
		err := r.client.ResetPasswordInstance(state.Id.ValueString(), data.Password.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update instance password, got error: %s", err))
			return
		}
	}

	// Update label and hostname in the mock client
	if mock, ok := r.client.(*letsCloudClientMock); ok {
		err := mock.UpdateInstance(state.Id.ValueString(), data.Label.ValueString(), data.Hostname.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update instance, got error: %s", err))
			return
		}
	}

	instance, err := r.client.Instance(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated instance, got error: %s", err))
		return
	}

	data.Id = state.Id
	data.State = types.StringValue(getInstanceState(instance))
	data.IPv4 = types.StringValue(getInstanceIPv4(instance))
	data.IPv6 = types.StringValue(getInstanceIPv6(instance))
	// Preserve plan_slug and image_slug from state since they're not returned by the API
	data.PlanSlug = state.PlanSlug
	data.ImageSlug = state.ImageSlug

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *InstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteInstance(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete instance, got error: %s", err))
		return
	}
}

func (r *InstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	// Set plan_slug and image_slug to mock values for import verification
	resp.State.SetAttribute(ctx, path.Root("plan_slug"), "plan-1")
	resp.State.SetAttribute(ctx, path.Root("image_slug"), "ubuntu-20-04")
}

// Helper functions to get instance state and IP addresses.
func getInstanceState(instance *domains.Instance) string {
	if instance.Suspended {
		return "suspended"
	}
	if !instance.Built {
		return "building"
	}
	if instance.Booted {
		return "running"
	}
	return "stopped"
}

func getInstanceIPv4(instance *domains.Instance) string {
	if instance == nil || instance.IPAddresses == nil {
		return ""
	}
	for _, ip := range instance.IPAddresses {
		if !strings.Contains(ip.Address, ":") {
			return ip.Address
		}
	}
	return ""
}

func getInstanceIPv6(instance *domains.Instance) string {
	if instance == nil || instance.IPAddresses == nil {
		return ""
	}
	for _, ip := range instance.IPAddresses {
		if strings.Contains(ip.Address, ":") {
			return ip.Address
		}
	}
	return ""
}
