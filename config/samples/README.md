Samples
-------

# Prerequisites

1. The cert-manager controller is running and ready in the current cluster.
2. The digicert-issuer controller is running and ready in the current cluster

# Usage

Note: In this context *deploy* means `kubectl apply -f $file.yaml`.

1. Replace the `DIGICERT_API_TOKEN` in the [secret](token-secret.yaml) and deploy it.
2. Adjust the configuration and deploy the [DigiCert issuer](certmanager_v1beta1_digicertissuer.yaml) and wait until it becomes `ready`.
3. Deploy the [sample certificate](certificate.yaml).
4. The cert-manager will create the corresponding *CertificateRequest*
5. This controller will ensure the certificate is being issued by the Digicert API.
