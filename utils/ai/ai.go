package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/moozik/nfo-spider/utils"
	"io"
	"net/http"
	"time"
)

// 定义请求和响应的结构体
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func Ask(content string) string {
	// 1. 构造请求数据
	requestData := ChatRequest{
		Model:  utils.GetEnv("MODEL"),
		Stream: false,
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	return Request(requestData)
}

func Translate(content string) string {
	// 1. 构造请求数据
	requestData := ChatRequest{
		Model:  utils.GetEnv("MODEL"),
		Stream: false,
		Messages: []Message{
			{
				Role:    "system",
				Content: "翻译日语为简体中文，输出单行结果`",
			},
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	return Request(requestData)
}

func Request(requestData ChatRequest) string {

	// 2. 序列化为 JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		panic(fmt.Sprintf("JSON marshaling failed: %v", err))
	}

	// 3. 创建 HTTP 请求
	client := &http.Client{Timeout: 60 * time.Second} // 设置超时，大模型响应可能较慢
	req, err := http.NewRequest("POST", utils.GetEnv("AI_PATH"), bytes.NewBuffer(jsonData))
	if err != nil {
		panic(fmt.Sprintf("Creating request failed: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")

	// 4. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Request failed: %v", err))
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 5. 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("Request failed with status: %s, body: %s", resp.Status, string(body)))
	}

	// 6. 解析响应
	var chatResponse ChatResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&chatResponse)
	if err != nil {
		panic(fmt.Sprintf("Parsing response failed: %v", err))
	}

	return chatResponse.Message.Content
}
