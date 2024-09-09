package test

import (
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"library/app/model"
	"library/app/tools"
	"os"
	"testing"
)

func TestBorrowBook(t *testing.T) {
	model.NewMysql()
	err := model.BorrowBook(1, 1)
	fmt.Println(err)
}
func TestSendCode(t *testing.T) {
	// 替换为您的阿里云 AccessKey
	accessKeyId := "YourAccessKeyId"
	accessKeySecret := "YourAccessKeySecret"
	// 替换为您的短信签名和模板代码
	signature := "YourSMSignature"
	templateCode := "YourSMSTemplateCode"
	// 获取手机号码
	phoneNumber := "18963508449"
	code := tools.GenerateVerificationCode().Code
	model.SendSMS(accessKeyId, accessKeySecret, signature, templateCode, phoneNumber, code)
}
func TestSendCodeV1(t *testing.T) {
	fmt.Println(os.Args[1:])
	err := model.SendCodeV1(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
func TestWechat(t *testing.T) {

}
