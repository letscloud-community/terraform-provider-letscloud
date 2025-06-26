// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sshkey_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/letscloud-community/terraform-provider-letscloud/internal/provider"
)

func TestAccSSHKeyDataSource(t *testing.T) {
	// Configura o mock client
	mockClient := provider.NewLetsCloudClientMock()
	provider.MockLetsCloudClient = mockClient

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create SSH key first
			{
				Config: providerConfig + `
resource "letscloud_ssh_key" "test" {
  label = "test-data-source-key"
  key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIi41q4x71KvIciZZq1cU8DADrcqbn3ySVhgfTrZG1yUloU3Jq7Nn8/7YzCMr0CFSf/7ZgRGW9P1QxJZz8K3mWJ2z2K3Uf1x1+9Z9L3Qw== test@example.com"
}

data "letscloud_ssh_key_lookup" "test_by_id" {
  id = letscloud_ssh_key.test.id
}

data "letscloud_ssh_key_lookup" "test_by_label" {
  label = letscloud_ssh_key.test.label
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check data source by ID
					resource.TestCheckResourceAttr("data.letscloud_ssh_key_lookup.test_by_id", "label", "test-data-source-key"),
					resource.TestCheckResourceAttrSet("data.letscloud_ssh_key_lookup.test_by_id", "id"),

					// Check data source by label
					resource.TestCheckResourceAttr("data.letscloud_ssh_key_lookup.test_by_label", "label", "test-data-source-key"),
					resource.TestCheckResourceAttrSet("data.letscloud_ssh_key_lookup.test_by_label", "id"),

					// Ensure both data sources return the same values
					resource.TestCheckResourceAttrPair("data.letscloud_ssh_key_lookup.test_by_id", "id", "data.letscloud_ssh_key_lookup.test_by_label", "id"),
					resource.TestCheckResourceAttrPair("data.letscloud_ssh_key_lookup.test_by_id", "label", "data.letscloud_ssh_key_lookup.test_by_label", "label"),
				),
			},
		},
	})
}

func TestAccSSHKeysDataSource(t *testing.T) {
	// Configura o mock client
	mockClient := provider.NewLetsCloudClientMock()
	provider.MockLetsCloudClient = mockClient

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create multiple SSH keys and test data source
			{
				Config: providerConfig + `
resource "letscloud_ssh_key" "test1" {
  label = "test-keys-1"
  key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIi41q4x71KvIciZZq1cU8DADrcqbn3ySVhgfTrZG1yUloU3Jq7Nn8/7YzCMr0CFSf/7ZgRGW9P1QxJZz8K3mWJ2z2K3Uf1x1+9Z9L3Qw== test1@example.com"
}

resource "letscloud_ssh_key" "test2" {
  label = "test-keys-2"
  key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIi41q4x71KvIciZZq1cU8DADrcqbn3ySVhgfTrZG1yUloU3Jq7Nn8/7YzCMr0CFSf/7ZgRGW9P1QxJZz8K3mWJ2z2K3Uf1x1+9Z9L3Qw== test2@example.com"
}

data "letscloud_ssh_keys" "all" {
  depends_on = [letscloud_ssh_key.test1, letscloud_ssh_key.test2]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.letscloud_ssh_keys.all", "ssh_keys.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("data.letscloud_ssh_keys.all", "ssh_keys.*", map[string]string{
						"label": "test-keys-1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("data.letscloud_ssh_keys.all", "ssh_keys.*", map[string]string{
						"label": "test-keys-2",
					}),
				),
			},
		},
	})
}
