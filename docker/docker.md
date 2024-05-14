# docker

## Install

```bash
apt install docker.io 
```



## default info

- setting file 
  -  /etc/docker/ 目录下，其中包括 daemon.json 文件

- Images volumes 
  -  /var/lib/docker/ 
  - /var/lib/docker/image/ 目录存储镜像的元数据和层信息，
  - /var/lib/docker/overlay2/ 目录存储容器的可写层。`
- Volume_data
  - /var/lib/docker/volumes
- Container_logs
  - /var/lib/docker/containers

- Network_type

  - host 
    - 共用宿主机host的网络资源，容器的网络配置将和宿主机host一模一样，优点就是传输效率高
    - 缺点就是容易和host造成端口冲突， 不进行dns解析

  - bridge
    -  系统默认会有docker0， 目的是为了实容器与宿主机的网络连接
  - None
    - 一个隔离的容器，不与外界进行交互。

- 

## settting

### Daemon.json

```bash
vi /etc/docker/daemon.json
{******其他配置
   "insecure-registries": ["192.168.66.102:85"]
   你的Harbor地址
}
```

### enable&start

```bash
systemctl restart docker
systemctl enable docker
```



### login

```bash
docker login -u 用户名 -p 密码 IP # ip为Harbor地址
```



## command

### network

- 创建网络 `docker network create -d bridge <netNasme> `
- 查看网络 `docker network ls // 查看网络`
- 详情 ` docker network inspect  <netName>  `



### container 

查看帮助信息`docker --help`

#### run 

```go
docker run -d --name xxx -p 9001:9001 images_id bash
// 宿主机使用端口:容器内端口
```



#### ps

```go
docker ps | grep container_id 或 container_name

docker ps -a // 查看所有
```



#### exec

```go
docker exec -it <PODNAME> bash 
// 进入容器
```

#### log

```go
docker logs -f --tail=100 ContainerdId 或 contain_name
```



#### update

```go
docker container update --publish-add <host-port>:<container-port> <container-id> 

// 一个运行中的容器增加一个端口docker container update --publish-add 8082:8082 <container-id>
```

```go
// 给运行中的容器增加一个挂载目录
docker cp /path/to/local_dir <container_id_or_name>:/path/to/mounted_dir
```

```go
// 现有容器名称名称
docker rename <container_id_or_name> new_containerName

```



### images 

#### build 

```go
docker build -t imageName:ImageTag -f Dockerfile.base . 
```



#### pull 

```go
docker pull octahub.8lab.cn:5000/poc211/assets_found_core:t1103
```



#### commit

```go
docker commit containerId newImage:tag // 将现有容器 打成一个新的镜像
```



#### push 

```go
docker push octahub.8lab.cn:5000/poc211/assets_found_core:t1103
// 需要登录
```



#### prune

```go
/usr/bin/docker image prune -f
// 清除镜像状态为 none 的镜像
```



### command_Detail

```go
attach    Attach to a running container                 # 当前 shell 下 attach 连接指定运行镜像
build     Build an image from a Dockerfile              # 通过 Dockerfile 定制镜像
commit    Create a new image from a container changes   # 提交当前容器为新的镜像
cp        Copy files/folders from the containers filesystem to the host path   #从容器中拷贝指定文件或者目录到宿主机中
create    Create a new container                        # 创建一个新的容器，同 run，但不启动容器
diff      Inspect changes on a container's filesystem   # 查看 docker 容器变化
events    Get real time events from the server          # 从 docker 服务获取容器实时事件
exec      Run a command in an existing container        # 在已存在的容器上运行命令
export    Stream the contents of a container as a tar archive   # 导出容器的内容流作为一个 tar 归档文件[对应 import ]
history   Show the history of an image                  # 展示一个镜像形成历史
images    List images                                   # 列出系统当前镜像
import    Create a new filesystem image from the contents of a tarball # 从tar包中的内容创建一个新的文件系统映像[对应export]
info      Display system-wide information               # 显示系统相关信息
inspect   Return low-level information on a container   # 查看容器详细信息
kill      Kill a running container                      # kill 指定 docker 容器
load      Load an image from a tar archive              # 从一个 tar 包中加载一个镜像[对应 save]
login     Register or Login to the docker registry server    # 注册或者登陆一个 docker 源服务器
logout    Log out from a Docker registry server          # 从当前 Docker registry 退出
logs      Fetch the logs of a container                 # 输出当前容器日志信息
port      Lookup the public-facing port which is NAT-ed to PRIVATE_PORT    # 查看映射端口对应的容器内部源端口
pause     Pause all processes within a container        # 暂停容器
ps        List containers                               # 列出容器列表
pull      Pull an image or a repository from the docker registry server   # 从docker镜像源服务器拉取指定镜像或者库镜像
push      Push an image or a repository to the docker registry server    # 推送指定镜像或者库镜像至docker源服务器
restart   Restart a running container                   # 重启运行的容器
rm        Remove one or more containers                 # 移除一个或者多个容器
rmi       Remove one or more images             # 移除一个或多个镜像[无容器使用该镜像才可删除，否则需删除相关容器才可继续或 -f 强制删除]
run       Run a command in a new container              # 创建一个新的容器并运行一个命令
save      Save an image to a tar archive                # 保存一个镜像为一个 tar 包[对应 load]
search    Search for an image on the Docker Hub         # 在 docker hub 中搜索镜像
start     Start a stopped containers                    # 启动容器
stop      Stop a running containers                     # 停止容器
tag       Tag an image into a repository                # 给源中镜像打标签
top       Lookup the running processes of a container   # 查看容器中运行的进程信息
unpause   Unpause a paused container                    # 取消暂停容器
version   Show the docker version information           # 查看 docker 版本号
wait      Block until a container stops, then print its exit code   # 截取容器停止时的退出状态值
```



## 启动docker设置登录

### Dockerfile

```dockerfile
# 使用基础镜像
FROM ubuntu:22.04

# 更新软件包列表
RUN apt-get update

# 安装 JDK 8
RUN apt-get install -y openjdk-17-jdk ssh

# 设置 JAVA_HOME 环境变量
ENV JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64/

# 设置 PATH 环境变量
ENV PATH="$JAVA_HOME/bin:$PATH"

# 验证 Java 安装
RUN java -version

# 设置容器启动时的默认命令
CMD ["bash"]
```



### Dockerfile

```go
FROM ubuntu:22.04

# 安装必要的软件包
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    openjdk-17-jdk \
    openssh-server && \
    rm -rf /var/lib/apt/lists/*

# 配置 SSH
RUN mkdir /var/run/sshd
RUN echo 'root:octa8lab123' | chpasswd
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN echo 'PermitEmptyPasswords yes' >> /etc/ssh/sshd_config

# 配置 SSH 映射端口
RUN sed -i 's/#Port 22/Port 22/' /etc/ssh/sshd_config

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D"]
```





### build && start

```go
//  build
docker build -t ubuntu22jdk17 .

// 运行容器
docker run -itd -p 9009:22 --name=ubuntu22jdk17 ubuntu22jdk17
```



### container set sshd_config

```go
// 进入启动的容器
docker exec -it ubuntu22jdk17 /bin/bash

// 设置密码， 登录后是root用户
passwd 你的密码， 输入两次

// 配置如下
vim /etc/ssh/sshd_config
PermitRootLogin yes

// 重载服务
service ssh reload


```

### connect_container

```go
sshpass -p "设置的密码" -p 9009 root@宿主机ip
```



## docker-compose

```bash
version: '3'
networkss:
  default:
    name: my_network

services:
  myapp-db:
    container_name: myapp-pg-db
    image: postgres:15
    restart: always
    ports:
      - 5432:5432
    env_file:
      -  ./myapp-env.env
    environment:
      - POSTGRES_USER=myapp_db_user
      - POSTGRES_PASSWORD=myapp_db_password
      - POSTGRES_DB=myapp_db
    volumes:
      - /data/pg/data/myapp_db:/var/lib/postgresql/data
    command: [ "postgres", "-N", "500" ]

  myapp-backend:
    image: myapp-backend:latest
    ports:
      - 9001:9001
    volumes:
      - ./data/log:/app/log/backend
    logging:     # 日志输入信息
      driver: json-file  
      options:
        max-size: 100m
        max-file: 3
        path: /app/log/backend/backend.log 
        
  myapp-backend:
    depends_on:  # 这两个启动后， 才会启动这个
      - myapp-backend
      - myapp-db
    image: myapp-frontend:latest
    ports:
      - 9001:9001
    volumes:
      - ./data/log:/app/log/frontend
    logging:
      driver: json-file  
      options:
        max-size: 100m
        max-file: 3
        path: /app/log/frontend/frontend.log 
        
    
```

myapp-env.env

```go
# myapp-env.env
POSTGRES_USER=myapp_db_user
POSTGRES_PASSWORD=myapp_db_password
POSTGRES_DB=myapp_db
```



### 启动停止

- 启动 `docker-compose up -d `

- 停止 `docker-compose down` 
- 重启 `docker-compose restart`
