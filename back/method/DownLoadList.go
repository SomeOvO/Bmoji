package method

import (
	"encoding/json"
	"io"
	"main/config"
	"main/loger"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
)

var list struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data struct {
		Packages []struct {
			ID      int    `json:"id"`
			Text    string `json:"text"`
			Type    int    `json:"type"`
			Url     string `json:"url"`
			UrlType int
		} `json:"all_packages"`
	} `json:"data"`
}

func DownloadEmoteList() bool {
	baseurl := "https://api.bilibili.com/x/emote/setting/panel?business=reply"
	client := http.Client{}
	req, err := http.NewRequest("GET", baseurl, nil)
	if err != nil {
		loger.Loger.Error("[Download]无法创建请求", zap.Error(err))
		return false
	}
	req.Header.Set("Cookie", "SESSDATA="+config.Cfg.BiliBili.SESSDATA)
	resp, err := client.Do(req)
	if err != nil {
		loger.Loger.Error("[Download]启动请求失败", zap.Error(err))
		return false
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Loger.Error("[Download]读取响应体失败", zap.Error(err))
		return true
	}
	json.Unmarshal(data, &list)
	for i, a := range list.Data.Packages {
		surl, urltype := geturl(a.Url)
		list.Data.Packages[i].Url = surl
		list.Data.Packages[i].UrlType = urltype
	}
	file, err := json.Marshal(list)
	if err != nil {
		loger.Loger.Error("[Download]序列化Json失败", zap.Error(err))
		return false
	}
	l := len(list.Data.Packages)
	loger.Loger.Info("[Download]下载完毕", zap.Int("数量", l))
	if l == 0 {
		loger.Loger.Error("[Download]下载后数量为0，请检查登陆凭证")
		return false
	}
	os.WriteFile("resp.json", file, 0775)

	return true
}
func geturl(url string) (string, int) {

	//神人BiliBili,有四个域名作为图片域名，甚至还有一个是自定义域名
	urls := []string{"https://i0.hdslb.com/bfs/emote/", "https://i0.hdslb.com/bfs/emote/", "https://i0.hdslb.com/bfs/garb/", "https://i0.hdslb.com/bfs/garb/"}
	for i, v := range urls {
		arr := strings.Split(url, v)
		if len(arr) > 1 {
			return arr[len(arr)-1], i
		}
	}
	return url, -1
}
