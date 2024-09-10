#!/usr/bin/env bash

kubectl label nodes kind-worker topology.kubernetes.io/zone=az-a
kubectl label nodes kind-worker2 topology.kubernetes.io/zone=az-a
kubectl label nodes kind-worker3 topology.kubernetes.io/zone=az-b
kubectl label nodes kind-worker4 topology.kubernetes.io/zone=az-b
kubectl label nodes kind-worker5 topology.kubernetes.io/zone=az-c
kubectl label nodes kind-worker6 topology.kubernetes.io/zone=az-c

cilium install --version 1.16.1 \
  --namespace kube-system \
  --set routingMode=native \
  --set bpf.masquerade=true \
  --set kubeProxyReplacement=true \
  --set ipam.mode=kubernetes \
  --set autoDirectNodeRoutes=true \
  --set loadBalancer.serviceTopology=true \
  --set ipv4NativeRoutingCIDR="10.0.0.0/8"
