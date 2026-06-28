## Test if default service account can list pods and create a configmap

<details>
<summary>Show answer</summary>

#### Simple approach
```
kubectl auth can-i get pods --as=system:serviceaccount:default:default
```
This however only tells you yes or no, without giving any details


#### More verbose approach 

<code>kubectl --as=system:serviceaccount:default:default get pods</code>

Gives more information like: Error from server (Forbidden): pods is forbidden: User "system:serviceaccount:default:default" cannot list resource "pods" in API group "" in the namespace "default


For more verbose information add
<code>kubectl --as=system:serviceaccount:default:default get pods -v=8</code>


#### More complicated and fun - from within the pod running with default service account

1. Create nginx pod and exec into it
```
kubectl run ng --image nginx
kubectl exec -it ng -- /bin/bash
```


2. Prepare token, certificate, namespace for future tests. 

The paths can be found in volumeMount if you inspect the pod 

```
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
CACERT="/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
NAMESPACE=$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace)

```

1. Test listing pods

```
curl --cacert $CACERT -H "Authorization: Bearer $TOKEN" https://kubernetes.default.svc/api/v1/namespaces/$NAMESPACE/pods

```

This shows as expected 

```
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "pods is forbidden: User \"system:serviceaccount:default:default\" cannot list resource \"pods\" in API group \"\" in the namespace \"default\"",
  "reason": "Forbidden",
  "details": {
    "kind": "pods"
  },
  "code": 403
```

Default service account does not have permissions on K8s resources. 
You can verify this with <code>kubectl auth can-i --list --as=system:serviceaccount:default:default</code>

2. Add permissions to list pods


2. Test creating configMap

curl \
  --cacert $CACERT \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -X POST \
  --data @/tmp/cm.json \
  https://kubernetes.default.svc/api/v1/namespaces/$NAMESPACE/configmaps


Notice how the API server tells you exactly which identity it authenticated.

Test 2: Create a ConfigMap

A real controller would POST JSON.

For example:

cat <<EOF >/tmp/cm.json
{
  "apiVersion":"v1",
  "kind":"ConfigMap",
  "metadata":{
    "name":"test-cm"
  },
  "data":{
    "hello":"world"
  }
}
EOF

Then:



Again you'll either get:

201 Created (or the created object)
403 Forbidden
Test 3: Delete the ConfigMap
curl \
  --cacert $CACERT \
  -H "Authorization: Bearer $TOKEN" \
  -X DELETE \
  https://kubernetes.default.svc/api/v1/namespaces/$NAMESPACE/configmaps/test-cm
Watch the RBAC decisions

Suppose your Pod runs as:

system:serviceaccount:default:default

When it makes:

GET /api/v1/namespaces/default/pods

the API server performs these steps:

1. Verify JWT signature
2. Authenticate:
      User:
        system:serviceaccount:default:default

      Groups:
        system:authenticated
        system:serviceaccounts
        system:serviceaccounts:default

3. Authorize:
      Can this user/groups list Pods?

4. Return:
      200 OK
or
      403 Forbidden

This is the complete authentication → authorization flow you've been studying.

An even more realistic approach: kubectl inside the Pod

Many Kubernetes controllers and operators use the official Kubernetes client libraries rather than raw HTTP. Those libraries behave just like kubectl: they automatically detect they're running inside a Pod, read the ServiceAccount token and CA certificate from /var/run/secrets/kubernetes.io/serviceaccount, and talk to https://kubernetes.default.svc.

You can mimic that by copying kubectl into a Pod (or using an image that already contains it):

kubectl get pods
kubectl create configmap test --from-literal=a=b

without providing a kubeconfig. kubectl will use the in-cluster configuration automatically. This is a great way to see the exact same RBAC behavior while using the same client library that many real applications rely on.

</details>