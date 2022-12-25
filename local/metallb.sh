#!/usr/bin/env bash

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/master/manifests/namespace.yaml
kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/master/manifests/metallb.yaml
kubectl wait -n metallb-system --for=condition=Ready pods --all

NET=$(docker network inspect kind | jq '.[].IPAM.Config[0] | select(.Gateway) | .Subnet' -r | sed 's|/16||') && \
  kubectl -n metallb-system create configmap config --dry-run=client -oyaml --from-file=config=/dev/stdin \
    <<< "address-pools: [{name: default, protocol: layer2, addresses: ['${NET/%.0/.200}-${NET/%.0/.250}']}]" | kubectl apply -f-


### Install cert-manager
helm repo add jetstack https://charts.jetstack.io
helm upgrade --install cert-manager jetstack/cert-manager --version v1.7.1 --namespace cert-manager --set installCRDs=true --create-namespace

kubectl apply -f- <<EOF
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: self-signed
  namespace: default
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: ca
  namespace: default
spec:
  isCA: true
  privateKey:
    algorithm: ECDSA
    size: 256
  secretName: ca
  commonName: ca
  issuerRef:
    name: self-signed
    kind: Issuer
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ca-issuer
  namespace: default
spec:
  ca:
    secretName: ca
EOF