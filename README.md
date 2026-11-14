# k8s-playground

A personal lab for building a single-node Kubernetes cluster on Hetzner Cloud from scratch, and for practising Kubernetes administration and security (CKS-style labs) on it.

The setup has three layers:

1. **Terraform** provisions the Hetzner infrastructure (network, firewalls, server).
2. **Ansible** installs and configures Kubernetes on the server (containerd, runc, CNI, kubeadm, Calico).
3. **Labs and apps** run on top of the cluster.

## Repository layout

| Path | Purpose |
| --- | --- |
| [k8s-server-terraform/](k8s-server-terraform) | Hetzner infrastructure as code |
| [k8s-server-ansible/](k8s-server-ansible) | Playbooks that turn the server into a Kubernetes node |
| [apps/k8s-api-client/](apps/k8s-api-client) | Small Go web app that queries the Kubernetes API |

## Prerequisites

- A Hetzner Cloud account and API token
- [Terraform](https://developer.hashicorp.com/terraform) and [Ansible](https://docs.ansible.com/)
- An SSH key pair for the server (the path is configured in [ansible.cfg](k8s-server-ansible/ansible.cfg))
- `kubectl` on your laptop (installed by playbook `42`)
- Go and Docker, only for the example app

## 1. Provision infrastructure (Terraform)

Located in [k8s-server-terraform/](k8s-server-terraform). Each directory is a separate root and they must be applied in this order:

1. [`global/`](k8s-server-terraform/global): global Hetzner artifacts (SSH keys) required by all servers.
2. [`core_network/`](k8s-server-terraform/core_network): parent network and subnet.
3. [`k8s_server/`](k8s-server-terraform/k8s_server): the `k8s-playground` server (Ubuntu 24.04, `cx23`, `nbg1`) with firewalls attached.

Reusable modules live in [`modules/`](k8s-server-terraform/modules): `core_network`, `hetzner_server`, `hetzner_firewall_ssh_only` and `hetzner_firewall_k8s_cluster`.

```bash
cd k8s-server-terraform/<root>
terraform init
terraform apply
```

Each root takes its inputs from `variables.tf`. For example, `hcloud_token` is the Hetzner API token. `k8s_server` also requires `admin_ssh_subnet_cidr`, which limits SSH access to that range, and falls back to your detected public IP when it is empty.

## 2. Configure Kubernetes (Ansible)

Located in [k8s-server-ansible/](k8s-server-ansible). The inventory ([hosts.yaml](k8s-server-ansible/inventory/hosts.yaml)) defines one host, `k8s-playground`. Versions and settings (Kubernetes, Calico, crictl, pod CIDR, and so on) are in [host_vars/k8s-playground.yaml](k8s-server-ansible/inventory/host_vars/k8s-playground.yaml).

Update the server address in the inventory and the SSH key and user in `ansible.cfg` to match your environment, then run the playbooks in numeric order from `k8s-server-ansible/`:

| Playbook | What it does |
| --- | --- |
| `00`, `01` | Install containerd and its systemd service |
| `10` | Install runc |
| `20` | Install CNI plugins |
| `30` | Swap configuration |
| `31` | Add the Kubernetes apt repository |
| `32`, `33`, `34` | Install kubelet, kubeadm and kubectl |
| `35` | Enable bridge networking |
| `40` | `kubeadm init` (uses the [kubeadm-init.yaml.j2](k8s-server-ansible/playbooks/kubeadm-init.yaml.j2) template) |
| `41` | Allow scheduling workloads on the control plane |
| `42` | Copy `admin.conf` to your laptop and install `kubectl` locally |
| `50`, `51` | Install and verify Calico |
| `52` | Calico cluster test: default-deny, egress and ingress NetworkPolicies (see [calico-cluster-test/](k8s-server-ansible/playbooks/calico-cluster-test)) |

```bash
cd k8s-server-ansible
ansible-playbook playbooks/00-containerd.yaml
# ...continue in order
```

### Accessing the cluster from your laptop

After playbook `42`, point `kubectl` at the fetched kubeconfig:

```bash
export KUBECONFIG=~/kubernetes/hetzner-playground/admin.conf
kubectl get nodes
kubectl proxy                                                       # API on http://localhost:8001
kubectl port-forward -n calico-system service/whisker 8081:8081     # Calico Whisker UI
```

See [from-laptop-README.md](k8s-server-ansible/playbooks/from-laptop-README.md) and [calico-README.md](k8s-server-ansible/playbooks/calico-README.md) for details.

## 3. Labs

[k8s-server-ansible/playbooks/cks-labs/](k8s-server-ansible/playbooks/cks-labs) holds hands-on notes and manifests for CKS (Certified Kubernetes Security Specialist) topics. Topics covered include:

- RBAC and ABAC, service accounts, authentication and kubeconfig
- Admission controllers, auditing, API groups, secrets
- Network policies, Cilium, Istio
- Security contexts, container sandboxing, host access, multi-tenancy isolation
- Kubelet API, node metadata, certificates, cluster install and upgrade
- Benchmarking with kube-bench and CIS-CAT Lite

Start with [useful-commands.md](k8s-server-ansible/playbooks/cks-labs/useful-commands.md).

## Example app: k8s-api-client

[apps/k8s-api-client/](apps/k8s-api-client) is a Go web app that lists pods through the Kubernetes API and renders them as HTML (for example at `/pods?namespace=default`). It uses in-cluster credentials when deployed and falls back to your kubeconfig locally. Starter RBAC and Deployment manifests are in [deploy-example/](apps/k8s-api-client/deploy-example).

See its [README](apps/k8s-api-client/README.md) for run, build and deploy instructions.

## Notes

This is a learning environment, not a production setup. The single node runs both the control plane and workloads, and several paths (home directory, SSH key location) are specific to the author's machine and need adjusting.
