package api

import (
	"encoding/json"
	"fmt"
	"main/config"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetEmoteInfo(c *gin.Context) {

	switch config.Cfg.App.SaveType {
	case "local":
		findlocal(c)
	}
}

var metalist []struct {
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

var MetalistByte []byte

func findlocal(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(403, Send(403, "id不正确"))
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if MetalistByte == nil {
		fmt.Println("nil")
		file, err := os.ReadFile(wd + "/download.json")
		if err != nil {
			panic(err)
		}
		MetalistByte = file
		json.Unmarshal(MetalistByte, &metalist)
	}

	for _, v := range metalist {
		if v.ID == id {
			c.JSON(200, v)
		}
	}
}
