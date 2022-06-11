package errormsg

const (
	SUCCESS = 200
	ERROR   = 500

	// code =1000...用户模块错误
	ERROR_USERNAME_USED    = 1001
	ERROR_PASSWORD_WRONG   = 1002
	ERROR_USER_NOT_EXIST   = 1003


	ERROP_DATABASE_SCAN_FALL = 5001
)

var codeMsg = map[int]string{
	SUCCESS:                "OK",
	ERROR:                  "FALL",
	ERROR_USERNAME_USED:    "用户名已存在",
	ERROR_PASSWORD_WRONG:   "密码错误",
	ERROR_USER_NOT_EXIST:   "用户不存在",

	ERROP_DATABASE_SCAN_FALL: "绑定数据库数据失败",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
