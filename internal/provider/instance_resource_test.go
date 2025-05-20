// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInstanceResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test")
	}

	// Configura o mock client
	mockClient := &letsCloudClientMock{}
	mockLetsCloudClient = mockClient

	log.Printf("[DEBUG] Using mock client: %T", mockClient)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccInstanceResourceConfig("test-instance"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("letscloud_instance.test", "label", "test-instance"),
					resource.TestCheckResourceAttr("letscloud_instance.test", "hostname", "test-instance.example.com"),
					resource.TestCheckResourceAttr("letscloud_instance.test", "location_slug", "us-east-1"),
					resource.TestCheckResourceAttr("letscloud_instance.test", "plan_slug", "plan-1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "letscloud_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccInstanceResourceConfig("test-instance-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("letscloud_instance.test", "label", "test-instance-updated"),
					resource.TestCheckResourceAttr("letscloud_instance.test", "hostname", "test-instance-updated.example.com"),
				),
			},
		},
	})
}

func testAccInstanceResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "letscloud" {
  api_token = "mock-token-for-testing"
}

resource "letscloud_instance" "test" {
  label         = %[1]q
  hostname      = "%[1]s.example.com"
  location_slug = "us-east-1"
  plan_slug     = "plan-1"
  image_slug    = "ubuntu-20-04"
}
`, name)
}
