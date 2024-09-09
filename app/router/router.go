package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"library/app/logic"
	"library/app/model"
	"library/app/tools"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	r.GET("/upload", logic.UploadIndex)
	r.POST("/upload", logic.UploadHandler)

	r.GET("/index", logic.Index) //静态页面
	r.GET("user_index", logic.UserIndex)
	r.Static("/static", "app/images")
	user := r.Group("/user")
	//user.Use(CheckUser)
	user.GET("", logic.Wechatcall)
	user.GET("/wechat", logic.WechatVerify)
	user.GET("/wechat/login", logic.Redirect)
	user.GET("wechat/callback", logic.Callback)
	user.GET("/buybook", logic.GetBuyBook)
	user.POST("/buybook", logic.BuyBook)
	user.GET("/returnpay", logic.GetRefundBook)
	user.POST("/returnpay", logic.RefundBook)
	user.POST("/aliyun_yanzheng", logic.SendCode)
	user.GET("/create_email", logic.GetEmailIndex)
	user.POST("verify_email", logic.VerifyEmailCode)
	user.POST("/create_email", logic.GetEmail)
	user.GET("/create_user", logic.CreateUserIndex)
	user.POST("/create_user", logic.CreateUser)
	user.GET("/login", logic.GetLogin)
	user.POST("/login", logic.Login)
	//user.GET("/book/borrow", logic.GetBorrowBook)
	user.POST("/book/borrow", logic.BorrowBook)
	user.POST("/book/return", logic.ReturnBook)

	admin := r.Group("/admin")
	admin.GET("/login", logic.GetAdminLogin)
	admin.POST("/login", logic.AdminLogin)
	admin.Use(CheckAdmin)
	admin.GET("/logout", logic.AdminLogout)

	book := r.Group("/book")
	book.GET("", logic.GetBook)
	book.GET("list", logic.GetBooks)
	book.GET("listinfo", logic.GetBookInfo)
	book.POST("", logic.AddBook)
	book.PUT("", logic.SaveBook)
	book.DELETE("", logic.DelBook)

	//验证码
	{
		r.GET("/captcha", logic.GetCaptcha)
		r.POST("/captcha/verify", func(context *gin.Context) {
			var param tools.CaptchaData
			if err := context.ShouldBind(&param); err != nil {
				context.JSON(http.StatusOK, tools.ParamErr)
				return
			}

			fmt.Printf("参数为：%+v", param)
			if !tools.CaptchaVerify(param) {
				context.JSON(http.StatusOK, tools.Ecode{
					Code:    10008,
					Message: "验证失败",
				})
				return
			}
			context.JSON(http.StatusOK, tools.OK)
		})

	}

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 创建通道来接收关闭信号
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// 启动服务器（非阻塞）
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %s\n", err.Error())
		}
	}()

	// 等待关闭信号
	<-shutdown

	// 创建一个上下文对象，设置超时时间以优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("关闭服务器失败: %s\n", err.Error())
	}

	log.Println("服务器已优雅关闭")
}

func CheckUser(c *gin.Context) {
	var name string
	var id int64
	jwt := c.GetString("auth")
	if jwt == "" {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10003,
			Message: "缺少token",
		})
		//c.Abort()
		return
	}

	d, err := model.CheckJwt(jwt)
	if err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10003,
			Message: "token校验失败",
		})
		c.Abort()
		return
	}

	name = d.Name
	id = d.Id

	if id <= 0 || len(name) == 0 {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10003,
			Message: "用户信息错误",
		})
		c.Abort()
		return
	}

	c.Next()
}

func CheckAdmin(context *gin.Context) {
	var name string
	var id int64
	values := model.GetSession(context)
	if v, ok := values["id"]; ok {
		id = v.(int64)
	}
	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if name == "" || id <= 0 {
		context.JSON(http.StatusUnauthorized, tools.NotLogin)
		context.Abort()
	}
	context.Next()
}

//如果强制的关闭gin框架：请求中断：如果服务器被强制关闭，正在处理的请求将会被中断，导致客户端接收到不完整或错误的响应。
//连接泄漏：直接关闭服务器可能导致持续活动的连接被意外终止，而不进行适当地清理。这可能导致资源泄漏，例如未关闭的数据库连接或文件描述符。
