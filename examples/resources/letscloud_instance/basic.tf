terraform {
  required_providers {
    letscloud = {
      source  = "letscloud-community/letscloud"
      version = "1.1.0"
    }
  }
}

provider "letscloud" {
  # API token can be set via LETSCLOUD_API_TOKEN environment variable
  api_token = "your-api-token"
}

# Basic instance with minimal configuration
resource "letscloud_instance" "basic" {
  label         = "basic-instance"
  plan_slug     = "1vcpu-1gb-10ssd"    # 1 vCPU, 1GB RAM, 10GB SSD
  image_slug    = "ubuntu-24.04-x86_64" # Ubuntu 24.04 LTS
  location_slug = "MIA1"               # Miami, USA
  hostname      = "basic-instance.example.com"
  password      = "P@ssw0rd123!Secure" # Must meet password requirements
} 