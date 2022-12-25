#!/usr/bin/env bash

cilium install --version -ci:"$(git rev-parse HEAD)" --config enable-envoy-config=true,enable-ingress-controller=true --kube-proxy-replacement=probe

# HTTP ingress
kubectl apply -n default -f https://raw.githubusercontent.com/istio/istio/release-1.11/samples/bookinfo/platform/kube/bookinfo.yaml
kubectl apply -n default -f https://raw.githubusercontent.com/cilium/cilium-service-mesh-beta/main/kubernetes-ingress/basic-ingress.yaml

lb=$(kg ingress basic-ingress -o json | jq '.status.loadBalancer.ingress[0].ip' | jq -r .)
curl http://"$lb"/details/1

# GRPC ingress
kubectl apply -n default -f https://raw.githubusercontent.com/GoogleCloudPlatform/microservices-demo/master/release/kubernetes-manifests.yaml
kubectl apply -n default -f https://raw.githubusercontent.com/cilium/cilium-service-mesh-beta/main/kubernetes-ingress/grpc-ingress.yaml

grpc=$(kg ingress grpc-ingress -o json | jq '.status.loadBalancer.ingress[0].ip' | jq -r .)
grpcurl -plaintext -proto ./local/demo.proto "$grpc":80 hipstershop.CurrencyService/GetSupportedCurrencies

# TLS ingress

kubectl apply -f https://raw.githubusercontent.com/cilium/cilium-service-mesh-beta/main/kubernetes-ingress/tls-ingress.yaml
kubectl create secret tls demo-cert --key=local/_.cilium.rocks/key.pem --cert=local/_.cilium.rocks/cert.pem

sudo perl -ni -e 'print if !/\.cilium\.rocks$/d' /etc/hosts; sudo tee -a /etc/hosts \
  <<<"$(kubectl get svc/cilium-ingress-tls-ingress -o=jsonpath='{.status.loadBalancer.ingress[0].ip}') bookinfo.cilium.rocks hipstershop.cilium.rocks"

kubectl wait --for=condition=Ready pods -n default --all

tlslb=$(kg ingress tls-ingress -o json | jq '.status.loadBalancer.ingress[0].hostname' | jq -r .)