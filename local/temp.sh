helm upgrade cilium cilium/cilium --version 1.11.3 --install \
    --namespace kube-system \
    --set ipam.mode=kubernetes \
    --set prometheus.enabled=true \
    --set operator.prometheus.enabled=true \
    --set hubble.enabled=true \
    --set hubble.ui.enabled=true \
    --set hubble.relay.enabled=true \
    --set hubble.metrics.enabled="{dns,drop,tcp:destinationContext=pod|ip;sourceContext=pod|ip,flow:destinationContext=pod|ip;sourceContext=pod|ip,port-distribution,icmp,http:destinationContext=pod|ip;sourceContext=pod|ip}" \
    --set hubble.metrics.port=9096

helm upgrade -i cilium install/kubernetes/cilium \
  --set debug.enabled=true \
  --set debug.verbose=flow \
  --set prometheus.enabled=true \
  --set ingressController.enabled=true \
  --set hubble.enabled=true \
  --set hubble.relay.enabled=true \
  --set hubble.relay.prometheus.enabled=true
