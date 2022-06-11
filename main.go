package main

import (
	"chess/model"
	"chess/routers"
	"chess/service"
	"fmt"
)

func main() {
	err := model.InitDb()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("连接数据库成功！")
	}
	go service.Manager.StartChat()
	routers.Start()

}