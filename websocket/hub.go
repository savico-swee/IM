package websocket

import "fmt"

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
			h.clients[client] = true
			fmt.Println("注册")
			fmt.Println(client)
		case client := <- h.unregister://取消注册
			if _, ok := h.clients[client]; ok {
				fmt.Println("取消注册")
				fmt.Println(client)
				delete(h.clients, client)//删除客户端
				close(client.send)//关闭管道
			}
		case message := <-h.broadcast://接收消息并转发消息
			room := <-h.room
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
