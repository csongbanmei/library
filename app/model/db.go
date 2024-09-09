package model

import (
	"context"
	"fmt"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Conn *gorm.DB
var Rdb *redis.Client
var Mdb *mongo.Client

func NewMysql() {
	username := "root"           //账号
	password := "chenchiyuan123" //密码
	host := "127.0.0.1"          //数据库地址，可以是Ip或者域名
	port := 3306                 //数据库端口
	Dbname := "library"          //数据库名
	timeout := "10s"             //连接超时，10秒
	my := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。此操作返回两个参数：*gormDB,和error
	//*gormDB表示数据库连接的对象，可以用于执行数据库操作，如：增删改查。error表示函数执行中的错误信息，如果数据库连接成功，则错误
	//值为nil，否则会返回相应的错误信息。
	conn, err := gorm.Open(mysql.Open(my))
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	Conn = conn
}
func NewRdb() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.141.98:6379", //本机的配置，可以修改
		Password: "",
		DB:       0,
	})
	Rdb = rdb
	//初始化session
	store, _ = redisstore.NewRedisStore(context.TODO(), Rdb)
	return
}
func NewMongDB() {
	// 设置MongoDB连接选项
	clientOptions := options.Client().ApplyURI("mongodb://192.168.141.98:27017")

	// 连接到MongoDB
	mdb, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = mdb.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// 将MongoDB客户端赋值给全局变量
	Mdb = mdb
}
func Close() {
	db, _ := Conn.DB()
	_ = db.Close()
	_ = Rdb.Close()
	//关闭mongdb数据库
	if Mdb != nil {
		err := Mdb.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("MongoDB connection closed.")
	}
}
