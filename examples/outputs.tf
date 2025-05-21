output "basic_ssh_key_id" {
  description = "The ID of the basic SSH key"
  value       = letscloud_ssh_key.basic.id
}

output "described_ssh_key_id" {
  description = "The ID of the SSH key with description"
  value       = letscloud_ssh_key.with_description.id
}

output "production_ssh_key_id" {
  description = "The ID of the production SSH key"
  value       = letscloud_ssh_key.production.id
}

output "staging_ssh_key_id" {
  description = "The ID of the staging SSH key"
  value       = letscloud_ssh_key.staging.id
}

output "jenkins_ssh_key_id" {
  description = "The ID of the Jenkins deployment SSH key"
  value       = letscloud_ssh_key.jenkins_deploy.id
}

output "all_ssh_key_ids" {
  description = "Map of all SSH key IDs"
  value = {
    basic      = letscloud_ssh_key.basic.id
    described  = letscloud_ssh_key.with_description.id
    production = letscloud_ssh_key.production.id
    staging    = letscloud_ssh_key.staging.id
    jenkins    = letscloud_ssh_key.jenkins_deploy.id
  }
} 