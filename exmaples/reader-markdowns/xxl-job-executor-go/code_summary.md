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
- .github
  - FUNDING.yml
- .gitignore
- LICENSE
- README.md
- constants.go
  Summarized code for constants.go

这个代码文件定义了一个名为 `xxl` 的包，并包含了一些常量，用于表示响应码。以下是对代码的详细解释：

### 包声明
```go
package xxl
```
这行代码声明了当前文件属于 `xxl` 包。在 Go 语言中，包是代码组织和重用的基本单位。通过将相关功能的代码放在同一个包中，可以方便地管理和调用这些功能。

### 常量定义
```go
const (
	SuccessCode = 200
	FailureCode = 500
)
```
这部分代码定义了两个常量：

1. **`SuccessCode`**:
   - 值为 `200`。
   - 通常用于表示请求成功的情况。在 HTTP 协议中，200 是一个标准的成功响应码。

2. **`FailureCode`**:
   - 值为 `500`。
   - 通常用于表示请求失败的情况。在 HTTP 协议中，500 是一个标准的内部服务器错误响应码。

### 功能和实现细节
- **常量声明**:
  - 使用 `const` 关键字定义常量。常量在声明后不能被修改，这有助于确保代码的安全性和稳定性。
  - 常量名采用驼峰命名法，`SuccessCode` 和 `FailureCode` 都是自描述的名称，便于理解其用途。

- **响应码**:
  - 这两个常量通常用于网络请求的响应处理中，帮助开发者快速判断请求的结果是成功还是失败。
  - 在实际应用中，可以根据这些常量来决定后续的业务逻辑，例如在成功时执行某些操作，在失败时进行错误处理或重试。

### 示例应用
假设有一个处理 HTTP 请求的函数，可以使用这些常量来返回响应码：

```go
package main

import (
	"net/http"
	"xxl" // 假设 xxl 包在同一目录下
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 模拟处理请求
	result := processRequest(r)

	if result {
		w.WriteHeader(xxl.SuccessCode)
	} else {
		w.WriteHeader(xxl.FailureCode)
	}
}

func processRequest(r *http.Request) bool {
	// 模拟请求处理逻辑
	return true // 或 false，根据实际逻辑
}

func main() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}
```

在这个示例中，`handleRequest` 函数根据 `processRequest` 的返回结果来设置 HTTP 响应码。如果请求处理成功，返回 `SuccessCode`（200），否则返回 `FailureCode`（500）。

### 总结
这个代码文件通过定义两个常量 `SuccessCode` 和 `FailureCode`，为处理 HTTP 响应提供了一种简洁且标准化的方式。通过使用这些常量，可以提高代码的可读性和可维护性，同时也便于在不同模块之间统一处理响应码。
- dto.go
  Summarized code for dto.go

该代码文件定义了一系列结构体和常量，主要用于处理任务调度和日志记录的请求和响应。以下是对代码的详细解释：

### 1. 通用响应结构体
```go
type res struct {
	Code int64       `json:"code"` // 200 表示正常、其他失败
	Msg  interface{} `json:"msg"`  // 错误提示消息
}
```
- `res` 结构体用于表示通用的响应信息，包含一个状态码 `Code` 和一个消息 `Msg`。状态码 `200` 表示操作成功，其他值表示失败。

### 2. 注册参数结构体
```go
type Registry struct {
	RegistryGroup string `json:"registryGroup"`
	RegistryKey   string `json:"registryKey"`
	RegistryValue string `json:"registryValue"`
}
```
- `Registry` 结构体用于表示注册参数，包含注册组 `RegistryGroup`、注册键 `RegistryKey` 和注册值 `RegistryValue`。

### 3. 回调任务结果结构体
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
- `call` 是一个 `callElement` 结构体的切片，用于表示多个回调任务结果。
- `callElement` 结构体包含日志ID `LogID`、日志时间 `LogDateTim`、执行结果 `ExecuteResult`、处理代码 `HandleCode` 和处理消息 `HandleMsg`。
- `ExecuteResult` 结构体用于表示任务执行结果，包含状态码 `Code` 和消息 `Msg`。

### 4. 阻塞处理策略常量
```go
const (
	serialExecution = "SERIAL_EXECUTION" // 单机串行
	discardLater    = "DISCARD_LATER"    // 丢弃后续调度
	coverEarly      = "COVER_EARLY"      // 覆盖之前调度
)
```
- 定义了三种阻塞处理策略：单机串行 `SERIAL_EXECUTION`、丢弃后续调度 `DISCARD_LATER` 和覆盖之前调度 `COVER_EARLY`。

### 5. 触发任务请求参数结构体
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
- `RunReq` 结构体用于表示触发任务的请求参数，包含任务ID `JobID`、任务标识 `ExecutorHandler`、任务参数 `ExecutorParams`、任务阻塞策略 `ExecutorBlockStrategy`、任务超时时间 `ExecutorTimeout`、日志ID `LogID`、日志时间 `LogDateTime`、任务模式 `GlueType`、GLUE脚本代码 `GlueSource`、GLUE脚本更新时间 `GlueUpdatetime`、分片参数 `BroadcastIndex` 和总分片 `BroadcastTotal`。

### 6. 终止任务请求参数结构体
```go
type killReq struct {
	JobID int64 `json:"jobId"`
}
```
- `killReq` 结构体用于表示终止任务的请求参数，包含任务ID `JobID`。

### 7. 忙碌检测请求参数结构体
```go
type idleBeatReq struct {
	JobID int64 `json:"jobId"`
}
```
- `idleBeatReq` 结构体用于表示忙碌检测的请求参数，包含任务ID `JobID`。

### 8. 日志请求参数结构体
```go
type LogReq struct {
	LogDateTim  int64 `json:"logDateTim"`
	LogID       int64 `json:"logId"`
	FromLineNum int   `json:"fromLineNum"`
}
```
- `LogReq` 结构体用于表示日志请求参数，包含日志时间 `LogDateTim`、日志ID `LogID` 和日志开始行号 `FromLineNum`。

### 9. 日志响应结构体
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
- `LogRes` 结构体用于表示日志响应，包含状态码 `Code`、消息 `Msg` 和日志响应内容 `Content`。
- `LogResContent` 结构体用于表示日志响应内容，包含日志开始行号 `FromLineNum`、日志结束行号 `ToLineNum`、日志内容 `LogContent` 和日志是否全部加载完 `IsEnd`。

### 总结
该代码文件定义了一系列用于任务调度和日志记录的请求和响应结构体，以及一些常量。这些结构体和常量用于在系统中传递和处理任务相关的数据，确保任务的正确执行和日志的记录。通过这些定义，可以清晰地了解系统中任务调度和日志记录的数据结构和处理逻辑。
- example
  - main.go
    Summarized code for main.go

这个代码文件是一个使用 `xxl-job` 任务调度框架的 Go 语言执行器示例。`xxl-job` 是一个分布式任务调度平台，支持多种编程语言实现任务执行器。以下是对代码文件的详细解析：

### 1. 包和导入
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
- `package main`：定义了主包，表示这是一个可执行程序。
- `import`：导入了所需的包，包括 `context`、`fmt`、`log` 以及 `xxl-job` 相关的包。

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
- `xxl.NewExecutor`：创建一个新的 `xxl-job` 执行器实例，并配置相关参数：
  - `ServerAddr`：指定 `xxl-job` 管理端的地址。
  - `AccessToken`：请求令牌，默认为空。
  - `ExecutorIp`：执行器 IP 地址，这里设置为 `127.0.0.1`。
  - `ExecutorPort`：执行器端口，默认为 `9999`。
  - `RegistryKey`：执行器名称，这里设置为 `golang-jobs`。
  - `SetLogger`：设置自定义日志处理器。
- `exec.Init()`：初始化执行器。
- `exec.Use(customMiddleware)`：使用自定义中间件。
- `exec.LogHandler(customLogHandle)`：设置自定义日志处理函数。
- `exec.RegTask`：注册任务处理函数，分别注册了 `task.test`、`task.test2` 和 `task.panic` 三个任务。
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
- `customLogHandle`：自定义日志处理函数，接收一个日志请求 `LogReq`，返回一个日志响应 `LogRes`。

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
- `logger`：自定义日志结构体，实现了 `xxl.Logger` 接口。
- `Info` 和 `Error` 方法：分别用于记录信息和错误日志。

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
- `customMiddleware`：自定义中间件函数，用于在任务执行前后添加额外的处理逻辑。

### 总结
这个代码文件实现了一个 `xxl-job` 执行器，配置了执行器的相关参数，注册了任务处理函数，并使用了自定义的日志处理器和中间件。通过这个示例，可以学习如何使用 `xxl-job` 框架在 Go 语言中实现任务调度功能。
  - task
    - panic.go
      Summarized code for panic.go

这个代码文件定义了一个名为 `Panic` 的函数，该函数用于在执行时触发一个 panic。以下是对代码的详细解释：

### 包声明和导入

```go
package task

import (
	"context"
	xxl "github.com/xxl-job/xxl-job-executor-go"
)
```

- `package task`：声明了代码文件属于 `task` 包。
- `import`：导入了两个包：
  - `"context"`：用于处理上下文，通常用于控制goroutine的生命周期和传递请求范围的数据。
  - `xxl "github.com/xxl-job/xxl-job-executor-go"`：导入了 `xxl-job-executor-go` 库，这是一个用于与 XXL-JOB 调度系统集成的库。

### 函数定义

```go
func Panic(cxt context.Context, param *xxl.RunReq) (msg string) {
	panic("test panic")
	return
}
```

- `func Panic(cxt context.Context, param *xxl.RunReq) (msg string)`：定义了一个名为 `Panic` 的函数。
  - `cxt context.Context`：函数的第一个参数是一个 `context.Context` 类型的上下文对象，用于传递上下文信息。
  - `param *xxl.RunReq`：函数的第二个参数是一个指向 `xxl.RunReq` 类型的指针，这个类型通常用于传递任务执行的请求参数。
  - `(msg string)`：函数的返回值是一个字符串，表示返回的消息。

- `panic("test panic")`：在函数体内，调用了 `panic` 函数，并传递了一个字符串 `"test panic"`。这会导致程序立即停止执行，并抛出一个运行时错误。

- `return`：由于 `panic` 会立即终止函数的执行，这里的 `return` 语句实际上不会被执行。

### 总结

这个代码文件的主要功能是定义了一个名为 `Panic` 的函数，该函数在执行时会触发一个 panic，并输出 `"test panic"` 的错误信息。这个函数通常用于测试或演示如何在 Go 程序中触发和处理 panic。由于 `panic` 会立即终止程序的执行，因此这个函数的返回值 `msg` 实际上不会被使用。
    - test.go
      Summarized code for test.go

这段代码是一个简单的Go语言包，用于在XXL-JOB执行器中注册并执行一个测试任务。下面是对代码的详细解释：

### 包声明和导入
```go
package task

import (
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
)
```
- `package task`：声明了代码所在的包名为`task`。
- `import`：导入了所需的包：
  - `"context"`：用于处理上下文。
  - `"fmt"`：用于格式化输入输出。
  - `xxl "github.com/xxl-job/xxl-job-executor-go"`：导入了XXL-JOB执行器的Go语言客户端库。

### 函数定义
```go
func Test(cxt context.Context, param *xxl.RunReq) (msg string) {
	fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + " log_id:" + xxl.Int64ToStr(param.LogID))
	return "test done"
}
```
- `func Test(cxt context.Context, param *xxl.RunReq) (msg string)`：定义了一个名为`Test`的函数，该函数接受两个参数：
  - `cxt context.Context`：一个上下文对象，用于控制函数的生命周期和传递请求范围的数据。
  - `param *xxl.RunReq`：一个指向`xxl.RunReq`结构体的指针，包含了任务执行所需的参数。
- `(msg string)`：声明了函数的返回值类型为`string`，并命名为`msg`。

### 函数实现
- `fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + " log_id:" + xxl.Int64ToStr(param.LogID))`：
  - 使用`fmt.Println`函数打印一条日志信息，包含了任务处理器名称(`param.ExecutorHandler`)、任务参数(`param.ExecutorParams`)和日志ID(`param.LogID`)。
  - `xxl.Int64ToStr(param.LogID)`：将`param.LogID`从`int64`类型转换为字符串类型。
- `return "test done"`：返回字符串`"test done"`，表示任务执行完成。

### 总结
这段代码定义了一个名为`Test`的函数，该函数在XXL-JOB执行器中注册并执行一个测试任务。函数接收任务参数，并打印一条日志信息，最后返回一个表示任务完成的字符串。这个函数可以作为XXL-JOB任务的一个示例，展示了如何在Go语言中实现一个简单的任务处理器。
    - test2.go
      Summarized code for test2.go

这个代码文件定义了一个名为 `Test2` 的函数，该函数用于执行一个任务，并在特定条件下终止任务。以下是对该代码文件的详细解释：

### 包和导入

```go
package task

import (
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"time"
)
```

- `package task`：定义了代码所在的包名为 `task`。
- `import`：导入了以下包：
  - `context`：用于处理上下文，管理任务的生命周期。
  - `fmt`：用于格式化输入和输出。
  - `xxl "github.com/xxl-job/xxl-job-executor-go"`：导入 XXL-JOB 执行器的 Go 客户端库。
  - `time`：用于处理时间相关的操作。

### 函数 `Test2`

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
			fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + "执行行", num)
			if num > 10 {
				fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + "执行完毕！")
				return
			}
		}
	}
}
```

- `func Test2(cxt context.Context, param *xxl.RunReq) (msg string)`：定义了一个名为 `Test2` 的函数，该函数接受两个参数：
  - `cxt context.Context`：上下文对象，用于管理任务的生命周期。
  - `param *xxl.RunReq`：XXL-JOB 运行请求的参数。
  - `msg string`：函数的返回值，表示任务的执行结果。

- `num := 1`：初始化一个计数器 `num`，用于记录任务执行的次数。

- `for { ... }`：无限循环，用于持续执行任务。

- `select { ... }`：使用 `select` 语句监听多个通道。
  - `case <-cxt.Done():`：当上下文被取消时（例如，任务被手动终止），执行以下操作：
    - 打印任务被手动终止的提示信息。
    - 返回，终止任务。
  - `default:`：默认情况下，执行以下操作：
    - `num++`：计数器加一。
    - `time.Sleep(10 * time.Second)`：暂停 10 秒钟。
    - 打印当前任务的执行信息，包括任务处理程序和参数，以及当前的执行次数。
    - `if num > 10 { ... }`：如果计数器超过 10，表示任务执行完毕，打印任务执行完毕的提示信息，并返回，终止任务。

### 总结

该代码文件定义了一个名为 `Test2` 的函数，用于执行一个任务，并在特定条件下终止任务。任务会每隔 10 秒钟执行一次，最多执行 10 次。如果任务被手动终止，或者执行次数超过 10 次，任务将终止。通过使用上下文对象，可以有效地管理任务的生命周期。
- executor.go
  [File too long to summarize]
- go.mod
- log.go
  Summarized code for log.go

这个代码文件定义了一个简单的日志系统，包括应用日志和系统日志的功能。以下是对代码的详细解释：

### 包声明和导入
```go
package xxl

import (
	"fmt"
	"log"
)
```
- `package xxl`：声明了代码所在的包名为 `xxl`。
- `import`：导入了两个标准库包 `fmt` 和 `log`，分别用于格式化输出和日志记录。

### 应用日志类型定义
```go
// LogFunc 应用日志
type LogFunc func(req LogReq, res *LogRes) []byte
```
- `LogFunc`：定义了一个函数类型 `LogFunc`，该函数接受两个参数 `req` 和 `res`，分别表示日志请求和日志响应，并返回一个字节数组。

### 系统日志接口定义
```go
// Logger 系统日志
type Logger interface {
	Info(format string, a ...interface{})
	Error(format string, a ...interface{})
}
```
- `Logger`：定义了一个接口类型 `Logger`，包含两个方法 `Info` 和 `Error`，用于记录信息和错误日志。

### 系统日志实现
```go
type logger struct {
}

func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf(format, a...))
}
```
- `logger`：定义了一个结构体 `logger`，实现了 `Logger` 接口。
- `Info` 方法：使用 `fmt.Println` 打印格式化后的信息日志。
- `Error` 方法：使用 `log.Println` 打印格式化后的错误日志。

### 总结
这个代码文件主要实现了以下功能：
1. 定义了一个应用日志函数类型 `LogFunc`，用于处理日志请求和响应。
2. 定义了一个系统日志接口 `Logger`，包含 `Info` 和 `Error` 两个方法。
3. 实现了一个简单的系统日志结构体 `logger`，通过 `fmt` 和 `log` 包来记录信息和错误日志。

通过阅读这个代码文件，读者可以学习到如何定义函数类型、接口以及如何实现接口方法。同时，也可以了解到如何使用 `fmt` 和 `log` 包进行日志记录。
- log_handler.go
  Summarized code for log_handler.go

这段代码文件定义了一个用于日志查询的功能，主要用于在xxl-job-admin后台显示日志信息。以下是对代码的详细解释：

### 包和导入
```go
package xxl

import (
	"encoding/json"
	"net/http"
)
```
- `package xxl`：定义了包名为`xxl`。
- `import`：导入了两个标准库包，`encoding/json`用于JSON编解码，`net/http`用于HTTP处理。

### 类型定义
```go
type LogHandler func(req *LogReq) *LogRes
```
- `LogHandler`：定义了一个函数类型，该函数接收一个`LogReq`类型的指针参数，并返回一个`LogRes`类型的指针。

### 默认日志处理器
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
- `defaultLogHandler`：默认的日志处理器函数。它接收一个`LogReq`类型的指针参数，并返回一个`LogRes`类型的指针。
  - `Code: SuccessCode`：表示操作成功。
  - `Msg: ""`：消息为空。
  - `Content`：包含日志内容的详细信息。
    - `FromLineNum`：从请求中获取的起始行号。
    - `ToLineNum: 2`：返回的结束行号为2。
    - `LogContent`：默认的日志内容。
    - `IsEnd: true`：表示日志结束。

### 请求错误处理器
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
- `reqErrLogHandler`：处理请求错误的函数。它接收一个`http.ResponseWriter`、一个`LogReq`类型的指针和一个`error`类型的参数。
  - `Code: FailureCode`：表示操作失败。
  - `Msg: err.Error()`：错误信息。
  - `Content`：包含错误日志内容的详细信息。
    - `FromLineNum`：从请求中获取的起始行号。
    - `ToLineNum: 0`：返回的结束行号为0。
    - `LogContent`：错误信息。
    - `IsEnd: true`：表示日志结束。
  - `json.Marshal(res)`：将`LogRes`对象编码为JSON字符串。
  - `w.Write(str)`：将JSON字符串写入HTTP响应。

### 总结
这段代码定义了一个用于日志查询的功能，包括一个默认的日志处理器和一个处理请求错误的日志处理器。这些处理器主要用于在xxl-job-admin后台显示日志信息，并通过HTTP响应返回日志内容或错误信息。
- middleware.go
  Summarized code for middleware.go

这个代码文件定义了一个中间件链的实现，主要用于处理任务函数（TaskFunc）。以下是对代码的详细解释：

### 1. 包声明
```go
package xxl
```
这行代码声明了代码文件所属的包名为 `xxl`。

### 2. 中间件类型定义
```go
// Middleware 中间件构造函数
type Middleware func(TaskFunc) TaskFunc
```
这行代码定义了一个名为 `Middleware` 的函数类型。`Middleware` 是一个函数，它接受一个 `TaskFunc` 类型的参数，并返回一个 `TaskFunc` 类型的值。`TaskFunc` 是一个任务函数类型，具体定义在其他地方。

### 3. 中间件链构造函数
```go
func (e *executor) chain(next TaskFunc) TaskFunc {
	for i := range e.middlewares {
		next = e.middlewares[len(e.middlewares)-1-i](next)
	}
	return next
}
```
这行代码定义了一个名为 `chain` 的方法，属于 `executor` 类型的指针接收者 `e`。该方法的功能是构造一个中间件链。

#### 方法参数和返回值
- `next TaskFunc`：这是一个任务函数，表示中间件链的下一个处理步骤。
- `TaskFunc`：方法返回一个任务函数，即经过所有中间件处理后的最终任务函数。

#### 方法实现细节
- `for i := range e.middlewares`：这是一个循环，遍历 `e.middlewares` 切片中的所有中间件。
- `e.middlewares[len(e.middlewares)-1-i](next)`：这行代码从后向前依次调用中间件函数，并将 `next` 作为参数传递给每个中间件。这样做的目的是确保中间件的执行顺序是从后向前的。
- `next = ...`：每次循环中，`next` 都会被更新为当前中间件处理后的任务函数。
- `return next`：最后返回经过所有中间件处理后的任务函数。

### 总结
这个代码文件的主要功能是定义了一个中间件链的构造方法。通过 `Middleware` 类型定义了中间件函数的形式，并通过 `chain` 方法实现了中间件链的构造。这种方法可以灵活地添加和组合多个中间件，以实现对任务函数的复杂处理逻辑。
- optinos.go
  Summarized code for optinos.go

这段代码定义了一个用于配置执行器（Executor）的选项结构体和相关函数。以下是对代码的详细解释：

### 包和导入
- 代码位于 `xxl` 包中。
- 导入了两个外部包：
  - `github.com/go-basic/ipv4`：用于获取本地IP地址。
  - `time`：用于处理时间相关的操作。

### 结构体 `Options`
- `Options` 结构体用于存储执行器的配置选项。
- 包含以下字段：
  - `ServerAddr`：调度中心地址。
  - `AccessToken`：请求令牌。
  - `Timeout`：接口超时时间。
  - `ExecutorIp`：本地（执行器）IP。
  - `ExecutorPort`：本地（执行器）端口。
  - `RegistryKey`：执行器名称。
  - `LogDir`：日志目录。
  - `l`：日志处理器（类型为 `Logger`，具体实现未在代码中展示）。

### 函数 `newOptions`
- `newOptions` 函数用于创建并初始化 `Options` 结构体。
- 接受一个可变参数 `opts`，类型为 `Option`。
- 初始化 `Options` 结构体时，默认设置 `ExecutorIp` 为本地IP，`ExecutorPort` 为 `"9999"`，`RegistryKey` 为 `"golang-jobs"`。
- 遍历 `opts` 参数，调用每个 `Option` 函数来设置 `Options` 结构体的字段。
- 如果 `l` 字段为空，则默认设置为 `&logger{}`（具体实现未在代码中展示）。

### 类型 `Option`
- `Option` 是一个函数类型，接受一个指向 `Options` 结构体的指针，并对其进行修改。

### 变量
- `DefaultExecutorPort`：默认执行器端口，值为 `"9999"`。
- `DefaultRegistryKey`：默认执行器标识，值为 `"golang-jobs"`。

### 函数（Option 函数）
- 这些函数用于创建 `Option` 类型的函数，以便在 `newOptions` 函数中使用。
- 每个函数都接受一个字符串参数，并返回一个 `Option` 函数，该函数会修改 `Options` 结构体的相应字段。
  - `ServerAddr`：设置调度中心地址。
  - `AccessToken`：设置请求令牌。
  - `ExecutorIp`：设置执行器IP。
  - `ExecutorPort`：设置执行器端口。
  - `RegistryKey`：设置执行器标识。
  - `SetLogger`：设置日志处理器。

### 总结
这段代码通过定义 `Options` 结构体和一系列 `Option` 函数，提供了一种灵活的方式来配置执行器的各项参数。通过使用这些函数，可以方便地设置和修改执行器的配置，而不需要直接操作结构体字段。这种设计模式在Go语言中非常常见，特别是在需要配置复杂对象时。
- task.go
  Summarized code for task.go

这段代码定义了一个任务执行框架，主要功能是管理和执行任务，并在任务执行过程中处理异常情况。以下是对代码的详细解释：

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
- `import`：导入了 `context`、`fmt` 和 `runtime/debug` 包，用于上下文管理、格式化输出和调试信息打印。

### 类型定义
#### TaskFunc
```go
// TaskFunc 任务执行函数
type TaskFunc func(cxt context.Context, param *RunReq) string
```
- `TaskFunc` 是一个函数类型，定义了任务执行函数的签名。该函数接收一个 `context.Context` 和一个 `*RunReq` 类型的参数，并返回一个字符串。

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
- `Task` 结构体定义了一个任务的详细信息，包括任务ID、名称、上下文、参数、执行函数、取消函数、开始和结束时间以及日志记录器。

### 方法定义
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
- `Run` 方法用于执行任务。它使用了 `defer` 语句来捕获任务执行过程中可能发生的 panic，并进行相应的处理：
  - 如果捕获到 panic，会记录日志并打印堆栈跟踪信息。
  - 调用 `callback` 函数通知任务执行结果，并调用 `cancel` 函数取消任务。
  - 如果没有发生 panic，则正常执行任务函数 `t.fn`，并调用 `callback` 函数通知任务执行成功。

#### Info
```go
// Info 任务信息
func (t *Task) Info() string {
	return fmt.Sprintf("任务ID[%d]任务名称[%s]参数:%s", t.Id, t.Name, t.Param.ExecutorParams)
}
```
- `Info` 方法返回任务的详细信息字符串，包括任务ID、名称和参数。

### 总结
这段代码实现了一个简单的任务执行框架，通过定义 `Task` 结构体和相关方法，可以方便地管理和执行任务，并在任务执行过程中处理异常情况。通过使用 `context` 包，可以有效地管理任务的生命周期和取消操作。
- task_list.go
  Summarized code for task_list.go

这个代码文件定义了一个用于管理任务列表的数据结构和相关操作方法。以下是对代码的详细解释：

### 包和导入
```go
package xxl

import "sync"
```
- 代码位于 `xxl` 包中。
- 导入了 `sync` 包，用于实现并发安全的操作。

### 数据结构
```go
// 任务列表 [JobID]执行函数,并行执行时[+LogID]
type taskList struct {
	mu   sync.RWMutex
	data map[string]*Task
}
```
- `taskList` 是一个结构体，用于存储任务列表。
- `mu` 是一个读写锁（`sync.RWMutex`），用于保证对 `data` 的并发访问安全。
- `data` 是一个映射（`map`），键为字符串类型（`string`），值为指向 `Task` 结构体的指针（`*Task`）。

### 方法
#### Set
```go
// Set 设置数据
func (t *taskList) Set(key string, val *Task) {
	t.mu.Lock()
	t.data[key] = val
	t.mu.Unlock()
}
```
- `Set` 方法用于向任务列表中添加或更新一个任务。
- 使用 `mu.Lock()` 和 `mu.Unlock()` 确保对 `data` 的写操作是线程安全的。

#### Get
```go
// Get 获取数据
func (t *taskList) Get(key string) *Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data[key]
}
```
- `Get` 方法用于从任务列表中获取一个任务。
- 使用 `mu.RLock()` 和 `mu.RUnlock()` 确保对 `data` 的读操作是线程安全的。

#### GetAll
```go
// GetAll 获取数据
func (t *taskList) GetAll() map[string]*Task {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.data
}
```
- `GetAll` 方法用于获取整个任务列表的副本。
- 使用 `mu.RLock()` 和 `mu.RUnlock()` 确保对 `data` 的读操作是线程安全的。

#### Del
```go
// Del 设置数据
func (t *taskList) Del(key string) {
	t.mu.Lock()
	delete(t.data, key)
	t.mu.Unlock()
}
```
- `Del` 方法用于从任务列表中删除一个任务。
- 使用 `mu.Lock()` 和 `mu.Unlock()` 确保对 `data` 的写操作是线程安全的。

#### Len
```go
// Len 长度
func (t *taskList) Len() int {
	return len(t.data)
}
```
- `Len` 方法用于获取任务列表的长度（即任务的数量）。
- 由于 `len` 函数是原子的，不需要加锁。

#### Exists
```go
// Exists Key是否存在
func (t *taskList) Exists(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.data[key]
	return ok
}
```
- `Exists` 方法用于检查某个键是否存在于任务列表中。
- 使用 `mu.RLock()` 和 `mu.RUnlock()` 确保对 `data` 的读操作是线程安全的。

### 总结
这个代码文件定义了一个并发安全的任务列表管理结构 `taskList`，并提供了添加、获取、删除、检查存在性和获取长度的方法。通过使用读写锁 `sync.RWMutex`，确保了在多线程环境下对任务列表的操作是安全的。
- util.go
  Summarized code for util.go

这个代码文件定义了一些用于处理任务调度和管理的函数，主要功能是将不同类型的任务状态转换为JSON格式的字符串，便于网络传输和日志记录。以下是对每个函数的详细解释：

### 1. `Int64ToStr` 函数
```go
// Int64ToStr int64 to str
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}
```
- **功能**：将 `int64` 类型的整数转换为字符串。
- **实现细节**：使用 `strconv.FormatInt` 函数，将 `int64` 类型的整数转换为十进制字符串。

### 2. `returnCall` 函数
```go
//执行任务回调
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
- **功能**：将任务执行结果封装为JSON格式的字符串。
- **实现细节**：
  - 创建一个 `call` 类型的结构体，包含一个 `callElement` 类型的元素。
  - `callElement` 结构体包含日志ID、日志时间、执行结果等信息。
  - 使用 `json.Marshal` 函数将 `data` 转换为JSON格式的字节数组。

### 3. `returnKill` 函数
```go
//杀死任务返回
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
- **功能**：将任务杀死的结果封装为JSON格式的字符串。
- **实现细节**：
  - 根据 `code` 的值判断是否成功杀死任务，并设置相应的消息。
  - 创建一个 `res` 类型的结构体，包含代码和消息。
  - 使用 `json.Marshal` 函数将 `data` 转换为JSON格式的字节数组。

### 4. `returnIdleBeat` 函数
```go
//忙碌返回
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
- **功能**：将任务忙碌状态封装为JSON格式的字符串。
- **实现细节**：
  - 根据 `code` 的值判断任务是否忙碌，并设置相应的消息。
  - 创建一个 `res` 类型的结构体，包含代码和消息。
  - 使用 `json.Marshal` 函数将 `data` 转换为JSON格式的字节数组。

### 5. `returnGeneral` 函数
```go
//通用返回
func returnGeneral() []byte {
	data := &res{
		Code: SuccessCode,
		Msg:  "",
	}
	str, _ := json.Marshal(data)
	return str
}
```
- **功能**：返回一个通用的成功状态的JSON格式的字符串。
- **实现细节**：
  - 创建一个 `res` 类型的结构体，包含成功的代码和空消息。
  - 使用 `json.Marshal` 函数将 `data` 转换为JSON格式的字节数组。

### 总结
这个代码文件主要用于将任务的不同状态（执行结果、杀死结果、忙碌状态等）转换为JSON格式的字符串，便于在网络中传输和记录日志。每个函数都使用了 `json.Marshal` 函数来实现JSON格式的转换。
