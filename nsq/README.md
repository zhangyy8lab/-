# Nsq队列

> 本节代码为oscro平台摘取部分代码



## Producer代码逻辑

- 主体接口接收创建namespace创建请求
  - 解析结构体
  - 根据参数， 获取workspace与ClusterNode绑定关系（记录资源配额信息）
  - 创建namespace数据库记录（状态为creating）
  - 创建namespace.set及相关数据库记录
  - 发消息

## Consumer代码逻辑

- 异步  goroutine 持续监控消息

  go namespace.NsCreate()

- 获取消息进行处理

  - nsqLookupD 实例化， for 循环 解析message.Body

  - 获取create.namespace.Obj 

  - 实例化 Meta 结构体

  - 实例化 http.client, 需要二次请求clusterNodeApi， 在clusterNodeServer 创建namespace

  -  处理响应结果， 同步namespace状态

    
