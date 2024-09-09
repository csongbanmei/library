package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library/app/model"
	"library/app/tools"
	_ "net/http"
	_ "time"
)

func SendCode(c *gin.Context) {
	// 替换为您的阿里云 AccessKey
	accessKeyId := "你自己的id"
	accessKeySecret := "你自己的密码"
	// 替换为您的短信签名和模板代码
	signature := "阿里云短信测试"
	templateCode := "SMS_154950909"
	// 获取手机号码
	phoneNumber := c.PostForm("phone")
	fmt.Println(phoneNumber)
	//随机生成一个六位数的验证码
	code := tools.GenerateVerificationCode().Code
	//fmt.Println(code)
	// 发送短信
	err := model.SendSMS(accessKeyId, accessKeySecret, signature, templateCode, phoneNumber, code)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to send SMS",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "SMS sent successfully!",
	})
}
