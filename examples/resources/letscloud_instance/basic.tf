terraform {
  required_providers {
    letscloud = {
      source = "letscloud-community/letscloud"
    }
  }
}

provider "letscloud" {
  # API token can be set via LETSCLOUD_API_TOKEN environment variable
  # api_token = "your-api-token"
}

resource "letscloud_instance" "basic" {
  label         = "basic-instance"
  plan_slug     = "1vcpu-1gb-10ssd"
  image_slug    = "ubuntu-24.04-x86_64"
  location_slug = "MIA1"
  hostname      = "basic-instance.example.com"
  password      = "SenhaSegura123!" # Required when no SSH keys are provided
} 