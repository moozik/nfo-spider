package javbus

import (
	"fmt"
	"github.com/moozik/nfo-spider/spider/av/av_base"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const SITE_NAME = "javbus"
const JAVBUS_HOST = "www.busfan.cfd"

const (
	ReleaseDateWord = "發行日期:"
	LenWord         = "長度:"
	DirectorWord    = "導演:"
	StudioWord      = "製作商:"
	LabelWord       = "發行商:"
	SeriesWord      = "系列:"
)

type AvJavbus struct {
	av_base.AvBase
}

func NewAvJavbus() *AvJavbus {
	return &AvJavbus{}
}
func (a *AvJavbus) AppName() string {
	return SITE_NAME
}
func (a *AvJavbus) Host() string {
	return JAVBUS_HOST
}
func (a *AvJavbus) Url(s string) string {
	if s == "" {
		return ""
	}
	if s[:4] == "http" {
		return s
	}
	if s[0] == '/' {
		return fmt.Sprintf("https://%s%s", a.Host(), s)
	}
	return fmt.Sprintf("https://%s/%s", a.Host(), s)
}

func (a *AvJavbus) GetOne(avId string) *av_base.AvItem {
	avId = strings.ToUpper(avId)
	doc, err := a.DocGet(a.Url(avId))
	if err != nil {
		log.Printf("docGet fail,err:" + err.Error())
		return nil
	}

	ret := av_base.AvItem{
		Scheme: "https",
		Link:   "/" + avId,
		AvId:   avId,
		Title:  doc.Find("body > div.container > h3").First().Text(),
		Genre:  getGenre(doc),
		Stars:  a.getStars(doc),
	}

	bigImage := doc.Find("a.bigImage > img").First().AttrOr("src", "")

	if strings.Contains(bigImage, "/digital/video") {
		//完整路径
		ret.Poster = a.Url(strings.Replace(bigImage, "pl.jpg", "ps.jpg", 1))
		ret.Clearart = bigImage
	} else if strings.Contains(bigImage, "/pics/cover") {
		//部分路径
		ret.Poster = fmt.Sprintf("/pics/thumb/%s.jpg", bigImage[12:16])
		ret.Clearart = fmt.Sprintf("/pics/cover/%s_b.jpg", bigImage[12:16])
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

func (a *AvJavbus) getStars(doc *goquery.Document) []av_base.Stars {
	var ret []av_base.Stars
	doc.Find("div.star-box > li > a > img").Each(func(i int, s *goquery.Selection) {
		ret = append(ret, av_base.Stars{
			Name:  s.First().AttrOr("title", ""),
			Image: s.First().AttrOr("src", ""),
		})
	})
	return ret
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
