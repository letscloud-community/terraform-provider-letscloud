# Terraform Provider for LetsCloud

This is the Terraform Provider for [LetsCloud](https://www.letscloud.io). It allows you to manage your LetsCloud resources using Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Installation

### Using Terraform Registry

Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    letscloud = {
      source = "letscloud-community/letscloud"
      version = "~> 1.0"
    }
  }
}
```

### Manual Installation

1. Clone the repository
```sh
git clone git@github.com:letscloud-community/terraform-provider-letscloud.git
```

2. Enter the repository directory
```sh
cd terraform-provider-letscloud
```

3. Build the provider
```sh
go build -o terraform-provider-letscloud
```

## Using the provider

To use the LetsCloud Provider in your Terraform configuration, you'll need to configure the provider with your API token. You can do this either in your Terraform configuration or by setting the `LETSCLOUD_API_TOKEN` environment variable.

```hcl
terraform {
  required_providers {
    letscloud = {
      source = "letscloud-community/letscloud"
      version = "~> 1.0"
    }
  }
}

provider "letscloud" {
  api_token = "your-api-token" # Optional: can also use LETSCLOUD_API_TOKEN env variable
}
```

### Example: Creating an Instance

```hcl
resource "letscloud_instance" "example" {
  label         = "example-instance"
  plan_slug     = "1vcpu-1gb-10ssd"
  image_slug    = "ubuntu-24.04-x86_64"
  location_slug = "MIA1"
  hostname      = "example.example.com"
  password      = "YourSecurePassword123!"
}
```

### Example: Creating an Instance with SSH Keys

```hcl
resource "letscloud_instance" "example_with_ssh" {
  label         = "example-instance-ssh"
  plan_slug     = "1vcpu-1gb-10ssd"
  image_slug    = "ubuntu-24.04-x86_64"
  location_slug = "MIA1"
  hostname      = "example-ssh.example.com"
  ssh_keys      = ["ssh-rsa AAAA..."]
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go build`. This will build the provider and put the provider binary in your current working directory.

```sh
go build -o terraform-provider-letscloud
```

### Running Tests

```sh
# Unit tests
go test ./...

# Acceptance tests
TF_ACC=1 go test ./...
```

## Documentation

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/letscloud-community/letscloud/latest/docs).

## License

This provider is licensed under the MIT License. See the LICENSE file for details.
