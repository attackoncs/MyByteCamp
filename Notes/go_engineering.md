# 1语言进阶

## 并发

go可以充分发挥多核优势，高效运行，介绍原理之前，先介绍几个概念：

+ 并发：单CPU同时间段切换交错执行多任务
+ 并发：多CPU同时刻都运行任务
+ 进程：资源管理的最小单位，进程虚拟地址空间分成用户和内核空间
+ 线程：资源调度最小单位，再内核态，共享进程中的资源，比进程轻量级
+ 协程：再用户态的轻量级线程，调度和切换都在用户态，因此高效

![image-20220509180825891](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509180825891.png)

### goroutine

go中的协程goroutine使用简单，只需在调用函数（普通或匿名函数）前加go关键字，没有返回值，因此通过channel通信。提倡通过**通信共享内存**，而不是通过共享内存而实现通信。

![image-20220509180118669](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509180118669.png)



由于goroutine是异步执行，因此需要同步，否则主程序退出而goroutine还没执行完，同步的主要方式有：Sleep、Channel、Sync。若知道每个协程执行时间，则可Sleep等待所有协程执行完再退出，实际不常用，常用后两者

### Channel

go是CSP并发模型，没有对内存加减锁减少性能消耗，chanel可让协程发送特定值到另外协程，遵循FIFO且保证收发顺序。通过make创建，分无缓冲（未指定大小）和有缓冲，前者会阻塞直到接收或发送，使用“<-”发送和接收。常用for range或if ok判断channel关闭时机，当然defer close不用判断，图是生产者-消费者模型

![image-20220509182722766](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509182722766.png)

### Sync

Sync保重WaitGroup内部维护计数器，通过主程序Add和协程Done增加较少计数器，主程序Wait阻塞等待任务执行完，进行同步，针对只执行一次场景使用Once

![image-20220509183407491](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509183407491.png)

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

![image-20220509184909390](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509184909390.png)

依赖版本包括语义化版本和基于commit伪版本版本

![image-20220509185135396](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509185135396.png)

语义版本中的MAJOR不同表示API不兼容，即使同个库，MAJOR不同也被认为是不同模块，MINOR通常是新增函数或功能，path一般是修复bug。基于commit则是commit的时间戳和12位的哈希前缀校验码，每次commit就默认生成一个版本号

![image-20220509200131126](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509200131126.png)

indirect表示间接依赖

![image-20220509200722048](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509200722048.png)

go是11版本提出go module，主版本大于等于2的包，应在路径中体现出版本，而很多包在此之前打上更高版本的tag，为兼容这些包会在版本号后奖赏+incompatible，如图上因该是lib6/v3 v3.2.0，若未遵守则打上incompatible标签

![image-20220509201146421](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509201146421.png)

若同包不同版本，则选择最低兼容版本

2. 管理依赖库的中心仓库

go proxy解决无法保证构建稳定性、依赖可用性、第三方托管平台压力的问题，它是服务站点，会缓存源站中软件内容，构建时直接从proxy站点拉取依赖

![image-20220509201824825](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509201824825.png)

GOPROXY环境变量中是go proxy站点URL列表，按顺序查找，direct表示源站下载

3. 本地工具

go get获取包，go mod初始化（init）、下载（download）、增加减少依赖（tidy）

# 3测试

测试可极大避免事故发生，是避免事故最后一道屏障，从上到下是回归测试、集成测试、单元测试，覆盖率逐层变大成本逐层降低

## 单元测试

![image-20220509204331176](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509204331176.png)

单元测试包括：输入、测试单元、输出三部分，单元包括函数、模块等，通过校对保证代码功能和预期相同，既能保证质量又能提升效率（定位和修复bug）,常用测试包assert，单元测试规则：

+ 测试文件以_test.go结尾

+ 测试函数以Test开头且连接的首字母大写

+ 初始化逻辑放在TestMain中

  

  ![image-20220509204805964](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509204805964.png)

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

![image-20220509212308404](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509212308404.png)

### ER图

思考实体的属性及之间的联系，对后续开发提供清晰思路

![image-20220509212447612](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509212447612.png)

### 分层结构 

代码结构采用分层结构设计，包括数据层、逻辑层、视图层。数据层封装外部数据增删改查，对逻辑层屏蔽底层数据差异，即不管底层是文件还是数据库还是微服务，同时对逻辑层接口模型不变。逻辑层处理核心业务逻辑并上送给视图层。视图层负责交互，以视图形式返回给客户端

![image-20220509213438571](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509213438571.png)

### 组件工具

高性能web框架gin v1.3.0版，主要涉及路由分发，使用go module依赖管理，按reposity、service、controller逐步实现

### Repository

定义Topic和Post结构体如下，如何高效查询

![image-20220509214046408](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509214046408.png)

为简单使用map实现内存索引，用文件元数据初始化全局内存索引，可o(1)查找

![image-20220509214325670](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509214325670.png)

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

![image-20220509215735509](C:/Users/yl/AppData/Roaming/Typora/typora-user-images/image-20220509215735509.png)

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

# 5 代码示例

以上语法和实战代码示例，都能再[这里](https://github.com/attackoncs/MyByteCamp/tree/main/Code/go-project-example-0)找到

