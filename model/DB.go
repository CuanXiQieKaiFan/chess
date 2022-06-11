package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB
var err error

//初始化连接数据库
func InitDb() (err error) {
	dsn := "root:254092@tcp(127.0.0.1:3306)/chess?charset=utf8mb4&parseTime=True"

	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("mysql 连接芭比Q啦")
		panic(err)
	}

	err = Db.Ping()
	if err != nil {
		return err
	}
	return nil
}
