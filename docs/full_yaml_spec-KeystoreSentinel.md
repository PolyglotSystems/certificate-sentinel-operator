# Full KeystoreSentinel YAML Spec Example

```yaml
apiVersion: config.polyglot.systems/v1
kind: KeystoreSentinel
metadata:
  name: keystoresentinel-sample
  namespace: cert-sentinel
spec:
  logLevel: 2 # [optional] logLevel is the verbosity of the logger, 1-4 with 4 being the most verbose, defaults to 2
  scanningInterval: 60 # [optional] scanningInterval is the number of seconds the Operator will scan the cluster - default is 60
  alert: # alerts is a list of alerting endpoints associated with these below targets
    # Log report to Stdout once a day, useful for Elastic/Splunk/etc environments
    name: secrets-mailer # must be a unique dns/k8s compliant name
    type: mailer # type can be: `logger` or `smtp`
    config: # optional on `logger` types, required for `smtp`
      reportInterval: daily # [optional] reportInterval can be `daily`, `weekly`, `monthly`, or `debug`, defaults to `daily`
      smtp_destination_addresses: # where is the emailed report being sent to, a list of emails
        - "infosec@example.com"
        - "certificates@example.com"
      smtp_sender_addresses: "ocp-certificate-sentinel+cluster-name@example.com" # what address is it being sent from
      smtp_sender_hostname: "cluster-name.example.com" # client hostname of the sender
      smtp_endpoint: "smtp.example.com:25" # SMTP endpoint, hostname:port format
      smtp_auth_secret: my-smtp-secret-name # name of the Secret containing the SMTP log in credentials
      smtp_auth_type: plain # SMTP authentication type, can be `plain`, `login`, or `cram-md5`
      smtp_use_ssl: false # [optional] Enable or disable SMTP TLS - defaults to `true`
      smtp_use_starttls: false # [optional] Enable or disable SMTP TLS - defaults to `true`
  target: # target is a Kubernetes object being targeted and scanned for Java Keystore data
    # Target Secrets/v1, looking for Keystores with certificates that have expirations coming in 30, 60, 90, 9000, and 9001 days across all namespaces with a specific serviceaccount
    apiVersion: v1 # Corresponds to the apiVersion of the object being targeted - likely just v1 for Secrets & ConfigMaps
    daysOut: # [optional] Expiration thresholds for 30, 60, 90, 9000, and 9001 days out - 9000/9001 are for testing.  Defaults to 30, 60, and 90
      - 30
      - 60
      - 90
      - 9001
      - 9000
    kind: Secret # Corresponds to the kind of the object being targeted - Secret or ConfigMap
    name: all-secrets # must be a unique dns/k8s compliant name
    namespaces: # list of namespaces to watch for Keystores in Secrets - can be a single wildcard or a list of specific namespaces
      - '*'
    serviceAccount: some-service-account # the ServiceAccount in tis namespace to use against the K8s/OCP API
    # [optional] targetLabels let you filter to targets such as Secrets and ConfigMaps that match a label filter
    targetLabels:
      - key: polyglot.systems/asset
        value:
          - certificate
    # [optional] namespaceLabels let you filter namespaces that match a label filter - will stack against the list of .spec.target.namespaces
    namespaceLabels:
      - key: polyglot.systems/certificate-sentinel-namespace
        value:
          - 'true'
    # keystorePassword is where the Keystore password should be sourced from - if the keystorePassword is inaccessible then the namespace will not be scanned for Keystore objects
    keystorePassword:
      plaintext: changeit
      type: plaintext

      #type: labels
      #labelRef:
      #  key: keystore-pass
      #  labelSelectors:
      #    - key: polyglot.systems/asset
      #      value:
      #      - keystore-password
      
      #type: secret
      #secretRef:
      #  key: keystore-pass
      #  name: keystore-secret-password
status: # .status is not user-defined, it will be updated at the end of a full scan/operator reconciliation and will list any Keystore Certificates found, the ones expiring within our designated daysOut thresholds, and when the last reports were sent for each alert
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