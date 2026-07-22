package api

import (
	"fmt"
	"main/config"
	"os"

	"github.com/gin-gonic/gin"
)

func GetRemoteList(c *gin.Context) {
	switch config.Cfg.App.SaveType {
	case "local":
		local(c)
	}
}

var List []byte

func local(c *gin.Context) {
	wd, err := os.Getwd()
	if err != nil {
		c.JSON(502, Send(502, "获取目录失败"))
		return
	}
	if List != nil {
		c.Data(200, "application/json", List)
		fmt.Println("内存已有")
		return
	}
	file, err := os.ReadFile(wd + "/resp.json")
	List = file
	if err != nil {
		c.JSON(502, Send(502, "获取目录失败"))
		return
	}
	c.Data(200, "application/json", file)
}
