package routers

import (
	"chess/api"
	"chess/service"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	r.Use(service.Cors())
	user := r.Group("/user") //用户相关操作
	{
		user.POST("/register", api.Register)    //注册用户
		user.GET("/login", api.Login)           //登录用户
	}

	chat:=r.Group("/chat")  //聊天相关操作
	{
		chat.GET("ws", service.Handle)
	}

	_ = r.Run(":8000")
}
