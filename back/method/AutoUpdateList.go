package method

import (
	"main/loger"
	"time"
)

var LastUpdateTime *time.Time

func AutoUpdate() {
	for {
		if DownloadEmoteList() {
			if DownloadEmoteMeta() {
				loger.Loger.Info("[Auto]已自动更新")
			} else {
				loger.Loger.Error("[Auto]自动更新出错")
			}
		} else {
			loger.Loger.Error("[Auto]自动更新出错")
		}
		time.Sleep(time.Hour * 24)
	}
}
