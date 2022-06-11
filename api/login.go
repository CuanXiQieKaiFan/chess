package api

import (
	"chess/model"
	"chess/util/errormsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

//登录
func Login(c *gin.Context) {
	uname := c.PostForm("username")
	upwd := c.PostForm("password")
	var code int
	code = model.CheckLogin(uname, upwd)

	c.JSON(http.StatusOK, gin.H{
		"status":       code,
		"message":      errormsg.GetErrMsg(code),
	})
}
