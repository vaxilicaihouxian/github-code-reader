### Repository Summary

#### Purpose
# JSONata in Go

Package jsonata is a query and transformation language for JSON.
It's a Go port of the JavaScript library [JSONata](http://jsonata.org/).

It currently has feature parity with jsonata-js 1.5.4. As well as a most of the functions added in newer versions. You can see potentially missing functions by looking at the [jsonata-js changelog](https://github.com/jsonata-js/jsonata/blob/master/CHANGELOG.md).

## Install

    go get github.com/blues/jsonata-go

## Usage

```Go
import (
	"encoding/json"
	"fmt"
	"log"

	jsonata "github.com/blues/jsonata-go"
)

const jsonString = `
    {
        "orders": [
            {"price": 10, "quantity": 3},
            {"price": 0.5, "quantity": 10},
            {"price": 100, "quantity": 1}
        ]
    }
`

func main() {

	var data interface{}

	// Decode JSON.
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		log.Fatal(err)
	}

	// Create expression.
	e := jsonata.MustCompile("$sum(orders.(price*quantity))")

	// Evaluate.
	res, err := e.Eval(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	// Output: 135
}
```

## JSONata Server
A locally hosted version of [JSONata Exerciser](http://try.jsonata.org/)
for testing is [available here](https://github.com/blues/jsonata-go/jsonata-server).

## JSONata tests
A CLI tool for running jsonata-go against the [JSONata test suite](https://github.com/jsonata-js/jsonata/tree/master/test/test-suite) is [available here](./jsonata-test).



## Contributing

We love issues, fixes, and pull requests from everyone. Please run the
unit-tests, staticcheck, and goimports prior to submitting your PR. By participating in this project, you agree to abide by
the Blues Inc [code of conduct](https://blues.github.io/opensource/code-of-conduct).

For details on contributions we accept and the process for contributing, see our
[contribution guide](CONTRIBUTING.md).

In addition to the Go unit tests there is also a test runner that will run against the jsonata-js test
suite in the [jsonata-test](./jsonata-test) directory. A number of these tests currently fail, but we're working towards feature parity with the jsonata-js reference implementation. Pull requests welcome!

If you would like to contribute to this library a good first issue would be to run the jsonata-test suite,
and fix any of the tests not passing.


#### Structure
- .github
  - workflows
    - test.yml
- .gitignore
- CODE_OF_CONDUCT.md
- CONTRIBUTING.md
- LICENSE
- README.md
- callable.go
  [File too long to summarize]
- callable_test.go
  [File too long to summarize]
- doc.go
  Summarized code for doc.go

### 代码文件总结

#### 文件信息
- **版权声明**：该代码文件的版权归 Blues Inc. 所有，使用该源代码需遵守版权持有者授予的许可，具体许可内容可在 `LICENSE` 文件中找到。
- **包名**：`jsonata`

#### 包描述
- **功能**：`jsonata` 包是一个用于 JSON 的查询和转换语言。它是 JavaScript 库 JSONata 的 Go 语言移植版本。
- **参考资料**：建议用户参考官方 JSONata 网站以获取语言参考。
- **官方网站**：[http://jsonata.org/](http://jsonata.org/)

#### 详细功能和实现细节
1. **查询和转换语言**：
   - JSONata 是一种强大的查询和转换语言，专门用于处理 JSON 数据。它允许用户通过简洁的语法从复杂的 JSON 结构中提取和转换数据。

2. **Go 语言移植**：
   - 该包是 JSONata 的 Go 语言实现，这意味着它将 JavaScript 库的功能和语法移植到了 Go 语言环境中。这使得 Go 开发者可以使用 JSONata 的功能来处理 JSON 数据。

3. **使用方法**：
   - 用户可以通过导入 `jsonata` 包来使用其提供的功能。具体的使用方法和示例可以参考官方 JSONata 网站上的文档和教程。

4. **许可和版权**：
   - 使用该包的源代码需要遵守版权持有者授予的许可，具体许可内容可在 `LICENSE` 文件中找到。这确保了代码的合法使用和分发。

#### 学习建议
- **阅读官方文档**：建议读者访问 [http://jsonata.org/](http://jsonata.org/) 以获取 JSONata 语言的详细文档和示例，这将有助于更好地理解和使用该包。
- **实践操作**：通过实际编写和运行代码来熟悉 JSONata 的语法和功能，这将加深对包的理解和应用能力。

### 总结
该代码文件是一个 Go 语言包，名为 `jsonata`，它是 JSONata 查询和转换语言的 Go 语言移植版本。该包允许 Go 开发者使用 JSONata 的功能来处理 JSON 数据。使用该包时，需遵守版权持有者授予的许可，并建议参考官方 JSONata 网站以获取语言参考和使用指南。
- env.go
  Summarized code for env.go

这个代码文件定义了一个用于处理JSONata表达式的环境（environment）和一系列函数扩展（extensions）。JSONata是一种轻量级的查询和转换JSON数据的语言。以下是对代码文件的详细解析：

### 1. 包和导入

```go
package jsonata

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/blues/jsonata-go/jlib"
	"github.com/blues/jsonata-go/jparse"
	"github.com/blues/jsonata-go/jtypes"
)
```

- **包名**：`jsonata`
- **导入的包**：包括标准库中的`errors`、`math`、`reflect`、`strings`、`unicode/utf8`，以及第三方包`jlib`、`jparse`、`jtypes`。

### 2. 环境（Environment）

```go
type environment struct {
	parent  *environment
	symbols map[string]reflect.Value
}

func newEnvironment(parent *environment, size int) *environment {
	return &environment{
		parent:  parent,
		symbols: make(map[string]reflect.Value, size),
	}
}

func (s *environment) bind(name string, value reflect.Value) {
	if s.symbols == nil {
		s.symbols = make(map[string]reflect.Value)
	}
	s.symbols[name] = value
}

func (s *environment) bindAll(values map[string]reflect.Value) {
	if len(values) == 0 {
		return
	}
	for name, value := range values {
		s.bind(name, value)
	}
}

func (s *environment) lookup(name string) reflect.Value {
	if v, ok := s.symbols[name]; ok {
		return v
	}
	if s.parent != nil {
		return s.parent.lookup(name)
	}
	return undefined
}
```

- **环境结构体**：`environment`包含一个父环境指针和一个符号表（`map[string]reflect.Value`）。
- **创建环境**：`newEnvironment`函数用于创建一个新的环境，可以指定父环境和初始大小。
- **绑定符号**：`bind`方法用于将一个符号（名称）绑定到一个值。
- **批量绑定符号**：`bindAll`方法用于批量绑定多个符号和值。
- **查找符号**：`lookup`方法用于查找符号的值，如果当前环境找不到，会递归查找父环境。

### 3. 默认处理函数

```go
var (
	defaultUndefinedHandler = jtypes.ArgUndefined(0)
	defaultContextHandler   = jtypes.ArgCountEquals(0)
	argCountEquals1         = jtypes.ArgCountEquals(1)
)
```

- **默认未定义处理函数**：`defaultUndefinedHandler`用于处理未定义参数。
- **默认上下文处理函数**：`defaultContextHandler`用于处理上下文参数数量等于0的情况。
- **参数数量等于1的处理函数**：`argCountEquals1`用于处理参数数量等于1的情况。

### 4. 基础环境初始化

```go
var baseEnv = initBaseEnv(map[string]Extension{
	// 省略具体扩展定义
})

func initBaseEnv(exts map[string]Extension) *environment {
	env := newEnvironment(nil, len(exts))
	for name, ext := range exts {
		fn := mustGoCallable(name, ext)
		env.bind(name, reflect.ValueOf(fn))
	}
	return env
}

func mustGoCallable(name string, ext Extension) *goCallable {
	callable, err := newGoCallable(name, ext)
	if err != nil {
		panicf("%s is not a valid function: %s", name, err)
	}
	return callable
}
```

- **基础环境**：`baseEnv`是一个全局变量，通过`initBaseEnv`函数初始化。
- **初始化基础环境**：`initBaseEnv`函数创建一个新的环境，并将所有扩展函数绑定到环境中。
- **确保可调用**：`mustGoCallable`函数确保扩展函数是有效的可调用对象。

### 5. 本地函数和处理函数

```go
func lookup(v reflect.Value, name string) (interface{}, error) {
	res, err := evalName(&jparse.NameNode{Value: name}, v, nil)
	if err != nil {
		return nil, err
	}
	if seq, ok := asSequence(res); ok {
		res = seq.Value()
	}
	if res.IsValid() && res.CanInterface() {
		return res.Interface(), nil
	}
	return nil, nil
}

func throw(msg string) (interface{}, error) {
	return nil, errors.New(msg)
}

// 省略未定义处理函数和上下文处理函数
```

- **查找函数**：`lookup`函数用于查找指定名称的值。
- **抛出错误**：`throw`函数用于抛出错误消息。
- **未定义处理函数**：`undefinedHandlerAppend`等函数用于处理未定义参数的情况。
- **上下文处理函数**：`contextHandlerSubstring`等函数用于处理上下文参数的情况。

### 6. 扩展函数定义

```go
// 省略具体扩展定义
```

- **字符串函数**：如`string`、`length`、`substring`等。
- **数字函数**：如`number`、`abs`、`floor`等。
- **布尔函数**：如`boolean`、`not`、`exists`等。
- **数组函数**：如`distinct`、`count`、`reverse`等。
- **对象函数**：如`each`、`sift`、`keys`等。
- **日期函数**：如`fromMillis`、`toMillis`等。
- **其他函数**：如`type`、`error`等。

### 总结

这个代码文件定义了一个用于处理JSONata表达式的环境，并提供了一系列函数扩展。通过环境结构体和相关方法，可以管理和查找符号。通过初始化基础环境，将所有扩展函数绑定到环境中。此外，还定义了一些本地函数和处理函数，用于处理未定义参数和上下文参数的情况。
- error.go
  Summarized code for error.go

### 代码文件内容总结

#### 文件信息
- **版权声明**：该文件的版权归 Blues Inc. 所有，使用该源代码需遵守相关许可证规定。

#### 包声明
- **包名**：`jsonata`

#### 导入的包
- `errors`：用于创建和操作错误。
- `fmt`：用于格式化输入和输出。
- `regexp`：用于正则表达式操作。
- `github.com/blues/jsonata-go/jtypes`：用于处理 JSONata 表达式的类型。

#### 常量和变量
- **错误类型**：定义了一系列错误类型常量，用于标识不同类型的错误。
- **错误消息映射**：`errmsgs` 是一个映射，将错误类型与错误消息模板关联起来。
- **正则表达式**：`reErrMsg` 用于匹配错误消息模板中的占位符。

#### 错误类型定义
- **ErrType**：错误类型的枚举类型。
- **EvalError**：表示 JSONata 表达式评估过程中的错误。
  - **字段**：
    - `Type`：错误类型。
    - `Token`：相关的令牌（如函数名、操作符等）。
    - `Value`：相关的值（如操作数、键值等）。
  - **方法**：
    - `newEvalError`：创建一个新的 `EvalError` 实例。
    - `Error`：返回错误消息字符串。

#### 其他错误类型
- **ArgCountError**：表示函数调用时参数数量错误的错误类型。
  - **字段**：
    - `Func`：函数名。
    - `Expected`：期望的参数数量。
    - `Received`：实际接收的参数数量。
  - **方法**：
    - `newArgCountError`：创建一个新的 `ArgCountError` 实例。
    - `Error`：返回错误消息字符串。

- **ArgTypeError**：表示函数调用时参数类型错误的错误类型。
  - **字段**：
    - `Func`：函数名。
    - `Which`：错误的参数位置。
  - **方法**：
    - `newArgTypeError`：创建一个新的 `ArgTypeError` 实例。
    - `Error`：返回错误消息字符串。

### 功能和实现细节

#### 错误处理
- **错误类型定义**：通过枚举类型 `ErrType` 定义了多种错误类型，每种类型对应一种特定的错误情况。
- **错误消息模板**：使用 `errmsgs` 映射将错误类型与错误消息模板关联起来，便于根据错误类型生成具体的错误消息。
- **错误实例化**：通过 `newEvalError`、`newArgCountError` 和 `newArgTypeError` 函数创建具体的错误实例，这些函数会根据传入的参数生成相应的错误对象。
- **错误消息生成**：通过 `EvalError`、`ArgCountError` 和 `ArgTypeError` 的 `Error` 方法生成具体的错误消息字符串，这些方法会根据错误对象的字段值替换错误消息模板中的占位符。

#### 正则表达式
- **占位符替换**：使用正则表达式 `reErrMsg` 匹配错误消息模板中的占位符 `{{token}}` 和 `{{value}}`，并在生成错误消息时进行替换。

#### 类型转换
- **字符串化**：在 `newEvalError` 函数中，使用 `stringify` 函数将传入的参数转换为字符串，以便在错误消息中使用。

### 总结
该代码文件主要定义了 JSONata 表达式评估过程中可能遇到的各种错误类型，并提供了相应的错误消息生成机制。通过枚举类型、错误消息模板和正则表达式等技术，实现了错误类型的定义、错误实例的创建和错误消息的生成，便于在 JSONata 表达式评估过程中进行错误处理和调试。
- eval.go
  [File too long to summarize]
- eval_test.go
  [File too long to summarize]
- example_eval_test.go
  Summarized code for example_eval_test.go

这段代码文件是一个使用 JSONata 表达式进行 JSON 数据处理的示例。JSONata 是一种用于查询和转换 JSON 数据的强大工具。以下是对该代码文件的详细解释：

### 文件信息
- **版权声明**：代码由 Blues Inc. 所有，使用该源代码需遵守相关许可证。

### 包声明
- **包名**：`jsonata_test`，表示这是一个测试包。

### 导入的包
- `encoding/json`：用于 JSON 的编码和解码。
- `fmt`：用于格式化输入输出。
- `log`：用于记录日志。
- `github.com/blues/jsonata-go`：JSONata 的 Go 语言实现。

### 常量
- `jsonString`：一个包含订单信息的 JSON 字符串。

### 示例函数
- **函数名**：`ExampleExpr_Eval`，这是一个示例函数，展示了如何使用 JSONata 表达式进行计算。

#### 详细步骤
1. **解码 JSON**：
    - 使用 `json.Unmarshal` 将 JSON 字符串解码为 Go 语言的 `interface{}` 类型。
    - 如果解码失败，使用 `log.Fatal` 记录错误并终止程序。

2. **创建表达式**：
    - 使用 `jsonata.MustCompile` 编译 JSONata 表达式 `$sum(orders.(price*quantity))`。
    - 该表达式的含义是：计算所有订单中 `price` 和 `quantity` 的乘积之和。

3. **评估表达式**：
    - 使用 `e.Eval` 方法对解码后的 JSON 数据进行评估。
    - 如果评估失败，使用 `log.Fatal` 记录错误并终止程序。

4. **输出结果**：
    - 使用 `fmt.Println` 输出评估结果。
    - 预期输出为 `135`，即所有订单的总金额。

### 代码逻辑
1. **JSON 解码**：
    ```go
    var data interface{}
    err := json.Unmarshal([]byte(jsonString), &data)
    if err != nil {
        log.Fatal(err)
    }
    ```
    这段代码将 JSON 字符串解码为 Go 语言的 `interface{}` 类型，以便后续处理。

2. **编译 JSONata 表达式**：
    ```go
    e := jsonata.MustCompile("$sum(orders.(price*quantity))")
    ```
    这段代码编译了一个 JSONata 表达式，用于计算所有订单的总金额。

3. **评估表达式**：
    ```go
    res, err := e.Eval(data)
    if err != nil {
        log.Fatal(err)
    }
    ```
    这段代码对解码后的 JSON 数据进行评估，并获取结果。

4. **输出结果**：
    ```go
    fmt.Println(res)
    ```
    这段代码将评估结果输出到控制台。

### 总结
该代码文件展示了如何使用 JSONata 表达式对 JSON 数据进行查询和计算。通过解码 JSON 数据、编译 JSONata 表达式、评估表达式并输出结果，实现了对订单总金额的计算。这个示例对于学习和理解 JSONata 在 Go 语言中的应用非常有帮助。
- example_exts_test.go
  Summarized code for example_exts_test.go

这段代码文件展示了如何使用 Go 语言中的 JSONata 库来扩展其功能，具体是通过添加自定义函数。以下是对代码的详细解释：

### 1. 导入依赖包
```go
import (
	"fmt"
	"log"
	"strings"
	jsonata "github.com/blues/jsonata-go"
)
```
- `fmt` 和 `log` 是 Go 标准库中的包，用于格式化输出和日志记录。
- `strings` 是 Go 标准库中的包，提供了字符串操作的函数。
- `jsonata` 是来自 `github.com/blues/jsonata-go` 的包，用于处理 JSONata 表达式。

### 2. 定义自定义函数
```go
var exts = map[string]jsonata.Extension{
	"titlecase": {
		Func: strings.Title,
	},
}
```
- `exts` 是一个映射，定义了一个名为 `titlecase` 的函数，该函数映射到 Go 标准库中的 `strings.Title` 函数。
- `jsonata.Extension` 是一个结构体，用于定义 JSONata 扩展函数。每个扩展函数必须返回一个或两个值，第二个值必须是错误类型。

### 3. 示例函数
```go
func ExampleExpr_RegisterExts() {
	// Create an expression that uses the titlecase function.
	e := jsonata.MustCompile(`$titlecase("beneath the underdog")`)

	// Register the titlecase function.
	err := e.RegisterExts(exts)
	if err != nil {
		log.Fatal(err)
	}

	// Evaluate.
	res, err := e.Eval(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	// Output: Beneath The Underdog
}
```
- `ExampleExpr_RegisterExts` 是一个示例函数，展示了如何使用和注册自定义的 JSONata 函数。
- `jsonata.MustCompile` 用于编译 JSONata 表达式，这里编译了一个使用 `titlecase` 函数的表达式。
- `e.RegisterExts(exts)` 用于注册自定义函数 `titlecase`。
- `e.Eval(nil)` 用于评估编译后的表达式，并返回结果。
- 最后，使用 `fmt.Println` 打印评估结果，预期输出为 `Beneath The Underdog`。

### 总结
这段代码的主要功能是展示如何通过自定义函数扩展 JSONata 的功能。具体来说，它定义了一个名为 `titlecase` 的函数，并将其映射到 Go 标准库中的 `strings.Title` 函数。然后，通过编译和评估一个使用该自定义函数的 JSONata 表达式，展示了整个过程。
- go.mod
- go.sum
- jlib
  - aggregate.go
    Summarized code for aggregate.go

这个代码文件定义了一个名为 `jlib` 的 Go 包，提供了一些基本的数组操作函数，包括求和、求最大值、求最小值和求平均值。这些函数都接受一个 `reflect.Value` 类型的参数，并返回一个 `float64` 类型的结果和一个 `error` 类型的错误信息。以下是对每个函数的详细解释：

### 1. `Sum` 函数

```go
// Sum returns the total of an array of numbers. If the array is
// empty, Sum returns 0.
func Sum(v reflect.Value) (float64, error) {
    if !jtypes.IsArray(v) {
        if n, ok := jtypes.AsNumber(v); ok {
            return n, nil
        }
        return 0, fmt.Errorf("cannot call sum on a non-array type")
    }

    v = jtypes.Resolve(v)

    var sum float64

    for i := 0; i < v.Len(); i++ {
        n, ok := jtypes.AsNumber(v.Index(i))
        if !ok {
            return 0, fmt.Errorf("cannot call sum on an array with non-number types")
        }
        sum += n
    }

    return sum, nil
}
```

- **功能**：计算数组中所有数字的总和。
- **参数**：`v` 是一个 `reflect.Value` 类型的参数，表示一个数组或单个数字。
- **返回值**：返回一个 `float64` 类型的总和和一个 `error` 类型的错误信息。
- **实现细节**：
  - 首先检查 `v` 是否是一个数组，如果不是数组，尝试将其转换为数字。如果转换成功，直接返回该数字；否则返回错误信息。
  - 调用 `jtypes.Resolve(v)` 解析 `v`，确保它是一个数组。
  - 遍历数组中的每个元素，将其转换为数字并累加到 `sum` 中。如果遇到非数字类型的元素，返回错误信息。
  - 最后返回累加的总和。

### 2. `Max` 函数

```go
// Max returns the largest value in an array of numbers. If the
// array is empty, Max returns 0 and an undefined error.
func Max(v reflect.Value) (float64, error) {
    if !jtypes.IsArray(v) {
        if n, ok := jtypes.AsNumber(v); ok {
            return n, nil
        }
        return 0, fmt.Errorf("cannot call max on a non-array type")
    }

    v = jtypes.Resolve(v)
    if v.Len() == 0 {
        return 0, jtypes.ErrUndefined
    }

    var max float64

    for i := 0; i < v.Len(); i++ {
        n, ok := jtypes.AsNumber(v.Index(i))
        if !ok {
            return 0, fmt.Errorf("cannot call max on an array with non-number types")
        }
        if i == 0 || n > max {
            max = n
        }
    }

    return max, nil
}
```

- **功能**：找出数组中最大的数字。
- **参数**：`v` 是一个 `reflect.Value` 类型的参数，表示一个数组或单个数字。
- **返回值**：返回一个 `float64` 类型的最大值和一个 `error` 类型的错误信息。
- **实现细节**：
  - 首先检查 `v` 是否是一个数组，如果不是数组，尝试将其转换为数字。如果转换成功，直接返回该数字；否则返回错误信息。
  - 调用 `jtypes.Resolve(v)` 解析 `v`，确保它是一个数组。
  - 如果数组为空，返回 `0` 和一个未定义的错误。
  - 遍历数组中的每个元素，将其转换为数字并更新 `max` 的值。如果遇到非数字类型的元素，返回错误信息。
  - 最后返回找到的最大值。

### 3. `Min` 函数

```go
// Min returns the smallest value in an array of numbers. If the
// array is empty, Min returns 0 and an undefined error.
func Min(v reflect.Value) (float64, error) {
    if !jtypes.IsArray(v) {
        if n, ok := jtypes.AsNumber(v); ok {
            return n, nil
        }
        return 0, fmt.Errorf("cannot call min on a non-array type")
    }

    v = jtypes.Resolve(v)
    if v.Len() == 0 {
        return 0, jtypes.ErrUndefined
    }

    var min float64

    for i := 0; i < v.Len(); i++ {
        n, ok := jtypes.AsNumber(v.Index(i))
        if !ok {
            return 0, fmt.Errorf("cannot call min on an array with non-number types")
        }
        if i == 0 || n < min {
            min = n
        }
    }

    return min, nil
}
```

- **功能**：找出数组中最小的数字。
- **参数**：`v` 是一个 `reflect.Value` 类型的参数，表示一个数组或单个数字。
- **返回值**：返回一个 `float64` 类型的最小值和一个 `error` 类型的错误信息。
- **实现细节**：
  - 首先检查 `v` 是否是一个数组，如果不是数组，尝试将其转换为数字。如果转换成功，直接返回该数字；否则返回错误信息。
  - 调用 `jtypes.Resolve(v)` 解析 `v`，确保它是一个数组。
  - 如果数组为空，返回 `0` 和一个未定义的错误。
  - 遍历数组中的每个元素，将其转换为数字并更新 `min` 的值。如果遇到非数字类型的元素，返回错误信息。
  - 最后返回找到的最小值。

### 4. `Average` 函数

```go
// Average returns the mean of an array of numbers. If the array
// is empty, Average returns 0 and an undefined error.
func Average(v reflect.Value) (float64, error) {
    if !jtypes.IsArray(v) {
        if n, ok := jtypes.AsNumber(v); ok {
            return n, nil
        }
        return 0, fmt.Errorf("cannot call average on a non-array type")
    }

    v = jtypes.Resolve(v)
    if v.Len() == 0 {
        return 0, jtypes.ErrUndefined
    }

    var sum float64

    for i := 0; i < v.Len(); i++ {
        n, ok := jtypes.AsNumber(v.Index(i))
        if !ok {
            return 0, fmt.Errorf("cannot call average on an array with non-number types")
        }
        sum += n
    }

    return sum / float64(v.Len()), nil
}
```

- **功能**：计算数组中所有数字的平均值。
- **参数**：`v` 是一个 `reflect.Value` 类型的参数，表示一个数组或单个数字。
- **返回值**：返回一个 `float64` 类型的平均值和一个 `error` 类型的错误信息。
- **实现细节**：
  - 首先检查 `v` 是否是一个数组，如果不是数组，尝试将其转换为数字。如果转换成功，直接返回该数字；否则返回错误信息。
  - 调用 `jtypes.Resolve(v)` 解析 `v`，确保它是一个数组。
  - 如果数组为空，返回 `0` 和一个未定义的错误。
  - 遍历数组中的每个元素，将其转换为数字并累加到 `sum` 中。如果遇到非数字类型的元素，返回错误信息。
  - 最后返回累加的总和除以数组长度得到的平均值。

### 总结

这个代码文件通过使用 `reflect` 包和 `jtypes` 包中的函数，实现了对数组或单个数字的求和、求最大值、求最小值和求平均值的功能。每个函数都包含了类型检查、错误处理和具体的计算逻辑，确保在处理不同类型的输入时能够正确地返回结果或错误信息。
  - array.go
    Summarized code for array.go

这个代码文件定义了一系列函数，用于处理和操作Go语言中的反射值（`reflect.Value`）。这些函数主要用于处理JSON数据，提供了计数、去重、追加、反转、排序、洗牌和压缩等操作。以下是对每个函数的详细解释：

### 1. `Count` 函数
- **功能**：计算给定值中元素的数量。
- **实现细节**：
  - 使用 `jtypes.Resolve` 解析值。
  - 如果值不是数组，则返回1（如果有效）或0（如果无效）。
  - 如果是数组，则返回数组的长度。

### 2. `Distinct` 函数
- **功能**：返回给定值中的唯一元素。
- **实现细节**：
  - 使用 `jtypes.Resolve` 解析值。
  - 如果值是字符串，则返回整个字符串（不处理单个字符）。
  - 如果值是数组，则遍历数组，使用 `map` 记录已访问的元素，并将唯一元素添加到结果切片中。
  - 如果元素是 `map`，则将其转换为字符串进行去重。

### 3. `Append` 函数
- **功能**：将两个值追加到一起。
- **实现细节**：
  - 使用 `arrayify` 将值转换为数组。
  - 遍历两个数组，将有效元素添加到结果切片中。

### 4. `Reverse` 函数
- **功能**：反转给定值中的元素顺序。
- **实现细节**：
  - 使用 `arrayify` 将值转换为数组。
  - 从后向前遍历数组，将有效元素添加到结果切片中。

### 5. `Sort` 函数
- **功能**：对给定值进行排序。
- **实现细节**：
  - 使用 `jtypes.Resolve` 解析值。
  - 根据值的类型（数组、数字数组、字符串数组）进行排序。
  - 如果提供了自定义排序函数，则使用自定义排序函数进行排序。

### 6. `sortNumberArray` 函数
- **功能**：对数字数组进行排序。
- **实现细节**：
  - 遍历数组，将数字添加到结果切片中。
  - 使用 `sort.SliceStable` 对结果切片进行排序。

### 7. `sortStringArray` 函数
- **功能**：对字符串数组进行排序。
- **实现细节**：
  - 遍历数组，将字符串添加到结果切片中。
  - 使用 `sort.SliceStable` 对结果切片进行排序。

### 8. `sortArrayFunc` 函数
- **功能**：使用自定义排序函数对数组进行排序。
- **实现细节**：
  - 遍历数组，将元素添加到结果切片中。
  - 使用 `mergeSort` 和 `merge` 函数进行排序。

### 9. `mergeSort` 函数
- **功能**：实现归并排序。
- **实现细节**：
  - 递归地将数组分成两部分，分别进行排序，然后合并结果。

### 10. `merge` 函数
- **功能**：合并两个已排序的数组。
- **实现细节**：
  - 使用自定义排序函数比较元素，将较小的元素添加到结果数组中。

### 11. `Shuffle` 函数
- **功能**：随机打乱给定值中的元素顺序。
- **实现细节**：
  - 使用 `forceArray` 将值转换为数组。
  - 遍历数组，随机交换元素位置。

### 12. `Zip` 函数
- **功能**：将多个数组合并成一个二维数组。
- **实现细节**：
  - 使用 `forceArray` 将每个值转换为数组。
  - 遍历每个数组，将对应位置的元素合并到结果数组中。

### 辅助函数
- **`forceArray`**：将值转换为数组。
- **`arrayLen`**：获取数组的长度。
- **`arrayify`**：将值转换为数组，如果值不是数组，则将其包装成数组。

### 常量
- **`typeInterface`**：定义了一个接口类型的反射类型。

这些函数通过反射机制处理不同类型的值，提供了丰富的操作功能，适用于处理JSON数据等场景。
  - boolean.go
    Summarized code for boolean.go

这段代码文件定义了一个名为 `jlib` 的 Go 包，主要功能是提供一些用于处理不同类型值的布尔转换和存在性检查的函数。以下是对代码的详细解释：

### 包声明和导入

```go
package jlib

import (
	"reflect"

	"github.com/blues/jsonata-go/jtypes"
)
```

- `package jlib`：声明了包名为 `jlib`。
- `import`：导入了两个包：
  - `"reflect"`：用于反射操作，处理不同类型的值。
  - `"github.com/blues/jsonata-go/jtypes"`：一个外部包，提供了一些类型转换和检查的辅助函数。

### 函数 `Boolean`

```go
// Boolean (golint)
func Boolean(v reflect.Value) bool {

	v = jtypes.Resolve(v)

	if b, ok := jtypes.AsBool(v); ok {
		return b
	}

	if s, ok := jtypes.AsString(v); ok {
		return s != ""
	}

	if n, ok := jtypes.AsNumber(v); ok {
		return n != 0
	}

	if jtypes.IsArray(v) {
		for i := 0; i < v.Len(); i++ {
			if Boolean(v.Index(i)) {
				return true
			}
		}
		return false
	}

	if jtypes.IsMap(v) {
		return v.Len() > 0
	}

	return false
}
```

- `Boolean` 函数接收一个 `reflect.Value` 类型的参数 `v`，并返回一个布尔值。
- `v = jtypes.Resolve(v)`：调用 `jtypes.Resolve` 函数处理 `v`，可能是为了解析指针或接口等类型。
- 接下来是一系列的类型检查和转换：
  - 如果 `v` 可以转换为布尔类型，则返回该布尔值。
  - 如果 `v` 可以转换为字符串类型，则返回该字符串是否非空。
  - 如果 `v` 可以转换为数字类型，则返回该数字是否非零。
  - 如果 `v` 是数组类型，则递归检查数组中的每个元素，只要有一个元素转换为布尔值为 `true`，则返回 `true`，否则返回 `false`。
  - 如果 `v` 是映射类型，则返回映射是否非空。
  - 如果以上条件都不满足，则返回 `false`。

### 函数 `Not`

```go
// Not (golint)
func Not(v reflect.Value) bool {
	return !Boolean(v)
}
```

- `Not` 函数接收一个 `reflect.Value` 类型的参数 `v`，并返回 `Boolean(v)` 的逻辑非结果。

### 函数 `Exists`

```go
// Exists (golint)
func Exists(v reflect.Value) bool {
	return v.IsValid()
}
```

- `Exists` 函数接收一个 `reflect.Value` 类型的参数 `v`，并返回 `v.IsValid()` 的结果，即检查 `v` 是否是一个有效的值。

### 总结

- `jlib` 包提供了三个主要函数：
  - `Boolean`：将任意类型的值转换为布尔值。
  - `Not`：返回 `Boolean` 函数的逻辑非结果。
  - `Exists`：检查一个值是否有效。
- 这些函数主要利用反射和 `jtypes` 包提供的辅助函数来处理不同类型的值。

通过阅读和理解这些代码，读者可以学习到如何使用反射来处理不同类型的值，并了解如何进行类型转换和检查。
  - date.go
    Summarized code for date.go

这个代码文件 `jlib` 主要提供了时间转换的功能，包括从毫秒数转换为时间字符串和从时间字符串转换为毫秒数。以下是对代码的详细解释：

### 包和导入

```go
package jlib

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/blues/jsonata-go/jlib/jxpath"
	"github.com/blues/jsonata-go/jtypes"
)
```

- `package jlib`：定义了包名为 `jlib`。
- `import`：导入了多个外部包，包括标准库中的 `fmt`、`regexp`、`strconv` 和 `time`，以及第三方包 `jxpath` 和 `jtypes`。

### 常量和变量

```go
const defaultFormatTimeLayout = "[Y]-[M01]-[D01]T[H01]:[m]:[s].[f001][Z01:01t]"

var defaultParseTimeLayouts = []string{
	"[Y]-[M01]-[D01]T[H01]:[m]:[s][Z01:01t]",
	"[Y]-[M01]-[D01]T[H01]:[m]:[s][Z0100t]",
	"[Y]-[M01]-[D01]T[H01]:[m]:[s]",
	"[Y]-[M01]-[D01]",
	"[Y]",
}
```

- `defaultFormatTimeLayout`：定义了一个默认的时间格式布局字符串。
- `defaultParseTimeLayouts`：定义了一个默认的时间解析布局字符串数组。

### 函数 `FromMillis`

```go
func FromMillis(ms int64, picture jtypes.OptionalString, tz jtypes.OptionalString) (string, error) {
	t := msToTime(ms).UTC()

	if tz.String != "" {
		loc, err := parseTimeZone(tz.String)
		if err != nil {
			return "", err
		}

		t = t.In(loc)
	}

	layout := picture.String
	if layout == "" {
		layout = defaultFormatTimeLayout
	}

	return jxpath.FormatTime(t, layout)
}
```

- `FromMillis`：将毫秒数转换为时间字符串。
  - `ms`：输入的毫秒数。
  - `picture`：可选的时间格式字符串。
  - `tz`：可选的时区字符串。
  - 首先将毫秒数转换为 `time.Time` 对象，并将其转换为 UTC 时间。
  - 如果提供了时区字符串，则解析时区并调整时间。
  - 使用提供的或默认的时间格式布局字符串，调用 `jxpath.FormatTime` 函数格式化时间。

### 函数 `parseTimeZone`

```go
func parseTimeZone(tz string) (*time.Location, error) {
	if len(tz) != 5 {
		return nil, fmt.Errorf("invalid timezone")
	}

	plusOrMinus := string(tz[0])

	var offsetMultiplier int
	switch plusOrMinus {
	case "-":
		offsetMultiplier = -1
	case "+":
		offsetMultiplier = 1
	default:
		return nil, fmt.Errorf("invalid timezone")
	}

	hours, err := strconv.Atoi(tz[1:3])
	if err != nil {
		return nil, fmt.Errorf("invalid timezone")
	}

	minutes, err := strconv.Atoi(tz[3:5])
	if err != nil {
		return nil, fmt.Errorf("invalid timezone")
	}

	offsetSeconds := offsetMultiplier * (60 * ((60 * hours) + minutes))

	loc := time.FixedZone(tz, offsetSeconds)

	return loc, nil
}
```

- `parseTimeZone`：解析 JSONata 时区字符串。
  - `tz`：输入的时区字符串。
  - 检查时区字符串的长度是否为 5。
  - 解析时区字符串的第一个字符，确定是加号还是减号。
  - 解析时区字符串的后四位，分别表示小时和分钟偏移量。
  - 计算总偏移秒数，并创建一个 `time.Location` 对象。

### 函数 `ToMillis`

```go
func ToMillis(s string, picture jtypes.OptionalString, tz jtypes.OptionalString) (int64, error) {
	layouts := defaultParseTimeLayouts
	if picture.String != "" {
		layouts = []string{picture.String}
	}

	for _, l := range layouts {
		if t, err := parseTime(s, l); err == nil {
			return timeToMS(t), nil
		}
	}

	return 0, fmt.Errorf("could not parse time %q", s)
}
```

- `ToMillis`：将时间字符串转换为毫秒数。
  - `s`：输入的时间字符串。
  - `picture`：可选的时间格式字符串。
  - `tz`：可选的时区字符串。
  - 使用提供的或默认的时间解析布局字符串数组，尝试解析时间字符串。
  - 如果解析成功，则将时间转换为毫秒数并返回。

### 函数 `parseTime`

```go
var reMinus7 = regexp.MustCompile("-(0*7)")

func parseTime(s string, picture string) (time.Time, error) {
	refTime := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*60*60))

	layout, err := jxpath.FormatTime(refTime, picture)
	if err != nil {
		return time.Time{}, fmt.Errorf("the second argument of the toMillis function must be a valid date format")
	}

	layout = reMinus7.ReplaceAllString(layout, "Z$1")

	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse time %q", s)
	}

	return t, nil
}
```

- `parseTime`：解析时间字符串。
  - `s`：输入的时间字符串。
  - `picture`：时间格式字符串。
  - 使用参考时间 `refTime` 和提供的格式字符串，生成一个布局字符串。
  - 使用正则表达式替换 `-07:00` 为 `Z07:00`。
  - 使用 `time.Parse` 函数解析时间字符串。

### 辅助函数

```go
func msToTime(ms int64) time.Time {
	return time.Unix(ms/1000, (ms%1000)*int64(time.Millisecond))
}

func timeToMS(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
```

- `msToTime`：将毫秒数转换为 `time.Time` 对象。
- `timeToMS`：将 `time.Time` 对象转换为毫秒数。

### 总结

这个代码文件 `jlib` 提供了时间转换的核心功能，包括从毫秒数转换为时间字符串和从时间字符串转换为毫秒数。它使用了多种时间格式布局字符串和时区解析逻辑，以支持灵活的时间转换需求。
  - date_test.go
    Summarized code for date_test.go

这段代码是一个Go语言的单元测试文件，用于测试`jlib`包中的`FromMillis`函数。该函数的主要功能是将一个表示毫秒时间戳的整数转换为特定格式的日期时间字符串。以下是对代码的详细解析：

### 文件头部信息
```go
// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib_test
```
- 版权声明和许可证信息。
- 声明了测试包`jlib_test`，通常用于对`jlib`包进行测试。

### 导入包
```go
import (
	"reflect"
	"testing"
	"time"

	"github.com/blues/jsonata-go/jlib"
	"github.com/blues/jsonata-go/jtypes"
)
```
- `reflect`：用于反射操作，这里用于设置`OptionalString`类型的值。
- `testing`：Go语言标准库中的测试包，用于编写和运行测试。
- `time`：用于处理日期和时间。
- `github.com/blues/jsonata-go/jlib`：被测试的包。
- `github.com/blues/jsonata-go/jtypes`：包含一些自定义类型，如`OptionalString`。

### 测试函数 `TestFromMillis`
```go
func TestFromMillis(t *testing.T) {
```
- 定义了一个名为`TestFromMillis`的测试函数，参数为`t *testing.T`，这是Go语言测试函数的典型签名。

### 初始化日期时间
```go
	date := time.Date(2018, time.September, 30, 15, 58, 5, int(762*time.Millisecond), time.UTC)
	input := date.UnixNano() / int64(time.Millisecond)
```
- 创建一个具体的日期时间对象`date`，表示2018年9月30日15时58分5秒762毫秒，时区为UTC。
- 将该日期时间对象转换为毫秒时间戳`input`。

### 测试用例
```go
	data := []struct {
		Picture       string
		TZ            string
		Output        string
		ExpectedError bool
	}{
		// 多个测试用例...
	}
```
- 定义了一个包含多个测试用例的切片`data`，每个测试用例包含以下字段：
  - `Picture`：日期时间格式字符串。
  - `TZ`：时区字符串。
  - `Output`：期望的输出结果。
  - `ExpectedError`：是否期望错误。

### 测试循环
```go
	for _, test := range data {
		var picture jtypes.OptionalString
		var tz jtypes.OptionalString

		if test.Picture != "" {
			picture.Set(reflect.ValueOf(test.Picture))
		}

		if test.TZ != "" {
			tz.Set(reflect.ValueOf(test.TZ))
		}

		got, err := jlib.FromMillis(input, picture, tz)

		if test.ExpectedError && err == nil {
			t.Errorf("%s: Expected error, got nil", test.Picture)
		} else if got != test.Output {
			t.Errorf("%s: Expected %q, got %q", test.Picture, test.Output, got)
		}
	}
```
- 遍历每个测试用例。
- 初始化`OptionalString`类型的变量`picture`和`tz`。
- 使用反射设置`picture`和`tz`的值。
- 调用`jlib.FromMillis`函数，传入毫秒时间戳`input`、格式字符串`picture`和时区`tz`。
- 检查返回结果`got`是否与期望的输出`test.Output`一致，以及是否符合期望的错误情况。

### 总结
该测试文件主要用于验证`jlib`包中的`FromMillis`函数是否能正确地将毫秒时间戳转换为特定格式的日期时间字符串。通过多个测试用例，涵盖了不同的日期时间格式和时区情况，确保函数的正确性和鲁棒性。
  - error.go
    Summarized code for error.go

这段代码定义了一个用于处理特定类型错误的包 `jlib`。以下是对代码的详细解释：

### 包声明
```go
package jlib
```
这行代码声明了当前代码文件属于 `jlib` 包。

### 导入依赖
```go
import "fmt"
```
这行代码导入了 `fmt` 包，用于格式化字符串和输出。

### 定义错误类型
```go
// ErrType (golint)
type ErrType uint
```
这里定义了一个名为 `ErrType` 的类型，它是 `uint` 类型的别名。这个类型用于表示不同的错误类型。

### 定义常量
```go
// ErrNanInf (golint)
const (
	_ ErrType = iota
	ErrNaNInf
)
```
这段代码使用 `iota` 定义了一个常量序列，其中 `_` 表示忽略第一个值，`ErrNaNInf` 表示第二个值。`iota` 是一个自增的枚举器，从 0 开始。这里 `ErrNaNInf` 的值为 1。

### 定义错误结构体
```go
// Error (golint)
type Error struct {
	Type ErrType
	Func string
}
```
这里定义了一个名为 `Error` 的结构体，包含两个字段：
- `Type`：表示错误的类型，类型为 `ErrType`。
- `Func`：表示发生错误的函数名，类型为 `string`。

### 实现 `error` 接口
```go
// Error (golint)
func (e Error) Error() string {
	var msg string

	switch e.Type {
	case ErrNaNInf:
		msg = "cannot convert NaN/Infinity to string"
	default:
		msg = "unknown error"
	}

	return fmt.Sprintf("%s: %s", e.Func, msg)
}
```
这段代码实现了 `error` 接口的 `Error` 方法。根据 `Error` 结构体中的 `Type` 字段，生成相应的错误信息字符串。具体逻辑如下：
- 如果 `Type` 是 `ErrNaNInf`，则错误信息为 `"cannot convert NaN/Infinity to string"`。
- 否则，错误信息为 `"unknown error"`。
- 最后，使用 `fmt.Sprintf` 将函数名和错误信息格式化为一个字符串并返回。

### 创建错误实例的函数
```go
func newError(name string, typ ErrType) *Error {
	return &Error{
		Func: name,
		Type: typ,
	}
}
```
这个函数用于创建一个新的 `Error` 实例。它接受两个参数：
- `name`：表示发生错误的函数名。
- `typ`：表示错误的类型。

函数返回一个指向新创建的 `Error` 实例的指针。

### 总结
这段代码定义了一个用于处理特定类型错误的包 `jlib`，主要功能包括：
- 定义了一个错误类型 `ErrType`。
- 定义了一个常量 `ErrNaNInf` 表示特定的错误类型。
- 定义了一个错误结构体 `Error`，包含错误类型和发生错误的函数名。
- 实现了 `error` 接口的 `Error` 方法，用于生成错误信息字符串。
- 提供了一个函数 `newError` 用于创建新的错误实例。

通过这些定义和实现，`jlib` 包可以方便地处理和表示特定类型的错误，并提供友好的错误信息。
  - hof.go
    Summarized code for hof.go

这个代码文件定义了一些高阶函数，如 `Map`、`Filter`、`Reduce` 和 `Single`，这些函数在处理数组或切片时非常有用。代码使用了反射机制来处理不同类型的输入，并调用了 `jtypes` 包中的一些类型和函数。以下是对每个函数的详细解释：

### 1. `Map` 函数

```go
func Map(v reflect.Value, f jtypes.Callable) (interface{}, error) {
	v = forceArray(jtypes.Resolve(v))
	var results []interface{}
	argc := clamp(f.ParamCount(), 1, 3)
	for i := 0; i < arrayLen(v); i++ {
		argv := []reflect.Value{v.Index(i), reflect.ValueOf(i), v}
		res, err := f.Call(argv[:argc])
		if err != nil {
			return nil, err
		}
		if res.IsValid() && res.CanInterface() {
			results = append(results, res.Interface())
		}
	}
	return results, nil
}
```

- **功能**：对数组或切片中的每个元素应用一个函数，并返回一个新的数组，其中包含应用函数后的结果。
- **实现细节**：
  - `v = forceArray(jtypes.Resolve(v))`：确保输入值是一个数组或切片。
  - `argc := clamp(f.ParamCount(), 1, 3)`：确定函数参数的数量，并限制在1到3之间。
  - 遍历数组中的每个元素，调用传入的函数 `f`，并将结果添加到 `results` 数组中。

### 2. `Filter` 函数

```go
func Filter(v reflect.Value, f jtypes.Callable) (interface{}, error) {
	v = forceArray(jtypes.Resolve(v))
	var results []interface{}
	argc := clamp(f.ParamCount(), 1, 3)
	for i := 0; i < arrayLen(v); i++ {
		item := v.Index(i)
		argv := []reflect.Value{item, reflect.ValueOf(i), v}
		res, err := f.Call(argv[:argc])
		if err != nil {
			return nil, err
		}
		if Boolean(res) && item.IsValid() && item.CanInterface() {
			results = append(results, item.Interface())
		}
	}
	return results, nil
}
```

- **功能**：根据传入的函数过滤数组或切片中的元素，返回一个新的数组，其中包含满足条件的元素。
- **实现细节**：
  - 与 `Map` 函数类似，但这里使用 `Boolean(res)` 来判断函数返回值是否为真，只有为真时才将元素添加到结果数组中。

### 3. `Reduce` 函数

```go
func Reduce(v reflect.Value, f jtypes.Callable, init jtypes.OptionalValue) (interface{}, error) {
	v = forceArray(jtypes.Resolve(v))
	var res reflect.Value
	if f.ParamCount() != 2 {
		return nil, fmt.Errorf("second argument of function \"reduce\" must be a function that takes two arguments")
	}
	i := 0
	switch {
	case init.IsSet():
		res = jtypes.Resolve(init.Value)
	case arrayLen(v) > 0:
		res = v.Index(0)
		i = 1
	}
	var err error
	for ; i < arrayLen(v); i++ {
		res, err = f.Call([]reflect.Value{res, v.Index(i)})
		if err != nil {
			return nil, err
		}
	}
	if !res.IsValid() || !res.CanInterface() {
		return nil, jtypes.ErrUndefined
	}
	return res.Interface(), nil
}
```

- **功能**：对数组或切片中的元素进行累积计算，返回一个单一的值。
- **实现细节**：
  - 确保传入的函数 `f` 有两个参数。
  - 根据是否有初始值 `init` 来确定初始的 `res` 值。
  - 遍历数组中的每个元素，调用传入的函数 `f`，并将结果累积到 `res` 中。

### 4. `Single` 函数

```go
func Single(v reflect.Value, f jtypes.Callable) (interface{}, error) {
	filteredValue, err := Filter(v, f)
	if err != nil {
		return nil, err
	}
	switch reflect.TypeOf(filteredValue).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(filteredValue)
		if s.Len() != 1 {
			return nil, fmt.Errorf("number of matching values returned by single() must be 1, got: %d", s.Len())
		}
		return s.Index(0).Interface(), nil
	default:
		return reflect.ValueOf(filteredValue).Interface(), nil
	}
}
```

- **功能**：从数组或切片中找到唯一一个满足条件的元素，如果满足条件的元素不止一个或没有，则返回错误。
- **实现细节**：
  - 使用 `Filter` 函数过滤出满足条件的元素。
  - 检查过滤后的结果是否为切片，如果是切片且长度不为1，则返回错误。
  - 如果过滤后的结果是单个值，则直接返回该值。

### 辅助函数 `clamp`

```go
func clamp(n, min, max int) int {
	switch {
	case n < min:
		return min
	case n > max:
		return max
	default:
		return n
	}
}
```

- **功能**：将一个整数 `n` 限制在 `min` 和 `max` 之间。
- **实现细节**：
  - 如果 `n` 小于 `min`，则返回 `min`。
  - 如果 `n` 大于 `max`，则返回 `max`。
  - 否则返回 `n`。

### 总结

这个代码文件通过反射机制实现了 `Map`、`Filter`、`Reduce` 和 `Single` 这些高阶函数，使得这些函数可以处理不同类型的数组或切片，并提供了灵活的参数处理和错误处理机制。这些函数在处理数据时非常有用，特别是在需要对数据进行转换、过滤和累积计算的场景中。
  - jlib.go
    Summarized code for jlib.go

这个代码文件定义了一个名为 `jlib` 的包，该包实现了 JSONata 函数库。以下是对代码的详细解释：

### 包声明和导入
```go
// Package jlib implements the JSONata function library.
package jlib

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/blues/jsonata-go/jtypes"
)
```
- **包声明**：声明了包名为 `jlib`，用于实现 JSONata 函数库。
- **导入包**：导入了多个标准库包和第三方包，包括 `fmt`、`math/rand`、`reflect`、`time` 以及 `github.com/blues/jsonata-go/jtypes`。

### 初始化函数
```go
func init() {
	// Seed random numbers for Random() and Shuffle().
	rand.Seed(time.Now().UnixNano())
}
```
- **初始化函数**：在包初始化时调用，用于为随机数生成器设置种子，确保每次运行时生成的随机数序列不同。

### 类型定义
```go
var typeBool = reflect.TypeOf((*bool)(nil)).Elem()
var typeCallable = reflect.TypeOf((*jtypes.Callable)(nil)).Elem()
var typeString = reflect.TypeOf((*string)(nil)).Elem()
var typeNumber = reflect.TypeOf((*float64)(nil)).Elem()
```
- **类型定义**：定义了一些常用的反射类型，用于后续的类型检查。

### 自定义类型和方法
```go
// StringNumberBool (golint)
type StringNumberBool reflect.Value

// ValidTypes (golint)
func (StringNumberBool) ValidTypes() []reflect.Type {
	return []reflect.Type{
		typeBool,
		typeString,
		typeNumber,
	}
}

// StringCallable (golint)
type StringCallable reflect.Value

// ValidTypes (golint)
func (StringCallable) ValidTypes() []reflect.Type {
	return []reflect.Type{
		typeString,
		typeCallable,
	}
}

func (s StringCallable) toInterface() interface{} {
	if v := reflect.Value(s); v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return nil
}
```
- **自定义类型**：定义了 `StringNumberBool` 和 `StringCallable` 类型，它们都是 `reflect.Value` 的别名。
- **方法**：为这两个类型定义了 `ValidTypes` 方法，返回它们支持的有效类型。`StringCallable` 类型还定义了一个 `toInterface` 方法，用于将 `reflect.Value` 转换为 `interface{}`。

### TypeOf 函数
```go
// TypeOf implements the jsonata $type function that returns the data type of
// the argument
func TypeOf(x interface{}) (string, error) {
	v := reflect.ValueOf(x)
	if jtypes.IsCallable(v) {
		return "function", nil
	}
	if jtypes.IsString(v) {
		return "string", nil
	}
	if jtypes.IsNumber(v) {
		return "number", nil
	}
	if jtypes.IsArray(v) {
		return "array", nil
	}
	if jtypes.IsBool(v) {
		return "boolean", nil
	}
	if jtypes.IsMap(v) {
		return "object", nil
	}

	switch x.(type) {
	case *interface{}:
		return "null", nil
	}

	xType := reflect.TypeOf(x).String()
	return "", fmt.Errorf("unknown type %s", xType)
}
```
- **TypeOf 函数**：实现了 JSONata 的 `$type` 函数，用于返回参数的数据类型。
- **类型检查**：使用反射检查参数的类型，并返回相应的字符串表示。如果类型未知，则返回错误。

### 总结
这个包主要用于实现 JSONata 函数库，提供了类型检查和反射操作的功能。通过定义自定义类型和方法，以及使用反射库，实现了对不同数据类型的识别和处理。
  - jxpath
    - formatdate.go
      [File too long to summarize]
    - formatdate_test.go
      Summarized code for formatdate_test.go

这个代码文件是一个用于测试时间格式化的Go语言测试文件，主要测试了年份、时区、星期几的格式化功能。以下是对代码文件内容的详细解释：

### 文件信息
- **版权声明**：代码由Blues Inc. 所有，使用该源代码需遵守相关许可。

### 包声明
- **包名**：`jxpath`，表示这个包是用于处理JXPath格式化的。

### 导入的包
- `reflect`：用于深度比较对象。
- `testing`：用于编写测试用例。
- `time`：用于处理时间相关的操作。

### 测试函数

#### `TestFormatYear`
- **功能**：测试年份的格式化。
- **实现细节**：
  - 定义了一个输入时间 `input`，为2018年4月1日12:00:00 UTC。
  - 定义了一个测试数据结构 `data`，包含不同的格式化字符串 `Picture` 和预期的输出 `Output`。
  - 使用循环遍历每个测试用例，调用 `FormatTime` 函数进行格式化，并比较结果和预期输出。

#### `TestFormatTimezone`
- **功能**：测试时区的格式化。
- **实现细节**：
  - 定义了一些时区信息 `timezones`，包括时区名称和偏移量。
  - 创建了一个时间数组 `times`，每个时间对应一个时区。
  - 定义了一个测试数据结构 `data`，包含不同的格式化字符串 `Picture` 和预期的输出数组 `Outputs`。
  - 使用嵌套循环遍历每个测试用例和每个时间，调用 `FormatTime` 函数进行格式化，并比较结果和预期输出。

#### `TestFormatDayOfWeek`
- **功能**：测试星期几的格式化。
- **实现细节**：
  - 定义了一个起始时间 `startTime`，为2018年4月1日12:00:00 UTC。
  - 创建了一个时间数组 `times`，包含从起始时间开始的7天。
  - 定义了一个测试数据结构 `data`，包含不同的格式化字符串 `Picture` 和预期的输出数组 `Outputs`。
  - 使用嵌套循环遍历每个测试用例和每个时间，调用 `FormatTime` 函数进行格式化，并比较结果和预期输出。

### 总结
这个代码文件通过编写详细的测试用例，确保了 `FormatTime` 函数在处理年份、时区和星期几的格式化时能够正确工作。每个测试函数都包含了多个测试用例，覆盖了不同的格式化需求和边界条件。通过这些测试，可以确保代码的正确性和健壮性。
    - formatnumber.go
      [File too long to summarize]
    - formatnumber_test.go
      Summarized code for formatnumber_test.go

这个代码文件主要用于测试一个名为 `FormatNumber` 的函数，该函数的功能是将浮点数按照指定的格式字符串进行格式化。以下是对代码文件的详细解析：

### 1. 版权声明
```go
// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.
```
这段注释声明了代码的版权信息和使用许可。

### 2. 包声明
```go
package jxpath
```
声明了代码所在的包名为 `jxpath`。

### 3. 导入依赖包
```go
import (
	"reflect"
	"testing"
)
```
导入了 `reflect` 和 `testing` 包，分别用于深度比较和单元测试。

### 4. 测试数据结构定义
```go
type formatNumberTest struct {
	Value   float64
	Picture string
	Output  string
	Error   error
}
```
定义了一个结构体 `formatNumberTest`，用于存储测试用例的数据。每个测试用例包含一个浮点数值 `Value`、一个格式字符串 `Picture`、预期的输出字符串 `Output` 和一个预期的错误 `Error`。

### 5. 测试函数 `TestExamples`
```go
func TestExamples(t *testing.T) {
	tests := []formatNumberTest{
		// 省略具体测试用例
	}

	testFormatNumber(t, tests)
}
```
定义了一个测试函数 `TestExamples`，该函数初始化了一系列测试用例，并调用 `testFormatNumber` 函数来执行这些测试。

### 6. 测试用例
```go
tests := []formatNumberTest{
	{
		Value:   12345.6,
		Picture: "#,###.00",
		Output:  "12,345.60",
	},
	// 省略其他测试用例
}
```
定义了一系列测试用例，每个用例包含一个浮点数值、一个格式字符串和预期的输出字符串。

### 7. 执行测试的函数 `testFormatNumber`
```go
func testFormatNumber(t *testing.T, tests []formatNumberTest) {
	df := NewDecimalFormat()

	for i, test := range tests {
		output, err := FormatNumber(test.Value, test.Picture, df)

		if output != test.Output {
			t.Errorf("%d. FormatNumber(%v, %q): expected %s, got %s", i+1, test.Value, test.Picture, test.Output, output)
		}

		if !reflect.DeepEqual(err, test.Error) {
			t.Errorf("%d. FormatNumber(%v, %q): expected error %v, got %v", i+1, test.Value, test.Picture, test.Error, err)
		}
	}
}
```
`testFormatNumber` 函数接收一个 `testing.T` 对象和一个测试用例切片。它首先创建一个 `DecimalFormat` 对象 `df`，然后遍历每个测试用例，调用 `FormatNumber` 函数进行格式化，并检查输出和错误是否与预期一致。如果不一致，则使用 `t.Errorf` 输出错误信息。

### 总结
这个代码文件通过定义一系列测试用例，使用 `FormatNumber` 函数对浮点数进行格式化，并验证输出和错误是否符合预期。通过这种方式，可以确保 `FormatNumber` 函数的正确性和稳定性。
    - language.go
      Summarized code for language.go

这段代码定义了一个用于处理日期和时间格式的包 `jxpath`，并提供了一个多语言支持的日期和时间格式化功能。以下是对代码的详细解释：

### 1. 包声明
```go
package jxpath
```
这行代码声明了当前代码文件属于 `jxpath` 包。

### 2. 导入依赖
```go
import (
	"time"
)
```
导入了 `time` 包，用于处理日期和时间相关的操作。

### 3. 定义 `dateLanguage` 结构体
```go
type dateLanguage struct {
	days     [7][]string
	months   [13][]string
	am       []string
	pm       []string
	tzPrefix string
}
```
- `days`: 一个包含7个元素的数组，每个元素是一个字符串切片，表示一周中每一天的不同名称形式。
- `months`: 一个包含13个元素的数组，每个元素是一个字符串切片，表示一年中每个月的不同名称形式。
- `am`: 一个字符串切片，表示上午的时间前缀。
- `pm`: 一个字符串切片，表示下午的时间前缀。
- `tzPrefix`: 一个字符串，表示时区的前缀。

### 4. 定义 `dateLanguages` 映射
```go
var dateLanguages = map[string]dateLanguage{
	"en": {
		days: [...][]string{
			time.Sunday: {
				"Sunday",
				"Sun",
				"Su",
			},
			time.Monday: {
				"Monday",
				"Mon",
				"Mo",
			},
			time.Tuesday: {
				"Tuesday",
				"Tues",
				"Tue",
				"Tu",
			},
			time.Wednesday: {
				"Wednesday",
				"Weds",
				"Wed",
				"We",
			},
			time.Thursday: {
				"Thursday",
				"Thurs",
				"Thur",
				"Thu",
				"Th",
			},
			time.Friday: {
				"Friday",
				"Fri",
				"Fr",
			},
			time.Saturday: {
				"Saturday",
				"Sat",
				"Sa",
			},
		},
		months: [...][]string{
			time.January: {
				"January",
				"Jan",
				"Ja",
			},
			time.February: {
				"February",
				"Feb",
				"Fe",
			},
			time.March: {
				"March",
				"Mar",
				"Mr",
			},
			time.April: {
				"April",
				"Apr",
				"Ap",
			},
			time.May: {
				"May",
				"My",
			},
			time.June: {
				"June",
				"Jun",
				"Jn",
			},
			time.July: {
				"July",
				"Jul",
				"Jl",
			},
			time.August: {
				"August",
				"Aug",
				"Au",
			},
			time.September: {
				"September",
				"Sept",
				"Sep",
				"Se",
			},
			time.October: {
				"October",
				"Oct",
				"Oc",
			},
			time.November: {
				"November",
				"Nov",
				"No",
			},
			time.December: {
				"December",
				"Dec",
				"De",
			},
		},
		am: []string{
			"am",
			"a",
		},
		pm: []string{
			"pm",
			"p",
		},
		tzPrefix: "GMT",
	},
}
```
- `dateLanguages` 是一个映射，键是语言代码（如 "en" 表示英语），值是 `dateLanguage` 结构体。
- 当前只定义了英语（"en"）的语言配置。

### 5. 定义默认语言
```go
var defaultLanguage = dateLanguages["en"]
```
这行代码将默认语言设置为英语。

### 总结
这段代码的主要功能是提供一个多语言支持的日期和时间格式化功能。通过定义 `dateLanguage` 结构体和 `dateLanguages` 映射，可以方便地扩展支持更多语言的日期和时间格式。当前只实现了英语的配置，但通过添加更多的语言代码和对应的 `dateLanguage` 结构体，可以轻松扩展到其他语言。
  - number.go
    Summarized code for number.go

这个代码文件定义了一个名为 `jlib` 的 Go 包，主要提供了一些与数值处理相关的功能。以下是对每个函数和实现细节的详细解释：

### 1. `Number` 函数

```go
// Number converts values to numbers. Numeric values are returned
// unchanged. Strings in legal JSON number format are converted
// to the number they represent. Boooleans are converted to 0 or 1.
// All other types trigger an error.
func Number(value StringNumberBool) (float64, error) {
	v := reflect.Value(value)
	if b, ok := jtypes.AsBool(v); ok {
		if b {
			return 1, nil
		}
		return 0, nil
	}

	if n, ok := jtypes.AsNumber(v); ok {
		return n, nil
	}

	s, ok := jtypes.AsString(v)
	if ok && reNumber.MatchString(s) {
		if n, err := strconv.ParseFloat(s, 64); err == nil {
			return n, nil
		}
	}

	return 0, fmt.Errorf("unable to cast %q to a number", s)
}
```

- **功能**：将不同类型的值转换为浮点数。
- **实现细节**：
  - 使用反射 (`reflect.Value`) 处理输入值。
  - 如果值是布尔类型，返回 `1` 或 `0`。
  - 如果值是数字类型，直接返回该数字。
  - 如果值是字符串类型，并且符合 JSON 数字格式，将其转换为浮点数。
  - 其他类型会触发错误。

### 2. `Round` 函数

```go
// Round rounds its input to the number of decimal places given
// in the optional second parameter. By default, Round rounds to
// the nearest integer. A negative precision specifies which column
// to round to on the left hand side of the decimal place.
func Round(x float64, prec jtypes.OptionalInt) float64 {
	// Adapted from gonum's floats.RoundEven.
	// https://github.com/gonum/gonum/tree/master/floats

	if x == 0 {
		// Make sure zero is returned
		// without the negative bit set.
		return 0
	}
	// Fast path for positive precision on integers.
	if prec.Int >= 0 && x == math.Trunc(x) {
		return x
	}
	intermed := multByPow10(x, prec.Int)
	if math.IsInf(intermed, 0) {
		return x
	}
	if isHalfway(intermed) {
		correction, _ := math.Modf(math.Mod(intermed, 2))
		intermed += correction
		if intermed > 0 {
			x = math.Floor(intermed)
		} else {
			x = math.Ceil(intermed)
		}
	} else {
		if x < 0 {
			x = math.Ceil(intermed - 0.5)
		} else {
			x = math.Floor(intermed + 0.5)
		}
	}

	if x == 0 {
		return 0
	}

	return multByPow10(x, -prec.Int)
}
```

- **功能**：将输入值四舍五入到指定的小数位数。
- **实现细节**：
  - 如果输入值为零，直接返回零。
  - 如果输入值是整数且精度非负，直接返回该值。
  - 使用 `multByPow10` 函数将输入值乘以 10 的指定次方。
  - 根据中间值是否为半值进行四舍五入。
  - 最后将结果除以 10 的指定次方，返回四舍五入后的值。

### 3. `Power` 函数

```go
// Power returns x to the power of y.
func Power(x, y float64) (float64, error) {
	res := math.Pow(x, y)
	if math.IsInf(res, 0) || math.IsNaN(res) {
		return 0, fmt.Errorf("the power function has resulted in a value that cannot be represented as a JSON number")
	}
	return res, nil
}
```

- **功能**：计算 `x` 的 `y` 次方。
- **实现细节**：
  - 使用 `math.Pow` 计算幂。
  - 如果结果是无限大或 NaN，返回错误。

### 4. `Sqrt` 函数

```go
// Sqrt returns the square root of a number. It returns an error
// if the number is less than zero.
func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, fmt.Errorf("the sqrt function cannot be applied to a negative number")
	}
	return math.Sqrt(x), nil
}
```

- **功能**：计算一个数的平方根。
- **实现细节**：
  - 如果输入值小于零，返回错误。
  - 否则使用 `math.Sqrt` 计算平方根。

### 5. `Random` 函数

```go
// Random returns a random floating point number between 0 and 1.
func Random() float64 {
	return rand.Float64()
}
```

- **功能**：返回一个介于 0 和 1 之间的随机浮点数。
- **实现细节**：
  - 使用 `rand.Float64` 生成随机数。

### 6. `multByPow10` 函数

```go
// multByPow10 multiplies a number by 10 to the power of n.
// It does this by converting back and forth to strings to
// avoid floating point rounding errors, e.g.
//
//     4.525 * math.Pow10(2) returns 452.50000000000006
func multByPow10(x float64, n int) float64 {
	if n == 0 || math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}

	s := fmt.Sprintf("%g", x)

	chunks := strings.Split(s, "e")
	switch len(chunks) {
	case 1:
		s = chunks[0] + "e" + strconv.Itoa(n)
	case 2:
		e, _ := strconv.Atoi(chunks[1])
		s = chunks[0] + "e" + strconv.Itoa(e+n)
	default:
		return x
	}

	x, _ = strconv.ParseFloat(s, 64)
	return x
}
```

- **功能**：将一个数乘以 10 的指定次方。
- **实现细节**：
  - 使用字符串格式化和解析来避免浮点数舍入误差。
  - 将输入值转换为科学计数法表示的字符串，然后调整指数部分。
  - 最后将字符串解析回浮点数。

### 7. `isHalfway` 函数

```go
func isHalfway(x float64) bool {
	_, frac := math.Modf(x)
	frac = math.Abs(frac)
	return frac == 0.5 || (math.Nextafter(frac, math.Inf(-1)) < 0.5 && math.Nextafter(frac, math.Inf(1)) > 0.5)
}
```

- **功能**：判断一个数是否为半值。
- **实现细节**：
  - 使用 `math.Modf` 获取小数部分。
  - 判断小数部分是否等于 0.5 或在其附近。

### 总结

这个代码文件主要提供了数值转换、四舍五入、幂运算、平方根计算和随机数生成等功能。通过使用反射、字符串处理和数学函数，实现了对不同类型输入值的处理和数值运算的精确控制。
  - number_test.go
    Summarized code for number_test.go

这段代码是一个Go语言的单元测试文件，主要用于测试`jlib`包中的`Round`函数。以下是对代码的详细解释：

### 文件头部信息
```go
// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.
```
这部分是版权声明，表明代码的版权归属和使用许可。

### 包声明和导入
```go
package jlib_test

import (
	"fmt"
	"testing"

	"github.com/blues/jsonata-go/jlib"
	"github.com/blues/jsonata-go/jtypes"
)
```
- `package jlib_test`：声明这是一个测试包，通常用于测试主包`jlib`。
- `import`：导入了多个包，包括标准库的`fmt`和`testing`，以及第三方库`github.com/blues/jsonata-go/jlib`和`github.com/blues/jsonata-go/jtypes`。

### 测试函数
```go
func TestRound(t *testing.T) {
```
定义了一个名为`TestRound`的测试函数，用于测试`jlib`包中的`Round`函数。

### 测试数据
```go
	data := []struct {
		Value     float64
		Precision jtypes.OptionalInt
		Output    float64
	}{
		// 省略具体数据
	}
```
定义了一个结构体切片`data`，包含多个测试用例。每个测试用例包含三个字段：
- `Value`：要进行四舍五入的浮点数。
- `Precision`：精度，类型为`jtypes.OptionalInt`，表示可选的整数精度。
- `Output`：期望的四舍五入结果。

### 测试用例
```go
	for _, test := range data {
		got := jlib.Round(test.Value, test.Precision)
		if got != test.Output {
			s := fmt.Sprintf("round(%g", test.Value)
			if test.Precision.IsSet() {
				s += fmt.Sprintf(", %d", test.Precision.Int)
			}
			s += ")"
			t.Errorf("%s: Expected %g, got %g", s, test.Output, got)
		}
	}
```
- 使用`for`循环遍历`data`切片中的每个测试用例。
- 调用`jlib.Round`函数，传入`Value`和`Precision`，得到四舍五入的结果`got`。
- 如果`got`与期望的`Output`不一致，则构建错误信息字符串`s`，并使用`t.Errorf`输出错误信息。

### 具体测试用例
```go
		{
			Value:  11.5,
			Output: 12,
		},
		{
			Value:  -11.5,
			Output: -12,
		},
		// 省略其他用例
```
每个测试用例都包含一个输入值`Value`和一个期望的输出值`Output`，部分用例还包含一个精度`Precision`。

### 总结
这段代码通过多个测试用例验证了`jlib`包中的`Round`函数的正确性，确保该函数能够按照预期进行四舍五入操作。通过这些测试用例，可以覆盖不同情况下的四舍五入行为，包括正负数、不同精度等。
  - object.go
    Summarized code for object.go

这个代码文件定义了一个名为 `jlib` 的 Go 包，主要用于处理 JSON 对象和结构体的高级操作。以下是对代码文件的详细解析：

### 1. 包声明和导入
```go
package jlib

import (
	"fmt"
	"reflect"

	"github.com/blues/jsonata-go/jtypes"
)
```
- 包名为 `jlib`。
- 导入了 `fmt`、`reflect` 和 `github.com/blues/jsonata-go/jtypes` 包。

### 2. 常量和变量
```go
var typeInterfaceMap = reflect.MapOf(typeString, jtypes.TypeInterface)
```
- `typeInterfaceMap` 是一个 `reflect.Type` 类型的变量，表示 `map[string]interface{}` 类型。

### 3. 函数 `toInterfaceMap`
```go
func toInterfaceMap(v reflect.Value) (map[string]interface{}, bool) {
	if v.Type() == typeInterfaceMap && v.CanInterface() {
		return v.Interface().(map[string]interface{}), true
	}
	return nil, false
}
```
- 尝试将 `reflect.Value` 转换为 `map[string]interface{}` 类型。
- 如果成功，返回转换后的 map 和 `true`；否则返回 `nil` 和 `false`。

### 4. 函数 `Each`
```go
func Each(obj reflect.Value, fn jtypes.Callable) (interface{}, error) {
	var each func(reflect.Value, jtypes.Callable) ([]interface{}, error)

	obj = jtypes.Resolve(obj)

	switch {
	case jtypes.IsMap(obj):
		each = eachMap
	case jtypes.IsStruct(obj) && !jtypes.IsCallable(obj):
		each = eachStruct
	default:
		return nil, fmt.Errorf("argument must be an object")
	}

	if argc := fn.ParamCount(); argc < 1 || argc > 3 {
		return nil, fmt.Errorf("function must take 1, 2 or 3 arguments")
	}

	results, err := each(obj, fn)
	if err != nil {
		return nil, err
	}

	switch len(results) {
	case 0:
		return nil, jtypes.ErrUndefined
	case 1:
		return results[0], nil
	default:
		return results, nil
	}
}
```
- 对对象中的每个键值对应用函数 `fn`，并返回结果数组。
- 支持 map 和 struct 类型。
- 函数 `fn` 的参数数量必须在 1 到 3 之间。

### 5. 辅助函数 `eachMap` 和 `eachStruct`
```go
func eachMap(v reflect.Value, fn jtypes.Callable) ([]interface{}, error) {
	size := v.Len()
	if size == 0 {
		return nil, nil
	}

	var results []interface{}
	argv := make([]reflect.Value, fn.ParamCount())

	for _, k := range v.MapKeys() {
		for i := range argv {
			switch i {
			case 0:
				argv[i] = v.MapIndex(k)
			case 1:
				argv[i] = k
			case 2:
				argv[i] = v
			}
		}

		res, err := fn.Call(argv)
		if err != nil {
			return nil, err
		}

		if res.IsValid() && res.CanInterface() {
			if results == nil {
				results = make([]interface{}, 0, size)
			}
			results = append(results, res.Interface())
		}
	}

	return results, nil
}

func eachStruct(v reflect.Value, fn jtypes.Callable) ([]interface{}, error) {
	size := v.NumField()
	if size == 0 {
		return nil, nil
	}

	var results []interface{}
	t := v.Type()
	argv := make([]reflect.Value, fn.ParamCount())

	for i := 0; i < size; i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		for j := range argv {
			switch j {
			case 0:
				argv[j] = v.Field(i)
			case 1:
				argv[j] = reflect.ValueOf(field.Name)
			case 2:
				argv[j] = v
			}
		}

		res, err := fn.Call(argv)
		if err != nil {
			return nil, err
		}

		if res.IsValid() && res.CanInterface() {
			if results == nil {
				results = make([]interface{}, 0, size)
			}
			results = append(results, res.Interface())
		}
	}

	return results, nil
}
```
- `eachMap` 和 `eachStruct` 分别处理 map 和 struct 类型的对象。
- 对每个键值对调用函数 `fn`，并将结果收集到数组中。

### 6. 函数 `Sift`
```go
func Sift(obj reflect.Value, fn jtypes.Callable) (interface{}, error) {
	var sift func(reflect.Value, jtypes.Callable) (map[string]interface{}, error)

	obj = jtypes.Resolve(obj)

	switch {
	case jtypes.IsMap(obj):
		sift = siftMap
	case jtypes.IsStruct(obj) && !jtypes.IsCallable(obj):
		sift = siftStruct
	default:
		return nil, fmt.Errorf("argument must be an object")
	}

	if argc := fn.ParamCount(); argc < 1 || argc > 3 {
		return nil, fmt.Errorf("function must take 1, 2 or 3 arguments")
	}

	results, err := sift(obj, fn)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, jtypes.ErrUndefined
	}

	return results, nil
}
```
- 从对象中筛选出满足条件的键值对，并返回新的 map。
- 支持 map 和 struct 类型。
- 函数 `fn` 的参数数量必须在 1 到 3 之间。

### 7. 辅助函数 `siftMap` 和 `siftStruct`
```go
func siftMap(v reflect.Value, fn jtypes.Callable) (map[string]interface{}, error) {
	size := v.Len()
	if size == 0 {
		return nil, nil
	}

	var results map[string]interface{}
	argv := make([]reflect.Value, fn.ParamCount())

	for _, k := range v.MapKeys() {
		key, ok := jtypes.AsString(k)
		if !ok {
			return nil, fmt.Errorf("object key must evaluate to a string, got %v (%s)", k, k.Kind())
		}

		val := v.MapIndex(k)
		if !val.IsValid() || !val.CanInterface() {
			continue
		}

		for i := range argv {
			switch i {
			case 0:
				argv[i] = val
			case 1:
				argv[i] = k
			case 2:
				argv[i] = v
			}
		}

		res, err := fn.Call(argv)
		if err != nil {
			return nil, err
		}

		if Boolean(res) {
			if results == nil {
				results = make(map[string]interface{}, size)
			}
			results[key] = val.Interface()
		}
	}

	return results, nil
}

func siftStruct(v reflect.Value, fn jtypes.Callable) (map[string]interface{}, error) {
	size := v.NumField()
	if size == 0 {
		return nil, nil
	}

	var results map[string]interface{}
	t := v.Type()
	argv := make([]reflect.Value, fn.ParamCount())

	for i := 0; i < size; i++ {
		key := t.Field(i).Name
		val := v.Field(i)
		if !val.IsValid() || !val.CanInterface() {
			continue
		}

		for j := range argv {
			switch j {
			case 0:
				argv[j] = val
			case 1:
				argv[j] = reflect.ValueOf(key)
			case 2:
				argv[j] = v
			}
		}

		res, err := fn.Call(argv)
		if err != nil {
			return nil, err
		}

		if Boolean(res) {
			if results == nil {
				results = make(map[string]interface{}, size)
			}
			results[key] = val.Interface()
		}
	}

	return results, nil
}
```
- `siftMap` 和 `siftStruct` 分别处理 map 和 struct 类型的对象。
- 对每个键值对调用函数 `fn`，并将满足条件的结果收集到新的 map 中。

### 8. 函数 `Keys`
```go
func Keys(obj reflect.Value) (interface{}, error) {
	results, err := keys(obj)
	if err != nil {
		return nil, err
	}

	switch len(results) {
	case 0:
		return nil, jtypes.ErrUndefined
	case 1:
		return results[0], nil
	default:
		return results, nil
	}
}
```
- 返回对象的键名数组。
- 支持 map、struct 和 array 类型。

### 9. 辅助函数 `keys`、`keysMap`、`keysMapFast`、`keysStruct` 和 `keysArray`
```go
func keys(v reflect.Value) ([]string, error) {
	v = jtypes.Resolve(v)

	switch {
	case jtypes.IsMap(v):
		return keysMap(v)
	case jtypes.IsStruct(v) && !jtypes.IsCallable(v):
		return keysStruct(v)
	case jtypes.IsArray(v):
		return keysArray(v)
	default:
		return nil, nil
	}
}

func keysMap(v reflect.Value) ([]string, error) {
	if v.Len() == 0 {
		return nil, nil
	}

	if m, ok := toInterfaceMap(v); ok {
		return keysMapFast(m), nil
	}

	results := make([]string, v.Len())

	for i, k := range v.MapKeys() {
		key, ok := jtypes.AsString(k)
		if !ok {
			return nil, fmt.Errorf("object key must evaluate to a string, got %v (%s)", k, k.Kind())
		}

		results[i] = key
	}

	return results, nil
}

func keysMapFast(m map[string]interface{}) []string {
	results := make([]string, 0, len(m))
	for key := range m {
		results = append(results, key)
	}

	return results
}

func keysStruct(v reflect.Value) ([]string, error) {
	size := v.NumField()
	if size == 0 {
		return nil, nil
	}

	var results []string
	t := v.Type()
	for i := 0; i < size; i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		if results == nil {
			results = make([]string, 0, size)
		}
		results = append(results, field.Name)
	}

	return results, nil
}

func keysArray(v reflect.Value) ([]string, error) {
	size := v.Len()
	if size == 0 {
		return nil, nil
	}

	kresults := make([][]string, 0, size)

	for i := 0; i < size; i++ {
		results, err := keys(v.Index(i))
		if err != nil {
			return nil, err
		}
		kresults = append(kresults, results)
	}

	size = 0
	for _, k := range kresults {
		size += len(k)
	}

	if size == 0 {
		return nil, nil
	}

	seen := map[string]bool{}
	results := make([]string, 0, size)

	for _, k := range kresults {
		for _, s := range k {
			if !seen[s] {
				seen[s] = true
				results = append(results, s)
			}
		}
	}

	return results, nil
}
```
- `keys` 根据对象类型调用相应的辅助函数。
- `keysMap` 和 `keysMapFast` 处理 map 类型的对象。
- `keysStruct` 处理 struct 类型的对象。
- `keysArray` 处理 array 类型的对象，返回所有对象的键名集合。

### 10. 函数 `Merge`
```go
func Merge(objs reflect.Value) (interface{}, error) {
	var size int
	var merge func(map[string]interface{}, reflect.Value) error

	objs = jtypes.Resolve(objs)

	switch {
	case jtypes.IsMap(objs):
		size = objs.Len()
		merge = mergeMap
	case jtypes.IsStruct(objs) && !jtypes.IsCallable(objs):
		size = objs.NumField()
		merge = mergeStruct
	case jtypes.IsArray(objs):
		for i := 0; i < objs.Len(); i++ {
			obj := jtypes.Resolve(objs.Index(i))
			switch {
			case jtypes.IsMap(obj):
				size += obj.Len()
			case jtypes.IsStruct(obj):
				size += obj.NumField()
			default:
				return nil, fmt.Errorf("argument must be an object or an array of objects")
			}
		}
		merge = mergeArray
	default:
		return nil, fmt.Errorf("argument must be an object or an array of objects")
	}

	results := make(map[string]interface{}, size)
	if err := merge(results, objs); err != nil {
		return nil, err
	}

	return results, nil
}
```
- 将多个对象合并成一个对象。
- 支持 map、struct 和 array 类型。

### 11. 辅助函数 `mergeMap`、`mergeMapFast`、`mergeStruct` 和 `mergeArray`
```go
func mergeMap(dest map[string]interface{}, src reflect.Value) error {
	if m, ok := toInterfaceMap(src); ok {
		mergeMapFast(dest, m)
		return nil
	}

	for _, k := range src.MapKeys() {
		key, ok := jtypes.AsString(k)
		if !ok {
			return fmt.Errorf("object key must evaluate to a string, got %v (%s)", k, k.Kind())
		}

		if val := src.MapIndex(k); val.IsValid() && val.CanInterface() {
			dest[key] = val.Interface()
		}
	}

	return nil
}

func mergeMapFast(dest, src map[string]interface{}) {
	for k, v := range src {
		if v != nil {
			dest[k] = v
		}
	}
}

func mergeStruct(dest map[string]interface{}, src reflect.Value) error {
	t := src.Type()

	for i := 0; i < src.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}

		if val := src.Field(i); val.IsValid() && val.CanInterface() {
			dest[field.Name] = val.Interface()
		}
	}

	return nil
}

func mergeArray(dest map[string]interface{}, src reflect.Value) error {
	var merge func(map[string]interface{}, reflect.Value) error

	for i := 0; i < src.Len(); i++ {
		item := jtypes.Resolve(src.Index(i))

		switch {
		case jtypes.IsMap(item):
			merge = mergeMap
		case jtypes.IsStruct(item) && !jtypes.IsCallable(item):
			merge = mergeStruct
		default:
			continue
		}

		if err := merge(dest, item); err != nil {
			return err
		}
	}

	return nil
}
```
- `mergeMap` 和 `mergeMapFast` 处理 map 类型的对象。
- `
  - object_test.go
    [File too long to summarize]
  - string.go
    [File too long to summarize]
  - string_test.go
    [File too long to summarize]
- jparse
  - doc.go
    Summarized code for doc.go

### 代码文件内容总结

#### 文件信息
- **版权声明**：该代码文件的版权归 Blues Inc. 所有，使用该源代码需遵守版权持有者授予的许可，具体许可内容见 `LICENSE` 文件。

#### 包描述
- **包名**：`jparse`
- **功能**：该包用于将 JSONata 表达式转换为抽象语法树（Abstract Syntax Tree, AST）。大多数客户端不需要直接使用此包。

#### 使用说明
- **函数调用**：通过调用 `Parse` 函数，传入一个 JSONata 表达式字符串。
- **返回值**：如果解析过程中发生错误，将返回一个 `Error` 类型的错误；否则，返回 AST 的根节点 `Node`。

### 功能和实现细节

#### 功能概述
- **JSONata 表达式解析**：该包的主要功能是将 JSONata 表达式解析成抽象语法树（AST）。JSONata 是一种用于查询和转换 JSON 数据的强大表达式语言。
- **抽象语法树（AST）**：AST 是一种树状结构，用于表示源代码的语法结构。通过将 JSONata 表达式转换为 AST，可以更方便地进行后续的分析和处理。

#### 实现细节
- **入口函数**：`Parse` 函数是该包的入口点，负责接收 JSONata 表达式字符串并进行解析。
- **错误处理**：在解析过程中，如果遇到错误，将返回一个 `Error` 类型的错误。这种错误类型可能包含了详细的错误信息，如错误位置、错误原因等。
- **AST 节点**：解析成功后，返回的根节点 `Node` 是 AST 的起点。通过遍历 `Node` 及其子节点，可以访问和操作整个 AST。

### 代码示例
以下是一个简单的代码示例，展示如何使用 `jparse` 包解析 JSONata 表达式：

```go
package main

import (
    "fmt"
    "jparse"
)

func main() {
    expr := `{"name": $string(firstName & " " & lastName)}`
    ast, err := jparse.Parse(expr)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("AST Root Node:", ast)
}
```

### 总结
- **包名**：`jparse`
- **主要功能**：将 JSONata 表达式转换为抽象语法树（AST）。
- **使用方法**：调用 `Parse` 函数，传入 JSONata 表达式字符串，返回 AST 的根节点或错误。
- **错误处理**：解析过程中发生错误时，返回 `Error` 类型的错误。
- **AST 节点**：解析成功后，返回的 `Node` 是 AST 的根节点，用于后续的分析和处理。

通过阅读该代码文件的注释和描述，读者可以了解到该包的主要功能、使用方法以及错误处理机制，从而更好地学习和使用该包。
  - error.go
    Summarized code for error.go

这个代码文件定义了一个用于解析器的错误处理系统。以下是对该文件的详细解释：

### 1. 包声明和导入
```go
package jparse

import (
	"fmt"
	"regexp"
)
```
- `package jparse`：声明了该文件属于 `jparse` 包。
- `import`：导入了 `fmt` 和 `regexp` 包，分别用于格式化字符串和正则表达式处理。

### 2. 错误类型定义
```go
// ErrType describes the type of an error.
type ErrType uint

// Error types returned by the parser.
const (
	_ ErrType = iota
	ErrSyntaxError
	ErrUnexpectedEOF
	ErrUnexpectedToken
	ErrMissingToken
	ErrPrefix
	ErrInfix
	ErrUnterminatedString
	ErrUnterminatedRegex
	ErrUnterminatedName
	ErrIllegalEscape
	ErrIllegalEscapeHex
	ErrInvalidNumber
	ErrNumberRange
	ErrEmptyRegex
	ErrInvalidRegex
	ErrGroupPredicate
	ErrGroupGroup
	ErrPathLiteral
	ErrIllegalAssignment
	ErrIllegalParam
	ErrDuplicateParam
	ErrParamCount
	ErrInvalidUnionType
	ErrUnmatchedOption
	ErrUnmatchedSubtype
	ErrInvalidSubtype
	ErrInvalidParamType
)
```
- `ErrType`：定义了一个无符号整数类型 `ErrType`，用于表示不同的错误类型。
- `const`：使用 `iota` 枚举了一系列的错误类型常量，每个常量代表一种特定的错误。

### 3. 错误消息映射
```go
var errmsgs = map[ErrType]string{
	ErrSyntaxError:        "syntax error: '{{token}}'",
	ErrUnexpectedEOF:      "unexpected end of expression",
	ErrUnexpectedToken:    "expected token '{{hint}}', got '{{token}}'",
	ErrMissingToken:       "expected token '{{hint}}' before end of expression",
	ErrPrefix:             "the symbol '{{token}}' cannot be used as a prefix operator",
	ErrInfix:              "the symbol '{{token}}' cannot be used as an infix operator",
	ErrUnterminatedString: "unterminated string literal (no closing '{{hint}}')",
	ErrUnterminatedRegex:  "unterminated regular expression (no closing '{{hint}}')",
	ErrUnterminatedName:   "unterminated name (no closing '{{hint}}')",
	ErrIllegalEscape:      "illegal escape sequence \\{{hint}}",
	ErrIllegalEscapeHex:   "illegal escape sequence \\{{hint}}: \\u must be followed by a 4-digit hexadecimal code point",
	ErrInvalidNumber:      "invalid number literal {{token}}",
	ErrNumberRange:        "invalid number literal {{token}}: value out of range",
	ErrEmptyRegex:         "invalid regular expression: expression cannot be empty",
	ErrInvalidRegex:       "invalid regular expression {{token}}: {{hint}}",
	ErrGroupPredicate:     "a predicate cannot follow a grouping expression in a path step",
	ErrGroupGroup:         "a path step can only have one grouping expression",
	ErrPathLiteral:        "invalid path step {{hint}}: paths cannot contain nulls, strings, numbers or booleans",
	ErrIllegalAssignment:  "illegal assignment: {{hint}} is not a variable",
	ErrIllegalParam:       "illegal function parameter: {{token}} is not a variable",
	ErrDuplicateParam:     "duplicate function parameter: {{token}}",
	ErrParamCount:         "invalid type signature: number of types must match number of function parameters",
	ErrInvalidUnionType:   "invalid type signature: unsupported union type '{{hint}}'",
	ErrUnmatchedOption:    "invalid type signature: option '{{hint}}' must follow a parameter",
	ErrUnmatchedSubtype:   "invalid type signature: subtypes must follow a parameter",
	ErrInvalidSubtype:     "invalid type signature: parameter type {{hint}} does not support subtypes",
	ErrInvalidParamType:   "invalid type signature: unknown parameter type '{{hint}}'",
}
```
- `errmsgs`：定义了一个映射，将错误类型 `ErrType` 映射到相应的错误消息字符串。这些消息字符串中包含占位符 `{{token}}` 和 `{{hint}}`，用于在运行时替换为具体的错误信息。

### 4. 正则表达式用于替换错误消息中的占位符
```go
var reErrMsg = regexp.MustCompile("{{(token|hint)}}")
```
- `reErrMsg`：定义了一个正则表达式，用于匹配错误消息中的占位符 `{{token}}` 和 `{{hint}}`。

### 5. 错误结构体定义
```go
// Error describes an error during parsing.
type Error struct {
	Type     ErrType
	Token    string
	Hint     string
	Position int
}
```
- `Error`：定义了一个结构体 `Error`，用于表示解析过程中发生的错误。包含错误类型、错误相关的 token、提示信息和错误发生的位置。

### 6. 创建错误的方法
```go
func newError(typ ErrType, tok token) error {
	return newErrorHint(typ, tok, "")
}

func newErrorHint(typ ErrType, tok token, hint string) error {
	return &Error{
		Type:     typ,
		Token:    tok.Value,
		Position: tok.Position,
		Hint:     hint,
	}
}
```
- `newError`：创建一个错误对象，不包含提示信息。
- `newErrorHint`：创建一个错误对象，包含提示信息。

### 7. 错误字符串方法
```go
func (e Error) Error() string {
	s := errmsgs[e.Type]
	if s == "" {
		return fmt.Sprintf("parser.Error: unknown error type %d", e.Type)
	}

	return reErrMsg.ReplaceAllStringFunc(s, func(match string) string {
		switch match {
		case "{{token}}":
			return e.Token
		case "{{hint}}":
			return e.Hint
		default:
			return match
		}
	})
}
```
- `Error` 方法：实现了 `error` 接口，返回错误消息字符串。该方法首先从 `errmsgs` 映射中获取错误消息模板，然后使用正则表达式替换模板中的占位符 `{{token}}` 和 `{{hint}}`。

### 8. 辅助函数
```go
func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
```
- `panicf`：一个辅助函数，用于格式化字符串并抛出 panic。

### 总结
该代码文件定义了一个用于解析器的错误处理系统，包括错误类型、错误消息映射、错误结构体以及创建和格式化错误的方法。通过这些定义和方法，解析器能够在遇到不同类型的错误时，生成并返回相应的错误消息，便于调试和问题定位。
  - jparse.go
    Summarized code for jparse.go

该代码文件实现了一个基于Pratt算法的JSONata表达式解析器，将JSONata表达式转换为抽象语法树（AST）。以下是对代码的详细解释：

### 功能概述
该解析器的主要功能是将JSONata表达式解析为抽象语法树。它基于Pratt的Top Down Operator Precedence算法，通过定义一系列的nud（null denotation）和led（left denotation）函数来处理不同类型的token，并使用绑定权力（binding powers）来确定操作符的优先级。

### 实现细节

#### 类型定义
1. **nud**: 一个函数类型，用于处理前缀位置的token，返回一个表示该token值的节点。
2. **led**: 一个函数类型，用于处理中缀位置的token，返回一个表示中缀操作的节点。

#### 全局变量
1. **nuds**: 定义了不同token类型的nud函数。
2. **leds**: 定义了不同token类型的led函数。
3. **bps**: 定义了不同token类型的绑定权力（操作符优先级）。

#### 函数和方法
1. **lookupNud**: 根据token类型查找对应的nud函数。
2. **lookupLed**: 根据token类型查找对应的led函数。
3. **lookupBp**: 根据token类型查找对应的绑定权力。
4. **Parse**: 主函数，用于解析JSONata表达式并返回抽象语法树的根节点。
5. **parser**: 解析器结构体，包含lexer、当前token以及查找nud、led和绑定权力的函数指针。
6. **newParser**: 初始化解析器实例，设置初始token并避免初始化循环。
7. **parseExpression**: 核心函数，实现Pratt算法，处理token的前缀和中缀操作。
8. **advance**: 获取下一个token并更新当前token。
9. **consume**: 检查当前token是否为预期类型，并获取下一个token。
10. **bp**: 返回给定token类型的绑定权力。
11. **initBindingPowers**: 计算并初始化绑定权力值。
12. **validateBindingPowers**: 验证绑定权力值的正确性。

### 代码逻辑
1. **初始化**: 通过`newParser`函数初始化解析器实例，设置lexer和初始token。
2. **解析表达式**: 通过`parseExpression`函数递归解析表达式，根据当前token的类型调用相应的nud或led函数，并根据绑定权力处理操作符优先级。
3. **错误处理**: 使用panic机制处理错误，并在`Parse`函数中捕获并返回错误。

### 关键点
- **nud和led函数**: 根据token的类型和位置，将token转换为AST节点。
- **绑定权力**: 用于确定操作符的优先级，确保表达式解析的正确性。
- **递归解析**: 通过递归调用`parseExpression`函数，处理嵌套和复杂的表达式结构。

### 总结
该代码文件实现了一个功能强大的JSONata表达式解析器，通过定义nud和led函数以及绑定权力，能够准确地将JSONata表达式解析为抽象语法树。代码结构清晰，逻辑严谨，适合学习和阅读。
  - jparse_test.go
    [File too long to summarize]
  - lexer.go
    Summarized code for lexer.go

该代码文件定义了一个用于解析JSONata表达式的词法分析器（lexer）。JSONata是一种用于查询和转换JSON数据的语言。以下是代码的详细功能和实现细节：

### 1. 包声明和导入
```go
package jparse

import (
	"fmt"
	"unicode/utf8"
)
```
- 包名为`jparse`，表示这是一个用于解析JSONata表达式的包。
- 导入了`fmt`和`unicode/utf8`包，用于格式化输出和处理UTF-8编码的字符。

### 2. 常量和类型定义
```go
const eof = -1

type tokenType uint8

const (
	typeEOF tokenType = iota
	typeError

	typeString   // string literal, e.g. "hello"
	typeNumber   // number literal, e.g. 3.14159
	typeBoolean  // true or false
	typeNull     // null
	typeName     // field name, e.g. Price
	typeNameEsc  // escaped field name, e.g. `Product Name`
	typeVariable // variable, e.g. $x
	typeRegex    // regular expression, e.g. /ab+/

	// Symbol operators
	typeBracketOpen
	typeBracketClose
	typeBraceOpen
	typeBraceClose
	typeParenOpen
	typeParenClose
	typeDot
	typeComma
	typeColon
	typeSemicolon
	typeCondition
	typePlus
	typeMinus
	typeMult
	typeDiv
	typeMod
	typePipe
	typeEqual
	typeNotEqual
	typeLess
	typeLessEqual
	typeGreater
	typeGreaterEqual
	typeApply
	typeSort
	typeConcat
	typeRange
	typeAssign
	typeDescendent

	// Keyword operators
	typeAnd
	typeOr
	typeIn
)
```
- 定义了一个常量`eof`表示文件结束。
- 定义了一个枚举类型`tokenType`，用于表示不同类型的token。
- 使用`iota`定义了一系列的`tokenType`常量，包括各种类型的token，如字符串、数字、布尔值、null、字段名、变量、正则表达式以及各种符号和关键字操作符。

### 3. tokenType的字符串表示
```go
func (tt tokenType) String() string {
	switch tt {
	case typeEOF:
		return "(eof)"
	case typeError:
		return "(error)"
	case typeString:
		return "(string)"
	case typeNumber:
		return "(number)"
	case typeBoolean:
		return "(boolean)"
	case typeName, typeNameEsc:
		return "(name)"
	case typeVariable:
		return "(variable)"
	case typeRegex:
		return "(regex)"
	default:
		if s := symbolsAndKeywords[tt]; s != "" {
			return s
		}
		return "(unknown)"
	}
}
```
- 为`tokenType`类型定义了一个`String`方法，用于返回该类型的字符串表示。

### 4. 符号和关键字的映射
```go
var symbols1 = [...]tokenType{
	'[': typeBracketOpen,
	']': typeBracketClose,
	'{': typeBraceOpen,
	'}': typeBraceClose,
	'(': typeParenOpen,
	')': typeParenClose,
	'.': typeDot,
	',': typeComma,
	';': typeSemicolon,
	':': typeColon,
	'?': typeCondition,
	'+': typePlus,
	'-': typeMinus,
	'*': typeMult,
	'/': typeDiv,
	'%': typeMod,
	'|': typePipe,
	'=': typeEqual,
	'<': typeLess,
	'>': typeGreater,
	'^': typeSort,
	'&': typeConcat,
}

type runeTokenType struct {
	r  rune
	tt tokenType
}

var symbols2 = [...][]runeTokenType{
	'!': {{'=', typeNotEqual}},
	'<': {{'=', typeLessEqual}},
	'>': {{'=', typeGreaterEqual}},
	'.': {{'.', typeRange}},
	'~': {{'>', typeApply}},
	':': {{'=', typeAssign}},
	'*': {{'*', typeDescendent}},
}

const (
	symbol1Count = rune(len(symbols1))
	symbol2Count = rune(len(symbols2))
)

func lookupSymbol1(r rune) tokenType {
	if r < 0 || r >= symbol1Count {
		return 0
	}
	return symbols1[r]
}

func lookupSymbol2(r rune) []runeTokenType {
	if r < 0 || r >= symbol2Count {
		return nil
	}
	return symbols2[r]
}

func lookupKeyword(s string) tokenType {
	switch s {
	case "and":
		return typeAnd
	case "or":
		return typeOr
	case "in":
		return typeIn
	case "true", "false":
		return typeBoolean
	case "null":
		return typeNull
	default:
		return 0
	}
}
```
- `symbols1`数组将单字符符号映射到相应的`tokenType`。
- `symbols2`数组将双字符符号映射到相应的`tokenType`。
- `lookupSymbol1`和`lookupSymbol2`函数用于查找单字符和双字符符号对应的`tokenType`。
- `lookupKeyword`函数用于查找关键字对应的`tokenType`。

### 5. token结构体
```go
type token struct {
	Type     tokenType
	Value    string
	Position int
}
```
- `token`结构体表示一个token，包含类型、值和位置信息。

### 6. lexer结构体
```go
type lexer struct {
	input   string
	length  int
	start   int
	current int
	width   int
	err     error
}

func newLexer(input string) lexer {
	return lexer{
		input:  input,
		length: len(input),
	}
}
```
- `lexer`结构体表示一个词法分析器，包含输入字符串、长度、当前位置、宽度等信息。
- `newLexer`函数用于创建一个新的词法分析器。

### 7. 词法分析器的核心方法
```go
func (l *lexer) next(allowRegex bool) token {
	l.skipWhitespace()
	ch := l.nextRune()
	if ch == eof {
		return l.eof()
	}
	if allowRegex && ch == '/' {
		l.ignore()
		return l.scanRegex(ch)
	}
	if rts := lookupSymbol2(ch); rts != nil {
		for _, rt := range rts {
			if l.acceptRune(rt.r) {
				return l.newToken(rt.tt)
			}
		}
	}
	if tt := lookupSymbol1(ch); tt > 0 {
		return l.newToken(tt)
	}
	if ch == '"' || ch == '\'' {
		l.ignore()
		return l.scanString(ch)
	}
	if ch >= '0' && ch <= '9' {
		l.backup()
		return l.scanNumber()
	}
	if ch == '`' {
		l.ignore()
		return l.scanEscapedName(ch)
	}
	l.backup()
	return l.scanName()
}
```
- `next`方法是词法分析器的核心方法，用于获取下一个token。
- 根据当前字符的不同，调用不同的扫描方法（如`scanRegex`、`scanString`、`scanNumber`、`scanEscapedName`、`scanName`）来生成相应的token。

### 8. 扫描方法
```go
func (l *lexer) scanRegex(delim rune) token {
	var depth int
Loop:
	for {
		switch l.nextRune() {
		case delim:
			if depth == 0 {
				break Loop
			}
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		case '\\':
			if r := l.nextRune(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.error(ErrUnterminatedRegex, string(delim))
		}
	}
	l.backup()
	t := l.newToken(typeRegex)
	l.acceptRune(delim)
	l.ignore()
	if l.acceptAll(isRegexFlag) {
		flags := l.newToken(0)
		t.Value = fmt.Sprintf("(?%s)%s", flags.Value, t.Value)
	}
	return t
}

func (l *lexer) scanString(quote rune) token {
Loop:
	for {
		switch l.nextRune() {
		case quote:
			break Loop
		case '\\':
			if r := l.nextRune(); r != eof {
				break
			}
			fallthrough
		case eof:
			return l.error(ErrUnterminatedString, string(quote))
		}
	}
	l.backup()
	t := l.newToken(typeString)
	l.acceptRune(quote)
	l.ignore()
	return t
}

func (l *lexer) scanNumber() token {
	if !l.acceptRune('0') {
		l.accept(isNonZeroDigit)
		l.acceptAll(isDigit)
	}
	if l.acceptRune('.') {
		if !l.acceptAll(isDigit) {
			l.backup()
			return l.newToken(typeNumber)
		}
	}
	if l.acceptRunes2('e', 'E') {
		l.acceptRunes2('+', '-')
		l.acceptAll(isDigit)
	}
	return l.newToken(typeNumber)
}

func (l *lexer) scanEscapedName(quote rune) token {
Loop:
	for {
		switch l.nextRune() {
		case quote:
			break Loop
		case eof, '\n':
			return l.error(ErrUnterminatedName, string(quote))
		}
	}
	l.backup()
	t := l.newToken(typeNameEsc)
	l.acceptRune(quote)
	l.ignore()
	return t
}

func (l *lexer) scanName() token {
	isVar := l.acceptRune('$')
	if isVar {
		l.ignore()
	}
	for {
		ch := l.nextRune()
		if ch == eof {
			break
		}
		if isWhitespace(ch) {
			l.backup()
			break
		}
		if lookupSymbol1(ch) > 0 || lookupSymbol2(ch) != nil {
			l.backup()
			break
		}
	}
	t := l.newToken(typeName)
	if isVar {
		t.Type = typeVariable
	} else if tt := lookupKeyword(t.Value); tt > 0 {
		t.Type = tt
	}
	return t
}
```
- `scanRegex`方法用于扫描正则表达式。
- `scanString`方法用于扫描字符串。
- `scanNumber`方法用于扫描数字。
- `scanEscapedName`方法用于扫描转义字段名。
- `scanName`方法用于扫描字段名、变量或关键字。

### 9. 辅助方法
```go
func (l *lexer) eof() token {
	return token{
		Type:     typeEOF,
		Position: l.current,
	}
}

func (l *lexer) error(typ ErrType, hint string) token {
	t := l.newToken(typeError)
	l.err = newErrorHint(typ, t, hint)
	return t
}

func (l *lexer) newToken(tt tokenType) token {
	t := token{
		Type:     tt,
		Value:    l.input[l.start:l.current],
		Position: l.start,
	}
	l.width = 0
	l.start = l.current
	return t
}

func (l *lexer) nextRune() rune {
	if l.err != nil || l.current >= l.length {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.current:])
	l.width = w
	l.current += w
	return r
}

func (l *lexer) backup() {
	l.current -= l.width
}

func (l *lexer) ignore() {
	l.start = l.current
}

func (l *lexer) acceptRune(r rune) bool {
	return l.accept(func(c rune) bool {
		return c == r
	})
}

func (l *lexer) acceptRunes2(r1, r2 rune) bool {
	return l.accept(func(c rune) bool {
		return c == r1 || c == r2
	})
}

func (l *lexer) accept(isValid func(rune) bool) bool {
	if isValid(l.nextRune()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptAll(isValid func(rune) bool) bool {
	var b bool
	for l.accept(isValid) {
		b = true
	}
	return b
}

func (l *lexer) skipWhitespace() {
	l.acceptAll(isWhitespace)
	l.ignore()
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', '\v':
		return true
	default:
		return false
	}
}

func isRegexFlag(r rune) bool {
	switch r {
	case 'i', 'm', 's':
		return true
	default:
		return false
	}
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isNonZeroDigit(r rune) bool {
	return r >= '1' && r <= '9'
}
```
- `eof`方法用于返回文件结束的token。
- `error`方法用于返回错误token。
- `newToken`方法用于创建一个新的token。
- `nextRune`方法用于获取下一个字符。
- `backup`方法用于回退一个字符。
- `ignore`方法用于忽略当前字符。
- `acceptRune`、`acceptRunes2`、`accept`和`acceptAll`方法用于接受特定字符。
- `skipWhitespace`方法用于跳过空白字符。
- `isWhitespace`、`isRegexFlag`、`isDigit`和`isNonZeroDigit`函数用于判断字符类型。

### 10. 符号和关键字的字符串表示
```go
var symbolsAndKeywords = func() map[tokenType]string {
	m := map[tokenType]string{
		typeAnd:  "and",
		typeOr:   "or",
		typeIn:   "in",
		typeNull: "null",
	}
	for r, tt := range symbols1 {
		if tt > 0 {
			m[tt] = fmt.Sprintf("%c", r)
		}
	}
	for r, rts := range symbols2 {
		for _, rt := range rts {
			m[rt.tt] = fmt.Sprintf("%c", r) + fmt.Sprintf("%c", rt.r)
		}
	}
	return m
}()
```
- `symbolsAndKeywords`映射将操作符token类型映射回它们的字符串表示。

### 总结
该代码文件实现了一个用于解析JSONata表达式的词法分析器。它定义了各种token类型，并通过一系列的扫描方法将输入字符串转换为token序列。词法分析器通过识别不同的字符模式来生成相应的token，并处理各种特殊情况，如正则表达式、字符串、数字等。
  - lexer_test.go
    Summarized code for lexer_test.go

该代码文件是一个用于解析JSON的词法分析器（lexer）的测试文件。词法分析器的主要功能是将输入的字符串分解成一系列的标记（tokens），这些标记代表了输入中的不同元素，如数字、字符串、符号等。以下是对该代码文件的详细解析：

### 1. 包声明和导入
```go
package jparse

import (
	"reflect"
	"testing"
)
```
- `package jparse`：声明了该文件属于`jparse`包。
- `import`：导入了`reflect`和`testing`包，分别用于深度比较和测试。

### 2. 类型定义
```go
type lexerTestCase struct {
	Input      string
	AllowRegex bool
	Tokens     []token
	Error      error
}
```
- `lexerTestCase`：定义了一个结构体，用于表示词法分析器的测试用例。包含输入字符串、是否允许正则表达式、期望的标记列表和期望的错误。

### 3. 测试函数
#### 3.1. 测试空白字符
```go
func TestLexerWhitespace(t *testing.T) {
	testLexer(t, []lexerTestCase{
		{
			Input: "",
		},
		{
			Input: "       ",
		},
		{
			Input: "\v\t\r\n",
		},
		{
			Input: `


			`,
		},
	})
}
```
- `TestLexerWhitespace`：测试空白字符的处理。输入为空字符串或仅包含空白字符的字符串。

#### 3.2. 测试正则表达式
```go
func TestLexerRegex(t *testing.T) {
	testLexer(t, []lexerTestCase{
		{
			Input: `//`,
			Tokens: []token{
				tok(typeDiv, "/", 0),
				tok(typeDiv, "/", 1),
			},
		},
		{
			Input:      `//`,
			AllowRegex: true,
			Tokens: []token{
				tok(typeRegex, "", 1),
			},
		},
		{
			Input:      `/ab+/`,
			AllowRegex: true,
			Tokens: []token{
				tok(typeRegex, "ab+", 1),
			},
		},
		{
			Input:      `/(ab+/)/`,
			AllowRegex: true,
			Tokens: []token{
				tok(typeRegex, "(ab+/)", 1),
			},
		},
		{
			Input:      `/ab+/i`,
			AllowRegex: true,
			Tokens: []token{
				tok(typeRegex, "(?i)ab+", 1),
			},
		},
		{
			Input:      `/ab+/ i`,
			AllowRegex: true,
			Tokens: []token{
				tok(typeRegex, "ab+", 1),
				tok(typeName, "i", 6),
			},
		},
		{
			Input:      `/ab+/I`,
			AllowRegex: true,
			Tokens: []token{
				tok(typeRegex, "ab+", 1),
				tok(typeName, "I", 5),
			},
		},
		{
			Input:      `/ab+`,
			AllowRegex: true,
			Tokens: []token{
				tok(typeError, "ab+", 1),
			},
			Error: &Error{
				Type:     ErrUnterminatedRegex,
				Token:    "ab+",
				Hint:     "/",
				Position: 1,
			},
		},
	})
}
```
- `TestLexerRegex`：测试正则表达式的处理。包括不同形式的正则表达式和错误处理。

#### 3.3. 测试字符串
```go
func TestLexerStrings(t *testing.T) {
	testLexer(t, []lexerTestCase{
		{
			Input: `""`,
			Tokens: []token{
				tok(typeString, "", 1),
			},
		},
		{
			Input: `''`,
			Tokens: []token{
				tok(typeString, "", 1),
			},
		},
		{
			Input: `"double quotes"`,
			Tokens: []token{
				tok(typeString, "double quotes", 1),
			},
		},
		{
			Input: "'single quotes'",
			Tokens: []token{
				tok(typeString, "single quotes", 1),
			},
		},
		{
			Input: `"escape\t"`,
			Tokens: []token{
				tok(typeString, "escape\\t", 1),
			},
		},
		{
			Input: `'escape\u0036'`,
			Tokens: []token{
				tok(typeString, "escape\\u0036", 1),
			},
		},
		{
			Input: `"超明體繁"`,
			Tokens: []token{
				tok(typeString, "超明體繁", 1),
			},
		},
		{
			Input: `'日本語'`,
			Tokens: []token{
				tok(typeString, "日本語", 1),
			},
		},
		{
			Input: `"No closing quote...`,
			Tokens: []token{
				tok(typeError, "No closing quote...", 1),
			},
			Error: &Error{
				Type:     ErrUnterminatedString,
				Token:    "No closing quote...",
				Hint:     "\"",
				Position: 1,
			},
		},
		{
			Input: `'No closing quote...`,
			Tokens: []token{
				tok(typeError, "No closing quote...", 1),
			},
			Error: &Error{
				Type:     ErrUnterminatedString,
				Token:    "No closing quote...",
				Hint:     "'",
				Position: 1,
			},
		},
	})
}
```
- `TestLexerStrings`：测试字符串的处理。包括不同形式的字符串和错误处理。

#### 3.4. 测试数字
```go
func TestLexerNumbers(t *testing.T) {
	testLexer(t, []lexerTestCase{
		{
			Input: "1",
			Tokens: []token{
				tok(typeNumber, "1", 0),
			},
		},
		{
			Input: "3.14159",
			Tokens: []token{
				tok(typeNumber, "3.14159", 0),
			},
		},
		{
			Input: "1e10",
			Tokens: []token{
				tok(typeNumber, "1e10", 0),
			},
		},
		{
			Input: "1E-10",
			Tokens: []token{
				tok(typeNumber, "1E-10", 0),
			},
		},
		{
			Input: "-100",
			Tokens: []token{
				tok(typeMinus, "-", 0),
				tok(typeNumber, "100", 1),
			},
		},
		{
			Input: "007",
			Tokens: []token{
				tok(typeNumber, "0", 0),
				tok(typeNumber, "0", 1),
				tok(typeNumber, "7", 2),
			},
		},
		{
			Input: ".5",
			Tokens: []token{
				tok(typeDot, ".", 0),
				tok(typeNumber, "5", 1),
			},
		},
		{
			Input: "5. ",
			Tokens: []token{
				tok(typeNumber, "5", 0),
				tok(typeDot, ".", 1),
			},
		},
	})
}
```
- `TestLexerNumbers`：测试数字的处理。包括不同形式的数字和错误处理。

#### 3.5. 测试名称
```go
func TestLexerNames(t *testing.T) {
	testLexer(t, []lexerTestCase{
		{
			Input: "hello",
			Tokens: []token{
				tok(typeName, "hello", 0),
			},
		},
		{
			Input: "hello world",
			Tokens: []token{
				tok(typeName, "hello", 0),
				tok(typeName, "world", 6),
			},
		},
		{
			Input: "hello, world.",
			Tokens: []token{
				tok(typeName, "hello", 0),
				tok(typeComma, ",", 5),
				tok(typeName, "world", 7),
				tok(typeDot, ".", 12),
			},
		},
		{
			Input: "HELLO!",
			Tokens: []token{
				tok(typeName, "HELLO", 0),
				tok(typeName, "!", 5),
			},
		},
		{
			Input: "`hello, world.`",
			Tokens: []token{
				tok(typeNameEsc, "hello, world.", 1),
			},
		},
		{
			Input: "`true or false`",
			Tokens: []token{
				tok(typeNameEsc, "true or false", 1),
			},
		},
		{
			Input: "`no closing quote...",
			Tokens: []token{
				tok(typeError, "no closing quote...", 1),
			},
			Error: &Error{
				Type:     ErrUnterminatedName,
				Token:    "no closing quote...",
				Hint:     "`",
				Position: 1,
			},
		},
	})
}
```
- `TestLexerNames`：测试名称的处理。包括不同形式的名称和错误处理。

#### 3.6. 测试变量
```go
func TestLexerVariables(t *testing.T) {
	testLexer(t, []lexerTestCase{
		{
			Input: "$",
			Tokens: []token{
				tok(typeVariable, "", 1),
			},
		},
		{
			Input: "$$",
			Tokens: []token{
				tok(typeVariable, "$", 1),
			},
		},
		{
			Input: "$var",
			Tokens: []token{
				tok(typeVariable, "var", 1),
			},
		},
		{
			Input: "$uppercase",
			Tokens: []token{
				tok(typeVariable, "uppercase", 1),
			},
		},
	})
}
```
- `TestLexerVariables`：测试变量的处理。包括不同形式的变量。

#### 3.7. 测试符号和关键字
```go
func TestLexerSymbolsAndKeywords(t *testing.T) {
	var tests []lexerTestCase

	for tt, s := range symbolsAndKeywords {
		tests = append(tests, lexerTestCase{
			Input: s,
			Tokens: []token{
				tok(tt, s, 0),
			},
		})
	}

	testLexer(t, tests)
}
```
- `TestLexerSymbolsAndKeywords`：测试符号和关键字的处理。通过遍历`symbolsAndKeywords`生成测试用例。

### 4. 辅助函数
#### 4.1. `testLexer`
```go
func testLexer(t *testing.T, data []lexerTestCase) {
	for _, test := range data {
		l := newLexer(test.Input)
		eof := tok(typeEOF, "", len(test.Input))

		for _, exp := range test.Tokens {
			compareTokens(t, test.Input, exp, l.next(test.AllowRegex))
		}

		compareErrors(t, test.Input, test.Error, l.err)

		for i := 0; i < 3; i++ {
			compareTokens(t, test.Input, eof, l.next(test.AllowRegex))
		}
	}
}
```
- `testLexer`：执行词法分析器的测试。遍历测试用例，创建词法分析器实例，并比较生成的标记和错误。

#### 4.2. `compareTokens`
```go
func compareTokens(t *testing.T, prefix string, exp, got token) {
	if got.Type != exp.Type {
		t.Errorf("%s: expected token with Type '%s', got '%s'", prefix, exp.Type, got.Type)
	}

	if got.Value != exp.Value {
		t.Errorf("%s: expected token with Value %q, got %q", prefix, exp.Value, got.Value)
	}

	if got.Position != exp.Position {
		t.Errorf("%s: expected token with Position %d, got %d", prefix, exp.Position, got.Position)
	}
}
```
- `compareTokens`：比较两个标记是否相同。

#### 4.3. `compareErrors`
```go
func compareErrors(t *testing.T, prefix string, exp, got error) {
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("%s: expected error %v, got %v", prefix, exp, got)
	}
}
```
- `compareErrors`：比较两个错误是否相同。

#### 4.4. `tok`
```go
func tok(typ tokenType, value string, position int) token {
	return token{
		Type:     typ,
		Value:    value,
		Position: position,
	}
}
```
- `tok`：创建一个新的标记实例。

### 总结
该代码文件通过一系列的测试用例，全面测试了词法分析器的功能，包括空白字符、正则表达式、字符串、数字、名称、变量、符号和关键字的处理。通过这些测试，可以确保词法分析器在各种输入情况下都能正确地生成标记，并处理错误情况。
  - node.go
    [File too long to summarize]
- jsonata-server
  - .gitignore
  - README.md
  - bench.go
    Summarized code for bench.go

### 代码文件内容总结

#### 概述
该代码文件实现了一个简单的HTTP服务器，该服务器使用JSONata表达式对预定义的JSON数据进行处理和验证。主要功能包括：
1. 解析JSON数据。
2. 编译和执行JSONata表达式。
3. 处理HTTP请求，返回处理结果或错误信息。

#### 详细功能和实现细节

1. **包声明和导入**
   ```go
   package main

   import (
   	"log"
   	"net/http"
   	"encoding/json"
   	jsonata "github.com/blues/jsonata-go"
   )
   ```
   - 该文件属于`main`包。
   - 导入了多个标准库和第三方库：
     - `log`：用于日志记录。
     - `net/http`：用于HTTP服务器功能。
     - `encoding/json`：用于JSON解析。
     - `jsonata "github.com/blues/jsonata-go"`：用于JSONata表达式的编译和执行。

2. **全局变量**
   ```go
   var (
   	benchData = []byte(`...`) // 预定义的JSON数据
   	benchExpression = `...`   // JSONata表达式
   )
   ```
   - `benchData`：包含一个预定义的JSON数据，用于后续处理。
   - `benchExpression`：包含一个JSONata表达式，用于对`benchData`进行处理。

3. **数据解析**
   ```go
   var data interface{}

   func init() {
   	if err := json.Unmarshal(benchData, &data); err != nil {
   		panic(err)
   	}
   }
   ```
   - 在初始化阶段，将`benchData`解析为Go语言的`interface{}`类型，存储在`data`变量中。

4. **HTTP处理函数**
   ```go
   func benchmark(w http.ResponseWriter, r *http.Request) {
   	expr, err := jsonata.Compile(benchExpression)
   	if err != nil {
   		bencherr(w, err)
   	}

   	_, err = expr.Eval(data)
   	if err != nil {
   		bencherr(w, err)
   	}

   	if _, err := w.Write([]byte("success")); err != nil {
   		log.Fatal(err)
   	}
   }
   ```
   - `benchmark`函数是HTTP请求的处理函数。
   - 首先，编译`benchExpression`表达式。
   - 然后，使用编译后的表达式对`data`进行求值。
   - 如果处理成功，向客户端返回字符串"success"。

5. **错误处理函数**
   ```go
   func bencherr(w http.ResponseWriter, err error) {
   	log.Println(err)
   	http.Error(w, err.Error(), http.StatusInternalServerError)
   }
   ```
   - `bencherr`函数用于处理错误情况。
   - 记录错误日志，并向客户端返回HTTP 500错误。

### 总结
该代码文件实现了一个基于HTTP的JSON数据处理服务，使用JSONata表达式对预定义的JSON数据进行处理。主要步骤包括：
1. 解析JSON数据。
2. 编译和执行JSONata表达式。
3. 处理HTTP请求，返回处理结果或错误信息。

通过阅读该代码，读者可以学习到如何使用Go语言实现一个简单的HTTP服务器，以及如何使用JSONata表达式处理JSON数据。
  - exts.go
    Summarized code for exts.go

这段代码文件主要实现了两个功能：将Unix时间戳转换为格式化的时间字符串，以及将格式化的时间字符串转换为Unix时间戳。以下是对代码的详细解释：

### 1. 包声明和导入
```go
package main

import (
	"github.com/blues/jsonata-go/jlib"
	"github.com/blues/jsonata-go/jtypes"
)
```
- `package main`：声明这是一个可执行程序的入口包。
- `import`：导入了两个外部包，`jlib` 和 `jtypes`，这两个包提供了与JSONata相关的功能和类型。

### 2. 常量声明
```go
const defaultDateFormat = "[Y]-[M01]-[D01] [H01]:[m] [ZN]"
```
- `defaultDateFormat`：定义了一个默认的日期格式字符串，用于时间格式化和解析。

### 3. `formatTime` 函数
```go
func formatTime(secs int64, picture jtypes.OptionalString, tz jtypes.OptionalString) (string, error) {
	if picture.String == "" {
		picture = jtypes.NewOptionalString(defaultDateFormat)
	}
	return jlib.FromMillis(secs*1000, picture, tz)
}
```
- `formatTime`：将Unix时间戳（以秒为单位）转换为格式化的时间字符串。
  - `secs`：输入的Unix时间戳（秒）。
  - `picture`：可选的日期格式字符串，如果为空，则使用默认格式。
  - `tz`：可选的时区字符串。
  - 函数内部首先检查`picture`是否为空，如果为空则使用默认的日期格式。
  - 调用`jlib.FromMillis`函数，将秒转换为毫秒，并根据给定的格式和时区生成时间字符串。

### 4. `parseTime` 函数
```go
func parseTime(value string, picture jtypes.OptionalString, tz jtypes.OptionalString) (int64, error) {
	if picture.String == "" {
		picture = jtypes.NewOptionalString(defaultDateFormat)
	}
	ms, err := jlib.ToMillis(value, picture, tz)
	if err != nil {
		return 0, err
	}
	return ms / 1000, nil
}
```
- `parseTime`：将格式化的时间字符串转换为Unix时间戳（以秒为单位）。
  - `value`：输入的格式化时间字符串。
  - `picture`：可选的日期格式字符串，如果为空，则使用默认格式。
  - `tz`：可选的时区字符串。
  - 函数内部首先检查`picture`是否为空，如果为空则使用默认的日期格式。
  - 调用`jlib.ToMillis`函数，将格式化的时间字符串转换为毫秒。
  - 将毫秒转换为秒并返回。

### 总结
这段代码主要提供了两个核心功能：
1. **时间格式化**：将Unix时间戳转换为指定格式的字符串。
2. **时间解析**：将指定格式的字符串转换为Unix时间戳。

通过使用外部包`jlib`和`jtypes`，代码实现了对时间格式和解析的灵活处理，包括可选的日期格式和时区支持。
  - main.go
    Summarized code for main.go

### 代码文件内容总结

#### 1. 版权声明和包声明
```go
// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main
```
- 代码文件的版权声明和包声明部分，表明代码属于 Blues Inc. 并受特定许可证约束。
- 包声明为 `main`，表示这是一个可执行程序。

#### 2. 导入依赖包
```go
import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"

	jsonata "github.com/blues/jsonata-go"
	"github.com/blues/jsonata-go/jtypes"
)
```
- 导入了多个标准库包和第三方包，用于处理 HTTP 请求、JSON 编解码、命令行参数解析等。
- `jsonata` 和 `jtypes` 是第三方包，用于支持 JSONata 表达式的编译和求值。

#### 3. 初始化函数
```go
func init() {
	argUndefined0 := jtypes.ArgUndefined(0)

	exts := map[string]jsonata.Extension{
		"formatTime": {
			Func:             formatTime,
			UndefinedHandler: argUndefined0,
		},
		"parseTime": {
			Func:             parseTime,
			UndefinedHandler: argUndefined0,
		},
	}

	if err := jsonata.RegisterExts(exts); err != nil {
		panic(err)
	}
}
```
- `init` 函数在程序启动时自动执行，用于注册 JSONata 扩展函数。
- 定义了两个扩展函数 `formatTime` 和 `parseTime`，并使用 `jtypes.ArgUndefined(0)` 作为未定义参数的处理函数。
- 通过 `jsonata.RegisterExts` 注册这些扩展函数，如果注册失败则抛出异常。

#### 4. 主函数
```go
func main() {
	port := flag.Uint("port", 8080, "The port `number` to serve on")
	flag.Parse()

	http.HandleFunc("/eval", evaluate)
	http.HandleFunc("/bench", benchmark)
	http.Handle("/", http.FileServer(http.Dir("site")))

	log.Printf("Starting JSONata Server on port %d:\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
```
- 主函数定义了一个命令行参数 `port`，默认值为 8080，用于指定服务器监听的端口。
- 注册了三个 HTTP 处理函数：
  - `/eval` 处理 JSONata 表达式的求值请求。
  - `/bench` 处理性能基准测试请求（代码中未实现）。
  - `/` 处理静态文件请求，提供 `site` 目录下的文件。
- 启动 HTTP 服务器并监听指定端口，如果启动失败则记录错误并退出。

#### 5. evaluate 函数
```go
func evaluate(w http.ResponseWriter, r *http.Request) {
	input := strings.TrimSpace(r.FormValue("json"))
	if input == "" {
		http.Error(w, "Input is empty", http.StatusBadRequest)
		return
	}

	expression := strings.TrimSpace(r.FormValue("expr"))
	if expression == "" {
		http.Error(w, "Expression is empty", http.StatusBadRequest)
		return
	}

	b, status, err := eval(input, expression)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), status)
		return
	}

	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}
```
- `evaluate` 函数处理 `/eval` 请求，从请求中获取 JSON 数据和 JSONata 表达式。
- 检查输入是否为空，如果为空则返回 400 错误。
- 调用 `eval` 函数进行表达式求值，并根据结果返回响应。

#### 6. eval 函数
```go
func eval(input, expression string) (b []byte, status int, err error) {
	defer func() {
		if r := recover(); r != nil {
			b = nil
			status = http.StatusInternalServerError
			err = fmt.Errorf("PANIC: %v", r)
			return
		}
	}()

	var data interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("input error: %s", err)
	}

	expr, err := jsonata.Compile(expression)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("compile error: %s", err)
	}

	result, err := expr.Eval(data)
	if err != nil {
		if err == jsonata.ErrUndefined {
			return []byte("No results found"), http.StatusOK, nil
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("eval error: %s", err)
	}

	b, err = jsonify(result)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("encode error: %s", err)
	}

	return b, http.StatusOK, nil
}
```
- `eval` 函数负责 JSON 数据的解析、JSONata 表达式的编译和求值。
- 使用 `defer` 捕获并处理可能的 panic 异常。
- 解析输入的 JSON 数据，编译 JSONata 表达式，并进行求值。
- 如果求值结果未定义，返回 "No results found"。
- 将求值结果转换为 JSON 格式并返回。

#### 7. jsonify 函数
```go
func jsonify(v interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := json.NewEncoder(&b)
	e.SetIndent("", "    ")
	if err := e.Encode(v); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
```
- `jsonify` 函数将任意类型的数据转换为格式化的 JSON 字节数组。
- 使用 `json.NewEncoder` 和 `SetIndent` 方法设置 JSON 编码器的缩进格式。
- 将数据编码为 JSON 并返回字节数组。

### 功能总结
- 该代码实现了一个基于 JSONata 表达式的 HTTP 服务器。
- 支持通过 `/eval` 接口对输入的 JSON 数据进行 JSONata 表达式求值。
- 提供了扩展函数 `formatTime` 和 `parseTime`，用于处理时间格式。
- 支持静态文件服务，提供 `site` 目录下的文件。

### 实现细节
- 使用 `flag` 包解析命令行参数，支持自定义端口。
- 使用 `http` 包处理 HTTP 请求和响应。
- 使用 `jsonata` 包编译和求值 JSONata 表达式。
- 使用 `json` 包进行 JSON 数据的编解码。
- 使用 `log` 包记录日志信息。
- 使用 `bytes` 包处理字节缓冲区。
- 使用 `strings` 包处理字符串操作。
- 使用 `panic` 和 `recover` 处理异常情况。
  - site
    - assets
      - css
        - codemirror.min.css
        - normalize.min.css
        - styles.css
      - js
        - codemirror.min.js
          [File too long to summarize]
        - javascript.min.js
          Summarized code for javascript.min.js

这个代码文件定义了一个用于语法高亮的JavaScript模式，适用于CodeMirror编辑器。以下是代码的详细功能和实现细节：

### 功能概述

1. **模式定义**：定义了一个名为`javascript`的CodeMirror模式，支持JavaScript、ECMAScript、JSON、JSON-LD和TypeScript的语法高亮。
2. **词法分析**：通过定义多个函数（如`a`、`i`、`o`等）来处理不同类型的token，如字符串、数字、注释、关键字等。
3. **状态管理**：使用状态对象来跟踪当前的词法分析状态，包括当前的token类型、缩进级别、上下文信息等。
4. **缩进逻辑**：根据当前的语法结构和输入内容来计算缩进级别。
5. **辅助函数**：定义了一些辅助函数来处理特定情况，如处理箭头函数、类型声明、类定义等。

### 实现细节

#### 1. 模式定义

```javascript
e.defineMode("javascript", function(t, r) {
    // 模式的具体实现
});
```

#### 2. 词法分析

- **字符串和字符**：
  ```javascript
  function a(e, t) {
      var r = e.next();
      if ('"' == r || "'" == r) {
          t.tokenize = function(e) {
              // 处理字符串
          };
      }
      // 其他token的处理
  }
  ```

- **数字**：
  ```javascript
  if ("." == r && e.match(/^\d+(?:[eE][+\-]?\d+)?/)) {
      return n("number", "number");
  }
  ```

- **注释**：
  ```javascript
  function i(e, t) {
      for (var r, i = !1; r = e.next();) {
          if ("/" == r && i) {
              t.tokenize = a;
              break;
          }
          i = "*" == r;
      }
      return n("comment", "comment");
  }
  ```

#### 3. 状态管理

- **状态对象**：
  ```javascript
  function u(e, t, r, n, a, i) {
      this.indented = e, this.column = t, this.type = r, this.prev = a, this.info = i, null != n && (this.align = n);
  }
  ```

- **初始状态**：
  ```javascript
  startState: function(e) {
      var t = {
          tokenize: a,
          lastType: "sof",
          cc: [],
          lexical: new u((e || 0) - Ee, 0, "block", !1),
          localVars: r.localVars,
          context: r.localVars && { vars: r.localVars },
          indented: e || 0
      };
      // 其他初始化
      return t;
  }
  ```

#### 4. 缩进逻辑

- **缩进计算**：
  ```javascript
  indent: function(t, n) {
      if (t.tokenize == i) return e.Pass;
      if (t.tokenize != a) return 0;
      // 缩进计算逻辑
  }
  ```

#### 5. 辅助函数

- **处理箭头函数**：
  ```javascript
  function c(e, t) {
      var r = e.string.indexOf("=>", e.start);
      if (r >= 0) {
          // 处理箭头函数
      }
  }
  ```

- **处理类型声明**：
  ```javascript
  function G(e, t) {
      if ("variable" == e || "void" == t) {
          // 处理类型声明
      }
      // 其他处理
  }
  ```

### 总结

这个代码文件通过定义词法分析函数、状态管理和缩进逻辑，实现了对JavaScript、ECMAScript、JSON、JSON-LD和TypeScript的语法高亮支持。通过这些功能的组合，CodeMirror编辑器能够准确地显示不同类型的代码元素，并根据语法结构调整缩进级别，从而提高代码的可读性和编辑体验。
        - jsonata-codemirror.js
          Summarized code for jsonata-codemirror.js

这段代码定义了一个用于 JSONata 语言的 CodeMirror 语法高亮模式。JSONata 是一种用于查询和转换 JSON 数据的语言。以下是代码的详细功能和实现细节：

### 功能概述

1. **定义 JSONata 模式**：使用 `CodeMirror.defineMode` 方法定义一个名为 "jsonata" 的语法高亮模式。
2. **配置参数**：接受两个配置参数 `config` 和 `parserConfig`，其中 `parserConfig` 包含 `template` 和 `jsonata` 两个属性。
3. **操作符和转义序列**：定义了 JSONata 语言中的操作符和字符串转义序列。
4. **词法分析器**：实现了一个词法分析器 `tokenizer`，用于将输入的字符串分解成不同类型的词法单元（tokens）。
5. **模板解析器**：实现了一个模板解析器 `templatizer`，用于处理包含模板语法的文本。
6. **词法单元类型映射**：定义了词法单元类型到 CodeMirror 样式名称的映射关系。
7. **语法高亮实现**：通过 `token` 方法实现语法高亮，根据不同的词法单元类型返回相应的样式名称。

### 实现细节

#### 1. 操作符和转义序列

- **操作符**：定义了一个 `operators` 对象，包含 JSONata 语言中的各种操作符及其优先级。
- **转义序列**：定义了一个 `escapes` 对象，包含 JSON 字符串中的转义序列及其对应的字符。

#### 2. 词法分析器 (`tokenizer`)

- **初始化**：词法分析器接受一个字符串 `path`，并初始化位置和长度。
- **创建词法单元**：定义了一个 `create` 函数，用于创建词法单元对象。
- **获取下一个词法单元**：定义了一个 `next` 函数，用于获取下一个词法单元。
  - **跳过空白字符**：跳过输入字符串中的空白字符。
  - **处理双字符操作符**：识别并处理双字符操作符，如 `..`、`||`、`:=` 等。
  - **处理单字符操作符**：识别并处理单字符操作符。
  - **处理字符串字面量**：识别并处理字符串字面量，包括转义序列和 Unicode 转义序列。
  - **处理数字字面量**：识别并处理数字字面量。
  - **处理名称**：识别并处理变量引用和关键字。

#### 3. 模板解析器 (`templatizer`)

- **初始化**：模板解析器接受一个字符串 `text`，并初始化位置和长度。
- **创建词法单元**：定义了一个 `create` 函数，用于创建词法单元对象。
- **获取下一个词法单元**：定义了一个 `next` 函数，用于获取下一个词法单元。
  - **处理模板语法**：识别并处理模板语法中的 `{{` 和 `}}`。

#### 4. 词法单元类型映射 (`TOKEN_NAMES`)

- **映射关系**：定义了一个 `TOKEN_NAMES` 对象，将词法单元类型映射到 CodeMirror 样式名称。

#### 5. 语法高亮实现 (`token` 方法)

- **选择词法分析器**：根据 `templateMode` 选择使用 `templatizer` 或 `tokenizer`。
- **获取词法单元**：调用选择的词法分析器获取下一个词法单元。
- **跳过字符**：根据词法单元的位置跳过相应的字符。
- **返回样式名称**：根据词法单元的类型返回相应的样式名称。

### 总结

这段代码通过定义操作符、转义序列、词法分析器和模板解析器，实现了 JSONata 语言的语法高亮功能。通过 `token` 方法，根据不同的词法单元类型返回相应的样式名称，从而实现代码的高亮显示。
        - split.min.js
          Summarized code for split.min.js

### 代码文件内容总结

#### 文件信息
- **文件名**: `Split.js`
- **版本**: `v1.3.5`

#### 功能概述
该代码文件实现了一个用于创建可调整大小的分割布局的JavaScript库。用户可以通过拖动分割条来调整布局中各个部分的大小。

#### 实现细节

1. **模块导出**
   - 代码首先检查当前环境是CommonJS模块还是AMD模块，或者直接将其挂载到全局对象上。

2. **严格模式**
   - 使用`"use strict";`启用严格模式，以避免一些常见的JavaScript陷阱。

3. **变量声明**
   - 定义了一些全局变量和辅助函数，如`e`（`window`对象）、`t`（`document`对象）、`n`（`addEventListener`）、`i`（`removeEventListener`）等。

4. **辅助函数**
   - `s`：返回`false`的函数，用于阻止默认事件。
   - `l`：根据传入的参数返回DOM元素，支持选择器字符串或直接传入DOM元素。

5. **主函数**
   - 接受两个参数：`u`（包含需要分割的元素的选择器或DOM元素的数组）和`c`（配置对象）。

6. **配置处理**
   - 处理默认配置，如`sizes`、`minSize`、`gutterSize`、`snapOffset`、`direction`、`cursor`等。

7. **CSS处理**
   - 根据配置的方向（`horizontal`或`vertical`）设置相关的CSS属性和计算方式。

8. **分割条处理**
   - 创建分割条元素，并为其添加事件监听器（`mousedown`和`touchstart`）以启动拖动操作。

9. **拖动操作**
   - 在拖动过程中，计算并更新各个分割部分的大小。
   - 拖动开始和结束时分别调用`onDragStart`和`onDragEnd`回调函数。

10. **API方法**
    - `setSizes`：设置各个分割部分的大小。
    - `getSizes`：获取各个分割部分的大小。
    - `collapse`：折叠某个分割部分。
    - `destroy`：销毁分割布局，移除所有分割条和事件监听器。

#### 代码结构

- **模块导出**：根据环境选择合适的导出方式。
- **严格模式**：启用严格模式。
- **辅助函数**：定义了一些辅助函数，如`s`和`l`。
- **主函数**：接受参数并处理分割布局的创建和操作。
- **配置处理**：处理默认配置和用户自定义配置。
- **CSS处理**：根据方向设置CSS属性和计算方式。
- **分割条处理**：创建分割条并添加事件监听器。
- **拖动操作**：处理拖动过程中的大小调整。
- **API方法**：提供设置大小、获取大小、折叠和销毁的方法。

#### 使用示例

```javascript
Split(['#one', '#two'], {
  sizes: [25, 75],
  minSize: 100,
  gutterSize: 10,
  direction: 'horizontal',
  cursor: 'ew-resize'
});
```

#### 总结
该代码文件实现了一个功能强大的分割布局库，支持水平和垂直方向的分割，提供了丰富的配置选项和API方法，方便用户自定义和操作分割布局。
    - favicon.ico
    - index.html
- jsonata-test
  - .gitignore
  - README.md
  - main.go
    Summarized code for main.go

这个代码文件是一个用于测试JSONata表达式的工具。JSONata是一种轻量级的查询和转换JSON数据的语言。以下是代码的详细功能和实现细节：

### 功能概述

1. **命令行工具**：提供一个命令行工具，用于运行JSONata测试用例。
2. **测试用例结构**：定义了一个`testCase`结构体，用于表示单个测试用例。
3. **测试运行**：从指定目录加载测试用例并运行它们，统计通过和失败的测试用例数量。
4. **表达式评估**：评估JSONata表达式，并与预期结果进行比较。
5. **错误处理**：处理加载和评估过程中的错误，并输出详细的错误信息。

### 实现细节

#### 1. 导入包

```go
import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	jsonata "github.com/blues/jsonata-go"
	types "github.com/blues/jsonata-go/jtypes"
)
```

- `encoding/json`：用于JSON的编码和解码。
- `flag`：用于解析命令行参数。
- `fmt`：用于格式化输入输出。
- `io`：用于处理输入输出操作。
- `io/ioutil`：用于读取文件内容。
- `os`：用于操作系统相关的操作。
- `path/filepath`：用于处理文件路径。
- `reflect`：用于反射操作。
- `regexp`：用于正则表达式操作。
- `strings`：用于字符串操作。
- `jsonata`：JSONata库，用于评估JSONata表达式。
- `types`：JSONata库中的类型处理工具。

#### 2. 定义测试用例结构体

```go
type testCase struct {
	Expr        string
	ExprFile    string `json:"expr-file"`
	Category    string
	Data        interface{}
	Dataset     string
	Description string
	TimeLimit   int
	Depth       int
	Bindings    map[string]interface{}
	Result      interface{}
	Undefined   bool
	Error       string `json:"code"`
	Token       string
	Unordered   bool
}
```

- `Expr`：JSONata表达式。
- `ExprFile`：包含JSONata表达式的文件路径。
- `Category`：测试用例的分类。
- `Data`：测试用例的数据。
- `Dataset`：数据集文件的路径。
- `Description`：测试用例的描述。
- `TimeLimit`：时间限制。
- `Depth`：深度。
- `Bindings`：绑定变量。
- `Result`：预期结果。
- `Undefined`：是否未定义。
- `Error`：预期错误代码。
- `Token`：令牌。
- `Unordered`：是否无序。

#### 3. 主函数

```go
func main() {
	var group string
	var verbose bool

	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.StringVar(&group, "group", "", "restrict to one or more test groups")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Syntax: jsonata-test [options] <directory>")
		os.Exit(1)
	}

	root := flag.Arg(0)
	testdir := filepath.Join(root, "groups")
	datadir := filepath.Join(root, "datasets")

	err := run(testdir, datadir, group)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while running: %s\n", err)
		os.Exit(2)
	}

	fmt.Fprintln(os.Stdout, "OK")
}
```

- 解析命令行参数，包括`verbose`和`group`。
- 检查命令行参数是否正确。
- 设置测试目录和数据目录。
- 调用`run`函数运行测试。
- 输出运行结果。

#### 4. 运行测试

```go
func run(testdir string, datadir string, filter string) error {
	var numPassed, numFailed int
	err := filepath.Walk(testdir, func(path string, info os.FileInfo, walkFnErr error) error {
		var dirName string

		if info.IsDir() {
			if path == testdir {
				return nil
			}
			dirName = filepath.Base(path)
			if filter != "" && !strings.Contains(dirName, filter) {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) == ".jsonata" {
			return nil
		}

		testCases, err := loadTestCases(path)
		if err != nil {
			return fmt.Errorf("walk %s: %s", path, err)
		}

		for _, testCase := range testCases {
			failed, err := runTest(testCase, datadir, path)

			if err != nil {
				return err
			}
			if failed {
				numFailed++
			} else {
				numPassed++
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walk %s: ", err)
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, numPassed, "passed", numFailed, "failed")
	return nil
}
```

- 遍历测试目录中的所有文件和子目录。
- 忽略以`.jsonata`结尾的文件。
- 加载测试用例并运行它们。
- 统计通过和失败的测试用例数量。
- 输出测试结果。

#### 5. 运行单个测试用例

```go
func runTest(tc testCase, dataDir string, path string) (bool, error) {
	if tc.Unordered {
		return false, nil
	}

	if tc.TimeLimit != 0 {
		return false, nil
	}

	data := tc.Data
	if tc.Dataset != "" {
		var dest interface{}
		err := readJSONFile(filepath.Join(dataDir, tc.Dataset+".json"), &dest)
		if err != nil {
			return false, err
		}
		data = dest
	}

	var failed bool
	expr, unQuoted := replaceQuotesInPaths(tc.Expr)
	got, _ := eval(expr, tc.Bindings, data)

	if !equalResults(got, tc.Result) {
		failed = true
		printTestCase(os.Stderr, tc, strings.TrimSuffix(filepath.Base(path), ".json"))
		fmt.Fprintf(os.Stderr, "Test file: %s \n", path)

		if tc.Category != "" {
			fmt.Fprintf(os.Stderr, "Category: %s \n", tc.Category)
		}
		if tc.Description != "" {
			fmt.Fprintf(os.Stderr, "Description: %s \n", tc.Description)
		}

		fmt.Fprintf(os.Stderr, "Expression: %s\n", expr)
		if unQuoted {
			fmt.Fprintf(os.Stderr, "Unquoted: %t\n", unQuoted)
		}
		fmt.Fprintf(os.Stderr, "Expected Result: %v [%T]\n", tc.Result, tc.Result)
		fmt.Fprintf(os.Stderr, "Actual Result:   %v [%T]\n", got, got)
	}

	return failed, nil
}
```

- 处理无序测试用例和时间限制。
- 加载数据集文件。
- 替换表达式中的引号。
- 评估表达式并比较结果。
- 输出详细的测试用例信息和结果。

#### 6. 加载测试用例

```go
func loadTestCases(path string) ([]testCase, error) {
	var tc testCase
	err := readJSONFile(path, &tc)
	if err != nil {
		var tcs []testCase
		err := readJSONFile(path, &tcs)
		if err != nil {
			return nil, err
		}

		for _, testCase := range tcs {
			if testCase.ExprFile != "" {
				expr, err := loadTestExprFile(path, testCase.ExprFile)
				if err != nil {
					return nil, err
				}
				testCase.Expr = expr
			}
		}
		return tcs, nil
	}

	if tc.ExprFile != "" {
		expr, err := loadTestExprFile(path, tc.ExprFile)
		if err != nil {
			return nil, err
		}
		tc.Expr = expr
	}

	return []testCase{tc}, nil
}
```

- 尝试将文件内容解码为单个测试用例。
- 如果失败，尝试解码为多个测试用例。
- 加载表达式文件并添加到测试用例中。

#### 7. 辅助函数

- `loadTestExprFile`：加载表达式文件。
- `printTestCase`：输出测试用例信息。
- `eval`：评估JSONata表达式。
- `equalResults`：比较两个结果是否相等。
- `readJSONFile`：读取JSON文件并解码。
- `replaceQuotesInPaths`：替换表达式中的引号。

### 总结

这个代码文件实现了一个用于测试JSONata表达式的命令行工具。它定义了测试用例的结构，加载测试用例并运行它们，评估表达式并与预期结果进行比较，处理错误并输出详细的测试结果。通过阅读和理解这个代码文件，读者可以学习如何使用Go语言编写命令行工具，处理JSON数据，以及评估和测试JSONata表达式。
  - main_test.go
    Summarized code for main_test.go

这个代码文件是一个Go语言的单元测试文件，主要用于测试一个名为`replaceQuotesInPaths`的函数。该函数的功能是将输入字符串中的某些特定字符串路径中的双引号替换为反引号，以确保路径中的特殊字符不会被解释为语法元素。以下是对代码文件的详细解释：

### 文件结构

1. **包声明**：
   ```go
   package main
   ```
   声明该文件属于`main`包。

2. **导入测试包**：
   ```go
   import "testing"
   ```
   导入Go语言的`testing`包，用于编写和运行测试。

### 测试函数

#### 1. `TestReplaceQuotesInPaths`

这个测试函数用于验证`replaceQuotesInPaths`函数是否能正确地将双引号替换为反引号。

- **输入和输出定义**：
  ```go
  inputs := []string{
      // 一系列包含双引号的字符串路径
  }

  outputs := []string{
      // 对应的一系列包含反引号的字符串路径
  }
  ```
  定义了一组输入字符串和期望的输出字符串。

- **测试循环**：
  ```go
  for i := range inputs {
      got, ok := replaceQuotesInPaths(inputs[i])
      if got != outputs[i] {
          t.Errorf("\n     Input: %s\nExp. Output: %s\nAct. Output: %s", inputs[i], outputs[i], got)
      }
      if !ok {
          t.Errorf("%s: Expected true, got %t", inputs[i], ok)
      }
  }
  ```
  遍历输入字符串，调用`replaceQuotesInPaths`函数，并检查返回的结果是否与期望的输出一致。如果结果不一致，则使用`t.Errorf`输出错误信息。

#### 2. `TestReplaceQuotesInPathsNoOp`

这个测试函数用于验证`replaceQuotesInPaths`函数在不需要替换的情况下是否能正确地返回原始输入。

- **输入定义**：
  ```go
  inputs := []string{
      // 一系列不需要替换的字符串
  }
  ```
  定义了一组不需要替换的字符串。

- **测试循环**：
  ```go
  for i := range inputs {
      got, ok := replaceQuotesInPaths(inputs[i])
      if got != inputs[i] {
          t.Errorf("\n     Input: %s\nExp. Output: %s\nAct. Output: %s", inputs[i], inputs[i], got)
      }
      if ok {
          t.Errorf("%s: Expected false, got %t", inputs[i], ok)
      }
  }
  ```
  遍历输入字符串，调用`replaceQuotesInPaths`函数，并检查返回的结果是否与原始输入一致。如果结果不一致，则使用`t.Errorf`输出错误信息。

### 总结

这个代码文件通过两个测试函数`TestReplaceQuotesInPaths`和`TestReplaceQuotesInPathsNoOp`，全面验证了`replaceQuotesInPaths`函数的功能。第一个测试函数确保函数能正确替换双引号为反引号，第二个测试函数确保函数在不需要替换的情况下能正确返回原始输入。通过这些测试，可以确保`replaceQuotesInPaths`函数在各种情况下都能正确工作。
- jsonata.go
  Summarized code for jsonata.go

该代码文件定义了一个用于处理JSONata表达式的Go包，名为`jsonata`。JSONata是一种轻量级的查询和转换JSON数据的语言。以下是对该代码文件的详细功能和实现细节的总结：

### 包导入
首先，代码导入了多个外部包，包括标准库中的`encoding/json`、`fmt`、`reflect`、`sync`、`time`和`unicode`，以及三个外部依赖包`github.com/blues/jsonata-go/jlib`、`github.com/blues/jsonata-go/jparse`和`github.com/blues/jsonata-go/jtypes`。

### 全局变量
定义了两个全局变量：
- `globalRegistryMutex`：一个读写锁，用于保护`globalRegistry`的并发访问。
- `globalRegistry`：一个映射，用于存储全局注册的自定义函数和变量。

### 类型定义
#### `Extension`
`Extension`结构体描述了添加到JSONata表达式中的自定义功能：
- `Func`：一个Go函数，实现自定义功能，返回一个或两个值，第二个值（如果有）必须是错误。
- `UndefinedHandler`：一个处理未定义参数的函数。
- `EvalContextHandler`：一个处理缺失参数的函数。

### 注册函数
#### `RegisterExts`
`RegisterExts`函数用于注册自定义函数，这些函数在程序启动时（例如在`init`函数中）调用一次。注册的函数对所有`Expr`对象可用。

#### `RegisterVars`
`RegisterVars`函数用于注册自定义变量，这些变量在程序启动时（例如在`init`函数中）调用一次。注册的变量对所有`Expr`对象可用。

### `Expr`类型
`Expr`结构体表示一个JSONata表达式：
- `node`：解析后的JSONata表达式节点。
- `registry`：一个映射，存储特定`Expr`对象的自定义函数和变量。

### 编译和执行函数
#### `Compile`
`Compile`函数解析一个JSONata表达式并返回一个`Expr`对象，该对象可以针对JSON数据进行评估。如果输入不是有效的JSONata表达式，则返回错误。

#### `MustCompile`
`MustCompile`函数类似于`Compile`，但如果给定无效表达式，则会引发恐慌。

#### `Eval`
`Eval`方法执行JSONata表达式，针对给定的数据源进行评估。输入通常是解组JSON字符串的结果。输出是一个适合编组为JSON字符串的对象。

#### `EvalBytes`
`EvalBytes`方法类似于`Eval`，但它接受和返回字节切片而不是对象。

### 注册方法
#### `RegisterExts`
`RegisterExts`方法为特定`Expr`对象注册自定义函数。

#### `RegisterVars`
`RegisterVars`方法为特定`Expr`对象注册自定义变量。

### 辅助方法
#### `String`
`String`方法返回`Expr`对象的字符串表示。

#### `updateRegistry`
`updateRegistry`方法更新`Expr`对象的注册表。

#### `newEnv`
`newEnv`方法创建一个新的评估环境。

### 时间相关函数
定义了两个时间相关的函数：
- `milisT`：返回当前时间的毫秒数。
- `nowT`：根据给定的毫秒数、格式和时区返回当前时间。

### 处理函数
#### `processExts`
`processExts`函数处理自定义扩展，确保它们是有效的函数。

#### `processVars`
`processVars`函数处理自定义变量，确保它们是有效的变量。

### 全局注册表更新
#### `updateGlobalRegistry`
`updateGlobalRegistry`函数更新全局注册表。

### 辅助函数
#### `validName`
`validName`函数验证名称是否有效。

#### `validVar`
`validVar`函数验证变量是否有效（TODO：待实现）。

#### `isLetter`
`isLetter`函数检查字符是否为字母。

#### `isDigit`
`isDigit`函数检查字符是否为数字。

### 总结
该代码文件实现了一个用于处理JSONata表达式的Go包，提供了编译、执行和注册自定义函数及变量的功能。通过全局和局部注册表，支持在不同层级上注册自定义功能，并提供了时间相关的内置函数。
- jsonata_test.go
  [File too long to summarize]
- jtypes
  - funcs.go
    Summarized code for funcs.go

这个代码文件定义了一个名为 `jtypes` 的 Go 包，主要用于通过反射机制检查和处理不同类型的 Go 值。以下是对该代码文件的详细解析：

### 包声明和导入
```go
// Package jtypes (golint)
package jtypes

import (
	"reflect"
)
```
- **包声明**：定义了包名为 `jtypes`。
- **导入**：导入了 `reflect` 包，用于 Go 语言的反射机制。

### 函数 `Resolve`
```go
// Resolve (golint)
func Resolve(v reflect.Value) reflect.Value {
	for {
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr:
			if !v.IsNil() {
				v = v.Elem()
				break
			}
			fallthrough
		default:
			return v
		}
	}
}
```
- **功能**：该函数用于解析 `reflect.Value` 类型的值，直到解析到非指针或非接口类型。
- **实现细节**：
  - 使用无限循环 `for` 循环，直到返回值。
  - 通过 `switch` 语句检查 `v.Kind()` 是否为 `reflect.Interface` 或 `reflect.Ptr`。
  - 如果 `v` 不是 `nil`，则通过 `v.Elem()` 获取其指向的值，并继续循环。
  - 如果 `v` 是 `nil` 或其他类型，则直接返回 `v`。

### 类型检查函数
以下函数用于检查 `reflect.Value` 类型的值是否为特定类型：

#### `IsBool`
```go
// IsBool (golint)
func IsBool(v reflect.Value) bool {
	return v.Kind() == reflect.Bool || resolvedKind(v) == reflect.Bool
}
```
- **功能**：检查 `v` 是否为布尔类型。

#### `IsString`
```go
// IsString (golint)
func IsString(v reflect.Value) bool {
	return v.Kind() == reflect.String || resolvedKind(v) == reflect.String
}
```
- **功能**：检查 `v` 是否为字符串类型。

#### `IsNumber`
```go
// IsNumber (golint)
func IsNumber(v reflect.Value) bool {
	return isFloat(v) || isInt(v) || isUint(v)
}
```
- **功能**：检查 `v` 是否为数字类型（包括浮点数、整数和无符号整数）。

#### `IsCallable`
```go
// IsCallable (golint)
func IsCallable(v reflect.Value) bool {
	v = Resolve(v)
	return v.IsValid() &&
		(v.Type().Implements(TypeCallable) || reflect.PtrTo(v.Type()).Implements(TypeCallable))
}
```
- **功能**：检查 `v` 是否为可调用类型（即实现了 `TypeCallable` 接口）。

#### `IsArray`
```go
// IsArray (golint)
func IsArray(v reflect.Value) bool {
	return isArrayKind(v.Kind()) || isArrayKind(resolvedKind(v))
}

func isArrayKind(k reflect.Kind) bool {
	return k == reflect.Slice || k == reflect.Array
}
```
- **功能**：检查 `v` 是否为数组或切片类型。

#### `IsArrayOf`
```go
// IsArrayOf (golint)
func IsArrayOf(v reflect.Value, hasType func(reflect.Value) bool) bool {
	if !IsArray(v) {
		return false
	}

	v = Resolve(v)
	for i := 0; i < v.Len(); i++ {
		if !hasType(v.Index(i)) {
			return false
		}
	}

	return true
}
```
- **功能**：检查 `v` 是否为特定类型的数组或切片。

#### `IsMap`
```go
// IsMap (golint)
func IsMap(v reflect.Value) bool {
	return resolvedKind(v) == reflect.Map
}
```
- **功能**：检查 `v` 是否为映射类型。

#### `IsStruct`
```go
// IsStruct (golint)
func IsStruct(v reflect.Value) bool {
	return resolvedKind(v) == reflect.Struct
}
```
- **功能**：检查 `v` 是否为结构体类型。

### 类型转换函数
以下函数用于将 `reflect.Value` 类型的值转换为特定类型：

#### `AsBool`
```go
// AsBool (golint)
func AsBool(v reflect.Value) (bool, bool) {
	v = Resolve(v)

	switch {
	case IsBool(v):
		return v.Bool(), true
	default:
		return false, false
	}
}
```
- **功能**：尝试将 `v` 转换为布尔类型，并返回转换结果和是否成功。

#### `AsString`
```go
// AsString (golint)
func AsString(v reflect.Value) (string, bool) {
	v = Resolve(v)

	switch {
	case IsString(v):
		return v.String(), true
	default:
		return "", false
	}
}
```
- **功能**：尝试将 `v` 转换为字符串类型，并返回转换结果和是否成功。

#### `AsNumber`
```go
// AsNumber (golint)
func AsNumber(v reflect.Value) (float64, bool) {
	v = Resolve(v)

	switch {
	case isFloat(v):
		return v.Float(), true
	case isInt(v), isUint(v):
		return v.Convert(typeFloat64).Float(), true
	default:
		return 0, false
	}
}
```
- **功能**：尝试将 `v` 转换为浮点数类型，并返回转换结果和是否成功。

#### `AsCallable`
```go
// AsCallable (golint)
func AsCallable(v reflect.Value) (Callable, bool) {
	v = Resolve(v)

	if v.IsValid() && v.Type().Implements(TypeCallable) && v.CanInterface() {
		return v.Interface().(Callable), true
	}

	if v.IsValid() && reflect.PtrTo(v.Type()).Implements(TypeCallable) && v.CanAddr() && v.Addr().CanInterface() {
		return v.Addr().Interface().(Callable), true
	}

	return nil, false
}
```
- **功能**：尝试将 `v` 转换为可调用类型，并返回转换结果和是否成功。

### 辅助函数
以下函数用于检查特定类型的子类型：

#### `isInt`
```go
func isInt(v reflect.Value) bool {
	return isIntKind(v.Kind()) || isIntKind(resolvedKind(v))
}

func isIntKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}
```
- **功能**：检查 `v` 是否为整数类型。

#### `isUint`
```go
func isUint(v reflect.Value) bool {
	return isUintKind(v.Kind()) || isUintKind(resolvedKind(v))
}

func isUintKind(k reflect.Kind) bool {
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}
```
- **功能**：检查 `v` 是否为无符号整数类型。

#### `isFloat`
```go
func isFloat(v reflect.Value) bool {
	return isFloatKind(v.Kind()) || isFloatKind(resolvedKind(v))
}

func isFloatKind(k reflect.Kind) bool {
	switch k {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
```
- **功能**：检查 `v` 是否为浮点数类型。

#### `resolvedKind`
```go
func resolvedKind(v reflect.Value) reflect.Kind {
	return Resolve(v).Kind()
}
```
- **功能**：返回 `v` 解析后的类型。

### 总结
该代码文件通过反射机制提供了丰富的类型检查和转换功能，适用于需要动态处理不同类型值的场景。通过这些函数，可以方便地检查和转换 `reflect.Value` 类型的值，提高了代码的灵活性和可维护性。
  - types.go
    Summarized code for types.go

### 代码文件内容总结

#### 1. 包声明和导入
```go
// Package jtypes provides types and utilities for third party extension functions.
package jtypes

import (
	"errors"
	"reflect"
)
```
- **包声明**：定义了包名为 `jtypes`，用于提供第三方扩展函数所需的类型和工具。
- **导入包**：导入了 `errors` 和 `reflect` 包，用于处理错误和反射操作。

#### 2. 常量和变量定义
```go
var undefined reflect.Value

var (
	typeBool    = reflect.TypeOf((*bool)(nil)).Elem()
	typeInt     = reflect.TypeOf((*int)(nil)).Elem()
	typeFloat64 = reflect.TypeOf((*float64)(nil)).Elem()
	typeString  = reflect.TypeOf((*string)(nil)).Elem()

	TypeOptional    = reflect.TypeOf((*Optional)(nil)).Elem()
	TypeCallable    = reflect.TypeOf((*Callable)(nil)).Elem()
	TypeConvertible = reflect.TypeOf((*Convertible)(nil)).Elem()
	TypeVariant     = reflect.TypeOf((*Variant)(nil)).Elem()
	TypeValue       = reflect.TypeOf((*reflect.Value)(nil)).Elem()
	TypeInterface   = reflect.TypeOf((*interface{})(nil)).Elem()
)

var ErrUndefined = errors.New("undefined")
```
- **undefined**：定义了一个未定义的 `reflect.Value`。
- **类型定义**：使用 `reflect.TypeOf` 定义了一系列基本类型和接口类型的反射类型。
- **ErrUndefined**：定义了一个错误常量 `ErrUndefined`，表示未定义的错误。

#### 3. 接口定义
```go
type Variant interface {
	ValidTypes() []reflect.Type
}

type Callable interface {
	Name() string
	ParamCount() int
	Call([]reflect.Value) (reflect.Value, error)
}

type Convertible interface {
	ConvertTo(reflect.Type) (reflect.Value, bool)
}

type Optional interface {
	IsSet() bool
	Set(reflect.Value)
	Type() reflect.Type
}
```
- **Variant**：定义了一个接口，包含 `ValidTypes` 方法，返回有效类型的列表。
- **Callable**：定义了一个接口，包含 `Name`、`ParamCount` 和 `Call` 方法，用于可调用对象。
- **Convertible**：定义了一个接口，包含 `ConvertTo` 方法，用于类型转换。
- **Optional**：定义了一个接口，包含 `IsSet`、`Set` 和 `Type` 方法，用于可选值。

#### 4. 结构体和方法定义
```go
type isSet bool

func (opt *isSet) IsSet() bool {
	return bool(*opt)
}
```
- **isSet**：定义了一个布尔类型的结构体，用于表示是否设置。
- **IsSet**：实现了 `Optional` 接口的 `IsSet` 方法。

```go
type OptionalBool struct {
	isSet
	Bool bool
}

func NewOptionalBool(value bool) OptionalBool {
	opt := OptionalBool{}
	opt.Set(reflect.ValueOf(value))
	return opt
}

func (opt *OptionalBool) Set(v reflect.Value) {
	opt.isSet = true
	opt.Bool = v.Bool()
}

func (opt *OptionalBool) Type() reflect.Type {
	return typeBool
}
```
- **OptionalBool**：定义了一个可选的布尔类型结构体，包含 `isSet` 和 `Bool` 字段。
- **NewOptionalBool**：构造函数，创建并设置 `OptionalBool` 实例。
- **Set**：实现了 `Optional` 接口的 `Set` 方法。
- **Type**：实现了 `Optional` 接口的 `Type` 方法。

类似地，定义了 `OptionalInt`、`OptionalFloat64`、`OptionalString`、`OptionalInterface`、`OptionalValue` 和 `OptionalCallable` 结构体，分别对应不同的可选类型。

#### 5. 函数定义
```go
type ArgHandler func([]reflect.Value) bool

func ArgCountEquals(n int) ArgHandler {
	return func(argv []reflect.Value) bool {
		return len(argv) == n
	}
}

func ArgUndefined(i int) ArgHandler {
	return func(argv []reflect.Value) bool {
		return len(argv) > i && argv[i] == undefined
	}
}
```
- **ArgHandler**：定义了一个函数类型，用于处理参数。
- **ArgCountEquals**：返回一个 `ArgHandler`，检查参数数量是否等于 `n`。
- **ArgUndefined**：返回一个 `ArgHandler`，检查指定位置的参数是否未定义。

### 功能和实现细节

1. **类型定义**：
   - 使用 `reflect.TypeOf` 定义了一系列基本类型和接口类型的反射类型，便于后续的反射操作。

2. **接口定义**：
   - 定义了多个接口，如 `Variant`、`Callable`、`Convertible` 和 `Optional`，用于抽象不同类型的行为。

3. **结构体和方法定义**：
   - 定义了多个可选类型的结构体，如 `OptionalBool`、`OptionalInt` 等，实现了 `Optional` 接口的方法。
   - 使用 `isSet` 布尔类型结构体来表示是否设置，简化了实现。

4. **函数定义**：
   - 定义了 `ArgHandler` 函数类型，用于处理参数。
   - 提供了 `ArgCountEquals` 和 `ArgUndefined` 函数，用于检查参数数量和是否未定义。

通过这些定义和实现，`jtypes` 包提供了一套灵活的类型和工具，便于第三方扩展函数的使用和处理。
- testdata
  - account.json
  - account2.json
  - account3.json
  - account4.json
  - account5.json
  - account6.json
  - account7.json
  - address.json
  - foobar.json
  - foobar2.json
  - library.json
  - nest1.json
  - nest2.json
  - nest3.json
