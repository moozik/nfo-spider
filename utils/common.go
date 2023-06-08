package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Encode(a any) string {
	v, _ := json.Marshal(a)
	return string(v)
}

func Url(host, src string) string {
	if src != "" {
		if src[:4] == "http" {
			return src
		}
		if src[0] == '/' {
			return host + src
		}
		return ""
	}
	return src
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

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func IsRelease() bool {
	return !strings.Contains(os.Args[0], "go-build")
}
