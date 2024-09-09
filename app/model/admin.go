package model

import "fmt"

func GetAdmin(name string) *Admin {
	var ret Admin
	err := Conn.Table("admin").Where("name= ?", name).Find(&ret).Error
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	return &ret
}
