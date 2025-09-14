package main

import (
	"flag"
	"github.com/moozik/nfo-spider/spider/av/av_base"
	"github.com/moozik/nfo-spider/spider/av/javbus"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/moozik/nfo-spider/define"
	"github.com/moozik/nfo-spider/utils"
)

func main() {
	log.Println(utils.GetCurrentDirectory())
	log.Println(utils.IsRelease())

	cacheDir := path.Join(utils.GetCurrentDirectory(), string(os.PathSeparator), "cache")
	if !utils.Exists(cacheDir) {
		log.Println("mkdir cache")
		_ = os.Mkdir(cacheDir, 0755)
	}
	// filePtr, err := os.OpenFile(fmt.Sprintf("logs/log_%s.log", time.Now().Format("2006_01_02")), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(filePtr)

	log.Printf("args:%+v\n", os.Args)

	avIdFlag := flag.String("a", "", "avid")
	avDirFlag := flag.String("d", "", "影片读取目录")
	writeDirFlag := flag.String("r", "", "文件写入路径")
	nasDirFlag := flag.String("n", "", "文件内路径（nas内部路径）")
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()
	if *avIdFlag != "" {
		log.Println(utils.EncodeString(av_base.GetOneWithCache(javbus.NewAvJavbus().SetDebug(*debug), *avIdFlag)))
	} else if *avDirFlag != "" && *writeDirFlag != "" && *nasDirFlag != "" {
		walkDir(*nasDirFlag, *writeDirFlag, *avDirFlag, *avDirFlag)
	}
}

func walkDir(nasDir, writeDir, dirRoot, dirNow string) {
	err := filepath.Walk(dirNow, func(path string, info os.FileInfo, err error) error {
		log.Printf("path:%s", path)
		//if info.IsDir() {
		//	if dirNow == path {
		//		return nil
		//	}
		//	walkDir(writeDir, dirRoot, path)
		//	return nil
		//}
		if define.PattrenAvFile.MatchString(info.Name()) {
			avIdList := define.PattrenAvName.FindStringSubmatch(info.Name())
			if len(avIdList) == 0 {
				return nil
			}
			avData := av_base.GetOneWithCache(javbus.NewAvJavbus(), avIdList[0])
			av_base.XMLBuild(nasDir, writeDir, dirRoot, path, avData)
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
	}
}
