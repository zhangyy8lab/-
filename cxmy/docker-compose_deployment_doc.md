# 持续免疫服务部署

## IMA

### 创建并开启ima度量

#### Ima 功能配置

> 添加至 l文件末尾
>
> ***\*注意：\****
>
> ***\*如果文件\*cat不存在， 则需要查看是否存在/boot/grub2/grub.cfg 这个文件，存在则追加文件内容

```bash
ima=on ima_tcb ima_template=ima ima_hash=sha1 ima_appraise=off
```



#### 配置 ima 策略

> /etc/ima/ima-policy
>
> **ima 在内核 2.6.30 kernel 以上版本默认启用**

```bash
# PROC_SUPER_MAGIC
dont_measure fsmagic=0x9fa0
dont_appraise fsmagic=0x9fa0
# SYSFS_MAGIC
dont_measure fsmagic=0x62656572
dont_appraise fsmagic=0x62656572
# DEBUGFS_MAGIC
dont_measure fsmagic=0x64626720
dont_appraise fsmagic=0x64626720
# TMPFS_MAGIC
dont_measure fsmagic=0x01021994
dont_appraise fsmagic=0x01021994
# RAMFS_MAGIC
dont_appraise fsmagic=0x858458f6
# DEVPTS_SUPER_MAGIC
dont_measure fsmagic=0x1cd1
dont_appraise fsmagic=0x1cd1
# BINFMTFS_MAGIC
dont_measure fsmagic=0x42494e4d
dont_appraise fsmagic=0x42494e4d
# SECURITYFS_MAGIC
dont_measure fsmagic=0x73636673
dont_appraise fsmagic=0x73636673
# SELINUX_MAGIC
dont_measure fsmagic=0xf97cff8c
dont_appraise fsmagic=0xf97cff8c
# CGROUP_SUPER_MAGIC
dont_measure fsmagic=0x27e0eb
dont_appraise fsmagic=0x27e0eb
# NSFS_MAGIC
dont_measure fsmagic=0x6e736673
dont_appraise fsmagic=0x6e736673
measure func=FILE_CHECK mask=MAY_EXEC
#measure func=FILE_MMAP mask=MAY_EXEC
```

#### ima 度量文件

> 位置：/sys/kernel/security/ima/ascii_runtime_measurements  #可以查看日志，如果有新的日志产生说明生效了



## mongoDB

### db_path

```bash
# 创建挂载目录  
mkdir -p /8lab/data/mongodb{0,1,2}
```



### docker-compose.yml

>/8lab/mongodb/docker-compose.yml

```yaml
version: '3'
services:
  mongo0:
    image: 1017127423/8lab:mongo_v210823
    container_name: mongo0
    hostname: mongo0
    restart: always
    ports:
      - "27017:27017"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=example
    command: mongod --replSet "rs0" --bind_ip 0.0.0.0 --port 27017
    volumes:
      - /8lab/data/mongodb0:/data/db

  mongo1:
    image: 1017127423/8lab:mongo_v210823
    container_name: mongo1
    hostname: mongo1
    restart: always
    ports:
      - "27018:27017"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=example
    command: mongod --replSet "rs0" --bind_ip 0.0.0.0 --port 27017
    volumes:
      - /8lab/data/mongodb1:/data/db

  mongo2:
    image: 1017127423/8lab:mongo_v210823
    container_name: mongo2
    hostname: mongo2
    restart: always
    ports:
      - "27019:27017"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=example
    command: mongod --replSet "rs0" --bind_ip 0.0.0.0 --port 27017
    volumes:
      - /8lab/data/mongodb2:/data/db

  mongo-init:
    image: 1017127423/8lab:mongo_v210823
    container_name: mongo-init
    depends_on:
      - mongo0
      - mongo1
      - mongo2
    command: >
      bash -c "
      sleep 10;
      echo 'Initiating replica set';
      mongo --host mongo0:27017 --eval '
        rs.initiate({
          _id: \"rs0\",
          members: [
            { _id: 0, host: \"mongo0:27017\" },
            { _id: 1, host: \"mongo1:27017\" },
            { _id: 2, host: \"mongo2:27017\" }
          ]
        });
      '"
```



### run&stop

```bash 
# 启动
docker-compose up -d 

# 停 
docker-compose down
```

### check_status

```bash
# 进入容器
docker exec -it monogo0 mongo

# 查看状态
rs.status()

# 查看配置
rs.conf()

# 退出
exit
```



## bitchainDb

### octachain.cfg

> /8lab/conf/bigchain/octachain.cfg

```json
{
    "database": {
        "certfile": null,
        "host": "172.19.23.177",
        "crlfile": null,
        "port": 28000,
        "backend": "mongodb",
        "keyfile_passphrase": null,
        "ca_cert": null,
        "ssl": false,
        "keyfile": null,
        "replicaset": "8lab-shard",
        "connection_timeout": 5000,
        "login": null,
        "max_tries": 3,
        "password": null,
        "name": "octachain"
    },
    "server": {
        "workers": null,
        "bind": "0.0.0.0:10070",
        "loglevel": "info"
    },
    "log": {
        "file": "/var/log/8lab/octachain.log",
        "level_console": "info",
        "level_logfile": "info",
        "datefmt_logfile": "%Y-%m-%d %H:%M:%S",
        "fmt_logfile": "[%(asctime)s] [%(levelname)s] %(message)s (%(processName)-10s - pid: %(process)d)",
        "granular_levels": {},
        "error_file": "/var/log/8lab/octachain-errors.log",
        "datefmt_console": "%Y-%m-%d %H:%M:%S",
        "fmt_console": "[%(asctime)s] [%(levelname)s] %(message)s (%(processName)-10s - pid: %(process)d)"
    },
    "backlog_reassign_delay": 120,
    "wsserver": {
        "scheme": "ws",
        "host": "0.0.0.0",
        "port": 10071
    },
    "graphite": {
        "host": "0.0.0.0"
    },
    "keypair": {
        "private": "4EC7mJs7R5jQpMmkLpQcyevZ8Nimm9xGkBJDKv2SuHej",
        "public": "6osKMk7gFcYhF9Xsc5KApVwLtKeJ2NRpzi6JP9755Zcs"
    },
    "keyring": []
}
```



### docker-compose.yml

```yaml
version: '3'
services:
  bigchaindb:
    image: 1017127423/8lab:octachain_v1020
    container_name: bigchaindb
    restart: always
    volumes:
      - /8lab/conf/bigchain/octachain.cfg:/etc/octachain/octachain.cfg
      - /data/logs/octachain:/app/logs 
    ports:
      - "10070:10070"
      - "10071:10071"
```

### run&stop

```bash 
# 启动
docker-compose up -d 

# 停 
docker-compose down
```



## mysql

### db_path

```bash
# 目录
mkdir /8lab/data/mysql

# 修改属主 ,根据容器中 mysql 的u_id 进行调整属主
chown <u_id>:<g_id> /
# - 启动时可能会提示无法创建 err.log ,先创建目录后， 启动， 查看日志可根据报错信息进行处理

# 也可以直接
chmod -R 777 /8lab/data/mysql
```

 

### mysql.cnf

> /8lab/conf/mysql/mysql.cnf

```bash
[client]
port            = 3306
socket          = /data/mysql/tmp/mysql.sock
user=mysql
password=mysql

[mysqld]
bind-address = 0.0.0.0
port            = 3306
user            = mysql
# basedir         = /usr/local/mysql
# datadir         = /data/mysql3306/data
# tmpdir          = /data/mysql3306/tmp
# socket          = /data/tmp/mysql.sock
# pid-file        = /data/tmp/mysql.pid
# log-bin         = /data/mysql3306/logs/bin-log
# log-error       = /data/mysql3306/tmp/err.log

basedir         = /usr/local/mysql
datadir         = /data/mysql/data
tmpdir          = /data/mysql/tmp
socket          = /data/tmp/mysql.sock
pid-file        = /data/tmp/mysql.pid
log-bin         = /data/logs/bin-log
log-error       = /data/logs/err.log
slow_query_log_file = /data/tmp/slow.log
binlog_format=row
default-storage-engine=innodb
character_set_server = utf8
lower_case_table_names = 1
skip_external_locking
skip-name-resolve
group_concat_max_len=102400000
#skip-networking
binlog-ignore-db=mysql
binlog-ignore-db=information_schema
replicate_ignore_db=mysql
replicate_ignore_db=information_schema
skip-slave-start
#skip-grant-tables
log-bin-trust-function-creators=1
#slave-skip-errors=i1032,1062,1053,1146
#slave-skip-errors=all
#read_only
explicit_defaults_for_timestamp=true

connect_timeout = 30
interactive_timeout = 1000
wait_timeout = 180
event_scheduler = 1
log_bin_trust_function_creators = 1
back_log = 500
##### binlog #####
#binlog_format = row
max_binlog_size = 512M
expire_logs_days = 10
binlog_cache_size = 2M

##### gtid #####
log-slave-updates = true
gtid-mode=on
enforce_gtid_consistency=on
master-info-repository = TABLE
relay-log-info-repository = TABLE
sync-master-info = 1
slave-parallel-workers = 4
binlog-checksum = CRC32
master-verify-checksum = 1
slave-sql-verify-checksum = 1
binlog-rows-query-log_events = 1
#report-port = 3306
#report-host = 192.168.1.193

##### replication #####
#skip-slave-start
auto_increment_increment = 2
auto_increment_offset = 1
log_slave_updates = 1
slave_net_timeout = 3600
relay_log_recovery = 1

replicate_wild_ignore_table=mysql.%
#replicate-ignore-db=information_schema
#replicate-ignore-db=mysql
#replicate-ignore-db=performance_schema
#replicate-ignore-db=test

##### slow log #####
slow_query_log = 1
#slow_query_log_file
long_query_time = 2

##### error log #####
#log_error

##### thread #####
max_connections = 3000
thread_stack = 256K
max_allowed_packet = 512M
table_open_cache = 2000
read_buffer_size = 4M
read_rnd_buffer_size = 2M
sort_buffer_size = 2M
join_buffer_size = 128M

##### InnoDB #####
innodb_data_home_dir = /8lab/data/mysqldb/data
innodb_log_group_home_dir = /8lab/data/mysqldb/logs
innodb_data_file_path = ibdata1:12M:autoextend
innodb_buffer_pool_size = 2048M
innodb_log_file_size = 1G
innodb_file_per_table = 1
innodb_flush_log_at_trx_commit = 2
sync_binlog = 0
innodb_thread_concurrency = 0
innodb_flush_method = O_DIRECT
innodb_lock_wait_timeout = 50
transaction_isolation=READ-COMMITTED

##### MyISAM #####
key_buffer_size = 128M

##### OTHER #####
#sql_mode=NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES
tmp_table_size = 256M
max_heap_table_size = 32M
thread_cache_size = 64
#thread_concurrency = 12
#bulk_insert_buffer_size
open_files_limit = 65535

# to avoid issues with 'bulk mode inserts' using autoinc
innodb_autoinc_lock_mode=2
query_cache_size=0
query_cache_type=0
bind-address=0.0.0.0
innodb_doublewrite=0

[mysql]
no-auto-rehash
```

### docker-compose.yml

> /8lab/mysql/docker-compose.yml

```yaml
version: '3.8'

services:
  mysql:
    image: octahub.8lab.cn:5000/octa-cis/mysql:v0826
    container_name: mysql
    ports:
      - "32143:3306"
    volumes:
      - /8lab/conf/mysql/mysql.cnf:/opt/mysql3306.cnf
      - /8lab/data/mysql/:/data/
    restart: always
```

## redis

### db_path

```bash
mkdir /8lab/data/redis
```



> /8lab/conf/redis/redis.conf

### redis.conf

```bash
daemonize no
port 6379
bind 0.0.0.0
requirepass octa8lab
masterauth octa8lab
pidfile /usr/local/redis/pid/redis.pid
logfile ""
dbfilename "redis.rdb"
dir "/data/redis"
tcp-keepalive 60
tcp-backlog 511
timeout 3000
loglevel notice
databases 16
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error no
rdbcompression yes
slave-serve-stale-data no
repl-diskless-sync yes
repl-diskless-sync-delay 5
repl-disable-tcp-nodelay no
slave-priority 100
appendonly no
no-appendfsync-on-rewrite yes
lua-time-limit 5000
slowlog-log-slower-than 10000
slowlog-max-len 128
latency-monitor-threshold 0
notify-keyspace-events ""
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-entries 512
list-max-ziplist-value 64
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit slave 256mb 64mb 60
client-output-buffer-limit pubsub 64mb 16mb 60
hz 10
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  redis:
    image: redis
    container_name: redis
    ports:
      - "32144:6379"
    volumes:
      - /8lab/conf/redis/redis.conf:/usr/local/etc/redis/redis.conf
      - /8lab/data/redis/:/data/redis/
    command: ["redis-server", "/usr/local/etc/redis/redis.conf"]
```



## kafka

### db_path

```bash 
mkdir /8lab/data/kafka && chmod 777 -R  /8lab/data/kafka
```

### docker-compose.yml

```yaml
version: '3'
services:
  kafka-server:
    image: bitnami/kafka:latest
    container_name: kafka-server
    ports:
      - '9092:9092'
    environment:
      - KAFKA_ENABLE_KRAFT=yes
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_CFG_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://172.19.23.177:9092
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_NODE_ID=1
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@localhost:9093
      - KAFKA_CFG_LOG_DIRS=/tmp/kraft-combined-logs
    volumes:
      - /8lab/data/kafka:/tmp/kraft-combined-logs
```

### start & stop & log

```bash
# start
docker-compose -f /8lab/kafka-server/docker-compose.yml up -d 

# stop
docker-compose -f /8lab/kafka-server/docker-compose.yml down

# log
docker logs -fn11 kafka-server  

# 查看日志如果存在有异常问题， 可将 /8lab/data/kafka/ 下的文件删除 命令为 rm -rf /8lab/data/fafka/*
```



### topic

```bash
# topic 作用 
csmp-ty-sysaudit-log logstash-5145 input filebeat写入kafka 
csmp-output-sysaudit-log logstash-5145 out 到kafka

csmp-output-sysaudit-log logstash-5156 input 从kafka 读
csmp-std-sandboxaudit-log logstash-5156 output

csmp-std-sandboxaudit-log  draven 读

# 以下为帮忙信息
# 查看规则
docker exec -it kafka-server kafka-topics.sh --create --topic <topic_name> --partitions 3 --replication-factor 1 --bootstrap-server localhost:9092

# 查看所有topic
docker exec -it kafka-server  kafka-topics.sh --list --bootstrap-server localhost:9092


# 查看topic消息
docker exec -it kafka-server kafka-console-consumer.sh --topic csmp-std-sandboxaudit-log --from-beginning --bootstrap-server localhost:9092

# 创建组， 默认会创建， 如何没有手动创建
docker exec -it kafka-server kafka-console-consumer.sh --bootstrap-server localhost:9092 \
  --topic csmp-std-sandboxaudit-log \
  --group GID_STD_CSMP_SANDBOXAUDIT_LOG_002 \
  --from-beginning

# 查看组
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --group GID_STD_CSMP_SANDBOXAUDIT_LOG_002 --topic csmp-std-sandboxaudit-log 


docker exec -it kafka-server  kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group GID_STD_CSMP_SANDBOXAUDIT_LOG_002 --topic csmp-std-sandboxaudit-log

kafka-consumer-groups.sh --bootstrap-server localhost:9092 --group GID_STD_CSMP_SANDBOXAUDIT_LOG_002 --describe --topic 
```

### demo-test.py

```python
# 测试kafka 生产/消费  demo.py
# apt install python3-pip 
# pip3 install confluent_kafka --break-system-packages

vim kafkaTest.py
from confluent_kafka import Producer, Consumer, KafkaError

# Kafka 配置
kafka_config = {
    'bootstrap.servers': '192.168.110.131:9092',  # Kafka 服务器地址
    'group.id': 'csmp',  # 消费者所属的组 ID
    'auto.offset.reset': 'earliest'  # 如果没有可用的偏移量，从最早的消息开始消费
}

# 生产者
def produce():
    producer = Producer({'bootstrap.servers': kafka_config['bootstrap.servers']})

    def delivery_report(err, msg):
        if err is not None:
            print(f'Message delivery failed: {err}')
        else:
            print(f'Message delivered to {msg.topic()} [{msg.partition()}]')

    # 发送消息
    for i in range(10):
        message = f'Message {i}'
        producer.produce('test-topic', message.encode('utf-8'), callback=delivery_report)

    # Wait up to 1 second for events
    producer.flush()

# 消费者
def consume():
    consumer = Consumer({
        'bootstrap.servers': kafka_config['bootstrap.servers'],
        'group.id': kafka_config['group.id'],
        'auto.offset.reset': kafka_config['auto.offset.reset']
    })
    consumer.subscribe(['test-topic'])

    try:
        while True:
            msg = consumer.poll(1.0)

            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    # End of partition event
                    print(f'End of partition reached {msg.topic()}/{msg.partition()}')
                else:
                    print(f'Error occurred: {msg.error().str()}')
                continue

            print(f'Received message: {msg.value().decode("utf-8")}')

    except KeyboardInterrupt:
        pass
    finally:
        consumer.close()

# 运行生产者和消费者
produce()
consume()
```

## auditd

> 部署: node1 node2
>
> 收集日志信息 /var/log/syslog

### install

```bash
# 安装 auditd
apt install auditd
```

### policy

```bash
# 规则文件解压到服务器
rm -rf /etc/audit

# /xx/xxx/audit_.tar.gz 注意路径，文件在
tar -xzf audit.tar.gz -C /etc/ 

# 这里会有很多规则文件
```

### start & Stop

```bash
systemctl reload auditd

# 开机启动
systemctl enable auditd

# 启动
systemctl start auditd

# 停止
systemctl stop auditd
```

## filebeat

> 部署：  node1 node2
>
> 将收集的日志 push 到kafka

```bash
8lab/filebeat
├── apt_update.sh
└── filebeat
    └── filebeat-7.17.18-x86_64.rpm
```

### Install 

```bash
cd /8lab/filebeat/ 
sh apt_update.sh 

# 配置文件生成路径：/etc/filebeat/filebeat.yml 
# 需要 注释es相关配置 line 66 
```

### start & stop & log

```bash
# start
systemctl start filebeat.service

# stop 
systemctl stop filebeat.service

# log
journalctl -u filebeat.service -f
```



## damper_server防篡改

> /8lab/conf/dtamper_server/server_configure.json

### server_configure.json

```json
{
  "service_type": ["web"],
  "client_address": "172.19.23.177",
  "server_address": "172.19.23.177",
  "bdb_host":"172.19.23.177",
  "bdb_port":"32145",
  "redis_host": "172.19.23.177",
  "redis_port": "32144",
  "redis_password": "octa8lab",
  "update_channel": "update",
  "report_channel": "return",
  "control_send_channel": "control_send",
  "control_return_channel": "control_return",
  "heartbeat_channel": "heartbeat",
  "heartbeat_timeout": "30",
  "heartbeat_log_freeze": 100,
  "aes_key": "ca$hc0w8L@6ExP0!",
  "aes_iv": "ca$hc0w8L@6ExP0!",
  "mysql_host": "172.19.23.177",
  "mysql_port": "32143",
  "mysql_username": "8lab",
  "mysql_password": "8lab",
  "token_switch": "off"
}
```



### docker-compose.yml

> /8lab/dtamper_server/docker-compose.yml

```yaml
version: '3'
services:
  restful_service:
    image: octahub.8lab.cn:5000/octa-cis/dtamper-server-4.3:202401231050
    container_name: dtamper_service
    volumes:
      - /8lab/conf/dtamper_server/server_configure.json:/usr/local/dtamper_server/bootstrap/server_configure.json
    ports:
      - "2345:2345"
    command: restful
    restart: unless-stopped

  heatbeat_monitor:
    image: octahub.8lab.cn:5000/octa-cis/dtamper-server-4.3:202401231050
    container_name: heatbeat_monitor
    volumes:
      - /8lab/conf/dtamper_server/server_configure.json:/usr/local/dtamper_server/bootstrap/server_configure.json
    command: heatbeat_monitor
    restart: unless-stopped

  task_monitor:
    image: octahub.8lab.cn:5000/octa-cis/dtamper-server-4.3:202401231050
    container_name: task_monitor
    volumes:
      - /8lab/conf/dtamper_server/server_configure.json:/usr/local/dtamper_server/bootstrap/server_configure.json
    command: task_monitor
    restart: unless-stopped

  task_dispatcher:
    image: octahub.8lab.cn:5000/octa-cis/dtamper-server-4.3:202401231050
    container_name: task_dispatcher
    volumes:
      - /8lab/conf/dtamper_server/server_configure.json:/usr/local/dtamper_server/bootstrap/server_configure.json
    command: task_dispatcher
    restart: unless-stopped
```

### run & stop & log

```bash
# run
docker-compose -f /8lab/dtamper_server/docker-compose.yml up -d 

# stop
docker-compose -f /8lab/dtamper_server/docker-compose.yml down

# log
docker logs -fn11 dtamper_service
```





## damper_client防篡改

### client_configure.json

> /8lab/dtamper_client/client_configure.json

```json
{
  "service_type": "web",
  "client_address": "172.19.23.177",
  "server_address": "172.19.23.177",
  "bdb_host":"172.19.23.177",
  "bdb_port":"32145",
  "redis_host": "172.19.23.177",
  "redis_port": "32144",
  "redis_password": "octa8lab",
  "update_channel": "update",
  "report_channel": "return",
  "control_send_channel": "control_send",
  "control_return_channel": "control_return",
  "heartbeat_channel": "heartbeat",
  "aes_key": "ca$hc0w8L@6ExP0!",
  "aes_iv": "ca$hc0w8L@6ExP0!",
  "mysql_host": "172.19.23.177",
  "mysql_port": "32143",
  "mysql_username": "8lab",
  "mysql_password": "8lab"
}
```

### client_db_generate_plain

> 生成db

```bash
cd /8lab/dtamper_client/ && ./client_db_generate_plain

# 会生成web_client.db
```

### dtamper_path_monitor.service

>  /etc/systemd/system/dtamper_path_monitor.service

```bash
[Unit]
Description= tamper_path_monitor service

[Service]
ExecStart=/data/yaml/dtamper_client/dtamper_path_monitor
Restart=always
User=root
WorkingDirectory=/data/yaml/dtamper_client
[Install]
WantedBy=multi-user.target

```

### dtamper_client.service

> /etc/systemd/system/dtamper_client.service

```bash
[Unit]
Description= dtamper_client service

[Service]
ExecStart=/8lab/dtamper_client/dtamper_client
Restart=always
User=root
WorkingDirectory=/8lab/dtamper_client

[Install]
WantedBy=multi-user.target
```

### start & stop & log

```bash
systemctl daemon-reload

# enable start
systemctl enable dtamper_path_monitor.service
systemctl enable dtamper_client.service

# start
systemctl start dtamper_path_monitor.service
systemctl start dtamper_client.service

# stop 
systemctl stop dtamper_path_monitor.service
systemctl stop dtamper_client.service

# log
journalctl -u dtamper_path_monitor.service
```



## oat_server 可信Server

### logging-linux.yaml

> /8lab/conf/oat_server/logging-linux.yaml

```bash
version: 1

formatters:
    defaultFormatter:
        format: '%(asctime)s [%(levelname)s] [%(module)s:%(funcName)s] [%(lineno)d] - %(message)s'
        datefmt: '%Y-%m-%d %H:%M:%S'

handlers:
    consoleHandler:
        class: logging.StreamHandler
        level: DEBUG
        stream: 'ext://sys.stdout'
        formatter: defaultFormatter

    fileHandler:
        class: logging.handlers.EnhancedRotatingFileHandler
        level: DEBUG
        formatter: defaultFormatter
        filename: '/8lab/log/oat.log'
        when: 'MIDNIGHT'
        interval: 1
        backupCount: 100000
        maxBytes: 209715200 #200M

loggers:
    root:
        level: DEBUG
        handlers: [consoleHandler, fileHandler]
    tx_module:
        level: DEBUG
        handlers: [consoleHandler, fileHandler]
        qualname: tx_module
```



### oat.yml

> /8lab/conf/oat_server/oat.yaml

```yaml
host:
  oat_port: 50051
  max_workers: 50
  max_send_message_lenght: 50        # 鏈€澶у彂閫佹秷鎭殑闀垮害 M
  max_receive_message_lenght: 50     # 鏈€澶ф帴鏀跺瓧鑺傜殑闀垮害 M

# 閫夋嫨mq涓棿浠剁被鍨� 鐢ㄦ潵鎺ユ敹agent鍙戦€佽繃鏉ョ殑鐧藉悕鍗� redis 鎴� RabbmitMQ
#queue_type: rabbitmq 涓篗Q妯″紡
#queue_type: redis 涓簉edis妯″紡
# 姝ら厤缃簲涓巓at_agent閰嶇疆鏂囦欢涓璠open_mq]鐨勯厤缃浉鍚屾
queue:
  queue_type: redis
  handler_threads: 5  # 鍗曞彴MQ澶勭悊鏁版嵁鐨勫惎鐢ㄧ殑绾跨▼鏁伴噺

# 蹇冭烦鐩稿叧閰嶇疆锛屽繀椤讳笌agent閰嶇疆鏂囦欢涓璠heartbeat]鐩稿悓姝�
# monitor_heartbeat鐨勫€煎簲涓巋b_interval鐨勫€肩浉鍚�
heartbeat:
  monitor_heartbeat: 10  # 鐩戝惉蹇冭烦棰戠巼榛樿涓�10s/娆� 绉掍负鍗曚綅
  expire_time: 15  # agent蹇冭烦杩囨湡鏃堕棿榛樿涓�15绉掕繃鏈� 绉掍负鍗曚綅 杩囨湡鏃堕棿蹇呴』澶т簬monitor_heartbeat鍜宎gent鐨刪b_interval鍊�

# 鏈ā鍧楀彧璁板綍杩炴帴鐨勭涓変腑闂翠欢
middleware:
  mysql:
    host: 172.19.23.177
    user: 8lab
    password: 8lab
    db_name: octa_cis
    port: 32143
  redis:
    host: 172.19.23.177
    port: 32144
    # 濡傛灉redis娌℃湁瀵嗙爜 閰嶇疆涓� no
    password: octa8lab
    db: 0
    # redis闆嗙兢鏍囧織 濡傛灉涓嶆槸闆嗙兢 鍊间负锛歯o 濡傛灉鏄泦缇ゅ€兼牱寮忥細192.168.1.1:9999,192.168.1.2:8888
    redis_cluster: no
  rabbmitmq:
    host: 192.168.3.100
    port: 5672
    vhost: /
    user: admin
    password: admin

# url example: 127.0.0.1:9999 涓嶅悜闈掍簯鎶ヨ 榛樿涓� no
qingcloud:
  qc_url: no

log:
  # linux 鏃ュ織閰嶇疆鏂囦欢璺緞
  log_config_file_linux: /8lab/conf/oat_server/logging-linux.yaml
  # windows 鏃ュ織閰嶇疆鏂囦欢璺緞
  log_config_file_windows: C:\8lab\conf\oat_server\logging-windows.yaml
  # mac 鏃ュ織閰嶇疆鏂囦欢璺緞
  log_config_file_mac: /Users/user_mac_name/Desktop/oat_server/logging-linux.yaml
```



### docker-compose.yml

> /8lab/oat_server/docker-compose.yml

```bash
version: '3'

services:
  oat_server:
    image: octahub.8lab.cn:5000/octa-cis/kx_oat_server:20231012001
    container_name: oat_server
    ports:
      - "50051:50051"
    volumes:
      - /8lab/conf/oat_server/:/8lab/conf/oat_server/
    restart: unless-stopped
```

### start & stop & log

```bash
# start
docker-compose -f /8lab/oat_server/docker-compose.yml up -d 

# stop
docker-compose -f /8lab/oat_server/docker-compose.yml down

# log
docker logs -fn11 oat_server
```



## oat_agent可信client

### pytagent

> /8lab/oat_tagent/pytagent



### oat_pytagent.service

> /etc/systemd/system/oat_pytagent.service


```bash
[Unit]
Description=oat_pytagent service

[Service]
ExecStart=/8lab/oat_tagent/pytagent
Restart=always
User=root
WorkingDirectory=/8lab/oat_tagent
[Install]
WantedBy=multi-user.target
```



### start & stop & log

```bash
# enable start 
systemctl enable pytagent.service

# start
systemctl start pytagent.service

# stop
systemctl stop pytagent.service

# log
tail -f /8lab/logs/oat_agent/tagent.log 
	or 
journalctl -u pytagent.service -f
```

## octa_cis  WebServer

### conf.json

> /8lab/conf/octa_cis/conf.json

```json
{
  "client_audit_hosts": [
    {
      "ip": "k8sproxy-service",
      "name": "client_b"
    }
  ],
  "product_uid": "bd43-430e-af62-e710",
  "server_ip": "172.19.23.177",
  "blackbox_ip": "192.168.1.169",
  "switch_waf_port": 8080,
  "es_server_ip_port": [{"host": "172.19.23.177", "port": 29200}],
  "es_server_user_name": "elastic",
  "es_server_password": "tarena",
  "bdb_host": "172.19.23.177",
  "bdb_port": "32145",
  "bdb_mongo_host": "172.19.23.177",
  "mongo_port": 27017,
  "rabbitmq_server": "192.168.1.184",
  "rabbitmq_port": 5672,
  "draven_server": "172.19.23.177",
  "draven_port": "8097",
  "db_ip_list": [
    {
      "ip": "192.168.1.236",
      "db_types": "mysql;PostgreSQL",
      "name": "mysql",
      "mask": "255.255.252.0"
    },
    {
      "ip": "192.168.1.181",
      "db_types": "PostgreSQL",
      "name": "PostgreSQL",
      "mask": "255.255.252.0"
    }
  ],
  "dvwa_address": "http://123.56.124.137/dvwa/login.php",
  "used_ports": {
    "waf": 5000,
    "clamav_rpc": 9090,
    "attack_ip": 5001,
    "whitelist": 5570,
    "blackbox_rpc": 5577,
    "mysql_audit": 6666,
    "clamav": 5555
  },
  "mysql_password": "8lab",
  "mysql_user": "8lab",
  "mysql_database": "octa_cis",
  "mysql_port": 32143 ,
  "mysql_host": "172.19.23.177",

  "nisa_mysql_user": "8lab",
  "nisa_mysql_password": "8lab",
  "nisa_mysql_database": "nisa",
  "nisa_mysql_port": 32143,
  "nisa_mysql_host": "172.19.23.177",

  "draven_mysql_user": "8lab",
  "draven_mysql_password": "8lab",
  "draven_mysql_database": "ueba_web",
  "draven_mysql_port": 32143,
  "draven_mysql_host": "172.19.23.177",

  "ueba_web_mysql_password": "8lab",
  "ueba_web_mysql_user": "8lab",
  "ueba_web_mysql_database": "ueba_web",
  "ueba_web_mysql_port": 32143,
  "ueba_web_mysql_host": "172.19.23.177",

  "alarm_second": 15,
  "alarm_enable": 1,
  "redis4bigchanidb_host":"172.19.23.177",
  "redis4bigchanidb_port": 32144,
  "redis4bigchanidb_password": "octa8lab",
  "eagle_host": "192.168.1.182",
  "eagle_port": 9099,
  "enable_register": "on",
  "ifr_urls": {
    "k": "http://192.168.1.163:5601",
    "b": "http://192.168.1.114:9099/eagle-service",
    "t": "http://192.168.1.246:1443/tpotweb.html"
  },
  "bak_mongo_host": "192.168.1.231",
  "bak_mongo_port": 27000,
  "trustlog_index": "trustlog192.168.1.169*",
  "tamper_proof_url": "http://172.19.23.177:2345",
  "chain_attach_mount_dir": "/8lab/upload/chain-attach/",
  "chain_attach_url": "/media/chain-attach/",
  "snort_rpc_port": 30303,
  "des_ip_addr": "110.88.128.28",
  "waf_index": "waf*",
  "face": {
    "model_detection": "hog",
    "tolerance": "0.45",
    "knn_bool": "False",
    "k": "1",
    "scale_ratio": "0.5",
    "addr": "rstp://xxxxxxxxxxx",
    "type": "usb",
    "size": "640,480"
  },

  "mysql_dtamper_web": "web_tamper_proof_db",
  "mysql_dtamper_svn": "svn_tamper_proof_db",

  "dtamper_redis_host": "172.19.23.177",
  "dtamper_redis_port": "32144",
  "dtamper_redis_password": "octa8lab",

  "aes_key": "ca$hc0w8L@6ExP0!",
  "aes_iv": "ca$hc0w8L@6ExP0!",

  "pay_alarm_key": "no_pay",
  "qingyun_user": "996996",
  "qingyun_pwd": "db45030336cc29a66d8404b3b14e7d5b10f8b09f",

  "email_from": "Warning<warn@8lab.cn>",
  "email_host": "smtp.ym.163.com",
  "email_port": 25,
  "email_user": "warn@8lab.cn",
  "email_pass": "8labtestinfo",

  "Alidayu_app": "23827310",
  "Alidayu_key": "24dbb9f199ea5fbc2826e4f2662c15df",
  "Alidayu_template1": "SMS_67765123",
  "Alidayu_template2": "SMS_67715193",
  "Alidayu_sms_free_sign_name": "八分量持续免疫系统",

  "clam_scheduler_url": "http://127.0.0.1:9000",

  "export_host": "http://127.0.0.1",
  "export_file_dir": "/data/export_file_dir/",

  "node_list": [
    {
      "ip": "192.168.1.98",
      "port": "28000",
      "show_ip": "192.168.1.98"
    },
    {
      "ip": "192.168.1.95",
      "port": "28000",
      "show_ip": "192.168.1.95"
    },
    {
      "ip": "192.168.1.97",
      "port": "28000",
      "show_ip": "192.168.1.97"
    }
  ]
}
```



### rpc.json

> /8lab/conf/octa_cis/rpc.json

```json

{
  "user_portrait": {
    "host": "101.251.211.205",
    "port": 9898
  },
  "grpc_param": {
    "host": "172.19.23.177",
    "port": 50051
  }
}
```



### docker-compose.yml

> /8lab/octa_cis/docker-compose.yml

```yaml
version: '3'

services:
  octa_cis:
    image: octahub.8lab.cn:5000/octa-cis/octa_cis_fc:202402021108
    container_name: octa_cis_server
    ports:
      - "8001:8001"
    volumes:
      - /8lab/conf/octa_cis/:/usr/local/octa_cis/conf/
    sysctls:
      net.core.somaxconn: 65535
    restart: unless-stopped
```



### start & stop & log

```bash
# start 
docker-compose -f /8lab/octa_cis/docker-compose.yml up -d

# stop 
docker-compose -f /8lab/octa_cis/docker-compose.yml down 

# log
docker logs -fn1 octa_cis_server
```



## octa_cis  Web

### octa_cis.conf

> /8lab/conf/octa_cis_web/octa_cis.conf

```bash

server {
        listen         8099;
        server_name    .*;
        charset UTF-8;
        access_log      /data/logs/nginx/nginx_access.log;
        error_log       /data/logs/nginx/error_access.log;

        client_max_body_size 75M;

        root /usr/local/octa_cis_web_4_3/app_views/static/prod;
        index index.html;

        location /static/prod {
            autoindex off;
            add_header Cache-Control private;
            alias /usr/local/octa_cis_web_4_3/app_views/static/prod/;
        }

        location /api {
            include uwsgi_params;
            uwsgi_pass 172.19.23.177:8001;
            uwsgi_read_timeout 30;
        }

        location / {
            try_files $uri $uri/ /index.html;
        }

        location /data/export_file_dir {
            autoindex off;
            add_header Cache-Control private;
            alias /usr/local/octa_cis_4_3/export_file_dir;
        }

    }
```



### app_views

>  web static 



### docker-compose.yml

> /8lab/octa_cis_web/docker-compose.yml

```yaml
version: '3'
services:
  octa_cis_web:
    image: octahub.8lab.cn:5000/octa-cis/octa_cis_web_fc_dongwu:202402021616
    container_name: octa_cis_web
    ports:
      - "8099:8099"
    volumes:
      - /8lab/conf/octa_cis_web/:/etc/nginx/conf.d/
      - /8lab/octa_cis_web/app_views/:/usr/local/octa_cis_web_4_3/app_views/
    restart: unless-stopped
```

### start & stop & log

```bash
# start 
docker-compose -f /8lab/octa_cis_web/docker-compose.yml up -d

# stop 
docker-compose -f /8lab/octa_cis_web/docker-compose.yml down 

# log
docker logs -fn1 octa_cis_web
```



## draven

### application-pro.properties

> /8lab/conf/draven/application-pro.properties

```bash
# kafka
kafka.consumer.servers=172.19.23.177:9092
#最早未被消费的offset,设置为earliest;当前的latest
kafka.consumer.auto.offset.reset=earliest
kafka.consumer.group.id=GID_STD_CSMP_SANDBOXAUDIT_LOG_002
#批量消费一次最大拉取的数据量
kafka.consumer.max.poll.records=2000
#是否开启自动提交
kafka.consumer.enable.auto.commit=true
#自动提交的间隔时间
kafka.consumer.auto.commit.interval=1000
#连接超时时间
kafka.consumer.session.timeout=300000
#手动提交设置与poll的心跳数
kafka.consumer.max.poll.interval=30000
#是否开启批量消费，true表示批量消费
kafka.listener.batch.listener=true
#设置消费的线程数
kafka.consumer.concurrency=5
kafka.listener.poll.timeout=5000
kafka.consumer.topic=csmp-std-sandboxaudit-log
# kafka.consumer.topic=sandbox_hdfs_audit_log
kafka.consumer.key-deserializer=org.apache.kafka.common.serialization.StringDeserializer
kafka.consumer.value-deserializer=org.apache.kafka.common.serialization.StringDeserializer

#认证开关 0开1关
kafka.switch=1
kafka.consumer.properties.sasl.mechanism=PLAIN
kafka.consumer.properties.sasl.jaas.config=org.apache.kafka.common.security.plain.PlainLoginModule required username="8lab" password="8lab";
kafka.consumer.properties.security.protocol=SASL_PLAINTEXT

#krb5认证开关 0开1关
kafka.krb5=1
kafka.krb5.conf=/data/storm4.3-topology/krb5.conf
kafka.client.jaas=/data/storm4.3-topology/kafka_client_jaas.conf
kafka.consumer.properties.sasl.kerberos.service.name=kafka


# ==================================== mysql 配置 ====================================
spring.datasource.driver-class-name=com.mysql.jdbc.Driver
spring.datasource.url=jdbc:mysql://172.19.23.177:32143/ueba_web?useSSL=false
spring.datasource.username=8lab
spring.datasource.password=8lab
spring.datasource.max-idle=10
spring.datasource.max-wait=10000
spring.datasource.min-idle=5
spring.datasource.initial-size=5

# ==================================== redis配置 ====================================
spring.redis.host=172.19.23.177
spring.redis.port=32144
spring.redis.password=octa8lab
spring.redis.database=0
spring.redis.lettuce.pool.max-active=10
spring.redis.lettuce.pool.max-idle=20
spring.redis.lettuce.pool.min-idle=10
spring.redis.lettuce.pool.max-wait=-1
spring.redis.timeout=30000

# ==================================== mybatis 配置 ====================================
mybatis.typeAliasesPackage: com.example.demo.entity
mybatis.mapperLocations: classpath:mapper/*.xml
#configLocation: classpath:/mybatis-config.xml
```

### docker-compose.yml

> /8lab/draven/docker-compose.yml

```yaml
version: '3.8'

services:
  draven-app:
    image: 1017127423/8lab:rhel-draven_202401021
    container_name: draven
    ports:
      - "8097:8097"  # 暴露的端口映射
    working_dir: /app
    environment:
      - JAVA_OPTS=-Xmx1024m  # 可以根据需要调整 Java 内存等参数
    volumes:
      - /8lab/conf/draven/application-pro.properties:/app/application-pro.properties
      - /8lab/draven/draven-dsl-engine-0.0.1-log.jar:/app/draven-dsl-engine-0.0.1-log.jar  
    command: >
      java -jar /app/draven-dsl-engine-0.0.1-log.jar
      --server.port=8097
      --spring.config.location=/app/application-pro.properties
      --spring.profiles.active=pro
```



### start && stop & log

```bash
# start 
docker-compose -f /8lab/draven/docker-compose.yml up -d

# stop 
docker-compose -f /8lab/draven/docker-compose.yml down 

# log
docker logs -fn1 draven
```



## Service Restart

### master

```bash
# 启动 draven
cd /8lab/draven && sh start.sh

# 检查 dtamper_server 状态
docker ps -a | grep dtamper | awk '{print $5, $6, $7, $NF}'  # 会有如下结果
  ``` 
  16 minutes ago dtamper_restful_service
  16 minutes ago dtamper_task_monitor
  16 minutes ago dtamper_heatbeat_monitor
  16 minutes ago dtamper_task_dispatcher 
  ```
# 


# 检查 dtamper_client.service 状态
systemctl status dtamper_client.service

# 检查 dtamper_path_monitor.service 状态
systemctl status dtamper_path_monitor.service

# 
```







## Other-info

### default web user

```bash
http://<ip>:8099
15948342592 testtest234 
```



### docker hub

```bash
https://hub.docker.com

1017127423@qq.com
119881220fang#
```

