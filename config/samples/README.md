<!--
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company

SPDX-License-Identifier: Apache-2.0
-->

Samples
-------

# Prerequisites

1. The cert-manager controller is running and ready in the current cluster.
2. The digicert-issuer controller is running and ready in the current cluster

# Usage

Note: In this context *deploy* means `kubectl apply -f $file.yaml`.

1. Replace the `DIGICERT_API_TOKEN` in the [secret](token-secret.yaml) and deploy it.

2. Adjust the configuration and deploy the [DigiCert issuer](certmanager_v1beta1_digicertissuer.yaml) and wait until it becomes `ready`.  
   Validate the current status via `kubectl describe digicertissuer digicert-issuer`.
   
3. Deploy the [sample certificate](certificate.yaml).

4. The cert-manager will create the corresponding *CertificateRequest*, which can be seen using `kubectl get certificaterequest`.

5. This controller will ensure the certificate is being issued by the Digicert API and stored in the specified secret.
   Verify the `tls.crt` and `tls.key` are present in the secret: `kubectl get secret -o yaml somednsname-tld`
