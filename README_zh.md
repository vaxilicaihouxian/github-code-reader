# github-code-reader
LLM辅助阅读Github代码库
1. 当前仅支持deepseek llm阅读代码
## 准备
1. github token
2. deepseek api key（这个比较便宜，注册方便）: https://platform.deepseek.com/
## 用法

```shell
export GITHUB_TOKEN=你的githubToken，需要自己去github生成，用于读取github上的库
export OPENAI_API_KEY=你的deepseek平台api key，当前只实现了用deepseek llm 阅读代码的功能
go run main.go summary "https://github.com/{owner}/{仓库名}"
```

会在当前目录下生成一个.txt文件，包含代码解读结果

### 示例
```
go run main.go summary "https://github.com/xxl-job/xxl-job-executor-go"
```

输出文件**code_summary.txt**，文件大概内容如下
```md
### Repository Summary

#### Purpose
# xxl-job-executor-go
很多公司java与go开发共存，java中有xxl-job做为任务调度引擎，为此也出现了go执行器(客户端)，使用起来比较简单：
# 支持
1.执行器注册
2.耗时任务取消
...
### 功能和实现细节
- **常量声明**:
  - 使用 `const` 关键字定义常量。常量在声明后不能被修改，这有助于确保代码的安全性和稳定性。
  - 常量名采用驼峰命名法，`SuccessCode` 和 `FailureCode` 都是自描述的名称，便于理解其用途。

- **响应码**:
  - 这两个常量通常用于网络请求的响应处理中，帮助开发者快速判断请求的结果是成功还是失败。
  - 在实际应用中，可以根据这些常量来决定后续的业务逻辑，例如在成功时执行某些操作，在失败时进行错误处理或重试。

### 示例应用
假设有一个处理 HTTP 请求的函数，可以使用这些常量来返回响应码：
...
```
更多代码阅读解析md文件示例，请查看**examples/reader-markdowns/{repo}/code_summary.md**