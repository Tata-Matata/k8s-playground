

### General strategy

<details>
<summary>Answer</summary>

First master, then workers. For workers, we can apply different strategies, for ex., upgrade one by one, moving workload to the rest


</details>


### Update first control plane

The instructions are very well documented in K8s docs.
These are just the general steps for Ubuntu.

While we are updating master nodes, the workloads on worker nodes still continue running, so we can not use API server for querying things or rolling out new resources, but existing apps will still be running.

#### Update version in package repository

<details>
<summary>Answer</summary>

update (minor) version to the desired one in **/etc/apt/sources.list.d/kubernetes.list**

<code>deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.35/deb/ / </code>

</details>


#### Determine patch release to upgrade to  

<details>
<summary>Answer</summary>


```
# Find the latest 1.35 version in the list.
# It should look like 1.35.x-*, where x is the latest patch.
sudo apt update
sudo apt-cache madison kubeadm

```

</details>



#### Update kubeadm

<details>
<summary>Answer</summary>


```
sudo apt-mark unhold kubeadm && \
sudo apt-get update && sudo apt-get install -y kubeadm='1.35.x-*' && \
sudo apt-mark hold kubeadm

kubeadm version

```

</details>


#### Use kubeadm to see the upgrade plan and then upgrade


<details>
<summary>Answer</summary>

This gives you a lot of valuable information about the cluster current state and what version you could be updating what components to

```
sudo kubeadm upgrade plan
sudo kubeadm upgrade apply v1.35.x

```

#### Manually upgrade your CNI provider plugin.





