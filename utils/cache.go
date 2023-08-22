package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Cache struct {
	SiteName string
	DataType string
	Tail     string
}

func NewCache(siteName, dataType, tail string) *Cache {
	return &Cache{
		SiteName: siteName,
		DataType: dataType,
		Tail:     tail,
	}
}

func (c *Cache) FilePath(key string) string {
	dir := GetCurrentDirectory()
	if !IsRelease() {
		dir = `F:\github\nfo-spider`
	}
	return filepath.Join(dir, string(os.PathSeparator), "cache", string(os.PathSeparator), fmt.Sprintf(`%s_%s_%s.%s`, c.SiteName, c.DataType, key, c.Tail))
}

func (c *Cache) Set(key string, data []byte) error {

	var file *os.File
	var err error
	fileName := c.FilePath(key)
	//文件是否存在
	if Exists(fileName) {
		//使用追加模式打开文件
		file, err = os.OpenFile(fileName, os.O_APPEND, 0666)
		if err != nil {
			log.Println("打开文件错误：", err)
			return err
		}
	} else {
		//不存在创建文件
		file, err = os.Create(fileName)
		if err != nil {
			log.Println("创建失败", err)
			return err
		}
	}

	defer file.Close()
	//写入文件
	_, err = file.Write(data)
	if err != nil {
		log.Println("写入错误：", err)
		return err
	}
	return nil
}

func (c *Cache) Get(key string) []byte {
	if !Exists(c.FilePath(key)) {
		return nil
	}
	byteData, _ := os.ReadFile(c.FilePath(key))
	return byteData
}
