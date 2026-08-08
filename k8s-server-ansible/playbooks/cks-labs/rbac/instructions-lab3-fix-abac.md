## Run 03-lab-fix-abac playbook

The playbook creates SA **abac** and configures ABAC authorization in kube-apiserver, but API server is not responding after the changes.
The ABAC policy gives the SA **abac** permission to list deployments

Initial API server manifest file (before ABAC changes) is backed up at /root/labs/abac/kube-apiserver-backup.yaml

IMPORTANT NOTE: there are several issues to be fixed, so once you get API server up and running again,
check if SA abac can get deployments

## Fix the cluster and ABAC

<details>
<summary>Show answer</summary>




</details>