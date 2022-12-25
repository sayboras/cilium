#!/usr/bin/env bash

gcloud compute instances create cilium-cilium-tam-vm \
    --labels "usage=cilium-cilium,owner=19042" \
    --machine-type e2-custom-2-4096 \
    --boot-disk-type pd-standard \
    --boot-disk-size 10GB \
    --preemptible \
    --image-project ubuntu-os-cloud \
    --image-family ubuntu-2004-lts \
    --metadata hostname=cilium-cilium-tam-vm \
    --metadata-from-file startup-script=.github/gcp-vm-startup.sh


gcloud container clusters create cilium-cilium-tam-vm \
  --labels "usage=cilium-cilium,owner=19042" \
  --node-taints node.cilium.io/agent-not-ready=true:NoExecute \
  --image-type COS_CONTAINERD \
  --num-nodes 2 \
  --machine-type e2-custom-2-4096 \
  --disk-type pd-standard \
  --disk-size 10GB \
  --preemptible

cilium install --cluster-name=cilium-cilium-tam-vm \
  --agent-image=quay.io/cilium/cilium-ci \
  --operator-image=quay.io/cilium/operator-generic-ci \
  --version=7c99e4d8799493af7151e15f679cba8862696111 \
  --wait=false \
  --rollback=false \
  --config monitor-aggregation=none \
  --config tunnel=vxlan \
  --kube-proxy-replacement=strict

cilium clustermesh enable \
  --apiserver-image=quay.io/cilium/clustermesh-apiserver-ci \
  --apiserver-version=7c99e4d8799493af7151e15f679cba8862696111

cilium clustermesh vm create cilium-cilium-tam-vm -n default \
  --ipv4-alloc-cidr 10.192.1.0/30

cilium clustermesh vm install install-external-workload.sh --config debug
gcloud compute scp install-external-workload.sh cilium-cilium-tam-vm:~/
gcloud compute ssh cilium-cilium-tam-vm \
    --command "~/install-external-workload.sh"