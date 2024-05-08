## 组件信息
```go
prometheus 负责采集数据, 通过rules 进行监控项进行判断
grafana   分析数据并展示
alertmanager  推送告警消息
webhook-dingtalk 告警消息到钉钉


// prometheus 官方参考。 
// https://www.prometheus.wang/alert/alert-manager-extension-with-webhook.html

// 告警规则
https://yunlzheng.gitbook.io/prometheus-book/parti-prometheus-ji-chu/alert/prometheus-alert-rule

// dashbord
https://grafana.com/grafana/dashboards/
```
## 架构
![image.png](https://cdn.nlark.com/yuque/0/2024/png/32532564/1707126450286-d43e0473-ebd6-48dd-a366-32e1e8118fea.png#averageHue=%23f6f4f3&clientId=u0dd1994e-2766-4&from=paste&height=594&id=u7293ce90&originHeight=594&originWidth=1043&originalType=binary&ratio=1&rotation=0&showTitle=false&size=101880&status=done&style=none&taskId=u1806e460-8d73-440d-b1d6-9dde1569076&title=&width=1043)
```go
// 参考连接
https://blog.csdn.net/crazymakercircle/article/details/127206293
```
## prometheus
### 配置
```yaml
// prometheus/prometheus.yml 
// 具体路径看 启动


global:
  scrape_interval:     5s    # 多久 收集 一次数据
  evaluation_interval: 5s    # 多久 评估 一次规则
  scrape_timeout:      4s    # 每次 收集数据的 超时时间

# 收集数据 配置 列表
scrape_configs:
  - job_name: prometheus            # 必须配置, 自动附加的job labels, 必须唯一
    static_configs:
      - targets: ['10.1.1.116:9090']       # 指定prometheus ip端口
        labels:
          instance: prometheus                 #标签

  - job_name: nodeExport    			// 每一个Job 都是一个监控项
    static_configs:								// 固定格式
      - targets: ['10.1.1.116:9010'].   // 这个地址是 explorer Node 的主机地址
        labels:   					 	   // labels下的东西为后面引用的时候用这些标签可根据自身情况进行配置
          instance: "10.1.1.116:9010" 
          location: "本地"


alerting:                         #Alertmanager相关的配置
  alertmanagers:
  - static_configs:
    - targets:
      - 10.1.1.116:9093         #指定告警模块

rule_files:                      #告警规则文件, 可以使用通配符 
  - rules/*.yml
```

```yaml
// prometheus/rules/alert.yml
// 具体路径看 启动
// 参考连接 rules https://www.cnblogs.com/namedgx/p/14919857.html
// 参考连接 rules https://www.cnblogs.com/blogof-fusu/p/17161554.html


groups:
  - name: prometheus-alert
    rules:
    - alert: instance-down
      expr: up == 0
      for: 5s
      labels:
        severity: P1
        team: "一条大河"
        instance: "{{ $labels.instance }}"
        location: "{{ $labels.location }}"
      annotations:
        title: "instance-down"
        description: "\n- instance:{{ $labels.instance }} down"


----------------------------------- 参考及说明 -----------------------------------------------------


# 相关的规则设置定义在一个group下。在每一个group中我们可以定义多个告警规则(rule)
groups:
  # 组名。报警规则组名称
  - name: 内存预警
    rules:
    - alert: 内存使用率预警
      # expr：基于PromQL表达式告警触发条件，用于计算是否有时间序列满足该条件。
      expr: (node_memory_MemTotal_bytes - (node_memory_MemFree_bytes+node_memory_Buffers_bytes+node_memory_Cached_bytes )) / node_memory_MemTotal_bytes * 100 > 98
      # for：评估等待时间，可选参数。用于表示只有当触发条件持续一段时间后才发送告警。在等待期间新产生告警的状态为pending。
      for: 1m # for语句会使 Prometheus 服务等待指定的时间, 然后执行查询表达式。（for 表示告警持续的时长，若持续时长小于该时间就不发给alertmanager了，大于该时间再发。for的值不要小于prometheus中的scrape_interval，例如scrape_interval为30s，for为15s，如果触发告警规则，则再经过for时长后也一定会告警，这是因为最新的度量指标还没有拉取，在15s时仍会用原来值进行计算。另外，要注意的是只有在第一次触发告警时才会等待(for)时长。）
      # labels：自定义标签，允许用户指定要附加到告警上的一组附加标签。
      labels:
        # severity: 指定告警级别。有三种等级，分别为 warning, critical 和 emergency 。严重等级依次递增。
        severity: critical
      # annotations: 附加信息，比如用于描述告警详细信息的文字等，annotations的内容在告警产生时会一同作为参数发送到Alertmanager。
      annotations:
        title: "内存使用率预警"
        serviceName: "{{ $labels.serviceName }}"
        instance: "{{ $labels.instance }}"
        value: "{{ $value }}"
        btn: "点击查看详情 :玫瑰:"
        link: "http://192.168.3.26:3000/grafana/d/aka/duo-job-ji-cheng-fu-wu-qi-jian-kong"
        template: "**${serviceName}**(${instance}) 内存使用率已经超过阈值 **98%**, 请及时处理！\n当前值: ${value}%"
  
  - name: 磁盘预警
    rules:
    - alert: 磁盘使用率预警
      expr: (node_filesystem_size_bytes - node_filesystem_avail_bytes) / node_filesystem_size_bytes * 100 > 90
  
      for: 1m
  
      labels:
        severity: critical
  
      annotations:
        title: "磁盘使用率预警"
        serviceName: "{{ $labels.serviceName }}"
        instance: "{{ $labels.instance }}"
        mountpoint: "{{ $labels.mountpoint }}"
        value: "{{ $value }}"
        btn: "点击查看详情 :玫瑰:"
        link: "http://192.168.3.26:3000/grafana/d/aka/duo-job-ji-cheng-fu-wu-qi-jian-kong"
        template: "**${serviceName}**(${instance}) 服务器磁盘设备使用率超过 **90%**, 请及时处理！\n挂载点: ${mountpoint}\n当前值: ${value}%!"
  
  - name: 实例存活报警
    rules:
    - alert: 实例存活报警
      expr: up == 0
      for: 30s
  
      labels:
        severity: emergency
  
      annotations:
        title: "节点宕机报警"
        serviceName: "{{ $labels.serviceName }}"
        instance: "{{ $labels.instance }}"
        btn: "点击查看详情 :玫瑰:"
        link: "http://192.168.3.26:9090/targets"
        template: "节点 **${serviceName}**(${instance}) 断联, 请及时处理!"
```

### rules cpu_mem_disk
```yaml
groups:
  # 组名。报警规则组名称
  - name: 内存预警
    rules:
    - alert: 内存使用率预警
      expr: (node_memory_MemTotal_bytes - (node_memory_MemFree_bytes+node_memory_Buffers_bytes+node_memory_Cached_bytes )) / node_memory_MemTotal_bytes * 100 > 90
      for: 1m
      labels:
        team: "tusima"
        severity: P2
        instance: "{{ $labels.instance }}"
        location: "{{ $labels.location }}"

      annotations:
        title: "内存使用率超出90%"
        description: "Memory usage exceeds 90%, current:{{ printf \"%.2f\" $value }}%"

  - name: 磁盘预警
    rules:
    - alert: 磁盘使用率预警
      expr: (node_filesystem_size_bytes - node_filesystem_avail_bytes) / node_filesystem_size_bytes * 100 > 30

      for: 1m

      labels:
        team: "tusima"
        severity: P2
        instance: "{{ $labels.instance }}"
        location: "{{ $labels.location }}"

      annotations:
        title: "磁盘使用率预警"
        description: "disk usage exceeds 90%, current:{{ printf \"%.2f\" $value }}%"
        detail: "\n - device:{{ $labels.device }}
        \n - fstype:{{ $labels.fstype }}
        \n - mountport: {{ $labels.mountpoint }}"

  - name: 实例存活报警
    rules:
    - alert: 实例存活报警
      expr: up == 0
      for: 30s

      labels:
        severity: P1
        team: "tusima"
        instance: "{{ $labels.instance }}"
        location: "{{ $labels.location }}"

      annotations:
        title: "节点宕机报警"
        instance: "{{ $labels.instance }}"
        description: "节点 **({{ $labels.instance}})** 断联, 请及时处理!"
```
### 启动
```shell
docker run -d \
-p 9090:9090 \
-v /root/prometheus/prometheus/prometheus/:/etc/prometheus/
--name prometheus prom/prometheus:latest \
--restart always
```

## NodeExport
### 启动
```shell
# 理解为prometheus 获取数据client
docker run -d --name node-exporter -p 9010:9100 --restart=always \
-h "nodeExport" \
-v "/proc:/host/proc:ro" \
-v "/sys:/host/sys:ro" \
-v "/:/rootfs:ro" \
prom/node-exporter
```

## Grafana
### 配置 && 启动
```go
docker run -d --name grafana11 \
-p 3000:3000 \
grafana/grafana:latest

---------
docker cp grafana11:/etc/grafana ./configs
docker cp grafana11:/var/lib/grafana ./data

---------------

docker run -d --name grafana \
-p 3000:3000 \
-v /xxx/configs/grafana/:/etc/grafana/ \
-v /xxx/volumes/grafana/:/var/lib/grafana/ \
grafana/grafana:latest
```

### Web 相关配置

-  首次登录时 账号和密码都是admin， 登录成功后需要修改密码 
-  添加数据源  Home -> Connections -> DataSources -> Prometheus 
-  Dashboard  
   - Home -> Dashboard -> New -> Import  (ID：8919)
   - 参考连接 [https://grafana.com/grafana/dashboards/](https://grafana.com/grafana/dashboards/)  
-  免密访问， 需要创建组织 
   -  Home -> Administration -> General -> Organizations +New org  名称随意起 暂叫： vister org 
   -  调整配置文件并启动服务 
```shell
 /xxx/configs/grafana.ini

#################################### Anonymous Auth ######################
[auth.anonymous]
# enable anonymous access
enabled = true

# specify organization name that should be used for unauthenticated users
org_name = vister org 

# specify role for unauthenticated users
org_role = Viewer
```
 

## AlertManager
### 配置
```yaml
// xxx/alerymanger.yml

global:
  # 每5分钟检查一次是否恢复
  resolve_timeout: 15s # 5m
# route用来设置报警的分发策略
route:
  # 采用哪个标签来作为分组依据
  group_by: ['alertname']
  # 组告警等待时间。也就是告警产生后等待30s，如果有同组告警一起发出
  group_wait: 10s # 30s
  # 两组告警的间隔时间
  group_interval: 10s # 30s
  # 重复告警的间隔时间，减少相同告警的发送频率
  repeat_interval: 10s # 1h
  # 设置默认接收人
  receiver: 'webhook'
receivers:
- name: 'webhook'
  webhook_configs:
  - url: 'http://172.19.23.177:9026/dingtalk/webhook1/send'
    send_resolved: true
----------------------------- 参考  -------------------------------------------

https://www.cnblogs.com/punchlinux/p/17035742.html
  #group_wait: 10s #第一次产生告警，等待 10s，组内有告警就一起发出，没有其它告警就单独发出。
  #group_interval: 2m #第二次产生告警，先等待 2 分钟，2 分钟后还没有恢复就进入 repeat_interval。
  #repeat_interval: 5m #在最终发送消息前再等待 5 分钟，5 分钟后还没有恢复就发送第二次告警。
  
  ----------------------------- 参考  -------------------------------------------
// 具体参考连接 https://yunlzheng.gitbook.io/prometheus-book/parti-prometheus-ji-chu/alert/alert-manager-route


route: # 根据标签匹配，确定当前告警应该如何处理；
  group_by: ['alertname'] # 告警应该根据那些标签进行分组，不分组可以指定 ...
  group_wait: 30s # 同一组的告警发出前要等待多少秒，这个是为了把更多的告警一个批次发出去
  group_interval: 5m # 同一组的多批次告警间隔多少秒后，才能发出
  repeat_interval: 1h # 重复的告警要等待多久后才能再次发出去
  receiver: 'webhook'
  routes:
  - receiver: webhook
    group_wait: 10s
    match:
      alertname: alertname

receivers: # 接收人是一个抽象的概念，它可以是一个邮箱也可以是微信，Slack或者Webhook等，接收人一般配合告警路由使用；
- name: 'webhook'
  webhook_configs:
  - url: http://10.1.1.116:8060/dingtalk/webhook1/send
    send_resolved: true

inhibit_rules: # 合理设置抑制规则可以减少垃圾告警的产生 比如说当我们的主机挂了，可能引起主机上的服务，数据库，中间件等一些告警，假如说后续的这些告警相对来说没有意义，我们可以用抑制项这个功能，让PrometheUS只发出主机挂了的告警。
  - source_match: 根据label匹配源告警
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'dev', 'instance'] # 处的集合的label，在源和目的里的值必须相等。如果该集合的内的值再源和目的里都没有，那么目的告警也会被抑制。
```

### 启动

```shell
docker run --name alertmanager -d -p 9093:9093 --restart=always \
prom/alertmanager

-------
docker cp alertmanager:/etc/alertmanager /xxx/configs


------- 
docker run --name alertmanager -d -p 9093:9093 --restart=always \
-v /xxx/configs/alertmanager/:/etc/alertmanager/ \
prom/alertmanager
```

## dingTalk

### 配置 Config

**vim config.yaml**
```yaml
templates:
  - /opt/dingtalk/template/dingtalk.tmpl
targets:
  webhook1:
    url: https://oapi.dingtalk.com/robot/send?access_token=4e9609518cb06b7177f8cf17cfec5ebffab46ae282a088530f5b39cac8d4a8ad
    # secret for signature
    secret: SEC3f84c04516050aaa54155e87c3b76e74541e03a1eae0ef394053d498e9642116
```

### temp1
**dingtalk.tmpl**
```go
{{ define "__subject" }}
[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}]
{{ end }}


{{ define "__alert_list" }}{{ range . }}
告警级别: {{ .Labels.severity }}

告警名称: {{ index .Annotations "title" }}

负责团队：{{ .Labels.team }}

告警主机: {{ .Labels.instance }}

区域主机: {{ .Labels.location }}

告警信息: {{ index .Annotations "description" }}

告警时间: {{ .StartsAt.Format "2006-01-02 15:04:05" }}

{{ end }}{{ end }}

{{ define "__resolved_list" }}{{ range . }}
告警级别: {{ .Labels.severity }}

告警名称: {{ index .Annotations "title" }}

负责团队：{{ .Labels.team }}

告警主机: {{ .Labels.instance }}

区域主机: {{ .Labels.location }}

告警信息: {{ index .Annotations "description" }}

告警时间: {{ .StartsAt.Format "2006-01-02 15:04:05" }}

恢复时间: {{ .StartsAt.Format "2006-01-02 15:04:05" }}

{{ end }}{{ end }}


{{ define "default.title" }}
{{ template "__subject" . }}
{{ end }}

{{ define "default.content" }}
{{ if gt (len .Alerts.Firing) 0 }}
==== **告警**  ====
{{ template "__alert_list" .Alerts.Firing }}
---
{{ end }}

{{ if gt (len .Alerts.Resolved) 0 }}
==== **告警恢复** ====
{{ template "__resolved_list" .Alerts.Resolved }}
{{ end }}
{{ end }}


{{ define "ding.link.title" }}{{ template "default.title" . }}{{ end }}
{{ define "ding.link.content" }}{{ template "default.content" . }}{{ end }}
{{ template "default.title" . }}
{{ template "default.content" . }}
```

### temp2
```yaml
{{ define "__subject" }}
[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}]
{{ end }}


{{ define "__alert_list" }}{{ range . }}
**告警级别:** {{ .Labels.severity }}

**告警名称:** {{ index .Annotations "title" }}

**负责团队:** {{ .Labels.team |toUpper}}

**告警主机:** {{ .Labels.instance }}

**区域主机:** {{ .Labels.location }}

**告警信息:** {{ index .Annotations "description" }}

**告警时间:** {{ .StartsAt.Format "2006-01-02 15:04:05" }}

{{ end }}{{ end }}

{{ define "__resolved_list" }}{{ range . }}
**告警级别:** {{ .Labels.severity }}

**告警名称:** {{ index .Annotations "title" }}

**负责团队:** {{ .Labels.team |toUpper}}

**告警主机:** {{ .Labels.instance }}

**区域主机:** {{ .Labels.location }}

**告警信息:** {{ index .Annotations "description" }}

**告警时间:** {{ .StartsAt.Format "2006-01-02 15:04:05" }}

**恢复时间:** {{ .StartsAt.Format "2006-01-02 15:04:05" }}

{{ end }}{{ end }}


{{ define "default.title" }}
{{ template "__subject" . }}
{{ end }}

{{ define "default.content" }}#### \[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}\] **[{{ index .GroupLabels "alertname" }}]({{ template "__alertmanagerURL" . }})**
---
{{ if gt (len .Alerts.Firing) 0 }}
{{ template "__alert_list" .Alerts.Firing }}
{{ end }}


{{ if gt (len .Alerts.Resolved) 0 }}
{{ template "__resolved_list" .Alerts.Resolved }}
{{ end }}
{{ end }}


{{ define "ding.link.title" }}{{ template "default.title" . }}{{ end }}
{{ define "ding.link.content" }}{{ template "default.content" . }}{{ end }}
{{ template "default.title" . }}
{{ template "default.content" . }}
```

### temp3
```yaml
{{ define "__subject" }}
[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}]
{{ end }}


{{ define "__alert_list" }}{{ range . }}
**告警级别:** {{ .Labels.severity }}

**告警名称:** {{ index .Annotations "title" }}

**负责团队:** {{ .Labels.team |toUpper}}

**告警主机:** {{ .Labels.instance }}

**区域主机:** {{ .Labels.location }}

**告警信息:** {{ index .Annotations "description" }}

{{ if eq .Labels.alertname "磁盘使用率预警" }}

**事件标签:** {{ index .Annotations "detail" }}

{{ end }}

**告警时间:** {{ .StartsAt.Format "2006-01-02 15:04:05" }}

---

{{ end }}{{ end }}

{{ define "__resolved_list" }}{{ range . }}
**告警级别:** {{ .Labels.severity }}

**告警名称:** {{ index .Annotations "title" }}

**负责团队:** {{ .Labels.team |toUpper}}

**告警主机:** {{ .Labels.instance }}

**区域主机:** {{ .Labels.location }}

**告警信息:** {{ index .Annotations "description" }}

**告警时间:** {{ .StartsAt.Format "2006-01-02 15:04:05" }}

**恢复时间:** {{ .StartsAt.Format "2006-01-02 15:04:05" }}

{{ end }}{{ end }}


{{ define "default.title" }}
{{ template "__subject" . }}
{{ end }}

{{ define "default.content" }}#### \[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}\] **[{{ index .GroupLabels "alertname" }}]({{ template "__alertmanagerURL" . }})**
---
{{ if gt (len .Alerts.Firing) 0 }}
{{ template "__alert_list" .Alerts.Firing }}
{{ end }}


{{ if gt (len .Alerts.Resolved) 0 }}
{{ template "__resolved_list" .Alerts.Resolved }}
{{ end }}
{{ end }}


{{ define "ding.link.title" }}{{ template "default.title" . }}{{ end }}
{{ define "ding.link.content" }}{{ template "default.content" . }}{{ end }}
{{ template "default.title" . }}
{{ template "default.content" . }}
```

## monitor_active_server

- 此服务是基于prometheus二次开发的
- [https://github.com/zhangyy8lab/tusimaServerMonitor](https://github.com/zhangyy8lab/tusimaServerMonitor)
```go
// 以下所有配置的说明 在服务运行时不可以增加对应描述信息
monitorServer:  // 固定格式
  server: [     // 固定格式， 列表中的元素对应docker ps | awk '{print $NF}' 取名称
    "prometheus",
    "monitor-active-server",
    "cadvisor",
    "nginx116",
    "cadvisor-aa",
    "cadvisor_aaa"
  ]

server:// 固定格式 
  port: 8000// 固定格式， 服务启动时占用的端口
```
## docker-compose.yml
```yaml
version: '3'
networks:
  default:
    name: monitor

services:
    prometheus:
        image: prom/prometheus
        container_name: monitor_prometheus
        hostname: prometheus
        restart: always
        command:
          - '--config.file=/etc/prometheus/prometheus.yml'
          - '--web.enable-lifecycle'
          - '--storage.tsdb.retention.time=30d'
        volumes:
            - ./prometheus/:/etc/prometheus/
        ports:
            - 9020:9090

    grafana:
        image: grafana/grafana
        container_name: monitor_grafana
        hostname: grafana
        restart: always
        ports:
            - "9022:3000"
        volumes:
            - ./grafana/config/:/etc/grafana/
            - ./grafana/data:/var/lib/grafana

    alertmanager:
        image: prom/alertmanager
        container_name: monitor_alertmanager
        hostname: alertmanager
        restart: always
        volumes:
            - ./alertmanager/alertmanager.yml:/etc/alertmanager/alertmanager.yml
        ports:
            - 9025:9093
        environment:
          - TZ=Asia/Shanghai


    dingtalk:
        image: timonwong/prometheus-webhook-dingtalk
        container_name: monitor_dingtalk
        hostname: dingtalk
        restart: always
        volumes:
          - ./dingtalk/config.yml:/etc/prometheus-webhook-dingtalk/config.yml
          - ./dingtalk/dingtalk.tmpl:/opt/dingtalk/template/dingtalk.tmpl

        ports:
          - "9026:8060"

        environment:
          - TZ=Asia/Shanghai

    nodeexporter:
        container_name: monitor_node_exporter
        image: prom/node-exporter
        restart: always
        hostname: nodeExport
        ports:
            - 9023:9100
        volumes:
            - /proc:/host/proc:ro
            - /sys:/host/sys:ro
            - /:/rootfs:ro

    activeserver:
        container_name: monitor_active_server
        image: tusimaservermonitor_monitor_active_server
        restart: always
        ports:
          - 9024:8000
        volumes:
          - /usr/bin/docker:/usr/bin/docker
          - ./monitor-server/service.yaml:/app/config/service.yaml
          - /var/run/docker.sock:/var/run/docker.sock
```

## 其他备注说明：

- 在 monitor_active_server 的配置文件中， 每增加一个配置项就需要在prometheus/rules/xx.yml 增加一个对应的监控项如下
```go
vim /monitor/prometheus/rules/active_server.yml
groups:
  - name: prometheus-alert
    rules:
    - alert: 容器服务存活检测
      expr: custom_check_alive_server{cadvisor_aa="stop"} == 0 // custom_check_alive_server是代码中固定的， 监测的容器名称为 cadvisor_aa 服务

      for: 7s  # 1m

      labels:
        severity: P1
        team: "tusima"
        instance: "{{ $labels.instance }}"
        location: "{{ $labels.location }}"

      annotations:
        title: "service-cadvisor-stop"
        description: "\n- server cadvisor stop"
```

- 

