package av_base

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/moozik/nfo-spider/define"
	"github.com/moozik/nfo-spider/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type avInter interface {
	AppName() string
	Host() string
	GetOne(string) *AvItem
}

type AvBase struct {
	avInter
}

func (a *AvBase) DocGet(url string) (*goquery.Document, error) {
	log.Printf("httpGet,url:%s\n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status:" + resp.Status)
	}
	return goquery.NewDocumentFromReader(resp.Body)
}

func GetOneWithCache(s avInter, id string) *AvItem {
	c := utils.NewCache(s.AppName(), define.CacheTypeAvItem)
	b := c.Get(id)
	var avItem AvItem
	if b != nil {
		_ = json.Unmarshal(b, &avItem)
		log.Printf("id:%s,read cache\n", id)
		return &avItem
	}
	pAvItem := s.GetOne(id)
	b, _ = json.Marshal(pAvItem)
	log.Printf("id:%s,save cache\n", id)
	c.Set(id, b)
	return pAvItem
}

func XMLBuild(dirRoot, avPath string, a *AvItem) {
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
	utils.ImageDownload(filePath, a.Poster)

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
		imagePath := path.Join(dirRoot, "images", "actor_"+act.Name+".jpg")
		utils.ImageDownload(imagePath, act.Image)
		ret.Actor = append(ret.Actor, define.Actor{
			Name:  act.Name,
			Type:  "actor",
			Thumb: imagePath,
		})
	}

	xmlByte, err := xml.Marshal(ret)
	if err != nil {
		log.Fatalf("xml.Marshal fail,err:%v", err)
		return
	}

	filePath = path.Join(path.Dir(avPath), patternFileRet[1]+".nfo")
	if utils.Exists(filePath) {
		os.Remove(filePath)
	}
	err = ioutil.WriteFile(filePath, append([]byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>`), xmlByte...), 0644)
	if err != nil {
		log.Fatal("nfo文件写入失败,", err)
		return
	}
}
