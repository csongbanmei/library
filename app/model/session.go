package model

//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/gorilla/sessions"
//)
//
//var store = sessions.NewCookieStore([]byte("香香编程喵喵喵"))
//var sessionName = "session-name"
//
//func GetSession(c *gin.Context) map[interface{}]interface{} {
//	session, _ := store.Get(c.Request, sessionName)
//	fmt.Printf("session:%+v\n", session.Values)
//	return session.Values
//}
//
//func SetSession(c *gin.Context, name string, id int64) error {
//	session, _ := store.Get(c.Request, sessionName)
//	session.Values["name"] = name
//	session.Values["id"] = id
//	return session.Save(c.Request, c.Writer)
//}
//
//// 清除会话数据中的闪存数据
//func FlushSession(c *gin.Context) error {
//	session, _ := store.Get(c.Request, sessionName)
//	fmt.Printf("session : %+v\n", session.Values)
//	session.Flashes()
//	return session.Save(c.Request, c.Writer)
//}
