DigiCert Issuer
---------------

[External Issuer](https://cert-manager.io/docs/configuration/external) extending the [cert-manager](https://cert-manager.io) with the [DigiCert cert-central API](https://dev.digicert.com/services-api/orders/).

# Prerequisites

The cert-manager and its `cert-manager.io/v1` CRDs needs to be installed in the selected cluster.

# Installation & Configuration

Use the [Kustomize resources](config) or run `make deploy` to install the DigiCert issuer in the current cluster.

# Documentation & Examples

For additional information see the [API documentation](docs/apidocs/api.md) and the provided [example](config/samples).

# Development

For development, it may be convenient to use `make deploy-local` to install all resource except the operator to the current cluster and run the operator locally via `make run`.

## Release management

Adjust the [version](VERSION) and run `make release` to build and publish a new version of the digicert-issuer.
