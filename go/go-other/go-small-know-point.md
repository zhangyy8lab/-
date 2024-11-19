# go smallPoint

## go-cmd run args



```go

// 定义用户信息结构体
type UserInfo struct {
  Name string
  Age int
}

// 客户端结构体
type Cli struct {
  U *UserInfo
}

// 提示信息
const Usage = `
	addUser --data "zhangsan, 21"
	printUser --data "username"
	printAll "print all user" 
`

// 客户端方法
func (cli *Cli) Run {
  cmds := os.Args
  if len(cmds) < 2 {
    fmt.Println("params error")
    fmt.Println(Usage)
  }
  // ....
  // 这里需要做一些其他的事情可根据自身需要进行处理， 添加数据库， init配置信息。。。。
  switch cmds[1] {
    // 输入的参数为 addUser 
    case "addUser": 
    	// addUser func...
    case "printUser":
    	// find for userName and print
    case "printAll":
    	// print all user
    default:
      fmt.Println(Usage)
  }
  
}
```



## go-print

```go
package main

import (
	"fmt"
)

type User struct {
	Name   string
	Age    int
	Active bool
}

func main() {
	// 布尔型
	t := true
	fmt.Printf("布尔值：%t\n", t) // 布尔值：true

	// 整数类型
	i := 42
	fmt.Printf("整数：%d\n", i)
	fmt.Printf("二进制：%b\n", i)
	fmt.Printf("八进制：%o\n", i)
	fmt.Printf("十六进制：%x\n", i)
	fmt.Printf("字符：%c\n", i)
	fmt.Printf("Unicode 码点：%U\n", i)

    // 整数：42
    // 二进制：101010
    // 八进制：52
    // 十六进制：2a
    // 字符：*
    // Unicode 码点：U+002A
    // ------------------------------------------------

	// 浮点数
	f := 3.14159
	fmt.Printf("浮点数：%f\n", f)
	fmt.Printf("科学计数法：%e\n", f)
	fmt.Printf("紧凑格式：%g\n", f)
    // 浮点数：3.141590
    // 科学计数法：3.141590e+00
    // 紧凑格式：3.14159
    // ------------------------------------------------

	// 字符串
	str := "Hello, Go!"
	fmt.Printf("字符串：%s\n", str)
	fmt.Printf("带引号的字符串：%q\n", str)
	fmt.Printf("字符串的十六进制：%x\n", str)
    // 字符串：Hello, Go!
    // 带引号的字符串："Hello, Go!"
    // 字符串的十六进制：48656c6c6f2c20476f21
    // ------------------------------------------------

	// 结构体
	u := User{Name: "Alice", Age: 30, Active: true}
	fmt.Printf("结构体：%v\n", u)
	fmt.Printf("带字段名的结构体：%+v\n", u)
	fmt.Printf("Go 语法格式表示的结构体：%#v\n", u)
    // 结构体：{Alice 30 true}
    // 带字段名的结构体：{Name:Alice Age:30 Active:true}
    // Go 语法格式表示的结构体：main.User{Name:"Alice", Age:30, Active:true}
    // ------------------------------------------------

	// 指针
	ptr := &u
	fmt.Printf("指针：%p\n", ptr)
    // 指针：0xc000010230
    // ------------------------------------------------

	// 变量类型
	fmt.Printf("变量类型：%T\n", u)
    // 变量类型：main.User
    // ------------------------------------------------
}
```



## go-chanel select 

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1)
		ch1 <- "one"
	}()

	go func() {
		time.Sleep(2)
		ch2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case ms1 := <-ch1: // chanel ch1 接收到消息 执行
			fmt.Println("ch1:", ms1)

		case ms2 := <-ch2:
			fmt.Println("ch2:", ms2)
		}
	}
}

```

## go-routine sync.Wait()

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker2(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("worker2 %d starting\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker2 %d done\n", i)
}

func worker1(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("worker1 %d starting\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker1 %d done\n", i)

}

func main() {
	var wg sync.WaitGroup
	for i := 1; i < 3; i++ {
		wg.Add(2) // 记数
		go worker1(i, &wg)
		go worker2(i, &wg)
	}

	wg.Wait() // 等所有协程结束

	fmt.Println("All worker done")
}


```

## go Q/R

### New/Make 区别

```go
	•	new：分配内存，但不初始化，用于值类型（如结构体、数组、整数）。返回零值的指针，不初始化，返回类型零值的指针
	•	make：分配并初始化内存，用于特定的引用类型（slice、map、channel）。初始化后的对象可直接使用

// New 
type Person struct {
    Name string
    Age  int
}

p := new(Person) // p 是 *Person 指针，p.Name == ""，p.Age == 0
fmt.Println(p)   // 输出: &{"" 0}


// Make
s := make([]int, 10) // 创建一个长度为 10 的 int 切片
m := make(map[string]int) // 创建一个空的 map
c := make(chan int) // 创建一个 int 类型的 channel

```

### Printf(), sprintf(), Fprintf() 区别

```go
•	控制台打印：使用 Printf。   	fmt.Printf("Name: %s, Age: %d\n", name, age) 
•	格式化字符串：使用 Sprintf。 result := fmt.Sprintf("Name: %s, Age: %d", name, age)
•	写入文件/网络等：使用 Fprintf。
	file, err := os.Create("output.txt")
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()

  name := "Charlie"
  age := 40
  fmt.Fprintf(file, "Name: %s, Age: %d\n", name, age)
```

### 命令作用

```go

go env:  			// 用于查看go的环境变量
go run:  			// 用于编译并运行go源码文件
go build: 		// 用于编译源码文件、代码包、依赖包go get:#用于动态获取远程代码包
go install: 	// 用于编译go文件，并将编译结构安装到bin、pkg目录
go clean: 		// 用于清理工作目录，删除编译和安装遗留的目标文件
go version:		// 用于查看go的版本信息

```

### go语言的协程

```go
在函数或方法调用前添加 go 关键字会启动一个新的协程（goroutine）来并发执行该函数或方法
	1.	非阻塞执行：在函数或方法前加上 go，不会阻塞主函数的执行。
	2.	生命周期：当主函数退出时，所有未完成的 goroutine 都会被终止。
	3.	并发问题：多个 goroutine 同时访问共享资源时，需通过通道（channel）或同步机制来避免数据竞争。

// 同步机制
	•	简单的并发增减、读取：考虑 sync/atomic。
	•	读多写少的共享资源：考虑 sync.RWMutex。
	•	一组 goroutine 的完成通知：使用 sync.WaitGroup。
	•	复杂的条件等待和通知：使用 sync.Cond。
	•	只执行一次的初始化操作：使用 sync.Once。
```

### GC 

```go 
Go 语言的垃圾回收机制（Garbage Collection，简称 GC）采用了**并发标记-清除（Concurrent Mark and Sweep）**的方式，自动管理内存，避免内存泄漏或使用未释放的内存。

会在满足一定条件时自动触发，以保证内存得到及时回收，防止内存占用过高。触发 GC 的主要条件如下：

// 堆内存使用量达到触发阈值: 
	• GOGC 环境变量的默认值是 100，代表当堆内存增加到上次回收后使用量的 100%（即两倍）时触发下一次 GC。
			例如，上次 GC 后堆内存使用量为 50 MB，那么当堆内存增长到 100 MB 时会触发下一次 GC。
	•	你可以调整 GOGC 的值来控制 GC 触发的频率。比如：
	•	GOGC=200 表示当堆内存使用量达到上次回收后使用量的两倍时触发 GC，GC 触发频率变低，内存使用量上升。
	•	GOGC=50 表示堆内存使用量达到上次回收后使用量的 1.5 倍时触发，GC 频率变高，但内存占用更少。

// 手动触发 GC 
  • 可以手动调用 runtime.GC() 来触发垃圾回收，但这通常不建议在生产环境中使用。手动调用不会改变自动 GC 的触发条件，而只是立即启动一次完整的 GC。

// 垃圾回收（GC）执行后，并不会直接“清除”内存中的数据，而是会标记并回收不再使用的对象，释放这些对象占用的内存空间。具体来说，GC 过程的作用是：

```

> • 垃圾回收（GC）执行后，并不会直接“清除”内存中的数据，而是会标记并回收不再使用的对象，释放这些对象占用的内存空间
>
> 1. 标记不再使用的对象：
>
>    • 在 GC 的标记阶段，Go 会遍历所有活动的对象（通过根对象，如全局变量、栈上的局部变量等）并标记它们。如果某个对象不能通过活动对象（引用链）访问，则认为它不再被使用。
>
> 2.回收不再使用的对象：
>
> ​	• 在标记阶段结束后，Go 会将未被标记的对象回收，释放其占用的内存。这些对象被认为是“垃圾”，不再对程序的运行产生影响。
>
> ​	• 这些回收的对象的内存将会被归还给操作系统，供后续的内存分配使用。
>
> 3. 已回收的对象：GC 会清除那些不再使用的对象，这些对象的内存会被标记为空闲状态，可以供新的内存分配使用。
>
> 4. 仍然被引用的对象：如果某些对象仍然被引用，它们不会被 GC 清除，仍会保持在内存中，直到它们不再被引用为止。

### 组和切片有什么区别 

```go
// 数组
	•	固定长度：数组的大小是编译时确定的，一旦定义就不能改变。
	•	值类型：数组是值类型，当数组被赋值或传递到函数时，会进行值拷贝，意味着副本会被创建。
	•	内存布局：数组元素在内存中是连续存储的，且数组本身也是一个类型，包含元素的类型和长度。

// 切片
	•	动态长度：切片的长度是可以动态改变的（例如，通过 append()）。
	•	引用类型：切片是引用类型，它包含一个指向数组的指针、切片的长度和容量。当切片被赋值或传递时，是传递的引用，而不是复制整个数据。
	•	内存布局：切片底层实现是基于数组的，切片本身不存储数据，而是指向一个数组的一个子数组（从切片的起始位置到其末尾）。切片的容量也可以比切片的长度大。
	•	灵活性：切片提供了比数组更高的灵活性，常常用于操作动态数据。
```

### Http&tcp/ip的区别

```go
// 协议层级
• TCP/IP 是一种网络协议套件，属于传输层和网络层协议的集合, 负责数据包的传输、路由和可靠传输
• HTTP 是应用层协议，位于 TCP/IP 协议栈的上层
// 功能和目的'
• TCP 提供可靠的数据传输服务，包括数据分段、重传、顺序确认等功能。
	IP 负责路由数据包，使数据能跨网络到达目标地址。
• TCP 通过三次握手建立连接、确认数据顺序、错误校正等，确保数据传输的完整性。
• HTTP 规定了客户端和服务器之间请求和响应的数据格式，比如 GET、POST 请求、状态码、响应头等。
•	请求-响应模式：HTTP 使用请求-响应模式进行通信，即客户端发送请求，服务器返回响应。
•	无状态性：HTTP 是无状态协议，每次请求都是独立的，不保留连接状态。

物理层，数据链路层, 网络层, 传输层，会话层, 表示层, 应用层
```

