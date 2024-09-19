

## 部署过程

> 直接将8lab目录下所有文件 移至 服务器 /8ab/下即可
>
> 直接将8lab目录下所有文件 移至 服务器 /8ab/下即可
>
> 直接将8lab目录下所有文件 移至 服务器 /8ab/下即可

### 创建并开启ima度量

#### Ima 功能配置

> 添加至 /boot/grub/grub.cfg 文件末尾
>
> ***\*注意：\****
>
> ***\*如果文件\****/boot/grub/grub.cfg不存在， 则需要查看是否存在/boot/grub2/grub.cfg 这个文件，存在则追加文件内容

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

### 配置文件目录结构

```bash
/8lab/conf
├── draven
│   └── application-pro.properties
├── dtamper_client
│   └── client_configure.json
├── dtamper_server
│   └── server_configure.json
├── oat_agent
│   ├── logging-linux.yaml
│   └── tagent.ini
├── oat_server
│   ├── logging-linux.yaml
│   └── oat.yaml
├── octa_cis
│   ├── conf.json
│   └── rpc.json
└── octa_cis_web
    └── octa_cis.conf
```





### mongo

> 部署： master节点， 采用集群架构
>
> 端口：28000

```bash
/8lab/mongodb
├── README.md
├── label
├── mongo-cm.yaml
├── mongo-job.yaml
├── mongo-svc.yaml
└── mongo.yaml
```



#### 添加label

> 给节点添加标签

```bash
# /8lab/mongodb/label
kubectl label node master node1 node2 mongo=mongo 
```

#### mongo-cm.yaml

> /8lab/mongodb/mongo-cm.yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mongo-init
  namespace: csmp
data:
  init.sh: |
    #!/bin/bash
    MONGO="$(which mongo)"
    PORT="$(grep -w "port" /usr/local/mongodb/etc/mongodb.cfg|awk '{print $2}')"
    RES_NAME=$(grep -w "replSetName" /usr/local/mongodb/etc/mongodb.cfg|awk '{print $2}')


    RS_STATUS=$(mongo --host mongo-0.mongo-svc:${PORT} --eval "rs.status()"|grep -w NotYetInitialized|awk -F [\"] '{print $(NF-1)}')
    if [ $RS_STATUS == NotYetInitialized ]; then
             echo "The cluster initialization starts."
             for ((i=0; i<=$RES_NUM-1; i ++)); do
                      MEMBERS+={_id:$i,host:\""mongo-${i}.mongo-svc.${NAMESPACE}.svc.cluster.local:28000\""},
             done
             $MONGO --host mongo-0.mongo-svc:${PORT} <<EOF
             rs.initiate({_id:"${RES_NAME}",members:[${MEMBERS}]})
             rs.status()
             exit;
    EOF
    else
             for ((i=0; i<=${RES_NUM}-1; i ++)); do
                      ${MONGO} --host mongo-${i}.mongo-svc:${PORT} --eval "rs.status()"|grep -C3 PRIMARY
                      if [ $? -eq 0 ]; then
                              PRIMARY=$(mongo --host mongo-${i}.mongo-svc:${PORT} --eval "rs.status()"|grep -C3 PRIMARY|awk -F[\"] 'NR==1{print $(NF-1)}')
                              echo $PRIMARY
                              break
                      fi
             done
             CS_NUM=$(${MONGO} --host $PRIMARY --eval "rs.status()"|grep -w "_id"|wc -l)
             if [ ${RES_NUM} -gt ${CS_NUM} ];then
             for ((i=0; i<=${RES_NUM}-1; i ++)); do
                      ${MONGO} --host mongo-${i}.mongo-svc:${PORT} --eval "rs.status()"|egrep -iw "Primary|Secondary|Recovering|Arbiter|STARTUP2"
                      if [ $? -eq 0 ]; then
                               echo "mongo-${i} already exists in the Replicaset. Not adding..."
                      else
                               $MONGO --host $PRIMARY --eval "rs.add('mongo-${i}.mongo-svc.${NAMESPACE}.svc.cluster.local:28000')"
                      fi
             done
             elif [ ${RES_NUM} -lt ${CS_NUM} ];then
               for ((i=${RES_NUM}; i<${CS_NUM}; i ++)); do
                        $MONGO --host $PRIMARY --eval "rs.remove('mongo-${i}.mongo-svc.${NAMESPACE}.svc.cluster.local:28000')"
               done
             fi
    fi

  mongodb.cfg: |
    storage:
        journal:
            enabled: true
            commitIntervalMs: 100
        dbPath: /data/mongodb/data
        directoryPerDB: false
        engine: mmapv1
        syncPeriodSecs: 60
        mmapv1:
            quota:
                enforced: false
                maxFilesPerDB: 8
            smallFiles: true
        wiredTiger:
            engineConfig:
                cacheSizeGB: 8
                journalCompressor: snappy
                directoryForIndexes: false
            collectionConfig:
                blockCompressor: snappy
            indexConfig:
                prefixCompression: true
    net:
        bindIp: 0.0.0.0
        port: 28000
        maxIncomingConnections: 65536
        wireObjectCheck: true
    operationProfiling:
        slowOpThresholdMs: 100
        mode: off
    replication:
        oplogSizeMB: 1000
        replSetName: 8lab-shard
```

#### mongo-job.yaml

> /8lab/mongodb/mongo-job.yaml

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: mongo-job
  namespace: csmp
spec:
  #ttlSecondsAfterFinished: 0  #job 执行完后等待100s 删除
  backoffLimit: 1  # 重试次数
  template:
    spec:
      restartPolicy: Never
      imagePullSecrets:
      - name: my-registry-secret
      containers:
      - name: mongo-init
        image: docker.io/1017127423/8lab:mongo_v210823
        command:
        - bash
        - -x
        - /usr/local/mongodb/bin/init.sh
        env:
        - name: RES_NUM
          value: "3"
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        volumeMounts:
        - name: config
          mountPath: /usr/local/mongodb/bin/init.sh
          subPath: init.sh
      volumes:
      - name: config
        configMap:
          name: mongo-init
```

#### mongo-svc.yaml

> /8lab/mongodb/mongo-svc.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mongo-svc
  namespace: csmp
  labels:
    app: mongo
spec:
  ports:
  - name: mongo
    port: 28000
    targetPort: 28000
    protocol: TCP
  selector:
    app: mongo
```

#### sts.yaml

> /8lab/mongodb/mongo.yaml

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mongo
  namespace: csmp
spec:
  selector:
     matchLabels:
        app: "mongo"
  serviceName: "mongo-svc"
  replicas: 3
  template:
    metadata:
      labels:
        app: mongo
    spec:
      imagePullSecrets:
      - name: my-registry-secret
      containers:
      - name: mongodb
        image: docker.io/1017127423/8lab:mongo_v210823
        command: ["/bin/sh", "-c"]
        args: ["sleep 10 && mongod -f /usr/local/mongodb/etc/mongodb.cfg"]
        ports:
        - containerPort: 28000
          #hostPort: 28000
          name: mongo
        volumeMounts:
        - name: mongo-data
          mountPath: /data/mongodb/data
        - name: config
          mountPath: /usr/local/mongodb/etc/mongodb.cfg
          subPath: mongodb.cfg
      volumes:
      - name: config
        configMap:
          name: mongo-init
          items:
          - key: mongodb.cfg
            path: mongodb.cfg
      - name: mongo-data
        hostPath:
          path: /data/mongo/data
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - mongo
            topologyKey: "kubernetes.io/hostname"
```

#### 启动& 停止

```bash
# 启动
kubectl apply -f .

# 停止
kubectl delete -f .

# 状态
kubectl get pod -n csmp 

# 查看集群信息
kubectl exec -it mongo-0 -- mongo
> use admin
switched to db admin
> rs.status()
```



### bigchainDb

> 部署： master节点
>
> 端口：32146

```bash
8lab/bigchainDb
├── first_deploy_yam.sh
├── octachain-conf
│   └── octachain.cfg
├── octachain-deployment.yaml
└── octachain-service.yaml
```



#### first_deploy_yam.sh

>  /8lab/bigchainDb/first_deploy_yam.sh

```bash
kubectl create configmap octachain-conf -n csmp --from-file=./octachain-conf/octachain.cfg
```



#### octachain.cfg

> /8lab/bigchainDb/octachain-conf/octachain.cfg

```bash
{
    "database": {
        "certfile": null,
        "host": "mongo-0.mongo-svc",
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

#### octachain-service.yaml

> 8lab/bigchainDb/octachain-service.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  name: octachain-service
  namespace: csmp
spec:
  type: NodePort
  ports:
  - port: 10070
    targetPort: 10070
    name: "octachainport"
    protocol: TCP
    nodePort: 32145
  - port: 10071
    targetPort: 10071
    name: "octawsport"
    protocol: TCP
    nodePort: 32146
  selector:
    app: octachain
```

#### octachain-deployment.yaml

> 8lab/bigchainDb/octachain-deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: octachain-deployment
  namespace: csmp
spec:
  selector:
    matchLabels:
      app: octachain
  replicas: 1
  template:
    metadata:
      labels:
        app: octachain
    spec:
      restartPolicy: Always
      imagePullSecrets:
      - name: my-registry-secret
      containers:
      - name: octachain
        image: docker.io/1017127423/8lab:octachain_v1020
        #command:
        #- sh
        #- -c
        #- sleep 999
        ports:
        - containerPort: 10070
          name: octachainport
        - containerPort: 10071
          name: octawsport
        volumeMounts:
          - mountPath: /var/log/8lab/
            name: octachain-logs
          - mountPath: /etc/octachain
            name: octachain-conf
      volumes:
      - name: octachain-logs
        hostPath:
          path: /data/logs/octachain/
      - name: octachain-conf
        configMap:
          name: octachain-conf
          items:
          - key: octachain.cfg
            path: octachain.cfg
```

#### Run&Stop

```bash
# 启动
kubectl apply -f . 

# 停止
kubectl delete -f . 

# 查看状态
kubectl get pod -n csmp 
```



### mysql

> 部署： master节点
>
> 端口：32143

```bash
/8lab/mysql
├── mysql-deploy.yaml
├── mysql-exporter-cm.yaml
├── mysql-master-cm.yaml
├── mysql-svc.yaml
└── sqlfiles
    ├── nisa.sql
    ├── octa_cis0607.sql
    ├── ueba_web.sql
    └── web_tamper_proof_db.sql
```



#### mysql-exporter-cm.yaml

> 8lab/mysql/mysql-exporter-cm.yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mysql-exporter-cm
  namespace: default
data:
  .my3306.cnf: |
    [client]
    port  = 3306
    user=8lab
    password=RZ8LYbh
    host=127.0.0.1
```

#### mysql-master-cm.yaml

> 8lab/mysql/mysql-exporter-cm.yaml

```yaml
apiVersion: v1
data:
  my3306.cnf: |
    [client]
    port            = 3306
    socket          = /data/mysql3306/tmp/mysql.sock
    user=mysql
    password=mysql

    [mysqld]
    server-id       = 99
    port            = 3306
    user            = mysql
    basedir         = /usr/local/mysql
    datadir         = /data/mysql3306/data
    tmpdir          = /data/mysql3306/tmp
    socket          = /data/tmp/mysql.sock
    pid-file        = /data/tmp/mysql.pid
    log-bin         = /data/mysql3306/logs/bin-log
    log-error       = /data/mysql3306/tmp/err.log
    slow_query_log_file = /data/mysql3306/tmp/slow.log
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
    innodb_data_home_dir = /data/mysql3306/data
    innodb_log_group_home_dir = /data/mysql3306/logs
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

    ## WSREP options
    #wsrep_provider=/usr/lib64/galera/libgalera_smm.so
    #wsrep_cluster_name="bjbs_citydb"
    #wsrep_cluster_address="gcomm://"
    #wsrep_cluster_address="gcomm://10.1.1.71,10.1.1.72,10.1.1.73"
    #wsrep_slave_threads=32
    #wsrep_node_address=10.1.1.71
    #wsrep_sst_method=xtrabackup
    #wsrep_sst_auth=sstuser:3tWbcWD5L6Z99s3h

    #slave_skip_errors = 1032

    [mysql]
    no-auto-rehash
kind: ConfigMap
metadata:
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:data:
        .: {}
        f:my3306.cnf: {}
    manager: kubectl
    operation: Update
  name: mysql-master
  namespace: default
```

#### mysql-svc.yaml

> 8lab/mysql/mysql-exporter-svc.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
  annotations:
    prometheus.io/http-probe: "true"
    prometheus.io/http-probe-port: "8080"
    prometheus.io/http-probe-path: "/healthCheck"
spec:
  type: NodePort
  ports:
  - port: 3306
    targetPort: 3306
    name: mysql
    protocol: TCP
    nodePort: 32143
  selector:
    app.oscro.io/name: mysql
```

#### mysql-deploy.yaml

> 8lab/mysql/mysql-exporter-deploy.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app.oscro.io/name: mysql
  replicas: 1
  template:
    metadata:
      labels:
        app.oscro.io/name: mysql
    spec:
      imagePullSecrets:
      - name: registry-secret
      containers:
      - name: mysql
        image: octahub.8lab.cn:5000/octa-cis/mysql:v0826
        #imagePullPolicy: Never
        ports:
        - containerPort: 3306
        volumeMounts:
        - name: host-time
          mountPath: /etc/localtime
          readOnly: true
        - name: mysql-tmp
          mountPath: /data/tmp
        - name: mysql-data
          mountPath: /data/mysql3306
        - name: mysql-conf
          mountPath: /opt/my3306.cnf
          subPath: my3306.cnf
        resources:
          requests:
            cpu: "1"
            memory: "6Gi"
          limits:
            cpu: "4"
            memory: "8Gi"
      volumes:
      - name: host-time
        hostPath:
          path: /etc/localtime
      - name: mysql-tmp
        emptyDir: {}
      - name: mysql-data
        hostPath:
          path: /data/mysql3306/data
      - name: mysql-logs
        hostPath:
          path: /data/mysql3306/logs
      - name: mysql-conf
        configMap:
          name: mysql-master

      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/hostname
                operator: In
                values:
                - master
```



#### Run&Stop

```bash
# 启动
kubectl apply -f . 

# 停止
kubectl delete -f . 
```

#### initDatabase

> 初始化sql文件在 /8lab/mysql/sqlfiles 下

```bash
# 初始化数据库

# 进入容器  
kubectl exec -it xx bash 

# 连接数据库
mysql -S /data/tmp/mysql.sock

# 创建库
create database nisa;
create database ueba_web;
create database web_tamper_proof_db;
create database octa;

# 导入表数据
use web_tamper_proof_db
source web_tamper_proof_db.sql 

# 
use nisa
source nisa.sql

# 
use octa_cis
source octa_cis0607.sql 

# 
use ueba_web
source ueba_web.sql 
```



### redis

> 部署： master节点
>
> 端口：32144

```bash
/8lab/redis
├── redis-cm.yaml
├── redis-deploy.yaml
└── redis-svc.yaml
```



#### redis-cm.yaml

```yaml
# vim /8lab/redis/redis-cm.yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: redis-master-cm
  namespace: default
data:
    redis.conf: |
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

#### redis-deploy.yaml

```bash
# vim /8lab/redis/redis-deploy.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-master-deploy
  namespace: default
  labels:
    app: redis-master
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis-master
  template:
    metadata:
      labels:
        app: redis-master
    spec:
      restartPolicy: Always
      securityContext:
        sysctls:
        - name: net.core.somaxconn
          value: "65535"
      initContainers:
      - name: never-thp
        image: octahub.8lab.cn:5000/octa-cis/busybox:1.28
        imagePullPolicy: IfNotPresent
        securityContext:
          privileged: true
        command: ["sh", "-c", "echo never > /host-sys/kernel/mm/transparent_hugepage/enabled"]
        volumeMounts:
        - name: host-sys
          mountPath: /host-sys
      imagePullSecrets:
      - name: registry-secret
      containers:
      - name: redis
        image: octahub.8lab.cn:5000/octa-cis/redis:v0421
        ports:
        - containerPort: 6379
          name: redis
        resources:
          limits:
            cpu: 300m
            memory: 200Mi
          requests:
            cpu: 200m
            memory: 200Mi
        volumeMounts:
        - name: redis-conf
          mountPath: /usr/local/redis/redis.conf
          subPath: redis.conf
        - name: redis-data
          mountPath: /data/redis
      volumes:
      - name: host-sys
        hostPath:
          path: /sys
      - configMap:
          items:
          - key: redis.conf
            path: redis.conf
          name: redis-master-cm
        name: redis-conf
      - name: redis-data
        hostPath:
          path: /data/redis
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/hostname
                operator: In
                values:
                - master
```

#### redis-svc.yaml

```yaml
# vim /8lab/redis/redis-svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-master-svc
  namespace: default
spec:
  type: NodePort
  ports:
  - port: 6379
    targetPort: 6379
    name: "redis"
    protocol: TCP
    nodePort: 32144
  selector:
    app: redis-master
  # externalIPs:
  # - 192.168.110.131
```

#### Run&Stop

```bash
# 启动
kubectl apply -f . 

# 停止
kubectl delete -f . 

# 状态
kubectl get pod -n default
```



### kafka

> 部署： master节点
>
> 端口：9092

```bash
kafka
└── docker-compose.yml
```



#### data

> 数据存储 目录， 需要给的777权限

```bash
mkdir -p /data/kafka/data  && chmod 777 -R /data/kafka/data  
```



#### docker-compose.yml

> /8lab/kafka/docker-compose.yml

```bash
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
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://192.168.110.131:9092
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_NODE_ID=1
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@localhost:9093
      - KAFKA_CFG_LOG_DIRS=/tmp/kraft-combined-logs
    volumes:
      - /data/kafka/data:/tmp/kraft-combined-logs
```



#### topic 

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


kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group GID_STD_CSMP_SANDBOXAUDIT_LOG_002 --topic csmp-std-sandboxaudit-log

kafka-consumer-groups.sh --bootstrap-server localhost:9092 --group GID_STD_CSMP_SANDBOXAUDIT_LOG_002 --describe --topic 
```

#### demo-test.py

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



### auditd

> 部署：主从节点

#### install

```bash
# 安装 auditd
apt install auditd
```

#### policy

```bash
# 规则文件解压到服务器
rm -rf /etc/audit

# /xx/xxx/audit_.tar.gz 注意路径，文件在
tar -xzf audit.tar.gz -C /etc/ 

# 这里会有很多规则文件
```

#### Run&&Stop

```bash
systemctl reload auditd

# 开机启动
systemctl enable auditd

# 启动
systemctl start auditd

# 停止
systemctl stop auditd
```



### filebeat

> 部署： master node1 node2

```bash
8lab/filebeat
├── apt_update.sh
└── filebeat
    └── filebeat-7.17.18-x86_64.rpm
```

#### Install && Run 

```bash
cd /8lab/filebeat/ 
sh apt_update.sh 
```



### logstash

> 部署： mster

```bash
# 服务说明
csmp-ty-sysaudit-log logstash-5145 input filebeat写入kafka 
csmp-output-sysaudit-log logstash-5145 out 到kafka

csmp-output-sysaudit-log logstash-5156 input 从kafka 读
csmp-std-sandboxaudit-log logstash-5156 output

csmp-std-sandboxaudit-log  draven 读
```

#### log5145

```bash
/8lab/log5145
├── csmp-log5145-deploy.yaml
└── csmp-log5145-getlogs-configs.yaml
```



##### csmp-log5145-getlogs-configs.yaml

> /8lab/log5145/csmp-log5145-getlogs-configs.yaml

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: csmp-log5145-getlogs-configs
  namespace: csmp
data:
  # logstash.conf: |
  logstash-sample.conf: |
    input {
        kafka {
            bootstrap_servers => "192.168.110.131:9092"
            topics => ["csmp-ty-sysaudit-log"]
            group_id => csmp
            codec => json {
             charset => "UTF-8"
            }
            partition_assignment_strategy => "round_robin"

        }
    }
    filter {
        drop {
                percentage => 0
             }
        mutate {
          add_field => {
              "message" => "%{[content][output]}"
              "source" => "%{filename}"
              "host" => "%{sysname}"
              "type_ip" => "%{ipaddress}"
          }
          rename => ["type","type_ip"]
          remove_field => ["type","content","kafka.consumer_group","kafka.key","kafka.offset","kafka.partition","kafka.topic","filename","topic","modelName","sysname","ipaddress"]
        }
        eaglebrief {}
        if ([rdsmsg] ==''){
            drop {}
        }
        split {
            field => "message"
            terminator => "#split#"
        }
    }

    output {
        stdout {}
        kafka {
            codec => json {
              charset => "UTF-8"
            }
            bootstrap_servers => ["192.168.110.131:9092"]
            topic_id => "csmp-output-sysaudit-log"
        }
    }
```

##### csmp-log5145-deploy.yaml

> /8lab/log5145/csmp-log5145-deploy.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: csmp-log5145-deploy
  namespace: csmp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: csmp-log5145-deploy
  template:
    metadata:
      labels:
        app: csmp-log5145-deploy
    spec:
      containers:
        - env:
            - name: HTTP_HOST
              value: 0.0.0.0
            - name: HTTP_PORT
              value: '9600'
            - name: LANG
              value: zh_CN.utf8
          image: 1017127423/8lab:logstash_update3_8.9.0
          name: rhel-logstash-sts
          resources:
            limits:
              cpu: '2'
              memory: 2Gi
            requests:
              cpu: '1'
              memory: 1Gi
          volumeMounts:
            - mountPath: /usr/share/logstash/pipeline/
              name: csmp-log5145-getlogs-configs
      hostAliases:
        - hostnames:
            - ntf-nrdp-mn-01.bigdata.com
          ip: 15.5.4.51
        - hostnames:
            - ntf-nrdp-mn-02.bigdata.com
          ip: 15.5.4.52
        - hostnames:
            - ntf-nrdp-kafka-01.bigdata.com
          ip: 15.5.4.53
        - hostnames:
            - ntf-it-nrdp-mn.cebbank.com
          ip: 15.5.4.77
        - hostnames:
            - ntf-it-nrdp-kafka-01.cebbank.com
          ip: 15.5.4.78
        - hostnames:
            - ntf-it-nrdp-kafka-02.cebbank.com
          ip: 15.5.4.79
        - hostnames:
            - ntf-it-nrdp-kafka-03.cebbank.com
          ip: 15.5.4.80
      imagePullSecrets:
        - name: my-registry-secret
      volumes:
        - configMap:
            defaultMode: 420
            name: csmp-log5145-getlogs-configs
          name: csmp-log5145-getlogs-configs
```



#### log5156

##### csmp-log5156-deploy.yaml

> /8lab/log5156/csmp-log5156-deploy.yaml

````yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: csmp-log5156-deploy
  namespace: csmp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: csmp-log5156-deploy
  template:
    metadata:
      labels:
        app: csmp-log5156-deploy
    spec:
      containers:
        - env:
            - name: HTTP_HOST
              value: 0.0.0.0
            - name: HTTP_PORT
              value: '9600'
            - name: LANG
              value: zh_CN.utf8
          image: 1017127423/8lab:logstash_update3_8.9.0
          imagePullPolicy: IfNotPresent
          name: rhel-logstash-sts
          resources:
            limits:
              cpu: '4'
              memory: 5Gi
            requests:
              cpu: '2'
              memory: 3Gi
          volumeMounts:
            - mountPath: /usr/share/logstash/pipeline/
              name: csmp-log5156-getlogs-configs
      imagePullSecrets:
        - name: my-registry-secret
      volumes:
        - configMap:
            defaultMode: 420
            name: csmp-log5156-getlogs-configs
          name: csmp-log5156-getlogs-configs
````

##### csmp-log5156-getlogs-configs.yaml

> /8lab/log5156/csmp-log5156-getlogs-configs.yaml

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: csmp-log5156-getlogs-configs
  namespace: csmp
data:
  logstash.conf: |
    input {
        kafka {
            bootstrap_servers => "192.168.110.131:9092"
            topics => "csmp-output-sysaudit-log"
            group_id => "csmp"
            codec => json {
             charset => "UTF-8"
            }
            type => "kafka"
            client_id => logstash0
            consumer_threads => 400
            fetch_min_bytes => 16768
            session_timeout_ms => 20000
            max_poll_interval_ms => 480000
            request_timeout_ms => 90000
            partition_assignment_strategy => "round_robin"
        }
    }

    filter {
        if [source] == "/var/log/ncolog" {
          mutate {
            remove_field => "type_ip"
          }
          grok {
            match => {
              "message" => "%{IPV4:type_ip}"
            }
          }
        }
        ruby {
            code => "event.set('timestamp', event.get('@timestamp').time.localtime + 8*60*60)"
        }
        eaglesys {}
        if ([rdsmsg] =='') {
          drop {}
        }
    }

    output {
       stdout {}
       kafka {
            codec => plain{format => "%{rdsmsg}"}
            bootstrap_servers => ["192.168.110.131:9092"]
            topic_id => "csmp-std-sandboxaudit-log"
        }
    }
```



### dtamper_server防篡改

> 部署： master节点

```bash
/8lab/dtamper_server
└── docker-compose.yml
```



#### server_configure.json

> d/server_configure.json

```bash
{
  "service_type": ["web"],
  "client_address": "192.168.110.131",
  "server_address": "192.168.110.131",
  "bdb_host":"192.168.110.131",
  "bdb_port":"32145",
  "redis_host": "192.168.110.131",
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
  "mysql_host": "192.168.110.131",
  "mysql_port": "32143",
  "mysql_username": "8lab",
  "mysql_password": "8lab",
  "token_switch": "off"
}
```

#### docker-compose.yml

> /8lab/dtamper_server/docker-compose.yml

```bash
version: '3'
services:
  restful_service:
    image: octahub.8lab.cn:5000/octa-cis/dtamper-server-4.3:202401231050
    container_name: restful_service
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



#### Run&stop

```bash
# 启动
docker-comepose up -d

# 停止
docker-compose down 
```



### dtamper_client防篡改

> 部署： master node1 node2

```bash
dtamper_client
├── client_db_generate_plain
├── dtamper_client
├── dtamper_path_monitor
├── key_pairs
│   ├── client_pri
│   ├── client_pub
│   ├── server_pri
│   └── server_pub
└── web_client.db  # 由 ./client_db_generate_plain 生成
```



#### client_configure.json

> /8lab/conf/dtamper_client/client_configure.json

```bash
{
  "service_type": "web",
  "client_address": "192.168.110.131",  # 当前主机
  "server_address": "192.168.110.131",  # dtamper_server 
  "bdb_host":"192.168.110.131",
  "bdb_port":"32145",
  "redis_host": "192.168.110.131",
  "redis_port": "32144",
  "redis_password": "octa8lab",
  "update_channel": "update",
  "report_channel": "return",
  "control_send_channel": "control_send",
  "control_return_channel": "control_return",
  "heartbeat_channel": "heartbeat",
  "aes_key": "ca$hc0w8L@6ExP0!",
  "aes_iv": "ca$hc0w8L@6ExP0!",
  "mysql_host": "192.168.110.131",
  "mysql_port": "32143",
  "mysql_username": "8lab",
  "mysql_password": "8lab"
}
```

####  web_client.db 

```bash
## 生成 web_client.db
cd /8lab/dtamper_client && ./client_db_generate_plain 

# 提示信息以下信息可忽略
# return
# mv: /tmp/_MEI7dbjaJ/libselinux.so.1: no version information available (required by mv)
# mv: 无法将 'web_client.db' 移动至 '../tamper_proof_main/': 不是目录
```



#### dtamper_path_monitor.service

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



#### dtamper_client.service

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



#### 启动并设置开机自启动

```bash
systemctl daemon-reload

# 设置为开机自启动
systemctl enable dtamper_path_monitor.service
systemctl enable dtamper_client.service

# 启动
systemctl start dtamper_path_monitor.service
systemctl start dtamper_client.service

# 停止 
systemctl stop dtamper_path_monitor.service
systemctl stop dtamper_client.service
```



### oat_server 可信Server

> 部署： master 节点

```bash
/8lab/oat_server
└── docker-compose.yml
```



#### logging-linux.yaml

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



#### oat.yaml

> /8lab/conf/oat_server/oat.yaml

```bash
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
    host: 192.168.110.131
    user: 8lab
    password: 8lab
    db_name: octa_cis
    port: 32143
  redis:
    host: 192.168.110.131
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



#### docker-compose.yml

> /8lab/oat_server/docker-compose.yml

```bash
ersion: '3'
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



### oat_agent  可信client

> 部署： master node1 node2 
>
> tagent.ini 配置中和网卡有绑定关系
>
> /8lab/logs/  日志目录

```bash
/8lab/tagent
├── agent_dbs.db
├── pytagent
├── start.sh
├── stop.sh
└── tmp
    └── alarm_file
```



#### logging-linux.yaml 

> /8lab/conf/oat_agent/logging-linux.yaml 

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
        filename: '/8lab/logs/oat_agent/tagent.log'
        when: 'MIDNIGHT'
        interval: 1
        backupCount: 20
        maxBytes: 209715200 #200M

loggers:
    root:
        level: DEBUG
        handlers: [fileHandler, consoleHandler,]
    tagent:
        level: DEBUG
        handlers: [fileHandler, consoleHandler,]
        qualname: tagent
```



#### tagent.ini

> /8lab/conf/oat_agent/logging-linux.yaml 

```bash
[agent]
; agent默认启动端口号
agent_port = 50052
; agent对外映射地址 没有为no， 启动异常无法获取ip时 可改为主机ip
agent_mapping_ip = no
; agent对外映射端口号 默认为0表示没有
agent_mapping_port = 0
; agent网卡名 获取指定网卡的ip 默认no_card
agent_network_card = ens2

[oat]
; 填写oat server对外映射地址
oat_ip = 192.168.110.131
oat_port = 50051

[open_options]
# 默认为1开启可信防御 0为关闭可信防御
trusted = 1
# 默认为1开启资产清点功能 0为关闭资产清点功能
assets_discover = 0
# 资产清点扫描的频率 以分钟为单位 默认为30分钟为一次
assets_discover_rate = 30

; 0 代表 开启redis传递模式 1 是开启RabbitMQ传递模式
[open_mq]
open_mq = 0

; 如果open_mq = 1时需配置[rabbitmq]填写地址等登录信息
[rabbitmq]
mq_ip = 192.168.3.100
mq_port = 5672
mq_user = admin
mq_pass = admin

# agent心跳配置
[heartbeat]
# agent发送心跳信号默认每隔10s发送一次，此处应与oat_server配置文件中heartbeat monitor_heartbeat值相同
hb_interval = 10

[whitelist]
;可自定义每次上传白名单及告警条数 单位是行
whitelist_num = 500
;可自定义检测告警频率 单位是秒，默认为10秒一次检测告警
detection_alarm_frequency = 10

[audit]
; audit_interval 0 表示不插入audit.log  其余数字表示时间间隔
audit_path = /var/log/audit/audit.log
audit_interval = 0

[limit]
; OFF 关  ON 开
open_limit = ON
; 上限
upper_limit = 20
; 下限
lower_limit = 5

[log]
; linux 日志文件路径
log_config_file_linux = /8lab/conf/oat_agent/logging-linux.yaml
; windows 日志文件路径
log_config_file_windows = C:\8lab\conf\oat_agent\logging-windows.yaml
log_dir_size = 300
```



#### start.sh

```bash
#!/bin/bash
nohup ./pytagent> /dev/null 2>&1 &
```



#### stop

```bash
#!/bin/bash
kill $(ps -ef | grep tagent | awk '{print $2}')
```



#### Run&stop

```bash 
#  启动
sh /8lab/tagent/start.sh

# 停止
sh /8lab/tagent/stop.sh

# 注
程序程序启动时：
- 配置文件目录:/8lab/conf/oat_agent/
- /8lab/logs/ /8lab/tagent/agent_dbs.db  程序启动时生成文件
```



### octa_cis_server  WebServer

> /8lab/octa_cis
> └── docker-compose.yml

#### conf.json

> /8lab/conf/octa_cis/conf.json

```bash
{
  "client_audit_hosts": [
    {
      "ip": "k8sproxy-service",
      "name": "client_b"
    }
  ],
  "product_uid": "bd43-430e-af62-e710",
  "server_ip": "192.168.110.131",
  "blackbox_ip": "192.168.1.169",
  "switch_waf_port": 8080,
  "es_server_ip_port": [{"host": "192.168.110.131", "port": 29200}],
  "es_server_user_name": "elastic",
  "es_server_password": "tarena",
  "bdb_host": "192.168.1.98",
  "bdb_port": "10070",
  "bdb_mongo_host": "192.168.1.98",
  "mongo_port": 28000,
  "rabbitmq_server": "192.168.1.184",
  "rabbitmq_port": 5672,
  "draven_server": "192.168.110.131",
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
  "mysql_host": "192.168.110.131",

  "nisa_mysql_user": "8lab",
  "nisa_mysql_password": "8lab",
  "nisa_mysql_database": "nisa",
  "nisa_mysql_port": 32143,
  "nisa_mysql_host": "192.168.110.131",

  "draven_mysql_user": "8lab",
  "draven_mysql_password": "8lab",
  "draven_mysql_database": "ueba_web",
  "draven_mysql_port": 32143,
  "draven_mysql_host": "192.168.110.131",

  "ueba_web_mysql_password": "8lab",
  "ueba_web_mysql_user": "8lab",
  "ueba_web_mysql_database": "ueba_web",
  "ueba_web_mysql_port": 32143,
  "ueba_web_mysql_host": "192.168.110.131",

  "alarm_second": 15,
  "alarm_enable": 1,
  "redis4bigchanidb_host":"192.168.110.131",
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
  "tamper_proof_url": "http://192.168.110.131:2345",
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

  "dtamper_redis_host": "192.168.110.131",
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

#### rpc.json

> /8lab/conf/octa_cis/rpc.json

```bash
{
  "user_portrait": {
    "host": "101.251.211.205",
    "port": 9898
  },
  "grpc_param": {
    "host": "192.168.110.131",
    "port": 50051
  }
}
```

#### docker-compose.yml

> /8lab/octa_cis/docker-compose.yml

```bash
# vim docker-compose.yml

version: '3'
services:
  octa_cis:
    image: octahub.8lab.cn:5000/octa-cis/octa_cis_fc:202402021108
    container_name: octa_cis
    ports:
      - "8001:8001"
    volumes:
      - /8lab/conf/octa_cis/:/usr/local/octa_cis/conf/
    sysctls:
      net.core.somaxconn: 65535
    restart: unless-stopped
```

#### Run&stop

```bash
# 启
docker-compose up -d 
# 停
docker-compose down 
```

### octa_cis_web Web页面

> 前端调整了web, 目录 app_views 为调整后的内容

```bash
├── app_views
│   └──....
│   └──....
└── docker-compose.yml
```



#### docker-compose.yml

> 8lab/octa_cis_web/docker-compose.yml

```bash
version: '3'
services:
  octa_cis_web:
    image: octahub.8lab.cn:5000/octa-cis/octa_cis_web_fc_dongwu:202402021616
    container_name: octa_cis_web
    ports:
      - "8099:8099"
    volumes:
      - /8lab/conf/octa_cis_web/:/etc/nginx/conf.d/
      - /data/yaml/octa_cis_web/app_views/:/usr/local/octa_cis_web_4_3/app_views/
    restart: unless-stopped
```



#### run&stop

```bash
# 启
docker-compose up -d 
# 停
docker-compose down 
```





### draven

> 手动启动

```bash
/8lab/draven
├── backup_draven_logs.sh
├── draven-dsl-engine-0.0.1-SNAPSHOT02.jar
├── draven-dsl-engine-0.0.1-log.jar
├── start.sh
└── stop.sh
```



#### application-pro.properties

> /8lab/draven/conf/draven/application-pro.properties

```bash
 kafka
kafka.consumer.servers=192.168.110.131:9092
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
spring.datasource.url=jdbc:mysql://192.168.110.131:32143/ueba_web?useSSL=false
spring.datasource.username=8lab
spring.datasource.password=8lab
spring.datasource.max-idle=10
spring.datasource.max-wait=10000
spring.datasource.min-idle=5
spring.datasource.initial-size=5

# ==================================== redis配置 ====================================
spring.redis.host=192.168.110.131
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



#### start.sh

> /8lab/draven/start.sh

```bash
#!/bin/bash
DIR="/8lab/logs/draven"
if [ ! -d "$DIR" ]; then
  mkdir -p "$DIR"
fi

nohup java -jar draven-dsl-engine-0.0.1-log.jar \
	--server.port=8097 \
	--spring.config.location=/8lab/conf/draven/application-pro.properties \
	--spring.profiles.active=pro \
	> $DIR/draven.log 2>&1 &



```

#### run&Stop

```bash
cd /8lab/draven

# 启动
sh start.sh

# 停止
sh /stop.sh

```



#### crontab

```bash
0 * * * * /8lab/logs/draven/draven.log
```

##### backup_draven_logs.sh

> /8lab/logs/draven/draven.log

```bash
# 备份日志文件， 每隔一小时进行一次备份， 删除2天以前的备份文件
#!/bin/bash

# 设置日志文件路径
LOG_DIR="/8lab/logs/draven"
LOG_FILE="${LOG_DIR}/draven.log"

# 备份日志的文件名，格式为 draven_YYYYMMDDHH.log
BACKUP_FILE="${LOG_DIR}/draven_$(date +%Y%m%d%H).log"

# 备份日志
cp $LOG_FILE $BACKUP_FILE

# 清空原始日志文件
cat /dev/null > $LOG_FILE

# 删除2天以前的备份文件
find $LOG_DIR -name "draven_*.log" -type f -mtime +2 -exec rm {} \;
```



## Other-info

### web

```bash
http://192.168.110.131:8099
15948342592 testtest234 
```

### docker.io

```bash
https://hub.docker.com

1017127423@qq.com
119881220fang#

```

### 8labhub.com

```bash
联系 运维
```

