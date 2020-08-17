package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
type Client struct{
	hub 	*Hub
	// websocket 连接
	conn	*websocket.Conn
	// 发送管道
	send	chan []byte
	// 客户端名称
	name 	string
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
			fmt.Println(string(message))
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
		if string(message) == "" {// 接收到空，则关闭连接
			c.hub.unregister <- c
			return
		}else{
			// 封装消息 map
			msg := map[string]string{"name":c.name,"message":string(message),"time":time.Now().String()}
			// map 转 bytes 发送 json 数据
			bt,err := json.Marshal(msg)
			if err != nil {
				break
			}
			c.hub.broadcast <- bt
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
	// 实例化新用户
	client := &Client{hub:hub,conn: conn,send: make(chan []byte),room: room,name: name}
	// 新用户注册 想管道发送注册信息
	client.hub.register <- client
	fmt.Println(client)
	fmt.Println(client.hub)
	go client.readMsg()
	go client.writeMsg()
}

//webSocket请求ping 返回pong
func Ping(c *gin.Context) {
	fmt.Println("sssss")
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		if string(message) == "ping" {
			message = []byte("pong")
		}
		fmt.Println(mt)
		fmt.Println(string(message))
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

