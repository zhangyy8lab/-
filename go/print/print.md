## fmt.Print的各种输出

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

