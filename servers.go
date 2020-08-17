package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wuyan94zl/IM/websocket"
)

func main(){
	hub := websocket.NewHub()
	go hub.Run()
	bindAddress := "localhost:2303"
	r := gin.Default()
	//r.GET("/ping", websocket.Ping)
	r.GET("/ws/:room", func(c *gin.Context) {
		websocket.RunWs(hub,c)
	})
	r.Run(bindAddress)

}
