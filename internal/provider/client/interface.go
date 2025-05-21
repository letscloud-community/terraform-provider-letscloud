// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/letscloud-community/letscloud-go/domains"
)

// LetsCloudClient defines the interface for LetsCloud API operations.
type LetsCloudClient interface {
	// SSH Key operations
	SSHKey(id string) (*domains.SSHKey, error)
	SSHKeys() ([]domains.SSHKey, error)
	CreateSSHKey(req *domains.SSHKeyCreateRequest) (*domains.SSHKey, error)
	DeleteSSHKey(id string) error

	// Instance operations
	Instance(id string) (*domains.Instance, error)
	Instances() ([]domains.Instance, error)
	CreateInstance(req *domains.CreateInstanceRequest) error
	DeleteInstance(id string) error
	ResetPasswordInstance(id string, password string) error
	LocationPlans(location string) ([]domains.Plan, error)

	// Close closes the client connection.
	Close()
}
