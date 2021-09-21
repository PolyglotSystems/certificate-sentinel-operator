# Certificate Sentinel Operator

[![Go Reference](https://pkg.go.dev/badge/github.com/PolyglotSystems/certificate-sentinel-operator.svg)](https://pkg.go.dev/github.com/PolyglotSystems/certificate-sentinel-operator) [![Go Report Card](https://goreportcard.com/badge/github.com/PolyglotSystems/certificate-sentinel-operator)](https://goreportcard.com/report/github.com/PolyglotSystems/certificate-sentinel-operator) [![License: GPL v3](https://img.shields.io/badge/License-Apache%20v2-blue.svg)](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/LICENSE)

#### Tested on OpenShift 4.8

> The Certificate Sentinel Operator allows for the scanning and reporting of SSL Certificates within a Kubernetes/OpenShift cluster.

This Operator provides two Custom Resource Definitions (CRDs):

- **CertificateSentinel** - This provides scanning of a cluster/namespace(es) for PEM-encoded x509 Certificates, generating an overall inventory list, list of expiring certificates, and produces STDOUT and SMTP reports.
- **KeystoreSentinel** - This provides scanning of a cluster/namespace(es) for x509 Certificates in Java Keystores, generating an overall inventory list, list of expiring certificates, and produces STDOUT and SMTP reports.

## Documentation

- [Quickstart](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/quickstart.md)
- [SMTP Configuration](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/smtp-configuration.md)
- [Examples - SSL Certificates for CertificateSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/examples/ssl_certificates/)
- [Examples - Keystore for KeystoreSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/examples/java_keystore/)
- [Full YAML Structure - CertificateSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/full_yaml_spec-CertificateSentinel.md)
- [Full YAML Structure - KeystoreSentinel](https://github.com/PolyglotSystems/certificate-sentinel-operator/tree/main/docs/full_yaml_spec-KeystoreSentinel.md)

## Deploying the Operator

PROD IDK yet, hold on, i'm doing it dev mode deployment bby

### Development & Testing Deployment

Requires Golang 1.16+ and the DevelopmentTools dnf group.

```bash
# plz be `oc login`'d already
# also also need @DevelopmentTools & golang installed
git clone https://github.com/PolyglotSystems/certificate-sentinel-operator
cd certificate-sentinel-operator/
make generate && make manifests && make install run
```