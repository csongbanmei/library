package model

import (
	"fmt"
)

func GetUser(name string) *User {
	var ret User
	err := Conn.Table("user").Where("name= ?", name).Find(&ret).Error
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	return &ret
}
func CreateUser(user *User) error {
	if err := Conn.Table("user").Where("uid", user.Uid).Create(user).Error; err != nil {
		fmt.Printf("err:%s", err)
		return err
	}
	return nil
}
