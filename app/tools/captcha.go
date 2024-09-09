package tools

import (
	"github.com/mojocn/base64Captcha"
)

type CaptchaData struct {
	CaptchaId string `json:"captcha_id"`
	Data      string `json:"data"`
}

type driverString struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverString  *base64Captcha.DriverString  //字符串
	DriverChinese *base64Captcha.DriverChinese //中文
	DriverMath    *base64Captcha.DriverMath    //数学
	DriverDigit   *base64Captcha.DriverDigit   //数字
}

// 数字驱动
var digitDriver = base64Captcha.DriverDigit{
	Height:   50,  //生成图片高度
	Width:    150, //生成图片宽度
	Length:   5,   //验证码长度
	MaxSkew:  1,   //文字的倾斜度 越大倾斜越狠，越不容易看懂
	DotCount: 1,   //背景的点数，越大，字体越模糊
}

// 使用内存驱动，相关数据会存在内存空间里，创建了一个基于内存存储的验证码存储对象
var store = base64Captcha.DefaultMemStore

// 生成验证码的功能
func CaptchaGenerate() (CaptchaData, error) {
	var ret CaptchaData
	c := base64Captcha.NewCaptcha(&digitDriver, store)
	//创建了一个验证码实例 c，用于生成和验证验证码。通过 &digitDriver 参数指定了验证码的类型为数字验证码，即验证码由数字组成。
	//通过 store 参数指定了验证码的存储方式，即验证码数据将通过 store 对象进行存储和管理。
	id, b64s, _, err := c.Generate()
	//id：验证码的唯一标识符，用于在存储中标识该验证码。
	//b64s：生成的验证码的base64编码字符串，可以是验证码图片的base64编码或验证码文本的base64编码。
	if err != nil {
		return ret, err
	}

	ret.CaptchaId = id
	ret.Data = b64s
	return ret, nil
}

func CaptchaVerify(data CaptchaData) bool {
	return store.Verify(data.CaptchaId, data.Data, true)
}
