output "admin_ssh_public_key" {
  description = "will be injected by hetzner terraform on the server for admin access via ssh"
  value       = hcloud_ssh_key.admin.id

}