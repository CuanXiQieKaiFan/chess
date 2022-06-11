package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func (message *ClientManager) StartChat() {
	for true {
		fmt.Println("#######监听管道#######")
		select {
		case conn := <-Manager.Register:
			fmt.Printf("有新连接:%v", conn.ID)
			Manager.Clients[conn.ID] = conn //用户连接放到用户管理
			replyMsg := ReplyMsg{
				Content: "连接服务器成功",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

		case conn := <-Manager.UnRegister:
			fmt.Printf("连接失败%s", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Content: "连接中断",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				//断开连接
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
			//广播消息
		case broadcast := <-Manager.Broadcast:
			message := broadcast.Message
			sendId := broadcast.Client.SendID
			flag := false //默认 不在线
			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				}
			}
			_ = broadcast.Client.ID
			if flag {
				replyMsg := &ReplyMsg{
					Content: "在线中...",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg) //消息进行广播
				//插入数据库
			}else {
				replyMsg :=ReplyMsg{
					Content:"不在线!",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)

			}

		}
	}
}
