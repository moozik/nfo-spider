package av_base

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
