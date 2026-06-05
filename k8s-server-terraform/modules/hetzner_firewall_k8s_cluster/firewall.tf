//Restricts access to K8s API

resource "hcloud_firewall" "hetzner_k8s_cluster_fw" {
  name = "hetzner_k8s_cluster"



  # K8s API only accessible from admin home network 
  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "6443"
    source_ips = [var.admin_ssh_subnet_cidr]
  }

  apply_to {
    label_selector = "role=k8s_playground"
  }

}
