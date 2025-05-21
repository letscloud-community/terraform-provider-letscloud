terraform {
  required_providers {
    letscloud = {
      source  = "local/letscloud/letscloud"
      version = "1.0.0"
    }
  }
}

provider "letscloud" {
  # API token can be set via LETSCLOUD_API_TOKEN environment variable
  # api_token = "your-token-here"
}

# Example instance
resource "letscloud_instance" "example" {
  label         = "example-instance"
  plan_slug     = "1vcpu-1gb-10ssd"
  image_slug    = "ubuntu-24.04-x86_64"
  location_slug = "MIA1"
  hostname      = "example.example.com"
  password      = "SecurePassword123!" # Required when no SSH keys are provided
} 