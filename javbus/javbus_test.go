package javbus

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/moozik/nfo-spider/utils"
)

func TestGetOne(t *testing.T) {
	var ret any
	//多女主 dmm主图
	// ret = GetOne("KYMI-031")
	// t.Log(utils.Encode(ret))

	//站内图
	ret = GetOne("JUQ-293")
	t.Log(utils.Encode(ret))
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
