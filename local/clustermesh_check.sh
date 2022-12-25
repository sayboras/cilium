#!/usr/bin/env bash

echo 'Checking in cluster 1'
for i in $(seq 1 10); do kubectl exec -ti --context kind-kind-1 -n default deployment/x-wing -- curl rebel-base; done

echo 'Checking in cluster 2'
for i in $(seq 1 10); do kubectl exec -ti --context kind-kind-2 -n default deployment/x-wing -- curl rebel-base; done
