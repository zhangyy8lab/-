# redis

简介
--

redis(REmote DIctionary Server)是一个由Salvatore Sanfilippo写key-value存储系统，它由C语言编写、遵守BSD协议、支持网络、可基于内存亦可持久化的日志型、Key-Value类型的数据库，并提供多种语言的API。和Memcached类似，它支持存储的value类型相对更多，包括string(字符串)、list(链表)、set(集合)、zset(sorted set --有序集合)和hash（哈希类型）。这些数据类型都支持push/pop、add/remove及取交集并集和差集及更丰富的操作，而且这些操作都是原子性的。在此基础上，redis支持各种不同方式的排序。与memcached一样，为了保证效率，数据都是缓存在内存中。区别的是redis会周期性的把更新的数据写入磁盘或者把修改操作写入追加的记录文件，并且在此基础上实现了master-slave(主从)同步，redis在3.0版本推出集群模式。

官方网站：https://redis.io/

### 源码部署

```bash
yum install gcc -y  #安装C依赖
wget http://download.redis.io/redis-stable.tar.gz  #下载稳定版本
tar zxvf redis-stable.tar.gz  #解压
cd redis-stable
make PREFIX=/opt/app/redis install   #指定目录编译
make install
mkdir /etc/redis   #建立配置目录
cp redis.conf /etc/redis/6379.conf # 拷贝配置文件
cp utils/redis_init_script /etc/init.d/redis  #拷贝init启动脚本针对6.X系统
chmod a+x  /etc/init.d/redis  #添加执行权限
```



修改配置文件：

```BASH
 vi /etc/redis/6379.conf

bind 0.0.0.0      #监听地址
maxmemory 4294967296   #限制最大内存（4G）：
daemonize yes   #后台运行

```



####启动与停止

```bash
/etc/init.d/redis start
/etc/init.d/redis stop
```



查看版本信息

#执行客户端工具
redis-cli 
#输入命令info
127.0.0.1:6379> info

# Server
redis\_version:4.0.10
redis\_git\_sha1:00000000
redis\_git\_dirty:0
redis\_build\_id:cf83e9c690dbed33
redis\_mode:standalone
os:Linux 2.6.32-642.el6.x86\_64 x86\_64
arch\_bits:64
multiplexing\_api:epoll



二、golang操作redis
---------------

### 安装

golang操作redis的客户端包有多个比如redigo、go-redis，github上Star最多的莫属redigo。

github地址：https://github.com/garyburd/redigo  目前已经迁移到：https://github.com/gomodule/redigo 

文档：https://godoc.org/github.com/garyburd/redigo/redis

go get github.com/garyburd/redigo/redis
import "github.com/garyburd/redigo/redis"

### 连接

Conn接口是与Redis协作的主要接口，可以使用Dial,DialWithTimeout或者NewConn函数来创建连接，当任务完成时，应用程序必须调用Close函数来完成操作。

```go
package main

import (
"github.com/garyburd/redigo/redis"
"fmt"
)

func main()  {
    conn,err := redis.Dial("tcp","10.1.210.69:6379")
    if err != nil {
        fmt.Println("connect redis error :",err)
        return
    }
    defer conn.Close()
}

```



### 命令操作

通过使用Conn接口中的do方法执行redis命令，redis命令大全参考：http://doc.redisfans.com/

go中发送与响应对应类型：

Do函数会必要时将参数转化为二进制字符串

| Go Type | Conversion |
| --- | --- |
| \[\]byte | Sent as is |
| string | Sent as is |
| int, int64 | strconv.FormatInt(v) |
| float64 | strconv.FormatFloat(v, 'g', -1, 64) |
| bool | true -> "1", false -> "0" |
| nil | "" |
| all other types | fmt.Print(v) |

Redis 命令响应会用以下Go类型表示：

| Redis type | Go type |
| --- | --- |
| error | redis.Error |
| integer | int64 |
| simple string | string |
| bulk string | \[\]byte or nil if value not present. |
| array | \[\]interface{} or nil if value not present. |

可以使用GO的类型断言或者reply辅助函数将返回的interface{}转换为对应类型。

操作示例：

#### get、set

```go
import (
  "github.com/garyburd/redigo/redis"
  "fmt"
)

func main()  {
    conn,err := redis.Dial("tcp","10.1.210.69:6379")
    if err != nil {
        fmt.Println("connect redis error :",err)
        return
    }
    defer conn.Close()
    _, err = conn.Do("SET", "name", "wd")
    if err != nil {
        fmt.Println("redis set error:", err)
    }
    name, err := redis.String(conn.Do("GET", "name"))
    if err != nil {
        fmt.Println("redis get error:", err)
    } else {
        fmt.Printf("Got name: %s \n", name)
    }
}

设置key过期时间

  _, err = conn.Do("expire", "name", 10) //10秒过期
    if err != nil {
        fmt.Println("set expire error: ", err)
        return
    }


```



#### 批量获取mget、批量设置mset

```go
_, err = conn.Do("MSET", "name", "wd","age",22)
    if err != nil {
        fmt.Println("redis mset error:", err)
    }
    res, err := redis.Strings(conn.Do("MGET", "name","age"))
    if err != nil {
        fmt.Println("redis get error:", err)
    } else {
        res_type := reflect.TypeOf(res)
        fmt.Printf("res type : %s \n", res_type)
        fmt.Printf("MGET name: %s \n", res)
        fmt.Println(len(res))
    }
//结果：
//res type : []string 
//MGET name: [wd 22] 
//2
```



#### 列表操作

```
gopackage main

import (
"github.com/garyburd/redigo/redis"
"fmt"
    "reflect"
)

func main()  {
    conn,err := redis.Dial("tcp","10.1.210.69:6379")
    if err != nil {
        fmt.Println("connect redis error :",err)
        return
    }
    defer conn.Close()
    _, err = conn.Do("LPUSH", "list1", "ele1","ele2","ele3")
    if err != nil {
        fmt.Println("redis mset error:", err)
    }
    res, err := redis.String(conn.Do("LPOP", "list1"))
    if err != nil {
        fmt.Println("redis POP error:", err)
    } else {
        res_type := reflect.TypeOf(res)
        fmt.Printf("res type : %s \n", res_type)
        fmt.Printf("res  : %s \n", res)
    }
}


```





#### hash操作

```go


package main

import (
"github.com/garyburd/redigo/redis"
"fmt"
    "reflect"
)

func main()  {
    conn,err := redis.Dial("tcp","10.1.210.69:6379")
    if err != nil {
        fmt.Println("connect redis error :",err)
        return
    }
    defer conn.Close()
    _, err = conn.Do("HSET", "student","name", "wd","age",22)
    if err != nil {
        fmt.Println("redis mset error:", err)
    }
    res, err := redis.Int64(conn.Do("HGET", "student","age"))
    if err != nil {
        fmt.Println("redis HGET error:", err)
    } else {
        res_type := reflect.TypeOf(res)
        fmt.Printf("res type : %s \n", res_type)
        fmt.Printf("res  : %d \n", res)
    }
}


```



### Pipelining(管道)

管道操作可以理解为并发操作，并通过Send()，Flush()，Receive()三个方法实现。客户端可以使用send()方法一次性向服务器发送一个或多个命令，命令发送完毕时，使用flush()方法将缓冲区的命令输入一次性发送到服务器，客户端再使用Receive()方法依次按照先进先出的顺序读取所有命令操作结果。

Send(commandName string, args ...interface{}) error
Flush() error
Receive() (reply interface{}, err error)

*   Send：发送命令至缓冲区
*   Flush：清空缓冲区，将命令一次性发送至服务器
*   Recevie：依次读取服务器响应结果，当读取的命令未响应时，该操作会阻塞。

示例：

```go
package main

import (
  "github.com/garyburd/redigo/redis"
  "fmt"
)

func main()  {
    conn,err := redis.Dial("tcp","10.1.210.69:6379")
    if err != nil {
        fmt.Println("connect redis error :",err)
        return
    }
    defer conn.Close()
    conn.Send("HSET", "student","name", "wd","age","22")
    conn.Send("HSET", "student","Score","100")
    conn.Send("HGET", "student","age")
    conn.Flush()

    res1, err :\= conn.Receive()
    fmt.Printf("Receive res1:%v \\n", res1)
    res2, err :\= conn.Receive()
    fmt.Printf("Receive res2:%v\\n",res2)
    res3, err :\= conn.Receive()
    fmt.Printf("Receive res3:%s\\n",res3)

}
//Receive res1:0 
//Receive res2:0
//Receive res3:22

```





### 发布/订阅

redis本身具有发布订阅的功能，其发布订阅功能通过命令SUBSCRIBE(订阅)／PUBLISH(发布)实现，并且发布订阅模式可以是多对多模式还可支持正则表达式，发布者可以向一个或多个频道发送消息，订阅者可订阅一个或者多个频道接受消息。

示意图：

发布者：

![](https://images2018.cnblogs.com/blog/1075473/201807/1075473-20180717142109641-448841196.png)

订阅者：

![](https://images2018.cnblogs.com/blog/1075473/201807/1075473-20180717142152579-840572610.png)

操作示例，示例中将使用两个goroutine分别担任发布者和订阅者角色进行演示：

```go


package main

import (
    "github.com/garyburd/redigo/redis"
    "fmt"
    "time"
)

func Subs() {  //订阅者
    conn, err := redis.Dial("tcp", "10.1.210.69:6379")
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
    conn, _ := redis.Dial("tcp", "10.1.210.69:6379")
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





### 事务操作

MULTI, EXEC,DISCARD和WATCH是构成Redis事务的基础，当然我们使用go语言对redis进行事务操作的时候本质也是使用这些命令。

MULTI：开启事务

EXEC：执行事务

DISCARD：取消事务

WATCH：监视事务中的键变化，一旦有改变则取消事务。

示例：



package main

import (
"github.com/garyburd/redigo/redis"
"fmt"
)


func main()  {
    conn,err :\= redis.Dial("tcp","10.1.210.69:6379")
    if err != nil {
        fmt.Println("connect redis error :",err)
        return
    }
    defer conn.Close()
    conn.Send("MULTI")
    conn.Send("INCR", "foo")
    conn.Send("INCR", "bar")
    r, err :\= conn.Do("EXEC")
    fmt.Println(r)
}
//\[1, 1\]



### 连接池使用

redis连接池是通过pool结构体实现，以下是源码定义，相关参数说明已经备注：





```go
type Pool struct {
    // Dial is an application supplied function for creating and configuring a
    // connection.
    //
    // The connection returned from Dial must not be in a special state
    // (subscribed to pubsub channel, transaction started, ...).
    Dial func() (Conn, error) //连接方法
    
    // TestOnBorrow is an optional application supplied function for checking
    // the health of an idle connection before the connection is used again by
    // the application. Argument t is the time that the connection was returned
    // to the pool. If the function returns an error, then the connection is
    // closed.
    TestOnBorrow func(c Conn, t time.Time) error

    // Maximum number of idle connections in the pool.
    MaxIdle int  //最大的空闲连接数，即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态

    // Maximum number of connections allocated by the pool at a given time.
    // When zero, there is no limit on the number of connections in the pool.
    MaxActive int //最大的激活连接数，同时最多有N个连接

    // Close connections after remaining idle for this duration. If the value
    // is zero, then idle connections are not closed. Applications should set
    // the timeout to a value less than the server's timeout.
    IdleTimeout time.Duration  //空闲连接等待时间，超过此时间后，空闲连接将被关闭

    // If Wait is true and the pool is at the MaxActive limit, then Get() waits
    // for a connection to be returned to the pool before returning.
    Wait bool  //当配置项为true并且MaxActive参数有限制时候，使用Get方法等待一个连接返回给连接池

    // Close connections older than this duration. If the value is zero, then
    // the pool does not close connections based on age.
    MaxConnLifetime time.Duration
    // contains filtered or unexported fields
  }
```




 示例：

```go


package main

import (
    "github.com/garyburd/redigo/redis"
    "fmt"
)

var Pool redis.Pool
func init()  {      //init 用于初始化一些参数，先于main执行
    Pool = redis.Pool{
        MaxIdle:     16,
        MaxActive:   32,
        IdleTimeout: 120,
        Dial: func() (redis.Conn, error) {
            return redis.Dial("tcp", "10.1.210.69:6379")
        },
    }
}

func main()  {

    conn :\=Pool.Get()
    res,err :\= conn.Do("HSET","student","name","jack")
    fmt.Println(res,err)
    res1,err :\= redis.String(conn.Do("HGET","student","name"))
    fmt.Printf("res:%s,error:%v",res1,err)

}
//0 <nil>
//res:jack,error:<nil>

```





三、golang操作mysql
---------------

mysql目前来说是使用最为流行的关系型数据库，golang操作mysql使用最多的包go-sql-driver/mysql。

sqlx包是作为database/sql包的一个额外扩展包，在原有的database/sql加了很多扩展，如直接将查询的数据转为结构体，大大简化了代码书写，当然database/sql包中的方法同样起作用。

github地址：

*   https://github.com/go-sql-driver/mysql
*   https://github.com/jmoiron/sqlx

golang sql使用：

*   [database/sql documentation](http://golang.org/pkg/database/sql/) 
*   [go-database-sql tutorial](http://go-database-sql.org/)

### 安装

go get "github.com/go-sql-driver/mysql"
go get "github.com/jmoiron/sqlx"

### 连接数据库

var Db \*sqlx.DB
db, err :\= sqlx.Open("mysql","username:password@tcp(ip:port)/database?charset=utf8")
Db \= db

### 处理类型（Handle Types)

sqlx设计和database/sql使用方法是一样的。包含有4中主要的handle types： 

*   sqlx.DB - 和sql.DB相似，表示数据库。 
*   sqlx.Tx - 和sql.Tx相似，表示事物。 
*   sqlx.Stmt - 和sql.Stmt相似，表示prepared statement。 
*   sqlx.NamedStmt - 表示prepared statement（支持named parameters）

所有的handler types都提供了对database/sql的兼容，意味着当你调用sqlx.DB.Query时，可以直接替换为sql.DB.Query.这就使得sqlx可以很容易的加入到已有的数据库项目中。

此外，sqlx还有两个cursor类型： 

*   sqlx.Rows - 和sql.Rows类似，Queryx返回。 
*   sqlx.Row - 和sql.Row类似，QueryRowx返回。

相比database/sql方法还多了新语法，也就是实现将获取的数据直接转换结构体实现。

*   Get(dest interface{}, …) error
*   Select(dest interface{}, …) error 

### 建表

以下所有示例均已以下表结构作为操作基础。

```go
CREATE TABLE `userinfo` (
    `uid` INT(10) NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64)  DEFAULT NULL,
    `password` VARCHAR(32)  DEFAULT NULL,
    `department` VARCHAR(64)  DEFAULT NULL,
    `email` varchar(64) DEFAULT NULL,
    PRIMARY KEY (`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8
```



### Exec使用

Exec和MustExec从连接池中获取一个连接然后只想对应的query操作。对于不支持ad-hoc query execution的驱动，在操作执行的背后会创建一个prepared statement。在结果返回前这个connection会返回到连接池中。

需要注意的是不同的数据库类型使用的占位符不同，mysql采用？作为占位符号。

*   MySQL 使用？ 
*   PostgreSQL 使用1,1,2等等 
*   SQLite 使用？或$1 
*   Oracle 使用:name

### Exec增删该示例

查询语法使用Query后续会提到

```go


package main

import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "fmt"
)

var Db *sqlx.DB

func init()  {
    db, err := sqlx.Open("mysql", "stu:1234qwer@tcp(10.0.0.241:3307)/test?charset=utf8")
    if err != nil {
        fmt.Println("open mysql failed,", err)
        return
    }
    Db = db
}

func main()  {
    result, err := Db.Exec("INSERT INTO userinfo (username, password, department,email) VALUES (?, ?, ?,?)","wd","123","it","wd@163.com")
    if err != nil{
        fmt.Println("insert failed,error： ", err)
        return
    }
    id,_ := result.LastInsertId()
    fmt.Println("insert id is :",id)
    _, err1 := Db.Exec("update userinfo set username = ? where uid = ?","jack",1)
    if err1 != nil{
        fmt.Println("update failed error:",err1)
    } else {
        fmt.Println("update success!")
    }
    _, err2 := Db.Exec("delete from userinfo where uid = ? ", 1)
    if err2 != nil{
        fmt.Println("delete error:",err2)
    }else{
        fmt.Println("delete success")
    }

}
//insert id is : 1
//update success!
//delete success

```





### sql预声明（Prepared Statements）

对于大部分的数据库来说，当一个query执行的时候，在sql语句数据库内部声明已经声明过了，其声明是在数据库中，我们可以提前进行声明，以便在其他地方重用。



stmt, err := db.Prepare(\`SELECT \* FROM place WHERE telcode=?\`)
row \= stmt.QueryRow(65)

tx, err :\= db.Begin()
txStmt, err :\= tx.Prepare(\`SELECT \* FROM place WHERE telcode=?\`)
row \= txStmt.QueryRow(852)



当然sqlx还提供了Preparex()进行扩展，可直接用于结构体转换

stmt, err := db.Preparex(\`SELECT \* FROM place WHERE telcode=?\`)
var p Place
err \= stmt.Get(&p, 852)

### Query

Query是database/sql中执行查询主要使用的方法，该方法返回row结果。Query返回一个sql.Rows对象和一个error对象。

在使用的时候应该吧Rows当成一个游标而不是一系列的结果。尽管数据库驱动缓存的方法不一样，通过Next()迭代每次获取一列结果，对于查询结果非常巨大的情况下，可以有效的限制内存的使用，Scan()利用reflect把sql每一列结果映射到go语言的数据类型如string，\[\]byte等。如果你没有遍历完全部的rows结果，一定要记得在把connection返回到连接池之前调用rows.Close()。

Query返回的error有可能是在server准备查询的时候发生的，也有可能是在执行查询语句的时候发生的。例如可能从连接池中获取一个坏的连级（尽管数据库会尝试10次去发现或创建一个工作连接）。一般来说，错误主要由错误的sql语句，错误的类似匹配，错误的域名或表名等。

在大部分情况下，Rows.Scan()会把从驱动获取的数据进行拷贝，无论驱动如何使用缓存。特殊类型sql.RawBytes可以用来从驱动返回的数据总获取一个zero-copy的slice byte。当下一次调用Next的时候，这个值就不在有效了，因为它指向的内存已经被驱动重写了别的数据。

Query使用的connection在所有的rows通过Next()遍历完后或者调用rows.Close()后释放。 

示例：



package main

import (
    \_ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "fmt"
)

var Db \*sqlx.DB

func init()  {
    db, err :\= sqlx.Open("mysql", "stu:1234qwer@tcp(10.0.0.241:3307)/test?charset=utf8")
    if err != nil {
        fmt.Println("open mysql failed,", err)
        return
    }
    Db \= db
}

func main()  {
    rows, err :\= Db.Query("SELECT username,password,email FROM userinfo")
    if err != nil{
        fmt.Println("query failed,error： ", err)
        return
    }
    for rows.Next() {  //循环结果
        var username,password,email string
        err \= rows.Scan(&username, &password, &email)
        println(username,password,email)
    }
    
}
//wd 123 wd@163.com
//jack 1222 jack@165.com



### Queryx

Queryx和Query行为很相似，不过返回一个sqlx.Rows对象，支持扩展的scan行为,同时可将对数据进行结构体转换。

示例：



package main

import (
    \_ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "fmt"
)

var Db \*sqlx.DB

type stu struct {
    Username string   \`db:"username"\`
    Password string      \`db:"password"\`
    Department string  \`db:"department"\`
    Email string        \`db:"email"\`
}

func init()  {
    db, err :\= sqlx.Open("mysql", "stu:1234qwer@tcp(10.0.0.241:3307)/test?charset=utf8")
    if err != nil {
        fmt.Println("open mysql failed,", err)
        return
    }
    Db \= db
}

func main()  {
    rows, err :\= Db.Queryx("SELECT username,password,email FROM userinfo")
    if err != nil{
        fmt.Println("Qeryx failed,error： ", err)
        return
    }
    for rows.Next() {  //循环结果
        var stu1 stu
        err \= rows.StructScan(&stu1)// 转换为结构体
        fmt.Println("stuct data：",stu1.Username,stu1.Password)
    }
}
//stuct data： wd 123
//stuct data： jack 1222



### QueryRow和QueryRowx

QueryRow和QueryRowx都是从数据库中获取一条数据，但是QueryRowx提供scan扩展，可直接将结果转换为结构体。



package main

import (
    \_ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "fmt"
)

var Db \*sqlx.DB

type stu struct {
    Username string   \`db:"username"\`
    Password string      \`db:"password"\`
    Department string  \`db:"department"\`
    Email string        \`db:"email"\`
}

func init()  {
    db, err :\= sqlx.Open("mysql", "stu:1234qwer@tcp(10.0.0.241:3307)/test?charset=utf8")
    if err != nil {
        fmt.Println("open mysql failed,", err)
        return
    }
    Db \= db
}

func main()  {
    row :\= Db.QueryRow("SELECT username,password,email FROM userinfo where uid = ?",1) // QueryRow返回错误，错误通过Scan返回
    var username,password,email string
    err :\=row.Scan(&username,&password,&email)
    if err != nil{
        fmt.Println(err)
    }
    fmt.Printf("this is QueryRow res:\[%s:%s:%s\]\\n",username,password,email)
    var s stu
    err1 :\= Db.QueryRowx("SELECT username,password,email FROM userinfo where uid = ?",2).StructScan(&s)
    if err1 != nil{
        fmt.Println("QueryRowx error :",err1)
    }else {
        fmt.Printf("this is QueryRowx res:%v",s)
    }
}
//this is QueryRow res:\[wd:123:wd@163.com\]
//this is QueryRowx res:{jack 1222  jack@165.com}



### Get 和Select（非常常用）

Get和Select是一个非常省时的扩展，可直接将结果赋值给结构体，其内部封装了StructScan进行转化。Get用于获取单个结果然后Scan，Select用来获取结果切片。

示例：



package main

import (
    \_ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "fmt"
)

var Db \*sqlx.DB

type stu struct {
    Username string   \`db:"username"\`
    Password string      \`db:"password"\`
    Department string  \`db:"department"\`
    Email string        \`db:"email"\`
}

func init()  {
    db, err :\= sqlx.Open("mysql", "stu:1234qwer@tcp(10.0.0.241:3307)/test?charset=utf8")
    if err != nil {
        fmt.Println("open mysql failed,", err)
        return
    }
    Db \= db
}

func main()  {
    var stus \[\]stu
    err :\= Db.Select(&stus,"SELECT username,password,email FROM userinfo")
    if err != nil{
        fmt.Println("Select error",err)
    }
    fmt.Printf("this is Select res:%v\\n",stus)
    var s stu
    err1 :\= Db.Get(&s,"SELECT username,password,email FROM userinfo where uid = ?",2)
    if err1 != nil{
        fmt.Println("GET error :",err1)
    }else {
        fmt.Printf("this is GET res:%v",s)
    }
}
//this is Select res:\[{wd 123  wd@163.com} {jack 1222  jack@165.com}\]
//this is GET res:{jack 1222  jack@165.com}



### 事务（Transactions）

事务操作是通过三个方法实现：

Begin()：开启事务

Commit()：提交事务（执行sql)

Rollback()：回滚

使用流程：



tx, err := db.Begin()
err \= tx.Exec(...)
err \= tx.Commit()

//或者使用sqlx扩展的事务
tx := db.MustBegin()
tx.MustExec(...)
err \= tx.Commit()



由于事务是一个一直连接的状态，所以Tx对象必须绑定和控制单个连接。一个Tx会在整个生命周期中保存一个连接，然后在调用commit或Rollback()的时候释放掉。在调用这几个函数的时候必须十分小心，否则连接会一直被占用直到被垃圾回收。 

使用示例：



package main

import (
    \_ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "fmt"
)

var Db \*sqlx.DB

func init()  {
    db, err :\= sqlx.Open("mysql", "stu:1234qwer@tcp(10.0.0.241:3307)/test?charset=utf8")
    if err != nil {
        fmt.Println("open mysql failed,", err)
        return
    }
    Db \= db
}

func main()  {
    tx, err :\= Db.Beginx()
    \_, err \= tx.Exec("insert into userinfo(username,password) values(?,?)", "Rose","2223")
    if err != nil {
        tx.Rollback()
    }
    \_, err \= tx.Exec("insert into userinfo(username,password) values(?,?)", "Mick",222)
    if err != nil {
        fmt.Println("exec sql error:",err)
        tx.Rollback()
    }
    err \= tx.Commit()
    if err != nil {
        fmt.Println("commit error")
    }

}



### 连接池设置

默认情况下，连接池增长无限制，并且只要连接池中没有可用的空闲连接，就会创建连接。我们可以使用DB.SetMaxOpenConns设置池的最大大小。未使用的连接标记为空闲，如果不需要则关闭。要避免建立和关闭大量连接，可以使用DB.SetMaxIdleConns设置最大空闲连接。

注意：该设置方法golang版本至少为1.2

*   DB.SetMaxIdleConns(n int)    设置最大空闲连接数
*   DB.SetMaxOpenConns(n int)  设置最大打开的连接数

示例：



package main

import (
    \_ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "fmt"
)

var Db \*sqlx.DB

func init()  {
    db, err :\= sqlx.Open("mysql", "stu:1234qwer@tcp(10.0.0.241:3307)/test?charset=utf8")
    if err != nil {
        fmt.Println("open mysql failed,", err)
        return
    }
    Db \= db
    Db.SetMaxOpenConns(30)
    Db.SetMaxIdleConns(15)

}



参考：http://jmoiron.github.io/sqlx/

本文转自 <https://www.cnblogs.com/wdliu/p/9330278.html>，如有侵权，请联系删除。
