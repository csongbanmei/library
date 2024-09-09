package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library/app/model"
	"library/app/tools"
	"net/http"
	"time"
)

type Admin struct {
	ID           int64     `json:"id" gorm:"id" form:"id"`
	Name         string    `json:"name" gorm:"name" form:"name"`
	Password     string    `json:"password" gorm:"password" form:"password"`
	CaptchaId    string    `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string    `json:"captcha_value" form:"captcha_value"`
	CreatedTime  time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime  time.Time `json:"updated_time" gorm:"updated_time"`
}

func GetAdminLogin(c *gin.Context) {
	c.HTML(200, "admin_login.tmpl", nil)
}
func AdminLogin(c *gin.Context) {
	var admin Admin
	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: err.Error(), //这里有风险
		})
	}

	fmt.Printf("user:%+v\n", admin)

	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: admin.CaptchaId,
		Data:      admin.CaptchaValue,
	}) {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10002,
			Message: "验证码校验失败！", //这里有风险
		})
		return
	}

	ret := model.GetAdmin(admin.Name)
	if ret.ID < 1 || tools.Encrypt(ret.Password) != tools.Encrypt(admin.Password) {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: "帐号密码错误！",
		})
		return
	}
	_ = model.SetSession(c, admin.Name, ret.ID)
	c.JSON(http.StatusOK, tools.Ecode{
		Message: "登录成功",
	})
}

func AdminLogout(c *gin.Context) {
	_ = model.FlushSession(c)
	c.JSON(http.StatusUnauthorized, tools.Ecode{
		Code:    0,
		Message: "您已退出登录",
	})
}
