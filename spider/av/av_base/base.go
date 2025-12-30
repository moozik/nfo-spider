package av_base

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/moozik/nfo-spider/utils/ai"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	IsDebug bool //打印debug
}

func (a *AvBase) DocGet(url string) (*goquery.Document, error) {
	log.Printf("httpGet,url:%s\n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	//设置cookie：PHPSESSID=7tmo4naqrpb01ah7nopts6dd80; existmag=mag
	req.Header.Set("Cookie", "PHPSESSID=7tmo4naqrpb01ah7nopts6dd80; existmag=mag")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "zh-CN,zh-TW;q=0.9,zh-HK;q=0.8,zh;q=0.7")
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("http status:%v", resp.Status)
		return nil, errors.New("status:" + resp.Status)
	}
	return goquery.NewDocumentFromReader(resp.Body)
}

func GetOneWithCache(s avInter, id string) *AvItem {
	id = strings.ToUpper(id)
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

func XMLBuild(nasDir, writeDir, realDir, avPath string, a *AvItem) {
	log.Printf("nasDir:%s,writeDir:%s,realDir:%s,avPath:%s", nasDir, writeDir, realDir, avPath)
	if a == nil {
		log.Println("nil input *define.AvItem")
		return
	}
	avFileName := path.Base(avPath)
	patternFileRet := define.PattrenAvFile.FindStringSubmatch(avFileName)
	if len(patternFileRet) != 3 {
		log.Printf("文件名错误:%s", avFileName)
		return
	}

	posterFilePath := path.Join(writeDir, "images", "poster_"+a.AvId+".jpg")
	utils.ImageDownload(posterFilePath, a.Poster)
	clearartFilePath := path.Join(writeDir, "images", "clearart_"+a.AvId+".jpg")
	utils.ImageDownload(clearartFilePath, a.ClearArt)

	nasDirPosterFilePath := path.Join(nasDir, "images", "poster_"+a.AvId+".jpg")
	nasDirClearartFilePath := path.Join(nasDir, "images", "clearart_"+a.AvId+".jpg")
	titleTransLate := ai.Translate(strings.Replace(a.Title, a.AvId, "", 1))
	ret := define.NfoMovie{
		Title:         a.AvId + " " + titleTransLate,
		Polt:          titleTransLate,
		Originaltitle: a.Title,
		Sorttitle:     a.AvId,
		Premiered:     a.ReleaseDate,
		Releasedate:   a.ReleaseDate,
		Writer:        a.Series,
		Credits:       a.Series,
		Genre:         a.Genre,
		Art: define.Art{
			Poster: nasDirPosterFilePath,
			Fanart: nasDirClearartFilePath,
		},
		Fanart: []define.Fanart{
			{
				Thumb: nasDirClearartFilePath,
			},
		},
		Actor: []define.Actor{},
	}

	for _, act := range a.Stars {
		//imagePath := path.Join(dirRoot, "images", "actor_"+act.Name+".jpg")
		//utils.ImageDownload(imagePath, act.Image)
		ret.Actor = append(ret.Actor, define.Actor{
			Name: act.Name,
			Type: "actor",
			//Thumb: imagePath,
		})
	}

	xmlByte, err := xml.Marshal(ret)
	if err != nil {
		log.Fatalf("xml.Marshal fail,err:%v", err)
		return
	}

	base := filepath.Base(avFileName)
	avFileNameHead := strings.TrimSuffix(base, filepath.Ext(base))
	posterFilePath = path.Join(writeDir, avFileNameHead+".nfo")
	if !utils.Exists(posterFilePath) {
		err = os.WriteFile(posterFilePath, append([]byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>`), xmlByte...), 0644)
		if err != nil {
			log.Fatal("nfo文件写入失败,", err)
			return
		}
	} else {
		log.Println("路径已存在", posterFilePath)
	}
}
