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

