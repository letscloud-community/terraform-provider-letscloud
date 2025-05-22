# LetsCloud SSH Key Example

This example demonstrates how to create and manage SSH keys in LetsCloud using Terraform.

## Prerequisites

- Terraform 0.13 or later
- A LetsCloud account and API token
- An SSH key pair (public and private)

## Usage

1. Set your LetsCloud API token:
```bash
export LETSCLOUD_API_TOKEN="your-api-token"
```

2. Initialize Terraform:
```bash
terraform init
```

3. Review the execution plan:
```bash
terraform plan
```

4. Apply the configuration:
```bash
terraform apply
```

## Example Configuration

```hcl
# Create an SSH key
resource "letscloud_ssh_key" "main" {
  label = "main-key"
  key   = file("~/.ssh/id_rsa.pub")
}
```

## Security Best Practices

1. **Key Management**:
   - Store private keys securely
   - Use different keys for different environments
   - Rotate keys periodically
   - Use descriptive labels

2. **Access Control**:
   - Remove unused keys
   - Limit key access to necessary instances
   - Use unique labels for each key

## Outputs

- `ssh_key_id`: The unique identifier for the SSH key
- `ssh_key_label`: The label used to identify the SSH key
- `ssh_key_fingerprint`: The fingerprint of the SSH key (if available)

## Notes

- SSH keys cannot be updated after creation
- The label must be unique within your account
- The public key must be in the correct format

## Cleanup

To remove the created SSH key:

```bash
terraform destroy
``` 