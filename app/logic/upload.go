package logic

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/app/model"
	"net/http"
	"time"
)

type Image struct {
	ID       string    `bson:"_id,omitempty"`
	Filename string    `bson:"filename"`
	FilePath string    `bson:"filepath"`
	UploadAt time.Time `bson:"uploadat"`
}

func UploadIndex(c *gin.Context) {
	c.HTML(200, "image.tmpl", nil)
}
func UploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving the file")
		return
	}

	files := form.File["image"]
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "No file uploaded")
		return
	}

	file := files[0]

	// 获取对应的集合
	collection := model.Mdb.Database("Image").Collection("images")

	// 构造要插入的文档
	image := Image{
		Filename: file.Filename,
		FilePath: file.Filename, // 这里将文件名作为图片地址存储
		UploadAt: time.Now(),
	}

	// 插入文档到 MongoDB
	_, err = collection.InsertOne(context.Background(), image)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving the file information")
		return
	}
	// 构建目标文件路径
	targetPath := "./uploads/" + file.Filename

	// 将文件保存到目标路径
	err = c.SaveUploadedFile(file, targetPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving the file")
		return
	}

	c.String(http.StatusOK, "File uploaded successfully")
}
