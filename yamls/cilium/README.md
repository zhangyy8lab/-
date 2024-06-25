#bandwidth
apiVersion: v1
kind: Pod
metadata:
  annotations:
    # Limits egress bandwidth to 10Mbit/s.
    kubernetes.io/egress-bandwidth: "10M"

#helm template cilium/cilium --debug ... After disbale the Hubble  
#because of: helm  install cilium cilium/cilium --version 1.13.0 --namespace kube-system --set kubeProxyReplacement=strict --set k8sServiceHost=192.168.1.249 --set k8sServicePort=6443 --set bandwidthManager.enabled=true --set bpf.masquerade=true --set-string extraConfig.enable-envoy-config=true --set tunnel=disabled --set autoDirectNodeRoutes=true  --set ipv4NativeRoutingCIDR=10.0.0.0/8 --set loadBalancer.mode=dsr --set hubble.relay.enabled=true --set hubble.ui.enabled=true  --set prometheus.enabled=true --set operator.prometheus.enabled=true --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"


#sed 修改 cilium-operator-deploy.yaml 和cilium-agent-ds.yaml中的KUBERNETES_SERVICE_HOST变量值
sed -i s%192.168.1.*%${APISERVER_ADDRESS}\"%g `grep -rl "192.168.1" *.yaml`
