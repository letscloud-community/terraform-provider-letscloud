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

# Fetch all SSH keys
data "letscloud_ssh_keys" "all" {}

# Example: Create instances using different existing SSH keys
resource "letscloud_instance" "web_servers" {
  count = length(data.letscloud_ssh_keys.all.ssh_keys)

  label         = "web-server-${count.index + 1}"
  plan_slug     = "1vcpu-1gb-10ssd"
  image_slug    = "ubuntu-24.04-x86_64"
  location_slug = "MIA1"
  hostname      = "web-server-${count.index + 1}.example.com"

  # Use each SSH key for different instances
  ssh_keys = [data.letscloud_ssh_keys.all.ssh_keys[count.index].id]

  password = "P@ssw0rd123!Secure"
}

# Filter SSH keys by label pattern
locals {
  production_keys = [
    for key in data.letscloud_ssh_keys.all.ssh_keys :
    key if can(regex("prod", key.label))
  ]

  development_keys = [
    for key in data.letscloud_ssh_keys.all.ssh_keys :
    key if can(regex("dev", key.label))
  ]
}

# Output filtered keys
output "production_ssh_keys" {
  description = "SSH keys with 'prod' in the label"
  value       = local.production_keys
}

output "development_ssh_keys" {
  description = "SSH keys with 'dev' in the label"
  value       = local.development_keys
}

output "total_ssh_keys_count" {
  description = "Total number of SSH keys in the account"
  value       = length(data.letscloud_ssh_keys.all.ssh_keys)
} 