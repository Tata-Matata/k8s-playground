#### External identity provider
<details>
<summary>Show answer</summary>

1. When user runs <code>kubectl</code>, kubectl needs token to issue authenticated request to K8s server
2. kubectl itself does not know how to talk to an identity provider. Instead, it delegates that job to an authentication plugin (or uses a token that has already been obtained).
   
There are three common scenarios.

##### Scenario 1: Exec plugin (most common today)

Your kubeconfig contains something like:

```
users:
- name: my-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1
      command: kubelogin
      args:
      - get-token
```

Now when user runs <code>kubectl get pods</code>, kubectl sends "I need credentials" to exec plugin (kubelogin),
the plugin opens browser (Keycloak / Entra ID / Okta), the external OIDC provides JWT. Now kubectl can use it to send requests to K8s server.

##### Scenario 2: Token already exists

kubeconfig simply contains

```
users:
- name: my-user
  user:
    token: eyJhbGc...
```



##### Scenario 3: Cloud CLI

Managed Kubernetes services often use their cloud CLI. For example:

- Microsoft (az)
- Google (gcloud)
- Amazon Web Services (aws)


Example for Azure: kubeconfig may contain:

```
exec:
  command: kubelogin
```
or

```
exec:
  command: az
```

Then kubectl --> Azure CLI --> Microsoft Entra ID --> JWT --> kubectl

##### How does the plugin know where the IdP is?

For example, the plugin might know:

```
Issuer:
https://login.microsoftonline.com/<tenant-id>

Client ID:
abcd-1234

Redirect URI:
http://localhost:8000

```

These values are usually stored in: your kubeconfig, or the plugin's configuration, or provided by the cloud provider.

##### Example with Keycloak

Administrator tells everyone: "Download this kubeconfig."

```
users:
- name: employee
  user:
    exec:
      command: kubelogin
      args:
      - get-token
      - --oidc-issuer-url=https://keycloak.company.com/realms/dev
      - --oidc-client-id=kubernetes
```
When user runs <code>kubectl get pods</code>, the plugin knows exactly where to send you: https://keycloak.company.com/realms/dev

User logs in. The plugin receives a JWT. It gives the JWT to kubectl. Kubectl includes it in every request:

```
Authorization: Bearer eyJ...
```

</details>

#### AuthenticationConfiguration
<details>
<summary>Show answer</summary>



Recent Kubernetes versions support authenticating JWTs directly via AuthenticationConfiguration. You can define how claims map to usernames and groups without running a full OIDC provider.
AuthenticationConfiguration is a Kubernetes API object that lets you configure how the API server authenticates users, primarily via JWT authenticators (OIDC and other JWT issuers).
This configuration tells Kubernetes **how to interpret the token.**

```
apiVersion: apiserver.config.k8s.io/v1
kind: AuthenticationConfiguration

jwt:
- issuer:
    url: https://issuer.example.com
    audiences:
    - kubernetes
  claimMappings:
    username:
      claim: sub
    groups:
      claim: groups
```

This tells the API server:

- Trust JWTs issued by https://issuer.example.com.
- Accept tokens intended for the audience kubernetes.
- Use the sub claim as the Kubernetes username.
- Use the groups claim as Kubernetes groups.


Previously, you configured OIDC authentication with many API server flags:

```
--oidc-issuer-url=...
--oidc-client-id=...
--oidc-username-claim=sub
--oidc-groups-claim=groups
```

With AuthenticationConfiguration, everything is in one YAML file:

```
kube-apiserver --authentication-config=/etc/kubernetes/authentication-config.yaml
```

#### Mapping to groups

Suppose an external identity provider issues this JWT:
```
{
  "sub": "alice",
  "groups": [
    "developers",
    "oncall"
  ]
}
```

and your configuration contains

```
claimMappings:
  username:
    claim: sub
  groups:
    claim: groups

```

Then Kubernetes authenticates the request as

```
Username: alice
Groups:
- developers
- oncall
```

RBAC then evaluates those groups.

#### Relation to TokenReview

A TokenReview returns exactly the identity produced by the configured authenticator.

For a JWT like the example above, the response would contain something similar to:

```
status:
  authenticated: true
  user:
    username: alice
    groups:
    - developers
    - oncall
```

This makes TokenReview an excellent troubleshooting tool because it lets you verify whether the API server is extracting the expected username and groups before RBAC is evaluated.

</details>