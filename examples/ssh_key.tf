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

# Basic SSH key example
resource "letscloud_ssh_key" "basic" {
  name       = "basic-key"
  public_key = "ssh-rsa AAAA... your-public-key-content"
}

# SSH key with description
resource "letscloud_ssh_key" "with_description" {
  name        = "described-key"
  public_key  = "ssh-rsa AAAA... your-public-key-content"
  description = "This key is used for development servers"
}

# Multiple SSH keys for different environments
resource "letscloud_ssh_key" "production" {
  name        = "prod-key"
  public_key  = "ssh-rsa AAAA... production-key-content"
  description = "SSH key for production servers"
}

resource "letscloud_ssh_key" "staging" {
  name        = "staging-key"
  public_key  = "ssh-rsa AAAA... staging-key-content"
  description = "SSH key for staging servers"
}

# SSH key with a more descriptive name
resource "letscloud_ssh_key" "jenkins_deploy" {
  name        = "jenkins-deploy-key-2024"
  public_key  = "ssh-rsa AAAA... jenkins-key-content"
  description = "Jenkins deployment key for CI/CD pipelines"
} 