export NAME="$(whoami)-$RANDOM"
export AZURE_RESOURCE_GROUP="${NAME}-group"
az group create --name "${AZURE_RESOURCE_GROUP}" -l australiacentral

# Create AKS cluster
az aks create \
  --resource-group "${AZURE_RESOURCE_GROUP}" \
  --name "${NAME}" \
  --network-plugin azure \
  --node-count 1 \
  --generate-ssh-keys

# Get name of initial system node pool
nodepool_to_delete=$(az aks nodepool list \
  --resource-group "${AZURE_RESOURCE_GROUP}" \
  --cluster-name "${NAME}" \
  --output tsv --query "[0].name")

# Create system node pool tainted with `CriticalAddonsOnly=true:NoSchedule`
az aks nodepool add \
  --resource-group "${AZURE_RESOURCE_GROUP}" \
  --cluster-name "${NAME}" \
  --name systempool \
  --mode system \
  --node-count 1 \
  --node-taints "CriticalAddonsOnly=true:NoSchedule" \
  --no-wait

# Create user node pool tainted with `node.cilium.io/agent-not-ready=true:NoExecute`
az aks nodepool add \
  --resource-group "${AZURE_RESOURCE_GROUP}" \
  --cluster-name "${NAME}" \
  --name userpool \
  --mode user \
  --node-count 2 \
  --node-taints "node.cilium.io/agent-not-ready=true:NoExecute" \
  --no-wait

# Delete the initial system node pool
az aks nodepool delete \
  --resource-group "${AZURE_RESOURCE_GROUP}" \
  --cluster-name "${NAME}" \
  --name "${nodepool_to_delete}" \
  --no-wait

# Get the credentials to access the cluster with kubectl
az aks get-credentials --resource-group "${AZURE_RESOURCE_GROUP}" --name "${NAME}"

#cilium install --version -service-mesh:v1.11.0-beta.1 --config enable-envoy-config=true --kube-proxy-replacement=probe --rollback false --azure-resource-group tammach-963-group
#cilium install --version -ci:276d9d3505561bd329fd554b8277d9c826d4d3c5 --config enable-envoy-config=true,enable-ingress-controller=true --rollback false --azure-resource-group tammach-20704-group


kubectl set image daemonsets/cilium -n kube-system cilium-agent="quay.io/cilium/cilium-ci:276d9d3505561bd329fd554b8277d9c826d4d3c5"
kubectl set image deployment/cilium-operator -n kube-system cilium-operator="quay.io/cilium/operator-azure-ci:276d9d3505561bd329fd554b8277d9c826d4d3c5"