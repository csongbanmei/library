package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library/app/model"
	"library/app/tools"
	"net/http"
	"strconv"
)

var Page int

func Index(c *gin.Context) {
	pageSize := 10                                   // 每页显示的记录数量
	pageNumber, err := strconv.Atoi(c.Query("page")) // 从请求参数中获取页码
	Page = pageNumber
	if err != nil || pageNumber < 1 {
		pageNumber = 1 // 默认为第一页
	}
	offset := (pageNumber - 1) * pageSize // 计算偏移量
	ret, _ := model.GetBookInfo(pageSize, offset)
	c.HTML(200, "index.tmpl", gin.H{"book": ret})
}
func UserIndex(context *gin.Context) {
	context.HTML(200, "user_index.tmpl", nil)
}
func GetBook(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" || idStr == "0" {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	ret := model.GetBook(id)
	c.JSON(http.StatusOK, tools.Ecode{
		Data: ret,
	})
	//fmt.Println(ret)

}

func GetBooks(c *gin.Context) {
	ret, err := model.GetBooks()
	if err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tools.Ecode{
		Data: ret,
	})
	fmt.Println(ret)
	return
}
func GetBookInfo(c *gin.Context) {
	pageSize := 10                                  // 每页显示的记录数量
	offset := Page                                  //拿到游标
	ret, err := model.GetBookInfo(pageSize, offset) // 调用 model.GetBookInfo() 函数，并传递分页参数
	if err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.Ecode{
		Data: ret,
	})
	return
}
func AddBook(c *gin.Context) {
	var data model.Book
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	//TODO:增加参数校验
	//id, _ := strconv.ParseInt(idStr, 10, 64)
	err := model.CreateBook(&data)
	if err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.OK)
	return
}

func DelBook(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" || idStr == "0" {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	err := model.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.OK)
	return
}

func SaveBook(c *gin.Context) {
	var data model.Book
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	//TODO:增加参数校验
	err := model.SaveBook(&data)
	if err != nil {
		c.JSON(http.StatusOK, tools.Ecode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.Ecode{
		Code: 0,
	})
	return
}
