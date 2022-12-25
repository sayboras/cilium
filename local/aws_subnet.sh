# Checking https://github.com/cilium/cilium/issues/18239

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: cni-configuration
  namespace: kube-system
data:
  cni-config: |-
    {
      "cniVersion":"0.3.1",
      "name":"cilium",
      "plugins": [
        {
          "cniVersion":"0.3.1",
          "type":"cilium-cni",
          "eni": {
            "subnet-tags":{
              "eksctl.cluster.k8s.io/v1alpha1/cluster-name":"tammach-17899"
            },
            "first-interface-index": 1
          }
        }
      ]
    }
EOF

#helm upgrade -i cilium install/kubernetes/cilium \
#  -n kube-system \
#  -f install/kubernetes/cilium/values/aws/aws-eni.yaml

helm upgrade -i cilium cilium/cilium --version 1.10.7 \
  -n kube-system \
  -f ~/go/src/github.com/cilium/cilium/install/kubernetes/cilium/values/aws/aws-eni-1-10.yaml