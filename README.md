# certificate-sentinel-operator

##### Tested on OpenShift 4.8

The Certificate Sentinel Operator allows for the scanning and reporting of SSL Certificates within a Kubernetes/OpenShift cluster.

## Deploying the Operator

PROD IDK yet, hold on, i'm doing it dev mode deployment bby

### Development & Testing Deployment

Requires Golang 1.16+ and the DevelopmentTools dnf group.

```bash
# plz be `oc login`'d already
# also also need @DevelopmentTools & golang installed
git clone https://github.com/kenmoini/certificate-sentinel-operator
cd certificate-sentinel-operator/
make generate && make manifests && make install run
```

## Quickstart

### 1. Create a Namespace

*THIS!  IS!  KUBERNETES!*

So ya know, make a Namespace to get started...

```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: cert-sentinel
spec: {}
```

### 2. Create ServiceAccount

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: some-service-account
  namespace: cert-sentinel
```

### 3. Create ClusterRoleBindings

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: namespace-reader
rules:
  - verbs:
      - get
      - watch
      - list
    apiGroups:
      - ''
    resources:
      - namespaces
```

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: secret-reader
  - verbs:
      - get
      - watch
      - list
    apiGroups:
      - ''
    resources:
      - secrets
```

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: configmap-reader
  - verbs:
      - get
      - watch
      - list
    apiGroups:
      - ''
    resources:
      - configmaps
```

### 4. Create RoleBindings

Your ServiceAccount needs to be able to query a Namespace List and the Secrets/ConfigMaps in those namespaces - you do this with a RoleBinding to associate the ClusterRoles we just defined with the some-service-account ServiceAccount.

#### Targeted to only allow the sa/some-service-account to read in a specific namespace, cert-sentinel

For other namespaces you would need to duplicate and variate the `.metadata.namespace`

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: list-namespaces-cert-sentinel
  namespace: cert-sentinel
subjects:
- kind: ServiceAccount
  name: some-service-account
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: namespace-reader
  apiGroup: rbac.authorization.k8s.io
```

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-secrets-cert-sentinel
  namespace: cert-sentinel
subjects:
- kind: ServiceAccount
  name: some-service-account
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-secrets-cert-sentinel
  namespace: openshift-kube-scheduler-operator
subjects:
- kind: ServiceAccount
  name: some-service-account
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

#### Cluster-wide access to secrets

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: list-namespaces-cert-sentinel
subjects:
- kind: ServiceAccount
  name: some-service-account
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: namespace-reader
  apiGroup: rbac.authorization.k8s.io
```

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: read-secrets-cert-sentinel
subjects:
- kind: ServiceAccount
  name: some-service-account
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

### 5. Create a CertificateSentinel

Now that you have a Namespace and a ServiceAccount that has access to other Namespaces and the ability to read their Secrets/ConfigMaps, you can create the an object with the type of a Custom Resource Definition (CRD) supplied by this Operator, CertificateSentinel.

CertificateSentinel will watch allowed (authorization by {Cluster}RoleBindings to the .targets[*].ServiceAccount, Kubernetes/OpenShift RBAC) Namespaces+Secrets/ConfigMaps.

It will then scan them for PEM base64 encoded x509 Certificates, such as ones used for client/server/user authentication and service security via SSL/TLS.

If the Secrets/ConfigMaps contain a valid x509 Certificate, it will check the expiration date of those certificates and check if they are to be soon expiring and if so fires off an Alert.  Current Alert Types are `logger` (just stdout via operator-controller log function, eg you just ship logs to Elastic/Splunk/etc and query/match/alert there) and `smtp` for email notifications.

The following CertificateSentinel will watch the whole cluster for Certificates in Secrets, accessing those it can and Alerting via logger to upcoming expirations:

```yaml
apiVersion: config.polyglot.systems/v1
kind: CertificateSentinel
metadata:
  name: certificatesentinel-sample
  namespace: cert-operator
spec:
  alerts:
    - name: secrets-logger
      type: logger
  targets:
    - apiVersion: v1
      kind: Secret
      name: all-secrets
      namespaces:
        - '*'
      serviceAccount: some-service-account
      daysOut:
        - 30
        - 60
        - 90
        - 9001
        - 9000
```

Once the Operator has found a series of Certificates, it will log the discovered and expired certificates and reflect the data in the `CertificateSentinel.status` as such:

```yaml
apiVersion: config.polyglot.systems/v1
kind: CertificateSentinel
metadata:
  creationTimestamp: '2021-08-31T02:53:12Z'
  generation: 4
  managedFields:
    ...
  name: certificatesentinel-sample
  namespace: cert-operator
  resourceVersion: '10437267'
  uid: 17db6400-2d6e-4c87-8b95-0a645ce211b9
spec:
  alerts:
    - name: secrets-logger
      type: logger
  targets:
    - apiVersion: v1
      daysOut:
        - 30
        - 60
        - 90
        - 9001
        - 9000
      kind: Secret
      name: all-secrets
      namespaces:
        - '*'
      serviceAccount: some-service-account
status:
  certificatesAtRisk:
    - triggeredDaysOut:
        - 9001
        - 9000
      certificateAuthorityCommonName: openshift-service-serving-signer@1630120637
      name: kube-scheduler-operator-serving-cert
      expiration: '2023-08-28 03:17:39 +0000 UTC'
      kind: Secret
      dataKey: tls.crt
      isCertificateAuthority: false
      namespace: openshift-kube-scheduler-operator
      apiVersion: v1
  discoveredCertificates:
    - triggeredDaysOut:
        - 9001
        - 9000
      certificateAuthorityCommonName: openshift-service-serving-signer@1630120637
      name: kube-scheduler-operator-serving-cert
      expiration: '2023-08-28 03:17:39 +0000 UTC'
      kind: Secret
      dataKey: tls.crt
      isCertificateAuthority: false
      namespace: openshift-kube-scheduler-operator
      apiVersion: v1
```

#### Full CertificateSentinel YAML Spec Example

```yaml
apiVersion: config.polyglot.systems/v1
kind: CertificateSentinel
metadata:
  name: certificatesentinel-sample
  namespace: cert-operator
spec:
  alerts:
    - name: secrets-logger
      type: logger
      config: # optional on logger types, required for SMTP
        reportInterval: daily
  targets:
    - apiVersion: v1
      daysOut:
        - 30
        - 60
        - 90
        - 9001
        - 9000
      kind: Secret
      name: all-secrets
      namespaces:
        - '*'
      serviceAccount: some-service-account
status:
  certificatesAtRisk:
      # triggeredDaysOut []int
    - triggeredDaysOut:
        - 9001
        - 9000
      # certificateAuthorityCommonName string
      certificateAuthorityCommonName: openshift-service-serving-signer@1630120637
      # name string
      name: kube-scheduler-operator-serving-cert
      # expiration string
      expiration: '2023-08-28 03:17:39 +0000 UTC'
      # kind string
      kind: Secret
      # dataKey string
      dataKey: tls.crt
      # isCertificateAuthority bool
      isCertificateAuthority: false
      # namespace string
      namespace: openshift-kube-scheduler-operator
      # apiVersion string
      apiVersion: v1
  discoveredCertificates:
      # triggeredDaysOut []int
    - triggeredDaysOut:
        - 9001
        - 9000
      # certificateAuthorityCommonName string
      certificateAuthorityCommonName: openshift-service-serving-signer@1630120637
      # name string
      name: kube-scheduler-operator-serving-cert
      # expiration string
      expiration: '2023-08-28 03:17:39 +0000 UTC'
      # kind string
      kind: Secret
      # dataKey string
      dataKey: tls.crt
      # isCertificateAuthority bool
      isCertificateAuthority: false
      # namespace string
      namespace: openshift-kube-scheduler-operator
      # apiVersion string
      apiVersion: v1
```