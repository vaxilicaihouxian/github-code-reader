### Repository Summary

#### Purpose
# xxl-job-executor-go
很多公司java与go开发共存，java中有xxl-job做为任务调度引擎，为此也出现了go执行器(客户端)，使用起来比较简单：
# 支持
```	
1.执行器注册
2.耗时任务取消
3.任务注册，像写http.Handler一样方便
4.任务panic处理
5.阻塞策略处理
6.任务完成支持返回执行备注
7.任务超时取消 (单位：秒，0为不限制)
8.失败重试次数(在参数param中，目前由任务自行处理)
9.可自定义日志
10.自定义日志查看handler
11.支持外部路由（可与gin集成）
12.支持自定义中间件
```

# Example
```go
package main

import (
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"github.com/xxl-job/xxl-job-executor-go/example/task"
	"log"
)

func main() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://127.0.0.1/xxl-job-admin"),
		xxl.AccessToken(""),            //请求令牌(默认为空)
		xxl.ExecutorIp("127.0.0.1"),    //可自动获取
		xxl.ExecutorPort("9999"),       //默认9999（非必填）
		xxl.RegistryKey("golang-jobs"), //执行器名称
		xxl.SetLogger(&logger{}),       //自定义日志
	)
	exec.Init()
	exec.Use(customMiddleware)
	//设置日志查看handler
	exec.LogHandler(customLogHandle)
	//注册任务handler
	exec.RegTask("task.test", task.Test)
	exec.RegTask("task.test2", task.Test2)
	exec.RegTask("task.panic", task.Panic)
	log.Fatal(exec.Run())
}

// 自定义日志处理器
func customLogHandle(req *xxl.LogReq) *xxl.LogRes {
	return &xxl.LogRes{Code: xxl.SuccessCode, Msg: "", Content: xxl.LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  "这个是自定义日志handler",
		IsEnd:       true,
	}}
}

// xxl.Logger接口实现
type logger struct{}

func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

// 自定义中间件
func customMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(cxt context.Context, param *xxl.RunReq) string {
		log.Println("I am a middleware start")
		res := tf(cxt, param)
		log.Println("I am a middleware end")
		return res
	}
}

```
# 示例项目
github.com/xxl-job/xxl-job-executor-go/example/
# 与gin框架集成
https://github.com/gin-middleware/xxl-job-executor
# xxl-job-admin配置
### 添加执行器
执行器管理->新增执行器,执行器列表如下：
```
AppName		名称		注册方式	OnLine 		机器地址 		操作
golang-jobs	golang执行器	自动注册 		查看 ( 1 ）   
```
查看->注册节点
```
http://127.0.0.1:9999
```
### 添加任务
任务管理->新增(注意，使用BEAN模式，JobHandler与RegTask名称一致)
```
1	测试panic	BEAN：task.panic	* 0 * * * ?	admin	STOP	
2	测试耗时任务	BEAN：task.test2	* * * * * ?	admin	STOP	
3	测试golang	BEAN：task.test		* * * * * ?	admin	STOP
```



#### Structure
```
.
├── .github
│   └── FUNDING.yml
├── .gitignore
├── LICENSE
├── README.md
├── constants.go
├── dto.go
├── example
│   ├── main.go
│   └── task
│       ├── panic.go
│       ├── test.go
│       └── test2.go
├── executor.go
├── go.mod
├── log.go
├── log_handler.go
├── middleware.go
├── optinos.go
├── task.go
├── task_list.go
└── util.go
```

#### Structure Detail
- .github
  - FUNDING.yml
- .gitignore
- LICENSE
- README.md
- constants.go
  Summarized code for constants.go

这个代码文件是一个Go语言的包，名为`xxl`。以下是对该文件内容的详细总结和讲解：

### 文件路径和包声明
```go
package xxl
```
- **包声明**：`package xxl` 表示这个文件属于 `xxl` 包。在Go语言中，包是代码的基本组织单位，同一个包内的代码可以互相引用。

### 常量定义
```go
// 响应码
const (
	SuccessCode = 200
	FailureCode = 500
)
```
- **常量定义**：使用 `const` 关键字定义了两个常量 `SuccessCode` 和 `FailureCode`。
  - `SuccessCode`：值为 `200`，通常用于表示HTTP请求成功。
  - `FailureCode`：值为 `500`，通常用于表示HTTP请求失败，服务器内部错误。

### 详细讲解
1. **包声明**：
   - `package xxl`：这行代码声明了当前文件属于 `xxl` 包。在Go语言中，包名通常与文件夹名一致，方便管理和引用。

2. **常量定义**：
   - `const`：Go语言中用于定义常量的关键字。常量一旦定义，其值在运行时不能被修改。
   - `SuccessCode` 和 `FailureCode`：这两个常量分别表示HTTP响应的成功和失败状态码。
     - `SuccessCode = 200`：HTTP状态码200表示请求成功。
     - `FailureCode = 500`：HTTP状态码500表示服务器内部错误。

### 代码示例
```go
package main

import (
	"fmt"
	"xxl"
)

func main() {
	fmt.Println("Success Code:", xxl.SuccessCode)
	fmt.Println("Failure Code:", xxl.FailureCode)
}
```
- **示例说明**：
  - 这个示例展示了如何在另一个包中使用 `xxl` 包中定义的常量。
  - `import "xxl"`：导入 `xxl` 包。
  - `fmt.Println("Success Code:", xxl.SuccessCode)`：打印 `SuccessCode` 常量的值。
  - `fmt.Println("Failure Code:", xxl.FailureCode)`：打印 `FailureCode` 常量的值。

### 总结
这个代码文件定义了一个名为 `xxl` 的包，并在其中定义了两个常量 `SuccessCode` 和 `FailureCode`，分别表示HTTP响应的成功和失败状态码。通过这种方式，可以方便地在其他代码中引用这些常量，提高代码的可读性和可维护性。
- dto.go
  Summarized code for dto.go

这个代码文件定义了一个用于任务调度和管理的结构体和常量，主要用于与任务调度系统进行交互。以下是对代码的详细讲解：

### 通用响应
```go
type res struct {
	Code int64       `json:"code"` // 200 表示正常、其他失败
	Msg  interface{} `json:"msg"`  // 错误提示消息
}
```
- `res` 结构体用于表示通用的响应信息。
- `Code` 字段表示响应状态码，200 表示正常，其他值表示失败。
- `Msg` 字段用于存储错误提示消息。

### 上行参数
#### Registry 注册参数
```go
type Registry struct {
	RegistryGroup string `json:"registryGroup"`
	RegistryKey   string `json:"registryKey"`
	RegistryValue string `json:"registryValue"`
}
```
- `Registry` 结构体用于注册任务或执行器。
- `RegistryGroup` 字段表示注册组。
- `RegistryKey` 字段表示注册键。
- `RegistryValue` 字段表示注册值。

#### 回调任务结果
```go
type call []*callElement

type callElement struct {
	LogID         int64          `json:"logId"`
	LogDateTim    int64          `json:"logDateTim"`
	ExecuteResult *ExecuteResult `json:"executeResult"`
	HandleCode    int            `json:"handleCode"` // 200表示正常,500表示失败
	HandleMsg     string         `json:"handleMsg"`
}

type ExecuteResult struct {
	Code int64       `json:"code"`
	Msg  interface{} `json:"msg"`
}
```
- `call` 是一个 `callElement` 结构体的切片，用于存储多个回调任务结果。
- `callElement` 结构体包含任务执行的日志信息和执行结果。
- `LogID` 和 `LogDateTim` 字段分别表示日志ID和日志时间。
- `ExecuteResult` 字段表示任务执行结果，包含状态码和消息。
- `HandleCode` 和 `HandleMsg` 字段用于表示回调处理的状态码和消息。

### 下行参数
#### 阻塞处理策略
```go
const (
	serialExecution = "SERIAL_EXECUTION" // 单机串行
	discardLater    = "DISCARD_LATER"    // 丢弃后续调度
	coverEarly      = "COVER_EARLY"      // 覆盖之前调度
)
```
- 定义了三种阻塞处理策略：
  - `serialExecution`：单机串行执行。
  - `discardLater`：丢弃后续调度。
  - `coverEarly`：覆盖之前调度。

#### 触发任务请求参数
```go
type RunReq struct {
	JobID                 int64  `json:"jobId"`
	ExecutorHandler       string `json:"executorHandler"`
	ExecutorParams        string `json:"executorParams"`
	ExecutorBlockStrategy string `json:"executorBlockStrategy"`
	ExecutorTimeout       int64  `json:"executorTimeout"`
	LogID                 int64  `json:"logId"`
	LogDateTime           int64  `json:"logDateTime"`
	GlueType              string `json:"glueType"`
	GlueSource            string `json:"glueSource"`
	GlueUpdatetime        int64  `json:"glueUpdatetime"`
	BroadcastIndex        int64  `json:"broadcastIndex"`
	BroadcastTotal        int64  `json:"broadcastTotal"`
}
```
- `RunReq` 结构体用于触发任务请求。
- `JobID` 字段表示任务ID。
- `ExecutorHandler` 字段表示任务标识。
- `ExecutorParams` 字段表示任务参数。
- `ExecutorBlockStrategy` 字段表示任务阻塞策略。
- `ExecutorTimeout` 字段表示任务超时时间。
- `LogID` 和 `LogDateTime` 字段表示日志ID和日志时间。
- `GlueType` 和 `GlueSource` 字段表示任务模式和脚本代码。
- `GlueUpdatetime` 字段表示脚本更新时间。
- `BroadcastIndex` 和 `BroadcastTotal` 字段表示分片参数。

#### 终止任务请求参数
```go
type killReq struct {
	JobID int64 `json:"jobId"`
}
```
- `killReq` 结构体用于终止任务请求。
- `JobID` 字段表示任务ID。

#### 忙碌检测请求参数
```go
type idleBeatReq struct {
	JobID int64 `json:"jobId"`
}
```
- `idleBeatReq` 结构体用于忙碌检测请求。
- `JobID` 字段表示任务ID。

#### 日志请求
```go
type LogReq struct {
	LogDateTim  int64 `json:"logDateTim"`
	LogID       int64 `json:"logId"`
	FromLineNum int   `json:"fromLineNum"`
}
```
- `LogReq` 结构体用于日志请求。
- `LogDateTim` 和 `LogID` 字段表示日志时间和日志ID。
- `FromLineNum` 字段表示日志开始行号。

#### 日志响应
```go
type LogRes struct {
	Code    int64         `json:"code"`
	Msg     string        `json:"msg"`
	Content LogResContent `json:"content"`
}

type LogResContent struct {
	FromLineNum int    `json:"fromLineNum"`
	ToLineNum   int    `json:"toLineNum"`
	LogContent  string `json:"logContent"`
	IsEnd       bool   `json:"isEnd"`
}
```
- `LogRes` 结构体用于日志响应。
- `Code` 字段表示响应状态码。
- `Msg` 字段表示错误提示消息。
- `Content` 字段包含日志响应内容。
- `LogResContent` 结构体包含日志的开始行号、结束行号、日志内容和是否加载完的标志。

通过这些结构体和常量的定义，代码文件提供了一套完整的任务调度和管理的数据结构，便于与任务调度系统进行交互。
- example
  - main.go
    Summarized code for main.go

这个Go代码文件实现了一个XXL-JOB执行器的示例，使用了`github.com/xxl-job/xxl-job-executor-go`库。以下是对代码的详细讲解：

### 1. 包声明和导入
```go
package main

import (
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"github.com/xxl-job/xxl-job-executor-go/example/task"
	"log"
)
```
- `package main`：声明这是一个可执行的Go程序。
- `import`：导入所需的包，包括标准库和第三方库。

### 2. 主函数
```go
func main() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://127.0.0.1/xxl-job-admin"),
		xxl.AccessToken(""),            //请求令牌(默认为空)
		xxl.ExecutorIp("127.0.0.1"),    //可自动获取
		xxl.ExecutorPort("9999"),       //默认9999（非必填）
		xxl.RegistryKey("golang-jobs"), //执行器名称
		xxl.SetLogger(&logger{}),       //自定义日志
	)
	exec.Init()
	exec.Use(customMiddleware)
	//设置日志查看handler
	exec.LogHandler(customLogHandle)
	//注册任务handler
	exec.RegTask("task.test", task.Test)
	exec.RegTask("task.test2", task.Test2)
	exec.RegTask("task.panic", task.Panic)
	log.Fatal(exec.Run())
}
```
- `xxl.NewExecutor`：创建一个新的XXL-JOB执行器实例，并配置相关参数。
  - `ServerAddr`：指定XXL-JOB管理端的地址。
  - `AccessToken`：请求令牌，默认为空。
  - `ExecutorIp`：执行器IP地址，这里设置为`127.0.0.1`。
  - `ExecutorPort`：执行器端口，默认为`9999`。
  - `RegistryKey`：执行器名称，这里设置为`golang-jobs`。
  - `SetLogger`：设置自定义日志处理器。
- `exec.Init()`：初始化执行器。
- `exec.Use(customMiddleware)`：使用自定义中间件。
- `exec.LogHandler(customLogHandle)`：设置自定义日志处理函数。
- `exec.RegTask`：注册任务处理函数，包括`task.test`、`task.test2`和`task.panic`。
- `log.Fatal(exec.Run())`：启动执行器，并在出现错误时记录日志。

### 3. 自定义日志处理器
```go
func customLogHandle(req *xxl.LogReq) *xxl.LogRes {
	return &xxl.LogRes{Code: xxl.SuccessCode, Msg: "", Content: xxl.LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  "这个是自定义日志handler",
		IsEnd:       true,
	}}
}
```
- `customLogHandle`：自定义日志处理函数，接收一个`LogReq`请求，返回一个`LogRes`响应。
  - `Code`：响应码，这里设置为成功码。
  - `Msg`：消息内容，为空。
  - `Content`：日志内容，包括起始行号、结束行号、日志内容和是否结束标志。

### 4. 自定义日志接口实现
```go
type logger struct{}

func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}
```
- `logger`：自定义日志结构体，实现了`xxl.Logger`接口。
  - `Info`：记录信息日志。
  - `Error`：记录错误日志。

### 5. 自定义中间件
```go
func customMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(cxt context.Context, param *xxl.RunReq) string {
		log.Println("I am a middleware start")
		res := tf(cxt, param)
		log.Println("I am a middleware end")
		return res
	}
}
```
- `customMiddleware`：自定义中间件函数，接收一个任务处理函数`tf`，返回一个新的任务处理函数。
  - 在任务处理函数执行前后分别记录日志。

### 总结
这个代码文件实现了一个XXL-JOB执行器的示例，包括配置执行器参数、注册任务处理函数、设置自定义日志处理器和中间件。通过阅读这个代码，读者可以学习如何使用`github.com/xxl-job/xxl-job-executor-go`库来创建和管理XXL-JOB执行器。
  - task
    - panic.go
      Summarized code for panic.go

这个代码文件是一个简单的Go语言包，主要用于在XXL-JOB执行器中测试panic情况。下面是对代码的详细讲解：

### 包声明
```go
package task
```
这行代码声明了当前文件属于`task`包。

### 导入依赖
```go
import (
	"context"
	xxl "github.com/xxl-job/xxl-job-executor-go"
)
```
这行代码导入了两个包：
1. `context`：用于处理上下文，通常用于控制goroutine的生命周期和传递请求范围的数据。
2. `github.com/xxl-job/xxl-job-executor-go`：这是XXL-JOB执行器的Go语言客户端库，用于与XXL-JOB调度中心进行交互。

### 函数定义
```go
func Panic(cxt context.Context, param *xxl.RunReq) (msg string) {
	panic("test panic")
	return
}
```
这个函数名为`Panic`，它有两个参数和一个返回值：
1. `cxt context.Context`：上下文参数，用于传递上下文信息。
2. `param *xxl.RunReq`：一个指向`xxl.RunReq`结构体的指针，这个结构体包含了任务执行的请求参数。
3. `msg string`：返回值，表示函数的返回消息。

函数体内部只有一行代码：
```go
panic("test panic")
```
这行代码会触发一个panic，抛出一个字符串消息`"test panic"`。在Go语言中，panic会导致程序立即停止执行，并开始回溯调用栈，直到被recover捕获或者程序终止。

### 返回值
```go
return
```
这行代码表示函数返回，但由于函数内部触发了panic，所以这里的返回值实际上不会被使用。

### 总结
这个代码文件的主要功能是定义了一个名为`Panic`的函数，用于在XXL-JOB执行器中测试panic情况。通过调用这个函数，可以触发一个panic，从而测试执行器在遇到panic时的处理机制。这个函数简单地抛出一个字符串消息`"test panic"`，没有实际的业务逻辑。
    - test.go
      Summarized code for test.go

这个代码文件是一个简单的Go语言程序，用于在XXL-JOB任务调度系统中执行一个测试任务。下面是对代码的详细讲解：

### 包声明和导入
```go
package task

import (
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
)
```
- `package task`：声明了该文件属于`task`包，通常用于组织相关的任务处理函数。
- `import`：导入了以下包：
  - `context`：用于处理上下文，通常用于控制goroutine的生命周期和传递请求范围的数据。
  - `fmt`：用于格式化输入和输出。
  - `github.com/xxl-job/xxl-job-executor-go`：XXL-JOB的Go语言执行器库，用于与XXL-JOB调度中心进行交互。

### 函数定义
```go
func Test(cxt context.Context, param *xxl.RunReq) (msg string) {
	fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + " log_id:" + xxl.Int64ToStr(param.LogID))
	return "test done"
}
```
- `func Test(cxt context.Context, param *xxl.RunReq) (msg string)`：定义了一个名为`Test`的函数，该函数接受两个参数：
  - `cxt context.Context`：上下文对象，用于控制函数的生命周期和传递数据。
  - `param *xxl.RunReq`：指向`xxl.RunReq`结构体的指针，包含了任务执行所需的参数。
- `(msg string)`：声明了函数的返回值类型为`string`，并命名为`msg`。

#### 函数体
- `fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + " log_id:" + xxl.Int64ToStr(param.LogID))`：
  - 使用`fmt.Println`函数打印一条日志信息，包含了任务处理器的名称（`param.ExecutorHandler`）、任务参数（`param.ExecutorParams`）和日志ID（`param.LogID`）。
  - `xxl.Int64ToStr(param.LogID)`：将`param.LogID`从`int64`类型转换为字符串类型，以便于拼接和打印。
- `return "test done"`：返回字符串`"test done"`，表示任务执行完成。

### 总结
这个代码文件定义了一个简单的任务处理函数`Test`，用于在XXL-JOB任务调度系统中执行一个测试任务。函数接受上下文和任务参数，打印任务相关信息，并返回一个表示任务完成的字符串。这个函数可以被XXL-JOB调度中心调用，用于执行具体的任务逻辑。
    - test2.go
      Summarized code for test2.go

这个代码文件定义了一个名为 `Test2` 的函数，该函数用于在 `xxl-job` 任务调度框架中执行一个任务。以下是对该代码的详细讲解：

### 包声明和导入
```go
package task

import (
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"time"
)
```
- `package task`：声明该文件属于 `task` 包。
- `import`：导入所需的包：
  - `context`：用于处理上下文，允许任务在需要时被终止。
  - `fmt`：用于格式化输入和输出。
  - `github.com/xxl-job/xxl-job-executor-go`：导入 `xxl-job` 执行器相关的包。
  - `time`：用于处理时间相关的操作，如睡眠。

### 函数定义
```go
func Test2(cxt context.Context, param *xxl.RunReq) (msg string) {
	num := 1
	for {
		select {
		case <-cxt.Done():
			fmt.Println("task" + param.ExecutorHandler + "被手动终止")
			return
		default:
			num++
			time.Sleep(10 * time.Second)
			fmt.Println("test one task"+param.ExecutorHandler+" param："+param.ExecutorParams+"执行行", num)
			if num > 10 {
				fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + "执行完毕！")
				return
			}
		}
	}
}
```
- `func Test2(cxt context.Context, param *xxl.RunReq) (msg string)`：定义了一个名为 `Test2` 的函数，该函数接受两个参数：
  - `cxt context.Context`：上下文对象，用于控制任务的生命周期。
  - `param *xxl.RunReq`：任务运行请求的参数，包含任务处理器的名称和参数。
  - `(msg string)`：函数的返回值，类型为字符串。

- `num := 1`：初始化一个计数器 `num`，用于记录任务执行的次数。

- `for { ... }`：无限循环，用于持续执行任务。

- `select { ... }`：多路复用语句，用于监听多个通道。
  - `case <-cxt.Done():`：监听上下文的 `Done` 通道，如果上下文被取消（例如任务被手动终止），则执行该分支：
    - `fmt.Println("task" + param.ExecutorHandler + "被手动终止")`：打印任务被手动终止的信息。
    - `return`：终止函数执行。

  - `default:`：如果没有收到上下文取消的信号，则执行默认分支：
    - `num++`：计数器加一。
    - `time.Sleep(10 * time.Second)`：让任务休眠 10 秒钟，模拟任务执行时间。
    - `fmt.Println("test one task"+param.ExecutorHandler+" param："+param.ExecutorParams+"执行行", num)`：打印任务执行的信息，包括任务处理器的名称、参数和当前执行次数。
    - `if num > 10 { ... }`：如果计数器 `num` 大于 10，则认为任务执行完毕：
      - `fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + "执行完毕！")`：打印任务执行完毕的信息。
      - `return`：终止函数执行。

### 总结
该代码定义了一个 `xxl-job` 任务调度框架中的任务函数 `Test2`，该函数会每隔 10 秒钟执行一次，并记录执行次数。如果任务被手动终止或执行次数超过 10 次，任务将终止执行并打印相应信息。这个函数展示了如何使用上下文来控制任务的生命周期，以及如何在任务中处理循环和休眠操作。
- executor.go
  Summarized code for executor.go

这个代码文件实现了一个任务执行器（Executor），主要用于管理和执行任务。以下是对代码的详细讲解：

### 包和导入
```go
package xxl

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)
```
- 导入了必要的包，包括标准库和一些第三方库。

### 接口定义
```go
// Executor 执行器
type Executor interface {
	// Init 初始化
	Init(...Option)
	// LogHandler 日志查询
	LogHandler(handler LogHandler)
	// Use 使用中间件
	Use(middlewares ...Middleware)
	// RegTask 注册任务
	RegTask(pattern string, task TaskFunc)
	// RunTask 运行任务
	RunTask(writer http.ResponseWriter, request *http.Request)
	// KillTask 杀死任务
	KillTask(writer http.ResponseWriter, request *http.Request)
	// TaskLog 任务日志
	TaskLog(writer http.ResponseWriter, request *http.Request)
	// Beat 心跳检测
	Beat(writer http.ResponseWriter, request *http.Request)
	// IdleBeat 忙碌检测
	IdleBeat(writer http.ResponseWriter, request *http.Request)
	// Run 运行服务
	Run() error
	// Stop 停止服务
	Stop()
}
```
- 定义了 `Executor` 接口，包含了任务执行器的各种方法。

### 执行器实现
```go
// NewExecutor 创建执行器
func NewExecutor(opts ...Option) Executor {
	return newExecutor(opts...)
}

func newExecutor(opts ...Option) *executor {
	options := newOptions(opts...)
	e := &executor{
		opts: options,
	}
	return e
}
```
- `NewExecutor` 函数用于创建一个新的执行器实例。

### 执行器结构体
```go
type executor struct {
	opts    Options
	address string
	regList *taskList //注册任务列表
	runList *taskList //正在执行任务列表
	mu      sync.RWMutex
	log     Logger

	logHandler  LogHandler   //日志查询handler
	middlewares []Middleware //中间件
}
```
- `executor` 结构体包含了执行器的各种属性和状态。

### 初始化方法
```go
func (e *executor) Init(opts ...Option) {
	for _, o := range opts {
		o(&e.opts)
	}
	e.log = e.opts.l
	e.regList = &taskList{
		data: make(map[string]*Task),
	}
	e.runList = &taskList{
		data: make(map[string]*Task),
	}
	e.address = e.opts.ExecutorIp + ":" + e.opts.ExecutorPort
	go e.registry()
}
```
- `Init` 方法用于初始化执行器，设置选项和启动注册过程。

### 日志处理和中间件
```go
// LogHandler 日志handler
func (e *executor) LogHandler(handler LogHandler) {
	e.logHandler = handler
}

func (e *executor) Use(middlewares ...Middleware) {
	e.middlewares = middlewares
}
```
- `LogHandler` 和 `Use` 方法分别用于设置日志处理函数和中间件。

### 运行服务
```go
func (e *executor) Run() (err error) {
	// 创建路由器
	mux := http.NewServeMux()
	// 设置路由规则
	mux.HandleFunc("/run", e.runTask)
	mux.HandleFunc("/kill", e.killTask)
	mux.HandleFunc("/log", e.taskLog)
	mux.HandleFunc("/beat", e.beat)
	mux.HandleFunc("/idleBeat", e.idleBeat)
	// 创建服务器
	server := &http.Server{
		Addr:         ":" + e.opts.ExecutorPort,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	// 监听端口并提供服务
	e.log.Info("Starting server at " + e.address)
	go server.ListenAndServe()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	e.registryRemove()
	return nil
}
```
- `Run` 方法启动 HTTP 服务器，监听指定端口并处理请求。

### 停止服务
```go
func (e *executor) Stop() {
	e.registryRemove()
}
```
- `Stop` 方法用于停止执行器并移除注册信息。

### 注册任务
```go
// RegTask 注册任务
func (e *executor) RegTask(pattern string, task TaskFunc) {
	var t = &Task{}
	t.fn = e.chain(task)
	e.regList.Set(pattern, t)
	return
}
```
- `RegTask` 方法用于注册任务。

### 运行任务
```go
// 运行一个任务
func (e *executor) runTask(writer http.ResponseWriter, request *http.Request) {
	e.mu.Lock()
	defer e.mu.Unlock()

	req, _ := ioutil.ReadAll(request.Body)
	param := &RunReq{}
	err := json.Unmarshal(req, &param)
	if err != nil {
		_, _ = writer.Write(returnCall(param, FailureCode, "params err"))
		e.log.Error("参数解析错误:" + string(req))
		return
	}
	e.log.Info("任务参数:%v", param)
	if !e.regList.Exists(param.ExecutorHandler) {
		_, _ = writer.Write(returnCall(param, FailureCode, "Task not registered"))
		e.log.Error("任务[" + Int64ToStr(param.JobID) + "]没有注册:" + param.ExecutorHandler)
		return
	}

	//阻塞策略处理
	if e.runList.Exists(Int64ToStr(param.JobID)) {
		if param.ExecutorBlockStrategy == coverEarly { //覆盖之前调度
			oldTask := e.runList.Get(Int64ToStr(param.JobID))
			if oldTask != nil {
				oldTask.Cancel()
				e.runList.Del(Int64ToStr(oldTask.Id))
			}
		} else { //单机串行,丢弃后续调度 都进行阻塞
			_, _ = writer.Write(returnCall(param, FailureCode, "There are tasks running"))
			e.log.Error("任务[" + Int64ToStr(param.JobID) + "]已经在运行了:" + param.ExecutorHandler)
			return
		}
	}

	cxt := context.Background()
	task := e.regList.Get(param.ExecutorHandler)
	if param.ExecutorTimeout > 0 {
		task.Ext, task.Cancel = context.WithTimeout(cxt, time.Duration(param.ExecutorTimeout)*time.Second)
	} else {
		task.Ext, task.Cancel = context.WithCancel(cxt)
	}
	task.Id = param.JobID
	task.Name = param.ExecutorHandler
	task.Param = param
	task.log = e.log

	e.runList.Set(Int64ToStr(task.Id), task)
	go task.Run(func(code int64, msg string) {
		e.callback(task, code, msg)
	})
	e.log.Info("任务[" + Int64ToStr(param.JobID) + "]开始执行:" + param.ExecutorHandler)
	_, _ = writer.Write(returnGeneral())
}
```
- `runTask` 方法用于运行一个任务，处理任务的阻塞策略和超时设置。

### 删除任务
```go
// 删除一个任务
func (e *executor) killTask(writer http.ResponseWriter, request *http.Request) {
	e.mu.Lock()
	defer e.mu.Unlock()
	req, _ := ioutil.ReadAll(request.Body)
	param := &killReq{}
	_ = json.Unmarshal(req, &param)
	if !e.runList.Exists(Int64ToStr(param.JobID)) {
		_, _ = writer.Write(returnKill(param, FailureCode))
		e.log.Error("任务[" + Int64ToStr(param.JobID) + "]没有运行")
		return
	}
	task := e.runList.Get(Int64ToStr(param.JobID))
	task.Cancel()
	e.runList.Del(Int64ToStr(param.JobID))
	_, _ = writer.Write(returnGeneral())
}
```
- `killTask` 方法用于删除一个正在运行的任务。

### 任务日志
```go
// 任务日志
func (e *executor) taskLog(writer http.ResponseWriter, request *http.Request) {
	var res *LogRes
	data, err := ioutil.ReadAll(request.Body)
	req := &LogReq{}
	if err != nil {
		e.log.Error("日志请求失败:" + err.Error())
		reqErrLogHandler(writer, req, err)
		return
	}
	err = json.Unmarshal(data, &req)
	if err != nil {
		e.log.Error("日志请求解析失败:" + err.Error())
		reqErrLogHandler(writer, req, err)
		return
	}
	e.log.Info("日志请求参数:%+v", req)
	if e.logHandler != nil {
		res = e.logHandler(req)
	} else {
		res = defaultLogHandler(req)
	}
	str, _ := json.Marshal(res)
	_, _ = writer.Write(str)
}
```
- `taskLog` 方法用于查询任务日志。

### 心跳检测和忙碌检测
```go
// 心跳检测
func (e *executor) beat(writer http.ResponseWriter, request *http.Request) {
	e.log.Info("心跳检测")
	_, _ = writer.Write(returnGeneral())
}

// 忙碌检测
func (e *executor) idleBeat(writer http.ResponseWriter, request *http.Request) {
	e.mu.Lock()
	defer e.mu.Unlock()
	defer request.Body.Close()
	req, _ := ioutil.ReadAll(request.Body)
	param := &idleBeatReq{}
	err := json.Unmarshal(req, &param)
	if err != nil {
		_, _ = writer.Write(returnIdleBeat(FailureCode))
		e.log.Error("参数解析错误:" + string(req))
		return
	}
	if e.runList.Exists(Int64ToStr(param.JobID)) {
		_, _ = writer.Write(returnIdleBeat(FailureCode))
		e.log.Error("idleBeat任务[" + Int64ToStr(param.JobID) + "]正在运行")
		return
	}
	e.log.Info("忙碌检测任务参数:%v", param)
	_, _ = writer.Write(returnGeneral())
}
```
- `beat` 和 `idleBeat` 方法分别用于心跳检测和忙碌检测。

### 注册和移除注册
```go
// 注册执行器到调度中心
func (e *executor) registry() {
	t := time.NewTimer(time.Second * 0) //初始立即执行
	defer t.Stop()
	req := &Registry{
		RegistryGroup: "EXECUTOR",
		RegistryKey:   e.opts.RegistryKey,
		RegistryValue: "http://" + e.address,
	}
	param, err := json.Marshal(req)
	if err != nil {
		log.Fatal("执行器注册信息解析失败:" + err.Error())
	}
	for {
		<-t.C
		t.Reset(time.Second * time.Duration(20)) //20秒心跳防止过期
		func() {
			result, err := e.post("/api/registry", string(param))
			if err != nil {
				e.log.Error("执行器注册失败1:" + err.Error())
				return
			}
			defer result.Body.Close()
			body, err := ioutil.ReadAll(result.Body)
			if err != nil {
				e.log.Error("执行器注册失败2:" + err.Error())
				return
			}
			res := &res{}
			_ = json.Unmarshal(body, &res)
			if res.Code != SuccessCode {
				e.log.Error("执行器注册失败3:" + string(body))
				return
			}
			e.log.Info("执行器注册成功:" + string(body))
		}()

	}
}

// 执行器注册摘除
func (e *executor) registryRemove() {
	t := time.NewTimer(time.Second * 0) //初始立即执行
	defer t.Stop()
	req := &Registry{
		RegistryGroup: "EXECUTOR",
		RegistryKey:   e.opts.RegistryKey,
		RegistryValue: "http://" + e.address,
	}
	param, err := json.Marshal(req)
	if err != nil {
		e.log.Error("执行器摘除失败:" + err.Error())
		return
	}
	res, err := e.post("/api/registryRemove", string(param))
	if err != nil {
		e.log.Error("执行器摘除失败:" + err.Error())
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	e.log.Info("执行器摘除成功:" + string(body))
}
```
- `registry` 和 `registryRemove` 方法分别用于注册执行器和移除注册信息。

### 回调任务列表
```go
// 回调任务列表
func (e *executor) callback(task *Task, code int64, msg string) {
	e.runList.Del(Int64ToStr(task.Id))
	res, err := e.post("/api/callback", string(returnCall(task.Param, code, msg)))
	if err != nil {
		e.log.Error("callback err : ", err.Error())
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e.log.Error("callback ReadAll err : ", err.Error())
		return
	}
	e.log.Info("任务回调成功:" + string(body))
}
```
- `callback` 方法用于回调任务结果。

### POST 请求
```go
// post
func (e *executor) post(action, body string) (resp *http.Response, err error) {
	request, err := http.NewRequest("POST", e.opts.ServerAddr+action, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("XXL-JOB-ACCESS-TOKEN", e.opts.AccessToken)
	client := http.Client{
		Timeout: e.opts.Timeout,
	}
	return client.Do(request)
}
```
- `post` 方法用于发送 POST 请求。

### 其他方法
```go
// RunTask 运行任务
func (e *executor) RunTask(writer http.ResponseWriter, request *http.Request) {
	e.runTask(writer, request)
}

// KillTask 删除任务
func (e *executor) KillTask(writer http.ResponseWriter, request *http.Request) {
	e.killTask(writer, request)
}

// TaskLog 任务日志
func (e *executor) TaskLog(writer http.ResponseWriter, request *http.Request) {
	e.taskLog(writer, request)
}

// Beat 心跳检测
func (e *executor) Beat(writer http.ResponseWriter, request *http.Request) {
	e.beat(writer, request)
}

// IdleBeat 忙碌检测
func (e *executor) IdleBeat(writer http.ResponseWriter, request *http.Request) {
	e.idleBeat(writer, request)
}
```
- 这些方法是对内部方法的简单封装，便于外部调用。

### 总结
这个代码文件实现了一个功能丰富的任务执行器，包括任务的注册、运行、删除、日志查询、心跳检测和忙碌检测等功能。通过 HTTP 接口与外部系统进行交互，支持任务的并发执行和阻塞策略处理。
- go.mod
- log.go
  Summarized code for log.go

这个代码文件定义了一个简单的日志系统，包括应用日志和系统日志的功能。以下是对代码的详细讲解：

### 包声明
```go
package xxl
```
这行代码声明了代码文件所属的包名，这里是 `xxl`。

### 导入依赖
```go
import (
	"fmt"
	"log"
)
```
这行代码导入了两个标准库包：
- `fmt`：用于格式化输入和输出。
- `log`：用于记录日志信息。

### 应用日志类型定义
```go
// LogFunc 应用日志
type LogFunc func(req LogReq, res *LogRes) []byte
```
这行代码定义了一个类型 `LogFunc`，它是一个函数类型，接受两个参数 `req` 和 `res`，并返回一个字节数组。具体参数类型 `LogReq` 和 `LogRes` 没有在代码中定义，可能是其他地方定义的类型。

### 系统日志接口定义
```go
// Logger 系统日志
type Logger interface {
	Info(format string, a ...interface{})
	Error(format string, a ...interface{})
}
```
这行代码定义了一个接口 `Logger`，包含两个方法：
- `Info`：用于记录信息级别的日志。
- `Error`：用于记录错误级别的日志。

### 系统日志结构体定义
```go
type logger struct {
}
```
这行代码定义了一个结构体 `logger`，用于实现 `Logger` 接口。

### 实现 `Info` 方法
```go
func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}
```
这行代码实现了 `Logger` 接口中的 `Info` 方法。它使用 `fmt.Sprintf` 格式化字符串，然后使用 `fmt.Println` 打印到标准输出。

### 实现 `Error` 方法
```go
func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf(format, a...))
}
```
这行代码实现了 `Logger` 接口中的 `Error` 方法。它使用 `fmt.Sprintf` 格式化字符串，然后使用 `log.Println` 打印到标准错误输出。

### 总结
这个代码文件定义了一个简单的日志系统，包括应用日志和系统日志的功能。应用日志通过 `LogFunc` 类型定义，系统日志通过 `Logger` 接口和 `logger` 结构体实现。`logger` 结构体实现了 `Logger` 接口中的 `Info` 和 `Error` 方法，分别用于记录信息和错误级别的日志。
- log_handler.go
  Summarized code for log_handler.go

这个代码文件定义了一个用于日志查询的HTTP处理程序，主要用于在xxl-job-admin后台显示日志信息。以下是对代码的详细讲解：

### 包声明和导入
```go
package xxl

import (
	"encoding/json"
	"net/http"
)
```
- `package xxl`：声明了包名为`xxl`。
- `import`：导入了两个标准库包，`encoding/json`用于JSON编码和解码，`net/http`用于HTTP请求和响应处理。

### 类型定义
```go
type LogHandler func(req *LogReq) *LogRes
```
- `LogHandler`：定义了一个函数类型，该函数接收一个`LogReq`类型的指针作为参数，并返回一个`LogRes`类型的指针。

### 默认日志处理函数
```go
func defaultLogHandler(req *LogReq) *LogRes {
	return &LogRes{Code: SuccessCode, Msg: "", Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  "这是日志默认返回，说明没有设置LogHandler",
		IsEnd:       true,
	}}
}
```
- `defaultLogHandler`：这是一个默认的日志处理函数，当没有设置自定义的`LogHandler`时，会使用这个函数。
  - 返回一个`LogRes`结构体，其中`Code`字段为`SuccessCode`，表示成功。
  - `Msg`字段为空字符串。
  - `Content`字段是一个`LogResContent`结构体，包含以下字段：
    - `FromLineNum`：从请求中获取的起始行号。
    - `ToLineNum`：固定为2，表示结束行号。
    - `LogContent`：默认的日志内容字符串。
    - `IsEnd`：固定为`true`，表示日志结束。

### 请求错误处理函数
```go
func reqErrLogHandler(w http.ResponseWriter, req *LogReq, err error) {
	res := &LogRes{Code: FailureCode, Msg: err.Error(), Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   0,
		LogContent:  err.Error(),
		IsEnd:       true,
	}}
	str, _ := json.Marshal(res)
	_, _ = w.Write(str)
}
```
- `reqErrLogHandler`：这是一个处理请求错误的函数。
  - 接收一个`http.ResponseWriter`、一个`LogReq`类型的指针和一个`error`类型的参数。
  - 创建一个`LogRes`结构体，其中`Code`字段为`FailureCode`，表示失败。
  - `Msg`字段为错误信息。
  - `Content`字段是一个`LogResContent`结构体，包含以下字段：
    - `FromLineNum`：从请求中获取的起始行号。
    - `ToLineNum`：固定为0，表示没有有效的结束行号。
    - `LogContent`：错误信息字符串。
    - `IsEnd`：固定为`true`，表示日志结束。
  - 将`LogRes`结构体序列化为JSON字符串，并通过`http.ResponseWriter`写回到客户端。

### 总结
这个代码文件主要定义了两个日志处理函数：
1. `defaultLogHandler`：用于在没有自定义日志处理函数时提供默认的日志返回。
2. `reqErrLogHandler`：用于处理请求错误，并将错误信息返回给客户端。

通过这些函数，可以实现日志查询功能，并将日志信息显示在xxl-job-admin后台。
- middleware.go
  Summarized code for middleware.go

这个代码文件定义了一个用于处理任务的中间件系统。以下是对代码的详细讲解：

### 包声明
```go
package xxl
```
这行代码声明了代码文件所属的包名，这里是 `xxl`。

### 中间件类型定义
```go
// Middleware 中间件构造函数
type Middleware func(TaskFunc) TaskFunc
```
这行代码定义了一个名为 `Middleware` 的类型，它是一个函数类型。这个函数接收一个 `TaskFunc` 类型的参数，并返回一个 `TaskFunc` 类型的值。这里的 `TaskFunc` 是一个任务处理函数，具体定义可能在其他地方。

### 中间件链式调用
```go
func (e *executor) chain(next TaskFunc) TaskFunc {
	for i := range e.middlewares {
		next = e.middlewares[len(e.middlewares)-1-i](next)
	}
	return next
}
```
这行代码定义了一个名为 `chain` 的方法，它属于 `executor` 类型的指针接收者。这个方法的作用是将多个中间件串联起来，形成一个中间件链。

#### 方法细节
1. **接收者**：`e *executor` 表示这个方法属于 `executor` 类型的指针接收者。
2. **参数**：`next TaskFunc` 表示这个方法接收一个 `TaskFunc` 类型的参数，这个参数是一个任务处理函数。
3. **返回值**：`TaskFunc` 表示这个方法返回一个 `TaskFunc` 类型的值，即经过中间件链处理后的任务处理函数。

#### 实现细节
- **循环遍历中间件**：
  ```go
  for i := range e.middlewares {
  ```
  这行代码遍历 `executor` 实例中的 `middlewares` 切片。`middlewares` 切片存储了多个中间件函数。

- **倒序应用中间件**：
  ```go
  next = e.middlewares[len(e.middlewares)-1-i](next)
  ```
  这行代码将中间件按倒序应用到 `next` 任务处理函数上。`len(e.middlewares)-1-i` 表示从最后一个中间件开始应用，依次向前。

- **返回处理后的任务函数**：
  ```go
  return next
  ```
  这行代码返回经过所有中间件处理后的任务处理函数。

### 总结
这个代码文件实现了一个中间件系统，允许将多个中间件函数串联起来，按倒序应用到任务处理函数上。通过这种方式，可以灵活地对任务处理进行扩展和定制。
- optinos.go
  Summarized code for optinos.go

这个代码文件定义了一个用于配置和管理执行器（Executor）的选项（Options）结构体及其相关功能。以下是对代码的详细讲解：

### 包和导入
```go
package xxl

import (
	"github.com/go-basic/ipv4"
	"time"
)
```
- `package xxl`：定义了包名为 `xxl`。
- `import`：导入了两个外部包：
  - `github.com/go-basic/ipv4`：用于获取本地IP地址。
  - `time`：用于处理时间相关的操作。

### 选项结构体
```go
type Options struct {
	ServerAddr   string        `json:"server_addr"`   //调度中心地址
	AccessToken  string        `json:"access_token"`  //请求令牌
	Timeout      time.Duration `json:"timeout"`       //接口超时时间
	ExecutorIp   string        `json:"executor_ip"`   //本地(执行器)IP(可自行获取)
	ExecutorPort string        `json:"executor_port"` //本地(执行器)端口
	RegistryKey  string        `json:"registry_key"`  //执行器名称
	LogDir       string        `json:"log_dir"`       //日志目录

	l Logger //日志处理
}
```
- `Options` 结构体定义了执行器的各种配置选项：
  - `ServerAddr`：调度中心地址。
  - `AccessToken`：请求令牌。
  - `Timeout`：接口超时时间。
  - `ExecutorIp`：本地执行器IP地址，默认会自动获取。
  - `ExecutorPort`：本地执行器端口，默认值为 `9999`。
  - `RegistryKey`：执行器名称，默认值为 `golang-jobs`。
  - `LogDir`：日志目录。
  - `l`：日志处理器，类型为 `Logger`。

### 默认选项
```go
var (
	DefaultExecutorPort = "9999"
	DefaultRegistryKey  = "golang-jobs"
)
```
- 定义了两个默认值：
  - `DefaultExecutorPort`：默认执行器端口为 `9999`。
  - `DefaultRegistryKey`：默认执行器名称为 `golang-jobs`。

### 创建选项
```go
func newOptions(opts ...Option) Options {
	opt := Options{
		ExecutorIp:   ipv4.LocalIP(),
		ExecutorPort: DefaultExecutorPort,
		RegistryKey:  DefaultRegistryKey,
	}

	for _, o := range opts {
		o(&opt)
	}

	if opt.l == nil {
		opt.l = &logger{}
	}

	return opt
}
```
- `newOptions` 函数用于创建 `Options` 实例，并应用传入的选项：
  - 初始化 `Options` 结构体，设置默认值。
  - 遍历传入的 `Option` 函数，应用每个选项。
  - 如果日志处理器 `l` 为空，则设置默认的日志处理器。

### 选项函数类型
```go
type Option func(o *Options)
```
- `Option` 是一个函数类型，用于修改 `Options` 结构体。

### 设置选项函数
```go
// ServerAddr 设置调度中心地址
func ServerAddr(addr string) Option {
	return func(o *Options) {
		o.ServerAddr = addr
	}
}

// AccessToken 请求令牌
func AccessToken(token string) Option {
	return func(o *Options) {
		o.AccessToken = token
	}
}

// ExecutorIp 设置执行器IP
func ExecutorIp(ip string) Option {
	return func(o *Options) {
		o.ExecutorIp = ip
	}
}

// ExecutorPort 设置执行器端口
func ExecutorPort(port string) Option {
	return func(o *Options) {
		o.ExecutorPort = port
	}
}

// RegistryKey 设置执行器标识
func RegistryKey(registryKey string) Option {
	return func(o *Options) {
		o.RegistryKey = registryKey
	}
}

// SetLogger 设置日志处理器
func SetLogger(l Logger) Option {
	return func(o *Options) {
		o.l = l
	}
}
```
- 这些函数用于设置 `Options` 结构体的各个字段：
  - `ServerAddr`：设置调度中心地址。
  - `AccessToken`：设置请求令牌。
  - `ExecutorIp`：设置执行器IP地址。
  - `ExecutorPort`：设置执行器端口。
  - `RegistryKey`：设置执行器标识。
  - `SetLogger`：设置日志处理器。

### 总结
这个代码文件定义了一个灵活的配置系统，用于管理执行器的各种选项。通过使用函数式选项模式（Functional Options Pattern），可以方便地设置和修改执行器的配置，同时保持代码的可扩展性和可读性。
- task.go
  Summarized code for task.go

这个代码文件定义了一个任务执行框架的基本结构和功能。以下是对代码的详细讲解：

### 包和导入
```go
package xxl

import (
	"context"
	"fmt"
	"runtime/debug"
)
```
- `package xxl`：定义了包名为 `xxl`。
- `import`：导入了 `context`、`fmt` 和 `runtime/debug` 包，用于上下文管理、格式化输出和堆栈跟踪。

### 类型定义
#### TaskFunc
```go
// TaskFunc 任务执行函数
type TaskFunc func(cxt context.Context, param *RunReq) string
```
- `TaskFunc` 是一个函数类型，接受一个 `context.Context` 和一个 `RunReq` 类型的指针作为参数，并返回一个字符串。

#### Task
```go
// Task 任务
type Task struct {
	Id        int64
	Name      string
	Ext       context.Context
	Param     *RunReq
	fn        TaskFunc
	Cancel    context.CancelFunc
	StartTime int64
	EndTime   int64
	//日志
	log Logger
}
```
- `Task` 结构体定义了一个任务的各个属性：
  - `Id`：任务的唯一标识。
  - `Name`：任务名称。
  - `Ext`：任务的上下文。
  - `Param`：任务的参数，类型为 `RunReq`。
  - `fn`：任务执行函数，类型为 `TaskFunc`。
  - `Cancel`：用于取消任务的函数。
  - `StartTime` 和 `EndTime`：任务开始和结束的时间戳。
  - `log`：用于记录日志的接口。

### 方法
#### Run
```go
// Run 运行任务
func (t *Task) Run(callback func(code int64, msg string)) {
	defer func(cancel func()) {
		if err := recover(); err != nil {
			t.log.Info(t.Info()+" panic: %v", err)
			debug.PrintStack() //堆栈跟踪
			callback(FailureCode, fmt.Sprintf("task panic:%v", err))
			cancel()
		}
	}(t.Cancel)
	msg := t.fn(t.Ext, t.Param)
	callback(SuccessCode, msg)
	return
}
```
- `Run` 方法用于运行任务：
  - 使用 `defer` 捕获任务执行过程中的 panic，并进行处理：
    - 记录 panic 信息和堆栈跟踪。
    - 调用 `callback` 函数返回失败代码和错误信息。
    - 调用 `cancel` 函数取消任务。
  - 如果没有 panic，则调用任务执行函数 `t.fn`，并将结果通过 `callback` 函数返回成功代码和消息。

#### Info
```go
// Info 任务信息
func (t *Task) Info() string {
	return fmt.Sprintf("任务ID[%d]任务名称[%s]参数:%s", t.Id, t.Name, t.Param.ExecutorParams)
}
```
- `Info` 方法返回任务的基本信息，包括任务 ID、任务名称和参数。

### 总结
这个代码文件定义了一个简单的任务执行框架，包括任务的定义、执行和日志记录。通过 `Task` 结构体和相关方法，可以方便地管理和执行任务，并在任务执行过程中捕获和处理 panic。这个框架可以进一步扩展，例如添加更多的任务参数、日志记录级别等。
- task_list.go
  Summarized code for task_list.go

这个代码文件定义了一个用于管理任务列表的数据结构和相关操作方法。以下是对代码的详细讲解：

### 包和导入
```go
package xxl

import "sync"
```
- `package xxl`：定义了代码所属的包名，这里假设包名为`xxl`。
- `import "sync"`：导入了Go语言标准库中的`sync`包，用于提供并发操作时的同步机制。

### 数据结构定义
```go
// 任务列表 [JobID]执行函数,并行执行时[+LogID]
type taskList struct {
	mu   sync.RWMutex
	data map[string]*Task
}
```
- `type taskList struct`：定义了一个名为`taskList`的结构体，用于存储任务列表。
  - `mu sync.RWMutex`：一个读写锁，用于在并发环境下保护对`data`的访问。
  - `data map[string]*Task`：一个映射表，键为字符串类型（通常是任务ID），值为指向`Task`结构体的指针。

### 方法定义
#### Set 方法
```go
// Set 设置数据
func (t *taskList) Set(key string, val *Task) {
	t.mu.Lock()
	t.data[key] = val
	t.mu.Unlock()
}
```
- `Set`方法：用于向任务列表中添加或更新任务。
  - `t.mu.Lock()`：获取写锁，确保在同一时间只有一个goroutine可以修改`data`。
  - `t.data[key] = val`：将键值对添加到映射表中。
  - `t.mu.Unlock()`：释放写锁。

#### Get 方法
```go
// Get 获取数据
func (t *taskList) Get(key string) *Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data[key]
}
```
- `Get`方法：用于从任务列表中获取指定键的任务。
  - `t.mu.RLock()`：获取读锁，允许多个goroutine同时读取`data`。
  - `defer t.mu.RUnlock()`：确保在函数返回前释放读锁。
  - `return t.data[key]`：返回指定键对应的任务。

#### GetAll 方法
```go
// GetAll 获取数据
func (t *taskList) GetAll() map[string]*Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data
}
```
- `GetAll`方法：用于获取任务列表中的所有任务。
  - `t.mu.RLock()`：获取读锁。
  - `defer t.mu.RUnlock()`：确保在函数返回前释放读锁。
  - `return t.data`：返回整个任务列表的映射表。

#### Del 方法
```go
// Del 设置数据
func (t *taskList) Del(key string) {
	t.mu.Lock()
	delete(t.data, key)
	t.mu.Unlock()
}
```
- `Del`方法：用于从任务列表中删除指定键的任务。
  - `t.mu.Lock()`：获取写锁。
  - `delete(t.data, key)`：从映射表中删除指定键的任务。
  - `t.mu.Unlock()`：释放写锁。

#### Len 方法
```go
// Len 长度
func (t *taskList) Len() int {
	return len(t.data)
}
```
- `Len`方法：用于获取任务列表的长度（即任务的数量）。
  - `return len(t.data)`：返回映射表的长度。

#### Exists 方法
```go
// Exists Key是否存在
func (t *taskList) Exists(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.data[key]
	return ok
}
```
- `Exists`方法：用于检查指定键的任务是否存在于任务列表中。
  - `t.mu.RLock()`：获取读锁。
  - `defer t.mu.RUnlock()`：确保在函数返回前释放读锁。
  - `_, ok := t.data[key]`：检查指定键是否存在于映射表中。
  - `return ok`：返回检查结果。

### 总结
这个代码文件定义了一个并发安全的任务列表管理器，提供了添加、获取、删除、检查存在性和获取长度的功能。通过使用读写锁，确保了在多goroutine环境下对任务列表的安全访问。
- util.go
  Summarized code for util.go

这个代码文件主要定义了一些用于处理任务回调、任务杀死、任务忙碌状态和通用返回的函数。以下是对每个函数的详细讲解：

### 1. `Int64ToStr` 函数
```go
// Int64ToStr int64 to str
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}
```
- **功能**：将 `int64` 类型的整数转换为字符串。
- **实现细节**：使用 `strconv.FormatInt` 函数，将 `int64` 类型的整数 `i` 转换为十进制字符串。

### 2. `returnCall` 函数
```go
// 执行任务回调
func returnCall(req *RunReq, code int64, msg string) []byte {
	data := call{
		&callElement{
			LogID:      req.LogID,
			LogDateTim: req.LogDateTime,
			ExecuteResult: &ExecuteResult{
				Code: code,
				Msg:  msg,
			},
			HandleCode: int(code),
			HandleMsg:  msg,
		},
	}
	str, _ := json.Marshal(data)
	return str
}
```
- **功能**：生成一个执行任务回调的 JSON 响应。
- **实现细节**：
  - 创建一个 `call` 结构体，包含一个 `callElement` 元素。
  - `callElement` 包含任务的日志 ID、日志时间、执行结果（包含状态码和消息）、处理状态码和处理消息。
  - 使用 `json.Marshal` 将 `data` 结构体转换为 JSON 字符串。

### 3. `returnKill` 函数
```go
// 杀死任务返回
func returnKill(req *killReq, code int64) []byte {
	msg := ""
	if code != SuccessCode {
		msg = "Task kill err"
	}
	data := res{
		Code: code,
		Msg:  msg,
	}
	str, _ := json.Marshal(data)
	return str
}
```
- **功能**：生成一个杀死任务的 JSON 响应。
- **实现细节**：
  - 根据状态码 `code` 判断是否成功，如果不成功则设置消息为 "Task kill err"。
  - 创建一个 `res` 结构体，包含状态码和消息。
  - 使用 `json.Marshal` 将 `data` 结构体转换为 JSON 字符串。

### 4. `returnIdleBeat` 函数
```go
// 忙碌返回
func returnIdleBeat(code int64) []byte {
	msg := ""
	if code != SuccessCode {
		msg = "Task is busy"
	}
	data := res{
		Code: code,
		Msg:  msg,
	}
	str, _ := json.Marshal(data)
	return str
}
```
- **功能**：生成一个任务忙碌状态的 JSON 响应。
- **实现细节**：
  - 根据状态码 `code` 判断是否成功，如果不成功则设置消息为 "Task is busy"。
  - 创建一个 `res` 结构体，包含状态码和消息。
  - 使用 `json.Marshal` 将 `data` 结构体转换为 JSON 字符串。

### 5. `returnGeneral` 函数
```go
// 通用返回
func returnGeneral() []byte {
	data := &res{
		Code: SuccessCode,
		Msg:  "",
	}
	str, _ := json.Marshal(data)
	return str
}
```
- **功能**：生成一个通用的成功 JSON 响应。
- **实现细节**：
  - 创建一个 `res` 结构体，包含成功的状态码和空消息。
  - 使用 `json.Marshal` 将 `data` 结构体转换为 JSON 字符串。

### 总结
这个代码文件主要用于生成不同场景下的 JSON 响应，包括任务回调、任务杀死、任务忙碌状态和通用成功响应。每个函数都通过创建相应的结构体并使用 `json.Marshal` 将其转换为 JSON 字符串来实现功能。
