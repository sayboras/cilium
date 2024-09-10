# -*- mode: Python -*-
source = load_dynamic('./Tiltfile.source')

helm_values = [
    # 'install/kubernetes/cilium/values/minikube.yaml',
    # 'install/kubernetes/cilium/values/kind.yaml',
    # 'install/kubernetes/cilium/values/aws/aws.yaml',
    'contrib/testing/kind-custom.yaml'
]

helm_sets = [
    'kubeProxyReplacement=true',
    'nodeinit.enabled=true',
    'ipam.mode=kubernetes',
    'image.repository=cilium',
    'operator.replicas=1',
    'operator.image.repository=operator',
    'hubble.relay.enabled=false',
    'hubble.relay.prometheus.enabled=false',
    'hubble.ui.enabled=false',
    'hubble.peerService.enabled=false',
    'ingressController.enabled=false',
    'ingressController.loadbalancerMode=dedicated',
    'envoy.enabled=false',
    'gatewayAPI.enabled=true',
    'debug.enabled=true',
    'extraArgs={--envoy-config-retry-interval=5s}'
    # 'debug.verbose=envoy',
]

version = 'main'
cilium_dockerfile = 'images/cilium/Dockerfile.debug'
operator_dockerfile = 'images/operator/Dockerfile.debug'
operator_name = 'operator-generic-ci'

if version != 'main':
    cilium_dockerfile = 'images/cilium/Dockerfile%s.debug' % version
    operator_dockerfile = 'images/operator/Dockerfile%s.debug' % version
    operator_name = 'operator-generic'

# Local resources
local_resource('cilium-compile', 'make -j $(nproc) build && make -C hubble', deps=source['cilium_source'],
               ignore=source['ignore_files'],
               labels=['build'])

local_resource('cilium-operator-compile', 'make -C operator -j $(nproc) cilium-operator-generic',
               deps=source['operator_source'],
               ignore=source['ignore_files'],
               labels=['build'])

# local_resource('hubble-relay-compile', 'make -j $(nproc) hubble-relay', deps=source['hubble_relay_source'],
#                ignore=source['ignore_files'],
#                labels=['build'])

# Generate helm template
k8s_yaml(helm('install/kubernetes/cilium', name='cilium', namespace='kube-system', values=helm_values, set=helm_sets))

# Docker build
docker_build('cilium', '.', dockerfile=cilium_dockerfile, only=source['cilium_source'])

docker_build(operator_name, '.', dockerfile=operator_dockerfile,
             build_args={
                 "OPERATOR_VARIANT": 'operator-generic',
             },
             entrypoint=['/usr/bin/cilium-operator-generic'],
             only=source['operator_source'],
             )

# Deployment manifest
k8s_resource(workload='cilium', resource_deps=['cilium-compile'], labels=['deployment'])
k8s_resource(workload='cilium-operator', resource_deps=['cilium-operator-compile'], labels=['deployment'])
# k8s_resource(workload='hubble-ui', port_forwards=8081)
