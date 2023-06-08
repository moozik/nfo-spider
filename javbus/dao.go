package javbus

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var once sync.Once
var client *http.Client

func init() {
	once.Do(func() {
		client = &http.Client{
			Timeout: time.Second * 5,
		}
	})
}

func docGet(url string) (*goquery.Document, error) {
	log.Printf("httpGet,url:%s\n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
