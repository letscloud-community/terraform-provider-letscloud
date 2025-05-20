terraform {
  required_providers {
    letscloud = {
      source = "letscloud-community/letscloud"
    }
  }
}

provider "letscloud" {
  # Configure the LetsCloud Provider
  api_token = "your-api-token" # or use LETSCLOUD_API_KEY env variable
}
