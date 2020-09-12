#!/usr/bin/env bash

kind delete cluster || true
kind create cluster --wait 3m --config=.github/kind-config.yaml

# Install cilium
echo "Installing cilium"
helm install cilium ./install/kubernetes/cilium \
   --wait \
   --namespace kube-system \
   --set global.debug.enabled=true \
   --set global.nodeinit.enabled=true \
   --set global.kubeProxyReplacement=partial \
   --set global.hostServices.enabled=false \
   --set global.externalIPs.enabled=true \
   --set global.nodePort.enabled=true \
   --set global.hostPort.enabled=true \
   --set config.bpfMasquerade=false \
   --set config.ipam=kubernetes \
   --set global.prometheus.enabled=true \
   --set global.operatorPrometheus.enabled=true \
   --set global.hubble.enabled=true \
   --set global.hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"

echo "Waiting for cilium"
kubectl wait -n kube-system --for=condition=Ready --all pod --timeout=5m
kubectl wait --for condition=Established crd/ciliumnetworkpolicies.cilium.io --timeout=5m

echo "Installing connectivity"
kubectl apply -f examples/kubernetes/connectivity-check/connectivity-check.yaml
