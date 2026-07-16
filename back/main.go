package main

import (
	"main/config"
	"main/loger"
	"main/method"
	"main/router"
	"os"
)

func main() {
	config.GetConfig()
	loger.InitLog()
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "meta":
			if !method.DownloadEmoteMeta() {
				loger.Loger.Fatal("[Meta]下载失败")
			}
		case "list":
			if !method.DownloadEmoteList() {
				loger.Loger.Fatal("[List]下载失败")
			}
		}
		return
	}
	go method.AutoUpdate()
	router.StartUp()
}
