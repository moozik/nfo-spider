package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

func PathConvert(dirRoot, fileName string) string {
	return path.Join(dirRoot, "images", fileName)
}

func ImageDownload(localPath, url string) {
	log.Println("ImageDownload", localPath, url)
	if Exists(localPath) {
		return
	}
	if url == "" {
		log.Println("url为空")
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("图片请求失败,", err)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("图片读取失败,", err)
		return
	}
	err = ioutil.WriteFile(localPath, data, 0644)
	if err != nil {
		log.Fatal("图片下载失败,", err)
		return
	}
}
