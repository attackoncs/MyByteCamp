这是字节青训营第三课：高质量编程与性能调优实战的笔记和总结

# 概要

![image-20220516155127481](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516155127481.png)

# 准备

- 尝试使用 test 命令，编写并运行简单[测试](https://go.dev/doc/tutorial/add-a-test)
- 尝试使用 -bench参数，对函数进行[性能测试](https://pkg.go.dev/testing#hdr-Benchmarks)
- 推荐阅读[Go代码Review建议](https://github.com/golang/go/wiki/CodeReviewComments)、[Uber的Go编码规范](https://github.com/uber-go/guide)

# 高质量编程

## 简介

高质量编程简言之就是代码达到正确、可靠、简洁清晰的目标：

正确性：考虑各种边界条件，错误的调用正确处理

可靠性：异常或错误处理策略保障依赖的服务出现异常能够处理

简洁：逻辑简单，后续调整或新增功能能快速支持

清晰：代码易于阅读理解，重构或修改功能不易出问题

## 编程原则

实际应用场景千变万化，各语法特性和语法各不相同，但原则相通

+ 简单性：逻辑清晰简单，无多余复杂性，易于理解改进
+ 可读性：代码给人看而非机器，可维护性前提是可读性
+ 生产力：团队整体效率非常重要

## 编码规范

### 代码格式

用 gofmt和goimports格式化代码和包，保证代码与官方推荐格式一致

### 注释

> Good code has lots of comments,bad code requires lots of comments.

+ 注释应该解释代码作用，适合注释公共符合，参考官方[代码](https://github.com/golang/go/blob/master/src/os/file.go#L313)
+ 注释应该解释代码如何做的，适合注释方法，参考官方[代码](https://github.com/golang/go/blob/master/src/net/http/client.go#L678)
+ 注释应该解释代码实现的原因，解释代码外部因素，参考官方[代码](https://github.com/golang/go/blob/master/src/net/http/client.go#L521)

+ 注释应该解释代码什么情况会出错
+ 包中每个公共符合：变量、常量、函数及结构体都要注释，参考官方[代码](https://github.com/golang/go/blob/master/src/io/io.go#L455)


#### 场景一

![image-20220516113722630](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516113722630.png)

如图，Open应解释作用，IsTableFull解释则无必要，因为已见名知意

#### 场景二

![image-20220516113753689](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516113753689.png)

第一个注释因逻辑较为复杂，需要注释，而第二个则完全没必要

#### 场景三

![image-20220516113958990](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516113958990.png)

如图shouldRedirect=false若脱离上下文后很难理解，需注释说明原因

#### 场景四

![image-20220516114240548](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516114240548.png)

注释应提醒潜在限制条件或无法处理情况，让使用者无需了解细节

#### 场景五

![image-20220516114545754](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516114545754.png)

![image-20220516114905948](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516114905948.png)

**包中每个公共符号始如变量、常量、函数及结构体终要注释**，唯一例外是不要注释实现接口的方法

#### 小结

- 代码是最好的注释
- 注释应提供代码未表达出的上下文信息，包括作用、实现、原因、出错情况等

### 命名规范

> Good naming is like good joke,if you have to explain it,it's not funny  --Dave Cheney

核心是降低阅读和理解代码的成本，重点考虑设计简洁清晰名称并考虑上下文信息

#### 变量

+ 简洁胜于冗长

![image-20220516120022895](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516120022895.png)

index和i仅作用于for，index并未增加对程序的理解

![image-20220516120223199](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516120223199.png)

deadline指截止时间，函数作用更明确

- 缩略词全大写，如HTTP不要Http，但当其位于变量开头不导出时用全小写，如xmlHTTPRequest代替XMLHTTPRequest
- 变量距离使用的地方越远，命名要越详细携带越多上下文信息

#### 函数

+ 函数名尽量简洁

- 函数名不携带包名上下文信息，因为包名和函数名总成对出现

> 如http包有Serve和ServeHTTP方法两个命名，应选择Serve命名，因为使用时会携带包名。类似C++命名空间，Java的包名

+ 名为foo的包某函数返回类型是Foo，可省略类型信息而不歧义
+ 名为foo的包某函数返回T（非Foo）可在函数名中加类型信息

#### 包

- 简短并包含一定上下文信息，但也要谨慎用缩写
- 只由小写字母组成（不包含大写、下划线等字符
- 不要和标准库同名冲突
- 不用常用变量名做包名如bufio而非buf
- 用单数而非复数，如encoding

### 控制流程

+ 避免嵌套，保持正常流程清晰

![image-20220516124608566](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516124608566.png)

+ 尽量保持正常代码路径为最小缩进，参考官方[代码](https://github.com/golang/go/blob/master/src/bufio/bufio.go#L277)

![image-20220516132503789](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516132503789.png)

嵌套使逻辑理复杂，调整后简单清晰，易于新增代码

![image-20220516132602548](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516132602548.png)

+ 互斥条件表驱动

如有并列的if嵌套逻辑：

```go
func CalculateByCmd(cmd string,a,b int)(int,error){
	if strings.EqualFold(cmd,"add"){
		return a+b,nil
	}
	if strings.EqualFold(cmd,"sub"){
		return a-b,nil
	}
	if strings.EqualFold(cmd,"mul"){
		return a*b,nil
	}
	return 0,errors.New("cmd not exist")
}
```

通过表驱动做出以下优化：

```go
var mapCalculate = map[string]func(a,b int) int{
	"add": func(a, b int) int {
		return a+b
	},
	"sub": func(a, b int) int {
		return a-b
	},
	"mul": func(a, b int) int {
		return a*b
	},
}

func CalculateByCmd(cmd string,a,b int)(int,error){
	if v,ok := mapCalculate[cmd];ok{
		return v(a,b),nil
	}
	return 0,errors.New("cmd not exist")
}
```

功能通过多个功能线性组合更简单，避免复杂嵌套分支，因为故障大多出现在复杂条件和循环语句，不易维护

### 错误处理

#### 错误的Wrap和Unwrap

fmt.Errorf用%w将错误关联到错误链，使每层调用方补充自己上下文，生成error跟踪链，参考官方[代码](https://github.com/golang/go/blob/master/src/cmd/go/internal/work/exec.go#L983)

#### error相关的函数

+ **errors.New()**：创建匿名变量直接表示错误，参考官方[代码](https://github.com/golang/go/blob/master/src/net/http/client.go#L802)

+ **errors.Is()**：判断错误断言，不同==，它能判断错误链中是否包含它，参考官方[代码](https://github.com/golang/go/blob/master/src/cmd/go/internal/modfetch/sumdb.go#L208)

![image-20220516134543214](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516134543214.png)

+ **errors.As()**：从错误链中提取想要的错误，参考官方[代码](https://github.com/golang/go/blob/master/src/errors/wrap_test.go#L255)

![image-20220516134607997](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516134607997.png)

#### panic和recover

+ panic：不推荐用panic，因为会向上传递到调用栈顶，若协程中所有被defer函数都不包含 recover 就会造成程序崩溃，启动阶段发生不可逆转错误时，可在 init 或 main 中用 panic，参考[代码](https://github.com/Shopify/sarama/blob/main/examples/consumergroup/main.go#L94)

![image-20220516135436991](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516135436991.png)

+ recover只能在被defer的函数中使用，嵌套无法生效，只在当前 goroutine 生效，参考官方[代码](https://github.com/golang/go/blob/master/src/fmt/scan.go#L247)，若需要更多上下文信息，可recover后在log中记录当前调用栈，参考官方[代码](https://github.com/golang/website/blob/master/internal/gitfs/fs.go#L228)

![image-20220516135537749](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516135537749.png)

## 性能优化

满足正确性、可靠性、健壮性、可读性等质量前提下，设法提高程序的效率，性能对比测试代码，可参考[代码](https://github.com/RaymondCode/go-practice)

### benchmark测试

性能需实际数据衡量，Go内置性能评估工具

![image-20220516135933452](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516135933452.png)

结果

![image-20220516135958879](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516135958879.png)

### slice

+ 参考[博客](https://ueokande.github.io/go-slice-tricks/)

+ 切片本质是一个有指针、长度、容量属性的数组的描述，切片操作时不复制切片指向的元素，而复用切片底层数组，尽可能用make() 初始化时提供容量信息，特别是append时防止拷贝：
  + append后长度<=cap时，直接利用原底层数组剩余空间
  + append后长度>cap时，分配更大区域容纳新底层数组
+ 在已有切片的基础上切片，若只用很小一段，但底层数组在内存仍占用大量空间无法释放，推荐用 copy 替代 re-slice

![image-20220516141608520](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516141608520.png)

### map

- 不断向map添加元素会触发扩容
- 根据实际需求提前预估好需要的空间
- 提前分配好空间可以减少内存拷贝和 Rehash

![image-20220516141742442](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516141742442.png)

### 字符串

- 常见的字符串拼接方式：strings.Builder、bytes.Buffer、+，效率递减
- 字符串在Go中是不可变类型，占用内存大小固定，用+拼接字符串会生成新的字符串，会开辟两字符串大小之和的新的空间，另两个内存是以倍数申请，底层都是[]byte 数组，bytes.Buffer 转化字符串时重新申请一块空间，存放生成的字符串，而 strings.Builder直接将底层[]byte 转换成字符串类型返回

![image-20220516142648934](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516142648934.png)

+ 字符串同样支持预分配，可进一步提高拼接性能

### 空结构体的使用

+ 空结构体不占内存仅作为占位符
+ 可以作为map实现简单set

![image-20220516142946389](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516142946389.png)

### atomic包

- 锁通过OS实现，属于系统调用，atomic通过硬件实现更高效
- sync.Mutex 应该用来保护一段逻辑，不仅仅用于保护一个变量
- 非数值可用 atomic.Value，它能承载一个interface{}

![image-20220516143302653](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516143302653.png)

### 总结

- 避免常见的性能陷阱可保证大部分程序的性能

- 普通应用代码，不要一味追求性能，应当在满足正确可靠、简洁清晰等质量前提下调优

## 性能调优实战

### 性能调优原则

- 要依靠数据不是猜测
- 要定位最大瓶颈而不是细枝末节
- 不要过早优化
- 不要过度优化

### 性能分析工具

性能调优的核心是性能瓶颈的分析，对于Go程序，最方便是 pprof 工具

- ##### pprof 功能说明

  - pprof 是用于可视化和分析性能分析数据的工具
  - 可以知道应用在什么地方耗费了多少 CPU、memory 等运行指标

![image-20220516145658070](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516145658070.png)

+ pprof实践

  - 排查 CPU 问题
    - 命令行分析
      - go tool pprof "[http://localhost:6060/debug/pprof/profile?seconds=10](https://link.juejin.cn?target=http%3A%2F%2Flocalhost%3A6060%2Fdebug%2Fpprof%2Fprofile%3Fseconds%3D10)"
    - top 命令
    - list 命令
    - 熟悉 web 页面分析
    - 调用关系图，火焰图
    - go tool pprof -http=:8080 "[http://localhost:6060/debug/pprof/cpu](https://link.juejin.cn?target=http%3A%2F%2Flocalhost%3A6060%2Fdebug%2Fpprof%2Fcpu)"

  - 排查堆内存问题
    - go tool pprof -http=:8080 "[http://localhost:6060/debug/pprof/heap](https://link.juejin.cn?target=http%3A%2F%2Flocalhost%3A6060%2Fdebug%2Fpprof%2Fheap)"

  - 排查协程问题
    - go tool pprof -http=:8080 "[http://localhost:6060/debug/pprof/goroutine](https://link.juejin.cn?target=http%3A%2F%2Flocalhost%3A6060%2Fdebug%2Fpprof%2Fgoroutine)"

  - 排查锁问题
    - go tool pprof -http=:8080 "[http://localhost:6060/debug/pprof/mutex](https://link.juejin.cn?target=http%3A%2F%2Flocalhost%3A6060%2Fdebug%2Fpprof%2Fmutex)"

  - 排查阻塞问题
    - go tool pprof -http=:8080 "[http://localhost:6060/debug/pprof/block](https://link.juejin.cn?target=http%3A%2F%2Flocalhost%3A6060%2Fdebug%2Fpprof%2Fblock)"

### pprof 的采样过程和原理

#### CPU 采样

![image-20220516150246286](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516150246286.png)

![image-20220516150336461](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516150336461.png)

启动采样时，进程向OS注册一个定时器，OS会每10ms向进程发送一个SIGPROF信号，进程接收到信号后就对当前调用栈进行记录。同时进程启动一个写缓冲的goroutine，它每隔100ms从进程中读取已记录的堆栈信息，并写入到输出流。当采样停止时，进程向OS取消定时器，不再接收信号，写缓冲读取不到新的堆栈时，结束输出。

#### 堆内存采样

![image-20220516150834598](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516150834598.png)

堆内存采样在实现上依赖内存分配器的记录，一些其他的内存分配，如栈内存、一些更底层使cgo调分配的内存，不会被采样记录，采样率是默认每分配512KB内存采样一次，采样率可以调整，设为1则每次分配都会记录。与CPU和goroutine都不同的是，内存的采样是个持续的过程，它记录从程序运行起的所有分配或释放的内存大小和对象数量，并在采样时谝历这些结果进行汇总

#### 协程和系统线程采样

![image-20220516151420613](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516151420613.png)

Goroutie采样会记录所有用户发起，也就是入口不是runtime开头的goroutine,以及main所在goroutine的信息和创建这些goroutine的调用栈
它们都是会在STW后，漏历所有goroutine/线程的列表〔图中的m就是GMP模型中的m，在golang中和线程对应）并输出堆栈，最后STW继续运行。该采样是立刻触发的全量记录，可以比较两个时间点的差值来得到某一时间段的指标

#### 阻塞操作和锁竞争采样

![image-20220516152319525](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220516152319525.png)

两指标在流程和原理上相似，不过指标的采样率含义不同：

+ 阻塞操作的采样率是个阈值，超过阈值时间的阻塞操作才会被记录，1为每次操作都会记录。炸弹程序的main里面设置rate=1

+ 锁竟争的采样率是个比例，运行时会通过随机数来记录固定比例的锁操作，1为每次操作都记录

  

实现也基本相同，在阻塞或锁操作发生时，会算出消耗的时间，连同调用栈一起主动上报给采样器，采样时，采样器会遍历已记录的信息，统计出具体操作次数、调用栈和总耗时。同样可以算两个时间点的差值算出段时间内的操作指标

### 性能调优案例

#### 基本概念

- 服务：能单独部署，承载一定功能的程序
- 依赖：Service A 的功能实现依赖 Service B 的响应结果，称为 Service A 依赖 Service B
- 调用链路：能支持一个接口请求的相关服务集合及其相互之间的依赖关系
- 基础库：公共的工具包、中间件

#### 业务优化

- 流程
  - 建立服务性能评估手段
  - 分析性能数据，定位性能瓶颈
  - 重点优化项改造
  - 优化效果验证
- 建立压测评估链路
  - 服务性能评估
  - 构造请求流量
  - 压测范围
  - 性能数据采集
- 分析性能火焰图，定位性能瓶颈
  - pprof 火焰图
- 重点优化项分析
  - 规范组件库使用
  - 高并发场景优化
  - 增加代码检查规则避免增量劣化出现
  - 优化正确性验证
- 上线验证评估
  - 逐步放量，避免出现问题
- 进一步优化，服务整体链路分析
  - 规范上游服务调用接口，明确场景需求
  - 分析业务流程，通过业务流程优化提升服务性能

#### 基础库优化

适应范围更广，覆盖更多服务，包括：

- AB 实验 SDK 的优化
  - 分析基础库核心逻辑和性能瓶颈
  - 完善改造方案，按需获取，序列化协议优化
  - 内部压测验证
  - 推广业务服务落地验证

- ##### Go 语言优化

  - 适应范围广通用性强，接入简单只需调整编译配置
  - 优化方式
    - 优化内存分配策略
    - 优化代码编译流程，生成更高效的程序
    - 内部压测验证
    - 推广业务服务落地验证

## 代码示例

以上语法和实战代码示例，都能再[这里](https://github.com/attackoncs/MyByteCamp/tree/main/Code/go-pprof-practice)找到

## 课后

- 了解下其他语言的编码规范，是否和 Go 语言编码规范有相通之处，注重理解哪些共同点

- 编码规范或者性能优化建议大部分是通用的，有没有方式能够自动化对代码进行检测？

- 从[链接](https://github.com/golang/go/tree/master/src)中选择感兴趣的包，看看官方代码是如何编写的


- 使用 Go 进行并发编程时有哪些性能陷阱或者优化手段？

- 真实线上环境中，遇到的性能问题各种各样，搜索下知名公司（如[uber](https://eng.uber.com/category/oss-projects/oss-go/)）的官方公众号或者博客，里面有哪些性能优化的案例？


- Go 语言本身在持续更新迭代，每个版本在性能上有哪些重要的优化点？

## 参考资料

- 熟悉 Go 语言基础后的必读[内容](https://go.dev/doc/effective_go)

- Dave Cheney 关于 Go 语言编程实践的演讲[记录](https://dave.cheney.net/practical-go/presentations/qcon-china.html)

- [《编程的原则：改善代码质量的101个方法》](https://mp.weixin.qq.com/s/vXSZOl2Gt7wcgq1OL9Cwow)，总结了很多编程原则，按照是什么 -> 为什么 -> 怎么做进行了说明

- [如何编写整洁的 Go 代码](https://github.com/Pungyeon/clean-go-article)

- [Go 官方博客](https://go.dev/blog/)，有关于 Go 的最新进展

- Dave Cheney 关于 Go 语言编程高性能编程的[介绍](https://go.dev/blog/)

- Go 语言高性能编程，博主总结了 Go 编程的一些性能[建议](https://geektutu.com/post/high-performance-go.html)

- Google 其他编程语言[编码规范](https://zh-google-styleguide.readthedocs.io/en/latest/)

