package utils

import (
	"os"
)

// defaults 存储环境变量默认值，程序启动初始化后应保持只读
var defaults = map[string]string{
	"JAVBUS_HOST": "www.seedmm.cyou",
	"IMAGE_HOST":  "http://awsimgsrc.dmm.co.jp/pics_dig",
	"PROXY":       "",                                        //socks5://127.0.0.1:1080
	"MODEL":       "huihui_ai/hunyuan-mt-abliterated:latest", //模型名字 deepseek-r1:8b hunyuan-mt-abliterated:latest
	"AI_PATH":     "http://192.168.31.61:11434/api/chat",     //AI路径
}

// GetEnv 获取环境变量值，遵循以下优先级：
// 1. 存在环境变量时返回环境变量值（包括空值）
// 2. 不存在时返回预设的默认值
// 3. 都没有时返回空字符串
func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaults[key]
}
