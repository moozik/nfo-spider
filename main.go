package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/moozik/nfo-spider/define"
	"github.com/moozik/nfo-spider/javbus"
	"github.com/moozik/nfo-spider/utils"
)

func main() {
	log.Println(utils.GetCurrentDirectory())
	log.Println(utils.IsRelease())
	// filePtr, err := os.OpenFile(fmt.Sprintf("logs/log_%s.log", time.Now().Format("2006_01_02")), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(filePtr)

	log.Printf("args:%+v\n", os.Args)

	avId := flag.String("a", "", "avid")
	avDir := flag.String("d", "", "影片目录")
	flag.Parse()

	if *avId != "" {
		log.Println(utils.Encode(javbus.GetOneWithCache(*avId)))
	} else if *avDir != "" {
		walkDir(*avDir, *avDir)
	}
}

func walkDir(dirRoot, dirNow string) {
	err := filepath.Walk(dirNow, func(path string, info os.FileInfo, err error) error {
		log.Printf("path:%s", path)
		if info.IsDir() {
			if dirNow == path {
				return nil
			}
			walkDir(dirRoot, path)
			return nil
		}
		if define.PattrenAvFile.MatchString(info.Name()) {
			avIdList := define.PattrenAvName.FindStringSubmatch(info.Name())
			if len(avIdList) == 0 {
				return nil
			}
			avData := javbus.GetOneWithCache(avIdList[0])
			utils.XMLBuild(dirRoot, path, avData)
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
	}
}
