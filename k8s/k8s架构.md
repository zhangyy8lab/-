# K8S架构



```
Master：集群控制节点，每个集群需要至少一个master节点负责集群的管控

Node：工作负载节点，由master分配容器到这些node工作节点上，然后node节点上的docker负责容器的运行

Pod：kubernetes的最小控制单元，容器都是运行在pod中的，一个pod中可以有1个或者多个容器

Controller：控制器，通过它来实现对pod的管理，比如启动pod、停止pod、伸缩pod的数量等等

Service：pod对外服务的统一入口，下面可以维护者同一类的多个pod

Label：标签，用于对pod进行分类，同一类pod会拥有相同的标签

NameSpace：命名空间，用来隔离pod的运行环境
```





## Master

- api-server
  - 负责处理接受请求的工作
- etcd
  - 一致且高可用的键值存储，用作 Kubernetes 所有集群数据的后台数据库。
- Controller
  - 控制器
- Scheduler
  - 调度决策考虑的因素包括单个 Pod 及 Pods 集合创建在哪个node上



## Node

- kube-let
  - 每个节点上运行的主要 “节点代理”。它可以使用以下方式之一向 API 服务器注册：
    - 主机名（hostname）；
    - 覆盖主机名的参数；
    - 特定于某云驱动的逻辑。

- kube-proxy
  - 网络代理
- pod
  - Pod 中可以运行多个containers
  - 

## POD生命周期

在 Kubernetes 中，Pod 是最小的可部署单元，表示一个或多个容器的组合。Pod 具有以下生命周期阶段：

### Pending（挂起）

- 当创建 Pod 时，它会进入 Pending 阶段。
- 在此阶段，Kubernetes 正在为 Pod 分配资源（如 CPU 和内存），并等待这些资源可用。
- 如果所有的资源分配都成功，Pod 将进入下一个阶段。

### ContainerCreating

- 容器正在创建中。

### Running

- 在 Running 阶段，Pod 中的容器正在运行。
- 容器将在 Pod 中运行，并根据定义的规范执行其任务。
- 此阶段中的容器可以被创建、启动、重启等。

### Successed

- 当 Pod 中的所有容器成功完成其任务并退出时，Pod 将进入 Succeeded 阶段。
- 在此阶段，Pod 将保持运行状态，但不再重启或执行其他操作。
- 可以使用命令 `kubectl logs <pod名称>` 检查容器的日志输出。

### Failed

- 如果 Pod 中的任何容器失败并退出，Pod 将进入 Failed 阶段。
- 在此阶段，Pod 将保持运行状态，并尝试重启容器以解决问题。
- 可以使用命令 `kubectl describe pod <pod名称>` 获取有关失败原因的详细信息。

### Unknown

- 如果无法获取 Pod 的状态信息，则将其标记为 Unknown 阶段。
- 这可能是由于与 Pod 通信的问题导致的。
- 一旦通信恢复，Pod 的状态将被更新。



### Terminating

- Pod 正在终止中。

### CrashLoopBackOff 

- 容器在短时间内连续崩溃和重启。
