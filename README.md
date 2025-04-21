# kube-create-secret

`kubectl create secret` but with extras. Built on top of [vals](https://github.com/variantdev/vals), so all back-ends handled by `vals` can be used (AWS secret manager, Hashicorp Vault, Azure Key Vault, ...)

Create a secret from a `SecretTemplate` containing `vals` expression to resolve the secrets. Extras on top of what `vals` already gives out of the box:

- create a `Secret` of type `kubernetes.io/tls` directly from a pkcs12 keystore
- optionally output directly as a `SealedSecret`
- adds an annotation `kube-create-secret/source` on the generated `Secret` (or `SealedSecret` template) containing the original template, to keep track of where the secret values originate from. Easily recreate the secret based on this annotation.

```shell
$ kube-create-secret
Usage:
  kube-create-secret [command]

Examples:
  kube-create-secret create -f template.yaml
  kube-create-secret re-create -f secret.yaml
  kube-create-secret show -f secret.yaml
  kube-create-secret new


Available Commands:
  create      Create a secret from a SecretTemplate definition.
  re-create   Re-create a secret from a Secret that was previously created with kube-create-secret.
  show        Show the template for a Secret that was previously created with kube-create-secret.
  new         Print starter template.
```

- [Examples](#examples)
- [Tls secrets](#tls-secrets)
- [Sealed secrets](#sealed-secrets)
- [Re-creating secrets](#re-creating-secrets)

## Examples

With a template:

```shell
$ cat <<EOF | VAR1=foo kube-create-secret create -f -
---
apiVersion: kube-create-secret/v1
kind: SecretTemplate
metadata:
  name: my-secret-template-1
spec:
  metadata:
    name: created-secret-1
    namespace: kube-system
  data:
    POSTGRES_PASSWORD: ref+azurekeyvault://my-vault/postgres-password
    POSTGRES_USERNAME: postgres
    VAR1: ref+envsubst://\$VAR1
  stringData:
    API_TOKEN: ref+azurekeyvault://my-vault/api-token
  type: Opaque
EOF
```

Results in:

```yaml
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  annotations:
    kube-create-secret/source: |-
      {
        "apiVersion": "kube-create-secret/v1",
        "kind": "SecretTemplate",
        "metadata": {
          "name": "my-secret-template-1"
        },
        "spec": {
          "metadata": {
            "name": "created-secret-1",
            "namespace": "kube-system"
          },
          "data": {
            "POSTGRES_PASSWORD": "ref+azurekeyvault://my-vault/postgres-password",
            "POSTGRES_USERNAME": "postgres",
            "VAR1": "ref+envsubst://"
          },
          "stringData": {
            "API_TOKEN": "ref+azurekeyvault://my-vault/api-token"
          },
          "type": "Opaque"
        }
  name: created-secret-1
  namespace: kube-system
data:
  POSTGRES_PASSWORD: my-postgres-password
  POSTGRES_USERNAME: postgres
  VAR1: foo
stringData:
  API_TOKEN: my-secet-token
type: Opaque
```

Output format is by default the format of the input but can be set via the `-o` option.

Multiple yaml documents in one yaml file is supported. E.g.:

```shell
$ cat <<EOF | VAR1=foo VAR2=bar kube-create-secret create -f -
---
apiVersion: kube-create-secret/v1
kind: SecretTemplate
metadata:
  name: my-secret-template-1
spec:
  metadata:
    name: created-secret-1
  stringData:
    VAR: ref+envsubst://\$VAR1
  type: Opaque
---
apiVersion: kube-create-secret/v1
kind: SecretTemplate
metadata:
  name: my-secret-template-2
spec:
  metadata:
    name: created-secret-2
  stringData:
    VAR: ref+envsubst://\$VAR2
  type: Opaque
EOF
```

Gives:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  annotations:
    kube-create-secret/source: |-
      {
        "kind": "SecretTemplate",
        "apiVersion": "kube-create-secret/v1",
        "metadata": {
          "name": "my-secret-template-1",
          "creationTimestamp": null
        },
        "spec": {
          "kind": "Secret",
          "apiVersion": "v1",
          "metadata": {
            "name": "created-secret-1",
            "creationTimestamp": null
          },
          "type": "Opaque",
          "stringData": {
            "VAR": "ref+envsubst://$VAR1"
          }
        }
      }
  name: created-secret-1
stringData:
  VAR: foo
type: Opaque
---
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  annotations:
    kube-create-secret/source: |-
      {
        "kind": "SecretTemplate",
        "apiVersion": "kube-create-secret/v1",
        "metadata": {
          "name": "my-secret-template-2",
          "creationTimestamp": null
        },
        "spec": {
          "kind": "Secret",
          "apiVersion": "v1",
          "metadata": {
            "name": "created-secret-2",
            "creationTimestamp": null
          },
          "type": "Opaque",
          "stringData": {
            "VAR": "ref+envsubst://$VAR2"
          }
        }
      }
  name: created-secret-2
stringData:
  VAR: bar
type: Opaque
```

For multiple documents in json the kubernetes format for lists can be used:

```shell
$ cat <<EOF | VAR1=foo VAR2=bar kube-create-secret create -f -
{
  "apiVersion": "v1",
  "kind": "List",
  "items": [
    {
      "apiVersion": "kube-create-secret/v1",
      "kind": "SecretTemplate",
      "metadata": {
        "name": "my-secret-template-1"
      },
      "spec": {
        "metadata": {
          "name": "created-secret-1"
        },
        "stringData": {
          "VAR": "ref+envsubst://\$VAR1"
        },
        "type": "Opaque"
      }
    },
    {
      "apiVersion": "kube-create-secret/v1",
      "kind": "SecretTemplate",
      "metadata": {
        "name": "my-secret-template-2"
      },
      "spec": {
        "metadata": {
          "name": "created-secret-2"
        },
        "stringData": {
          "VAR": "ref+envsubst://\$VAR2"
        },
        "type": "Opaque"
      }
    }
  ]
}
EOF
```

Which gives:

```json
{
  "kind": "List",
  "apiVersion": "v1",
  "metadata": {},
  "items": [
    {
      "kind": "Secret",
      "apiVersion": "v1",
      "metadata": {
        "name": "created-secret-1",
        "creationTimestamp": null,
        "annotations": {
          "kube-create-secret/source": "{\"kind\":\"SecretTemplate\",\"apiVersion\":\"kube-create-secret/v1\",\"metadata\":{\"name\":\"my-secret-template-1\",\"creationTimestamp\":null},\"spec\":{\"kind\":\"Secret\",\"apiVersion\":\"v1\",\"metadata\":{\"name\":\"created-secret-1\",\"creationTimestamp\":null},\"type\":\"Opaque\",\"stringData\":{\"VAR\":\"ref+envsubst://$VAR1\"}}}"
        }
      },
      "stringData": {
        "VAR": "foo"
      },
      "type": "Opaque"
    },
    {
      "kind": "Secret",
      "apiVersion": "v1",
      "metadata": {
        "name": "created-secret-2",
        "creationTimestamp": null,
        "annotations": {
          "kube-create-secret/source": "{\"kind\":\"SecretTemplate\",\"apiVersion\":\"kube-create-secret/v1\",\"metadata\":{\"name\":\"my-secret-template-2\",\"creationTimestamp\":null},\"spec\":{\"kind\":\"Secret\",\"apiVersion\":\"v1\",\"metadata\":{\"name\":\"created-secret-2\",\"creationTimestamp\":null},\"type\":\"Opaque\",\"stringData\":{\"VAR\":\"ref+envsubst://$VAR2\"}}}"
        }
      },
      "stringData": {
        "VAR": "bar"
      },
      "type": "Opaque"
    }
  ]
}
```

The `data` or `stringData` of the secret can be populated with an entire object. E.g.:

```shell
$ cat <<EOF | kube-create-secret create -f -
---
apiVersion: kube-create-secret/v1
kind: SecretTemplate
metadata:
  name: my-expanded-secret-template-1
spec:
  metadata:
    name: created-expanded-secret-1
    namespace: kube-system
  data: ref+azurekeyvault://my-vault#/*
  stringData: ref+azurekeyvault://my-vault#/*
  type: Opaque
EOF
```

Results in:

```yaml
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  annotations:
    kube-create-secret/source: |-
      {
        "kind": "SecretTemplate",
        "apiVersion": "kube-create-secret/v1",
        "metadata": {
          "name": "my-expanded-secret-template-1",
          "creationTimestamp": null
        },
        "spec": {
          "kind": "Secret",
          "apiVersion": "v1",
          "metadata": {
            "name": "created-expanded-secret-1",
            "namespace": "kube-system",
            "creationTimestamp": null
          },
          "type": "Opaque",
          "data": "ref+azurekeyvault://my-vault#/*",
          "stringData": "ref+azurekeyvault://my-vault#/*"
        }
      }
  name: created-expanded-secret-1
  namespace: kube-system
data:
  api-token: bXktc2VjcmV0LXRva2Vu
  postgres-password: bXktcG9zdGdyZXNzLXBhc3N3b3Jk
stringData:
  api-token: my-secret-token
  postgres-password: my-postgress-password
type: Opaque
```

## Tls secrets

`kubectl create secret tls NAME --cert=path/to/cert/file --key=path/to/key/file` allows creating a `Secret` of type `kubernetes.io/tls`, but you need the key and certificate already in PEM format. Often the key and cert are in a keystore. For example in Azure Key Vault certificates can be stored and retrieved in pkcs12 format. To create a tls secret from such a file:

```shell
$ cat <<EOF | kube-create-secret create -f -
---
apiVersion: kube-create-secret/v1
kind: SecretTemplate
metadata:
  name: my-pkcs12-secret-template
spec:
  metadata:
    name: my-tls-secret
    namespace: default
  tls:
    pkcs12: ref+azurekeyvault://my-vault/my-vault-certificate
  type: kubernetes.io/tls
EOF
```

Results in:

```yaml
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  annotations:
    kube-create-secret/source: |-
      {
        "kind": "SecretTemplate",
        "apiVersion": "kube-create-secret/v1",
        "metadata": {
          "name": "my-pkcs12-secret-template",
          "creationTimestamp": null
        },
        "spec": {
          "kind": "Secret",
          "apiVersion": "v1",
          "metadata": {
            "name": "my-tls-secret",
            "namespace": "default",
            "creationTimestamp": null
          },
          "type": "kubernetes.io/tls",
          "tls": {
            "pkcs12": "ref+azurekeyvault://my-vault/my-vault-certificate"
          }
        }
      }
  name: my-tls-secret
  namespace: default
data:
  tls.crt: LS0tL...........SUNBVEUtLS0tLQo=
  tls.key: LS0tL...SBLRVktLS0tLQo=
type: kubernetes.io/tls
```

If the key store has a password, the password can be defined in the template, e.g.:

```yaml
tls:
  pkcs12: ref+azurekeyvault://my-vault/my-vault-certificate
  password: ref+envsubst://$PASSWORD
```

To specify the name of the key and crt in the secret:

```yaml
tls:
  pkcs12: ref+azurekeyvault://my-vault/my-vault-certificate
  name: foo
```

Which will result in:

```yaml
spec:
  ...
  data:
    foo.key: ...
    foo.crt: ...
```

Alternatively, a different name for the key and certificate can be given:

```yaml
tls:
  pkcs12: ref+azurekeyvault://my-vault/my-vault-certificate
  key:
    name: foo.key
  crt:
    name: bar.crt
```

Resulting in:

```yaml
spec:
  ...
  data:
    foo.key: ...
    bar.crt: ...
```

If the certficate contains a chain, it can in certain circumstances be useful to add a delimiter between the certificates in the chain:

```yaml
tls:
  pkcs12: ref+azurekeyvault://my-vault/my-vault-certificate
  crt:
    delimiter: ","
```

If the key and cert are known in PEM format then nothing special to do, e.g.:

```yaml
apiVersion: kube-create-secret/v1
kind: SecretTemplate
metadata:
  name: my-pkcs12-pems-secret-template
spec:
  metadata:
    name: my-tls-secret
    namespace: default
  data:
    tls.key: ref+file://my-key.key
    tls.crt: ref+file://my-certificate.crt
  type: kubernetes.io/tls
```

## Sealed secrets

If you're using [sealed secrets](https://github.com/bitnami-labs/sealed-secrets) the output of `kube-create-secret` can be piped to `kubeseal`, but for convenience `kube-create-secret` can directly output sealed secrets. E.g.:

```shell
$ cat <<EOF | VAR1=foo kube-create-secret create -f - -o sealed-secret
---
apiVersion: kube-create-secret/v1
kind: SecretTemplate
metadata:
  name: my-sealed-secret-template-1
spec:
  apiVersion: v1
  kind: Secret
  metadata:
    name: created-sealed-secret-1
    namespace: default
  data:
    VAR: ref+envsubst://\$VAR1
  type: Opaque
EOF
```

Gives:

```yaml
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  creationTimestamp: null
  name: created-sealed-secret-1
  namespace: default
spec:
  encryptedData:
    VAR: AgBZuBgjEr...koGVkfwhlQ=
  template:
    metadata:
      creationTimestamp: null
      annotations:
        kube-create-secret/source: |-
          {
            "kind": "SecretTemplate",
            "apiVersion": "kube-create-secret/v1",
            "metadata": {
              "name": "my-sealed-secret-template-1",
              "creationTimestamp": null
            },
            "spec": {
              "kind": "Secret",
              "apiVersion": "v1",
              "metadata": {
                "name": "created-sealed-secret-1",
                "creationTimestamp": null
              },
              "type": "Opaque",
              "data": {
                "VAR": "ref+envsubst://$VAR1"
              }
            }
          }
      name: created-sealed-secret-1
      namespace: default
    type: Opaque
```

The input parameters for `kubeseal` (e.g. to give the sealed secret controller name, the kubernetes context, ...) can be given via the input flags starting with `--kubeseal-`.

## Re-creating secrets

The template from which a secret is generated is saved as an annotation on the the secret. This annotation can be used to re-create the secret if the values in the secret have changed. For example:

```shell
$ kubectl get secret my-secret -o yaml | kube-create-secret re-create -f - | kubectl apply -f -
```

It retrieves the secret from the cluster, re-evaluates the template of it with the latest values for the `vals` references, and then re-applies it on the cluster.

If an extra secret value needs to be added, the template can first be retrieved back via `kube-create-secret show` after which the template can be updated and the secret re-generated via `kube-create-secret create`.

```shell
$ kubectl get secret my-secret -o yaml | kube-create-secret show -f -
```

The `re-create` command is the combination of the `show` command and the `create` command, this is the same as the `re-create` example from above:

```shell
$ kubectl get secret my-secret -o yaml | kube-create-secret show -f - | kube-create-secret create -f | kubectl apply -f -
```

When using GitOps with sealed secrets you can either store the secret template next to the generated sealed secret in git, or just the generated sealed secret and then use `re-create` when the sealed secret needs to be updated. For example, for a sealed secret generated by `kube-create-secret` stored in git, you can do `kube-create-secret re-create -f my-sealed-secret.yml` to see what is in the sealed secret (assuming the `vals` reference values didn't change).
