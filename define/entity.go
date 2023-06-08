package define

import "encoding/xml"

type AvItem struct {
	Link string `json:"link"`

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
	Clearart    string   `json:"clearart"` //横着海报
}

type Stars struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Actor struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name"`
	Type  string `xml:"type"`
	Thumb string `xml:"thumb"`
}
type Art struct {
	Text   string `xml:",chardata"`
	Poster string `xml:"poster"`
	Fanart string `xml:"fanart"`
}
type Fanart struct {
	Text  string `xml:",chardata"`
	Thumb string `xml:"thumb"`
}
type NfoMovie struct {
	// Lockdata      string   `xml:"lockdata"`
	XMLName       xml.Name `xml:"movie"`
	Text          string   `xml:",chardata"`
	Title         string   `xml:"title"`
	Originaltitle string   `xml:"originaltitle"`
	Year          string   `xml:"year"`
	Sorttitle     string   `xml:"sorttitle"`
	Premiered     string   `xml:"premiered"`
	Releasedate   string   `xml:"releasedate"`
	Rating        string   `xml:"rating"`
	Criticrating  string   `xml:"criticrating"`
	Writer        string   `xml:"writer"`
	Credits       string   `xml:"credits"`
	Genre         []string `xml:"genre"`
	Tag           []string `xml:"tag"`
	Director      []string `xml:"director"`
	Art           Art      `xml:"art"`
	Fanart        []Fanart `xml:"fanart"`
	Actor         []Actor  `xml:"actor"`
}
