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
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(LetsCloudClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected client.LetsCloudClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *InstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *InstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting instance creation", map[string]interface{}{
		"label":         data.Label.ValueString(),
		"location_slug": data.LocationSlug.ValueString(),
		"plan_slug":     data.PlanSlug.ValueString(),
		"image_slug":    data.ImageSlug.ValueString(),
		"hostname":      data.Hostname.ValueString(),
	})

	// Convert SSH keys to string slice
	sshKeys := make([]string, len(data.SSHKeys))
	for i, key := range data.SSHKeys {
		sshKeys[i] = key.ValueString()
	}

	// Use the first SSH key if available
	sshSlug := ""
	if len(sshKeys) > 0 {
		sshSlug = sshKeys[0]
		tflog.Debug(ctx, "Using SSH key", map[string]interface{}{
			"ssh_key_prefix": sshSlug[:10] + "...",
		})
	}

	createRequest := &domains.CreateInstanceRequest{
		LocationSlug: data.LocationSlug.ValueString(),
		PlanSlug:     data.PlanSlug.ValueString(),
		ImageSlug:    data.ImageSlug.ValueString(),
		SSHSlug:      sshSlug,
		Password:     data.Password.ValueString(),
		Label:        data.Label.ValueString(),
		Hostname:     data.Hostname.ValueString(),
	}

	tflog.Info(ctx, "Preparing instance creation request", map[string]interface{}{
		"label":         createRequest.Label,
		"location_slug": createRequest.LocationSlug,
		"plan_slug":     createRequest.PlanSlug,
		"image_slug":    createRequest.ImageSlug,
		"hostname":      createRequest.Hostname,
		"has_ssh_key":   createRequest.SSHSlug != "",
		"has_password":  createRequest.Password != "",
	})

	// Validate required fields
	if createRequest.Label == "" {
		resp.Diagnostics.AddError("Validation Error", "Label is required")
		return
	}
	if createRequest.LocationSlug == "" {
		resp.Diagnostics.AddError("Validation Error", "Location slug is required")
		return
	}
	if createRequest.PlanSlug == "" {
		resp.Diagnostics.AddError("Validation Error", "Plan slug is required")
		return
	}
	if createRequest.ImageSlug == "" {
		resp.Diagnostics.AddError("Validation Error", "Image slug is required")
		return
	}
	if createRequest.Hostname == "" {
		resp.Diagnostics.AddError("Validation Error", "Hostname is required")
		return
	}

	// Check if label already exists
	existingInstances, listErr := r.client.Instances()
	if listErr != nil {
		tflog.Error(ctx, "Error checking for existing instances", map[string]interface{}{"error": listErr.Error()})
		resp.Diagnostics.AddError("Client Error", "Error checking for existing instances: "+listErr.Error())
		return
	}

	for _, inst := range existingInstances {
		if inst.Label == createRequest.Label {
			tflog.Error(ctx, "Label already exists", map[string]interface{}{
				"label": createRequest.Label,
				"id":    inst.Identifier,
			})
			resp.Diagnostics.AddError("Validation Error", fmt.Sprintf("Label '%s' already exists. Please choose a different label.", createRequest.Label))
			return
		}
	}

	// Create the instance
	tflog.Debug(ctx, "Sending create instance request to LetsCloud API", map[string]interface{}{
		"request": fmt.Sprintf("%+v", createRequest),
	})

	// Try to create the instance with retries
	maxRetries := 3
	retryCount := 0
	var createErr error

	for retryCount < maxRetries {
		tflog.Info(ctx, "Attempting to create instance", map[string]interface{}{
			"attempt":     retryCount + 1,
			"max_retries": maxRetries,
			"label":       createRequest.Label,
			"location":    createRequest.LocationSlug,
			"plan":        createRequest.PlanSlug,
			"image":       createRequest.ImageSlug,
			"hostname":    createRequest.Hostname,
		})

		createErr = r.client.CreateInstance(createRequest)
		if createErr == nil {
			tflog.Info(ctx, "Instance creation request sent successfully", map[string]interface{}{
				"label": createRequest.Label,
			})
			break
		}

		tflog.Warn(ctx, "Failed to create instance, retrying...", map[string]interface{}{
			"error":       createErr.Error(),
			"retry_count": retryCount + 1,
			"max_retries": maxRetries,
			"label":       createRequest.Label,
		})

		retryCount++
		if retryCount < maxRetries {
			time.Sleep(5 * time.Second) // Wait 5 seconds between retries
		}
	}

	if createErr != nil {
		tflog.Error(ctx, "Failed to create instance after retries", map[string]interface{}{
			"error":   createErr.Error(),
			"request": fmt.Sprintf("%+v", createRequest),
		})
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create instance, got error: %s", createErr))
		return
	}

	// Now wait for the instance to be ready using the label and hostname
	instance, err := waitForInstanceReady(ctx, r.client, createRequest.Label, createRequest.Hostname)
	if err != nil {
		tflog.Error(ctx, "Error waiting for instance to be ready", map[string]interface{}{
			"error":    err.Error(),
			"label":    createRequest.Label,
			"hostname": createRequest.Hostname,
		})
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error waiting for instance to be ready: %s", err))
		return
	}

	if instance == nil {
		tflog.Error(ctx, "Instance is nil after waiting for ready state", map[string]interface{}{
			"label": createRequest.Label,
		})
		resp.Diagnostics.AddError("Client Error", "Instance is nil after waiting for ready state")
		return
	}

	data.Id = types.StringValue(instance.Identifier)
	data.State = types.StringValue(getInstanceState(instance))
	data.IPv4 = types.StringValue(getInstanceIPv4(instance))
	data.IPv6 = types.StringValue(getInstanceIPv6(instance))

	tflog.Info(ctx, "Instance created successfully", map[string]interface{}{
		"id":           instance.Identifier,
		"state":        getInstanceState(instance),
		"ipv4":         getInstanceIPv4(instance),
		"ipv6":         getInstanceIPv6(instance),
		"label":        instance.Label,
		"built":        instance.Built,
		"booted":       instance.Booted,
		"raw_instance": fmt.Sprintf("%+v", instance),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func waitForInstanceReady(ctx context.Context, client LetsCloudClient, label, hostname string) (*domains.Instance, error) {
	maxAttempts := 400 // 20 minutes total (400 attempts × 3 seconds)
	attempt := 0
	lastError := ""
	lastState := ""
	lastIPs := ""
	lastResponse := ""
	var instanceID string

	for attempt < maxAttempts {
		tflog.Info(ctx, "Checking instance status", map[string]interface{}{
			"label":           label,
			"hostname":        hostname,
			"attempt":         attempt + 1,
			"max_attempts":    maxAttempts,
			"seconds_elapsed": attempt * 3,
			"last_error":      lastError,
			"last_state":      lastState,
			"last_ips":        lastIPs,
			"last_response":   lastResponse,
		})

		// List all instances to find our target
		instances, err := client.Instances()
		if err != nil {
			lastError = err.Error()
			tflog.Warn(ctx, "Error listing instances", map[string]interface{}{
				"error":           err.Error(),
				"attempt":         attempt + 1,
				"max_attempts":    maxAttempts,
				"seconds_elapsed": attempt * 3,
			})
			attempt++
			time.Sleep(3 * time.Second)
			continue
		}

		// Find our instance by label and hostname
		var instance *domains.Instance
		for _, inst := range instances {
			if inst.Label == label && inst.Hostname == hostname {
				instanceCopy := inst
				instance = &instanceCopy
				instanceID = instance.Identifier
				break
			}
		}

		if instance == nil {
			tflog.Info(ctx, "Instance not found yet, waiting...", map[string]interface{}{
				"label":           label,
				"hostname":        hostname,
				"attempt":         attempt + 1,
				"max_attempts":    maxAttempts,
				"seconds_elapsed": attempt * 3,
			})
			attempt++
			time.Sleep(3 * time.Second)
			continue
		}

		// Log the current state
		currentState := getInstanceState(instance)
		lastState = currentState
		currentIPs := fmt.Sprintf("IPv4: %s, IPv6: %s", getInstanceIPv4(instance), getInstanceIPv6(instance))
		lastIPs = currentIPs
		lastResponse = fmt.Sprintf("%+v", instance)

		tflog.Info(ctx, "Current instance state", map[string]interface{}{
			"instance_id":     instanceID,
			"built":           instance.Built,
			"booted":          instance.Booted,
			"state":           currentState,
			"ipv4":            getInstanceIPv4(instance),
			"ipv6":            getInstanceIPv6(instance),
			"attempt":         attempt + 1,
			"max_attempts":    maxAttempts,
			"seconds_elapsed": attempt * 3,
			"raw_instance":    lastResponse,
		})

		// Check if instance is in an error state
		if instance.Suspended {
			return nil, fmt.Errorf("instance %s is suspended", instanceID)
		}

		// Check if instance has IP addresses assigned
		hasIPv4 := getInstanceIPv4(instance) != ""
		hasIPv6 := getInstanceIPv6(instance) != ""

		// Instance is considered ready when:
		// 1. It is built and booted
		// 2. It has at least one IP address assigned
		if instance.Built && instance.Booted && (hasIPv4 || hasIPv6) {
			tflog.Info(ctx, "Instance is ready!", map[string]interface{}{
				"instance_id":  instanceID,
				"state":        currentState,
				"ipv4":         getInstanceIPv4(instance),
				"ipv6":         getInstanceIPv6(instance),
				"has_ipv4":     hasIPv4,
				"has_ipv6":     hasIPv6,
				"raw_instance": lastResponse,
			})
			return instance, nil
		}

		// Log progress information
		tflog.Info(ctx, "Instance still not ready", map[string]interface{}{
			"instance_id":     instanceID,
			"built":           instance.Built,
			"booted":          instance.Booted,
			"has_ipv4":        hasIPv4,
			"has_ipv6":        hasIPv6,
			"state":           currentState,
			"attempt":         attempt + 1,
			"max_attempts":    maxAttempts,
			"seconds_elapsed": attempt * 3,
			"raw_instance":    lastResponse,
		})

		attempt++
		time.Sleep(3 * time.Second) // Wait 3 seconds between checks
	}

	return nil, fmt.Errorf("timeout waiting for instance with label %s and hostname %s to be ready after %d seconds. Last known state: %s, Last error: %s, Last IPs: %s, Last response: %s",
		label, hostname, maxAttempts*3, lastState, lastError, lastIPs, lastResponse)
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
		instance, err := mock.Instance(state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance for update, got error: %s", err))
			return
		}
		instance.Label = data.Label.ValueString()
		instance.Hostname = data.Hostname.ValueString()
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
