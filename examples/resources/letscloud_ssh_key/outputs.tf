output "ssh_key_id" {
  description = "The ID of the created SSH key"
  value       = letscloud_ssh_key.example.id
}

output "ssh_key_name" {
  description = "The name of the created SSH key"
  value       = letscloud_ssh_key.example.name
}

output "ssh_key_description" {
  description = "The description of the created SSH key"
  value       = letscloud_ssh_key.example.description
} 