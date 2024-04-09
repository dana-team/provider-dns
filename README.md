# Provider DNS

`provider-dns` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
DNS API.

## Getting Started

Install the provider on your cluster:

```bash
$ up ctp provider install quay.io/danateamorg/provider-dns:v0.1.0
```

Alternatively, you can use declarative installation:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-dns
spec:
  package: quay.io/danateamorg/provider-dns:v0.1.0
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
      "rfc": "<3645|2845>",
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

The following table summarizes the available resources:

| Name            | apiversion                               | namespaced | kind          |
|-----------------|------------------------------------------|------------|---------------|
| ptrs            | record.dns.crossplane.io/v1alpha1       | false      | PTR           |
| aaaarecordsets  | recordset.dns.crossplane.io/v1alpha1    | false      | AAAARecordSet |
| arecordsets     | recordset.dns.crossplane.io/v1alpha1    | false      | ARecordSet    |
| mxrecordsets    | recordset.dns.crossplane.io/v1alpha1    | false      | MXRecordSet   |
| nsrecordsets    | recordset.dns.crossplane.io/v1alpha1    | false      | NSRecordSet   |
| srvrecordsets   | recordset.dns.crossplane.io/v1alpha1    | false      | SRVRecordSet  |
| txtrecordsets   | recordset.dns.crossplane.io/v1alpha1    | false      | TXTRecordSet  |

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
  providerConfigRef:
    name: default
```

For details on how to configure the rest of the resources, use `kubectl explain` to see the available `spec` options, and advise with the the Terraform [provider-dns](https://registry.terraform.io/providers/hashicorp/dns/latest/docs) docs.

## Developing

Run code-generation pipeline:

```bash
$ go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```bash
$ make run
```

Build, push, and install:

```bash
$ make all
```

Build binary:

```bash
$ make build
```
