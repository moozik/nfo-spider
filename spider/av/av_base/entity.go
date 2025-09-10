package av_base

import (
	"fmt"
)

type AvItem struct {
	Scheme string `json:"scheme"`
	Link   string `json:"link"`

	AvId        string   `json:"av_id"`
	Title       string   `json:"title"`
	ReleaseDate string   `json:"release_date"`
	Len         string   `json:"len"`
	Director    string   `json:"director"`
	Studio      string   `json:"studio"`
	Label       string   `json:"label"`
	Series      string   `json:"series"`
	Genre       []string `json:"genre"`
	Stars       []Stars  `json:"stars"`
	Poster      string   `json:"poster"`   //竖着海报
	ClearArt    string   `json:"clearart"` //横着海报
}

type Stars struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (a *AvItem) BuildLink(host string) {
	if a.Poster != "" && a.Poster[0] == '/' {
		a.Poster = a.Url(host, a.Poster)
	}
	if a.ClearArt != "" && a.ClearArt[0] == '/' {
		a.ClearArt = a.Url(host, a.ClearArt)
	}
	for k, item := range a.Stars {
		a.Stars[k].Image = a.Url(host, item.Image)
	}
}

func (a *AvItem) Url(host, path string) string {
	if path == "" {
		return ""
	}
	if path[:4] == "http" {
		return path
	}
	if path[0] == '/' {
		return fmt.Sprintf("%s://%s%s", a.Scheme, host, path)
	}
	return fmt.Sprintf("%s://%s/%s", a.Scheme, host, path)
}
