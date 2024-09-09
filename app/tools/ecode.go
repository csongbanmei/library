package tools

import "fmt"

var (
	OK       = Ecode{Code: 0}
	NotLogin = Ecode{Code: 1001, Message: "用户未登录"}
	ParamErr = Ecode{Code: 1002, Message: "参数错误"}
	UserErr  = Ecode{Code: 1003, Message: "账号密码错误"}
)

type Ecode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	PayUrl  any    `json:"payUrl"`
}

func (e *Ecode) String() string {
	return fmt.Sprintf("ecode:%d,message:%s", e.Code, e.Message)
}
