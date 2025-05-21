// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sshkey_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/letscloud-community/terraform-provider-letscloud/internal/provider"
)

const (
	providerConfig = `
provider "letscloud" {
  api_token = "mock-token-for-testing"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"letscloud": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
)

func TestAccSSHKeyResource(t *testing.T) {
	// Configura o mock client
	mockClient := provider.NewLetsCloudClientMock()
	provider.MockLetsCloudClient = mockClient

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "letscloud_ssh_key" "test" {
  label = "test-key"
  key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIi41q4x71KvIciZZq1cU8DADrcqbn3ySVhgfTrZG1yUloU3Jq7Nn8/7YzCMr0CFSf/7ZgRGW9P1QxJZz8K3mWJ2z2K3Uf1x1+9Z9L3Qw== test@example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("letscloud_ssh_key.test", "label", "test-key"),
					resource.TestCheckResourceAttrSet("letscloud_ssh_key.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "letscloud_ssh_key.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The key field is sensitive and not returned by the API
				ImportStateVerifyIgnore: []string{"key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSSHKeyResource_DuplicateLabel(t *testing.T) {
	// Configura o mock client
	mockClient := provider.NewLetsCloudClientMock()
	provider.MockLetsCloudClient = mockClient

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create first SSH key
			{
				Config: providerConfig + `
resource "letscloud_ssh_key" "test1" {
  label = "duplicate-label"
  key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIi41q4x71KvIciZZq1cU8DADrcqbn3ySVhgfTrZG1yUloU3Jq7Nn8/7YzCMr0CFSf/7ZgRGW9P1QxJZz8K3mWJ2z2K3Uf1x1+9Z9L3Qw== test@example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("letscloud_ssh_key.test1", "label", "duplicate-label"),
				),
			},
			// Try to create second SSH key with same label
			{
				Config: providerConfig + `
resource "letscloud_ssh_key" "test1" {
  label = "duplicate-label"
  key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIi41q4x71KvIciZZq1cU8DADrcqbn3ySVhgfTrZG1yUloU3Jq7Nn8/7YzCMr0CFSf/7ZgRGW9P1QxJZz8K3mWJ2z2K3Uf1x1+9Z9L3Qw== test@example.com"
}

resource "letscloud_ssh_key" "test2" {
  label = "duplicate-label"
  key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIi41q4x71KvIciZZq1cU8DADrcqbn3ySVhgfTrZG1yUloU3Jq7Nn8/7YzCMr0CFSf/7ZgRGW9P1QxJZz8K3mWJ2z2K3Uf1x1+9Z9L3Qw== test@example.com"
}
`,
				ExpectError: regexp.MustCompile(`SSH key with label 'duplicate-label' already exists`),
			},
		},
	})
}

func TestAccSSHKeyResource_InvalidKey(t *testing.T) {
	// Configura o mock client
	mockClient := provider.NewLetsCloudClientMock()
	provider.MockLetsCloudClient = mockClient

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "letscloud_ssh_key" "test" {
  label = "invalid-key"
  key   = "invalid-key-format"
}
`,
				ExpectError: regexp.MustCompile(`Invalid SSH key format`),
			},
		},
	})
}
