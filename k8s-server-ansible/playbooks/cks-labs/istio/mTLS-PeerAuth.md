### Strict mTLS

<details>
<summary>Answer</summary>

All workloads in production must use mTLS when communicating with each other.

  
```
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: default
  namespace: production
spec:
  mtls:
    mode: STRICT
```

</details>

### For backend workloads, require mTLS from their peers

<details>
<summary>Answer</summary>

  
```
spec:
  selector:
    matchLabels:
      app: backend
  mtls:
    mode: STRICT
```

</details>

### Check what Istio is actually enforcing on the pod

<details>
<summary>Answer</summary>

<code>istioctl x describe pod <pod-name> -n production </code>


</details>