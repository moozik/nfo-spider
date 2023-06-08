package javbus

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/moozik/nfo-spider/define"
	"github.com/moozik/nfo-spider/utils"
)

const SITE_NAME = "javbus"
const JAVBUS_HOST = "https://www.cdnbus.lol"

const (
	ReleaseDateWord = "發行日期:"
	LenWord         = "長度:"
	DirectorWord    = "導演:"
	StudioWord      = "製作商:"
	LabelWord       = "發行商:"
	SeriesWord      = "系列:"
)

func GetHost() string {
	return JAVBUS_HOST
}

func GetOneWithCache(avId string) *define.AvItem {
	c := utils.NewCache(SITE_NAME, "nfo")
	b := c.Get(avId)
	var avItem define.AvItem
	if b != nil {
		_ = json.Unmarshal(b, &avItem)
		log.Printf("av_id:%s,read cache\n", avId)
		return &avItem
	}
	pAvItem := GetOne(avId)
	b, _ = json.Marshal(pAvItem)
	log.Printf("av_id:%s,save cache\n", avId)
	c.Set(avId, b)
	return pAvItem
}

func GetOne(avId string) *define.AvItem {
	avId = strings.ToUpper(avId)
	link := GetHost() + "/" + avId
	doc, err := docGet(link)
	if err != nil {
		log.Printf("docGet fail,err:" + err.Error())
		return nil
	}

	ret := define.AvItem{
		Link:  link,
		AvId:  avId,
		Title: doc.Find("body > div.container > h3").First().Text(),
		Genre: getGenre(doc),
		Stars: getStars(doc),
	}

	bigImage := doc.Find("a.bigImage > img").First().AttrOr("src", "")
	ret.Clearart = url(bigImage)
	if strings.Contains(bigImage, "/digital/video") {
		ret.Poster = url(strings.Replace(bigImage, "pl.jpg", "ps.jpg", 1))
	} else if strings.Contains(bigImage, "/pics/cover") {
		ret.Poster = url(fmt.Sprintf("/pics/thumb/%s.jpg", bigImage[12:16]))
	}

	doc.Find("span.header").Each(func(i int, s *goquery.Selection) {
		if s.Text() == ReleaseDateWord {
			ret.ReleaseDate = parseData(s.Parent().Text(), ReleaseDateWord)
		}
		if s.Text() == LenWord {
			ret.Len = parseData(s.Parent().Text(), LenWord)
		}
		if s.Text() == DirectorWord {
			ret.Director = parseData(s.Parent().Text(), DirectorWord)
		}
		if s.Text() == StudioWord {
			ret.Studio = parseData(s.Parent().Text(), StudioWord)
		}
		if s.Text() == LabelWord {
			ret.Label = parseData(s.Parent().Text(), LabelWord)
		}
		if s.Text() == SeriesWord {
			ret.Series = parseData(s.Parent().Text(), SeriesWord)
		}
	})
	return &ret
}

func parseData(target, keyword string) string {
	return strings.TrimSpace(strings.Replace(target, keyword, "", 1))
}

func getGenre(doc *goquery.Document) []string {
	var ret []string
	doc.Find("span.genre > label > a").Each(func(i int, s *goquery.Selection) {
		ret = append(ret, s.First().Text())
	})
	return ret
}

func getStars(doc *goquery.Document) []define.Stars {
	var ret []define.Stars
	doc.Find("div.star-box > li > a > img").Each(func(i int, s *goquery.Selection) {
		ret = append(ret, define.Stars{
			Name:  s.First().AttrOr("title", ""),
			Image: url(s.First().AttrOr("src", "")),
		})
	})
	return ret
}

func url(s string) string {
	return utils.Url(GetHost(), s)
}
