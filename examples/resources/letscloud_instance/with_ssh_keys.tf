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

# Create an SSH key
resource "letscloud_ssh_key" "main" {
  label = "main-key"
  key   = file("~/.ssh/id_rsa.pub")
}

# Create an instance with SSH key access
resource "letscloud_instance" "with_ssh" {
  label         = "ssh-instance"
  plan_slug     = "2vcpu-2gb-20ssd"     # 2 vCPU, 2GB RAM, 20GB SSD
  image_slug    = "ubuntu-24.04-x86_64" # Ubuntu 24.04 LTS
  location_slug = "MIA1"                # Miami, USA
  hostname      = "ssh-instance.example.com"
  ssh_keys      = [letscloud_ssh_key.main.id]
  password      = "P@ssw0rd123!Secure" # Must meet password requirements
  depends_on    = [letscloud_ssh_key.main]
} 