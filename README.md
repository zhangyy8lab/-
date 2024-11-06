## 个人信息

- 姓名：张洋洋

- 地址：https://github.com/zhangyy8lab/docs.git

  ​           https://blog.csdn.net/weixin_38507813?spm=1011.2124.3001.5343

- 电话：18510085710

- 邮箱：0xA2618@gmail.com

## 教育背景

- 石家庄信息工程职业学院｜计算机网络  |  专科   |   2011.09  -  2014.06
- 中国石油大学                   ｜计算机应用 ｜  本科  ｜ 2019.09  -  2022.07

## 技能

- 编程语言： Go, Python
- 容器化技术：kubernetes, Docker, docker-ompose
- 监控组件： prometheus, grafana, alertManage, dingtalk
- 消息队列：NSQ, RocketMQ
- Web框架：Gin, Django REST Framework, fastApi, flask
- ORM：GORM, djangoOrm
- 微服务注册与发现：Consul
- 反向代理服务器：Nginx
- 数据库：MySQL, Redis
- Web3: solana, solana-explorer, blockscout

## 工作经验

### 北京八分量信息科技有限公司 | devops | 2021.07 至今

#### Tusima / 运维 2024.01至今

- 基于polygon-cdk 实现的 区块链网络二层服务

  - 目前正在学习solana相关技术

- 配置和管理 Nginx 反向代理服务器，实现请求的负载均衡和反向代理。

- 使用 Docker 和 Docker Compose 部署和管理开发和测试环境，简化开发环境搭建和应用程序的部署过程。

- 通过 prometheus、grafana、alertManager、dingtalk 、二次prome-cli 组合方式实现主机及容器服务实时告警监控

  

#### 云原生平台后端开发/运维部 2021.07-2024.01

- 使用 Python 和 Go 开发和维护后端服务，采用 Gin 框架和 ORM/GORM  进行快速开发和数据库交互
- 采用Oauth2授权协议，实现平台登录
- 通过 Consul 实现微服务的注册和发现，确保服务的可用性和负载均衡
- 集成 NSQ 和 RocketMQ 实现高吞吐量的消息传递系统，确保可靠性和消息传递的有序性。
- 使用 MySQL 和 Redis 设计和优化关键业务功能的数据库结构和查询性能。
- K8sCluster部署和管理生产环境的容器化应用程序，实现高可用性、水平扩展和自动化部署。
- ci/cd自动化 pipeline 进行服务版本更新

### 亚信科技(成都)有限公司

#### 中国移动4A项目组 / 运维开发 2019.11 至 2021.07

-  日常系统版本更新/升级
- 使用Python和ansible进行主机管理与维护
- 使用crontab+sync实现文件增量备份，每日/周/月数据量统计并生成html格式，定时发邮件给甲方
- 集中化5.0项目微服务开发及部署，接口联调(4省份1中心3前置节点)

### 项目

- oscro云原生平台开发

  https://8labteam.github.io/OSCRO_DOCS/ 

  - 担任后端架设设计及后端开发

  - 平台架构分为两层 

    - 1).是用户交互层，即web对应的后端api接口； 
    - 2). resource层（基于k8s-client), 调用k8s底层接口

  - web交互层: 实现开发者 <==> dapp(应用) 状态管理；应用编排 application/svc/sts/depl/cm 等各类组件yaml编排及启动状态

  - resource层: 实现用户(dev/test/ops)在namespace下执行组件的接口调用，ci/cd执行，日志输出记录(clickhouse记录)

    

- Tusima zkevm-layer2

  - 使用 docker-compose 进行各服务组件更新与测试

  - 区块链浏览器部署； 优化高并发合约交易时后端服务负载过高问题

    

### 自我评价

- 具备扎实的 Linux 系统管理能力，熟悉常用的 Linux 发行版和命令行工具。
- 在 MySQL 和 Redis 数据库方面有丰富的实际经验，包括性能优化和故障排查。
- 熟悉消息队列的使用和配置，能够处理高吞吐量的消息传递场景。
- 熟练编写 Python 和 Go 语言的后端服务，熟悉 Gin 框架和 GORM ORM。
- 具备使用 Consul 实现微服务注册和发现的经验，保证了服务的可用性和弹性扩展。
- 精通 Nginx 反向代理服务器的配置和管理，了解负载均衡和反向代理的原理。
- 精通 Docker 和 Docker Compose 的使用，能够快速构建和管理容器化应用。
-  Kubernetes (K8s) 的基本概念和使用，具备在 K8s 上部署和管理容器化应用的能力。
- 具备良好的团队合作和沟通能力，能够适应快节奏的工作环境。

---