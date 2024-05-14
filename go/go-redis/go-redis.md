# go-redis



## Redis支持的数据结构

Redis支持诸如字符串（strings）、哈希（hashes）、列表（lists）、集合（sets）、带范围查询的排序集合（sorted sets）、位图（bitmaps）、hyperloglogs、带半径查询和流的地理空间索引等数据结构（geospatial indexes）。



## 启动

### 单节点模式

```bash
docker run -d -p 6379:6379 --name my-redis -e REDIS_PASSWORD=8lab redis
```

```go
docker run：用于启动一个新的容器。
-d：将容器以后台（detached）模式运行。
-p 6379:6379：这里是将容器的 6379 端口映射到主机的 6379 端口。hostPort:containerPort
--name my-redis：为容器指定一个名称，这里是 "my-redis"。
-e REDIS_PASSWORD=8lab：通过环境变量 REDIS_PASSWORD 设置 Redis 的密码为 "8lab"。
- redis：指定要运行的镜像为 Redis。
```



### 哨兵模式

- docker-compose.yml

  ```go
  version: '3'
  services:
    redis-master:
      image: redis
      command: redis-server --port 6379
      ports:
        - 6379:6379
    redis-slave:
      image: redis
      command: redis-server --port 6380 --replicaof redis-master 6379
      ports:
        - 6380:6380
    redis-sentinel:
      image: redis
      command: redis-sentinel /usr/local/etc/redis/sentinel.conf
      ports:
        - 26379:26379
      volumes:
        - ./sentinel.conf:/usr/local/etc/redis/sentinel.conf
  
  ```

- sentinel.conf

  ```go 
  port 26379
  sentinel monitor mymaster redis-master 6379 2
  sentinel down-after-milliseconds mymaster 5000
  sentinel failover-timeout mymaster 10000
  ```



### 集群模式

```yaml
version: '3'
services:
  redis-1:
    image: redis
    command: redis-server --port 7000 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000 --appendonly yes
    ports:
      - 7000:7000
    volumes:
      - ./data/redis-1:/data
  redis-2:
    image: redis
    command: redis-server --port 7001 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000 --appendonly yes
    ports:
      - 7001:7001
    volumes:
      - ./data/redis-2:/data
  redis-3:
    image: redis
    command: redis-server --port 7002 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000 --appendonly yes
    ports:
      - 7002:7002
    volumes:
      - ./data/redis-3:/data
  redis-4:
    image: redis
    command: redis-server --port 7003 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000 --appendonly yes
    ports:
      - 7003:7003
    volumes:
      - ./data/redis-4:/data
  redis-5:
    image: redis
    command: redis-server --port 7004 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000 --appendonly yes
    ports:
      - 7004:7004
    volumes:
      - ./data/redis-5:/data
  redis-6:
    image: redis
    command: redis-server --port 7005 --cluster-enabled yes --cluster-config-file nodes.conf --cluster-node-timeout 5000 --appendonly yes
    ports:
      - 7005:7005
    volumes:
      - ./data/redis-6:/data
```





redis命令大全参考：http://doc.redisfans.com/

## 操作

###	下载go-redis依赖包

```go
go get -u github.com/go-redis/redis
```



### 新版本特性

最新版本的`go-redis`库的相关命令都需要传递`context.Context`参数

#### 作用

它提供了一种在 Redis 模块和 Redis 服务器之间进行数据交互和状态管理的机制

- 数据交互：`context` 对象允许模块访问和操作 Redis 数据库。通过 `context`，模块可以执行诸如读取、写入、修改数据等操作，以及执行 Redis 命令或调用 Redis API。

- 状态管理：`context` 对象可以用于跟踪模块的状态和上下文信息。在模块的不同函数调用之间，可以使用 `context` 对象来传递和保持状态信息，以便在后续调用中使用。

- 错误处理：`context` 对象提供了处理错误和异常情况的机制。模块可以通过 `context` 对象报告错误、记录日志和处理异常，以便与 Redis 服务器进行适当的交互和响应。

- 事件处理：`context` 对象可以用于订阅和处理 Redis 服务器的事件。模块可以通过 `context` 对象注册回调函数，以便在特定事件发生时被调用，从而实现对事件的处理和响应。



#### 示例

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

// connect redis
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "10.1.1.116:6379",
		Password: "8lab",
		DB:       0,
		PoolSize: 100,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}

func Example() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}

	setErr := rdb.Set(ctx, "key1", "value1", 0).Err()
	if setErr != nil {
		panic(setErr)
	}

	val, getErr := rdb.Get(ctx, "key1").Result()
	fmt.Println(111)
	if getErr != nil {
		panic(getErr)
	}
	fmt.Println(val)
}

func main() {
	Example()
}

```







### 连接

#### 单节点

```go
// connect redis
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "10.1.1.116:6379",
		Password: "8lab",
		DB:       0,
		PoolSize: 100,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}

```



#### 哨兵模式

```go
func initClient()(err error){
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master",
		SentinelAddrs: []string{"10.1.1.116:26379"},
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
```



#### 集群模式

```go
func initClient()(err error){
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
```



### 操作

#### get/set 

```GO
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

// connect redis
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "10.1.1.116:6379",
		Password: "8lab",
		DB:       0,
		PoolSize: 100,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}

func Example() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}

	setErr := rdb.Set(ctx, "key1", "value1", 0).Err()
	if setErr != nil {
		panic(setErr)
	}

	val, getErr := rdb.Get(ctx, "key1").Result()
	fmt.Println(111)
	if getErr != nil {
		panic(getErr)
	}
	fmt.Println(val)
}

func main() {
	Example()
}

```

#### get-prefix

按通配符获取key

```go
vals, err := rdb.Keys(ctx, "prefix*").Result()
```



#### del-prefix 

按通配符删除key

```go
ctx := context.Background()
iter := rdb.Scan(ctx, 0, "prefix*", 0).Iterator()
for iter.Next(ctx) {
	err := rdb.Del(ctx, iter.Val()).Err()
	if err != nil {
		panic(err)
	}
}
if err := iter.Err(); err != nil {
	panic(err)
}
```



#### Set 集合

```go

func Example() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}

	// set add
	err := rdb.SAdd(ctx, "s1", "sv1", "sv2", "sv3").Err()
	if err != nil {
		panic(err)
	}

	// set get all
	vl, _ := rdb.SMembers(ctx, "s1").Result()

	// exist
	ok, _ := rdb.SIsMember(ctx, "s1", "sv1").Result()
	if ok {
		fmt.Println("sv1 is in s1")
	} else {
		fmt.Println("sv1 is not in s1")
	}
	
	// set remove members
	_ = rdb.SRem(ctx, "s1", "sv1").Err()
	
	// del key
	r1, _ := rdb.Del(ctx, "s1").Result()
	fmt.Println(r1)

	fmt.Println(vl)
}

```

#### mget、mset 批量操作

批量获取/设置

```go
func Example() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}

	err := rdb.MSet(ctx, "k1", "v1", "k2", "v2").Err()
	if err != nil {
		panic(err)
	}
	
	val := rdb.MGet(ctx, "k1", "k2")

	for _, v := range val.Val() {
		fmt.Println("v:", v)
	}

	fmt.Println(val)
}

```

#### list 列表

```go
func Example() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}

	// left push list
	err := rdb.LPush(ctx, "l1", "1", "2", "3").Err()
	if err != nil {
		panic(err)
	}
	
	// right push list
	err = rdb.LPush(ctx, "l1", "a", "b", "c").Err()
	if err != nil {
		panic(err)
	}

	// left pop list
	if err = rdb.LPop(ctx, "l1").Err(); err != nil {
		panic(err)
	}

	// right pop list
	if err = rdb.RPop(ctx, "l1").Err(); err != nil {
		panic(err)
	}

	// get all
	vl, _ := rdb.LRange(ctx, "l1", 0, -1).Result()
  // len(list)
	countVl, _ := rdb.LLen(ctx, "l1").Result()

	
	fmt.Println(vl)
	fmt.Println(countVl)
}

```

#### 发布/订阅

redis本身具有发布订阅的功能，其发布订阅功能通过命令SUBSCRIBE(订阅)／PUBLISH(发布)实现，并且发布订阅模式可以是多对多模式还可支持正则表达式，发布者可以向一个或多个频道发送消息，订阅者可订阅一个或者多个频道接受消息。



发布者：

![img](file:///Users/zhangyy/8lab/github/docs/go/go-redis/images/1075473-20180717142109641-448841196.png?lastModify=1715647828)



订阅者：

![img](file:///Users/zhangyy/8lab/github/docs/go/go-redis/images/1075473-20180717142152579-840572610.png?lastModify=1715647828)



操作示例，示例中将使用两个goroutine分别担任发布者和订阅者角色进行演示：

```
 
 package main
 
 import (
     "github.com/garyburd/redigo/redis"
     "fmt"
     "time"
 )
 
 func Subs() {  //订阅者
     conn, err := redis.Dial("tcp", "10.1.1.116:6379")
     if err != nil {
         fmt.Println("connect redis error :", err)
         return
     }
     defer conn.Close()
     psc := redis.PubSubConn{conn}
     psc.Subscribe("channel1") //订阅channel1频道
     for {
         switch v := psc.Receive().(type) {
         case redis.Message:
             fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
         case redis.Subscription:
             fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
         case error:
             fmt.Println(v)
             return
         }
     }
 }
 
 func Push(message string)  { //发布者
     conn, _ := redis.Dial("tcp", "10.1.1.116:6379")
     _,err1 := conn.Do("PUBLISH", "channel1", message)
        if err1 != nil {
              fmt.Println("pub err: ", err1)
                  return
             }
 
 }
 
 func main()  {
     go Subs()
     go Push("this is wd")
     time.Sleep(time.Second*3)
 }
 //channel1: subscribe 1
 //channel1: message: this is wd
 
```

#### Pipelining(管道)

 执行过程： 发送命令 －> 命令排队 －> 命令执行 －> 返回结果

Pipeline 可以减少对应的 时间

```go
func Example() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}

	p := rdb.Pipeline()
	p.Set(ctx, "name", "zhangSan2", 0)
	p.Exec(ctx)
	
	//fmt.Println(vl)
}

```



#### 事务操作

Redis 事务的本质是一组命令的集合。事务支持一次执行多个命令，一个事务中所有命令都会被序列化。在事务执行过程，会按照顺序串行化执行队列中的命令，其他客户端提交的命令请求不会插入到事务执行命令序列中。

总结说：redis事务就是一次性、顺序性、排他性的执行一个队列中的一系列命令。



- TxPipeline ：开启事务，redis会将后续的命令逐个放入队列中，然后使用EXEC命令来原子化执行这个命令系列。
- EXEC：执行事务中的所有操作命令。
- DISCARD：取消事务，放弃执行事务块中的所有命令。
- WATCH：监视一个或多个key,如果事务在执行前，这个key(或多个key)被其他命令修改，则事务被中断，不会执行事务中的任何命令。
- UNWATCH：取消WATCH对所有key的监视。

示例：

```GO

func Example() {
	ctx := context.Background()
	_ = rdb.Watch(ctx, func(tx *redis.Tx) error {

		// get current value c1 = v1
		val, _ := tx.Get(ctx, "c1").Result()
		fmt.Println(val)

    // 启动管道
		p := tx.Pipeline()
		time.Sleep(1 * time.Second)
		_, _ = p.Set(ctx, "c1", "v2", 0).Result()
		res, err := p.Exec(ctx)
		if err != nil {
			return err
		}
		fmt.Println("res", res)

		return nil
	})

}

func Example2() {
	ctx := context.Background()
	ret, _ := rdb.Set(ctx, "c1", "v3", 0).Result()
	fmt.Println("ret", ret)
}

func main() {
	_ = initClient()

	// 设置一个默认值
	rdb.Set(context.Background(), "c1", "v1", 0)

	// 开启协程
	go Example()
	
	// 正常代码更新
	Example2()

}
```

- ## Redis事务执行步骤

  通过上文命令执行，很显然Redis事务执行是三个阶段：

  - **开启**：以MULTI开始一个事务
  - **入队**：将多个命令入队到事务中，接到这些命令并不会立即执行，而是放到等待执行的事务队列里面
  - **执行**：由EXEC命令触发事务

  当一个客户端切换到事务状态之后， 服务器会根据这个客户端发来的不同命令执行不同的操作：

  - 如果客户端发送的命令为 EXEC 、 DISCARD 、 WATCH 、 MULTI 四个命令的其中一个， 那么服务器立即执行这个命令。
  - 与此相反， 如果客户端发送的命令是 EXEC 、 DISCARD 、 WATCH 、 MULTI 四个命令以外的其他命令， 那么服务器并不立即执行这个命令， 而是将这个命令放入一个事务队列里面， 然后向客户端返回 QUEUED 回复。

  



#### 事务与pipeline区别

- 原子性
  - 事务：事务提供原子性保证，即在事务执行期间的所有命令要么全部执行成功，要么全部回滚。事务中的命令在提交前不会立即执行，而是在 `EXEC` 命令被调用时才执行。如果在事务执行期间发生错误，可以回滚事务，使得之前执行的命令不会对数据库产生影响。
  - 管道：管道不提供原子性保证。它将多个命令一次性发送给 Redis 服务器，但服务器会立即执行这些命令并返回结果。如果在管道执行期间发生错误，仍然会继续执行后续的命令
- 通信方式：
  - 事务：事务使用 `MULTI`、`EXEC` 和 `DISCARD` 命令进行通信。`MULTI` 开始事务，`EXEC` 执行事务，`DISCARD` 取消事务。事务中的命令在 `EXEC` 被调用时才发送给 Redis 服务器执行。
  - 管道：管道使用一次性发送所有命令的方式进行通信。通过将多个命令打包成一个请求发送给 Redis 服务器，减少了网络通信的开销。
- 回复处理：
  - 事务：事务的回复是一个数组，包含每个命令的执行结果。可以通过遍历数组来获取每个命令的回复。
  - 管道：管道的回复是一个数组，按照命令发送的顺序返回每个命令的回复。可以通过索引访问数组中的回复。
- 错误处理：
  - 事务：事务中的命令在执行过程中出现错误并不会立即中断，而是继续执行剩余的命令。在 `EXEC` 被调用时，如果有命令执行失败，整个事务的回复将包含错误信息。没有错误的执行成功正常提交
  - 管道：管道中的命令在执行过程中出现错误会立即中断，并将错误返回给客户端

