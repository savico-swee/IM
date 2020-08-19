package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wuyan94zl/IM/database"
	"github.com/wuyan94zl/IM/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

//type SendMessage struct {
//	Name 		string	`json:name`
//	Message 	string	`json:message`
//	Time 		string	`json:time`
//	Type 		string  `json:type`
//}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var (
	newline = []byte{'\n'}
	space   = []byte{'<','b','r','>'}
)
type Client struct{
	hub 	*Hub
	// websocket 连接
	conn	*websocket.Conn
	// 发送管道
	send	chan []byte
	// 客户端名称
	name 	string
	id		int
	// 房间
	room	string
}
//客户端写消息
func (c *Client) writeMsg(){
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			//客户端关闭，退出
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
//客户端读消息
func (c *Client) readMsg(){
	defer func() {
		c.hub.unregister <- c
		// 推送 退出 消息
		msg := map[string]string{"name":c.name,"message":"退出聊天室","time":time.Now().Format("2006-01-02 15:04:05"),"type": "3","id":strconv.Itoa(c.id)}
		//bt := []byte("return im")
		// map 转 bytes 发送 json 数据
		bt,_ := json.Marshal(msg)
		c.hub.broadcast <- bt
		c.hub.room <- c.room
		models.RightLog(c.name,"退出聊天室",c.room,3)
		c.conn.Close()
	}()
	for  {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if string(message) == "" {// 接收到空，不处理
			//c.hub.unregister <- c
			//return
		}else{
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			msgtype := "1"
			//封装消息 map
			models.RightLog(c.name,string(message),c.room,1)
			msg := map[string]string{"name":c.name,"message":string(message),"time":time.Now().Format("2006-01-02 15:04:05"),"type": msgtype}
			// map 转 bytes 发送 json 数据
			message,_ := json.Marshal(msg)
			c.hub.broadcast <- message
			c.hub.room <- c.room
		}
	}
}

func RunWs(hub *Hub,c *gin.Context) {
	room := c.Param("room")
	name := c.DefaultQuery("name","wow")
	//升级get请求为webSocket协议
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	user := models.User{}
	database.DB.Where("name = ?",name).First(&user)
	// 实例化新用户
	client := &Client{hub:hub,conn: conn,send: make(chan []byte),room: room,name: name,id: user.Id}
	// 连接时休眠1秒  防止刷新页面 先连接后退出
	time.Sleep(time.Duration(1)*time.Second)

	// map 转 bytes 发送 json 数据
	msg := map[string]string{"name":client.name,"message":"进入聊天室","time":time.Now().Format("2006-01-02 15:04:05"),"type": "2","id":strconv.Itoa(user.Id)}
	bt,_ := json.Marshal(msg)
	client.hub.broadcast <- bt
	client.hub.room <- room
	models.RightLog(client.name,"进入聊天室",room,2)
	// 新用户注册 想管道发送注册信息
	client.hub.register <- client


	go client.readMsg()
	go client.writeMsg()
}

