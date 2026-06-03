module "k8s_playground_server" {
  source = "../modules/hetzner_server"

  #server config
  server_name     = "k8s-playground"
  server_location = "nbg1"
  os_image        = "ubuntu-22.04"
  server_type     = "cx23"

  #admin ssh access
  ssh_key_ids = [data.terraform_remote_state.global_ssh_keys.outputs.admin_ssh_public_key]


  #network config
  //enable public IP only for bastion 
  public_ip_enabled = true

  // Hetzner expects here ID of the parent network
  parent_network_id = data.terraform_remote_state.core_network.outputs.parent_network_id

  //But Hetzner also expects server IP that belongs to a subnet of the network
  subnet_cidr = data.terraform_remote_state.core_network.outputs.subnet_cidr

  // e.g., for 10.50.1.5 use offset 5
  host_offset = var.host_offset_k8s_playground

  //for attaching firewall
  server_labels = {
    role = "k8s_playground"
  }

}

module "hetzner_firewall_ssh_only" {
  source = "../modules/hetzner_firewall_ssh_only"

  #access to the server only from this subnet via ssh
  admin_ssh_subnet_cidr =  "${local.my_ip}/32"

}
