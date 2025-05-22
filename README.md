# Terraform Provider for LetsCloud

This is a Terraform provider for managing resources in LetsCloud. It allows you to create and manage instances, SSH keys, and other resources using Terraform.

## Features

- Create and manage LetsCloud instances
- Manage SSH keys for secure instance access
- Support for multiple regions and instance types
- Secure authentication using API tokens

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.0
- [Go](https://golang.org/doc/install) >= 1.18 (to build the provider plugin)

## Installation

### Using Terraform Registry

Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    letscloud = {
      source  = "letscloud-community/letscloud"
      version = "1.1.0"
    }
  }
}

provider "letscloud" {
  api_token = "your-api-token" # or use LETSCLOUD_API_TOKEN environment variable
}
```

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/letscloud-community/terraform-provider-letscloud
cd terraform-provider-letscloud
```

2. Build the provider:
```bash
make build
```

3. Install the provider:
```bash
make install
```

## Quick Start

1. Set your LetsCloud API token:
```bash
export LETSCLOUD_API_TOKEN="your-api-token"
```

2. Create a new Terraform configuration:
```hcl
# Create an SSH key
resource "letscloud_ssh_key" "main" {
  label = "main-key"
  key   = file("~/.ssh/id_rsa.pub")
}

# Create an instance
resource "letscloud_instance" "basic" {
  label         = "basic-instance"
  plan_slug     = "1vcpu-1gb-10ssd"    # 1 vCPU, 1GB RAM, 10GB SSD
  image_slug    = "ubuntu-24.04-x86_64" # Ubuntu 24.04 LTS
  location_slug = "MIA1"               # Miami, USA
  hostname      = "basic-instance.example.com"
  ssh_keys      = [letscloud_ssh_key.main.id]
  password      = "P@ssw0rd123!Secure" # Must meet password requirements
  depends_on    = [letscloud_ssh_key.main]
}
```

3. Initialize Terraform:
```bash
terraform init
```

4. Apply the configuration:
```bash
terraform apply
```

## Available Resources

### letscloud_instance

Manages a LetsCloud instance. Supports:
- Multiple instance types (plans)
- Various operating systems (images)
- Multiple regions (locations)
- SSH key authentication
- Custom hostnames

Example:
```hcl
resource "letscloud_instance" "example" {
  label         = "example-instance"
  plan_slug     = "1vcpu-1gb-10ssd"
  image_slug    = "ubuntu-24.04-x86_64"
  location_slug = "MIA1"
  hostname      = "example.example.com"
  ssh_keys      = [letscloud_ssh_key.main.id]
  password      = "P@ssw0rd123!Secure"
}
```

### letscloud_ssh_key

Manages SSH keys for secure instance access. Supports:
- Multiple SSH keys per account
- Unique labels for key identification
- Public key format validation

Example:
```hcl
resource "letscloud_ssh_key" "main" {
  label = "main-key"
  key   = file("~/.ssh/id_rsa.pub")
}
```

## Best Practices

1. **Security**:
   - Use SSH keys for instance access
   - Follow password requirements
   - Use unique labels for resources
   - Store API tokens securely

2. **Resource Management**:
   - Use descriptive labels
   - Follow naming conventions
   - Implement proper dependencies
   - Monitor resource usage

3. **Configuration**:
   - Use environment variables for sensitive data
   - Implement proper state management
   - Use version control for configurations
   - Document your infrastructure

## Examples

Check the [examples](./examples) directory for complete examples:
- [Basic Instance](./examples/resources/letscloud_instance/basic.tf)
- [Instance with SSH Keys](./examples/resources/letscloud_instance/with_ssh_keys.tf)
- [Complete Setup](./examples/resources/letscloud_instance/complete.tf)
- [SSH Key Management](./examples/resources/letscloud_ssh_key)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please:
1. Check the [documentation](./docs)
2. Open an issue in the GitHub repository
3. Contact LetsCloud support

## Acknowledgments

- [LetsCloud](https://letscloud.io) for providing the API
- [Terraform](https://terraform.io) for the provider framework
- All contributors who have helped improve this provider
