package main

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

var MaxLLMInputLength = 4096
var LLMAuthorizationToken = os.Getenv("OPENAI_API_KEY")

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

// GetRepositoryContents retrieves the contents of a GitHub repository at a specific path.
func GetRepositoryContents(client *github.Client, repo GitHubRepo, path string) ([]*github.RepositoryContent, error) {
	_, dirContents, _, err := client.Repositories.GetContents(context.Background(), repo.Owner, repo.Name, path, nil)
	if err != nil {
		return nil, err
	}
	return dirContents, nil
}

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

//// getStructure recursively retrieves the directory structure and returns it as a formatted string.
//func getStructure(client *github.Client, repo GitHubRepo, path string, level int) (string, error) {
//	contents, err := GetRepositoryContents(client, repo, path)
//	if err != nil {
//		return "", err
//	}
//
//	var structure strings.Builder
//	indent := strings.Repeat("  ", level)
//	for _, content := range contents {
//		structure.WriteString(fmt.Sprintf("%s- %s\n", indent, content.GetName()))
//		if content.GetType() == "dir" {
//			subStructure, err := getStructure(client, repo, content.GetPath(), level+1)
//			if err != nil {
//				return "", err
//			}
//			structure.WriteString(subStructure)
//		}
//	}
//
//	return structure.String(), nil
//}

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

// ParseGitHubURL parses a GitHub repository URL and extracts the owner and repo name.
func ParseGitHubURL(url string) (GitHubRepo, error) {
	parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
	if len(parts) != 2 {
		return GitHubRepo{}, fmt.Errorf("invalid GitHub URL: %s", url)
	}
	return GitHubRepo{Owner: parts[0], Name: parts[1]}, nil
}

// isCodeFile checks if a file is a code file based on its extension.
func isCodeFile(filename string) bool {
	extensions := []string{".go", ".py", ".c", ".js", ".ts", ".php", ".cpp"}
	for _, ext := range extensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

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
			Model: "deepseek-coder",
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
