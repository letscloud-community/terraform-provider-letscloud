// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/letscloud-community/letscloud-go"
	"github.com/letscloud-community/letscloud-go/domains"
)

// RealLetsCloudClient wraps the SDK client to implement our interface.
type RealLetsCloudClient struct {
	client *letscloud.LetsCloud
}

// NewRealLetsCloudClient creates a new real client wrapper.
func NewRealLetsCloudClient(lc *letscloud.LetsCloud) *RealLetsCloudClient {
	return &RealLetsCloudClient{client: lc}
}

// Close closes the gRPC connection.
func (c *RealLetsCloudClient) Close() {
	if c.client != nil {
		// The SDK should handle closing the connection
		c.client = nil
	}
}

// SSH Key methods.
func (c *RealLetsCloudClient) SSHKey(id string) (*domains.SSHKey, error) {
	return c.client.SSHKey(id)
}

func (c *RealLetsCloudClient) SSHKeys() ([]domains.SSHKey, error) {
	return c.client.SSHKeys()
}

func (c *RealLetsCloudClient) CreateSSHKey(req *domains.SSHKeyCreateRequest) (*domains.SSHKey, error) {
	return c.client.NewSSHKey(req.Title, req.Key)
}

func (c *RealLetsCloudClient) DeleteSSHKey(id string) error {
	return c.client.DeleteSSHKey(id)
}

// Instance methods.
func (c *RealLetsCloudClient) Instance(id string) (*domains.Instance, error) {
	return c.client.Instance(id)
}

func (c *RealLetsCloudClient) Instances() ([]domains.Instance, error) {
	return c.client.Instances()
}

func (c *RealLetsCloudClient) CreateInstance(req *domains.CreateInstanceRequest) error {
	return c.client.CreateInstance(req)
}

func (c *RealLetsCloudClient) DeleteInstance(id string) error {
	return c.client.DeleteInstance(id)
}

func (c *RealLetsCloudClient) ResetPasswordInstance(id string, password string) error {
	return c.client.ResetPasswordInstance(id, password)
}

func (c *RealLetsCloudClient) LocationPlans(location string) ([]domains.Plan, error) {
	return c.client.LocationPlans(location)
}
