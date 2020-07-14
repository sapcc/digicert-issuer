<br>
# API Docs
This Document documents the types introduced by the DigiCert Issuer to be consumed by users.
> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [DigicertIssuer](#digicertissuer)
* [DigicertIssuerCondition](#digicertissuercondition)
* [DigicertIssuerList](#digicertissuerlist)
* [DigicertIssuerSpec](#digicertissuerspec)
* [DigicertIssuerStatus](#digicertissuerstatus)

## DigicertIssuer

DigicertIssuer is the Schema for the digicertissuers API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [DigicertIssuerSpec](#digicertissuerspec) | false |
| status |  | [DigicertIssuerStatus](#digicertissuerstatus) | false |

[Back to TOC](#table-of-contents)

## DigicertIssuerCondition

DigicertIssuerCondition  ...

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type | Type of the condition, currently ('Ready'). | ConditionType | true |
| status | Status of the condition, one of ('True', 'False', 'Unknown'). | ConditionStatus | true |
| lastTransitionTime | LastTransitionTime is the timestamp corresponding to the last status change of this condition. | *metav1.Time | false |
| reason | Reason is a brief machine readable explanation for the condition's last transition. | string | false |
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
| url | URL is the DigiCert cert-central URL containing the token. | string | true |
| provisioner | Provisioner ... | [DigicertProvisioner](#digicertprovisioner) | true |

[Back to TOC](#table-of-contents)

## DigicertIssuerStatus

DigicertIssuerStatus defines the observed state of DigicertIssuer

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| conditions | Conditions ... | [][DigicertIssuerCondition](#digicertissuercondition) | false |

[Back to TOC](#table-of-contents)
