cilium uninstall

GIT_HASH=$(git rev-parse HEAD)
DOCKER_IMAGE_TAG=${GIT_HASH} DOCKER_DEV_ACCOUNT=localhost:5000/cilium make docker-operator-generic-image docker-cilium-image
docker push "localhost:5000/cilium/operator-generic:${GIT_HASH}"
docker push "localhost:5000/cilium/cilium:${GIT_HASH}"

cilium install \
  --agent-image "localhost:5000/cilium/cilium:${GIT_HASH}" \
  --operator-image "localhost:5000/cilium/operator-generic:${GIT_HASH}" \
  --encryption ipsec
