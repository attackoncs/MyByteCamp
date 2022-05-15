这是字节青训营的第四课：高性能Go语言发行版优化与落地实践的笔记和总结
# 概览

本节课主要介绍性能优化和实践，性能优化包括自动内存管理、Go内存管理及优化、编译器与静态分析、编译器优化，实践包括字节内部的Balanced GC优化对象分配以及编译器优化Beast mode

## 性能优化的基本问题

1. 性能优化是什么？

提升软件系统处理能力，减少不必要消耗，充分发掘计算机性能

2. 为什么要做性能优化？

+ 提升用户体验：让刷抖音更丝滑不卡顿
+ 高效利用资源：降低成本提高效率，小优化乘海量机器也会显著节约成本

## 性能优化的两个层面
![在这里插入图片描述](https://img-blog.csdnimg.cn/932679e5c6d14ccfbeb44f12c21ea6c0.png)
代码和软件层面如上图，重点关注**业务层代码**和**语言运行时**优化。前者针对特定场景和问题，容易获得较大性能收益。后者解决更通用的性能问题，考虑更多场景，因此也会有更多权衡取舍。优化的原则是以**数据驱动**，包括三点：优化应以数据为衡量标准而非猜测、善用自动化性能分析工具pprof、优先优化最大瓶颈

## 性能优化的可维护性

软件质量至关重要，因此性能优化应以**保证接口稳定的前提下改进**，即保证可维护性，具体措施包括：

+ 测试用例需尽可能多的覆盖场景，保证功能一致，方便回归
+ 文档需详细记录做了什么，没做什么以及能达到什么效果
+ 通过选项控制是否开启优化，达到隔离性
+ 必要的日志输出使得优化结果可观测

# 自动内存管理

自动内存管理指**由程序语言的运行时系统管理程序运行时按需分配的动态内存**，避免手动内存管理，避免CVE中大量出现的double-free和use-after-free问题，保证内存使用的正确性和安全性，从而专注业务逻辑，主要包括三个任务：

+ 为新对象分配空间
+ 找到存活对象
+ 回收死亡对象的内存空间

## 相关概念
![在这里插入图片描述](https://img-blog.csdnimg.cn/69e1dfe7134342819f12889f7d4680f0.png)
+ Mutator：业务线程（实际是oroutines），分配新对象，修改对象指向关系

+ Collector：GC线程，找到存活对象，回收死亡对象内存空间

+ Serial GC：只有一个collector

+ Parallel GC：多个collectors同时回收

+ Concurrent GC：mutator和collector同时执行，collector需感知对象指向关系的变化

- 评价GC算法
  - Safety：安全性，不能回收存活的对象，是基本要求
  - Throughput：花在GC上的时间，1-（GC时间/程序执行总时间）
  - Pause time：暂停时间，业务是否感知
  - Space overhead：内存开销，GC元数据的开销
## 追踪垃圾回收

指针指向关系不可达的对象是对象被回收的条件，**标记**找到可达对象，**清理**所有不可达对象
![在这里插入图片描述](https://img-blog.csdnimg.cn/5e3ef97cb19341b68e3c25645308e729.png)
+ 标记：先标记根对象（静态变量、全局变量、常量、线程栈等），从根对象出发找到所有可达对象

+ 清理：清理所有不可达对象，包括Copying GC、Mark-sweep GC、Mark-compact GC

  + Copying GC：将存活对象赋值到另外内存空间
![在这里插入图片描述](https://img-blog.csdnimg.cn/ef2c0cc586744c81ba619f4c970a8f9a.png)
+ Mark-sweep GC：使用free list管理空闲内存
![在这里插入图片描述](https://img-blog.csdnimg.cn/13027357f1a648b389c4914ab1c71f3c.png)
+ Compact GC：原地整理对象
![在这里插入图片描述](https://img-blog.csdnimg.cn/bf175a0837c943c6bee886a86d53feb0.png)
## 分代GC（Generational GC）

基于分代假说：大多对象分配后很快就不再使用，因此根据对象年龄（经历GC的次数），对不同代对象指定不同GC策略，不同年龄对象放于heap不同区域，降低整体内存管理开销，具体来说：

+ 年轻代：常规对象分配，由于存活对象少，可采用copying collection，GC吞吐率高
+ 老年代：对象趋向于一直活着，反复复制开销较大，可采用mark-sweep collection

## 引用计数

每个对象都有个与之关联的引用数，当且仅当引用数大于0对象存活

+ 优点：不需了解内存管理实现细节，操作被平摊到执行过程
![在这里插入图片描述](https://img-blog.csdnimg.cn/55d3f6a2111b42b6b3f78fc9b712b338.png)
+ 缺点：原子操作保证引用计数原子性和可见性，额外空间存储引用数，故维护引用计数开销大，需用弱引用回收环形引用，回收内存时依然可能引发暂停
![在这里插入图片描述](https://img-blog.csdnimg.cn/a981b54e7eab4596bf37d93cb65c144f.png)
# Go内存管理及优化

## Go内存分配

### 分块

- 目标：为对象在heap上分配内存

- 提前将内存分块

  - 调用系统调用`mmap()`，向OS申请一大块内存，如4MB
  - 先将内存划分为大块，如8KB，称作`mspan`
  - 再将大块继续划分为特定大小的小块，用于对象分配
  - `noscan mspan`：分配不包含指针的对象—GC不需扫描
  - `scan mspan`：分配包含指针的对象—GC需扫描
![在这里插入图片描述](https://img-blog.csdnimg.cn/d0069d99366a45b8b0973d3cd6f745e4.png)
- 对象分配：根据对象的大小，选择最合适的块返回

### 缓存

+ 类似TCMalloc：`thread caching`

+ 每个p包含一个`mcache`用于快速分配，用于为绑定于p上的g分配对象

+ mcache管理一组mspan

+ 当`mcache`中的`mspan`分配完毕，向`mcentral`申请带有未分配块的`mspan`

+ 当`mspan`中没有分配的对象，`mspan`会被缓存在`mcentral`中，而不是立即释放归还给OS
![在这里插入图片描述](https://img-blog.csdnimg.cn/ad2f74bfc7494305b851bb5a9d551102.png)
## Go内存管理优化

+ 对象的分配是非常高频的操作：线上业务每秒分配GB级内存

+ **小对象占比高**，大多都低于80B

+ Go内存分配比较耗时

  - 分配路径长：`g -> m -> p -> mcache -> mspan -> memory block -> return pointer`

  - pprof：对象分配的函数是最频繁调用的函数之一（占用很多的CPU）

## 字节优化方案Balanced GC

高峰期CPU是用来较低4.6%，核心接口时延下降4.5%-7.7%，具体细节如下：

+ 每个g都绑定一大块内存（1KB），称作`goroutine allocation buffer`（GAB）

+ GAB用于noscan类型的小对象分配：<128B

+ 使用三个指针维护GAB：base，end，top，指针碰撞风格对象分配

  - 无须和其他分配请求互斥

  - 分配动作简单高效
![在这里插入图片描述](https://img-blog.csdnimg.cn/6df1cc2e6c364b39a3159bf5cf4f06ed.png)
+ GAB对于Go内存管理来说是个大对象，本质是将多个小对象的分配合并成一次大对象的分配，会导致内存延迟释放
![在这里插入图片描述](https://img-blog.csdnimg.cn/a38c0c93daf04a54a6381388b673bfd9.png)
+ 方案：移动GAB中存活的对象

	+ GAB总的大小超过一定阈值时，将GAB中存活的对象复制到另外分配的GAB中
	+ 原先的GAB可以释放，避免内存泄漏
	+ 本质：用copying GC的算法管理小对象（根据对象的生命周期，使用不同的标记和清理策略）
![在这里插入图片描述](https://img-blog.csdnimg.cn/fcf0429e0d594b1482ddab58c205bfb1.png)
# 编译器和静态分析

## 基本介绍

编译器是识别符号语法和非法，生成正确且高效代码的重要系统软件，分前后端，主要学习后端优化
![在这里插入图片描述](https://img-blog.csdnimg.cn/6e57e49b45b44e82bf4f371f3533eee4.png)
+ 静态分析：**不执行代码**推导程序行为，分析程序性质
## 数据流和控制流

- 控制流：程序的执行流程

- 数据流：数据在控制流上的传递
![在这里插入图片描述](https://img-blog.csdnimg.cn/e55cf05185664e1585049035222d78c8.png)
对应的控制流图
![在这里插入图片描述](https://img-blog.csdnimg.cn/5dde9cae29484c7394e72d2ae42d1333.png)
通过分析控制流和数据流，可知道更多用来优化程序的性质
![在这里插入图片描述](https://img-blog.csdnimg.cn/83d7b05d5c2742e6bd678a5f0af14056.png)
通过分析知程序始终返回 4。编译器根据该结果优化
## 过程内和过程间分析

- Intra-procedural analysis: 过程内分析，仅在过程内分析控制流和数据流

- Inter-procedural analysis: 过程间分析需同时分析数据流和控制流，联合求解较复杂，如参数传递，函数返回值等
![在这里插入图片描述](https://img-blog.csdnimg.cn/39e3db0edff040d69637e948cd229151.png)
数据流分析i的类型，才确定新控制流即调用哪个foo

# Go编译器优化

## 函数内联

### 为什么做编译器优化

- 用户无感知，重新编译即可获得性能收益

- 通用性优化

### 现状

- 采用的优化较少

- 追求编译时间短，因此没有进行复杂的代码分析和优化

### 编译优化的思路

- 场景：面向后端长期执行的任务

- 取舍：用编译时间换取更高效的代码

### Beast mode

+ **函数内联**
+ **逃逸分析**
+ 默认栈大小调整
+ 边界检查消除
+ 循环展开

#### 函数内联

- 定义：将被调用函数的函数体（callee）的副本替换到调用位置（caller）上，同时重写代码以反映参数的绑定

- 优点
  - 消除调用开销，如参数传递、保存寄存器等
  - 将过程间分析的问题转换为过程内分析，帮助其他优化，如逃逸分析
![在这里插入图片描述](https://img-blog.csdnimg.cn/a22f5a60630f4de598fbdfe713dcac7f.png)
用micro-benchmark验证，内联后性能提升4.5倍

![在这里插入图片描述](https://img-blog.csdnimg.cn/7505c15f52ff4a428d31b2c899330bbe.png)
- 缺点
  - 函数体变大，instruction cache(icache)不友好
  - 编译生成的 Go 镜像文件变大

- 函数内联在大多数情况下是正向优化，即多内联，会提升性能

- 采取一定的策略决定是否内联
  - 调用和被调用函数的规模

- Go 内联的限制
  - 语言特性：interface, defer 等等，限制了内联优化
  - 内联策略非常保守

- 字节跳动的优化方案Beast Mode
  - 修改了内联策略，降低函数调用开销，让更多函数被内联
  - 增加其他优化的机会：逃逸分析

- 开销
  - Go 镜像大小略有增加~10%
  - 编译时间增加

#### 逃逸分析

- 定义：分析代码中指针的动态作用域，即指针在何处可以被访问

- 大致思路
  - 从对象分配处出发，沿着控制流，观察数据流。若发现指针 p 在当前作用域 s:
    - 作为参数传递给其他函数；
    - 传递给全局变量；
    - 传递给其他的 goroutine;
    - 传递给已逃逸的指针指向的对象；
  - 则指针 p 逃逸出 s，反之则没有逃逸出 s.
- Beast Mode：函数内联拓展函数边界，更多对象不逃逸，高峰CPU usage降低9%，时延降低10%，内存使用降低3%

- **优化：未逃逸出当前函数的指针指向的对象可以在栈上分配**
  - 对象在栈上分配和回收很快：移动 sp 即可分配和回收
  - 减少在堆上分配对象，降低 GC 负担

# 课后

1. 从业务层和语言运行时层进行优化分别有什么特点？

1. 从软件工程的角度出发，为了保证语言 SDK 的可维护性和可拓展性，在进行运行时优化时需要注意什么？

1. 自动内存管理技术从大类上分为哪两种，每一种技术的特点以及优缺点有哪些？

1. 什么是分代假说？分代 GC 的初衷是为了解决什么样的问题？

1. Go 是如何管理和组织内存的？

1. 为什么采用 bump-pointer 的方式分配内存会很快？

1. 为什么我们需要在编译器优化中进行静态代码分析？

1. 函数内联是什么，这项优化的优缺点是什么？

1. 什么是逃逸分析？逃逸分析是如何提升代码性能的？

# 参考文献

1. The Garbage Collection Handbook -- the art of automatic memory management

自动内存管理领域的集大成之作。把自动内存管理的问题、动机、方案、以及最新研究进展和方向进行了非常详尽的阐述。整个书很好读，参考文献非常充实，推荐阅读英文版。

2. JEP 333: ZGC: A Scalable Low-Latency Garbage Collector [openjdk.java.net/jeps/333](https://link.juejin.cn?target=https%3A%2F%2Fopenjdk.java.net%2Fjeps%2F333)

目前 HotSpot JVM 上 pauseless GC 实现的 proposal，可以看作 GC 领域比较新的工程方面的进展。

3. 数据密集型应用系统设计 Designing Data-Intensive Applications: The Big Ideas Behind Reliable, Scalable, and Maintainable Systems

通过例子带大家理解互联网产品需要解决的问题以及方案。

4. 编译原理 The Dragon book, Compilers: Principles, Techniques, and Tools

在编译器前端着墨较多。本书第二版的第九章 **机器无关优化**，推荐大家反复仔细阅读。这一章主要讲述的是编译优化中常见的数据流分析问题，建议大家跟随书本仔细推导书中的例子，会帮助你对数据流分析有个大致的认识。这一章给出的引用文献大多是编译和静态分析领域非常有影响力的论文，有兴趣的同学可以阅读。

5. 编译原理 Principles and Techniques of Compilers [silverbullettt.bitbucket.io/courses/com…](https://link.juejin.cn?target=https%3A%2F%2Fsilverbullettt.bitbucket.io%2Fcourses%2Fcompiler-2022%2Findex.html)

南京大学编译原理课程。

6. 静态程序分析 Static Program Analysis [pascal-group.bitbucket.io/teaching.ht…](https://link.juejin.cn?target=https%3A%2F%2Fpascal-group.bitbucket.io%2Fteaching.html)

南京大学静态程序分析课程。参考文献 4 数据流分析读不懂的地方可以参考本课程的课件。

7. 编译器设计 Engineering a Compiler

在编译器后端优化着墨较多。可以帮助大家理解后端优化的问题。

8. JVM Anatomy Quark #4: TLAB allocation [shipilev.net/jvm/anatomy…](https://link.juejin.cn?target=https%3A%2F%2Fshipilev.net%2Fjvm%2Fanatomy-quarks%2F4-tlab-allocation%2F)

Goroutine allocation buffer (GAB) 的优化思路在 HotSopt JVM 也能找到类似的实现。

9. Constant folding, [en.wikipedia.org/wiki/Consta…](https://link.juejin.cn?target=https%3A%2F%2Fen.wikipedia.org%2Fwiki%2FConstant_folding)

常量折叠数据流分析。

10. Choi, Jong-Deok, et al. "Escape analysis for Java." *Acm Sigplan* *Notices* 34.10 (1999): 1-19.

逃逸分析的 Java 实现。

11. Zhao, Wenyu, Stephen M. Blackburn, and Kathryn S. McKinley. "Low-Latency, High-Throughput Garbage Collection." (PLDI 2022). 学术界和工业界在一直在致力于解决自动内存管理技术的不足之处，感兴趣的同学可以阅读。

# 代码示例
以上语法和实战代码示例，都能再[这里](https://github.com/attackoncs/MyByteCamp)找到