# Output the SSH key ID
output "ssh_key_id" {
  description = "The unique identifier for the SSH key"
  value       = letscloud_ssh_key.main.id
}

# Output the SSH key label
output "ssh_key_label" {
  description = "The label used to identify the SSH key"
  value       = letscloud_ssh_key.main.label
}

# Output the SSH key fingerprint (if available)
output "ssh_key_fingerprint" {
  description = "The fingerprint of the SSH key (if available)"
  value       = letscloud_ssh_key.main.id
}

output "admin_key_id" {
  description = "ID of the admin SSH key"
  value       = letscloud_ssh_key.admin.id
}

output "admin_key_label" {
  description = "Label of the admin SSH key"
  value       = letscloud_ssh_key.admin.label
}

output "developer_key_id" {
  description = "ID of the developer SSH key"
  value       = letscloud_ssh_key.developer.id
}

output "developer_key_label" {
  description = "Label of the developer SSH key"
  value       = letscloud_ssh_key.developer.label
} 