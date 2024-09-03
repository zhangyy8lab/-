



```bash
ubuntu18.04 amdx86
k8s 1.28.2

master 192.168.1.185
node1 192.168.1.186
```



```bash
echo > /etc/hosts <<EOF
master 192.168.1.185
node1 192.168.1.186
EOF

```

## 系统设置

```bash
// 内核模块 
cat << EOF > /etc/modules-load.d/containerd.conf
overlay
br_netfilter
EOF

// 使生效
modprobe overlay
modprobe br_netfilter

// 
cat << EOF > /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
user.max_user_namespaces=28633
EOF
// 使生效
// “99” 代表文件的优先级或顺序。
sysctl -p /etc/sysctl.d/99-kubernetes-cri.conf


// ipvs已经加入到了内核的主干，所以为kube-proxy开启ipvs的前提需要加载
cat > /etc/modules-load.d/ipvs.conf <<EOF
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
EOF

// 使生效
modprobe ip_vs
modprobe ip_vs_rr
modprobe ip_vs_wrr
modprobe ip_vs_sh

// 禁swap
swapoff -a 

// 设置时区
timedatectl set-timezone Asia/Shanghai

// 设置主机名
 hostnamectl set-hostname master
```

## apt源配置

```bash
cat > /etc/apt/sources.list <<EOF
deb http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
EOF

// 生效
apt update 
```

## apt install

```bash
apt install -y ipset ipvsadm
```

## containerd 

### 下载

```bash
// containerd 下载
wget https://github.com/containerd/containerd/releases/download/v1.7.3/containerd-1.7.3-linux-amd64.tar.gz
tar Cxzf /usr/local containerd-1.7.3-linux-amd64.tar.gz

// containerd 配置文件路径
mkdir -p /etc/containerd

// 生成containerd默认配置文件
containerd config default > /etc/containerd/config.toml
```

### 修改配置

```bash
// 调整 containerd 相应配置
....
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
    SystemdCgroup = true  // 这个改为true 
....

....
# sandbox_image = "registry.k8s.io/pause:3.8"
  sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.9"
....

```



### systemd 启动 

```bash
cat << EOF > /etc/systemd/system/containerd.service
[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target local-fs.target

[Service]
#uncomment to enable the experimental sbservice (sandboxed) version of containerd/cri integration
#Environment="ENABLE_CRI_SANDBOXES=sandboxed"
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/containerd

Type=notify
Delegate=yes
KillMode=process
Restart=always
RestartSec=5
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNPROC=infinity
LimitCORE=infinity
LimitNOFILE=infinity
# Comment TasksMax if your systemd version does not supports it.
# Only systemd 226 and above support this version.
TasksMax=infinity
OOMScoreAdjust=-999

[Install]
WantedBy=multi-user.target
EOF

// 使生效
systemctl daemon-reload
systemctl enable containerd --now 
systemctl status containerd
```

## runC

```bash
// runc
wget https://github.com/opencontainers/runc/releases/download/v1.1.9/runc.amd64
install -m 755 runc.amd64 /usr/local/sbin/runc
```



## crictl

```bash
wget https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.28.0/crictl-v1.28.0-linux-amd64.tar.gz
tar -zxvf crictl-v1.28.0-linux-amd64.tar.gz
install -m 755 crictl /usr/local/bin/crictl
```



> test crictl --runtime-endpoint=unix:///run/containerd/containerd.sock  version
>
> Version:  0.1.0
> RuntimeName:  containerd
> RuntimeVersion:  v1.7.3
> RuntimeApiVersion:  v1

## kubeadm

### 源

```bash
tee /etc/apt/sources.list.d/kubernetes.list <<-'EOF'
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF

// 
curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -

```



```bash
// 更新
apt-get update

// 安装
apt install kubeadm kubectl kubelet
// 指定版本
apt install kubeadm=1.28.2 kubectl=1.28.2 kubelet=1.28.2

// 固定版本 为了防止依赖组件更新导致异常
apt-mark hold kubelet kubeadm kubectl
```



### init 配置

```bash
// kubelet 开机自启动
systemctl enable kubelet.service

// 输出默认配置
kubeadm config print init-defaults --component-configs KubeletConfiguration

// 写入配置文件
echo > kubeadm-init.yaml << EOF
apiVersion: kubeadm.k8s.io/v1beta3
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 192.168.1
  bindPort: 6443
nodeRegistration:
  criSocket: unix:///run/containerd/containerd.sock
  taints:
  - effect: PreferNoSchedule
    key: node-role.kubernetes.io/master
---
apiVersion: kubeadm.k8s.io/v1beta3
kind: ClusterConfiguration
kubernetesVersion: 1.28.2
imageRepository: registry.aliyuncs.com/google_containers
networking:
  podSubnet: 10.244.0.0/16
---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
cgroupDriver: systemd
failSwapOn: false
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: ipvs
EOF


```

### init k8s

```bash

// 初始化
kubeadm init --config kubeadm.yaml

// 以下为输出
....
To start using your cluster, you need to run the following as a regular user:

// 需要操作如下
  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Alternatively, if you are the root user, you can run:

  export KUBECONFIG=/etc/kubernetes/admin.conf

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

// node 节点加入时使用此token
kubeadm join 192.168.1.185:6443 --token ag6egz.xjq1zz01meq8iboq \
        --discovery-token-ca-cert-hash sha256: axx........


```



### token-generate Again

```bash
// 如果忘记 join token ，可以在master节点上运行:
kubeadm token create --print-join-command
```

### 网络插件calico

```bash
// 下载
wget https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml

// 修改  4800 行 CIDR 
4800             - name: CALICO_IPV4POOL_CIDR
4801               value: "192.168.10.0/16"

// 运行
kubectl apply -f calico.yaml


// 获取相关信息
kubectl get node
kubectl get pod -A
```

### 命令行补全

```bash
// 安装bash-completion工具
apt install bash-completion

// 执行bash_completion
source /usr/share/bash-completion/bash_completion

// 在bash shell 中永久的添加自动补全（永久有效，不受切换终端影响）
echo "source <(kubectl completion bash)" >> ~/.bashrc
```

