package utils

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/moozik/nfo-spider/define"
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

func XMLBuild(dirRoot, avPath string, a *define.AvItem) {
	log.Printf("dirRoot:%s,avPath:%s", dirRoot, avPath)
	if a == nil {
		log.Fatal("nil input *define.AvItem")
		return
	}
	avFileName := path.Base(avPath)
	patternFileRet := define.PattrenAvFile.FindStringSubmatch(avFileName)
	if len(patternFileRet) != 3 {
		log.Printf("文件名错误:%s", avFileName)
		return
	}

	filePath := path.Join(dirRoot, "images", "poster_"+a.AvId+".jpg")
	ImageDownload(filePath, a.Poster)

	ret := define.NfoMovie{
		Title:         a.Title,
		Originaltitle: a.Title,
		Sorttitle:     a.AvId,
		Premiered:     a.ReleaseDate,
		Releasedate:   a.ReleaseDate,
		Writer:        a.Series,
		Credits:       a.Series,
		Genre:         a.Genre,
		Art: define.Art{
			Poster: filePath,
			Fanart: filePath,
		},
		Fanart: []define.Fanart{
			{
				Thumb: filePath,
			},
		},
		Actor: []define.Actor{},
	}

	for _, act := range a.Stars {
		filePath := path.Join(dirRoot, "images", "actor_"+act.Name+".jpg")
		ImageDownload(filePath, act.Image)
		ret.Actor = append(ret.Actor, define.Actor{
			Name:  act.Name,
			Type:  "actor",
			Thumb: filePath,
		})
	}

	xmlByte, err := xml.Marshal(ret)
	if err != nil {
		log.Fatalf("xml.Marshal fail,err:%v", err)
		return
	}

	filePath = path.Join(path.Dir(avPath), patternFileRet[1]+".nfo")
	if Exists(filePath) {
		os.Remove(filePath)
	}
	err = ioutil.WriteFile(filePath, append([]byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>`), xmlByte...), 0644)
	if err != nil {
		log.Fatal("nfo文件写入失败,", err)
		return
	}
}
