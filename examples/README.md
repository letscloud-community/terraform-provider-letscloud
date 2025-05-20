# Examples

This directory contains examples that demonstrate how to use the LetsCloud provider. These examples can be used for documentation and can also be run/tested manually via the Terraform CLI.

## Instance Examples

### Basic Example (`resources/letscloud_instance/basic.tf`)
A minimal example showing how to create an instance with just the required fields. This example uses password authentication.

### SSH Keys Example (`resources/letscloud_instance/with_ssh_keys.tf`)
Demonstrates how to create SSH keys and use them to authenticate with your instance. This is the recommended approach for secure access.

### Complete Example (`resources/letscloud_instance/complete.tf`)
A comprehensive example showing all available options for the instance resource, including:
- Instance configuration (plan, image, location)
- Authentication (SSH keys and password)
- Lifecycle rules
- Output values

## Provider Configuration

The provider configuration examples show how to:
- Configure the LetsCloud provider
- Set up authentication using API tokens
- Use environment variables for sensitive data

## Running the Examples

To run any of these examples:

1. Set your API token:
   ```bash
   export LETSCLOUD_API_TOKEN="your-api-token"
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Review the planned changes:
   ```bash
   terraform plan
   ```

4. Apply the configuration:
   ```bash
   terraform apply
   ```

## Notes

- Replace placeholder values (like SSH key paths) with your actual values
- The examples use the Miami (MIA1) location by default, but you can change it to any available location
- IPv6 is only available in certain locations that support it
- When using password authentication, ensure the password is at least 8 characters long
