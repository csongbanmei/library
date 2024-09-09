package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library/app/model"
	"library/app/tools"
	"net/http"
	"regexp"
	"time"
)

type User struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func GetLogin(c *gin.Context) {
	c.HTML(200, "login.tmpl", nil)
}
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: err.Error(), //这里有风险
		})
	}

	fmt.Printf("user:%+v\n", user)

	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: user.CaptchaId,
		Data:      user.CaptchaValue,
	}) {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10002,
			Message: "验证码校验失败！", //这里有风险
		})
		return
	}

	ret := model.GetUser(user.Name)
	if ret.ID < 1 || ret.Password != tools.Encrypt(user.Password) {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: "帐号密码错误！",
		})
		return
	}

	//生成TOKEN
	token, _ := model.GetJwt(ret.ID, user.Name)
	c.JSON(http.StatusOK, tools.Ecode{
		Message: "登录成功",
		Data:    token,
	})
	return
}

func Logout(c *gin.Context) {
	_ = model.FlushSession(c)
	c.Redirect(302, "/login")
}

func GetCaptcha(context *gin.Context) {
	if !tools.CheckXYZ(context) {
		//fmt.Println("您的手速真是太快了，老师er")
		context.JSON(http.StatusOK, tools.Ecode{
			Code:    10005,
			Message: "您的手速真的是太快了！",
		})
		return
	}
	captcha, err := tools.CaptchaGenerate()
	if err != nil {
		context.JSON(http.StatusOK, tools.Ecode{
			Code:    10005,
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, tools.Ecode{
		Data: captcha,
	})
}

type CUser struct {
	Name      string `json:"name" form:"name"`
	Password  string `json:"password" form:"password"`
	Password2 string `json:"password2" form:"password2"`
	Email     string `json:"email" form:"email"`
	Captcha   string `json:"captcha" form:"captcha"`
}

func CreateUserIndex(c *gin.Context) {
	c.HTML(200, "user_register.tmpl", nil)
}
func CreateUser(c *gin.Context) {
	var user CUser
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(200, tools.Ecode{Message: err.Error()})
	}
	if user.Name == "" || user.Password == "" || user.Password2 == "" {
		c.JSON(200, tools.ParamErr)
	}
	if user.Password != user.Password2 {
		c.JSON(200, tools.Ecode{
			Message: "两次密码不匹配！"})
	}
	if oldUser := model.GetUser(user.Name); oldUser.ID > 0 {
		c.JSON(200, tools.Ecode{
			Message: "用户名已经存在！"}) //这里有风险，且很大，并发安全
	}
	nameLen := len(user.Name)
	//fmt.Println(nameLen)
	password := len(user.Password)
	//fmt.Println(password)
	if nameLen > 16 || nameLen < 8 || password > 16 || password < 8 {
		c.JSON(200, tools.Ecode{
			Message: "账户和密码的长度必须大于8小于16！"})
	}
	//判断密码是不是纯数字，使用正则表达式来进行判断
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(user.Password) {
		c.JSON(200, tools.Ecode{
			Message: "密码不能为纯数字"})
	}
	newUser := model.User{
		Name:        user.Name,
		Password:    tools.Encrypt(user.Password),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		Uid:         tools.GetUid(),
	}
	if err := model.CreateUser(&newUser); err != nil {
		c.JSON(200, tools.Ecode{
			Message: "新用户创建失败"})
	} else {
		c.JSON(200, tools.Ecode{Message: "用户创建成功"})
	}
}
