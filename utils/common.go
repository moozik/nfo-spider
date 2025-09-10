package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var client *http.Client

func init() {
	var transport *http.Transport
	if GetEnv("PROXY") != "" {
		proxy, err := url.Parse(GetEnv("PROXY"))
		if err != nil {
			log.Println("代理地址解析失败:", err)
			return
		}

		transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}

	client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func EncodeString(a any) string {
	v, _ := json.Marshal(a)
	return string(v)
}

func Encode(a any) []byte {
	v, _ := json.Marshal(a)
	return v
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

var currentPath string

func GetCurrentDirectory() string {
	if currentPath != "" {
		return currentPath
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("dir:", dir)
	currentPath = strings.Replace(dir, "\\", "/", -1)
	log.Println("currentPath:", dir)
	return currentPath
}

func IsRelease() bool {
	return !strings.Contains(os.Args[0], "go-build")
}

func ImageDownload(localPath, url string) {
	if Exists(localPath) {
		return
	}
	if url == "" {
		log.Println("ImageDownload url为空")
		return
	}
	log.Println("ImageDownload", localPath, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("图片请求失败,", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("图片读取失败,", err)
		return
	}

	err = os.WriteFile(localPath, data, 0644)
	if err != nil {
		log.Println("图片下载失败,", err)
		return
	}
}
