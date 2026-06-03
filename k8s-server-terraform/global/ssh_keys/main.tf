//for ssh into server, private key is local, 
//public key is registered with Hetzner Cloud and injected into the server upon creation
resource "hcloud_ssh_key" "admin" {
  name       = "admin_ssh_key"
  public_key = file(var.ssh_public_key)
}