#!/usr/bin/env bash

#GIT_HASH=$(git rev-parse HEAD)
GIT_HASH=v1.12.${RANDOM}
#DOCKER_DEV_ACCOUNT=localhost:5000/cilium
DOCKER_DEV_ACCOUNT=quay.io/sayboras

DOCKER_IMAGE_TAG=${GIT_HASH} DOCKER_DEV_ACCOUNT=${DOCKER_DEV_ACCOUNT} make \
  docker-operator-generic-image \
  docker-cilium-image

docker push "${DOCKER_DEV_ACCOUNT}/operator-generic:${GIT_HASH}"
docker push "${DOCKER_DEV_ACCOUNT}/cilium:${GIT_HASH}"

cilium uninstall

sleep 10

cilium install \
  --chart-directory ./install/kubernetes/cilium \
  --version v1.12.0 \
  --kube-proxy-replacement=strict \
  --helm-set securityContext.privileged=true \
  --helm-set debug.enabled=true \
  --helm-set debug.verbose=datapath \
  --helm-set prometheus.enabled=true \
  --helm-set ingressController.enabled=true \
  --helm-set hubble.enabled=true \
  --helm-set hubble.relay.enabled=true \
  --helm-set hubble.relay.prometheus.enabled=true \
  --helm-set hostFirewall.enabled=true \
  --agent-image "${DOCKER_DEV_ACCOUNT}/cilium:${GIT_HASH}" \
  --operator-image "${DOCKER_DEV_ACCOUNT}/operator-generic:${GIT_HASH}"

#kubectl set image daemonsets/cilium -n kube-system cilium-agent="${DOCKER_DEV_ACCOUNT}/cilium:${GIT_HASH}"
#kubectl set image deployment/cilium-operator -n kube-system cilium-operator="${DOCKER_DEV_ACCOUNT}/operator-generic:${GIT_HASH}"



