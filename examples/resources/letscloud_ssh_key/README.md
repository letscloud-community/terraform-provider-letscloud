# SSH Key Resource Example

This example demonstrates how to create and manage SSH keys using the LetsCloud provider.

## Usage

1. Replace `your-api-token` with your actual LetsCloud API token
2. Replace `your-public-key-content` with your actual SSH public key content
3. Run the following commands:

```bash
terraform init
terraform plan
terraform apply
```

## Example Configuration

The example creates a basic SSH key with:
- A name
- A public key
- An optional description

## Outputs

After applying the configuration, you can get the SSH key ID using:

```bash
terraform output ssh_key_id
```

## Cleanup

To remove the created SSH key:

```bash
terraform destroy
``` 