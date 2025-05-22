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

# Create an SSH key for secure instance access
resource "letscloud_ssh_key" "main" {
  label = "main-key"
  key   = file("~/.ssh/id_rsa.pub")
}

# Output the SSH key ID for reference
output "ssh_key_id" {
  value = letscloud_ssh_key.main.id
} 