#!/usr/bin/env bash
set -aeuo pipefail

echo "Running setup.sh"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "Waiting until provider is healthy..."
${KUBECTL} wait provider.pkg ${PROJECT_NAME} --for condition=Healthy --timeout 5m

echo "Waiting for all pods to come up..."
${KUBECTL} -n ${CROSSPLANE_NAMESPACE} wait --for=condition=Available deployment --all --timeout=5m

echo "Creating cloud credential secret..."
${KUBECTL} -n ${CROSSPLANE_NAMESPACE} create secret generic provider-secret --from-literal=credentials="${UPTEST_CLOUD_CREDENTIALS}" \
    --dry-run=client -o yaml | ${KUBECTL} apply -f -

echo "Creating a default provider config..."
cat <<EOF | ${KUBECTL} apply -f -
apiVersion: zitadel.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: provider-secret
      namespace: ${CROSSPLANE_NAMESPACE}
      key: credentials
EOF
