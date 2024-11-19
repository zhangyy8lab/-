package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 交换输出数字和字符串
// 12ab34cd ...
func PrintStrNumAlternately() {
	var wg sync.WaitGroup
	charChan := make(chan bool)
	numChan := make(chan bool)
	done := make(chan bool)

	wg.Add(2)
	go func() {
		for i := 1; i <= 28; i += 2 {
			<-numChan
			fmt.Printf("%d%d", i, i+1)
			charChan <- true
		}
	}()

	go func() {
		for ch := 'A'; ch <= 'Z'; ch += 2 {
			<-charChan
			fmt.Printf("%c%c", ch, ch+1)
			numChan <- true
		}
		done <- true
	}()

	numChan <- true
	<-done // 等待打印完成
}

// 输出不重复字符
func PrintUnisqueStr(s string) bool {
	if len(s) > 300 {
		return false
	}

	for _, ch := range s {
		if ch > 127 {
			return false
		}

		if strings.Count(s, string(ch)) > 1 {
			return false // 如果某个字符出现超过1次，说明有重复
		}
	}
	return true

}

// 翻转字符串
func swapStr(s string) string {
	runes := []rune(s)
	leftIndex, rightIndex := 0, len(runes)-1
	for leftIndex < rightIndex {
		runes[leftIndex], runes[rightIndex] = runes[rightIndex], runes[leftIndex]
		leftIndex++
		rightIndex--
	}
	return string(runes)
}

// 替换空格
func replaceSpace(s string) {
	s2 := strings.Replace(s, " ", "%20", -1)
	fmt.Printf("%s", s2)
}

// 指针与接口
func Param1() {
	type Param map[string]interface{}
	type Show struct {
		Param
	}

	s := new(Show)
	s.Param = make(Param)
	s.Param["rmb"] = 1000
	fmt.Println(s.Param)
}

// channel
func Channel1() {
	ch := make(chan int, 10)
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			fmt.Println("p i:", i)
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for {
			i, ok := <-ch
			if !ok {
				fmt.Println("close")
				os.Exit(0)
			}
			fmt.Println("c i:", i)
		}
	}()

	wg.Wait()
	fmt.Println("process started ")
	time.Sleep(time.Second * 3)
}

func defer_call() {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()

	// panic("触发异常")
}

func main() {
	// PrintStrNumAlternately()  // 交换输出字符串和数字

	// s := "abcdaaa a"
	// bool := PrintUnisqueStr(s) // 判断一个字符串是否有重复
	// fmt.Printf("%v", bool)

	// newStr := swapStr(s) // 翻转字符串
	// fmt.Println(newStr)

	// replaceSpace(s)

	// channel // 通道的使用
    //	defer_call()
	var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    fmt.Printf("GC 总次数: %d\n", memStats.NumGC)
    fmt.Printf("最后一次 GC 暂停时间 (纳秒): %d\n", memStats.PauseNs[(memStats.NumGC+255)%256])
    fmt.Printf("总分配的内存 (字节): %d\n", memStats.TotalAlloc)
    fmt.Printf("堆上分配的对象数: %d\n", memStats.HeapObjects)
    fmt.Printf("下一次 GC 的内存目标 (字节): %d\n", memStats.NextGC)
}
