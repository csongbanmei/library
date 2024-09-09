package app

import (
	"library/app/model"
	"library/app/router"
)

func Start() {
	model.NewMysql()
	model.NewRdb()
	//model.NewMongDB()
	defer func() {
		model.Close()
	}()
	router.New()
}
