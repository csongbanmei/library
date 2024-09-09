package logic

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/skip2/go-qrcode"
	"library/app/model"
	"library/app/tools"
	"library/config"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func BorrowBook(c *gin.Context) {
	//获取用户信息
	// 获取用户ID和图书ID
	userId, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	bookId, _ := strconv.ParseInt(c.PostForm("book_id"), 10, 64)
	//执行借书逻辑
	err := model.BorrowBook(userId, bookId)
	if err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10002,
			Message: err.Error(),
		})
		return
	}
	//返回成功
	c.JSON(http.StatusOK, tools.OK)
}
func ReturnBook(c *gin.Context) {
	// 获取用户ID和图书ID
	userId, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 64)
	bookId, _ := strconv.ParseInt(c.PostForm("book_id"), 10, 64)
	//执行借书逻辑
	err := model.ReturnBook(userId, bookId)
	if err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10002,
			Message: err.Error(),
		})
		return
	}
	//返回成功
	c.JSON(http.StatusOK, tools.OK)
}
func GetEmailIndex(c *gin.Context) {
	c.HTML(200, "email_verification.tmpl", nil)
}

type Cap struct {
	Email   string `json:"email" form:"email"`
	Captcha string `json:"captcha" form:"captcha"'`
}

func GetEmail(c *gin.Context) {
	var captcha Cap
	if err := c.ShouldBind(&captcha); err != nil {
		c.JSON(200, tools.Ecode{Message: err.Error()})
	}
	youxiang := captcha.Email
	//调用model中的GenerateVerificationCode()方法，生成一个随机的六位字符串
	code := tools.GenerateVerificationCode().Code
	// 将验证码存储到 Redis
	err := model.Rdb.Set(context.Background(), youxiang, code, 5*time.Minute).Err()
	if err != nil {
		// 处理存储错误
		c.JSON(200, tools.Ecode{
			Message: "验证码存储失败"})
		return
	}
	err = tools.SendVerificationCode(youxiang, code)
	if err != nil {
		c.JSON(200, tools.Ecode{
			Message: "验证码发送失败"})
	} else {
		c.JSON(200, tools.Ecode{
			Message: "验证码发送成功成功"})
	}
}

func VerifyEmailCode(c *gin.Context) {
	var captcha Cap
	if err := c.ShouldBind(&captcha); err != nil {
		c.JSON(200, tools.Ecode{Message: err.Error()})
		return
	}
	fmt.Println(captcha.Captcha)
	// 从 Redis 中获取存储的验证码
	storedCode, err := model.Rdb.Get(context.Background(), captcha.Email).Result()
	if err != nil {
		// 处理获取验证码失败的情况
		c.JSON(200, tools.Ecode{
			Message: "对不起，您获取验证码失败了"})
		return
	}
	// 验证输入的验证码与存储的验证码是否匹配
	if captcha.Captcha != storedCode {
		c.JSON(200, tools.Ecode{
			Message: "您的验证码不匹配，请您重新输入"})
		return
	}
	// 返回验证成功的响应
	c.JSON(200, tools.Ecode{
		Message: "邮箱验证码验证成功"})
}
func GetBuyBook(c *gin.Context) {
	c.HTML(200, "buy_book.tmpl", nil)
}
func BuyBook(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.PostForm("id"), 10, 64)
	bookId, _ := strconv.ParseInt(c.PostForm("uid"), 10, 64)
	num, _ := strconv.ParseInt(c.PostForm("num"), 10, 64)
	price, _ := strconv.ParseInt(c.PostForm("price"), 10, 64)
	price = price * num
	jiage := strconv.FormatInt(price, 10)
	// 获取url进行支付
	client, err := alipay.NewClient(config.AppId, config.PrivateKey, config.IsProduction)
	if err != nil {
		log.Println("支付宝初始化错误")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "支付宝初始化错误"})
		return
	}
	outTrade := tools.GreateoutTradeNo()
	//生成一个用户的购买记录,用来进行用户的退款操作
	model.PayRecord(price, userId, outTrade)
	client.SetCharset("utf-8").SetSignType(alipay.RSA2).SetNotifyUrl(config.NotifyURL).SetReturnUrl(config.ReturnURL)
	bm := make(gopay.BodyMap)
	//这段代码创建了一个gopay.BodyMap类型的变量bm，用于存储支付宝支付接口的请求参数。
	//gopay.BodyMap是一个用于存储HTTP请求参数的字典类型，常用于构建支付宝接口请求的参数。
	//通过make(gopay.BodyMap)，我们创建了一个空的gopay.BodyMap变量bm，可以使用它来添加和获取支付宝接口的请求参数。
	bm.Set("subject", "这里是小陈的支付页面")
	bm.Set("out_trade_no", outTrade)
	bm.Set("total_amount", jiage)
	bm.Set("product_code", config.ProductCode)
	payUrl, err := client.TradePagePay(context.Background(), bm)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "支付链接生成失败"})
		return
	}

	// 更新购买逻辑，例如生成订单、更新库存等
	err = model.BuyBook(userId, bookId, num)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}

	// 返回购买成功的 JSON 响应，并返回支付链接
	c.JSON(http.StatusOK, gin.H{"message": "恭喜您，买书成功", "payUrl": payUrl})
}
func GetRefundBook(c *gin.Context) {
	c.HTML(200, "return_pay.tmpl", nil)
}
func RefundBook(c *gin.Context) {
	//userId, _ := strconv.ParseInt(c.PostForm("id"), 10, 64)
	bookId, _ := strconv.ParseInt(c.PostForm("uid"), 10, 64)
	num, _ := strconv.ParseInt(c.PostForm("num"), 10, 64)
	// 初始化支付宝客户端
	client, err := alipay.NewClient(config.AppId, config.PrivateKey, config.IsProduction)
	if err != nil {
		log.Println("支付宝初始化错误")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "支付宝初始化错误"})
		return
	}
	// 构建退款请求参数
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", "171275678974")
	bm.Set("refund_amount", "60")
	bm.Set("out_request_no", tools.GreateoutTradeNo())
	// 发起退款请求
	resp, err := client.TradePageRefund(context.Background(), bm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to request refund"})
		return
	}
	//进行增加图书库存的操作
	model.RefundBook(bookId, num)
	// 处理退款响应
	c.JSON(http.StatusOK, gin.H{"message": "Refund request successful", "response": resp})
}
func Wechatcall(c *gin.Context) {
	c.HTML(200, "wechatlogin.tmpl", nil)
}

// TOKEN 假设您在Go代码中定义了一个名为TOKEN的常量，用于存储您的令牌值
const TOKEN = "1234"

// 配置公众号的token
func WechatVerify(c *gin.Context) {
	// 获取微信发送的参数
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	// 将 token, timestamp 和 nonce 按照字典序排序
	params := []string{TOKEN, timestamp, nonce}
	sort.Strings(params)

	// 将排序后的三个参数拼接成一个字符串
	rawString := strings.Join(params, "")

	// 对拼接后的字符串进行SHA1加密
	hash := sha1.New()
	hash.Write([]byte(rawString))
	hashString := hex.EncodeToString(hash.Sum(nil))

	// 对比微信传递过来的签名和我们生成的签名
	if hashString == signature {
		// 签名验证成功，返回echostr
		c.String(http.StatusOK, echostr)
	} else {
		// 签名验证失败，返回错误信息
		c.String(http.StatusUnauthorized, "Signature verification failed")
	}
}

func RandString(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func Redirect(c *gin.Context) {
	path := c.Query("Url")
	fmt.Sprintf(path)
	//防止跨站请求伪造攻击 增加安全性
	redirectURL := url.QueryEscape("http://p622ma.natappfree.cc/user/buybook") //userinfo,
	wechatLoginURL := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&state=%s&scope=snsapi_userinfo#wechat_redirect", "wxb81f9c8acb700297", redirectURL, "aaaaa")
	wechatLoginURL, _ = url.QueryUnescape(wechatLoginURL)
	// 生成二维码
	qrCode, err := qrcode.Encode(wechatLoginURL, qrcode.Medium, 256)
	if err != nil {
		// 错误处理
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}
	// 将二维码图片作为响应返回给用户
	c.Header("Content-Type", "image/png")
	c.Writer.Write(qrCode)
}

func Callback(c *gin.Context) {
	// 获取微信返回的授权码
	code := c.Query("code")
	// 向微信服务器发送请求，获取access_token和openid
	tokenResp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", "wxb81f9c8acb700297", "cbc21569030393ad2f2932602773f0e9", code))
	if err != nil {
		fmt.Println(err)
		resp := tools.Ecode{
			Data:    nil,
			Message: "error,获取token失败",
			Code:    10000,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	// 解析响应中的access_token和openid
	var tokenData struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
	}
	if err1 := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err1 != nil {
		resp := &tools.Ecode{
			Data:    nil,
			Message: "error,获取token失败",
			Code:    10000,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", tokenData.AccessToken, tokenData.OpenID)
	userInfoResp, err := http.Get(userInfoURL)
	if err != nil {
		// 错误处理
		log.Println("获取失败")
		return
	}
	defer userInfoResp.Body.Close()

	//------------------------------------
	var userData struct {
		OpenID   string `json:"openid"`
		Nickname string `json:"nickname"`
	}
	if err1 := json.NewDecoder(userInfoResp.Body).Decode(&userData); err1 != nil {
		// 错误处理
		log.Println("获取用户信息失败")
		return
	}
	////用户的名字
	//var user1 model.User
	//nickname := userData.Nickname
	//if err2 := mysql.DB.Where("user_name=?", nickname).First(&user1).Error; err2 != nil {
	//	if errors.Is(err2, gorm.ErrRecordNotFound) {
	//		user1.UserName = nickname
	//		user1.UserID, _ = snowflake.GetID()
	//		user1.Identity = "普通用户"
	//	} else {
	//		zap.L().Error("验证登录信息过程中出错")
	//		ResponseError(c, CodeServerBusy)
	//		return
	//	}
	//}
	////添加jwt验证
	//atoken, rtoken, err3 := jwt.Token.GetToken(uint64(user1.UserID), user1.UserName, user1.Identity)
	//if err3 != nil {
	//	zap.L().Error("生成JWT令牌失败")
	//	return
	//}
	//c.Header("Authorization", atoken)
	////发送成功响应
	//ResponseSuccess(c, &model.LoginData{
	//	AccessToken:  atoken,
	//	RefreshToken: rtoken,
	//})
	//zap.L().Info("登录成功")
	//return
}
