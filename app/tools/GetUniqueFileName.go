package tools

import (
	"path/filepath"
	"strconv"
	"time"
)

// GetUniqueFileName 生成唯一的文件名
func GetUniqueFileName(originalFilename string) string {
	// 获取文件的扩展名
	ext := filepath.Ext(originalFilename)

	// 生成时间戳作为文件名的一部分
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	// 组合文件名
	uniqueFilename := timestamp + ext

	return uniqueFilename
}
