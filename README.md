# Provider DNS

`provider-dns` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
DNS API.

## Getting Started

### Install with Helm

Helm chart docs are available on `charts/provider-dns`. 

```
$ helm upgrade --install provider-dns --namespace crossplane-system --create-namespace oci://ghcr.io/dana-team/helm-charts/provider-dns --version <release>
```

### Install Manually

#### Set Up

First, create a `ConfigMap` which includes the content of the relevant `krb5.conf` file. This `ConfigMap` then needs to be mounted to the provider pod.

```bash
$ kubectl create configmap krb5-config --from-file=krb5.conf=/etc/krb5.conf -n crossplane-system
```

The `krb5.conf` file should look something like this:

```ini
[libdefaults]
    default_realm = DANA-DEV.COM

[realms]
 DANA-DEV.COM = {
     kdc = dana-wdc-1.dana-dev.com
     admin_server = dana-wdc-1.dana-dev.com
 }

[domain_realm]
 .dana-dev.com = DANA-DEV.COM
 dana-dev.com = DANA-DEV.COM
```

#### Install the provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-dns
spec:
  package: ghcr.io/dana-team/provider-dns:<release>
  runtimeConfigRef:
    apiVersion: pkg.crossplane.io/v1beta1
    kind: DeploymentRuntimeConfig
    name: config
```

```yaml
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: config
spec:
  deploymentTemplate:
    spec:
      selector:
        matchLabels:
          pkg.crossplane.io/provider: provider-dns
      template:
        spec:
          containers:
          - args:
            - --debug
            name: package-runtime
            volumeMounts:
            - mountPath: /etc/krb5.conf
              name: krb5-config
              readOnly: true
              subPath: krb5.conf
          volumes:
          - configMap:
              name: krb5-config
            name: krb5-config
```

## Configuration

The provider supports both `RFC 2845` and `RFC 3645` authentication models, but was only tested with `RFC 3645`. Each authentication model has different required parameters, refer to the Terraform [provider-dns](https://registry.terraform.io/providers/hashicorp/dns/latest/docs) for more details.

To connect to the provider, create the following `secret`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: example-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "rfc": "3645",
      "server": "<DNS-SERVER-FQDN>",
      "realm": "<DOMAIN-NAME-IN-CAPS>,
      "username": "<DOMAIN-USER>",
      "password": "<DOMAIN-USER-PASSWORD>"
    }
```

For example:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: example-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "rfc": "3645",
      "server": "dana-wdc-1.dana-dev.com",
      "realm": "DANA-DEV.COM",
      "username": "dana",
      "password": "KLm&x7Cv%GT@k!"
    }
```

Then create the `ProviderConfig`:

```yaml
apiVersion: dns.dns.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: example-creds
      namespace: crossplane-system
      key: credentials
```

## Resources

To Install the CRDs manually, run:

```bash
$ make generate
$ kubectl apply -f package/crds
```

The following table summarizes the available resources:

| Name            | apiVersion                               | Namespaced | Kind          |
|-----------------|------------------------------------------|------------|---------------|
| `ptrs`           | `record.dns.crossplane.io/v1alpha1`       | false      | `PTRRecord`     |
| `cnamerecords`   | `record.dns.crossplane.io/v1alpha1`       | false      | `CNAMERecord`   |
| `aaaarecordsets` | `recordset.dns.crossplane.io/v1alpha1`    | false      | `AAAARecordSet` |
| `arecordsets`     | `recordset.dns.crossplane.io/v1alpha1`    | false      | `ARecordSet`    |
| `mxrecordsets`    | `recordset.dns.crossplane.io/v1alpha1`    | false      | `MXRecordSet`   |
| `nsrecordsets`    | `recordset.dns.crossplane.io/v1alpha1`    | false      | `NSRecordSet`   |
| `srvrecordsets`   | `recordset.dns.crossplane.io/v1alpha1`    | false      | `SRVRecordSet`  |
| `txtrecordsets`   | `recordset.dns.crossplane.io/v1alpha1`    | false      | `TXTRecordSet`  |

## Examples

### ARecordSet

```yaml
apiVersion: recordset.dns.crossplane.io/v1alpha1
kind: ARecordSet
metadata:
  name: crossplane-test
spec:
  forProvider:
    addresses:
      - 10.1.30.1
      - 10.1.30.2
      - 10.1.30.3
    ttl: 3600
    zone: crossplane.dana-dev.com.
    name: testy-test # actual name of the record
  providerConfigRef:
    name: default
```

In order to create a record in a subdomain, include the subdomain in the name:

```yaml
apiVersion: recordset.dns.crossplane.io/v1alpha1
kind: ARecordSet
metadata:
  name: crossplane-test-sub
spec:
  forProvider:
    addresses:
      - 10.1.30.1
      - 10.1.30.2
      - 10.1.30.3
    ttl: 3600
    zone: crossplane.dana-dev.com.
    name: testy-test.example-sub # record will be called testy-test in subdomain example-sub
  providerConfigRef:
    name: default
```

### CNAMERecord

```yaml
apiVersion: record.dns.crossplane.io/v1alpha1
kind: CNAMERecord
metadata:
  name: crossplane-test-caname
spec:
  forProvider:
    cname: testy-test.crossplane.dana-dev.com.
    ttl: 3600
    zone: crossplane.dana-dev.com.
    name: cname-testy-test
  providerConfigRef:
    name: default
```

For details on how to configure the rest of the resources, use `kubectl explain` to see the available `spec` options, and advise with the Terraform [provider-dns](https://registry.terraform.io/providers/hashicorp/dns/latest/docs) docs.
