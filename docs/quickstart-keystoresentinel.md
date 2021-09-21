# Quickstart - KeystoreSentinel

This document will help demonstrate the workflow for using this operator to scan a cluster for expiring x509 certificates.

## 1. Create a Namespace

*THIS!  IS!  KUBERNETES!*

So, ya know, make a Namespace to get started...

```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: cert-sentinel
  labels:
    polyglot.systems/certificate-sentinel-namespace: "true"
    polyglot.systems/keystore-sentinel-namespace: "true"
spec: {}
```

## 2. Create ServiceAccount

This ServiceAccount will be the RBAC object that will access the K8s/OCP API in order to scan the cluster for Certificates

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: some-service-account
  namespace: cert-sentinel
```

## 3. Create ClusterRoleBindings

These ClusterRoles will define the RBAC permissions required order to access Namespaces, Secrets, and ConfigMaps

#### namespace-reader

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

#### secret-reader

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: secret-reader
rules:
  - verbs:
      - get
      - watch
      - list
    apiGroups:
      - ''
    resources:
      - secrets
```

#### configmap-reader

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: configmap-reader
rules:
  - verbs:
      - get
      - watch
      - list
    apiGroups:
      - ''
    resources:
      - configmaps
```

#### certificate-reader

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: certificate-reader
rules:
  - verbs:
      - get
      - watch
      - list
    apiGroups:
      - ''
    resources:
      - certificates
```

#### sentinel-reader

This ClusterRole has all the objects defined together

```yaml
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: sentinel-reader
rules:
  - verbs:
      - get
      - watch
      - list
    apiGroups:
      - ''
    resources:
      - certificates
      - configmaps
      - namespaces
      - secrets
```

## 4. Create RoleBindings

Your ServiceAccount needs to be able to query a Namespace List and the Secrets/ConfigMaps in those namespaces - you do this with a RoleBinding to associate the ClusterRoles we just defined with the some-service-account ServiceAccount.

### Targeted to only allow the sa/some-service-account to read in a specific namespace, cert-sentinel

For other namespaces you would need to duplicate and variate the `.metadata.namespace`

#### Allow the serviceaccount/some-service-account to access Namespaces on the cluster

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
  namespace: cert-sentinel
roleRef:
  kind: ClusterRole
  name: namespace-reader
  apiGroup: rbac.authorization.k8s.io
```

#### Allow the serviceaccount/some-service-account to access Secrets in namespace/cert-sentinel

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
  namespace: cert-sentinel
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

#### Allow the serviceaccount/some-service-account to access Secrets in namespace/openshift-kube-scheduler-operator

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
  namespace: cert-sentinel
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

### Cluster-wide access to secrets

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: list-namespaces-cert-sentinel
subjects:
- kind: ServiceAccount
  name: some-service-account
  namespace: cert-sentinel
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
  namespace: cert-sentinel
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

### Cluster-wide access to all needed objects

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: read-objects-cert-sentinel
subjects:
- kind: ServiceAccount
  name: some-service-account
  namespace: cert-sentinel
roleRef:
  kind: ClusterRole
  name: sentinel-reader
  apiGroup: rbac.authorization.k8s.io
```

## 5. Create a KeystoreSentinel

Now that you have a Namespace and a ServiceAccount that has access to other Namespaces and the ability to read their Secrets/ConfigMaps, you can create the an object with the type of a Custom Resource Definition (CRD) supplied by this Operator, KeystoreSentinel.

KeystoreSentinel will watch allowed (authorization by {Cluster}RoleBindings to the .targets[*].ServiceAccount, Kubernetes/OpenShift RBAC) Namespaces+Secrets/ConfigMaps.

It will then scan them for Java Keystores and if parsed, will scan for x509 Certificates, such as ones used for providing trusted root stores to Java applications.

If the Secrets/ConfigMaps contain a valid Java Keystore that has valid certificates, it will check the expiration date of those certificates and check if they are to be soon expiring and if so fires off an Alert.  Current Alert Types are `logger` (just stdout via operator-controller log function, eg you just ship logs to Elastic/Splunk/etc and query/match/alert there) and `smtp` for email notifications.

The following KeystoreSentinel will watch the whole cluster for Certificates in Secrets, accessing those it can and Alerting via logger to upcoming expirations:

```yaml
apiVersion: config.polyglot.systems/v1
kind: KeystoreSentinel
metadata:
  name: keystoresentinel-sample
  namespace: cert-sentinel
spec:
  alert:
    name: secrets-logger
    type: logger
    config:
      reportInterval: debug
  target:
    apiVersion: v1
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
    keystorePassword:
      plaintext: changeit
      type: plaintext
```

Once the Operator has found a series of Certificates, it will log the discovered and expired certificates and reflect the data in the `KeystoreSentinel.status` as such:

```yaml
apiVersion: config.polyglot.systems/v1
kind: KeystoreSentinel
metadata:
  creationTimestamp: '2021-08-31T02:53:12Z'
  generation: 4
  managedFields:
    ...
  name: keystoresentinel-sample
  namespace: cert-sentinel
  resourceVersion: '10437267'
  uid: 17db6400-2d6e-4c87-8b95-0a645ce211b9
spec:
  alert:
    name: secrets-logger
    type: logger
    config:
      reportInterval: debug
  target:
    apiVersion: v1
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
    keystorePassword:
      plaintext: changeit
      type: plaintext
status:
  discoveredKeystoreCertificates:
    - triggeredDaysOut:
        - 9001
        - 9000
      certificateAuthorityCommonName: Example Labs Intermediate Certificate Authority
      name: keystore-secret-test
      keystoreAlias: examplelabsica
      expiration: '2024-09-06 00:00:00 +0000 UTC'
      kind: Secret
      dataKey: jks
      commonName: Example Labs Signing Certificate Authority
      isCertificateAuthority: true
      namespace: cert-sentinel
      apiVersion: v1
    - triggeredDaysOut:
        - 9001
        - 9000
      certificateAuthorityCommonName: Certificate Authority
      name: keystore-secret-test
      keystoreAlias: idmca
      expiration: '2041-04-05 16:25:15 +0000 UTC'
      kind: Secret
      dataKey: jks
      commonName: Certificate Authority
      isCertificateAuthority: true
      namespace: cert-sentinel
      apiVersion: v1
  expiringCertificates: 2
  keystoresAtRisk: 1
  lastReportSent: 1632232508
  totalKeystoresFound: 1
```

- For further information on the KeystoreSentinel YAML spec, see [full_yaml_spec-KeystoreSentinel.md](./full_yaml_spec-KeystoreSentinel.md)
- For additional examples of KeystoreSentinels and Keystore + Password sets represented in Secrets/ConfigMaps to test/use, see [examples/](../examples/)