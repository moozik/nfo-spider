package javbus

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"regexp"
	"testing"
)

func Test1(t *testing.T) {
	parse, _ := url.Parse("https://www.busfan.cfd:8081/pics/cover/9xns_b.jpg")
	log.Println("Path", parse.Path)
	log.Println("Host", parse.Host)
	log.Println("Scheme", parse.Scheme)
	log.Println("RawQuery", parse.RawQuery)
	log.Println("RawPath", parse.RawPath)
	log.Println("Fragment", parse.Fragment)
}

func TestSplit(t *testing.T) {
	//取文件名
	fmt.Println(filepath.Base(`G:\jp_sister\abw-238-C.mp4`))
	//取路径
	fmt.Println(filepath.Dir(`G:\jp_sister\abw-238-C.mp4`))

	s := "/pics/cover/9ok5_b.jpg"
	fmt.Println(s[12:16])
}

func TestRegexp(t *testing.T) {
	c, _ := regexp.Compile(`^[a-zA-Z]{3,5}-\d{3,5}`)

	avList := []string{
		"abw-350ch.mp4",
		"ABW-331-C.mp4",
		"435MFCS-041.mp4",
		"FC2-PPV-3136593.mp4",
	}
	for _, avFile := range avList {
		fmt.Println(avFile, c.FindStringSubmatch(avFile))
	}
}
