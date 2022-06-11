package api

import (
	"chess/model"
	"github.com/gin-gonic/gin"
)

//注册用户
func Register(c *gin.Context) {
	var data model.User
	username:=c.PostForm("username")
	password:=c.PostForm("password")
	data.Password=password
	data.UserName=username

	model.Register(&data)
}

