package main

import (
	"encoding/gob"
	"fmt"
	"github.com/bits-and-blooms/bloom"
	"os"
)

// 保存布隆过滤器的快照
func saveBloomFilterSnapshot(filter *bloom.BloomFilter, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(filter)
}

// 加载布隆过滤器的快照
func loadBloomFilterSnapshot(filename string) (*bloom.BloomFilter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var filter bloom.BloomFilter
	if err := decoder.Decode(&filter); err != nil {
		return nil, err
	}

	return &filter, nil
}

func main() {

	backupFile := "bloom_filter.gob"
	//
	//// 设置布隆过滤器参数
	//size := uint(1000000) // 位数组大小
	//
	//// 创建布隆过滤器实例
	//filter := bloom.NewWithEstimates(size, 0.01)
	//
	//// 添加一些元素
	//elements := []string{"apple", "banana", "cherry"}
	//for _, element := range elements {
	//	filter.Add([]byte(element))
	//}
	//
	//// 检查元素是否存在
	//testElements := []string{"apple", "grape", "cherry"}
	//for _, testElement := range testElements {
	//	if filter.Test([]byte(testElement)) {
	//		fmt.Printf("%s might be in the filter\n", testElement)
	//	} else {
	//		fmt.Printf("%s is definitely not in the filter\n", testElement)
	//	}
	//}
	//
	//if err := saveBloomFilterSnapshot(filter, backupFile); err != nil {
	//	fmt.Println("backup err:", err.Error())
	//
	//}

	filter, err := loadBloomFilterSnapshot(backupFile)
	if err != nil {
		fmt.Println("load backup err:", err.Error())
	}

	if filter.Test([]byte("apple")) {
		fmt.Println("apple exist")
	} else {
		fmt.Println("apple not exist")
	}

	if filter.Test([]byte("zhangSan")) {
		fmt.Println("zhangSan exist")
	} else {
		fmt.Println("zhangSan not exist")
	}

}
