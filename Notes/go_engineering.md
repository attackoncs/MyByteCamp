这是字节青训营的第二课：工程实践的笔记和总结
# 1语言进阶

## 并发

go可以充分发挥多核优势，高效运行，介绍原理之前，先介绍几个概念：

+ 并发：单CPU同时间段切换交错执行多任务
+ 并发：多CPU同时刻都运行任务
+ 进程：资源管理的最小单位，进程虚拟地址空间分成用户和内核空间
+ 线程：资源调度最小单位，再内核态，共享进程中的资源，比进程轻量级
+ 协程：再用户态的轻量级线程，调度和切换都在用户态，因此高效
![在这里插入图片描述](https://img-blog.csdnimg.cn/1923463d6d8e4e0cacfe9da3a44beaf3.png)
### goroutine

go中的协程goroutine使用简单，只需在调用函数（普通或匿名函数）前加go关键字，没有返回值，因此通过channel通信。提倡通过**通信共享内存**，而不是通过共享内存而实现通信。
![在这里插入图片描述](https://img-blog.csdnimg.cn/afc6088cfb034becb9ebc7a47bc806ca.png)
由于goroutine是异步执行，因此需要同步，否则主程序退出而goroutine还没执行完，同步的主要方式有：Sleep、Channel、Sync。若知道每个协程执行时间，则可Sleep等待所有协程执行完再退出，实际不常用，常用后两者

### Channel

go是CSP并发模型，没有对内存加减锁减少性能消耗，chanel可让协程发送特定值到另外协程，遵循FIFO且保证收发顺序。通过make创建，分无缓冲（未指定大小）和有缓冲，前者会阻塞直到接收或发送，使用“<-”发送和接收。常用for range或if ok判断channel关闭时机，当然defer close不用判断，图是生产者-消费者模型
![在这里插入图片描述](https://img-blog.csdnimg.cn/802f99c05a9840d4b56d2b4db45d329a.png)
### Sync

Sync保重WaitGroup内部维护计数器，通过主程序Add和协程Done增加较少计数器，主程序Wait阻塞等待任务执行完，进行同步，针对只执行一次场景使用Once
![在这里插入图片描述](https://img-blog.csdnimg.cn/516501324b2b4c609858c88052f58180.png)
## 总结

协程：通过高效的调度模型实现高并发；channel：通过通信实现共享内存；sync：实现并发安全操作和携程间的同步

# 2依赖管理

## 背景

实际开发需通常更关注业务逻辑的实现，常使用被封装好、经过验证的开发组件或工具提升开发效率（不然还不如使用C和C++自己造轮子:）），如框架、日志、driver等一系列依赖通过sdk引入，因此需要对依赖管理。

## Go依赖管理演进

Go的依赖管理经过GOPATH、Go Vendor、Go Module，GOPAHT无法实现包的多版本控制，Go Vendor无法控制依赖的版本，更新有可能出现依赖冲突等问题，因此常使用Go Module。整个演进路线围绕：

+ 不同项目依赖的版本不同
+ 控制依赖库的版本

### 依赖管理三要素

1. 描述依赖的配置文件
![在这里插入图片描述](https://img-blog.csdnimg.cn/b00f0c7c78964972b6b9077978971927.png)
依赖版本包括语义化版本和基于commit伪版本版本
![在这里插入图片描述](https://img-blog.csdnimg.cn/7d2201e0d4de4fa092e493b9d126a8e8.png)
语义版本中的MAJOR不同表示API不兼容，即使同个库，MAJOR不同也被认为是不同模块，MINOR通常是新增函数或功能，path一般是修复bug。基于commit则是commit的时间戳和12位的哈希前缀校验码，每次commit就默认生成一个版本号
![在这里插入图片描述](https://img-blog.csdnimg.cn/e3a650547971443f8e78a57bf9d0ce2d.png)
indirect表示间接依赖
![在这里插入图片描述](https://img-blog.csdnimg.cn/0a7cd644b3bb4a74bfdf86996cc5f4d3.png)
go是11版本提出go module，主版本大于等于2的包，应在路径中体现出版本，而很多包在此之前打上更高版本的tag，为兼容这些包会在版本号后奖赏+incompatible，如图上因该是lib6/v3 v3.2.0，若未遵守则打上incompatible标签
![在这里插入图片描述](https://img-blog.csdnimg.cn/0c20ebf3600747148d66781c288f1fb1.png)
若同包不同版本，则选择最低兼容版本

2. 管理依赖库的中心仓库

go proxy解决无法保证构建稳定性、依赖可用性、第三方托管平台压力的问题，它是服务站点，会缓存源站中软件内容，构建时直接从proxy站点拉取依赖
![在这里插入图片描述](https://img-blog.csdnimg.cn/f077393e32e74a979096c81b4b9895c8.png)
GOPROXY环境变量中是go proxy站点URL列表，按顺序查找，direct表示源站下载

3. 本地工具

go get获取包，go mod初始化（init）、下载（download）、增加减少依赖（tidy）

# 3测试

测试可极大避免事故发生，是避免事故最后一道屏障，从上到下是回归测试、集成测试、单元测试，覆盖率逐层变大成本逐层降低

## 单元测试
![在这里插入图片描述](https://img-blog.csdnimg.cn/c999bcead4a34d62a95bf860dd740a15.png)
单元测试包括：输入、测试单元、输出三部分，单元包括函数、模块等，通过校对保证代码功能和预期相同，既能保证质量又能提升效率（定位和修复bug）,常用测试包assert，单元测试规则：

+ 测试文件以_test.go结尾

+ 测试函数以Test开头且连接的首字母大写

+ 初始化逻辑放在TestMain中
![在这里插入图片描述](https://img-blog.csdnimg.cn/5d70dccd2a5544dc819640f79d1dc5ff.png)
单元测试覆盖率指代码执行量和总代码量间的比率，主要包括分支、行、方法、类四个指标，实际项目一般要求50-60%覆盖率，资金型重要服务需达到80%，为提升覆盖率：

+ 测试分支相互独立、全面覆盖
+ 测试单元粒度要够小，遵循函数单一职责

## Mock测试

复杂项目一般会依赖文件、数据库、缓存等，而单测需保证任何情况都能运行测试的稳定，每次测试结果都相同的幂等，因此需要mock机制，常用monkey库，运行时通过Go的unsafe包，将函数地址替换为运行时函数地址，跳转到待打桩函数或方法，摆脱依赖

## 基准测试

实际项目开发中经常遇到代码性能瓶颈，为定位问题需对代码做性能分析，也就是基准测试，测试程序运行即耗费cpu的程度

# 4项目实战

项目实战是通过项目实践讲解项目开发的思路和流程，包括需求分析、代码开发、测试运行三部分。

## 需求背景

### 需求描述

掘金社区话题页面，展示话题（标题、文字描述）和回帖列表，不考虑前端仅实现本地web服务，话题和回帖数据用文件存储

### 需求用例

主要是用户浏览话题内容和回帖列表，包含主题内容和回帖列表，想象每个实体的属性及它们之间的联系
![在这里插入图片描述](https://img-blog.csdnimg.cn/f3080c3243194c7cb12da54f1ba38393.png)
### ER图

思考实体的属性及之间的联系，对后续开发提供清晰思路
![在这里插入图片描述](https://img-blog.csdnimg.cn/d5b38b281c204dfca3954fed41cc97bb.png)
### 分层结构 

代码结构采用分层结构设计，包括数据层、逻辑层、视图层。数据层封装外部数据增删改查，对逻辑层屏蔽底层数据差异，即不管底层是文件还是数据库还是微服务，同时对逻辑层接口模型不变。逻辑层处理核心业务逻辑并上送给视图层。视图层负责交互，以视图形式返回给客户端
![在这里插入图片描述](https://img-blog.csdnimg.cn/036408917c334032a2b12f911831e3a3.png)
### 组件工具

高性能web框架gin v1.3.0版，主要涉及路由分发，使用go module依赖管理，按reposity、service、controller逐步实现

### Repository

定义Topic和Post结构体如下，如何高效查询
![在这里插入图片描述](https://img-blog.csdnimg.cn/4a279ca64c024e38a6feca653e09d4d0.png)
为简单使用map实现内存索引，用文件元数据初始化全局内存索引，可o(1)查找
![在这里插入图片描述](https://img-blog.csdnimg.cn/3b10ba6927be4541a6026ebb479cb55e.png)
迭代遍历数据行，转为结构体存储到map中

```go
func initTopicIndexMap(filePath string) error {
	open, err := os.Open(filePath + "topic")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(open)
	topicTmpMap := make(map[int64]*Topic)
	for scanner.Scan() {
		text := scanner.Text()
		var topic Topic
		if err := json.Unmarshal([]byte(text), &topic); err != nil {
			return err
		}
		topicTmpMap[topic.Id] = &topic
	}
	topicIndexMap = topicTmpMap
	return nil
}
```

查询则直接查询key获得value即可

### Service

Service主要包括PageInfo结构体，它包含Topic和PostList

```go
type PageInfo struct {
	Topic    *repository.Topic
	PostList []*repository.Post
}
```

Service实现流程包括参数校验、准备数据、组装实体三部分

```go
func (f *QueryPageInfoFlow) Do() (*PageInfo, error) {
	if err := f.checkParam(); err != nil {
		return nil, err
	}
	if err := f.prepareInfo(); err != nil {
		return nil, err
	}
	if err := f.packPageInfo(); err != nil {
		return nil, err
	}
	return f.pageInfo, nil
}
```

prepreinfo方法中，话题和回帖信息的获取都依赖tipicid，这样两者就可并行执行，提高效率，实际开发中要思考流程是否可并发，从而提高并发
![在这里插入图片描述](https://img-blog.csdnimg.cn/453d16ed4764403ba5bbaabed37d9f14.png)
### Controller

定义PageData结构体作为view对象，通过code和msg打包业务状态信息，data承载业务实体信息

```go
type PageData struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
func QueryPageInfo(topicIdStr string) *PageData {
	topicId, err := strconv.ParseInt(topicIdStr, 10, 64)
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	pageInfo, err := service.QueryPageInfo(topicId)
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	return &PageData{
		Code: 0,
		Msg:  "success",
		Data: pageInfo,
	}

}
```



### Router

最后是web服务的引擎配置，path映射到具体controller，通过path变量传递话题id。过程包括初始化数据索引、初始化引擎配置、构建路由、启动服务

```go
func main() {
	if err := Init("./data/"); err != nil {
		os.Exit(-1)
	}
	r := gin.Default()
	r.GET("/community/page/get/:id", func(c *gin.Context) {
		topicId := c.Param("id")
		data := cotroller.QueryPageInfo(topicId)
		c.JSON(200, data)
	})
	err := r.Run()
	if err != nil {
		return
	}
}
```

### 测试和运行

本地go run运行，并使用curl请求服务暴露的接口
![在这里插入图片描述](https://img-blog.csdnimg.cn/4c816a28277f4e7d8e3749b008723279.png)
# 5 代码示例

以上语法和实战代码示例，都能再[这里](https://github.com/attackoncs/MyByteCamp/tree/main/Code/go-project-example-0)找到