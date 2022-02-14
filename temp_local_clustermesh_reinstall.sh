#!/bin/bash

set -euo pipefail

CLUSTER_1_NAME="kind-1"
CLUSTER_2_NAME="kind-2"

CLUSTER_1_CONTEXT="kind-${CLUSTER_1_NAME}"
CLUSTER_2_CONTEXT="kind-${CLUSTER_2_NAME}"

function info() {
    echo "=> ${1}"
}

kubectl config use "${CLUSTER_1_CONTEXT}"
cilium uninstall
kubectl delete -f "examples/kubernetes/clustermesh/global-service-example/cluster1.yaml"

kubectl config use "${CLUSTER_2_CONTEXT}"
cilium uninstall
kubectl delete -f "examples/kubernetes/clustermesh/global-service-example/cluster1.yaml"


info "Build the docker images"
#GIT_HASH=$(git rev-parse HEAD)
GIT_HASH=$RANDOM

DOCKER_IMAGE_TAG=${GIT_HASH} DOCKER_DEV_ACCOUNT=localhost:5000/cilium make -j$(nproc) docker-operator-generic-image docker-cilium-image docker-clustermesh-apiserver-image
docker push "localhost:5000/cilium/operator-generic:${GIT_HASH}"
docker push "localhost:5000/cilium/cilium:${GIT_HASH}"
docker push "localhost:5000/cilium/clustermesh-apiserver:${GIT_HASH}"

info "Install cilium to cluster 1"
kubectl config use "${CLUSTER_1_CONTEXT}"
cilium install --cluster-name "${CLUSTER_1_NAME}" --cluster-id "1" \
  --agent-image "localhost:5000/cilium/cilium:${GIT_HASH}" \
  --operator-image "localhost:5000/cilium/operator-generic:${GIT_HASH}" \
  --ipam kubernetes \
  --debug-mode

info "Install cilium to cluster 2 with inherit CA"
kubectl config use "${CLUSTER_2_CONTEXT}"
cilium install --cluster-name "${CLUSTER_2_NAME}" --cluster-id "2" \
  --inherit-ca "${CLUSTER_1_CONTEXT}" \
  --agent-image "localhost:5000/cilium/cilium:${GIT_HASH}" \
  --operator-image "localhost:5000/cilium/operator-generic:${GIT_HASH}" \
  --ipam kubernetes \
  --debug-mode

info "Creating the cluster mesh..."
cilium clustermesh enable --context "${CLUSTER_1_CONTEXT}" --service-type NodePort \
  --apiserver-image "localhost:5000/cilium/clustermesh-apiserver:${GIT_HASH}" \
  --debug-mode

cilium clustermesh enable --context "${CLUSTER_2_CONTEXT}" --service-type NodePort \
  --apiserver-image "localhost:5000/cilium/clustermesh-apiserver:${GIT_HASH}" \
  --debug-mode

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
