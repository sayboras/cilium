#!/usr/bin/env bash

CLUSTER_NAME=tammach-ipam

eksctl create cluster --name "${CLUSTER_NAME}" \
  --nodes 2 \
  --instance-types "t3.medium,t3a.medium" \
  --node-volume-type gp3 \
  --node-volume-size 10 \
  --managed --spot
