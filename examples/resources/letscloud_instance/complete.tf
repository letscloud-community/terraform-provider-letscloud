terraform {
  required_providers {
    letscloud = {
      source  = "letscloud-community/letscloud"
      version = "1.0.1"
    }
  }
}

provider "letscloud" {
  # API token can be set via LETSCLOUD_API_TOKEN environment variable
  api_token = "your-api-token"
}

# Create multiple SSH keys for different access levels
resource "letscloud_ssh_key" "admin" {
  label = "admin-key"
  key   = file("~/.ssh/admin.pub")
}

resource "letscloud_ssh_key" "developer" {
  label = "developer-key"
  key   = file("~/.ssh/developer.pub")
}

# Create a production instance with multiple SSH keys
resource "letscloud_instance" "production" {
  label         = "prod-instance"
  plan_slug     = "4vcpu-4gb-40ssd"     # 4 vCPU, 4GB RAM, 40GB SSD
  image_slug    = "ubuntu-24.04-x86_64" # Ubuntu 24.04 LTS
  location_slug = "MIA1"                # Miami, USA
  hostname      = "prod-instance.example.com"
  ssh_keys = [
    letscloud_ssh_key.admin.id,
    letscloud_ssh_key.developer.id
  ]
  password = "P@ssw0rd123!Secure" # Must meet password requirements
  depends_on = [
    letscloud_ssh_key.admin,
    letscloud_ssh_key.developer
  ]
}

# Create a staging instance with developer access only
resource "letscloud_instance" "staging" {
  label         = "staging-instance"
  plan_slug     = "2vcpu-2gb-20ssd"     # 2 vCPU, 2GB RAM, 20GB SSD
  image_slug    = "ubuntu-24.04-x86_64" # Ubuntu 24.04 LTS
  location_slug = "MIA1"                # Miami, USA
  hostname      = "staging-instance.example.com"
  ssh_keys      = [letscloud_ssh_key.developer.id]
  password      = "P@ssw0rd123!Secure" # Must meet password requirements
  depends_on    = [letscloud_ssh_key.developer]
}

# Output the instance IPs for easy access
output "production_ip" {
  value = letscloud_instance.production.ipv4
}

output "staging_ip" {
  value = letscloud_instance.staging.ipv4
} 