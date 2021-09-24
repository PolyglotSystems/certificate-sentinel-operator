# Deployment - Operator

## Deploy to Red Hat OpenShift

If you already have an OpenShift cluster or Kubernetes cluster with OLM deployed, you can simply apply the Operator CatalogSource and then create the needed Subscription from there.

```bash
# Create the Operator CatalogSource in the openshift-marketplace Namespace/Project
oc apply -f https://raw.githubusercontent.com/PolyglotSystems/certificate-sentinel-operator/main/deploy/operator-install/01-catalogsource.yaml -n openshift-marketplace

# Create a Subscription for the Operator (this installs the Operator)
oc apply -f https://raw.githubusercontent.com/PolyglotSystems/certificate-sentinel-operator/main/deploy/operator-install/02-subscription.yaml -n openshift-operators
```

# Deployment - Examples

This directory houses some example Kubernetes manifests that when applied to a cluster will quickly deploy a simple demonstration of CertificateSentinel and KeystoreSentinel CRDs provided by this Certificate Sentinel Operator

## Initial Setup

The [examples-initial-setup](./example-initial-setup) directory houses a set of Kubernetes assets that will set up the basics required, such as

- [01-namespace.yaml](./example-initial-setup/01-namespace.yaml) - This will create a `cert-sentinel` Namespace to work in
- [02-serviceaccount.yaml](./example-initial-setup/02-serviceaccount.yaml) - This provides an atomic ClusterRole that allows a User/Group/ServiceAccount to read Secrets on the cluster

With these defined, you may then create a ClusterRole (or few) that defines access capabilities of your target User/Group/ServiceAccount

## RBAC ClusterRoles

Information on Roles and ClusterRoles: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-and-clusterrole

The [examples-rbac-clusterroles](./example-rbac-clusterroles) directory houses a number of Kubernetes ClusterRoles that will define access to a number of objects:

- [01-namespace-reader-clusterrole.yaml](./example-rbac-clusterroles/01-namespace-reader-clusterrole.yaml) - This provides an atomic ClusterRole that allows a User/Group/ServiceAccount to read Namespaces on the cluster
- [02-secret-reader-clusterrole.yaml](./example-rbac-clusterroles/02-secret-reader-clusterrole.yaml) - This provides an atomic ClusterRole that allows a User/Group/ServiceAccount to read Secrets on the cluster
- [03-configmap-reader-clusterrole.yaml](./example-rbac-clusterroles/03-configmap-reader-clusterrole.yaml) - This provides an atomic ClusterRole that allows a User/Group/ServiceAccount to read ConfigMaps on the cluster
- [04-certificate-reader-clusterrole.yaml](./example-rbac-clusterroles/04-certificate-reader-clusterrole.yaml) - This provides an atomic ClusterRole that allows a User/Group/ServiceAccount to read Certificates on the cluster
- [05-sentinel-reader-clusterrole.yaml](./example-rbac-clusterroles/05-sentinel-reader-clusterrole.yaml) - This provides a combined ClusterRole that allows a User/Group/ServiceAccount to read Namespaces, Secrets, ConfigMaps, and Certificates on the cluster

With these defined, you may then create a ClusterRoleBinding or RoleBinding to your target User/Group/ServiceAccount that provides bound access to any of these ClusterRole definitions.

## RBAC RoleBindings

Information on RoleBindings: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding

The [examples-rbac-rolebindings](./example-rbac-rolebindings) directory houses a number of Kubernetes RoleBindings:

- [01-list-namespaces-rolebinding.yaml](./example-rbac-clusterroles/01-list-namespaces-rolebinding.yaml) - This provides an atomic RoleBinding in `namespace/cert-sentinel` that binds the `ServiceAccount/some-service-account` to the `ClusterRole/namespace-reader` to read Namespaces
- [02-read-secrets-rolebinding.yaml](./example-rbac-clusterroles/02-read-secrets-rolebinding.yaml) - This provides an atomic RoleBinding in `namespace/cert-sentinel` that binds the `ServiceAccount/some-service-account` to the `ClusterRole/secret-reader` to read Secrets
- [03-read-configmaps-rolebinding.yaml](./example-rbac-clusterroles/03-read-configmaps-rolebinding.yaml) - This provides an atomic RoleBinding in `namespace/cert-sentinel` that binds the `ServiceAccount/some-service-account` to the `ClusterRole/configmap-reader` to read ConfigMaps
- [04-read-certificates-rolebinding.yaml](./example-rbac-clusterroles/04-read-certificates-rolebinding.yaml) - This provides an atomic RoleBinding in `namespace/cert-sentinel` that binds the `ServiceAccount/some-service-account` to the `ClusterRole/certificate-reader` to read Certificates
- [05-read-collection-rolebinding.yaml](./example-rbac-clusterroles/05-read-collection-rolebinding.yaml) - This provides a combined RoleBinding in `namespace/cert-sentinel` that binds the `ServiceAccount/some-service-account` to the `ClusterRole/sentinel-reader` to read Namespaces, Secrets, ConfigMaps, and Certificates

These RoleBindings are scoped to specific namespaces, `cert-sentinel` as they are defined.  In order to allow the Certificate Sentinel Operator's ServiceAccounts to scan across the whole cluster a ClusterRoleBinding would be more ideal.

## RBAC ClusterRoleBindings

Information on ClusterRoleBindings: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding

The [examples-rbac-clusterrolebindings](./example-rbac-clusterrolebindings) directory houses a number of Kubernetes RoleBindings:

- [01-list-namespaces-clusterrolebinding.yaml](./example-rbac-clusterrolebindings/01-list-namespaces-clusterrolebinding.yaml) - This provides an atomic ClusterRoleBinding that binds the `ServiceAccount/some-service-account` to the `ClusterRole/namespace-reader` to read Namespaces across the cluster
- [02-read-secrets-clusterrolebinding.yaml](./example-rbac-clusterrolebindings/02-read-secrets-clusterrolebinding.yaml) - This provides an atomic ClusterRoleBinding that binds the `ServiceAccount/some-service-account` to the `ClusterRole/secret-reader` to read Secrets across the cluster
- [03-read-configmaps-clusterrolebinding.yaml](./example-rbac-clusterrolebindings/03-read-configmaps-clusterrolebinding.yaml) - This provides an atomic ClusterRoleBinding that binds the `ServiceAccount/some-service-account` to the `ClusterRole/configmap-reader` to read ConfigMaps across the cluster
- [04-read-certificates-clusterrolebinding.yaml](./example-rbac-clusterrolebindings/04-read-certificates-clusterrolebinding.yaml) - This provides an atomic ClusterRoleBinding that binds the `ServiceAccount/some-service-account` to the `ClusterRole/certificate-reader` to read Certificates across the cluster
- [05-read-collection-clusterrolebinding.yaml](./example-rbac-clusterrolebindings/05-read-collection-clusterrolebinding.yaml) - This provides a combined ClusterRoleBinding that binds the `ServiceAccount/some-service-account` to the `ClusterRole/sentinel-reader` to read Namespaces, Secrets, ConfigMaps, and Certificates across the cluster

# CRD Deployment, CertificateSentinel and KeystoreSentinel

For deployment examples on how to create the CustomResourceDefinitions provided, CertificateSentinel and KeystoreSentinel, see the [Documentation](../docs/).