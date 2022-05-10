有幸参加了字节跳动举办的青训营活动，主要是go语言的编程实践，我将会整理课程的笔记和总结，欢迎关注！
# 1Go语言简介

由于Go语言有语法简单、高性能等特点，因此国内外各大公司如谷歌、腾讯、字节等都在使用，特别是字节全面拥抱Go，原因是最初因性能问题将Python换成Go，而Go学习简单，性能高，且部署简单。总的来说Go语言特性如下：

+ 语法简单、学习曲线平缓

+ 高性能、高并发

+ 丰富的标准库
+ 完善的工具链
+ 静态链接
+ 快速编译
+ 跨平台
+ 垃圾回收

个人而言，因我是C++出身，对C++复杂的语法感到麻木，因此特别喜欢Go语法简单、上手快的特点，几小时就能上手，再加上Go天生高并发，有丰富标准库等特点，几近完美，因此建议大家学下Go

# 2入门

## 2.1开发环境

本地编译环境：Golang、VSCode

云端开发环境：gitpod

## 2.2基础语法

### Hello World

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world")
}
```

第一行代表属于main包，它是程序入口包，import导入fmt包，并再main函数中调用其Println函数。使用go build或go run编译或编译并运行代码

### 变量

go是强类型，变量都有类型，声明变量时常用方式：

+ var varname type = value：需var关键字且类型放后面
+ varname := value：会自行推到类型

### 字符串

go的strings包有很多常用函数：contains判断是否包含另一字符串、count字符串计数、index查找字符串位置、join连接多字符串、repeat重复多个字符串、replace替换字符串

### 字符串格式化

fmt包中有很多字符串格式化相关的方法，但常用“%v”打印任意类型变量，不需区分类型，“%+v”和“%#v”则更详细

### 数字解析

字符串和数字间转换通过strconv包的函数，Atoi表示字符串转数字相反则Itoa

### if else判断

if后无括号，写括号的话编辑器会自动去掉，且后面必须接大括号

### switch

switch后变量名无括号，case中不加break，且可以使用任意变量类型，当判断较多时可用switch完全替代if else

### for循环

只有唯一的for循环，for中三语句任何语句都可省略且也没括号，循环中可用break或continue跳出或继续循环

### range

常用range快速遍历数组、slice、map等，返回索引和值，可用“_”忽略索引

### 数组

类似其他语言中固定长度数组，实际很少使用，更多使用切片

### 切片

slice长度不固定，用make创建，append时自动扩容，可像python中切片一样截取

### map

实际使用最频繁的数据结构，用make创建，常用val,ok:=m[key]写法，通过判断ok判断是否存在key对应的val

### 函数

参数和返回值类型后置，支持返回多个值，实际代码中几乎所有函数返回两个值：结果和错误信息，注意参数是都值传递，若想修改参数，需传递指针。首字母大写表示公共函数，类似C++中public函数，可包外访问

### 指针

类似C/C++指针，但操作有限，指针主要用途是对传入参数修改

### 结构体

结构体初始化时需传入每个字段初始值，也可键值对的方式初始化，通常用结构体的指针作为参数，既能修改结构体又能避免结构体的拷贝开销，字段首字母大写表示公开字段，类似C++中public成员

### 结构体方法

再函数名前带上结构体参数和括号，就实现类似其他语言中的类成员函数

### 错误处理

go常见做法是再传递结果的返回值的基础上，新增一个传递错误信息的返回值，这样能清晰知道哪个函数出错，且能用简单的if else处理错误

### defer处理

go中为避免资源泄露，需用defer手动关闭资源，会在函数运行结束后执行

### json处理

对已有结构体，保证每个字段首字母大写，是公开字段，就能用marshaler序列化变成json字符串，unmarshaler反序列化，可用json tag修改输出结果的字段名

### 时间处理

time包包含时间处理各种函数，用Date构造带时区时间，Now获取当前时间

### 进程处理

os包包含进程处理，Args得到命令参数，Getenv读环境变量，Command执行命令

# 3实战项目

实战项目主要是为巩固基础语法，并且是逐步开发、快速迭代改进，很有参考价值

## 猜数字游戏

设置随机数种子，随机生成100以内的数，提示用户输入，并和用户输入的数据比大小，根据大小提示信息，不断循环，最终猜正确才退出循环。通过该例子巩固变量循环、函数控制流和错误处理等知识

```go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	maxNum := 100
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(maxNum)
	// fmt.Println("The secret number is ", secretNumber)

	fmt.Println("Please input your guess")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			continue
		}
		input = strings.TrimSuffix(input, "\n")

		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter an integer value")
			continue
		}
		fmt.Println("You guess is", guess)
		if guess > secretNumber {
			fmt.Println("Your guess is bigger than the secret number. Please try again")
		} else if guess < secretNumber {
			fmt.Println("Your guess is smaller than the secret number. Please try again")
		} else {
			fmt.Println("Correct, you Legend!")
			break
		}
	}
}

```

## 简单字典

抓包彩云小译翻译的api，并拷贝curl命令，通过[网站](https://curlconverter.com/)转换成go代码，并拷贝服务端响应的json数据，通过[网站](https://oktools.net/json2go)生成go结构体，然后解析对应的字段输出，通过该例子学习发送http请求、解析json格式、使用代码生成提高开发效率等

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

func query(word string) {
	client := &http.Client{}
	request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	query(word)
}

```

## proxy

通过该例子，学习socks5协议工作原理，包含四个阶段：

1. 握手阶段：浏览器向socks5代理发请求，服务器选一个认证方式返回给浏览器
2. 认证阶段：开始认证流程，不概述
3. 请求阶段：认证通过后浏览器向socks5服务器发请求，代理服务器收到响应后和服务器建立连接，然后一个响应
4. replay阶段：浏览器发送请求，代理服务器接收到请求转发给服务器，接收到响应转发给浏览器
![在这里插入图片描述](https://img-blog.csdnimg.cn/aecf2440555c4cbbbb82612e779fbcb7.png)


```go
package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const socks5Ver = 0x05
const cmdBind = 0x01
const atypIPV4 = 0x01
const atypeHOST = 0x03
const atypeIPV6 = 0x04

func main() {
	server, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		panic(err)
	}
	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("Accept failed %v", err)
			continue
		}
		go process(client)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	err := auth(reader, conn)
	if err != nil {
		log.Printf("client %v auth failed:%v", conn.RemoteAddr(), err)
		return
	}
	err = connect(reader, conn)
	if err != nil {
		log.Printf("client %v auth failed:%v", conn.RemoteAddr(), err)
		return
	}
}

func auth(reader *bufio.Reader, conn net.Conn) (err error) {
	// +----+----------+----------+
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	// VER: 协议版本，socks5为0x05
	// NMETHODS: 支持认证的方法数量
	// METHODS: 对应NMETHODS，NMETHODS的值为多少，METHODS就有多少个字节。RFC预定义了一些值的含义，内容如下:
	// X’00’ NO AUTHENTICATION REQUIRED
	// X’02’ USERNAME/PASSWORD

	ver, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read ver failed:%w", err)
	}
	if ver != socks5Ver {
		return fmt.Errorf("not supported ver:%v", ver)
	}
	methodSize, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read methodSize failed:%w", err)
	}
	method := make([]byte, methodSize)
	_, err = io.ReadFull(reader, method)
	if err != nil {
		return fmt.Errorf("read method failed:%w", err)
	}

	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+
	_, err = conn.Write([]byte{socks5Ver, 0x00})
	if err != nil {
		return fmt.Errorf("write failed:%w", err)
	}
	return nil
}

func connect(reader *bufio.Reader, conn net.Conn) (err error) {
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER 版本号，socks5的值为0x05
	// CMD 0x01表示CONNECT请求
	// RSV 保留字段，值为0x00
	// ATYP 目标地址类型，DST.ADDR的数据对应这个字段的类型。
	//   0x01表示IPv4地址，DST.ADDR为4个字节
	//   0x03表示域名，DST.ADDR是一个可变长度的域名
	// DST.ADDR 一个可变长度的值
	// DST.PORT 目标端口，固定2个字节

	buf := make([]byte, 4)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("read header failed:%w", err)
	}
	ver, cmd, atyp := buf[0], buf[1], buf[3]
	if ver != socks5Ver {
		return fmt.Errorf("not supported ver:%v", ver)
	}
	if cmd != cmdBind {
		return fmt.Errorf("not supported cmd:%v", ver)
	}
	addr := ""
	switch atyp {
	case atypIPV4:
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return fmt.Errorf("read atyp failed:%w", err)
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case atypeHOST:
		hostSize, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("read hostSize failed:%w", err)
		}
		host := make([]byte, hostSize)
		_, err = io.ReadFull(reader, host)
		if err != nil {
			return fmt.Errorf("read host failed:%w", err)
		}
		addr = string(host)
	case atypeIPV6:
		return errors.New("IPv6: no supported yet")
	default:
		return errors.New("invalid atyp")
	}
	_, err = io.ReadFull(reader, buf[:2])
	if err != nil {
		return fmt.Errorf("read port failed:%w", err)
	}
	port := binary.BigEndian.Uint16(buf[:2])

	dest, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr, port))
	if err != nil {
		return fmt.Errorf("dial dst failed:%w", err)
	}
	defer dest.Close()
	log.Println("dial", addr, port)

	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER socks版本，这里为0x05
	// REP Relay field,内容取值如下 X’00’ succeeded
	// RSV 保留字段
	// ATYPE 地址类型
	// BND.ADDR 服务绑定的地址
	// BND.PORT 服务绑定的端口DST.PORT
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_, _ = io.Copy(dest, reader)
		cancel()
	}()
	go func() {
		_, _ = io.Copy(conn, dest)
		cancel()
	}()

	<-ctx.Done()
	return nil
}

```

# 4代码示例

以上语法和实战代码示例，都能再[这里](https://github.com/attackoncs/MyByteCamp)找到