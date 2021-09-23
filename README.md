# Certificate Sentinel Operator

[![Go Reference](https://pkg.go.dev/badge/github.com/PolyglotSystems/certificate-sentinel-operator.svg)](https://pkg.go.dev/github.com/PolyglotSystems/certificate-sentinel-operator) [![Go Report Card](https://goreportcard.com/badge/github.com/PolyglotSystems/certificate-sentinel-operator)](https://goreportcard.com/report/github.com/PolyglotSystems/certificate-sentinel-operator) [![License: GPL v3](https://img.shields.io/badge/License-Apache%20v2-blue.svg)](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/LICENSE)

#### Tested on OpenShift 4.8

> The Certificate Sentinel Operator allows for the scanning and reporting of SSL Certificates within a Kubernetes/OpenShift cluster.

This Operator provides two Custom Resource Definitions (CRDs):

- **CertificateSentinel** - This provides scanning of a cluster/namespace(es) for PEM-encoded x509 Certificates, generating an overall inventory list, list of expiring certificates, and produces STDOUT and SMTP reports.
- **KeystoreSentinel** - This provides scanning of a cluster/namespace(es) for x509 Certificates in Java Keystores, generating an overall inventory list, list of expiring certificates, and produces STDOUT and SMTP reports.

## Documentation

- [Quickstart - CertificateSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/quickstart-certificatesentinel.md)
- [Quickstart - KeystoreSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/quickstart-keystoresentinel.md)
- [SMTP Configuration](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/smtp-configuration.md)
- [Examples - SSL Certificates for CertificateSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/examples/ssl_certificates/)
- [Examples - Keystore for KeystoreSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/examples/java_keystore/)
- [Full YAML Structure - CertificateSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/full_yaml_spec-CertificateSentinel.md)
- [Full YAML Structure - KeystoreSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/full_yaml_spec-KeystoreSentinel.md)

## Deploying the Operator

### Deploy to Red Hat OpenShift

If you already have an OpenShift cluster or Kubernetes cluster with OLM deployed, you can simply apply the Operator CatalogSource and then create the needed Subscription from there.

```bash
# Create the Operator CatalogSource in the openshift-marketplace Namespace/Project
oc apply -f https://raw.githubusercontent.com/PolyglotSystems/certificate-sentinel-operator/main/deploy/01-catalogsource.yaml -n openshift-marketplace

# Create a Subscription for the Operator (this installs the Operator)
oc apply -f https://raw.githubusercontent.com/PolyglotSystems/certificate-sentinel-operator/main/deploy/02-subscription.yaml -n openshift-operators
```

From here you should just need to create the needed RBAC and deploy the CRDs!

### Development & Testing Deployment

Requires Golang 1.16+ and the DevelopmentTools dnf group.

```bash
# plz be `oc login`'d already
# also also need @DevelopmentTools, Podman, & Golang installed

# Clone
git clone https://github.com/PolyglotSystems/certificate-sentinel-operator
cd certificate-sentinel-operator/

# Update the VERSION variable
vi Makefile

# Test & Build Test
make generate && make manifests && make install run
```

At this point, you have the Operator built and interacting with the current context cluster - in order to release a new version, the following workflow would apply:

```bash
# Create the Operator Container
make podman-build IMG="quay.io/username/repo:vX.Y.Z"
make podman-push IMG="quay.io/username/repo:vX.Y.Z"

# Create the Operator Bundle - this is basically just meta data
make bundle
make bundle-build BUNDLE_IMG="quay.io/username/repo-bundle:vX.Y.Z"
make bundle-push BUNDLE_IMG="quay.io/username/repo-bundle:vX.Y.Z"

# For a new addition to the Operator Catalog
make catalog-build CATALOG_IMG="quay.io/username/operator-catalog:vX.Y.Z"
make catalog-push CATALOG_IMG="quay.io/username/operator-catalog:vX.Y.Z"

# Git workflow for Github Actions
git add .
git commit -m "new version vX.Y.Z"
git tag vX.Y.Z HEAD
git push origin vX.Y.Z
```

So long as you have the `REGISTRY_USERNAME`, `REGISTRY_TOKEN`, and `GHUB_TOKEN` set up as GitHub Action Secrets then the automation workflow should kick off.

## Sample Build Script & Release

It's pretty easy to figure out the release process for the repos in use:

```yaml

make bundle && \
sudo make podman-build && \
sudo make podman-push && \
sudo make bundle-build && \
sudo make bundle-push && \
sudo make catalog-build && \
sudo make catalog-push
```