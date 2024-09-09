package tools

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(pwd string) string {
	newPwd := pwd + "蜡笔小新"
	hash := md5.New()
	hash.Write([]byte(newPwd))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	//fmt.Printf("加密后的密码:%s", hashString)
	return hashString
}
