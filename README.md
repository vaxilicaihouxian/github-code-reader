# GitHub Code Reader

LLM-assisted reading of GitHub code repositories.

## Features

- Currently supports code reading using DeepSeek LLM.

## Prerequisites

1. GitHub token
2. DeepSeek API key (affordable and easy to register): [DeepSeek Platform](https://platform.deepseek.com/)

## Usage

```shell
export GITHUB_TOKEN=your_github_token # Generate this token on GitHub to access repositories
export OPENAI_API_KEY=your_deepseek_api_key # Currently supports DeepSeek LLM for code reading
go run main.go "https://github.com/{owner}/{repository_name}"
```

A `.txt` file containing the code interpretation results will be generated in the current directory.

### Example

```shell
go run main.go "https://github.com/xxl-job/xxl-job-executor-go"
```

Output file: **code_summary.txt**

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

For more examples of code reading and analysis in Markdown files, please refer to **examples/reader-markdowns/{repo}/code_summary.md**.