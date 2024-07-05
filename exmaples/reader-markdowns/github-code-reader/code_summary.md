### Repository Summary

#### Purpose
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
go run main.go "https://github.com/{owner}/{仓库名}"
```

会在当前目录下生成一个.txt文件，包含代码解读结果

### 示例
```
go run main.go "https://github.com/xxl-job/xxl-job-executor-go"
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

#### Structure
- .gitignore
- README.md
- go.mod
- main.go
  Summarized code for main.go

这段代码的主要功能是从GitHub仓库中获取文件内容，并生成仓库的结构和目的的总结。以下是对代码的详细解释：

### 1. 导入包
```go
import (
	"context"
	"fmt"
	"github.com/google/go-github/v45/github"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2"
	"os"
	"strconv"
	"strings"
)
```
- `context`: 用于处理上下文。
- `fmt`: 用于格式化输入输出。
- `github.com/google/go-github/v45/github`: 用于与GitHub API交互。
- `github.com/sashabaranov/go-openai`: 用于与OpenAI API交互。
- `golang.org/x/oauth2`: 用于处理OAuth2认证。
- `os`: 用于与操作系统交互。
- `strconv`: 用于字符串转换。
- `strings`: 用于字符串操作。

### 2. 全局变量
```go
var MaxLLMInputLength = 4096
var LLMAuthorizationToken = os.Getenv("OPENAI_API_KEY")
```
- `MaxLLMInputLength`: 定义了LLM输入的最大长度。
- `LLMAuthorizationToken`: 从环境变量中获取OpenAI API的授权令牌。

### 3. 数据结构
```go
// GitHubRepo represents the basic information about a GitHub repository.
type GitHubRepo struct {
	Owner string
	Name  string
}

// FileContent represents the content of a file in a repository.
type FileContent struct {
	Name    string
	Content string
}
```
- `GitHubRepo`: 表示GitHub仓库的基本信息。
- `FileContent`: 表示仓库中文件的内容。

### 4. 函数
#### 4.1. `GetRepositoryContents`
```go
// GetRepositoryContents retrieves the contents of a GitHub repository at a specific path.
func GetRepositoryContents(client *github.Client, repo GitHubRepo, path string) ([]*github.RepositoryContent, error) {
	_, dirContents, _, err := client.Repositories.GetContents(context.Background(), repo.Owner, repo.Name, path, nil)
	if err != nil {
		return nil, err
	}
	return dirContents, nil
}
```
- 从GitHub仓库中获取指定路径的内容。

#### 4.2. `GetFileContent`
```go
// GetFileContent retrieves the content of a file in a repository.
func GetFileContent(client *github.Client, repo GitHubRepo, filePath string) (string, error) {
	fileContent, _, _, err := client.Repositories.GetContents(context.Background(), repo.Owner, repo.Name, filePath, nil)
	if err != nil {
		return "", err
	}
	content, err := fileContent.GetContent()
	if err != nil {
		return "", err
	}
	return content, nil
}
```
- 从GitHub仓库中获取指定文件的内容。

#### 4.3. `SummarizeRepository`
```go
// SummarizeRepository generates a summary of the repository's structure and purpose.
func SummarizeRepository(readmeContent string, repo GitHubRepo, client *github.Client) (string, error) {
	structure, err := getStructure(client, repo, "", 0, MaxLLMInputLength)
	if err != nil {
		return "", err
	}

	summary := "### Repository Summary\n\n"
	summary += "#### Purpose\n"
	summary += readmeContent + "\n\n"
	summary += "#### Structure\n"
	summary += structure

	return summary, nil
}
```
- 生成仓库的结构和目的的总结。

#### 4.4. `getStructure`
```go
// getStructure recursively retrieves the directory structure and returns it as a formatted string.
func getStructure(client *github.Client, repo GitHubRepo, path string, level int, maxLength int) (string, error) {
	contents, err := GetRepositoryContents(client, repo, path)
	if err != nil {
		return "", err
	}

	var structure strings.Builder
	indent := strings.Repeat("  ", level)
	for _, content := range contents {
		structure.WriteString(fmt.Sprintf("%s- %s\n", indent, content.GetName()))
		if content.GetType() == "dir" {
			subStructure, err := getStructure(client, repo, content.GetPath(), level+1, maxLength)
			if err != nil {
				return "", err
			}
			structure.WriteString(subStructure)
		} else {
			if isCodeFile(content.GetName()) {
				code, err := GetFileContent(client, repo, content.GetPath())
				if err != nil {
					return "", err
				}
				if len(code) > maxLength {
					structure.WriteString(fmt.Sprintf("%s  [File too long to summarize]\n", indent))
				} else {
					summary := summarizeCode(content.GetName(), code)
					structure.WriteString(fmt.Sprintf("%s  %s\n", indent, summary))
				}
			}
		}
	}

	return structure.String(), nil
}
```
- 递归地获取目录结构并返回格式化的字符串。

#### 4.5. `main`
```go
func main() {
	// read repo url from command argument
	repoURL := os.Args[1]
	// parse repo url to GitHubRepo
	repo, err0 := ParseGitHubURL(repoURL)
	if err0 != nil {
		fmt.Printf("Error parsing repo url: %v\n", err0)
		return
	}
	// Replace with your GitHub token
	token := os.Getenv("GITHUB_TOKEN")
	maxLLMInputLength := os.Getenv("LLM_MAX_INPUT_LENGTH")
	if maxLLMInputLength != "" {
		// convert maxLLMInputLength to int
		maxInputLength, err := strconv.Atoi(maxLLMInputLength)
		if err != nil {
			fmt.Printf("Error converting max input length: %v\n", err)
			return
		}
		MaxLLMInputLength = maxInputLength
	}

	if LLMAuthorizationToken == "" {
		fmt.Println("No LLM authorization token provided")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	readmeContent, err := GetFileContent(client, repo, "README.md")
	if err != nil {
		fmt.Printf("Error getting README content: %v\n", err)
		return
	}

	summary, err := SummarizeRepository(readmeContent, repo, client)
	if err != nil {
		fmt.Printf("Error summarizing repository: %v\n", err)
		return
	}

	fmt.Println(summary)
	f, err := os.Create("./code_summary.txt")
	f.WriteString(summary)
	f.Close()
}
```
- 主函数，从命令行参数读取仓库URL，解析并获取仓库内容，生成总结并保存到文件。

#### 4.6. `ParseGitHubURL`
```go
// ParseGitHubURL parses a GitHub repository URL and extracts the owner and repo name.
func ParseGitHubURL(url string) (GitHubRepo, error) {
	parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
	if len(parts) != 2 {
		return GitHubRepo{}, fmt.Errorf("invalid GitHub URL: %s", url)
	}
	return GitHubRepo{Owner: parts[0], Name: parts[1]}, nil
}
```
- 解析GitHub仓库URL并提取所有者和仓库名称。

#### 4.7. `isCodeFile`
```go
// isCodeFile checks if a file is a code file based on its extension.
func isCodeFile(filename string) bool {
	extensions := []string{".go", ".c", ".js", ".ts", ".php", ".cpp"}
	for _, ext := range extensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
```
- 检查文件是否为代码文件基于其扩展名。

#### 4.8. `summarizeCode`
```go
// summarizeCode is a placeholder function to simulate calling an LLM to summarize code.
func summarizeCode(fileName, code string) string {
	fmt.Println("Summarizing code for", fileName)
	config := openai.DefaultConfig(LLMAuthorizationToken)
	config.BaseURL = "https://api.deepseek.com"
	client := openai.NewClientWithConfig(config)
	prompt := `总结以下代码文件内容，尽可能详细讲解功能和实现细节，便于读者学习阅读该代码: %s`
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "deepseek-chat",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant.Use Chinese to answer questions.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(prompt, code),
				},
			},
		},
	)

	if err != nil {
		return fmt.Sprintf("ChatCompletion error: %v\n", err)
	}

	// Simulate LLM call - replace with actual LLM call logic
	result := fmt.Sprintf("Summarized code for %s\n\n%s", fileName, resp.Choices[0].Message.Content)
	fmt.Println("process result:", "---------------\n\n", result, "\n\n---------------\n\n")
	return result
}
```
- 模拟调用LLM来总结代码文件内容。

### 总结
这段代码的主要功能是从GitHub仓库中获取文件内容，并生成仓库的结构和目的的总结。它通过递归地获取目录结构，并使用OpenAI API来总结代码文件内容。最终，生成的总结会被保存到一个文件中。
