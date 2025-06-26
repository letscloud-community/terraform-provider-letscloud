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

# Data source to fetch SSH key by ID
data "letscloud_ssh_key" "by_id" {
  id = "ssh-key-12345"
}

# Data source to fetch SSH key by label
data "letscloud_ssh_key" "by_label" {
  label = "production-key"
}

# Data source to fetch all SSH keys
data "letscloud_ssh_keys" "all" {}

# Create an instance using existing SSH key found by label
resource "letscloud_instance" "with_existing_key" {
  label         = "instance-with-existing-key"
  plan_slug     = "2vcpu-2gb-20ssd"
  image_slug    = "ubuntu-24.04-x86_64"
  location_slug = "MIA1"
  hostname      = "existing-key.example.com"
  
  # Use SSH key found by data source
  ssh_keys = [data.letscloud_ssh_key.by_label.id]
  
  password = "P@ssw0rd123!Secure"
}

# Output all SSH key information
output "ssh_key_by_id" {
  value = {
    id    = data.letscloud_ssh_key.by_id.id
    label = data.letscloud_ssh_key.by_id.label
  }
}

output "ssh_key_by_label" {
  value = {
    id    = data.letscloud_ssh_key.by_label.id
    label = data.letscloud_ssh_key.by_label.label
  }
}

output "all_ssh_keys" {
  value = data.letscloud_ssh_keys.all.ssh_keys
}

# Example: Find a specific key from all keys and use it
locals {
  production_key = [
    for key in data.letscloud_ssh_keys.all.ssh_keys :
    key if key.label == "production-key"
  ][0]
}

output "production_key_id" {
  value = local.production_key.id
} 