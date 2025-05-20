resource "letscloud_instance" "example" {
  name   = "web-server-1"
  size   = "small"
  image  = "ubuntu-22-04"
  region = "us-east"

  ssh_keys = [
    "ssh-key-1",
    "ssh-key-2"
  ]
} 