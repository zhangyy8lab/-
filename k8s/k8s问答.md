什么是k8s？说出你的理解
-------------

K8s是kubernetes的简称，其本质是一个开源的容器编排系统，主要用于管理容器化的应用，

其目标是让部署容器化的应用简单并且高效（powerful）,Kubernetes提供了应用部署，规划，更新，维护的一种机制。

说简单点：k8s就是一个编排容器的系统，一个可以管理容器应用全生命周期的工具，从创建应用，应用的部署，应用提供服务，扩容缩容应用，应用更新，都非常的方便，而且还可以做到故障自愈，

所以，k8s是一个非常强大的容器编排系统。

k8s的组件有哪些，作用分别是什么？
------------------

k8s主要由master节点和node节点构成。

master节点负责管理集群，node节点是容器应用真正运行的地方。

master节点包含的组件有：kube-api-server、kube-controller-manager、kube-scheduler、etcd。

node节点包含的组件有：kubelet、kube-proxy、container-runtime。

**kube-api-server：**

以下简称api-server，api-server是k8s最重要的核心组件之一，它是k8s集群管理的统一访问入口，提供了RESTful API接口, 实现了认证、授权和准入控制等安全功能；api-server还是其他组件之间的数据交互和通信的枢纽，其他组件彼此之间并不会直接通信，其他组件对资源对象的增、删、改、查和监听操作都是交由api-server处理后，api-server再提交给etcd数据库做持久化存储，只有api-server才能直接操作etcd数据库，其他组件都不能直接操作etcd数据库，其他组件都是通过api-server间接的读取，写入数据到etcd。

**kube-controller-manager：**

以下简称controller-manager，controller-manager是k8s中各种控制器的的管理者，是k8s集群内部的管理控制中心，也是k8s自动化功能的核心；controller-manager内部包含replication controller、node controller、deployment controller、endpoint controller等各种资源对象的控制器，每种控制器都负责一种特定资源的控制流程，而controller-manager正是这些controller的核心管理者。

**kube-scheduler：**

以下简称scheduler，scheduler负责集群资源调度，其作用是将待调度的pod通过一系列复杂的调度算法计算出最合适的node节点，然后将pod绑定到目标节点上。shceduler会根据pod的信息，全部节点信息列表，过滤掉不符合要求的节点，过滤出一批候选节点，然后给候选节点打分，选分最高的就是最佳节点，scheduler就会把目标pod安置到该节点。

**Etcd：**

etcd是一个分布式的键值对存储数据库，主要是用于保存k8s集群状态数据，比如，pod，service等资源对象的信息；etcd可以是单个也可以有多个，多个就是etcd数据库集群，etcd通常部署奇数个实例，在大规模集群中，etcd有5个或7个节点就足够了；另外说明一点，etcd本质上可以不与master节点部署在一起，只要master节点能通过网络连接etcd数据库即可。

**kubelet：**

每个node节点上都有一个kubelet服务进程，kubelet作为连接master和各node之间的桥梁，负责维护pod和容器的生命周期，当监听到master下发到本节点的任务时，比如创建、更新、终止pod等任务，kubelet 即通过控制docker来创建、更新、销毁容器；  
每个kubelet进程都会在api-server上注册本节点自身的信息，用于定期向master汇报本节点资源的使用情况。

**kube-proxy：**

kube-proxy运行在node节点上，在Node节点上实现Pod网络代理，维护网络规则和四层负载均衡工作，kube-proxy会监听api-server中从而获取service和endpoint的变化情况，创建并维护路由规则以提供服务IP和负载均衡功能。简单理解此进程是Service的透明代理兼负载均衡器，其核心功能是将到某个Service的访问请求转发到后端的多个Pod实例上。

container-runtime：容器运行时环境，即运行容器所需要的一系列程序，目前k8s支持的容器运行时有很多，如docker、rkt或其他，比较受欢迎的是docker，但是新版的k8s已经宣布弃用docker。

简述Kubernetes相关基础概念?
-------------------

答：

**master：**

k8s集群的管理节点，负责管理集群，提供集群的资源数据访问入口。拥有Etcd存储服务（可选），运行Api Server进程，Controller Manager服务进程及Scheduler服务进程；

**node（worker）：**

Node（worker）是Kubernetes集群架构中运行Pod的服务节点，是Kubernetes集群操作的单元，用来承载被分配Pod的运行，是Pod运行的宿主机。运行docker eninge服务，守护进程kunelet及负载均衡器kube-proxy；

**pod：**

运行于Node节点上，若干相关容器的组合。Pod内包含的容器运行在同一宿主机上，使用相同的网络命名空间、IP地址和端口，能够通过localhost进行通信。Pod是Kurbernetes进行创建、调度和管理的最小单位，它提供了比容器更高层次的抽象，使得部署和管理更加灵活。一个Pod可以包含一个容器或者多个相关容器；

**label：**

Kubernetes中的Label实质是一系列的Key/Value键值对，其中key与value可自定义。Label可以附加到各种资源对象上，如Node、Pod、Service、RC等。一个资源对象可以定义任意数量的Label，同一个Label也可以被添加到任意数量的资源对象上去。Kubernetes通过Label Selector（标签选择器）查询和筛选资源对象；

**Replication Controller：**

Replication Controller用来管理Pod的副本，保证集群中存在指定数量的Pod副本。

集群中副本的数量大于指定数量，则会停止指定数量之外的多余容器数量。反之，则会启动少于指定数量个数的容器，保证数量不变。

Replication Controller是实现弹性伸缩、动态扩容和滚动升级的核心；

**Deployment：**

Deployment在内部使用了RS来实现目的，Deployment相当于RC的一次升级，其最大的特色为可以随时获知当前Pod的部署进度；

**HPA（Horizontal Pod Autoscaler）：**

Pod的横向自动扩容，也是Kubernetes的一种资源，通过追踪分析RC控制的所有Pod目标的负载变化情况，来确定是否需要针对性的调整Pod副本数量；

**Service：**

Service定义了Pod的逻辑集合和访问该集合的策略，是真实服务的抽象。

Service提供了一个统一的服务访问入口以及服务代理和发现机制，关联多个相同Label的Pod，用户不需要了解后台Pod是如何运行；

**Volume：**

Volume是Pod中能够被多个容器访问的共享目录，Kubernetes中的Volume是定义在Pod上，可以被一个或多个Pod中的容器挂载到某个目录下；

**Namespace：**

Namespace用于实现多租户的资源隔离，可将集群内部的资源对象分配到不同的Namespace中，形成逻辑上的不同项目、小组或用户组，便于不同的Namespace在共享使用整个集群的资源的同时还能被分别管理；

简述Kubernetes和Docker的关系?
-----------------------

答：

Docker开源的容器引擎，一种更加轻量级的虚拟化技术；

K8s，容器管理工具，用来管理容器pod的集合，它可以实现容器集群的自动化部署、自动扩缩容、维护等功能；

简述Kubernetes如何实现集群管理?
---------------------

答：在集群管理方面，Kubernetes将集群中的机器划分为一个Master节点和一群工作节点Node。

其中，在Master节点运行着集群管理相关的一组进程kube-apiserver、kube-controller-manager和kube-scheduler，这些进程实现了整个集群的资源管理、Pod调度、弹性伸缩、安全控制、系统监控和纠错等管理能力，并且都是全自动完成的；

简述Kubernetes的优势、适应场景及其特点?
-------------------------

答：

优势：容器编排、轻量级、开源、弹性伸缩、负载均衡；

场景：快速部署应用、快速扩展应用、无缝对接新的应用功能、节省资源，优化硬件资源的使用；

特点：

可移植: 支持公有云、私有云、混合云、多重云（multi-cloud）、

可扩展: 模块化,、插件化、可挂载、可组合、

自动化: 自动部署、自动重启、自动复制、自动伸缩/扩展；

简述Kubernetes的缺点或当前的不足之处?
------------------------

答：

安装过程和配置相对困难复杂、管理服务相对繁琐、运行和编译需要很多时间、它比其他替代品更昂贵、对于简单的应用程序来说，可能不需要涉及Kubernetes即可满足；

简述Kubernetes中什么是Minikube、Kubectl、Kubelet?
-----------------------------------------

答：

Minikube 是一种可以在本地轻松运行一个单节点 Kubernetes 群集的工具；

Kubectl 是一个命令行工具，可以使用该工具控制Kubernetes集群管理器，如检查群集资源，创建、删除和更新组件，查看应用程序；

Kubelet 是一个代理服务，它在每个节点上运行，并使从服务器与主服务器通信；

kubelet的功能、作用是什么？（重点，经常会问）
--------------------------

答：kubelet部署在每个node节点上的，它主要有4个功能：  
1、节点管理。

kubelet启动时会向api-server进行注册，然后会定时的向api-server汇报本节点信息状态，资源使用状态等，这样master就能够知道node节点的资源剩余，节点是否失联等等相关的信息了。master知道了整个集群所有节点的资源情况，这对于 pod 的调度和正常运行至关重要。  
2、pod管理。

kubelet负责维护node节点上pod的生命周期，当kubelet监听到master的下发到自己节点的任务时，比如要创建、更新、删除一个pod，kubelet 就会通过CRI（容器运行时接口）插件来调用不同的容器运行时来创建、更新、删除容器；常见的容器运行时有docker、containerd、rkt等等这些容器运行时，我们最熟悉的就是docker了，但在新版本的k8s已经弃用docker了，k8s1.24版本中已经使用containerd作为容器运行时了。

3、容器健康检查。

pod中可以定义启动探针、存活探针、就绪探针等3种，我们最常用的就是存活探针、就绪探针，kubelet 会定期调用容器中的探针来检测容器是否存活，是否就绪，如果是存活探针，则会根据探测结果对检查失败的容器进行相应的重启策略；

4、Metrics Server资源监控。

在node节点上部署Metrics Server用于监控node节点、pod的CPU、内存、文件系统、网络使用等资源使用情况，而kubelet则通过Metrics Server获取所在节点及容器的上的数据。

kube-api-server的端口是多少？各个pod是如何访问kube-api-server的？
-------------------------------------------------

kube-api-server的端口是8080和6443，前者是http的端口，后者是https的端口，以我本机使用kubeadm安装的k8s为例：

在命名空间的kube-system命名空间里，有一个名称为kube-api-master的pod，

这个pod就是运行着kube-api-server进程，它绑定了master主机的ip地址和6443端口，但是在default命名空间下，存在一个叫kubernetes的服务，该服务对外暴露端口为443，目标端口6443，

这个服务的ip地址是clusterip地址池里面的第一个地址，同时这个服务的yaml定义里面并没有指定标签选择器，

也就是说这个kubernetes服务所对应的endpoint是手动创建的，该endpoint也是名称叫做kubernetes，该endpoint的yaml定义里面代理到master节点的6443端口，也就是kube-api-server的IP和端口。

这样一来，其他pod访问kube-api-server的整个流程就是：pod创建后嵌入了环境变量，pod获取到了kubernetes这个服务的ip和443端口，请求到kubernetes这个服务其实就是转发到了master节点上的6443端口的kube-api-server这个pod里面。

k8s中命名空间的作用是什么？
---------------

amespace是kubernetes系统中的一种非常重要的资源，namespace的主要作用是用来实现多套环境的资源隔离，或者说是多租户的资源隔离。

k8s通过将集群内部的资源分配到不同的namespace中，可以形成逻辑上的隔离，以方便不同的资源进行隔离使用和管理。

不同的命名空间可以存在同名的资源，命名空间为资源提供了一个作用域。

可以通过k8s的授权机制，将不同的namespace交给不同的租户进行管理，这样就实现了多租户的资源隔离，还可以结合k8s的资源配额机制，限定不同的租户能占用的资源，例如CPU使用量、内存使用量等等来实现租户可用资源的管理。

k8s提供了大量的REST接口，其中有一个是Kubernetes Proxy API接口，简述一下这个Proxy接口的作用，已经怎么使用。
---------------------------------------------------------------------

kubernetes proxy api接口，从名称中可以得知，proxy是代理的意思，其作用就是代理rest请求；

Kubernets API server 将接收到的rest请求转发到某个node上的kubelet守护进程的rest接口，由该kubelet进程负责响应。

我们可以使用这种Proxy接口来直接访问某个pod，这对于逐一排查pod异常问题很有帮助。  
下面是一些简单的例子：

```bash
http://<kube-api-server>:<api-sever-port>/api/v1/nodes/node名称/proxy/pods  	#查看指定node的所有pod信息
http://<kube-api-server>:<api-sever-port>/api/v1/nodes/node名称/proxy/stats  	#查看指定node的物理资源统计信息
http://<kube-api-server>:<api-sever-port>/api/v1/nodes/node名称/proxy/spec  	#查看指定node的概要信息

http://<kube-api-server>:<api-sever-port>/api/v1/namespace/命名名称/pods/pod名称/pod服务的url/  	#访问指定pod的程序页面
http://<kube-api-server>:<api-sever-port>/api/v1/namespace/命名名称/servers/svc名称/url/  	#访问指定server的url程序页面

```

pod是什么？
-------

在kubernetes的世界中，k8s并不直接处理容器，而是使用多个容器共存的理念，这组容器就叫做pod。

pod是k8s中可以创建和管理的最小单元，是资源对象模型中由用户创建或部署的最小资源对象模型，其他的资源对象都是用来支撑pod对象功能的，比如，pod控制器就是用来管理pod对象的，service或者imgress资源对象是用来暴露pod引用对象的，persistentvolume资源是用来为pod提供存储等等，

简而言之，k8s不会直接处理容器，而是pod，pod才是k8s中可以创建和管理的最小单元，也是基本单元。

pod的原理是什么？
----------

在微服务的概念里，一般的，一个容器会被设计为运行一个进程，除非进程本身产生子进程，

这样，由于不能将多个进程聚集在同一个单独的容器中，所以需要一种更高级的结构将容器绑定在一起，并将它们作为一个单元进行管理，这就是k8s中pod的背后原理。

pod有什么特点？
---------

1、每个pod就像一个独立的逻辑机器，k8s会为每个pod分配一个集群内部唯一的IP地址，所以每个pod都拥有自己的IP地址、主机名、进程等；  
2、一个pod可以包含1个或多个容器，1个容器一般被设计成只运行1个进程，1个pod只可能运行在单个节点上，即不可能1个pod跨节点运行，pod的生命周期是短暂，也就是说pod可能随时被消亡（如节点异常，pod异常等情况）；  
2、每一个pod都有一个特殊的被称为"根容器"的pause容器，也称info容器，pause容器对应的镜像属于k8s平台的一部分，除了pause容器，每个pod还包含一个或多个跑业务相关组件的应用容器；  
3、一个pod中的容器共享network命名空间；  
4、一个pod里的多个容器共享pod IP，这就意味着1个pod里面的多个容器的进程所占用的端口不能相同，否则在这个pod里面就会产生端口冲突；既然每个pod都有自己的IP和端口空间，那么对不同的两个pod来说就不可能存在端口冲突；  
5、应该将应用程序组织到多个pod中，而每个pod只包含紧密相关的组件或进程；  
6、pod是k8s中扩容、缩容的基本单位，也就是说k8s中扩容缩容是针对pod而言而非容器。

pod的重启策略有哪些？
------------

pod重启容器策略是指针对pod内所有容器的重启策略，不是重启pod，其可以通过restartPolicy字段配置pod重启容器的策略，如下：

*   Always: 当容器终止退出后，总是重启容器，默认策略就是Always。
    
*   OnFailure: 当容器异常退出，退出状态码非0时，才重启容器。
    
*   Never: 当容器终止退出，不管退出状态码是什么，从不重启容器。
    

pod的镜像拉取策略有哪几种？
---------------

pod镜像拉取策略可以通过imagePullPolicy字段配置镜像拉取策略，

主要有3中镜像拉取策略，如下：

*   IfNotPresent: 默认值，镜像在node节点宿主机上不存在时才拉取。
*   Always: 总是重新拉取，即每次创建pod都会重新从镜像仓库拉取一次镜像。
*   Never: 永远不会主动拉取镜像，仅使用本地镜像，需要你手动拉取镜像到node节点，如果node节点不存在镜像则pod启动失败。

kubenetes针对pod资源对象的健康监测机制?（必须记住3重探测方式，重点，经常问）
---------------------------------------------

提供了三类probe（探针）来执行对pod的健康监测：

*   livenessProbe探针 （存活探针）:

可以根据用户自定义规则来判定pod是否健康，用于判断容器是否处于Running状态，

如果不是，kubelet就会杀掉该容器，并根据重启策略做相应的处理。如果容器不包含该探针，那么kubelet就会默认返回值都是success;

*   ReadinessProbe探针:

同样是可以根据用户自定义规则来判断pod是否健康，容器服务是否可用（Ready），如果探测失败，控制器会将此pod从对应service的endpoint列表中移除，从此不再将任何请求调度到此Pod上，直到下次探测成功;

*   startupProbe探针:

启动检查机制，应用一些启动缓慢的业务，避免业务长时间启动而被上面两类探针kill掉，

这个问题也可以换另一种方式解决，就是定义上面两类探针机制时，初始化时间定义的长一些即可;

备注：每种探测方法能支持以下几个相同的检查参数，用于设置控制检查时间：

*   initialDelaySeconds：初始第一次探测间隔，用于应用启动的时间，防止应用还没启动而健康检查失败；
    
*   periodSeconds：检查间隔，多久执行probe检查，默认为10s；
    
*   timeoutSeconds：检查超时时长，探测应用timeout后为失败；
    
*   successThreshold：成功探测阈值，表示探测多少次为健康正常，默认探测1次。
    

就绪探针（ReadinessProbe探针）与存活探针（livenessProbe探针）区别是什么？
--------------------------------------------------

两者作用不一样，

存活探针是将检查失败的容器杀死，创建新的启动容器来保持pod正常工作；

就绪探针是，当就绪探针检查失败，并不重启容器，而是将pod移出endpoint，就绪探针确保了service中的pod都是可用的，确保客户端只与正常的pod交互并且客户端永远不会知道系统存在问题。

存活探针的属性参数有哪几个？
--------------

存活探针的附加属性参数有以下几个：

*   initialDelaySeconds：表示在容器启动后延时多久秒才开始探测；
    
*   periodSeconds：表示执行探测的频率，即间隔多少秒探测一次，默认间隔周期是10秒，最小1秒；
    
*   timeoutSeconds：表示探测超时时间，默认1秒，最小1秒，表示容器必须在超时时间范围内做出响应，否则视为本次探测失败；
    
*   successThreshold：表示最少连续探测成功多少次才被认定为成功，默认是1，对于liveness必须是1，最小值是1；
    
*   failureThreshold：表示连续探测失败多少次才被认定为失败，默认是3，连续3次失败，k8s 将根据pod重启策略对容器做出决定；
    

注意：定义存活探针时，一定要设置initialDelaySeconds属性，该属性为初始延时，如果不设置，默认容器启动时探针就开始探测了，这样可能会存在  
应用程序还未启动就绪，就会导致探针检测失败，k8s就会根据pod重启策略杀掉容器然后再重新创建容器的莫名其妙的问题。  
在生产环境中，一定要定义一个存活探针。

pod的就绪探针有哪几种？
-------------

我们知道，当一个pod启动后，就会立即加入service的endpoint ip列表中，并开始接收到客户端的链接请求，

假若此时pod中的容器的业务进程还没有初始化完毕，那么这些客户端链接请求就会失败，为了解决这个问题，kubernetes提供了就绪探针来解决这个问题的。

在pod中的容器定义一个就绪探针，就绪探针周期性检查容器，

如果就绪探针检查失败了，说明该pod还未准备就绪，不能接受客户端链接，则该pod将从endpoint列表中移除，

pod被剔除了, service就不会把请求分发给该pod，

然后就绪探针继续检查，如果随后容器就绪，则再重新把pod加回endpoint列表。

kubernetes提供了3种探测容器的存活探针，如下：

*   httpGet：通过容器的IP、端口、路径发送http 请求，返回200-400范围内的状态码表示成功。
    
*   exec：在容器内执行shell命令，根据命令退出状态码是否为0进行判断，0表示健康，非0表示不健康。
    
*   TCPSocket：与容器的IP、端口建立TCP Socket链接，能建立则说明探测成功，不能建立则说明探测失败
    

pod的就绪探针的属性参数有哪些
----------------

就绪探针的附加属性参数有以下几个：

*   initialDelaySeconds：延时秒数，即容器启动多少秒后才开始探测，不写默认容器启动就探测；
    
*   periodSeconds ：执行探测的频率（秒），默认为10秒，最低值为1；
    
*   timeoutSeconds ：超时时间，表示探测时在超时时间内必须得到响应，负责视为本次探测失败，默认为1秒，最小值为1；
    
*   failureThreshold ：连续探测失败的次数，视为本次探测失败，默认为3次，最小值为1次；
    
*   successThreshold ：连续探测成功的次数，视为本次探测成功，默认为1次，最小值为1次；
    

pod的重启策略是什么？
------------

答：通过命令“[kubectl](https://so.csdn.net/so/search?q=kubectl&spm=1001.2101.3001.7020) explain pod.spec”查看pod的重启策略;

*   Always：但凡pod对象终止就重启，此为默认策略;
    
*   OnFailure：仅在pod对象出现错误时才重启;
    

简单讲一下 pod创建过程
-------------

情况一、使用kubectl run命令创建的pod：

```dockerfile
注意：
kubectl run 在旧版本中创建的是deployment，
但在新的版本中创建的是pod则其创建过程不涉及deployment
```

如果是单独的创建一个pod，则其创建过程是这样的：  
1、首先，用户通过kubectl或其他api客户端工具提交需要创建的pod信息给apiserver；  
2、apiserver验证客户端的用户权限信息，验证通过开始处理创建请求生成pod对象信息，并将信息存入etcd，然后返回确认信息给客户端；  
3、apiserver开始反馈etcd中pod对象的变化，其他组件使用watch机制跟踪apiserver上的变动；  
4、scheduler发现有新的pod对象要创建，开始调用内部算法机制为pod分配最佳的主机，并将结果信息更新至apiserver；  
5、node节点上的kubelet通过watch机制跟踪apiserver发现有pod调度到本节点，尝试调用docker启动容器，并将结果反馈apiserver；  
6、apiserver将收到的pod状态信息存入etcd中。  
至此，整个pod创建完毕。

情况二、使用deployment来创建pod：

1、首先，用户使用kubectl create命令或者kubectl apply命令提交了要创建一个deployment资源请求；  
2、api-server收到创建资源的请求后，会对客户端操作进行身份认证，在客户端的~/.kube文件夹下，已经设置好了相关的用户认证信息，这样api-server会知道我是哪个用户，并对此用户进行鉴权，当api-server确定客户端的请求合法后，就会接受本次操作，并把相关的信息保存到etcd中，然后返回确认信息给客户端。  
3、apiserver开始反馈etcd中过程创建的对象的变化，其他组件使用watch机制跟踪apiserver上的变动。  
4、controller-manager组件会监听api-server的信息，controller-manager是有多个类型的，比如Deployment Controller, 它的作用就是负责监听Deployment，此时Deployment Controller发现有新的deployment要创建，那么它就会去创建一个ReplicaSet，一个ReplicaSet的产生，又被另一个叫做ReplicaSet Controller监听到了，紧接着它就会去分析ReplicaSet的语义，它了解到是要依照ReplicaSet的template去创建Pod, 它一看这个Pod并不存在，那么就新建此Pod，当Pod刚被创建时，它的nodeName属性值为空，代表着此Pod未被调度。  
5、调度器Scheduler组件开始介入工作，Scheduler也是通过watch机制跟踪apiserver上的变动，发现有未调度的Pod，则根据内部算法、节点资源情况，pod定义的亲和性反亲和性等等，调度器会综合的选出一批候选节点，在候选节点中选择一个最优的节点，然后将pod绑定该该节点，将信息反馈给api-server。  
6、kubelet组件布署于Node之上，它也是通过watch机制跟踪apiserver上的变动，监听到有一个Pod应该要被调度到自身所在Node上来，kubelet首先判断本地是否在此Pod，如果不存在，则会进入创建Pod流程，创建Pod有分为几种情况，第一种是容器不需要挂载外部存储，则相当于直接docker run把容器启动，但不会直接挂载docker网络，而是通过CNI调用网络插件配置容器网络，如果需要挂载外部存储，则还要调用CSI来挂载存储。kubelet创建完pod，将信息反馈给api-server，api-servier将pod信息写入etcd。  
7、Pod建立成功后，ReplicaSet Controller会对其持续进行关注，如果Pod因意外或被我们手动退出，ReplicaSet Controller会知道，并创建新的Pod，以保持replicas数量期望值。

k8s 创建一个pod的详细流程，涉及的组件怎么通信的？
----------------------------

答：

1）客户端提交创建请求，可以通过 api-server 提供的 restful 接口，或者是通过 kubectl 命令行工具，支持的数据类型包括 JSON 和 YAML；

2）api-server 处理用户请求，将 pod 信息存储至 etcd 中；

3）kube-scheduler 通过 api-server 提供的接口监控到未绑定的 pod，尝试为 pod 分配 node 节点，主要分为两个阶段，预选阶段和优选阶段，其中预选阶段是遍历所有的 node 节点，根据策略筛选出候选节点，而优选阶段是在第一步的基础上，为每一个候选节点进行打分，分数最高者胜出；

4）选择分数最高的节点，进行 pod binding 操作，并将结果存储至 etcd 中；

5）随后目标节点的 kubelet 进程通过 api-server 提供的接口监测到 kube-scheduler 产生的 pod 绑定事件，然后从 etcd 获取 pod 清单，下载镜像并启动容器；

简单描述一下pod的终止过程
--------------

1、用户向apiserver发送删除pod对象的命令；  
2、apiserver中的pod对象信息会随着时间的推移而更新，在宽限期内（默认30s），pod被视为dead；  
3、将pod标记为terminating状态；  
4、kubectl在监控到pod对象为terminating状态了就会启动pod关闭过程；  
5、endpoint控制器监控到pod对象的关闭行为时将其从所有匹配到此endpoint的server资源endpoint列表中删除；  
6、如果当前pod对象定义了preStop钩子处理器，则在其被标记为terminating后会意同步的方式启动执行；  
7、pod对象中的容器进程收到停止信息；  
8、宽限期结束后，若pod中还存在运行的进程，那么pod对象会收到立即终止的信息；  
9、kubelet请求apiserver将此pod资源的宽限期设置为0从而完成删除操作，此时pod对用户已不可见。

pod的生命周期有哪几种？
-------------

pod生命周期有的5种状态（也称5种相位），如下：

*   Pending（挂起）：API server已经创建pod，但是该pod还有一个或多个容器的镜像没有创建，包括正在下载镜像的过程；
    
*   Running（运行中）：Pod内所有的容器已经创建，且至少有一个容器处于运行状态、正在启动括正在重启状态；
    
*   Succeed（成功）：Pod内所有容器均已退出，且不会再重启；
    
*   Failed（失败）：Pod内所有容器均已退出，且至少有一个容器为退出失败状态
    
*   Unknown（未知）：某于某种原因apiserver无法获取该pod的状态，可能由于网络通行问题导致；
    

pod一致处于pending状态一般有哪些情况，怎么排查？（重点，持续更新）
--------------------------------------

（这个问题被问到的概率非常大）  
一个pod一开始创建的时候，它本身就是会处于pending状态，这时可能是正在拉取镜像，正在创建容器的过程。

如果等了一会发现pod一直处于pending状态，

那么我们可以使用kubectl describe命令查看一下pod的Events详细信息。一般可能会有这么几种情况导致pod一直处于pending状态：  
1、调度器调度失败。

Scheduer调度器无法为pod分配一个合适的node节点。

而这又会有很多种情况，比如，node节点处在cpu、内存压力，导致无节点可调度；pod定义了资源请求，没有node节点满足资源请求；node节点上有污点而pod没有定义容忍；pod中定义了亲和性或反亲和性而没有节点满足这些亲和性或反亲和性；以上是调度器调度失败的几种情况。  
2、pvc、pv无法动态创建。

如果因为pvc或pv无法动态创建，那么pod也会一直处于pending状态，比如要使用StatefulSet 创建redis集群，因为粗心大意，定义的storageClassName名称写错了，那么会造成无法创建pvc，这种情况pod也会一直处于pending状态，或者，即使pvc是正常创建了，但是由于某些异常原因导致动态供应存储无法正常创建pv，那么这种情况pod也会一直处于pending状态。

DaemonSet资源对象的特性？
-----------------

答：

DaemonSet这种资源对象会在每个k8s集群中的节点上运行，并且每个节点只能运行一个pod，这是它和deployment资源对象的最大也是唯一的区别。

所以，在其yaml文件中，不支持定义replicas，

除此之外，与Deployment、RS等资源对象的写法相同,

DaemonSet一般使用的场景有

*   在去做每个节点的日志收集工作；
*   监控每个节点的的运行状态;

删除一个Pod会发生什么事情？
---------------

答：

Kube-apiserver会接受到用户的删除指令，默认有30秒时间等待优雅退出，超过30秒会被标记为死亡状态，

此时Pod的状态Terminating，kubelet看到pod标记为Terminating就开始了关闭Pod的工作;

关闭流程如下：

1）pod从service的endpoint列表中被移除；

2)如果该pod定义了一个停止前的钩子，其会在pod内部被调用，停止钩子一般定义了如何优雅的结束进程；

3)进程被发送TERM信号（kill -14）;

4)当超过优雅退出的时间后，Pod中的所有进程都会被发送SIGKILL信号（kill -9）;

pod的共享资源？
---------

答：

1）PID 命名空间：Pod 中的不同应用程序可以看到其他应用程序的进程 ID；

2）网络命名空间：Pod 中的多个容器能够访问同一个IP和端口范围；

3）IPC 命名空间：Pod 中的多个容器能够使用 SystemV IPC 或 POSIX 消息队列进行通信；

4）UTS 命名空间：Pod 中的多个容器共享一个主机名；

5）Volumes（共享存储卷）：Pod 中的各个容器可以访问在 Pod 级别定义的 Volumes；

pod的初始化容器是干什么的？
---------------

init container，初始化容器用于在启动应用容器之前完成应用容器所需要的前置条件，

初始化容器本质上和应用容器是一样的，但是初始化容器是仅允许一次就结束的任务，初始化容器具有两大特征：

1、初始化容器必须运行完成直至结束，若某初始化容器运行失败，那么kubernetes需要重启它直到成功完成；  
2、初始化容器必须按照定义的顺序执行，当且仅当前一个初始化容器成功之后，后面的一个初始化容器才能运行；

pod的资源请求、限制如何定义？
----------------

pod的资源请求、资源限制可以直接在pod中定义

主要包括两块内容，

*   limits，限制pod能使用的最大cpu和内存，
*   requests，pod启动时申请的cpu和内存。

```bash
 resources:					#资源配额
      limits:					#限制最大资源，上限
        cpu: 2					#CPU限制，单位是code数
        memory: 2G				#内存最大限制
      requests:					#请求资源（最小，下限）
        cpu: 1					#CPU请求，单位是code数
        memory: 500G			#内存最小请求

```

pod的定义中有个command和args参数，这两个参数不会和docker镜像的entrypointc冲突吗？
--------------------------------------------------------

不会。

在pod中定义的command参数用于指定容器的启动命令列表，如果不指定，则默认使用Dockerfile打包时的启动命令，args参数用于容器的启动命令需要的参数列表；

特别说明：

kubernetes中的command、args其实是实现覆盖dockerfile中的ENTRYPOINT的功能的。

```bash
1、如果command和args均没有写，那么使用Dockerfile的配置；
2、如果command写了但args没写，那么Dockerfile默认的配置会被忽略，执行指定的command；
3、如果command没写但args写了，那么Dockerfile中的ENTRYPOINT的会被执行，使用当前args的参数；
4、如果command和args都写了，那么Dockerfile会被忽略，执行输入的command和args。


```

pause容器作用是什么？
-------------

每个pod里运行着一个特殊的被称之为pause的容器，也称根容器，而其他容器则称为业务容器；

创建pause容器主要是为了为业务容器提供 Linux命名空间，共享基础：包括 pid、icp、net 等，以及启动 init 进程，并收割僵尸进程；

这些业务容器共享pause容器的网络命名空间和volume挂载卷，

当pod被创建时，pod首先会创建pause容器，从而把其他业务容器加入pause容器，从而让所有业务容器都在同一个命名空间中，这样可以就可以实现网络共享。

pod还可以共享存储，在pod级别引入数据卷volume，业务容器都可以挂载这个数据卷从而实现持久化存储。

标签及标签选择器是什么，如何使用？
-----------------

标签是键值对类型，标签可以附加到任何资源对象上，主要用于管理对象，查询和筛选。

标签常被用于标签选择器的匹配度检查，从而完成资源筛选；一个资源可以定义一个或多个标签在其上面。

标签选择器，标签要与标签选择器结合在一起，标签选择器允许我们选择标记有特定标签的资源对象子集，如pod，并对这些特定标签的pod进行查询，删除等操作。

标签和标签选择器最重要的使用之一在于，在deployment中，在pod模板中定义pod的标签，然后在deployment定义标签选择器，这样就通过标签选择器来选择哪些pod是受其控制的，service也是通过标签选择器来关联哪些pod最后其服务后端pod。

service是如何与pod关联的？
------------------

答案是通过标签选择器，每一个由deployment创建的pod都带有标签，这样，service就可以定义标签选择器来关联哪些pod是作为其后端了，就是这样，service就与pod管联在一起了。

service的域名解析格式、pod的域名解析格式
-------------------------

service的DNS域名表示格式为`<servicename>.<namespace>.svc.<clusterdomain>`，

servicename是service的名称，namespace是service所处的命名空间，clusterdomain是k8s集群设置的域名后缀，一般默认为 cluster.local

pod的DNS域名格式为：`<pod-ip>.<namespace>.pod.<clusterdomain>` ，

其中，pod-ip需要使用-将ip直接的点替换掉，namespace为pod所在的命名空间，clusterdomain是k8s集群设置的域名后缀，一般默认为 cluster.local ，

演示如下：`10-244-1-223.default.pod.cluster.local`

对于deployment、daemonsets等创建的pod，还还可以通过`<pod-ip>.<deployment-name>.<namespace>.svc.<clusterdomain>` 这样的域名访问。

service的类型有哪几种
--------------

service的类型一般有4中，分别是：

*   ClusterIP：表示service仅供集群内部使用，默认值就是ClusterIP类型
    
*   NodePort：表示service可以对外访问应用，会在每个节点上暴露一个端口，这样外部浏览器访问地址为：任意节点的IP：NodePort就能连上service了
    
*   LoadBalancer：表示service对外访问应用，这种类型的service是公有云环境下的service，此模式需要外部云厂商的支持，需要有一个公网IP地址
    
*   ExternalName：这种类型的service会把集群外部的服务引入集群内部，这样集群内直接访问service就可以间接的使用集群外部服务了
    

一般情况下，service都是ClusterIP类型的，通过ingress接入的外部流量。

Pod到Service的通信？
---------------



1）k8s在创建服务时为服务分配一个虚拟IP，客户端通过该IP访问服务，服务则负责将请求转发到后端Pod上；

2）Service是通过kube-proxy服务进程实现，该进程在每个Node上均运行可以看作一个透明代理兼负载均衡器；

3）对每个TCP类型Service，kube-proxy都会在本地Node上建立一个SocketServer来负责接受请求，然后均匀发送到后端Pod默认采用Round Robin负载均衡算法；

4）Service的Cluster IP与NodePort等概念是kube-proxy通过Iptables的NAT转换实现，kube-proxy进程动态创建与Service相关的Iptables规则；

5）kube-proxy通过查询和监听API Server中Service与Endpoints的变化来实现其主要功能，包括为新创建的Service打开一个本地代理对象，接收请求针对针对发生变化的Service列表，kube-proxy会逐个处理；

一个应用pod是如何发现service的，或者说，pod里面的容器用于是如何连接service的？
-------------------------------------------------

答：有两种方式，一种是通过环境变量，另一种是通过service的dns域名方式。

1、环境变量：

当pod被创建之后，k8s系统会自动为容器注入集群内有效的service名称和端口号等信息为环境变量的形式，

这样容器应用直接通过取环境变量值就能访问service了，

如`curl http://${WEBAPP_SERVICE_HOST}:{WEBAPP_SERVICE_PORT}`

2、DNS方式：

使用dns域名解析的前提是k8s集群内有DNS域名解析服务器，

默认k8s中会有一个CoreDNS作为k8s集群的默认DNS服务器提供域名解析服务器；

service的DNS域名表示格式为`<servicename>.<namespace>.svc.<clusterdomain>`，

servicename是service的名称，namespace是service所处的命名空间，clusterdomain是k8s集群设置的域名后缀，一般默认为 cluster.local ，

这样容器应用直接通过service域名就能访问service了，

如`wget http://svc-deployment-nginx.default.svc.cluster.local:80`，

另外，service的port端口如果定义了名称，那么port也可以通过DNS进行解析，

格式为：`_<portname>._<protocol>.<servicename>.<namespace>.svc.<clusterdomain>`

如何创建一个service代理外部的服务，或者换句话来说，在k8s集群内的应用如何访问外部的服务，如数据库服务，缓存服务等?
--------------------------------------------------------------

答：可以通过创建一个没有标签选择器的service来代理集群外部的服务。

1、创建service时不指定selector标签选择器，但需要指定service的port端口、端口的name、端口协议等，这样创建出来的service因为没有指定标签选择器就不会自动创建endpoint；

2、手动创建一个与service同名的endpoint，endpoint中定义外部服务的IP和端口，endpoint的名称一定要与service的名称一样，端口协议也要一样，端口的name也要与service的端口的name一样，不然endpoint不能与service进行关联。

完成以上两步，k8s会自动将service和同名的endpoint进行关联，

这样，k8s集群内的应用服务直接访问这个service就可以相当于访问外部的服务了。

service、endpoint、kube-proxys三种的关系是什么？
-------------------------------------

**service**：

在kubernetes中，service是一种为一组功能相同的pod提供单一不变的接入点的资源。

当service被建立时，service的IP和端口不会改变，这样外部的客户端（也可以是集群内部的客户端）通过service的IP和端口来建立链接，这些链接会被路由到提供该服务的任意一个pod上。

通过这样的方式，客户端不需要知道每个单独提供服务的pod地址，这样pod就可以在集群中随时被创建或销毁。

**endpoint**：

service维护一个叫endpoint的资源列表，endpoint资源对象保存着service关联的pod的ip和端口。

从表面上看，当pod消失，service会在endpoint列表中剔除pod，当有新的pod加入，service就会将pod ip加入endpoint列表；

但是正在底层的逻辑是，endpoint的这种自动剔除、添加、更新pod的地址其实底层是由`endpoint controller`控制的，`endpoint controller`负责监听service和对应的pod副本的变化，如果监听到service被删除，则删除和该service同名的endpoint对象，如果监听到新的service被创建或者修改，则根据该service信息获取得相关pod列表，然后创建或更新service对应的endpoint对象，如果监听到pod事件，则更新它所对应的service的endpoint对象。

**kube-proxy**：

kube-proxy运行在node节点上，在Node节点上实现Pod网络代理，维护网络规则和四层负载均衡工作，

`kube-proxy`会监听`api-server`中从而获取service和endpoint的变化情况，创建并维护路由规则以提供服务IP和负载均衡功能。

简单理解此进程是Service的透明代理兼负载均衡器，其核心功能是将到某个Service的访问请求转发到后端的多个Pod实例上。

无头service和普通的service有什么区别，无头service使用场景是什么？
-------------------------------------------

答：

**无头service**没有cluster ip，在定义service时将 `service.spec.clusterIP：None`，就表示创建的是无头service。

**普通的service**是用于为一组后端pod提供请求连接的负载均衡，让客户端能通过固定的service ip地址来访问pod，这类的pod是没有状态的，同时service还具有负载均衡和服务发现的功能。普通service跟我们平时使用的nginx反向代理很相识。

试想这样一种情况，有6个redis pod ,它们相互之间要通信并要组成一个redis集群，

不需要所谓的service负载均衡，这时无头service就是派上用场了，

无头service由于没有cluster ip，kube-proxy就不会处理它也就不会对它生成规则负载均衡，无头service直接绑定的是pod 的ip。无头service仍会有标签选择器，有标签选择器就会有endpoint资源。

**无头service使用场景：**

无头service一般用于有状态的应用场景，如Kaka集群、Redis集群等，这类pod之间需要相互通信相互组成集群，不在需要所谓的service负载均衡。

deployment怎么扩容或缩容？
------------------

答：直接修改pod副本数即可，可以通过下面的方式来修改pod副本数：

1、直接修改yaml文件的replicas字段数值，然后`kubectl apply -f xxx.yaml`来实现更新；

2、使用`kubectl edit deployment xxx` 修改replicas来实现在线更新；

3、使用`kubectl scale --replicas=5 deployment/deployment-nginx`命令来扩容缩容。

deployment的更新升级策略有哪些？
---------------------

答：deployment的升级策略主要有两种。

1、Recreate 重建更新：这种更新策略会杀掉所有正在运行的pod，然后再重新创建的pod；

2、rollingUpdate 滚动更新：这种更新策略，deployment会以滚动更新的方式来逐个更新pod，同时通过设置滚动更新的两个参数`maxUnavailable、maxSurge`来控制更新的过程。

deployment的滚动更新策略有两个特别主要的参数，解释一下它们是什么意思？
----------------------------------------

答：deployment的滚动更新策略，rollingUpdate 策略，主要有两个参数，maxUnavailable、maxSurge。

*   maxUnavailable：最大不可用数，maxUnavailable用于指定deployment在更新的过程中不可用状态的pod的最大数量，maxUnavailable的值可以是一个整数值，也可以是pod期望副本的百分比，如25%，计算时向下取整。
    
*   maxSurge：最大激增数，maxSurge指定deployment在更新的过程中pod的总数量最大能超过pod副本数多少个，maxUnavailable的值可以是一个整数值，也可以是pod期望副本的百分比，如25%，计算时向上取整。
    

deployment更新的命令有哪些？
-------------------

答：可以通过三种方式来实现更新deployment。

1、直接修改yaml文件的镜像版本，然后`kubectl apply -f xxx.yaml`来实现更新；

2、使用`kubectl edit deployment xxx` 实现在线更新；

3、使用`kubectl set image deployment/nginx busybox=busybox nginx=nginx:1.9.1` 命令来更新。

简述一下deployment的更新过程?
--------------------

deployment是通过控制replicaset来实现，由replicaset真正创建pod副本，每更新一次deployment，都会创建新的replicaset，下面来举例deployment的更新过程：

假设要升级一个nginx-deployment的版本镜像为nginx:1.9，deployment的定义滚动更新参数如下：

```bash
replicas: 3
deployment.spec.strategy.type: RollingUpdate
maxUnavailable：25%
maxSurge：25%

```

通过计算我们得出，3\*25%=0.75，maxUnavailable是向下取整，则maxUnavailable=0，maxSurge是向上取整，则maxSurge=1，所以我们得出在整个deployment升级镜像过程中，不管旧的pod和新的pod是如何创建消亡的，pod总数最大不能超过3+maxSurge=4个，最大pod不可用数3-maxUnavailable=3个。

现在具体讲一下deployment的更新升级过程：

使用`kubectl set image deployment/nginx nginx=nginx:1.9 --record` 命令来更新；

1、deployment创建一个新的replaceset，先新增1个新版本pod，此时pod总数为4个，不能再新增了，再新增就超过pod总数4个了；旧=3，新=1，总=4；

2、减少一个旧版本的pod，此时pod总数为3个，这时不能再减少了，再减少就不满足最大pod不可用数3个了；旧=2，新=1，总=3；

3、再新增一个新版本的pod，此时pod总数为4个，不能再新增了；旧=2，新=2，总=4；

4、减少一个旧版本的pod，此时pod总数为3个，这时不能再减少了；旧=1，新=2，总=3；

5、再新增一个新版本的pod，此时pod总数为4个，不能再新增了；旧=1，新=3，总=4；

6、减少一个旧版本的pod，此时pod总数为3个，更新完成，pod都是新版本了；旧=0，新=3，总=3；

deployment的回滚使用什么命令
-------------------

在升级deployment时kubectl set image 命令加上 --record 参数可以记录具体的升级历史信息，

使用`kubectl rollout history deployment/deployment-nginx` 命令来查看指定的deployment升级历史记录，

如果需要回滚到某个指定的版本，可以使用`kubectl rollout undo deployment/deployment-nginx --to-revision=2` 命令来实现。

讲一下都有哪些存储卷，作用分别是什么?
-------------------

| 卷 | 作用 | 常用场景 |
| --- | --- | --- |
| emptyDir | 用于存储临时数据的简单空目录 | 一个pod中的多个容器需要共享彼此的数据 ，emptyDir的数据随着容器的消亡也会销毁 |
| hostPath | 用于将目录从工作节点的文件系统挂载到pod中 | 不常用，缺点是，pod的调度不是固定的，也就是当pod消失后deployment重新创建一个pod，而这pod如果不是被调度到之前pod的节点，那么该pod就不能访问之前的数据 |
| configMap | 用于将非敏感的数据保存到键值对中，使用时可以使用作为环境变量、命令行参数arg，存储卷被pods挂载使用 | 将应用程序的不敏感配置文件创建为configmap卷，在pod中挂载configmap卷，可是实现热更新 |
| secret | 主要用于存储和管理一些敏感数据，然后通过在 Pod 的容器里挂载 Volume 的方式或者环境变量的方式访问到这些 Secret 里保存的信息了，pod会自动解密Secret 的信息 | 将应用程序的账号密码等敏感信息通过secret卷的形式挂载到pod中使用 |
| downwardApi | 主要用于暴露pod元数据，如pod的名字 | pod中的应用程序需要指定pod的name等元数据，就可以通过downwardApi 卷的形式挂载给pod使用 |
| projected | 这是一种特殊的卷，用于将上面这些卷一次性的挂载给pod使用 | 将上面这些卷一次性的挂载给pod使用 |
| pvc | pvc是存储卷声明 | 通常会创建pvc表示对存储的申请，然后在pod中使用pvc |
| 网络存储卷 | pod挂载网络存储卷，这样就能将数据持久化到后端的存储里 | 常见的网络存储卷有nfs存储、glusterfs 卷、ceph rbd存储卷 |

> 注：本文以 PDF 持续更新，最新尼恩 架构笔记、面试题 的PDF文件，请从下面的链接获取：[码云](https://gitee.com/crazymaker/SimpleCrayIM/blob/master/%E7%96%AF%E7%8B%82%E5%88%9B%E5%AE%A2%E5%9C%88%E6%80%BB%E7%9B%AE%E5%BD%95.md) 或者 [语雀](https://www.yuque.com/crazymakercircle/gkkw8s/khigna)

pv的访问模式有哪几种
-----------

pv的访问模式有3种，如下：

*   ReadWriteOnce，简写：RWO 表示，只仅允许单个节点以读写方式挂载；
    
*   ReadOnlyMany，简写：ROX 表示，可以被许多节点以只读方式挂载；
    
*   ReadWriteMany，简写：RWX 表示，可以被多个节点以读写方式挂载；
    

pv的回收策略有哪几种
-----------

主要有3中回收策略：retain 保留、delete 删除、 Recycle回收。

*   Retain：保留，该策略允许手动回收资源，当删除PVC时，PV仍然存在，PV被视为已释放，管理员可以手动回收卷。
    
*   Delete：删除，如果Volume插件支持，删除PVC时会同时删除PV，动态卷默认为Delete，目前支持Delete的存储后端包括AWS EBS，GCE PD，Azure Disk，OpenStack Cinder等。
    
*   Recycle：回收，如果Volume插件支持，Recycle策略会对卷执行rm -rf清理该PV，并使其可用于下一个新的PVC，但是本策略将来会被弃用，目前只有NFS和HostPath支持该策略。（这种策略已经被废弃，不用记）
    

在pv的生命周期中，一般有几种状态
-----------------

pv一共有4中状态，分别是：

创建pv后，pv的的状态有以下4种：Available（可用）、Bound（已绑定）、Released（已释放）、Failed（失败）

```bash
Available，表示pv已经创建正常，处于可用状态；
Bound，表示pv已经被某个pvc绑定，注意，一个pv一旦被某个pvc绑定，那么该pvc就独占该pv，其他pvc不能再与该pv绑定；
Released，表示pvc被删除了，pv状态就会变成已释放；
Failed，表示pv的自动回收失败；

```

pv存储空间不足怎么扩容?
-------------

一般的，我们会使用动态分配存储资源，

在创建storageclass时指定参数 allowVolumeExpansion：true，表示允许用户通过修改pvc申请的存储空间自动完成pv的扩容，

当增大pvc的存储空间时，不会重新创建一个pv，而是扩容其绑定的后端pv。

这样就能完成扩容了。但是allowVolumeExpansion这个特性只支持扩容空间不支持减少空间。

存储类的资源回收策略:
-----------

主要有2中回收策略，delete 删除，默认就是delete策略、retain 保留。  
Retain：保留，该策略允许手动回收资源，当删除PVC时，PV仍然存在，PV被视为已释放，管理员可以手动回收卷。  
Delete：删除，如果Volume插件支持，删除PVC时会同时删除PV，动态卷默认为Delete，目前支持Delete的存储后端包括AWS EBS，GCE PD，Azure Disk，OpenStack Cinder等。

注意：使用存储类动态创建的pv默认继承存储类的回收策略，当然当pv创建后你也可以手动修改pv的回收策略。

怎么使一个node脱离集群调度，比如要停机维护单又不能影响业务应用
---------------------------------

使用kubectl drain 命令

k8s生产中遇到什么特别映像深刻的问题吗，问题排查解决思路是怎么样的？（重点）
---------------------------------------

（此问题被问到的概率高达90%，所以可以自己准备几个自己在生产环境中遇到的问题进行讲解）

答：前端的lb负载均衡服务器上的keepalived出现过脑裂现象。

1、当时问题现象是这样的，vip同时出现在主服务器和备服务器上，但业务上又没受到影响；  
2、这时首先去查看备服务器上的keepalived日志，发现有日志信息显示凌晨的时候备服务器出现了vrrp协议超时，所以才导致了备服务器接管了vip；查看主服务器上的keepalived日志，没有发现明显的报错信息，继续查看主服务器和备服务器上的keepalived进程状态，都是running状态的；查看主服务器上检测脚本所检测的进程，其进程也是正常的，也就是说主服务器根本没有成功执行检测脚本（成功执行检查脚本是会kill掉keepalived进程，脚本里面其实就是配置了检查nginx进程是否存活，如果检查到nginx不存活则kill掉keepalived，这样来实现备服务器接管vip）；  
3、排查服务器上的防火墙、selinux，防火墙状态和selinux状态都是关闭着的；  
4、使用tcpdump工具在备服务器上进行抓取数据包分析，分析发现，现在确实是备接管的vip，也确实是备服务器也在对外发送vrrp心跳包，所以现在外部流量应该都是流入备服务器上的vip；  
5、怀疑：主服务器上设置的vrrp心跳包时间间隔太长，以及检测脚本设置的检测时间设置不合理导致该问题；  
6、修改vrrp协议的心跳包时间间隔，由原来的2秒改成1秒就发送一次心跳包；检测脚本的检测时间也修改短一点，同时还修改检测脚本的检测失败的次数，比如连续检测2次失败才认定为检测失败；  
7、重启主备上的keepalived，现在keepalived是正常的，主服务器上有vip，备服务器上没有vip；  
8、持续观察：第二天又发现keepalived出现过脑裂现象，vip又同时出现在主服务器和备服务器上，又是凌晨的时候备服务器显示vrrp心跳包超时，所以才导致备服务器接管了vip；  
9、同样的时间，都是凌晨，vrrp协议超时；很奇怪，很有理由怀疑是网络问题，询问第三方厂家上层路由器是否禁止了vrrp协议，第三方厂家回复，没有禁止vrrp协议；  
10、百度、看官方文档求解；  
11、百度、看官网文档得知，keepalived有2种传播模式，一种是组播模式，一种是单播模式，keepalived默认在组播模式下工作，主服务器会往主播地址224.0.0.18发送心跳包，当局域网内有多个keepalived实例的时候，如果都用主播模式，会存在冲突干扰的情况，所以官方建议使用单播模式通信，单播模式就是点对点通行，即主向备服务器一对一的发送心跳包；

12、将keepalived模式改为单播模式，继续观察，无再发生脑裂现象。问题得以解决。

k8s生产中遇到什么特别映像深刻的问题吗，问题排查解决思路是怎么样的？（重点）
---------------------------------------

参考答案二：测试环境二进制搭建etcd集群，etcd集群出现2个leader的现象。  
1、问题现象就是：刚搭建的k8s集群，是测试环境的，搭建完成之后发现，使用kubectl get nodes 显示没有资源，kubectl get namespace 一会能正常显示全部的命名空间，一会又显示不了命名空间，这种奇怪情况。  
2、当时经验不是很足，第一点想到的是不是因为网络插件calico没装导致的，但是想想，即使没有安装网络插件，最多是node节点状态是notready，也不可能是没有资源发现呀；  
3、然后想到etcd数据库，k8s的资源都是存储在etcd数据库中的；  
4、查看etcd进程服务的启动状态，发现etcd服务状态是处于running状态，但是日志有大量的报错信息，日志大概报错信息就是集群节点的id不匹配，存在冲突等等报错信息；  
5、使用etcdctl命令查看etcd集群的健康状态，发现集群是health状态，但是居然显示有2个leader，这很奇怪（当初安装etcd的时候其实也只是简单看到了集群是健康状态，然后没注意到有2个leader，也没太关注etcd服务进程的日志报错信息，以为etcd集群状态是health状态就可以了）  
6、现在etcd出现了2个leader，肯定是存在问题的；  
7、全部检测一遍etcd的各个节点的配置文件，确认配置文件里面各个参数配置都没有问题，重启etcd集群，报错信息仍未解决，仍然存在2个leader；  
8、尝试把其中一个leader节点踢出集群，然后再重新添加它进入集群，仍然是报错，仍然显示有2个leader；  
9、尝试重新生成etcd的证书，重新颁发etcd的证书，问题仍然存在，仍然显示有2个leader；日志仍是报错集群节点的id不匹配，存在冲突；  
10、计算etcd命令的MD5值，确保各个节点的etcd命令是相同的，确保在scp传输的时候没有损耗等等，问题仍未解决；  
11、无解，请求同事，架构师介入帮忙排查问题，仍未解决；  
12、删除全部etcd相关的文件，重新部署etcd集群，etcd集群正常了，现在只有一个leader，使用命令kubectl get nodes 查看节点，也能正常显示了；  
13、最终问题的原因也没有定位出来，只能怀疑是环境问题了，由于是刚部署的k8s测试环境，etcd里面没有数据，所以可以删除重新创建etcd集群，如果是线上环境的etcd集群出现这种问题，就不能随便删除etcd集群了，必须要先进行数据备份才能进行其他方法的处理。

etcd集群节点可以设置为偶数个吗，为什么要设置为基数个呢？
------------------------------

不能，也不建议这么设置。

底层的原理，涉及到集群的脑裂 ，

具体的答案，请参考 《尼恩java面试宝典 专题14》

![在这里插入图片描述](https://img-blog.csdnimg.cn/e16f1bf77303495485b4df7e4584bdaa.png)

![在这里插入图片描述](https://img-blog.csdnimg.cn/f89139533bfa42a4a57447b276a19509.png)

etcd集群节点之间是怎么同步数据的？
-------------------

总体而言，是 通过Raft协议进行节点之间数据同步， 保证节点之间的数据一致性

在正式开始介绍 Raft 协议之间，我们有必要简单介绍一下其相关概念。

在现实的场景中，节点之间的一致性也就很难保证，这样就需要 Paxos、Raft 等一致性协议。

一致性协议可以保证在集群中大部分节点可用的情况下，集群依然可以工作并给出一个正确的结果，从而保证依赖于该集群的其他服务不受影响。

这里的“大部分节点可用”指的是集群中超过半数以上的节点可用，例如，集群中共有 5个节点，此时其中有 2 个节点出现故障宕机，剩余的可用节点数为 3，此时，集群中大多数节点处于可用的状态，从外部来看集群依然是可用的。

常见的一致性算法有Paxos、Raft等，

Paxos协议是Leslie Lamport于1990年提出的一种基于消息传递的、具有高度容错特性的一致性算法，Paxos 算法解决的主要问题是分布式系统内如何就某个值达成一致。在相当长的一段时间内，Paxos 算法几乎成为一致性算法的代名词，

但是 Paxos 有两个明显的缺点：第一个也是最明显的缺点就是 Paxos 算法难以理解，Paxos 算法的论文本身就比较晦涩难懂，要完全理解 Paxos 协议需要付出较大的努力，很多经验丰富的开发者在看完 Paxos 论文之后，无法将其有效地应用到具体工程实践中，这明显增加了工程化的门槛，也正因如此，才出现了几次用更简单的术语来解释 Paxos 的尝试。

Paxos算法的第二个缺点就是它没有提供构建现实系统的良好基础，也有很多工程化 Paxos 算法的尝试，但是它们对 Paxos 算法本身做了比较大的改动，彼此之间的实现差距都比较大，实现的功能和目的都有所不同，同时与Paxos算法的描述有很多出入。例如，著名Chubby，它实现了一个类Paxos的算法，但其中很多细节并未被明确。本章并不打算详细介绍 Paxos 协议的相关内容，如果读者对Paxos感兴趣，则可以参考Lamport发表的三篇论文：《The Part-Time Parliament》、《Paxos made simple》、《Fast Paxos》。

Raft算法是一种用于管理复制日志的一致性算法，其功能与Paxos算法相同类似，但其算法结构和Paxos算法不同，在设计Raft算法时设计者就将易于理解作为其目标之一，这使得Raft算法更易于构建实际的系统，大幅度减少了工程化的工作量，也方便开发者此基础上进行扩展。

Raft协议中，核心就是用于：

*   Leader选举
*   日志复制。

### Leader选举

Raft 协议的工作模式是一个 Leader 节点和多个 Follower 节点的模式，也就是常说的Leader-Follower 模式。

在 Raft 协议中，每个节点都维护了一个状态机，该状态机有三种状态，分别是Leader状态、Follower状态和Candidate状态，在任意时刻，集群中的任意一个节点都处于这三个状态之一。

各个状态和转换条件如图所示。  
![在这里插入图片描述](https://img-blog.csdnimg.cn/85663942ddfb437a8eeee92f79a531f4.png)

在多数情况下，集群中有一个Leader节点，其他节点都处于Follower状态，下面简单介绍一下每个状态的节点负责的主要工作。

*   Leader节点负责处理所有客户端的请求，当接收到客户端的写入请求时，Leader节点会在本地追加一条相应的日志，然后将其封装成消息发送到集群中其他的Follower节点。当Follower节点收到该消息时会对其进行响应。如果集群中多数（超过半数）节点都已收到该请求对应的日志记录时，则 Leader 节点认为该条日志记录已提交（committed），可以向客户端返回响应。Leader 还会处理客户端的只读请求，其中涉及一个简单的优化，后面介绍具体实现时，再进行详细介绍。Leader节点的另一项工作是定期向集群中的 Follower 节点发送心跳消息，这主要是为了防止集群中的其他Follower节点的选举计时器超时而触发新一轮选举。
    
*   Follower节点不会发送任何请求，它们只是简单地响应来自Leader或者Candidate 的请求；Follower节点也不处理Client的请求，而是将请求重定向给集群的Leader节点进行处理。
    
*   Candidate节点是由Follower节点转换而来的，当Follower节点长时间没有收到Leader节点发送的心跳消息时，则该节点的选举计时器就会过期，同时会将自身状态转换成Candidate，发起新一轮选举。选举的具体过程在下面详细描述。
    

了解了Raft协议中节点的三种状态及各个状态下节点的主要行为之后，我们通过一个示例介绍Raft协议中Leader选举的大致流程。为了方便描述，我们假设当前集群中有三个节点（A、B、C），如图所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/cd86035aff7640949ebe64082b63aa3c.png)

在Raft协议中有两个时间控制Leader选举发生，其中一个是选举超时时间（election timeout），每个Follower节点在接收不到Leader节点的心跳消息之后，并不会立即发起新一轮选举，而是需要等待一段时间之后才切换成Candidate状态发起新一轮选举。这段等待时长就是这里所说的election timeout（后面介绍etcd的具体实现时会提到，Follower节点等待的时长并不完全等于该配置）。之所以这样设计，主要是 Leader 节点发送的心跳消息可能因为瞬间的网络延迟或程序瞬间的卡顿而迟到（或是丢失），因此就触发新一轮选举是没有必要的。election timeout一般设置为150ms～300ms之间的随机数。另一个超时时间是心跳超时时间（heartbeat timeout），也就是Leader节点向集群中其他Follower节点发送心跳消息的时间间隔。

当集群初始化时，所有节点都处于 Follower 的状态，此时的集群中没有 Leader 节点。当Follower 节点一段时间（选举计时器超时）内收不到 Leader 节点的心跳消息，则认为 Leader节点出现故障导致其任期（Term）过期，Follower节点会转换成Candidate状态，发起新一轮的选举。所谓 “任期（Term）”，实际上就是一个全局的、连续递增的整数，在 Raft 协议中每进行一次选举，任期（Term）加一，在每个节点中都会记录当前的任期值（currentTerm）。每一个任期都是从一次选举开始的，在选举时，会出现一个或者多个 Candidate 节点尝试成为 Leader节点，如果其中一个Candidate节点赢得选举，则该节点就会切换为Leader状态并成为该任期的Leader节点，直到该任期结束。

回到前面的示例中，此时节点 A 由于长时间未收到 Leader 的心跳消息，就会切换成为Candidate状态并发起选举（节点A的选举计时器（election timer）已被重置）。

在选举过程中，节点A首先会将自己的选票投给自己，并会向集群中其他节点发送选举请求（Request Vote）以获取其选票，如图2-3（1）所示；此时的节点B和节点C还都是处于Term=0的任期之中，且都是Follower状态，均未投出Term=1任期中的选票，所以节点B和节点C在接收到节点A的选举请求后会将选票投给节点A，另外，节点B、C在收到节点A的选举请求的同时会将选举定时器重置，这是为了防止一个任期中同时出现多个Candidate节点，导致选举失败，如图2-3 （2）所示。

注意，节点B和节点C也会递增自身记录的Term值。  
![在这里插入图片描述](https://img-blog.csdnimg.cn/dc406fa0f9a2460fb765e48ff1c372e0.png)

在节点 A 收到节点 B、C 的投票之后，其收到了集群中超过半数的选票，所以在 Term=1这个任期中，该集群的Leader节点就是节点A，其他节点将切换成Follower状态，如图2-4所示。

另外需要读者了解的是，集群中的节点除了记录当期任期号（currentTerm），还会记录在该任期中当前节点的投票结果（VoteFor）。

![在这里插入图片描述](https://img-blog.csdnimg.cn/4ea4fde0b7f44bf38f1b159ae26aa0ef.png)

继续前面的示例，成为Term=1任期的Leader节点之后，节点A会定期向集群中的其他节点发送心跳消息，如图2-5（1）所示，

这样就可以防止节点B和节点C中的选举计时器（election timer）超时而触发新一轮的选举；当节点B和节点C（Follower）收到节点A的心跳消息之后会重置选举计时器，如图2-5（2）所示，由此可见，心跳超时时间（heartbeat timeout）需要远远小于选举超时时间（election timeout）

![在这里插入图片描述](https://img-blog.csdnimg.cn/d221517b5e794b6fac6f705a75225d27.png)

到这里读者可能会问，如果有两个或两个以上节点的选举计时器同时过期，则这些节点会同时由 Follower 状态切换成 Candidate 状态，然后同时触发新一轮选举，在该轮选举中，每个Candidate节点获取的选票都不到半数，无法选举出Leader节点，那么Raft协议会如何处理呢？这种情况确实存在，假设集群中有4个节点，其中节点A和节点B的选举计时器同时到期，切换到Candidate状态并向集群中其他节点发出选举请求，如图2-6（1）所示。

这里假设节点A发出的选举请求先抵达节点C，节点B发出的选举请求先抵达节点D，如图2-6（2）所示，节点A和节点B除了得到自身的选票之外，还分别得到了节点C和节点D投出的选票，得票数都是2，都没有超过半数。在这种情况下，Term=4这个任期会以选举失败结束，随着时间的流逝，当任意节点的选举计时器到期之后，会再次发起新一轮的选举。前面提到过election timeout是在一个时间区间内取的随机数，所以在配置合理的时候，像上述情况多次出现的概率并不大。

![在这里插入图片描述](https://img-blog.csdnimg.cn/e21e573adf37422096f4b7a7f26c66fe.png)

继续上面的示例，这里假设节点A的选举计时器再次到期（此次节点B、C、D 的选举计时器并未到期），它会切换成Candidate状态并发起新一轮选举（Term=5），如图2-7（1）所示，其中节点B虽然处于Candidate状态，但是接收到Term值比自身记录的Term值大的请求时，节点会切换成Follower状态并更新自身记录的Term值，所以该示例中的节点B也会将选票投给节点A，如图2-7（2）所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/5bfb92202e9d45a8b9d97a4830903d0b.png)

在获取集群中半数以上的选票并成为新任期（Term=5）的 Leader 之后，节点 A 会定期向集群中其他节点发送心跳消息；当集群中其他节点收到Leader节点的心跳消息的时候，会重置选举定时器，如图2-8所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/73a2af9ee6d74261b10f67288f9aca6a.png)

介绍完集群启动时的Leader选举流程之后，下面分析Leader节点宕机之后重新选举的场景。继续上述4节点集群的示例，在系统运行一段时间后，集群当前的Leader节点（A）因为故障而宕机，此时将不再有心跳消息发送到集群的其他Follower节点（节点B、C、D），一段时间后，会有一个Follower节点的选举计时器最先超时，这里假设节点D的选举计时器最先超时，然后它将切换为Candidate状态并发起新一轮选举，如图2-9所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/a4ce0a1aca334403b566633aa80b3b33.png)

当节点B和节点C收到节点D的选举请求后，会将其选票投给节点D，由于节点A已经宕机，没有参加此次选举，也就无法进行投票，但是在此轮选举中，节点D依然获得了半数以上的选票，故成为新任期（Term=6）的Leader节点，并开始向其他Follower节点发送心跳消息，如图2-10所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/58fcfb24b2ff4b5db20c5ec2a0fca763.png)

当节点A恢复之后，会收到节点D发来的心跳消息，该消息中携带的任期号（Term=6）大于节点A当前记录的任期号（Term=5），所以节点A会切换成Follower状态。在Raft协议中，当某个节点接收到的消息所携带的任期号大于当前节点本身记录的任期号，那么该节点会更新自身记录的任期号，同时会切换为Follower状态并重置选举计时器，这是Raft算法中所有节点最后请读者考虑一个场景：如果集群中选出的Leader节点频繁崩溃或是其他原因导致选举频繁发生，这会使整个集群中没有一个稳定的Leader节点，这样客户端无法与集群中的Leader节点正常交互，也就会导致整个集群无法正常工作。

Leader选举是Raft算法中对时间要求较为严格的一个点，一般要求整个集群中的时间满足如下不等式：  
广播时间 ＜＜ 选举超时时间 ＜＜ 平均故障间隔时间

在上述不等式中，广播时间指的是从一个节点发送心跳消息到集群中的其他节点并接收响应的平均时间；平均故障间隔时间就是对于一个节点而言，两次故障之间的平均时间。为了保证整个Raft集群可用，广播时间必须比选举超时时间小一个数量级，这样Leader节点才能够发送稳定的心跳消息来重置其他 Follower 节点的选举计时器，从而防止它们切换成 Candidate 状态，触发新一轮选举。在前面的描述中也提到过，选举超时时间是一个随机数，通过这种随机的方式，会使得多个Candidate节点瓜分选票的情况明显减少，也就减少了选举耗时。

另外，选举超时时间应该比平均故障间隔时间小几个数量级，这样Leader节点才能稳定存在，整个集群才能稳定运行。当Leader节点崩溃之后，整个集群会有大约相当于选举超时的时间不可用，这种情况占比整个集群稳定运行的时间还是非常小的。

广播时间和平均故障间隔时间是由网络和服务器本身决定的，但是选举超时时间是可以由我们自己调节的。

一般情况下，广播时间可以做到0.5ms～50ms，选举超时时间设置为200ms～1s之间，而大多数服务器的平均故障间隔时间都在几个月甚至更长，很容易满足上述不等式的时间需求。

### 日志复制

通过上一节介绍的Leader选举过程，集群中最终会选举出一个Leader节点，而集群中剩余的其他节点将会成为Follower节点。

Leader节点除了向Follower节点发送心跳消息，**还会处理客户端的请求**，并将客户端的更新操作以消息（Append Entries消息）的形式发送到集群中所有的Follower节点。

当Follower节点记录收到的这些消息之后，会向Leader节点返回相应的响应消息。当Leader节点在收到半数以上的Follower节点的响应消息之后，会对客户端的请求进行应答。

最后，Leader会提交客户端的更新操作，该过程会发送Append Entries消息到Follower节点，通知Follower节点该操作已经提交，同时Leader节点和Follower节点也就可以将该操作应用到自己的状态机中。

上面这段描述仅仅是Raft协议中日志复制部分的大致流程，下面我们依然通过一个示例描述该过程，为了方便描述，我们依然假设当前集群中有三个节点（A、B、C），其中A是Leader节点，B、C是Follower 节点，此时有一个客户端发送了一个更新操作到集群，如图 2-11（1）所示。前面提到过，集群中只有Leader节点才能处理客户端的更新操作，这里假设客户端直接将请求发给了节点A。当收到客户端的请求时，节点A会将该更新操作记录到本地的Log中，如图2-11（2）所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/4c36164331ed456c91ab9e2ab49ba0f5.png)

之后，节点A会向其他节点发送Append Entries消息，其中记录了Leader节点最近接收到的请求日志，如图2-12（1）所示。集群中其他Follower节点收到该Append Entries消息之后，会将该操作记录到本地的Log中，并返回相应的响应消息，如图2-12所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/e80d5392060e4919917f6f501a79cdc9.png)

当Leader节点收到半数以上的响应消息之后，会认为集群中有半数以上的节点已经记录了该更新操作，Leader 节点会将该更新操作对应的日志记录设置为已提交（committed），并应用到自身的状态机中。同时 Leader 节点还会对客户端的请求做出响应，如图 2-13（1）所示。同时，Leader节点也会向集群中的其他Follower节点发送消息，通知它们该更新操作已经被提交，Follower节点收到该消息之后，才会将该更新操作应用到自己的状态机中，如图2-13（2）所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/9dc41bd27ea142fc94d347b4fd0b1182.png)

在上述示例的描述中我们可以看到，集群中各个节点都会维护一个本地Log用于记录更新操作，除此之外，每个节点还会维护commitIndex和lastApplied两个值，它们是本地Log的索引值，其中commitIndex表示的是当前节点已知的、最大的、已提交的日志索引值，lastApplied表示的是当前节点最后一条被应用到状态机中的日志索引值。当节点中的 commitIndex 值大于lastApplied值时，会将lastApplied 加1，并将lastApplied对应的日志应用到其状态机中。

在Leader节点中不仅需要知道自己的上述信息，还需要了解集群中其他Follower节点的这些信息，例如，Leader节点需要了解每个Follower节点的日志复制到哪个位置，从而决定下次发送 Append Entries 消息中包含哪些日志记录。为此，Leader 节点会维护 nextIndex\[\]和matchIndex\[\]两个数组，这两个数组中记录的都是日志索引值，其中nextIndex\[\]数组记录了需要发送给每个 Follower 节点的下一条日志的索引值，matchIndex\[\]表示记录了已经复制给每个Follower节点的最大的日志索引值。

这里简单看一下 Leader 节点与某一个 Follower 节点复制日志时，对应 nextIndex 和matchIndex值的变化：Follower节点中最后一条日志的索引值大于等于该Follower节点对应的nextIndex 值，那么通过 Append Entries 消息发送从 nextIndex 开始的所有日志。之后，Leader节点会检测该 Follower 节点返回的相应响应，如果成功则更新相应该 Follower 节点对应的nextIndex值和matchIndex值；如果因为日志不一致而失败，则减少nextIndex值重试。

下面我们依然通过一个示例来说明nextIndex\[\]和matchIndex\[\]在日志复制过程中的作用，假设集群现在有三个节点，其中节点A是Leader节点（Term=1），而Follower节点C因为宕机导致有一段时间未与Leader节点同步日志。此时，节点C的Log中并不包含全部的已提交日志，而只是节点A的Log的子集，节点C故障排除后重新启动，当前集群的状态如图2-14所示（这里只关心Log、nextIndex\[\]、matchIndex\[\]，其他的细节省略，另外需要注意的是，图中的Term=1表示的是日志发送时的任期号，而非当前的任期号）。

![在这里插入图片描述](https://img-blog.csdnimg.cn/26bd306d2aed48b5a3a3e5e09454a0a9.png)

A作为Leader节点，记录了nextIndex\[\]和matchIndex\[\]，所以知道应该向节点C发送哪些日志，在本例中，Leader节点在下次发送Append Entries消息时会携带Index=2的消息（这里为了描述简单，每条消息只携带单条日志，Raft协议采用批量发送的方式，这样效率更高），如图2-15（1）所示。当节点C收到Append Entries消息后，会将日志记录到本地Log中，然后向Leader 节点返回追加日志成功的响应，当 Leader 节点收到响应之后，会递增节点 C 对应的nextIndex和matchIndex，这样Leader节点就知道下次发送日志的位置了，该过程如图2-15（2）所示。

在上例中，当Leader节点并未发生过切换，所以Leader节点始终准确地知道节点C对应nextIndex值和matchIndex值。

如果在上述示例中，在节点C故障恢复后，节点A宕机后重启，并且导致节点B成为新任期（Term=2）的 Leader 节点，则此时节点 B 并不知道旧 Leader 节点中记录的 nextIndex\[\]和matchIndex\[\]信息，所以新Leader节点会重置nextIndex\[\]和matchIndex\[\]，其中会将nextIndex\[\]全部重置为其自身Log的最后一条已提交日志的Index值，而matchIndex\[\]全部重置为0，如图2-16所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/3d6bbeecf7f54029b3e4e226053551e3.png)

![在这里插入图片描述](https://img-blog.csdnimg.cn/eb411fb6a81245bfae053de9ea950ff5.png)

随后，新任期中的Leader节点会向其他节点发送Append Entries消息，如图2-17（1）所示，节点A已经拥有了当前Leader的全部日志记录，所以会返回追加成功的响应并等待后续的日志，而节点C并没有Index=2和Index=3两条日志，所以返回追加日志失败的响应，在收到该响应后，Leader节点会将nextIndex前移，如图2-17（2）所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/5abf8298a4cc470f9306f20b123afcd6.png)

然后新 Leader 节点会再次尝试发送 Append Entries 消息，循环往复，不断减小 nextIndex值，直至节点C返回追加成功的响应，之后就进入了正常追加消息记录的流程，不再赘述。

了解了 Log 日志及节点中基本的数据结构之后，请读者回顾前面描述的选举过程，

其中Follower节点的投票过程并不像前面描述的那样简单（先收到哪个Candidate节点的选举请求，就将选票投给哪个Candidate节点），Follower节点还需要比较该Candidate节点的日志记录与自身的日志记录，拒绝那些日志没有自己新的Candidate节点发来的投票请求，确保将选票投给包含了全部已提交（committed）日志记录的 Candidate 节点。

这也就保证了已提交的日志记录不会丢失：Candidate节点为了成为Leader节点，必然会在选举过程中向集群中半数以上的节点发送选举请求，因为已提交的日志记录必须存在集群中半数以上的节点中，这也就意味着每一条已提交的日志记录肯定在这些接收到节点中的至少存在一份。也就是说，记录全部已提交日志的节点和接收到Candidate节点的选举请求的节点必然存在交集，如图2-18所示。

![在这里插入图片描述](https://img-blog.csdnimg.cn/e490ad24562848c18cdccfe5ce853753.png)

如果Candidate节点上的日志记录与集群中大多数节点上的日志记录一样新，那么其日志一定包含所有已经提交的日志记录，也就可以获得这些节点的投票并成为Leader。

在比较两个节点的日志新旧时，Raft 协议通过比较两节点日志中的最后一条日志记录的索引值和任期号，以决定谁的日志比较新：首先会比较最后一条日志记录的任期号，如果最后的日志记录的任期号不同，那么任期号大的日志记录比较新；如果最后一条日志记录的任期号相同，那么日志索引较大的 比较新。

这里只是大概介绍一下 Raft 协议的流程和节点使用的各种数据结构，读者需要了解的是Raft 协议的工作原理，如果对上述数据结构描述感到困惑，在后面介绍etcd-raft 模块时，还会再次涉及这些数据结构，到时候读者可以结合代码及这里的描述进一步进行分析。

请详述kube-proxy原理?
----------------

答：集群中每个Node上都会运行一个kube-proxy服务进程，他是Service的透明代理兼均衡负载器，其核心功能是将某个Service的访问转发到后端的多个Pod上。

kube-proxy通过监听集群状态变更，并对本机iptables做修改，从而实现网络路由。

而其中的负载均衡，也是通过iptables的特性实现的。

从V1.8版本开始，用IPVS（IP Virtual Server）模式，用于路由规则的配置，主要优势是：

1）为大型集群提供了更好的扩展性和性能。采用哈希表的数据结构，更高效；

2）支持更复杂的负载均衡算法；

3）支持服务器健康检查和连接重试；

4）可以动态修改ipset的集合；

flannel 和 ovs 网络的区别？
--------------------

答：

1）配置是否自动化：OpenvSwitch（ovs）作为开源的交换机软件，相对比较成熟和稳定，支持各种网络隧道和协议，经历了大型项目 OpenStack 的考验，而 flannel 除了支持建立覆盖网络来实现 Pod 到 Pod 之间的无缝通信之外，还跟 docker、k8s 的架构体系紧密结合，flannel 能感知 k8s 中的 service 对象，然后动态维护自己的路由表，并通过 etcd 来协助 docker 对整个 k8s 集群的 docker0 网段进行规范，而 ovs ，这些操作则需要手动完成，假如集群中有 N 个节点，则需要建立 N(N-1)/2 个 Vxlan 或者 gre 连接，这取决于集群的规模，如果集群的规模很大，则必须通过自动化脚本来初始化，避免出错。

2）是否支持隔离：flannel 虽然很方便实现 Pod 到 Pod 之间的通信，但不能实现多租户隔离，也不能很好地限制 Pod 的网络流量，而 ovs 网络有两种模式：单租户模式和多租户模式，单租户模式直接使用 openvswitch + vxlan 将 k8s 的 pod 网络组成一个大二层，所有的 pod 可以互相通信访问，多租户模式以 Namespace 为维度分配虚拟网络，从而形成一个网络独立用户，一个 Namespace 中的 pod 无法访问其他 Namespace 中的 pod 和 svc 对象；

k8s集群外流量怎么访问Pod？
----------------

答：

可以通过Service的NodePort方式访问，会在所有节点监听同一个端口，比如：30000，访问节点的流量会被重定向到对应的Service上面；

K8S 资源限制 QoS？
-------------

答：Quality of Service（Qos）

主要有三种类别：

1）BestEffort：什么都不设置（CPU or Memory），佛系申请资源；

2）Burstable：Pod 中的容器至少一个设置了CPU 或者 Memory 的请求；

3）Guaranteed：Pod 中的所有容器必须设置 CPU 和 Memory，并且 request 和 limit 值相等；

k8s数据持久化的方式有哪些？
---------------

答：

1)EmptyDir（空目录）：没有指定要挂载宿主机上的某个目录，直接由Pod内保部映射到宿主机上。类似于docker中的manager volume；场景有：a.只需要临时将数据保存在磁盘上，比如在合并/排序算法中；b.作为两个容器的共享存储，使得第一个内容管理的容器可以将生成的数据存入其中，同时由同一个webserver容器对外提供这些页面;emptyDir的特性：同个pod里面的不同容器，共享同一个持久化目录，当pod节点删除时，volume的数据也会被删除。如果仅仅是容器被销毁，pod还在，则不会影响volume中的数据。总结来说：emptyDir的数据持久化的生命周期和使用的pod一致。一般是作为临时存储使用。

2）Hostpath：将宿主机上已存在的目录或文件挂载到容器内部。类似于docker中的bind mount挂载方式；

3）PersistentVolume（简称PV）：基于NFS服务的PV，也可以基于GFS的PV。它的作用是统一数据持久化目录，方便管理，PVC是向PV申请应用所需的容量大小，K8s集群中可能会有多个PV，PVC和PV若要关联，其定义的访问模式必须一致。定义的storageClassName也必须一致，若群集中存在相同的（名字、访问模式都一致）两个PV，那么PVC会选择向它所需容量接近的PV去申请，或者随机申请；

K8S的基本组成部分？
-----------

答：

Master节点主要有五个组件，分别是kubectl、api-server、controller-manager、kube-scheduler 和 etcd；

node节点主要有三个组件，分别是 kubelet、kube-proxy 和 容器运行时 docker 或者 rkt；

kubectl：客户端命令行工具，作为整个系统的操作入口。  
apiserver：以REST API服务形式提供接口，作为整个系统的控制入口。  
controller-manager：执行整个系统的后台任务，包括节点状态状况、Pod个数、Pods和Service的关联等。  
kube-scheduler：负责节点资源管理，接收来自kube-apiserver创建Pods任务，并分配到某个节点。  
etcd：负责节点间的服务发现和配置共享。  
kube-proxy：运行在每个计算节点上，负责Pod网络代理。定时从etcd获取到service信息来做相应的策略。  
kubelet：运行在每个计算节点上，作为agent，接收分配该节点的Pods任务及管理容器，周期性获取容器状态，反馈给kube-apiserver。  
DNS：一个可选的DNS服务，用于为每个Service对象创建DNS记录，这样所有的Pod就可以通过DNS访问服务了。

K8s中镜像的下载策略是什么？
---------------

答：可通过命令“kubectl explain pod.spec.containers”来查看imagePullPolicy这行的解释，

K8s的镜像下载策略有三种：

Always：镜像标签为latest时，总是从指定的仓库中获取镜像；

Never：禁止从仓库中下载镜像，也就是说只能使用本地镜像；

IfNotPresent：仅当本地没有对应镜像时，才从目标仓库中下载；

标签与标签选择器的作用是什么？
---------------

答：标签：是当相同类型的资源对象越来越多的时候，为了更好的管理，可以按照标签将其分为一个组，为的是提升资源对象的管理效率；标签选择器：就是标签的查询过滤条件。

K8s的负载均衡器？
----------

答：负载均衡器是暴露服务的最常见和标准方式之一。

根据工作环境使用两种类型的负载均衡器，即内部负载均衡器或外部负载均衡器。内部负载均衡器自动平衡负载并使用所需配置分配容器，而外部负载均衡器将流量从外部负载引导至后端容器；

kubelet 监控 Node 节点资源使用是通过什么组件来实现的？
----------------------------------

答：用Metrics Server提供核心指标，包括Node、Pod的CPU和内存的使用。而Metrics Server需要采集node上的cAdvisor提供的数据资源，

当 kubelet 服务启动时，它会自动启动 cAdvisor 服务，然后 cAdvisor 会实时采集所在节点的性能指标及在节点上运行的容器的性能指标。

kubelet 的启动参数 --cadvisor-port 可自定义 cAdvisor 对外提供服务的端口号，默认是 4194；

Pod的状态？
-------

答：

1）Pending：已经创建了Pod，但是其内部还有容器没有创建；

2）Running：Pod内部的所有容器都已经创建，只有由一个容器还处于运行状态或者重启状态；

3）Succeeed：Pod内所有容器均已经成功执行并且退出，不会再重启；

4）Failed：Pod内所有容器都退出，但至少有一个为退出失败状态；

5）Unknown：由于某种原因不能获取该Pod的状态，可能是网络问题；

deployment/rs的区别？
-----------------

答：deployment是rs的超集，提供更多的部署功能，如：回滚、暂停和重启、 版本记录、事件和状态查看、滚动升级和替换升级。

如果能使用deployment，则不应再使用rc和rs；

rc/rs实现原理？
----------

答：

Replication Controller 可以保证Pod始终处于规定的副本数，

而当前推荐的做法是使用Deployment+ReplicaSet，

ReplicaSet 号称下一代的 Replication Controller，当前唯一区别是RS支持set-based selector，

RC是通过ReplicationManager监控RC和RC内Pod的状态，从而增删Pod，以实现维持特定副本数的功能，RS也是大致相同；

kubernetes服务发现？
---------------

答：

1）环境变量： 当你创建一个Pod的时候，kubelet会在该Pod中注入集群内所有Service的相关环境变量。**需要注意:** 要想一个Pod中注入某个Service的环境变量，则必须Service要先比该Pod创建；

2）DNS：可以通过cluster add-on方式轻松的创建KubeDNS来对集群内的Service进行服务发现；

k8s发布(暴露)服务，servcie的类型有那些？
--------------------------

答：

kubernetes原生的，一个Service的ServiceType决定了其发布服务的方式。

1） ClusterIP：这是k8s默认的ServiceType。通过集群内的ClusterIP在内部发布服务。

2）NodePort：这种方式是常用的，用来对集群外暴露Service，你可以通过访问集群内的每个NodeIP:NodePort的方式，访问到对应Service后端的Endpoint。

3）LoadBalancer: 这也是用来对集群外暴露服务的，不同的是这需要Cloud Provider的支持，比如AWS等。

4）ExternalName：这个也是在集群内发布服务用的，需要借助KubeDNS(version >= 1.7)的支持，就是用KubeDNS将该service和ExternalName做一个Map，KubeDNS返回一个CNAME记录；

简述ETCD及其特点?
-----------

答：etcd是一个分布式的、高可用的、一致的key-value存储数据库，基于Go语言实现，主要用于共享配置和服务发现。特点：

1）完全复制：集群中的每个节点都可以使用完整的存档；

2）高可用性：Etcd可用于避免硬件的单点故障或网络问题；

3）一致性：每次读取都会返回跨多主机的最新写入；

4）简单：包括一个定义良好、面向用户的API（gRPC）；

5）安全：实现了带有可选的客户端证书身份验证的自动化TLS；

6）快速：每秒10000次写入的基准速度；

7）可靠：使用Raft算法实现了强一致、高可用的服务存储目录；

简述ETCD适应的场景?
------------

答：

1）服务发现：服务发现要解决的也是分布式系统中最常见的问题之一，即在同一个分布式集群中的进程或服务，要如何才能找到对方并建立连接。本质上来说，服务发现就是想要了解集群中是否有进程在监听udp或tcp端口，并且通过名字就可以查找和连接。

2）消息发布与订阅：在分布式系统中，最实用对的一种组件间的通信方式：消息发布与订阅。构建一个配置共享中心，数据提供者在这个配置中心发布消息，而消息使用者订阅他们关心的主题，一旦主题有消息发布，就会实时通知订阅者。达成集中式管理与动态更新。应用中用到的一些配置信息放到etcd上进行集中管理。

3）负载均衡：分布式系统中，为了保证服务的高可用以及数据的一致性，通常都会把数据和服务部署多份，以此达到对等服务，即使其中的某一个服务失效了，也不影响使用。etcd本身分布式架构存储的信息访问支持负载均衡。

4）分布式通知与协调：通过注册与异步通知机制，实现分布式环境下不同系统之间的通知与协调，从而对数据变更做到实时处理。

5）分布式锁：因为etcd使用Raft算法保持了数据的强一致性，某次操作存储到集群中的值必然是全局一致的，所以很容易实现分布式锁。锁服务有两种使用方式，一是保持独占，二是控制时序。

6）分布式队列：分布式队列的常规用法与场景五中所描述的分布式锁的控制时序用法类似，即创建一个先进先出的队列，保证顺序。

7）集群监控与Leader精选：通过etcd来进行监控实现起来非常简单并且实时性强；

> 注：本文以 PDF 持续更新，最新尼恩 架构笔记、面试题 的PDF文件，请从下面的链接获取：[码云](https://gitee.com/crazymaker/SimpleCrayIM/blob/master/%E7%96%AF%E7%8B%82%E5%88%9B%E5%AE%A2%E5%9C%88%E6%80%BB%E7%9B%AE%E5%BD%95.md) 或者 [语雀](https://www.yuque.com/crazymakercircle/gkkw8s/khigna)

简述Kubernetes RC的机制?
-------------------

答：Replication Controller用来管理Pod的副本，保证集群中存在指定数量的Pod副本。当定义了RC并提交至Kubernetes集群中之后，Master节点上的Controller Manager组件获悉，并同时巡检系统中当前存活的目标Pod，并确保目标Pod实例的数量刚好等于此RC的期望值，若存在过多的Pod副本在运行，系统会停止一些Pod，反之则自动创建一些Pod；

简述kube-proxy作用?
---------------

答：kube-proxy 运行在所有节点上，它监听 apiserver 中 service 和 endpoint 的变化情况，创建路由规则以提供服务 IP 和负载均衡功能。

简单理解此进程是Service的透明代理兼负载均衡器，其核心功能是将到某个Service的访问请求转发到后端的多个Pod实例上；

简述kube-proxy iptables原理?
------------------------

答：Kubernetes从1.2版本开始，将iptables作为kube-proxy的默认模式。iptables模式下的kube-proxy不再起到Proxy的作用，其核心功能：通过API Server的Watch接口实时跟踪Service与Endpoint的变更信息，并更新对应的iptables规则，Client的请求流量则通过iptables的NAT机制“直接路由”到目标Pod；

简述kube-proxy ipvs原理?
--------------------

答：IPVS在Kubernetes1.11中升级为GA稳定版。

IPVS则专门用于高性能负载均衡，并使用更高效的数据结构（Hash表），允许几乎无限的规模扩张，因此被kube-proxy采纳为最新模式；

在IPVS模式下，使用iptables的扩展ipset，而不是直接调用iptables来生成规则链。

iptables规则链是一个线性的数据结构，ipset则引入了带索引的数据结构，因此当规则很多时，也可以很高效地查找和匹配；

可以将ipset简单理解为一个IP（段）的集合，这个集合的内容可以是IP地址、IP网段、端口等，iptables可以直接添加规则对这个“可变的集合”进行操作，这样做的好处在于可以大大减少iptables规则的数量，从而减少性能损耗；

简述kube-proxy ipvs和iptables的异同?
------------------------------

答：iptables与IPVS都是基于Netfilter实现的，但因为定位不同，二者有着本质的差别：

iptables是为防火墙而设计的；IPVS则专门用于高性能负载均衡，并使用更高效的数据结构（Hash表），允许几乎无限的规模扩张。

与iptables相比，IPVS拥有以下明显优势：为大型集群提供了更好的可扩展性和性能；支持比iptables更复杂的复制均衡算法（最小负载、最少连接、加权等）；支持服务器健康检查和连接重试等功能；可以动态修改ipset的集合，即使iptables的规则正在使用这个集合；

简述Kubernetes中什么是静态Pod?
----------------------

答：静态pod是由kubelet进行管理的仅存在于特定Node的Pod上，他们不能通过API Server进行管理，无法与ReplicationController、Deployment或者DaemonSet进行关联，并且kubelet无法对他们进行健康检查。

静态Pod总是由kubelet进行创建，并且总是在kubelet所在的Node上运行；

简述Kubernetes Pod的常见调度方式?
------------------------

答：

1）Deployment或RC：该调度策略主要功能就是自动部署一个容器应用的多份副本，以及持续监控副本的数量，在集群内始终维持用户指定的副本数量；

2）NodeSelector：定向调度，当需要手动指定将Pod调度到特定Node上，可以通过Node的标签（Label）和Pod的nodeSelector属性相匹配；

3）NodeAffinity亲和性调度：亲和性调度机制极大的扩展了Pod的调度能力，目前有两种节点亲和力表达：硬规则，必须满足指定的规则，调度器才可以调度Pod至Node上（类似nodeSelector，语法不同）；软规则，优先调度至满足的Node的节点，但不强求，多个优先级规则还可以设置权重值；

4）Taints和Tolerations（污点和容忍）：Taint：使Node拒绝特定Pod运行；Toleration：为Pod的属性，表示Pod能容忍（运行）标注了Taint的Node；

简述Kubernetes初始化容器（init container）?
----------------------------------

答：init container的运行方式与应用容器不同，它们必须先于应用容器执行完成，当设置了多个init container时，将按顺序逐个运行，并且只有前一个init container运行成功后才能运行后一个init container。

当所有init container都成功运行后，Kubernetes才会初始化Pod的各种信息，并开始创建和运行应用容器；

简述Kubernetes deployment升级过程?
----------------------------

答：

初始创建Deployment时，系统创建了一个ReplicaSet，并按用户的需求创建了对应数量的Pod副本；

当更新Deployment时，系统创建了一个新的ReplicaSet，并将其副本数量扩展到1，然后将旧ReplicaSet缩减为2；

之后，系统继续按照相同的更新策略对新旧两个ReplicaSet进行逐个调整；

最后，新的ReplicaSet运行了对应个新版本Pod副本，旧的ReplicaSet副本数量则缩减为0；

简述Kubernetes deployment升级策略?
----------------------------

答：

在Deployment的定义中，可以通过spec.strategy指定Pod更新的策略，

目前支持两种策略：Recreate（重建）和RollingUpdate（滚动更新），

默认值为RollingUpdate；

Recreate：设置spec.strategy.type=Recreate，表示Deployment在更新Pod时，会先杀掉所有正在运行的Pod，然后创建新的Pod；

RollingUpdate：设置spec.strategy.type=RollingUpdate，表示Deployment会以滚动更新的方式来逐个更新Pod。同时，可以通过设置spec.strategy.rollingUpdate下的两个参数（maxUnavailable和maxSurge）来控制滚动更新的过程；

简述Kubernetes DaemonSet类型的资源特性?
------------------------------

答：

DaemonSet资源对象会在每个Kubernetes集群中的节点上运行，并且每个节点只能运行一个pod，这是它和deployment资源对象的最大也是唯一的区别。

因此，在定义yaml文件中，不支持定义replicas。

它的一般使用场景如下：在去做每个节点的日志收集工作。监控每个节点的的运行状态。

简述Kubernetes自动扩容机制?
-------------------

答：

Kubernetes使用Horizontal Pod Autoscaler（HPA）的控制器实现基于CPU使用率进行自动Pod扩缩容的功能。

HPA控制器周期性地监测目标Pod的资源性能指标，并与HPA资源对象中的扩缩容条件进行对比，在满足条件时对Pod副本数量进行调整；

简述Kubernetes Service分发后端的策略?
----------------------------

答：

1）RoundRobin：默认为轮询模式，即轮询将请求转发到后端的各个Pod上；

2）SessionAffinity：基于客户端IP地址进行会话保持的模式，即第1次将某个客户端发起的请求转发到后端的某个Pod上，之后从相同的客户端发起的请求都将被转发到后端相同的Pod上；

简述Kubernetes Headless Service?
------------------------------

答：在某些应用场景中，若需要人为指定负载均衡器，不使用Service提供的默认负载均衡的功能，或者应用程序希望知道属于同组服务的其他实例。

Kubernetes提供了Headless Service来实现这种功能，即不为Service设置ClusterIP（入口IP地址），仅通过Label Selector将后端的Pod列表返回给调用的客户端；

简述Kubernetes外部如何访问集群内的服务?
-------------------------

答：

映射Pod到物理机：将Pod端口号映射到宿主机，即在Pod中采用hostPort方式，以使客户端应用能够通过物理机访问容器应用；

映射Service到物理机：将Service端口号映射到宿主机，即在Service中采用nodePort方式，以使客户端应用能够通过物理机访问容器应用；

映射Service到LoadBalancer：通过设置LoadBalancer映射到云服务商提供的LoadBalancer地址。这种用法仅用于在公有云服务提供商的云平台上设置Service的场景；

简述Kubernetes ingress?
---------------------

答：

K8s的Ingress资源对象，用于将不同URL的访问请求转发到后端不同的Service，以实现HTTP层的业务路由机制。

K8s使用了Ingress策略和Ingress Controller，两者结合并实现了一个完整的Ingress负载均衡器。

使用Ingress进行负载分发时，Ingress Controller基于Ingress规则将客户端请求直接转发到Service对应的后端Endpoint（Pod）上，从而跳过kube-proxy的转发功能，kube-proxy不再起作用，

全过程为：ingress controller + ingress 规则 ----> services；

简述Kubernetes镜像的下载策略?
--------------------

答：

1）Always：镜像标签为latest时，总是从指定的仓库中获取镜像；

2）Never：禁止从仓库中下载镜像，也就是说只能使用本地镜像；

3）IfNotPresent：仅当本地没有对应镜像时，才从目标仓库中下载；默认的镜像下载策略是：当镜像标签是latest时，默认策略是Always；当镜像标签是自定义时（也就是标签不是latest），那么默认策略是IfNotPresent；

简述Kubernetes的负载均衡器?
-------------------

答：

根据工作环境使用两种类型的负载均衡器，即内部负载均衡器或外部负载均衡器。

内部负载均衡器自动平衡负载并使用所需配置分配容器，而外部负载均衡器将流量从外部负载引导至后端容器；

简述Kubernetes各模块如何与API Server通信?
-------------------------------

答：K8s API Server作为集群的核心，负责集群各功能模块之间的通信。

集群内的各个功能模块通过API Server将信息存入etcd，当需要获取和操作这些数据时，则通过API Server提供的REST接口（用GET、LIST或WATCH方法）来实现，从而实现各模块之间的信息交互。

1）kubelet进程与API Server的交互：每个Node上的kubelet每隔一个时间周期，就会调用一次API Server的REST接口报告自身状态，API Server在接收到这些信息后，会将节点状态信息更新到etcd中；

2）kube-controller-manager进程与API Server的交互：kube-controller-manager中的Node Controller模块通过API Server提供的Watch接口实时监控Node的信息，并做相应处理

；3）kube-scheduler进程与API Server的交互：Scheduler通过API Server的Watch接口监听到新建Pod副本的信息后，会检索所有符合该Pod要求的Node列表，开始执行Pod调度逻辑，在调度成功后将Pod绑定到目标节点上；

简述Kubernetes Scheduler作用及实现原理?
------------------------------

答：

Scheduler是负责Pod调度的重要功能模块，负责接收Controller Manager创建的新Pod，为其调度至目标Node，调度完成后，目标Node上的kubelet服务进程接管后继工作，负责Pod接下来生命周期；

Scheduler的作用是将待调度的Pod，按照特定的调度算法和调度策略绑定（Binding）到集群中某个合适的Node上，并将绑定信息写入etcd中；

Scheduler通过调度算法调度为待调度Pod列表中的每个Pod从Node列表中选择一个最适合的Node来实现Pod的调度。随后，目标节点上的kubelet通过API Server监听到Kubernetes Scheduler产生的Pod绑定事件，然后获取对应的Pod清单，下载Image镜像并启动容器；

简述Kubernetes Scheduler使用哪两种算法将Pod绑定到worker节点?
---------------------------------------------

答：

1）预选（Predicates）：输入是所有节点，输出是满足预选条件的节点。kube-scheduler根据预选策略过滤掉不满足策略的Nodes。如果某节点的资源不足或者不满足预选策略的条件则无法通过预选；

2）优选（Priorities）：输入是预选阶段筛选出的节点，优选会根据优先策略为通过预选的Nodes进行打分排名，选择得分最高的Node。例如，资源越富裕、负载越小的Node可能具有越高的排名；

简述Kubernetes kubelet的作用?
------------------------

答：

在Kubernetes集群中，在每个Node（又称Worker）上都会启动一个kubelet服务进程。

该进程用于处理Master下发到本节点的任务，管理Pod及Pod中的容器。

每个kubelet进程都会在API Server上注册节点自身的信息，定期向Master汇报节点资源的使用情况，并通过cAdvisor监控容器和节点资源；

简述Kubernetes kubelet监控Worker节点资源是使用什么组件来实现的?
--------------------------------------------

答：

kubelet使用cAdvisor对worker节点资源进行监控。

在 Kubernetes 系统中，cAdvisor 已被默认集成到 kubelet 组件内，当 kubelet 服务启动时，它会自动启动 cAdvisor 服务，然后 cAdvisor 会实时采集所在节点的性能指标及在节点上运行的容器的性能指标；

简述Kubernetes如何保证集群的安全性?
-----------------------

答：

1）基础设施方面：保证容器与其所在宿主机的隔离；

2）用户权限：划分普通用户和管理员的角色；

3）API Server的认证授权：Kubernetes集群中所有资源的访问和变更都是通过Kubernetes API Server来实现的，因此需要建议采用更安全的HTTPS或Token来识别和认证客户端身份（Authentication），以及随后访问权限的授权（Authorization）环节；

4）API Server的授权管理：通过授权策略来决定一个API调用是否合法。对合法用户进行授权并且随后在用户访问时进行鉴权，建议采用更安全的RBAC方式来提升集群安全授权；

5）AdmissionControl（准入机制）：对kubernetes api的请求过程中，顺序为：先经过认证 & 授权，然后执行准入操作，最后对目标对象进行操作；

简述Kubernetes准入机制?
-----------------

答：

在对集群进行请求时，每个准入控制代码都按照一定顺序执行。

如果有一个准入控制拒绝了此次请求，那么整个请求的结果将会立即返回，并提示用户相应的error信息，准入控制（AdmissionControl）准入控制本质上为一段准入代码，在对kubernetes api的请求过程中，顺序为：先经过认证 & 授权，然后执行准入操作，最后对目标对象进行操作。

常用组件（控制代码）如下：

AlwaysAdmit：允许所有请求；

AlwaysDeny：禁止所有请求，多用于测试环境；

ServiceAccount：它将serviceAccounts实现了自动化，它会辅助serviceAccount做一些事情，比如如果pod没有serviceAccount属性，它会自动添加一个default，并确保pod的serviceAccount始终存在；

LimitRanger：观察所有的请求，确保没有违反已经定义好的约束条件，这些条件定义在namespace中LimitRange对象中；

NamespaceExists：观察所有的请求，如果请求尝试创建一个不存在的namespace，则这个请求被拒绝；

简述Kubernetes RBAC及其特点（优势）?
--------------------------

答：

RBAC是基于角色的访问控制，是一种基于个人用户的角色来管理对计算机或网络资源的访问的方法，

优势：

1）对集群中的资源和非资源权限均有完整的覆盖；

2）整个RBAC完全由几个API对象完成， 同其他API对象一样， 可以用kubectl或API进行操作；

3）可以在运行时进行调整，无须重新启动API Server；

简述Kubernetes Secret作用?
----------------------

答：

Secret对象，主要作用是保管私密数据，比如密码、OAuth Tokens、SSH Keys等信息。

将这些私密信息放在Secret对象中比直接放在Pod或Docker Image中更安全，也更便于使用和分发；

简述Kubernetes Secret有哪些使用方式?
---------------------------

答：

1）在创建Pod时，通过为Pod指定Service Account来自动使用该Secret；

2）通过挂载该Secret到Pod来使用它；

3）在Docker镜像下载时使用，通过指定Pod的spc.ImagePullSecrets来引用它；

简述Kubernetes PodSecurityPolicy机制?
---------------------------------

答：

Kubernetes PodSecurityPolicy是为了更精细地控制Pod对资源的使用方式以及提升安全策略。

在开启PodSecurityPolicy准入控制器后，Kubernetes默认不允许创建任何Pod，需要创建PodSecurityPolicy策略和相应的RBAC授权策略（Authorizing Policies），Pod才能创建成功；

简述Kubernetes PodSecurityPolicy机制能实现哪些安全策略?
------------------------------------------

1）特权模式：privileged是否允许Pod以特权模式运行；

2）宿主机资源：控制Pod对宿主机资源的控制，如hostPID：是否允许Pod共享宿主机的进程空间；

3）用户和组：设置运行容器的用户ID（范围）或组（范围）；

4）提升权限：AllowPrivilegeEscalation：设置容器内的子进程是否可以提升权限，通常在设置非root用户（MustRunAsNonRoot）时进行设置；

5）SELinux：进行SELinux的相关配置；

简述Kubernetes网络模型?
-----------------

答：Kubernetes网络模型中每个Pod都拥有一个独立的IP地址，不管它们是否运行在同一个Node（宿主机）中，都要求它们可以直接通过对方的IP进行访问；

同时为每个Pod都设置一个IP地址的模型使得同一个Pod内的不同容器会共享同一个网络命名空间，也就是同一个Linux网络协议栈。

这就意味着同一个Pod内的容器可以通过localhost来连接对方的端口；在Kubernetes的集群里，IP是以Pod为单位进行分配的。一个Pod内部的所有容器共享一个网络堆栈；

简述Kubernetes CNI模型?
-------------------

答：

Kubernetes CNI模型是对容器网络进行操作和配置的规范，通过插件的形式对CNI接口进行实现。

CNI仅关注在创建容器时分配网络资源，和在销毁容器时删除网络资源。

容器（Container）：是拥有独立Linux网络命名空间的环境，例如使用Docker或rkt创建的容器。容器需要拥有自己的Linux网络命名空间，这是加入网络的必要条件；

网络（Network）：表示可以互连的一组实体，这些实体拥有各自独立、唯一的IP地址，可以是容器、物理机或者其他网络设备（比如路由器）等；

简述Kubernetes网络策略?
-----------------

答：

为实现细粒度的容器间网络访问隔离策略，K8s引入Network Policy主要功能是对Pod间的网络通信进行限制和准入控制，设置允许访问或禁止访问的客户端Pod列表。

Network Policy定义网络策略，配合策略控制器（Policy Controller）进行策略的实现；

简述Kubernetes网络策略原理?
-------------------

答：

Network Policy的工作原理主要为：policy controller需要实现一个API Listener，监听用户设置的Network Policy定义，并将网络访问规则通过各Node的Agent进行实际设置（Agent则需要通过CNI网络插件实现）；

简述Kubernetes中flannel的作用?
------------------------

答：

1）它能协助Kubernetes，给每一个Node上的Docker容器都分配互相不冲突的IP地址；

2）它能在这些IP地址之间建立一个覆盖网络（Overlay Network），通过这个覆盖网络，将数据包原封不动地传递到目标容器内；

简述Kubernetes Calico网络组件实现原理?
----------------------------

答：

Calico是一个基于BGP的纯三层的网络方案，与OpenStack、Kubernetes、AWS、GCE等云平台都能够良好地集成，Calico在每个计算节点都利用Linux Kernel实现了一个高效的vRouter来负责数据转发。每个vRouter都通过BGP协议把在本节点上运行的容器的路由信息向整个Calico网络广播，并自动设置到达其他节点的路由转发规则；Calico保证所有容器之间的数据流量都是通过IP路由的方式完成互联互通的。

Calico节点组网时可以直接利用数据中心的网络结构（L2或者L3），不需要额外的NAT、隧道或者Overlay Network，没有额外的封包解包，能够节约CPU运算，提高网络效率；

简述Kubernetes共享存储的作用?
--------------------

答：

Kubernetes对于有状态的容器应用或者对数据需要持久化的应用，因此需要更加可靠的存储来保存应用产生的重要数据，以便容器应用在重建之后仍然可以使用之前的数据。因此需要使用共享存储；

简述Kubernetes PV和PVC?
--------------------

答：

PV是对底层网络共享存储的抽象，将共享存储定义为一种“资源”；

PVC则是用户对存储资源的一个“申请”；

简述Kubernetes PV生命周期内的阶段?
------------------------

答：

1）Available：可用状态，还未与某个PVC绑定；

2）Bound：已与某个PVC绑定；

3）Released：绑定的PVC已经删除，资源已释放，但没有被集群回收；

4）Failed：自动资源回收失败；

简述Kubernetes CSI模型?
-------------------

答：

CSI是Kubernetes推出与容器对接的存储接口标准，存储提供方只需要基于标准接口进行存储插件的实现，就能使用Kubernetes的原生存储机制为容器提供存储服务，CSI使得存储提供方的代码能和Kubernetes代码彻底解耦，部署也与Kubernetes核心组件分离；

CSI包括CSI Controller：的主要功能是提供存储服务视角对存储资源和存储卷进行管理和操作；Node的主要功能是对主机（Node）上的Volume进行管理和操作；

简述Kubernetes Worker节点加入集群的过程?
-----------------------------

答：在该Node上安装Docker、kubelet和kube-proxy服务； 然后配置kubelet和kubeproxy的启动参数，将Master URL指定为当前Kubernetes集群Master的地址，最后启动这些服务； 通过kubelet默认的自动注册机制，新的Worker将会自动加入现有的Kubernetes集群中； Kubernetes Master在接受了新Worker的注册之后，会自动将其纳入当前集群的调度范围；

简述Kubernetes Pod如何实现对节点的资源控制?
-----------------------------

答：

Kubernetes集群里的节点提供的资源主要是计算资源，计算资源是可计量的能被申请、分配和使用的基础资源。当前Kubernetes集群中的计算资源主要包括CPU、GPU及Memory。

CPU与Memory是被Pod使用的，因此在配置Pod时可以通过参数CPU Request及Memory Request为其中的每个容器指定所需使用的CPU与Memory量，Kubernetes会根据Request的值去查找有足够资源的Node来调度此Pod；

简述Kubernetes Requests和Limits如何影响Pod的调度?
---------------------------------------

答：

当一个Pod创建成功时，Kubernetes调度器（Scheduler）会为该Pod选择一个节点来执行。对于每种计算资源（CPU和Memory）而言，每个节点都有一个能用于运行Pod的最大容量值。调度器在调度时，首先要确保调度后该节点上所有Pod的CPU和内存的Requests总和，不超过该节点能提供给Pod使用的CPU和Memory的最大容量值；

简述Kubernetes Metric Service?
----------------------------

答：在Kubernetes从1.10版本后采用Metrics Server作为默认的性能数据采集和监控，主要用于提供核心指标（Core Metrics），包括Node、Pod的CPU和内存使用指标。

对其他自定义指标（Custom Metrics）的监控则由Prometheus等组件来完成；

简述Kubernetes中，如何使用EFK实现日志的统一管理？
-------------------------------

答：

在Kubernetes集群环境中，通常一个完整的应用或服务涉及组件过多，建议对日志系统进行集中化管理，EFK是 Elasticsearch、Fluentd 和 Kibana 的组合，

Elasticsearch：是一个搜索引擎，负责存储日志并提供查询接口；

Fluentd：负责从 Kubernetes 搜集日志，每个node节点上面的fluentd监控并收集该节点上面的系统日志，并将处理过后的日志信息发送给Elasticsearch；

Kibana：提供了一个 Web GUI，用户可以浏览和搜索存储在 Elasticsearch 中的日志；

简述Kubernetes如何进行优雅的节点关机维护?
--------------------------

答：由于Kubernetes节点运行大量Pod，因此在进行关机维护之前，建议先使用kubectl drain将该节点的Pod进行驱逐，然后进行关机维护；

简述Kubernetes集群联邦?
-----------------

答：Kubernetes集群联邦可以将多个Kubernetes集群作为一个集群进行管理。因此，可以在一个数据中心/云中创建多个Kubernetes集群，并使用集群联邦在一个地方控制/管理所有集群；

简述Helm及其优势?
-----------

答：Helm 是 Kubernetes 的软件包管理工具，Helm能够将一组K8S资源打包统一管理, 是查找、共享和使用为Kubernetes构建的软件的最佳方式。 Helm中通常每个包称为一个Chart，一个Chart是一个目录，优势：1）统一管理、配置和更新这些分散的 k8s 的应用资源文件；2）分发和复用一套应用模板；3）将应用的一系列资源当做一个软件包管理；4）对于应用发布者而言，可以通过 Helm 打包应用、管理应用依赖关系、管理应用版本并发布应用到软件仓库；5）对于使用者而言，使用 Helm 后不用需要编写复杂的应用部署文件，可以以简单的方式在 Kubernetes 上查找、安装、升级、回滚、卸载应用程序；

标签与标签选择器的作用是什么?
---------------

答：

1）标签可以附加在kubernetes任何资源对象之上的键值型数据，常用于标签选择器的匹配度检查，从而完成资源筛选；

2）标签选择器用于表达标签的查询条件或选择标准，Kubernetes API目前支持两个选择器：基于等值关系（equality-based）的标签选项器以及基于集合关系（set-based）的标签选择器；

什么是Google容器引擎?
--------------

答：Google Container Engine（GKE）是Docker容器和集群的开源管理平台。这个基于 Kubernetes的引擎仅支持在Google的公共云服务中运行的群集；

image的状态有那些？
------------

答：

1）Running：Pod所需的容器已经被成功调度到某个节点，且已经成功运行；

2）Pending：APIserver创建了pod资源对象，并且已经存入etcd中，但它尚未被调度完成或者仍然处于仓库中下载镜像的过程；

3）Unknown：APIserver无法正常获取到pod对象的状态，通常是其无法与所在工作节点的kubelet通信所致；

Service这种资源对象的作用是什么?
--------------------

答：

service就是将多个POD划分到同一个逻辑组中，并统一向外提供服务，POD是通过Label Selector加入到指定的service中。

Service相当于是一个负载均衡器，用户请求会先到达service，再由service转发到它内部的某个POD上，通过 services.spec.type 字段来指定：

1）ClusterIP：用于集群内部访问。该类型会为service分配一个IP，集群内部请求先到达service，再由service转发到其内部的某个POD上；

2）NodePort：用于集群外部访问。该类型会将Service的Port映射到集群的每个Node节点上，然后在集群之外，就能通过Node节点上的映射端口访问到这个Service；

3）LoadBalancer：用于集群外部访问。该类型是在所有Node节点前又挂了一个负载均衡器，作为集群外部访问的统一入口，外部流量会先到达LoadBalancer，再由它转发到集群的node节点上，通过nodePort再转发给对应的service，最后由service转发到后端Pod中；

4）ExternalName：创建一个DNS别名（即CNAME）并指向到某个Service Name上，也就是为某个Service Name添加一条CNAME记录，当有请求访问这个CNAME时会自动解析到这个Service Name上；

常用的标签分类有哪些?
-----------

答：release（版本）：stable（稳定版）、canary（金丝雀版本）、beta（测试版本）、environment（环境变量）：dev（开发）、qa（测试）、production（生产）、application（应用）：ui、as（application software应用软件）、pc、sc、tier（架构层级）：frontend（前端）、backend（后端）、cache（缓存）、partition（分区）：customerA（客户A）、customerB（客户B）、track（品控级别）：daily（每天）、weekly（每周）；

说说你对Job这种资源对象的了解?
-----------------

答：

Job控制一组Pod容器，可以通过Job这种资源对象定义并启动一个批处理任务的Job，其中Job所控制的Pod副本是短暂运行的，可以将其视为一组Docker容器，每个Docker容器都仅仅运行一次，当Job控制的所有Pod的副本都运行结束时，对应的Job也就结来。

Job生成的副本是不能自动重启的，对应的Pod副本的RestartPolicy都被设置为Never。

Job所控制的Pod副本的工作模式能够多实例并行计算。

k8s是怎么进行服务注册的?
--------------

答：

1）Service创建的时候会向 API Server 用 POST 方式提交一个新的 Service 定义，这个请求需要经过认证、鉴权以及其它的准入策略检查过程之后才会放行；

2）CoreDns 会为Service创建一个dns记录，Service 得到一个 ClusterIP（虚拟 IP 地址），并保存到集群数据仓库；

3）在集群范围内传播 Service 配置；

Kubernetes与Docker Swarm的区别如何?
-----------------------------

答：

1）安装和部署：k8s安装很复杂;但是一旦安装完毕，集群就非常强大，Docker Swarm安装非常简单;但是集群不是很强大;2)图形用户界面：k8s有，Docker Swarm无；

3）可伸缩性：k8s支持，Docker Swarm比k8s快5倍；

4）自动伸缩：k8s有，Docker Swarm无；

5）负载均衡：k8s在不同的Pods中的不同容器之间平衡负载流量，需要手动干预，Docker Swarm可以自动平衡集群中容器之间的流量；

6）滚动更新回滚：k8s支持，Docker Swarm可以部署滚动更新，但不能自动回滚；

7）数据量：k8s可以共享存储卷。只能与其他集装箱在同一Pod，Docker Swarm可以与任何其他容器共享存储卷；

8）日志记录和监控：k8s内置的日志和监控工具，Docker Swarm要用第三方工具进行日志记录和监控；

什么是Container Orchestration?
---------------------------

答：

1）资源编排 - 负责资源的分配，如限制 namespace 的可用资源，scheduler 针对资源的不同调度策略；

2）工作负载编排 - 负责在资源之间共享工作负载，如 Kubernetes 通过不同的 controller 将 Pod 调度到合适的 node 上，并且负责管理它们的生命周期；

3）服务编排 - 负责服务发现和高可用等，如 Kubernetes 中可用通过 Service 来对内暴露服务，通过 Ingress 来对外暴露服务；容器编排常用的控制器有：Deployment 经常被作为无状态实例控制器使用; StatefulSet 是一个有状态实例控制器; DaemonSet 可以指定在选定的 Node 上跑，每个 Node 上会跑一个副本，它有一个特点是它的 Pod 的调度不经过调度器，在 Pod 创建的时候就直接绑定 NodeName；最后一个是定时任务，它是一个上级控制器，和 Deployment 有些类似，当一个定时任务触发的时候，它会去创建一个 Job ，具体的任务实际上是由 Job 来负责执行的；

什么是Heapster?
------------

答：

Heapster 是 K8s 原生的集群监控方案。

Heapster 以 Pod 的形式运行，它会自动发现集群节点、从节点上的 Kubelet 获取监控数据。Kubelet 则是从节点上的 cAdvisor 收集数据；

k8s Architecture的不同组件有哪些?
-------------------------

答：

主要有两个组件 – 主节点和工作节点。

主节点具有kube-controller-manager，kube-apiserver，kube-scheduler等组件。

而工作节点具有kubelet和kube-proxy等组件；

能否介绍一下Kubernetes中主节点的工作情况?
--------------------------

答：

主节点是集群控制节点，负责集群管理和控制，包含：

1）apiserver: rest接口，资源增删改查入口；

2）controller-manager:所有资源对象的控制中心；

3）scheduler:负责资源调度，例如pod调度；

4）etcd: 保存资源对象数据；

kube-apiserver和kube-scheduler的作用是什么？
------------------------------------

答：

kube-apiserver: rest接口，增删改查接口，集群内模块通信；

kube-scheduler: 将待调度的pod按照调度算法绑定到合适的pod，并将绑定信息写入etcd；

你能简要介绍一下Kubernetes控制管理器吗？
-------------------------

Kubernetes控制管理器是集群内部的控制中心，负责node,pod,namespace等管理，

控制管理器负责管理各种控制器，每个控制器通过api server监控资源对象状态，将现有状态修正到期望状态；

Kubernetes有哪些不同类型的服务？
---------------------

*   ClusterIP、
*   NodePort、
*   LoadBalancer、
*   ExternalName；

你对Kubernetes的负载均衡器有什么了解？
------------------------

答：

1）内部负载均衡器: 自动平衡负载并使用所需配置分配容器；

2）外部负载均衡器: 将流量从外部负载引导至后端容器；

使用Kubernetes时可以采取哪些最佳安全措施?
--------------------------

1）确保容器本身安全；

2）锁定容器的Linux内核；

3）使用基于角色的访问控制（RBAC）；

4）保守秘密的辛勤工作；5）保持网络安全；

> 注：本文以 PDF 持续更新，最新尼恩 架构笔记、面试题 的PDF文件，请从下面的链接获取：[码云](https://gitee.com/crazymaker/SimpleCrayIM/blob/master/%E7%96%AF%E7%8B%82%E5%88%9B%E5%AE%A2%E5%9C%88%E6%80%BB%E7%9B%AE%E5%BD%95.md) 或者 [语雀](https://www.yuque.com/crazymakercircle/gkkw8s/khigna)

参考文献：
-----

[https://blog.csdn.net/qq\_21222149/article/details/89201744](https://blog.csdn.net/qq_21222149/article/details/89201744)

[https://blog.csdn.net/warrior\_0319/article/details/80073720](https://blog.csdn.net/warrior_0319/article/details/80073720)

[http://www.sel.zju.edu.cn/?p=840](http://www.sel.zju.edu.cn/?p=840)

[http://alexander.holbreich.org/docker-components-explained/](http://alexander.holbreich.org/docker-components-explained/)

[https://www.cnblogs.com/sparkdev/p/9129334.htmls](https://www.cnblogs.com/sparkdev/p/9129334.htmls)

推荐阅读：
-----

*   《[Docker面试题（史上最全 + 持续更新）](https://blog.csdn.net/crazymakercircle/article/details/128670335)》
    
*   《 [场景题：假设10W人突访，你的系统如何做到不 雪崩？](https://blog.csdn.net/crazymakercircle/article/details/128533821)》
    
*   《[尼恩Java面试宝典](https://blog.csdn.net/crazymakercircle/article/details/124790425)》
    
*   《[Springcloud gateway 底层原理、核心实战 (史上最全)](https://blog.csdn.net/crazymakercircle/article/details/125057567)》
    
*   《[Flux、Mono、Reactor 实战（史上最全）](https://blog.csdn.net/crazymakercircle/article/details/124120506)》
    
*   《[sentinel （史上最全）](https://blog.csdn.net/crazymakercircle/article/details/125059491)》
    
*   《[Nacos (史上最全)](https://blog.csdn.net/crazymakercircle/article/details/125057545)》
    
*   《[分库分表 Sharding-JDBC 底层原理、核心实战（史上最全）](https://blog.csdn.net/crazymakercircle/article/details/123420859)》
    
*   《[TCP协议详解 (史上最全)](https://blog.csdn.net/crazymakercircle/article/details/114527369)》
    
*   《[clickhouse 超底层原理 + 高可用实操 （史上最全）](https://blog.csdn.net/crazymakercircle/article/details/126992542)》
    
*   《[nacos高可用（图解+秒懂+史上最全）](https://blog.csdn.net/crazymakercircle/article/details/120702536)》
    
*   《[队列之王： Disruptor 原理、架构、源码 一文穿透](https://blog.csdn.net/crazymakercircle/article/details/128264803)》
    
*   《[环形队列、 条带环形队列 Striped-RingBuffer （史上最全）](https://blog.csdn.net/crazymakercircle/article/details/128264508)》
    
*   《[一文搞定：SpringBoot、SLF4j、Log4j、Logback、Netty之间混乱关系（史上最全）](https://blog.csdn.net/crazymakercircle/article/details/125135726)
    
*   《[单例模式（史上最全）](https://blog.csdn.net/crazymakercircle/article/details/128265067)
    
*   《[红黑树（ 图解 + 秒懂 + 史上最全）](https://blog.csdn.net/crazymakercircle/article/details/125017316)》
    
*   《[分布式事务 （秒懂）](https://blog.csdn.net/crazymakercircle/article/details/109459593)》
    
*   《[缓存之王：Caffeine 源码、架构、原理（史上最全，10W字 超级长文）](https://blog.csdn.net/crazymakercircle/article/details/128123114)》
    
*   《[缓存之王：Caffeine 的使用（史上最全）](https://blog.csdn.net/crazymakercircle/article/details/113751575)》
    
*   《[Java Agent 探针、字节码增强 ByteBuddy（史上最全）](https://blog.csdn.net/crazymakercircle/article/details/126579528)》
    
*   《[Docker原理（图解+秒懂+史上最全）](https://blog.csdn.net/crazymakercircle/article/details/120747767)》
    
*   《[Redis分布式锁（图解 - 秒懂 - 史上最全）](https://blog.csdn.net/crazymakercircle/article/details/116425814)》
    
*   《[Zookeeper 分布式锁 - 图解 - 秒懂](https://blog.csdn.net/crazymakercircle/article/details/85956246)》
    
*   《[Zookeeper Curator 事件监听 - 10分钟看懂](https://blog.csdn.net/crazymakercircle/article/details/85922561)》
    
*   《[Netty 粘包 拆包 | 史上最全解读](https://blog.csdn.net/crazymakercircle/article/details/83957259)》
    
*   《[Netty 100万级高并发服务器配置](https://blog.csdn.net/crazymakercircle/article/details/83758107)》
    
*   《[Springcloud 高并发 配置 （一文全懂）](https://blog.csdn.net/crazymakercircle/article/details/102557988)》

本文转自 <https://www.cnblogs.com/crazymakercircle/p/17052058.html>，如有侵权，请联系删除。