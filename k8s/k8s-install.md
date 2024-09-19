

## 基本配置

```bash
# 设置主机名
 hostnamectl set-hostname node2 <hostname>
 
#  禁swap
swapoff -a 

# 设置时区
timedatectl set-timezone Asia/Shanghai

# 添加host
cat >> /etc/hosts <<EOF
192.168.110.13 master 
192.168.110.11 node1 
192.168.110.12 node2 
EOF
```



## 系统设置

```bash
# 内核模块 
cat >> /etc/modules-load.d/containerd.conf <<EOF
overlay
br_netfilter
EOF

# 使生效
modprobe overlay
modprobe br_netfilter

#  
cat << EOF > /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
user.max_user_namespaces=28633
EOF

# 使生效 “99” 代表文件的优先级或顺序。
sysctl -p /etc/sysctl.d/99-kubernetes-cri.conf


# ipvs已经加入到了内核的主干，所以为kube-proxy开启ipvs的前提需要加载
cat >> /etc/modules-load.d/ipvs.conf <<EOF
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
EOF

# 使生效
modprobe ip_vs
modprobe ip_vs_rr
modprobe ip_vs_wrr
modprobe ip_vs_sh


```

## apt源配置

```bash
cat >> /etc/apt/sources.list <<EOF
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
apt install -y ipset ipvsadm bird auditd

# 其中docker-compose 从 https://github.com/docker/compose 下载对应的版本， 如 docker-compose-darwin-x86_64 并将其改名为 docker-compose 放到 /usr/bin/下

wget https://github.com/docker/compose/releases/download/v2.29.1/docker-compose-darwin-x86_64 
chmod +x docker-compose-darwin-x86_64 mv docker-compose-darwin-x86_64 /usr/bin/docker-compose
```

## containerd 

### 下载

```bash
# containerd 下载
wget https://github.com/containerd/containerd/releases/download/v1.7.3/containerd-1.7.3-linux-amd64.tar.gz
tar Cxzf /usr/local containerd-1.7.3-linux-amd64.tar.gz

# containerd 配置文件路径
mkdir -p /etc/containerd

# 生成containerd默认配置文件
containerd config default > /etc/containerd/config.toml
```

### 修改配置

```bash
vim /etc/containerd/config.toml 
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
cat >> /etc/systemd/system/containerd.service << EOF
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

# 使生效
systemctl daemon-reload
systemctl enable containerd --now 
systemctl status containerd
```

## runC

```bash
wget https://github.com/opencontainers/runc/releases/download/v1.1.9/runc.amd64
install -m 755 runc.amd64 /usr/local/sbin/runc
```



## crictl

```bash
wget https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.28.0/crictl-v1.28.0-linux-amd64.tar.gz
tar -zxvf crictl-v1.28.0-linux-amd64.tar.gz
install -m 755 crictl /usr/local/bin/crictl
```



> crictl --runtime-endpoint=unix:///run/containerd/containerd.sock  version
>
> Version:  0.1.0
> RuntimeName:  containerd
> RuntimeVersion:  v1.7.3
> RuntimeApiVersion:  v1



## kubeadm

```bash
# 配置源
tee /etc/apt/sources.list.d/kubernetes.list <<-'EOF'
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF

#
curl -s https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | sudo apt-key add -


# 更新
apt-get update

# 安装
apt install kubeadm kubectl kubelet -y

# 指定版本
apt install kubeadm=1.28.2 kubectl=1.28.2 kubelet=1.28.2

# 固定版本 为了防止依赖组件更新导致异常
apt-mark hold kubelet kubeadm kubectl

# kubelet 开机自启动
systemctl enable kubelet.service
# 输出默认配置
kubeadm config print init-defaults --component-configs KubeletConfiguration

# 初始化yaml配置文件
cat >>  kubeadm.yaml << EOF 
apiVersion: kubeadm.k8s.io/v1beta3
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 192.168.110.130  # 这一行需要调整
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
# 重置k8s
kubeadm reset

# 初始化
kubeadm init --config kubeadm.yaml

# 以下为输出
....
To start using your cluster, you need to run the following as a regular user:

#  需要操作如下
  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Alternatively, if you are the root user, you can run:

  export KUBECONFIG=/etc/kubernetes/admin.conf

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:

# node 节点加入时使用此token
# 注意ip
kubeadm join 192.168.1.185:6443 --token ag6egz.xjq1zz01meq8iboq \
        --discovery-token-ca-cert-hash sha256: axx........

```



### token-generate Again

```bash
# 如果忘记 join token ，可以在master节点上运行:
kubeadm token create --print-join-command
```

### 网络插件calico

```bash
# 下载 calico 网络插件
wget https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml
wget https://raw.githubusercontent.com/projectcalico/calico/v3.26.0/manifests/calico.yaml

# 修改  4800 行 CIDR 
4800             - name: CALICO_IPV4POOL_CIDR
4801               value: "192.168.10.0/16"  # 这个值和 kubeadm.yaml 中 10.244.0.0/16 有关系

# 运行
kubectl apply -f calico.yaml

# 获取相关信息
kubectl get node
kubectl get pod -A
```

### 命令行补全

```bash
# 安装bash-completion工具
apt install bash-completion

# 执行bash_completion
source /usr/share/bash-completion/bash_completion

# 在bash shell 中永久的添加自动补全（永久有效，不受切换终端影响）
echo "source <(kubectl completion bash)" >> ~/.bashrc
```



## 配置私有仓库

```bash
vim /etc/containerd/config.toml
....
[plugins."io.containerd.grpc.v1.cri".registry]
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
    [plugins."io.containerd.grpc.v1.cri".registry.mirrors."harbor-ip-or-domain:5000"]
      endpoint = ["http://octahub.8lab.cn:5000"]

    [plugins."io.containerd.grpc.v1.cri".registry.configs."harbor-ip-or-domain:5000".tls]
      insecure_skip_verify = true
          
  [plugins."io.containerd.grpc.v1.cri".registry.configs."harbor-ip-or-domain:5000".auth]
    username = "username"
    password = "password"
...

# 重启containerd
systemctl restart containerd
```



### 命令行下载

```bash
ctr -n=k8s.io image pull -u username:password --plain-http octahub.8lab.cn:5000/octa-cis/nisa-web-4.3:202202241657

# -n 指定下载的名称空间
# -u 使用指定的用户名和密码
# --plain-http ctr下载镜像默认为https 改为http

```

```bash
# 下载多个镜像
images="image:tag, \
image1:tag, \
image2:tag, \
image3:tag, \
image4:tag, \
...
"

# 使用逗号和空格作为分隔符进行循环
IFS=", "  # 设置内部字段分隔符 (Internal Field Separator)
for image in $images; do
    ctr -n=k8s.io image pull -u username:password --plain-http "$image"
done
```



## K8sMasterModifyIP

### master

#### 修改hosts 

> /etc/hosts 

```bash
# 如修改每台主机修改ip与主机对应关系
# 修改主机host 主机ip 由192.168.110.10 改为 192.168.110.131

vim /etc/hosts 
...
# 192.168.110.10 master 由这一行改为下一行内容
192.168.110.131 master
...
```

#### kubeadm.yaml

```yaml
# vim kubeadm.yaml
apiVersion: kubeadm.k8s.io/v1beta3
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 192.168.110.131 # 为修改后的ip
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

```

#### 调整 k8s 证书相关配置

```bash
# 停kubelet
systemctl stop kubelet

# 备份
mv /etc/kubernetes /etc/kubernetes-bak
mv /var/lib/kubelet/ /var/lib/kubelet-bak

# 保存pki 证书目录
mkdir -p /etc/kubernetes
cp -r /etc/kubernetes-bak/pki /etc/kubernetes
rm /etc/kubernetes/pki/{apiserver.*,etcd/peer.*}

# 重新初始化
# 通过 --ignore-preflight-errors=DirAvailable--var-lib-etcd 标志来告诉 kubeadm 使用预先存在的 etcd 数据。
kubeadm init --config kubeadm.yaml --ignore-preflight-errors=DirAvailable--var-lib-etcd

# 初始化后会得到join token 并记录

kubeadm join 192.168.110.131:6443 --token zgg6hj.ru2dhvnji8eml0z0 \
	--discovery-token-ca-cert-hash sha256:356f6382e894d517127acae9b0f756b13292f49f62a0b32627c5907302b34222

```

### Node 调整

```bash
# 修改主机名同上 略过 /etc/hosts

# 重置
kubeadm reset 

# 得新加入， 使用master 得到的token
kubeadm join 192.168.110.131:6443 --token zgg6hj.ru2dhvnji8eml0z0 \
	--discovery-token-ca-cert-hash sha256:356f6382e894d517127acae9b0f756b13292f49f62a0b32627c5907302b34222
```



## 外部镜像导入导出

> ctr 是 containerd 自带的工具，有命名空间的概念。Kubernetes 下使用的 containerd 默认命名空间是 k8s.io。所以在导入镜像时需要指定命令空间为 k8s.io，否则使用 crictl images 无法查询到。

```bash
# docat 导出 
docker save -o xx.tar image:tag 

# ctr 导出
ctr -n k8s.io  <output-file>.tar <image-name>:<tag>

# ctr 导入 , 这个才可以被k8s使用
ctr -n k8s.io image import xx.tar

# 普通导入
ctr image import xx.tar 

```

