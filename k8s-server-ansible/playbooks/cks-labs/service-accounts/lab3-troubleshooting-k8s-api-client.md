## Run 03-lab playbook under resources

The playbook deploys a web application in default namespace 
that queries K8s API and exposes on /pods endpoint list of pods.

After running the playbook, port forward the application with 

<code>kubectl port-forward --namespace default service/k8s-api-client 8888:8888</code>

and investigate the issue with the app on http://localhost:8888


<details>
<summary>Troubleshooting on the cluster</summary>

The application does not run at all. So we can check the pods status.
<code>kubectl get pods</code> does not show any pods

We need to check the deployment status with <code>kubectl get deploy</code>
We see that deploy *k8s-api-client* has 0 pods in status Ready
Inspect the deployment with 

<code>kubectl get deploy -o yaml > dep.yaml</code>

The yaml shows the error clearly 

```
 message: 'pods "k8s-api-client-8cbbc4fd8-" is forbidden: error looking up service
      account default/k8s-api-client: serviceaccount "k8s-api-client" not found'

```


</details>


