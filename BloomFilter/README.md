# bloom Filter

## 布隆过滤器原理

### BloomFilter 的算法

- 分配一块内存空间做 bit 数组（ m ），数组的 bit 位初始值全部设为 0。
- 加入==元素==时，使用 k个 Hash(元素) 计算, 将k个Hash结果映射到数组下标并全部设置为 1。
- 检测 key 是否存在，仍然用这 k 个 Hash 函数计算出 k 个位置，如果位置全部为 1，则表明 key 存在，否则不存在。如下图所示：

![](/Users/zhangyy/8lab/github/docs/images/bloomFilter.png)



### hash算法

> 灵活哈希算法， 常用有 hash256、sha128、md5、hash64(CityHash)等
>
> - CityHash 是一种高效的哈希函数，生成 64 位的哈希值，适用于性能要求较高的场景



### hash结果对应到bit数组

> 哈希函数的输出通常是一个整数,这个整数需要被映射到位数组的有效范围内
>
> - 哈希函数输出： 哈希函数将输入元素映射到一个大范围的整数值
> - 映射到位数组：将哈希函数输出的整数值通过取模操作（模运算）映射到位数组的有效下标范围。假设位数组的大小为 m，则下标的计算方式通常是 `index = hash_value % m`  hash_value 是个整数

```python
import hashlib
sha256_hash = hashlib.sha256(b'example').digest()

# 假设位数组大小为 10
bit_array_size = 10
index = int.from_bytes(sha256_hash[:4], 'big') % bit_array_size
print(index) == 8

# 大致原理 将字节转换为整数 %取余 数组大小 = 在数组中的下标
```





### 误报率（False Positive Rate）

>  布隆过滤器的一个关键特点是它可能会错误地判断一个不在集合中的元素为在集合中（即误报）。使用的哈希函数数量直接影响误报率。
>
> 布隆过滤器的误报率（ P ）与哈希函数数量（ k ）、位数组大小（ m ）、以及已插入的元素数量（ n ）有关。具体的误报率公式如下：

$$
P \approx \left(1 - e^{-\frac{k \cdot n}{m}}\right)^k
$$

> 其中：
> 	•	 n ：插入的元素数量
> 	•	 m ：位数组的大小
> 	•	 k ：哈希函数的数量
> 	•	 e ：自然对数的底数（约等于2.71828）



## go使用 bloomFilter

```go
package main

import (
    "fmt"
    "github.com/bits-and-blooms/bloom"
)

func main() {
    // 创建一个布隆过滤器，位数组大小为 1,000,000 位，使用 5 个哈希函数
    // Bloom filter size = 1,000,000 bits, with 5 hash functions
    n := uint(1000) // 预计插入的元素数量
    m := uint(1000000) // 位数组大小
    k := uint(5) // 哈希函数数量

    // 创建布隆过滤器
    filter := bloom.NewWithEstimates(m, k, float64(n))

    // 向布隆过滤器添加元素
    filter.Add([]byte("apple"))
    filter.Add([]byte("banana"))

    // 查询元素是否存在
    fmt.Println("apple in filter:", filter.Test([]byte("apple")))   // Output: true
    fmt.Println("banana in filter:", filter.Test([]byte("banana"))) // Output: true
    fmt.Println("grape in filter:", filter.Test([]byte("grape")))   // Output: false
}
```



## 布隆过滤器的内存占用

> 布隆过滤器（Bloom Filter）的内存占用主要由其位数组的大小决定。每个位占用 1 位（bit）的内存，但在实际计算中，我们通常按字节（byte）计算内存，因为大多数编程语言和计算机系统以字节为单位分配内存。
>
> 假设布隆过滤器的位数组大小为 10 亿位（10^9 位），我们可以计算其内存占用如下：
>
> - **计算位数到字节数**
>
>   1 字节 = 8 位,因此，内存占用（字节数） = 位数 / 8
>
>   对于 10 亿位：
>   $$
>   \text{内存占用（字节数）} = \frac{10^9 \text{ 位}}{8} = 125,000,000 \text{ 字节}
>   $$
>
> - **转换为更常用的单位**：
>
>   1 MB（兆字节） = 1,024 × 1,024 字节 = 1,048,576 字节
>
>   1 GB（千兆字节） = 1,024 × 1,024 × 1,024 字节 = 1,073,741,824 字节
>
>   将字节数转换为 MB
>   $$
>   \text{内存占用（MB）} = \frac{125,000,000 \text{ 字节}}{1,048,576} \approx 119.2 \text{ MB}
>   $$
>   将字节数转换为 GB：
>   $$
>   \text{内存占用（GB）} = \frac{125,000,000 \text{ 字节}}{1,073,741,824} \approx 0.116 \text{ GB}
>   $$
>
> - **总结**
>
>   ​	布隆过滤器的位数组大小为 10 亿位（10^9 位）大约占用 125,000,000 字节。
>
>   ​	这个内存大小大约是 119.2 MB，或者约 0.116 GB。



## 布隆过滤器的备份策略

- 文件导出/导入
- 分布式（待研究）

