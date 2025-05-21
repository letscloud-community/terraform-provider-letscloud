terraform {
  required_providers {
    letscloud = {
      source = "letscloud-community/letscloud"
    }
  }
}

provider "letscloud" {
  api_token = "your-api-token"
}

# Create a basic SSH key
resource "letscloud_ssh_key" "example" {
  label = "example-key"
  key   = "ssh-rsa AAAA... your-public-key-content"
} 