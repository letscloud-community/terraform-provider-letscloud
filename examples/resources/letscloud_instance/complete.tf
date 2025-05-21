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

# Create SSH key for instance access
resource "letscloud_ssh_key" "main" {
  label = "main-key"
  key   = file("~/.ssh/id_rsa.pub") # Replace with your public key path
}

# Create a complete instance with all available options
resource "letscloud_instance" "complete" {
  label         = "complete-instance"
  plan_slug     = "1vcpu-1gb-10ssd"     # Small instance with 1 vCPU, 1GB RAM, 10GB SSD
  image_slug    = "ubuntu-24.04-x86_64" # Ubuntu 24.04 LTS
  location_slug = "MIA1"                # Miami, USA
  hostname      = "complete-instance.example.com"

  # Authentication options - you can use either SSH keys or password
  ssh_keys = [letscloud_ssh_key.main.id]
  password = "SecurePassword123!" # Optional if using SSH keys, required if not

  # Lifecycle rules
  lifecycle {
    # Prevent accidental deletion
    prevent_destroy = true

    # Ignore changes to password
    ignore_changes = [password]
  }
}

# Output the instance details
output "instance_ipv4" {
  description = "The IPv4 address of the instance"
  value       = letscloud_instance.complete.ipv4
}

output "instance_ipv6" {
  description = "The IPv6 address of the instance (if available in the selected location)"
  value       = letscloud_instance.complete.ipv6
}

output "instance_state" {
  description = "The current state of the instance"
  value       = letscloud_instance.complete.state
} 