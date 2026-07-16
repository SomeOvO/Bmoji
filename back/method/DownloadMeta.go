package method

import (
	"encoding/json"
	"fmt"
	"io"
	"main/config"
	"main/loger"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func DownloadEmoteMeta() bool {
	var download []struct {
		ID     int    `json:"id"`
		Text   string `json:"text"`
		Emotes []struct {
			Text   string `json:"text"`
			Url    string `json:"url"`
			GifUrl string `json:"gif_url"`
			Type   int    `json:"type"`
			Meta   struct {
				Alias string `json:"alias"`
			} `json:"meta"`
		} `json:"emote"`
	}
	wd, err := os.Getwd()
	if err != nil {
		loger.Loger.Error("[Meta]获取工作目录失败", zap.Error(err))
		return false
	}
	file, err := os.ReadFile(wd + "/resp.json")
	if err != nil {
		loger.Loger.Error("[Meta]读取列表文件失败", zap.Error(err))
		return false
	}
	err = json.Unmarshal(file, &list)
	if err != nil {
		loger.Loger.Error("[Meta]反序列化失败", zap.Error(err))
		return false
	}
	//获取所有id字符串
	var idlist []string
	for _, v := range list.Data.Packages {
		idlist = append(idlist, strconv.Itoa(v.ID))
	}

	//50个为一组分组
	var ids [][]string
	lent := int(math.Floor(float64(len(idlist) / 50)))
	for i := range lent {
		s := i * 50
		o := (i + 1) * 50
		if i+1 == lent {
			ids = append(ids, idlist[s:])
		} else {
			ids = append(ids, idlist[s:o])
		}

	}
	//获取
	for i, v := range ids {
		var metalist struct {
			Code int    `json:"code"`
			Msg  string `json:"message"`
			Data struct {
				Package []struct {
					ID     int    `json:"id"`
					Text   string `json:"text"`
					Emotes []struct {
						Text   string `json:"text"`
						Url    string `json:"url"`
						GifUrl string `json:"gif_url"`
						Type   int    `json:"type"`
						Meta   struct {
							Alias string `json:"alias"`
						} `json:"meta"`
					} `json:"emote"`
				} `json:"packages"`
			} `json:"data"`
		}

		client := http.Client{}
		baseurl := "https://api.bilibili.com/x/emote/package"
		var str strings.Builder
		for _, val := range v {
			str.WriteString(val)
			str.WriteString(",")
		}
		resq, err := http.NewRequest("GET", baseurl+"?business=reply&ids="+str.String(), nil)
		resq.Header.Set("Cookie", "SESSDATA="+config.Cfg.BiliBili.SESSDATA)
		resp, err := client.Do(resq)
		if err != nil {
			loger.Loger.Error("[Meta]无法发起请求", zap.Error(err))
			return false
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			loger.Loger.Error("[Meta]读取响应体失败", zap.Error(err))
			return false
		}
		err = json.Unmarshal(data, &metalist)

		if err != nil {
			loger.Loger.Error("[Meta]反序列化失败", zap.Error(err))
			return false
		}

		//替换http为https，否则前端下载报错
		for _, val := range metalist.Data.Package {
			for i, vals := range val.Emotes {
				val.Emotes[i].Url = strings.Replace(vals.Url, "http://", "https://", 1)
				val.Emotes[i].GifUrl = strings.Replace(vals.GifUrl, "http://", "https://", 1)
			}
			download = append(download, val)
		}
		fmt.Printf("\r%v/%v", i, len(ids))
	}
	fmt.Println()
	file, err = json.Marshal(download)
	if err != nil {
		loger.Loger.Error("[Meta]无法序列化", zap.Error(err))
		return false
	}
	//写入
	err = os.WriteFile("download.json", file, 0775)
	if err != nil {
		loger.Loger.Error("[Meta]无法写入文件", zap.Error(err))
		return false
	}

	return true

}
