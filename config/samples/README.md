Samples
-------

# Prerequisites

1. The cert-manager controller is running and ready in the current cluster.
2. The digicert-issuer controller is running and ready in the current cluster

# Usage

1. Replace the `DIGICERT_API_TOKEN` in the [secret](token-secret.yaml).
2. Deploy the [DigiCert issuer](certmanager_v1beta1_digicertissuer.yaml) and wait until it becomes `ready`.
3. Deploy the [sample certificate](certificate.yaml).
