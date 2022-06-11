package model

import (
	"chess/util/errormsg"
	"fmt"
)

type User struct {
	ID int `json:"id",omitempty`
	UserName string `json:"username" validate:"required,min=4,max=12"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=20"`
}

//注册用户
func Register(data *User) int {
	//data.Password = ScryptPw(data.Password)
	dst := "insert into users (users.username,users.password) values (?,?)"
	_, err = Db.Exec(dst, data.UserName, data.Password)
	if err != nil {
		return errormsg.ERROR
	}
	return errormsg.SUCCESS
}


//登录验证
func CheckLogin(username string, password string) int {
	var user User
	row := Db.QueryRow("select username,password from users where username = ?", username)
	err := row.Scan(&user.UserName, &user.Password)
	if err != nil {
		fmt.Printf("err is:%v\n", err)
	}
	if user.UserName == "" {
		return errormsg.ERROR_USER_NOT_EXIST
	}
	if password != user.Password {
		return errormsg.ERROR_PASSWORD_WRONG
	}

	return errormsg.SUCCESS
}

