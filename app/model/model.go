package model

import "time"

// Admin 管理员
type Admin struct {
	ID          int64     `json:"id" gorm:"id" form:"id"`
	Name        string    `json:"name" gorm:"name" form:"name"`
	Password    string    `json:"password" gorm:"password" form:"password"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time" form:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time" form:"updated_time"`
}

// User 用户
type User struct {
	ID          int64     `json:"id" gorm:"id" form:"id"`
	Uid         int64     `json:"uid" gorm:"uid" form:"uid"`
	Name        string    `json:"name" gorm:"name" form:"name"`
	Password    string    `json:"password" gorm:"password" form:"password"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time" form:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time" form:"updated_time"`
}

// Book 图书情况
type Book struct {
	ID          int64     `json:"id" gorm:"id" form:"id"`
	Uid         int64     `json:"uid" gorm:"uid" form:"uid"`
	Name        string    `json:"name" gorm:"name" form:"name"`
	Cate        string    `json:"cate" gorm:"cate" form:"cate"`
	Status      int64     `json:"status" gorm:"status" form:"status"`
	Num         int64     `json:"num" gorm:"num" form:"num"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time" form:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time" form:"updated_time"`
}

// BookUser  图书用户信息
type BookUser struct {
	ID          int64     `json:"id" gorm:"id" form:"id"`
	UserId      int64     `json:"user_id" gorm:"user_id" form:"user_id"`
	BookId      int64     `json:"book_id" gorm:"book_id" form:"book_id"`
	Status      int64     `json:"status" gorm:"status" form:"status"`
	Time        int64     `json:"time" gorm:"time" form:"time"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time" form:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time" form:"updated_time"`
}

// BookInfo 图书详细信息
type BookInfo struct {
	ID                 int64     `json:"id" gorm:"id" form:"id"`                                                    // 书的id
	BookName           string    `json:"book_name" gorm:"book_name" form:"book_name"`                               // 书名
	Author             string    `json:"author" gorm:"author" form:"author"`                                        // 作者
	PublishingHouse    string    `json:"publishing_house" gorm:"publishing_house" form:"publishing_house"`          // 出版社
	Translator         string    `json:"translator" gorm:"translator" form:"translator"`                            // 译者
	PublishDate        time.Time `json:"publish_date" gorm:"publish_date" form:"publish_date"`                      // 出版时间
	Pages              int64     `json:"pages" gorm:"pages" form:"pages"`                                           // 页数
	Isbn               string    `json:"ISBN" gorm:"ISBN" form:"ISBN"`                                              // ISBN号码
	Price              float64   `json:"price" gorm:"price" form:"price"`                                           // 价格
	BriefIntroduction  string    `json:"brief_introduction" gorm:"brief_introduction" form:"brief_introduction"`    // 内容简介
	AuthorIntroduction string    `json:"author_introduction" gorm:"author_introduction" form:"author_introduction"` // 作者简介
	ImgUrl             string    `json:"img_url" gorm:"img_url" form:"img_url"`                                     // 封面地址
	DelFlg             int64     `json:"del_flg" gorm:"del_flg" form:"del_flg"`                                     // 删除标识
}
type BookWithInfo struct {
	Book     Book
	BookInfo []BookInfo
}
type UserBookBuy struct {
	ID          int64     `json:"id" gorm:"id"`
	UserId      int64     `json:"user_id" gorm:"user_id"`
	BookId      int64     `json:"book_id" gorm:"book_id"`
	BuyNum      int64     `json:"buy_num" gorm:"buy_num"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"`
}

// TableName 表名称
func (*UserBookBuy) TableName() string {
	return "user_book_buy"
}

type Verificationcode struct {
	Email        string    `gorm:"column:email;type:varchar(255);primary_key" json:"email"`
	Code         string    `gorm:"column:code;type:varchar(255)" json:"code"`
	Expirationat time.Time `gorm:"column:expirationat;type:date" json:"expirationat"`
}

func (m *Verificationcode) TableName() string {
	return "verificationcode"
}

// Alipay undefined
type Alipay struct {
	Outtradeno string  `json:"outtradeno" gorm:"outtradeno" form:"outtradeno"`
	Amount     float64 `json:"amount" gorm:"amount" form:"amount"`
	UserId     int64   `json:"user_id" gorm:"user_id" form:"user_id"`
}

// TableName 表名称
func (*Alipay) TableName() string {
	return "alipay"
}

func (*BookInfo) TableName() string {
	return "book_info"
}
func (*BookUser) TableName() string {
	return "book_user"
}
func (*Book) TableName() string {
	return "book"
}
func (*User) TableName() string {
	return "user"
}
func (*Admin) TableName() string {
	return "admin"
}
