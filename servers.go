package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wuyan94zl/IM/controllers"
	"github.com/wuyan94zl/IM/websocket"
)

func main(){
	hub := websocket.NewHub()
	go hub.Run()

	r := gin.Default()

	r.Static("css", "./html/css")
	r.LoadHTMLGlob("./html/**/*")

	//聊天室列表
	r.GET("/rooms", controllers.RoomList)

	//创建聊天室
	r.POST("/rooms/add", controllers.RoomAdd)

	//进入聊天室
	r.GET("/room/:number/:name", controllers.RoomInfo)

	//聊天室服务
	r.GET("/ws/:room", func(c *gin.Context) {
		websocket.RunWs(hub,c)
	})
	r.Run(":8303")

}
