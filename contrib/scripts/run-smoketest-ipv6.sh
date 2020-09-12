#!/usr/bin/env bash

kind delete cluster || true
kind create cluster --wait 2m --config=.github/kind-config-ipv6.yaml

# Install cilium
echo "Installing cilium"

helm upgrade -i cilium ./install/kubernetes/cilium \
  --namespace kube-system \
  --set global.debug.enabled=true \
  --set global.nodeinit.enabled=true \
  --set global.kubeProxyReplacement=strict \
  --set config.ipam=kubernetes \
  --set global.ipv6.enabled=true \
  --set global.ipv4.enabled=false \
  --set global.tunnel=disabled \
  --set global.autoDirectNodeRoutes=true \
  --set global.prometheus.enabled=true \
  --set global.operatorPrometheus.enabled=true \
  --set global.hubble.enabled=true \
  --set global.hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"

echo "Waiting for cilium"
kubectl wait -n kube-system --for=condition=Ready --all pod --timeout=3m
kubectl wait --for condition=Established crd/ciliumnetworkpolicies.cilium.io --timeout=3m

echo "Installing connectivity"
kubectl apply -f examples/kubernetes/connectivity-check/connectivity-check-internal.yaml
