package tools

import (
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

type Cap struct {
	Code    string `json:"code" form:"code"`
	Captcha string `json:"captcha" form:"captcha"`
}

// VerificationCode 结构体表示一个验证码
type VerificationCode struct {
	Code string
}

// GenerateVerificationCode 生成一个6位随机验证码
func GenerateVerificationCode() *VerificationCode {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return &VerificationCode{Code: code}
}
func SendVerificationCode(youxiang string, code string) error {
	e := email.NewEmail()
	e.From = "测试专用 <18963508449@163.com>"
	e.To = []string{youxiang}
	e.Subject = "验证码"
	e.Text = []byte("您的验证码是：" + code)
	err := e.Send("smtp.163.com:25", smtp.PlainAuth("", youxiang, "SEDPPHPSTFGRHBBK", "smtp.163.com"))
	if err != nil {
		return err
	}
	return nil
}

func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return strconv.Itoa(code)
}
