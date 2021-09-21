# SMTP Configuration

In order to have a CertificateSentinel/KeystoreSentinel communicate with an SMTP server, it must be provided some log in parameters - these are stored as a Secret and the name of that Secret is passed to the `.spec.alerts[*].config.smtp_auth_secret` definition.

```bash
## Plain SMTP Authentication Type
export SMTP_USERNAME="someUser"
export SMTP_PASSWORD="securePassword"

oc create secret generic my-smtp-secret-name --from-literal=username=${SMTP_USERNAME} --from-literal=password=${SMTP_PASSWORD}

# Login SMTP Authentication Type
export SMTP_IDENTITY="yourIdentity" # in addition to the SMTP_{USERNAME,PASSWORD} exported vars above
oc create secret generic my-smtp-secret-name --from-literal=username=${SMTP_USERNAME} --from-literal=password=${SMTP_PASSWORD} --from-literal=identity=${SMTP_IDENTITY}

## CRAM-MD5 SMTP Authentication Types
export SMTP_CRAM_MD5="challengeSecret" # in addition to the SMTP_USERNAME exported var above
oc create secret generic my-smtp-secret-name --from-literal=username=${SMTP_USERNAME} --from-literal=cram=${SMTP_CRAM_MD5}
```