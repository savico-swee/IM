package websocket

import (
	"fmt"
	"github.com/wuyan94zl/IM/database"
	"github.com/jinzhu/gorm"
	"github.com/wuyan94zl/IM/models"
	"time"
)

type Hub struct {
	// 已注册客户端
	clients 	map[*Client]bool
	// 客户端发送的消息
	broadcast chan []byte
	room chan string
	// 注册
	register chan *Client
	// 取消注册
	unregister chan *Client
}

func NewHub() *Hub  {
	return &Hub{
		broadcast:  make(chan []byte),
		room: make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run(){
	for {
		select {
		case client := <- h.register:// 注册
			database.DB.Table("rooms").Where("number = ?",client.room).UpdateColumn("user_num",gorm.Expr("user_num + ?", 1))
			roomHasUser := models.RoomHasUser{RoomNumber: client.room,UserName: client.name,CreatedAt: time.Now(),UserId: client.id}
			database.DB.Create(&roomHasUser)
			h.clients[client] = true
			fmt.Println("注册")
		case client := <- h.unregister://取消注册
			if _, ok := h.clients[client]; ok {
				database.DB.Table("rooms").Where("number = ?",client.room).UpdateColumn("user_num",gorm.Expr("user_num - ?", 1))
				database.DB.Delete(models.RoomHasUser{},"room_number = ? AND user_name = ?",client.room,client.name)
				fmt.Println("取消注册")
				delete(h.clients, client)//删除客户端
				close(client.send)//关闭管道
			}
		case message := <-h.broadcast://接收消息并转发消息

			room := <-h.room
			fmt.Println(string(message),string(room))
			for client := range h.clients {
				if room == client.room{
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
