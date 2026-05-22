<br>
# API Docs
This Document documents the types introduced by the DigiCert Issuer to be consumed by users.
> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
- [API Docs](#api-docs)
  - [Table of Contents](#table-of-contents)
  - [ClusterDigicertIssuer](#clusterdigicertissuer)
  - [ClusterDigicertIssuerList](#clusterdigicertissuerlist)
  - [DigicertIssuer](#digicertissuer)
  - [DigicertIssuerCondition](#digicertissuercondition)
  - [DigicertIssuerList](#digicertissuerlist)
  - [DigicertIssuerSpec](#digicertissuerspec)
  - [DigicertIssuerStatus](#digicertissuerstatus)
  - [DigicertProvisioner](#digicertprovisioner)
    - [Using `preferredChain` and `caCertID`](#using-preferredchain-and-cacertid)
  - [SecretKeySelector](#secretkeyselector)

## ClusterDigicertIssuer

ClusterDigicertIssuer is the Schema for the clusterdigicertissuers API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [DigicertIssuerSpec](#digicertissuerspec) | true |
| status |  | *[DigicertIssuerStatus](#digicertissuerstatus) | false |

[Back to TOC](#table-of-contents)

## ClusterDigicertIssuerList

ClusterDigicertIssuerList contains a list of ClusterDigicertIssuer

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][ClusterDigicertIssuer](#clusterdigicertissuer) | true |

[Back to TOC](#table-of-contents)

## DigicertIssuer

DigicertIssuer is the Schema for the digicertissuers API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [DigicertIssuerSpec](#digicertissuerspec) | true |
| status |  | *[DigicertIssuerStatus](#digicertissuerstatus) | false |

[Back to TOC](#table-of-contents)

## DigicertIssuerCondition

DigicertIssuerCondition  ...

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type | Type of the condition, currently ('Ready'). | ConditionType | true |
| status | Status of the condition, one of ('True', 'False', 'Unknown'). | ConditionStatus | true |
| lastTransitionTime | LastTransitionTime is the timestamp corresponding to the last status change of this condition. | *metav1.Time | false |
| reason | Reason is a brief machine readable explanation for the condition's last transition. | ConditionReason | false |
| message | Message is a human readable description of the details of the last transition, complementing reason. | string | false |

[Back to TOC](#table-of-contents)

## DigicertIssuerList

DigicertIssuerList contains a list of DigicertIssuer

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][DigicertIssuer](#digicertissuer) | true |

[Back to TOC](#table-of-contents)

## DigicertIssuerSpec

DigicertIssuerSpec defines the desired state of DigicertIssuer

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| url | Optional URL is the DigiCert cert-central API. | string | false |
| provisioner | Provisioner contains the DigiCert provisioner configuration. | [DigicertProvisioner](#digicertprovisioner) | true |

[Back to TOC](#table-of-contents)

## DigicertIssuerStatus

DigicertIssuerStatus defines the observed state of DigicertIssuer

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| conditions | Conditions is a list of DigicertIssuerConditions describing the current status. | [][DigicertIssuerCondition](#digicertissuercondition) | false |

[Back to TOC](#table-of-contents)

## DigicertProvisioner

DigicertProvisioner contains the DigiCert provisioner configuration.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| apiTokenReference | APITokenReference references a secret in the same namespace containing the DigiCert API token. | [SecretKeySelector](#secretkeyselector) | true |
| preferredChain | PreferredChain requests a preferred trust chain root common name. This is best-effort and falls back to the default chain when not available. | string | false |
| caCertID | CACertID is the ID of the CA if multiple CA certificates are configured in the (sub-)account. | string | false |
| organizationID | OrganizationID is the ID of the organization in Digicert. | *int | false |
| organizationName | OrganizationName is the name of the organization in Digicert. If specified takes precedence over OrganizationID. | string | false |
| organizationUnits | OrganizationUnits is the list of organizational units. | []string | false |
| validityDays | ValidityDays is the validity of the order and certificate in days. Overrides ValidityYears if set. | *int | false |
| validityYears | ValidityYears is the validity of the order and certificate in years. Defaults to 1 year if not set. Can be overridden by ValidityDays. | *int | false |
| disableRenewalNotifications | DisableRenewalNotifications disables email renewal notifications for expiring certificates. | *bool | false |
| paymentMethod | PaymentMethod is the configured payment method in the Digicert account. | string | false |
| skipApproval | SkipApproval skips the approval of the certificate. | *bool | false |
| orderType | OrderType is the certificate order type. | string | false |
| containerID | ContainerID is the ID of the division | *int | false |

### Using `preferredChain` and `caCertID`

DigiCert can issue certificates through multiple trust chains when cross-signed intermediates are available. This lets clients that trust **different** root CAs all validate the same certificate. By default the issuer uses whichever chain DigiCert returns first.

**`preferredChain`** lets you request a specific root CA by its common name. The issuer searches all chains returned by DigiCert and selects the one whose root matches the given name. If no match is found it falls back to the default chain silently, so this field handles all cases: single-chain, multi-chain, and cross-signed bundles without any special configuration on your side.

**`caCertID`** selects which intermediate CA DigiCert should use to *issue* the certificate. This is independent of `preferredChain`: `caCertID` controls the issuing CA, while `preferredChain` controls which root appears at the top of the returned chain. If not set, the account's default CA is used for issuance.

| Configuration | Behaviour |
|---|---|
| neither set | DigiCert default chain and default issuing CA |
| `preferredChain` only | Default issuing CA, chain anchored to the requested root (or default if not found) |
| `caCertID` only | Specific issuing CA, DigiCert default chain for that CA |
| both set | Specific issuing CA, chain anchored to the requested root |

**Example** — request a certificate issued by a specific ECC intermediate, anchored to an ECC root:

```yaml
apiVersion: certmanager.cloud.sap/v1beta1
kind: DigicertIssuer
metadata:
  name: my-issuer
spec:
  provisioner:
    # Select the ECC intermediate CA to issue the certificate.
    caCertID: "<your-ca-cert-id>"

    # Among all chains returned for that CA, pick the one rooted at this CN.
    # If this root is not present in the response, the default chain is used.
    preferredChain: "DigiCert TLS ECC P384 Root G5"
```

Without `preferredChain`, the first chain DigiCert returns is used. Without `caCertID`, the account's default issuing CA is used.

[Back to TOC](#table-of-contents)

## SecretKeySelector

SecretKeySelector references a secret in the same namespace containing sensitive configuration.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | The name of the secret. | string | true |
| key | The key in the secret. | string | true |

[Back to TOC](#table-of-contents)
