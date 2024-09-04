# README #

这是基于微服务架构设计的 OSCRO 平台其中一环 `用户管理` 服务

### 服务介绍

- 基于consul组件来进行注册和发现
- 实现平台用户的管理
- 所有接口由 OSCRO 整体服务内部调用(mainServer的Proxy方法)。调用前的请求会经由统一入口用户校验
- 调整用户与角色， 该角色为平台角色。 
- 4类角色分别是， 具体角色工作内容与范围或权限分配请参考用户管理设计文档: 
  - 一级管理员即平台最大权限
  - 二级管理员即workspace管理员
  - 普通用户
  - 未授权用户

### 镜像构建

- 基于开发人员编写的Dockerfile文件
- 可基于平台的pipeline功能 或 人员进行构建

### 服务启动

- 基于k8s以应用启动的方式进行启动，启动前需要将配置文件进行挂载进去，灵活配置管理
- 详情参数对应yaml文件


### 如何使用

- 平台配合web页面进行使用，由web页面请求相关接口及数据入参

### 配置说明

[service]
Mode = debug # 环境使用模式
Port = 4002  # 服务端口
Address = 10.1.1.116 # 服务ip
Name = account-dev # 服务名，需要注册到consul中的

[mysql]  # 数据库数据
Db = mysql
DbAddress = 192.168.1.178
DbPort = 3306
DbUser = 8lab
DbPassWord = 8lab
DbName = osCro23
Charset = utf8mb4

[consul] # consul节点信息
Address = 10.1.1.116
Port = 8500
Token = 6e34993c-16e9-c780-9095-377e5004980e
NodeSide = 1


[log]  # 日志存放路径
Path = /var/log/8lab/auth

[key]  # mainServer 调用accountServer时需要， 校验请求是否全法
AccountKey = 4dxdL1bOclQ63AUxGkB2ei1n2FwUTbQ3gEv

[auth] # authServe名称， 根据这个名称去consul中获取服务使用的ipPort信息
AuthServerName = authServer

[oauth2] # oauth2的服务地址，这个服务没有注册到consul
ServerHost = https://192.168.1.178:3000