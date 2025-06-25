<!--
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company

SPDX-License-Identifier: Apache-2.0
-->

DigiCert Issuer
---------------

[![REUSE status](https://api.reuse.software/badge/github.com/sapcc/digicert-issuer)](https://api.reuse.software/info/github.com/sapcc/digicert-issuer)

[External Issuer](https://cert-manager.io/docs/configuration/external) extending the [cert-manager](https://cert-manager.io) with the [DigiCert cert-central API](https://dev.digicert.com/services-api/orders/).

# Prerequisites

The cert-manager and its `cert-manager.io/v1` CRDs needs to be installed in the selected cluster.

# Installation & Configuration

The container image can be found here: [ghcr.io/sapcc/digicert-issuer](https://github.com/sapcc/digicert-issuer/pkgs/container/digicert-issuer).

1) Using [Helm](https://helm.sh)  
   The [DigiCert Issuer Helm chart](https://github.com/sapcc/helm-charts/tree/master/system/digicert-issuer) can be used for **production** environments.  
   Additional documentation on configuration options is provided within the chart.

2) Using [Kustomize](https://kustomize.io)  
   Use the [Kustomize resources](config) or run `make deploy` to install the DigiCert issuer in the current cluster.

# Documentation & Examples

For additional information see the [API documentation](docs/apidocs/api.md) and the provided [example](config/samples).

# Development

For development, it may be convenient to use `make deploy-local` to install all resource except the operator to the current cluster and run the operator locally via `make run`.

## Release management

Adjust the [version](VERSION) and run `make release` to build and publish a new version of the digicert-issuer.
