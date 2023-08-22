package av_base

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/moozik/nfo-spider/define"
	"github.com/moozik/nfo-spider/utils"
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
	if err != nil {
		return nil, err
	}
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
	c := utils.NewCache(s.AppName(), define.CacheTypeAvItem, "json")
	b := c.Get(id)
	var avItem AvItem
	if b != nil {
		_ = json.Unmarshal(b, &avItem)
		log.Printf("id:%s,read cache\n", id)
		avItem.BuildLink(s.Host())
		return &avItem
	}
	pAvItem := s.GetOne(id)
	if pAvItem == nil {
		return nil
	}
	log.Printf("id:%s,save cache\n", id)
	_ = c.Set(id, utils.Encode(pAvItem))
	pAvItem.BuildLink(s.Host())
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

	posterFilePath := path.Join(dirRoot, "images", "poster_"+a.AvId+".jpg")
	utils.ImageDownload(posterFilePath, a.Poster)
	clearartFilePath := path.Join(dirRoot, "images", "clearart_"+a.AvId+".jpg")
	utils.ImageDownload(clearartFilePath, a.Clearart)

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
			Poster: posterFilePath,
			Fanart: clearartFilePath,
		},
		Fanart: []define.Fanart{
			{
				Thumb: clearartFilePath,
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

	posterFilePath = path.Join(path.Dir(avPath), patternFileRet[1]+".nfo")
	if utils.Exists(posterFilePath) {
		_ = os.Remove(posterFilePath)
	}
	err = os.WriteFile(posterFilePath, append([]byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>`), xmlByte...), 0644)
	if err != nil {
		log.Fatal("nfo文件写入失败,", err)
		return
	}
}
