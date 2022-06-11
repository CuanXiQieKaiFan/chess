package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"
)

//发送消息
type SendMsg struct {
	Content string `json:"content"`
}

//回复消息
type ReplyMsg struct {
	Content string `json:"content"`  //即面板传递的消息
}

//用户端结构体
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

//广播类
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

//用户管理类
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	UnRegister chan *Client
}

//消息类
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

//实例管理者
var Manager = ClientManager{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	UnRegister: make(chan *Client),
}

//跨域处理，似乎不弄部署不行
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		cors.New(cors.Config{
			AllowAllOrigins: true, //允许所有跨域
			AllowMethods:    []string{"*"},
			AllowHeaders:    []string{"origin"},
			ExposeHeaders:   []string{"Content-Length", "Authorization"},
			MaxAge: 12 * time.Hour,
		})
	}
}

func Handle(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")

	//升级websocket协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { //处理跨域
			return true
		}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	//用户实例
	client := &Client{
		ID:    uid, //发送方 1->2
		SendID: toUid, //接受发送方 2-1
		Socket: conn,
		Send:   make(chan []byte),
	}
	//用户注册到用户管理
	Manager.Register <- client
	//开始读写
	go client.Read()
	go client.Write()
}


//读消息
func (c *Client) Read() {
	defer func() {
		Manager.UnRegister <- c //离线
		_ = c.Socket.Close()    //关闭socket
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		err := c.Socket.ReadJSON(&sendMsg)
		if err != nil {
			fmt.Println("数据格式不对", err)
			Manager.UnRegister <- c
			_ = c.Socket.Close() //数据不对断开连接
			break
		}

		fmt.Println(c.ID,"发送消息",sendMsg.Content)
		Manager.Broadcast <- &Broadcast{
			Client:  c,
			Message: []byte(sendMsg.Content),
		}
	}
}

//写消息
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for true {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := ReplyMsg{
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg) //序列化消息
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

