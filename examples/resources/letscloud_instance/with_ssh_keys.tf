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

# First, create SSH keys
resource "letscloud_ssh_key" "key1" {
  name       = "my-key-1"
  public_key = file("~/.ssh/id_rsa.pub") # Replace with your public key path
}

resource "letscloud_ssh_key" "key2" {
  name       = "my-key-2"
  public_key = file("~/.ssh/id_ed25519.pub") # Replace with your public key path
}

# Then create the instance with the SSH keys
resource "letscloud_instance" "with_ssh" {
  label         = "ssh-instance"
  plan_slug     = "1vcpu-1gb-10ssd"
  image_slug    = "ubuntu-24.04-x86_64"
  location_slug = "MIA1"
  hostname      = "ssh-instance.example.com"
  ssh_keys = [
    letscloud_ssh_key.key1.id,
    letscloud_ssh_key.key2.id
  ]
} 