variable "hcloud_token" {
  type        = string
  description = "Hetzner Cloud API token"
  sensitive   = true
}


locals {
  my_ip = var.admin_ssh_subnet_cidr != "" ? var.admin_ssh_subnet_cidr : trimspace(data.http.my_ip.response_body)
}

# from which servers admin is allowed to ssh into the k8s server
variable "admin_ssh_subnet_cidr" {
  type        = string
}

variable "host_offset_k8s_playground" {
  description = "Host offset for K8s server IP in the subnet"
  type        = number
  default     = 5

  validation {
    condition = alltrue([

      var.host_offset_k8s_playground >= var.min_ip,
      var.host_offset_k8s_playground <= var.max_ip
    ])

    error_message = "Host offset must be in range ${var.min_ip}-${var.max_ip} to avoid IP conflicts in the subnet."
  }

}

variable "min_ip" {
  description = "Minimum IP address in the subnet"
  type        = number
  default     = 1
}

variable "max_ip" {
  description = "Maximum IP address in the subnet"
  type        = number
  default     = 10

}

