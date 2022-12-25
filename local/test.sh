#!/bin/bash

set -euo pipefail

CLUSTER_1_NAME="tam-1"
CLUSTER_2_NAME="tam-2"

CLUSTER_1_CONTEXT="${CLUSTER_1_NAME}"
CLUSTER_2_CONTEXT="${CLUSTER_2_NAME}"

function info() {
    echo "=> ${1}"
}

kubectl config use "${CLUSTER_1_CONTEXT}"
cilium uninstall
kubectl delete --ignore-not-found -f "examples/kubernetes/clustermesh/global-service-example/cluster1.yaml"

kubectl config use "${CLUSTER_2_CONTEXT}"
cilium uninstall
kubectl delete --ignore-not-found -f "examples/kubernetes/clustermesh/global-service-example/cluster2.yaml"


info "Install cilium to cluster 1"
kubectl config use "${CLUSTER_1_CONTEXT}"
cilium install --cluster-name "${CLUSTER_1_NAME}" --cluster-id "1" \
  --ipam kubernetes \
  --config debug=true

info "Install cilium to cluster 2 with inherit CA"
kubectl config use "${CLUSTER_2_CONTEXT}"
cilium install --cluster-name "${CLUSTER_2_NAME}" --cluster-id "2" \
  --inherit-ca "${CLUSTER_1_CONTEXT}" \
  --ipam kubernetes \
  --config debug=true

info "Creating the cluster mesh..."
cilium clustermesh enable --context "${CLUSTER_1_CONTEXT}" --service-type NodePort

cilium clustermesh enable --context "${CLUSTER_2_CONTEXT}" --service-type NodePort

cilium clustermesh status --context "${CLUSTER_1_CONTEXT}" --wait
cilium clustermesh status --context "${CLUSTER_2_CONTEXT}" --wait
cilium clustermesh connect --context "${CLUSTER_1_CONTEXT}" --destination-context "${CLUSTER_2_CONTEXT}"
cilium clustermesh status --context "${CLUSTER_1_CONTEXT}" --wait
cilium clustermesh status --context "${CLUSTER_2_CONTEXT}" --wait

info "Deploying the test application..."
kubectl config use "${CLUSTER_1_CONTEXT}"
kubectl apply -f "examples/kubernetes/clustermesh/global-service-example/cluster1.yaml"
kubectl config use "${CLUSTER_2_CONTEXT}"
kubectl apply -f "examples/kubernetes/clustermesh/global-service-example/cluster2.yaml"

echo 'Checking in cluster 1'
for i in $(seq 1 10); do kubectl exec -ti --context ${CLUSTER_1_CONTEXT} deployment/x-wing -- curl rebel-base; done

echo 'Checking in cluster 2'
for i in $(seq 1 10); do kubectl exec -ti --context ${CLUSTER_2_CONTEXT} deployment/x-wing -- curl rebel-base; done
