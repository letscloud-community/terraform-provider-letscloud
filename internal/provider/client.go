// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/letscloud-community/letscloud-go/domains"

// LetsCloudClient defines the interface for LetsCloud client operations.
type LetsCloudClient interface {
	CreateInstance(req *domains.CreateInstanceRequest) error
	Instance(id string) (*domains.Instance, error)
	DeleteInstance(id string) error
	ResetPasswordInstance(id string, password string) error
	LocationPlans(location string) ([]domains.Plan, error)
}
