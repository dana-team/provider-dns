# provider-dns

![Version: 0.0.0](https://img.shields.io/badge/Version-0.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

A Helm chart for Crossplane provider-dns.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| deploymentRuntimeConfig | object | `{"container":{"args":["--debug"],"name":"package-runtime","volumeMounts":{"mountPath":"/etc/krb5.conf","name":"krb5-config","readOnly":true,"subPath":"krb5.conf"}},"dnsConfig":{"nameservers":null},"dnsPolicy":"","name":"dns-config","volumes":{"configMap":{"name":"krb5-config"},"name":"krb5-config"}}` | Configuration to be added to the provider deployment via the DeploymentRuntimeConfig resource |
| deploymentRuntimeConfig.dnsPolicy | string | `""` | Optional DNS fields for local development |
| fullnameOverride | string | `""` |  |
| image.repository | string | `"ghcr.io/dana-team/provider-dns"` | The repository of the provider container image. |
| image.tag | string | `""` | The tag of the manager container image. |
| krb5Config.name | string | `"krb5-config"` | Name of the configMap which contains Kerberos configuration |
| nameOverride | string | `""` |  |
| provider.name | string | `"provider-dns"` | Name of the provider |
| provider.runtimeConfigRef.name | string | `"dns-config"` | Name of the DeploymentRuntimeConfig object to use |
| providerConfig | object | `{"credentials":{"secretRef":{"key":"credentials","name":"dns-creds","namespace":"crossplane-system"},"source":"Secret"},"name":"dns-default"}` | Provider authentication configuration |
| realm.kdc | string | `"dana-wdc-1.dana-dev.com"` | Name of the Kerberos Key Distribution Center server |
| realm.name | string | `"DANA-DEV.COM"` | Name of the Kerberos Realm |
| secret | object | `{"name":"dns-creds","password":"passw0rd","rfc":3645,"type":"Opaque","username":"dana"}` | Secret values for the provider authentication. |
| secret.name | string | `"dns-creds"` | Name of the secret. |
| secret.password | string | `"passw0rd"` | Password to connect to authenticate with. |
| secret.rfc | int | `3645` | RFC authentication (3645 for GSS-TSIG; 2845 for secret key based transaction). |
| secret.type | string | `"Opaque"` | Type of the secret. |
| secret.username | string | `"dana"` | Username to connect to authenticate with. |

