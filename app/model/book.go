package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// 该部分表示创建一个新的记录并将其插入到指定的数据表中。通过调用 Create() 方法并传入结构体 b，可以将结构体中的字段值映射到数据表的列，并在数据库中插入一条新的记录。
func CreateBook(b *Book) error {
	return Conn.Table("book").Create(b).Error
}

func GetBook(id int64) BookWithInfo {
	var ret Book
	Conn.Table("book").Where("id = ?", id).Find(&ret)
	info := make([]BookInfo, 0)
	err := Conn.Table("book_info").Where("id=?", id).Find(&info).Error
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	return BookWithInfo{ret, info}
}

func GetBooks() ([]*Book, error) {
	var books []*Book
	if err := Conn.Table("book").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
func GetBookInfo(pageSize, offset int) ([]*BookInfo, error) {
	// 构建用于在Redis中存储和检索缓存数据的键
	cacheKey := fmt.Sprintf("bookinfo:%d:%d", pageSize, offset)

	// 在Redis中检查缓存数据是否存在
	cacheData, err := Rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		// 缓存数据存在，直接从Redis中获取数据并返回
		var bookinfo []*BookInfo
		err = json.Unmarshal([]byte(cacheData), &bookinfo)
		if err != nil {
			return nil, err
		}
		return bookinfo, nil
	}
	//缓存数据不存在，从数据库中获取数据并且返回
	var bookinfo []*BookInfo
	if err := Conn.Table("book_info").Where("id>?", offset).Limit(pageSize).Find(&bookinfo).Error; err != nil {
		return nil, err
	}
	fmt.Println(bookinfo)
	// 将获取到的数据存储到Redis中，以供后续使用
	jsonData, err := json.Marshal(bookinfo)
	if err != nil {
		return nil, err
	}
	err = Rdb.Set(context.Background(), cacheKey, jsonData, 240*time.Hour).Err() // 这里设置缓存过期时间为24小时
	if err != nil {
		return nil, err
	}
	return bookinfo, nil
}
func SaveBook(data *Book) error {
	return Conn.Table("book").Save(data).Error
}

func DeleteBook(id int64) error {
	return Conn.Table("book").Where("id = ?", id).Delete(nil).Error
}
func BorrowBook(userId, bookId int64) error {
	tx := Conn.Begin()
	//查询用户是否存在
	var user User
	if err := tx.First(&user, userId).Error; err != nil {
		tx.Rollback()
	}
	//查询图书是否存在，是否正常
	var book Book
	if err := tx.Where("id = ? AND num > 0", bookId).First(&book).Error; err != nil {
		tx.Rollback()
	}
	//创建借阅记录
	now := time.Now()
	bookUser := BookUser{
		UserId:      userId,
		BookId:      bookId,
		Status:      1,  // 1 表示借阅中
		Time:        60, // 借阅时长，以分钟为单位
		CreatedTime: now,
		UpdatedTime: now,
	}
	if err := tx.Create(&bookUser).Error; err != nil {
		tx.Rollback()
	}
	// 更新图书库存（使用乐观锁）gorm.Expr函数接受两个参数：表达式字符串和参数。在这里，表达式字符串是"num - ?"，表示将数据库中的num字段减去一个参数。问号?表示参数占位符，后面的参数1将替换问号。
	//通过使用gorm.Expr("num - ?", 1)，我们可以在更新图书库存时执行原始的SQL表达式，将库存数量减去指定的数量。这样，无论库存数量是多少，都可以通过这个表达式进行减法运算。
	result := tx.Model(&book).Where("id = ? AND num > 0", bookId).Update("num", gorm.Expr("num - ?", 1))
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("concurrent borrow error")
	} //result.RowsAffected == 0的判断用于检查更新操作是否影响了任何行。如果result.RowsAffected的值为0，表示在更新图书库存时没有影响到任何记录，即更新操作没有成功执行。
	tx.Commit()
	return nil
}
func ReturnBook(userId, bookId int64) error {
	tx := Conn.Begin()
	// 查询用户是否存在（加悲观锁）
	var user User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, userId).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 查询借阅记录是否存在（加悲观锁）
	var bookUser BookUser
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ? AND book_id = ? AND status = ?", userId, bookId, 1).First(&bookUser).Error; err != nil {
		tx.Rollback()
		return errors.New("借阅记录不存在或已归还")
	}

	// 查询图书是否存在，是否正常（加悲观锁）
	var book Book
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ? AND num > 0", bookId).First(&book).Error; err != nil {
		tx.Rollback()
		return errors.New("图书信息不存在或库存不足")
	}

	// 更新借阅记录状态为已归还
	if err := tx.Model(&bookUser).Where("user_id = ? AND book_id = ? AND status = ?", userId, bookId, 1).Update("status", 0).Error; err != nil {
		tx.Rollback()
		return errors.New("更新借阅记录失败")
	}

	// 更新图书库存
	if err := tx.Model(&Book{}).Where("id = ?", bookId).UpdateColumn("num", gorm.Expr("num + ?", 1)).Error; err != nil {
		tx.Rollback()
		return errors.New("更新图书库存失败")
	}

	tx.Commit()
	return nil
}
func BuyBook(userId, bookId, buynum int64) error {
	tx := Conn.Begin()
	var user User
	fmt.Println(userId)
	if err := tx.Table("user").Where("id=?", userId).Find(&user).Error; err != nil {
		tx.Rollback()
		return errors.New("用户不存在")
	}
	var book Book
	if err := tx.Table("book").Where("uid=?", bookId).Find(&book).Error; err != nil {
		tx.Rollback()
		return errors.New("您查询的书本不存在")
	}
	if book.Num < 1 || book.Num-buynum < 0 {
		tx.Rollback()
		return errors.New("库存不够了")
	}
	//更新book表中的库存
	if err := tx.Model(&Book{}).Where("uid=?", bookId).UpdateColumn("num", gorm.Expr("num - ?", buynum)).Error; err != nil {
		tx.Rollback()
		return errors.New("库存更新失败")
	}
	//创建买书记录
	booknum := UserBookBuy{
		UserId:      userId,
		BookId:      bookId,
		BuyNum:      buynum,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if err := tx.Table("user_book_buy").Create(&booknum).Error; err != nil {
		tx.Rollback()
		return errors.New("添加购买记录失败")
	}
	tx.Commit()
	return nil
}
func PayRecord(jiage, userId int64, outTrade string) error {
	tx := Conn.Begin()
	//创建购买的记录
	var pay = Alipay{
		UserId:     userId,
		Outtradeno: outTrade,
		Amount:     float64(jiage),
	}
	if err := tx.Table("alipay").Create(&pay).Error; err != nil {
		tx.Rollback()
		return errors.New("创建购买记录失败")
	}
	tx.Commit()
	return nil
}
func RefundBook(bookId, num int64) error {
	tx := Conn.Begin()
	//退款的时候将图书的加上之前购买的数量
	num1 := num * 2
	if err := tx.Model(&Book{}).Where("uid=?", bookId).UpdateColumn("num", gorm.Expr("num + ?", num1)).Error; err != nil {
		tx.Rollback()
		return errors.New("库存更新失败")
	}
	return nil
}
func QueryTradeNo(userId int64) *Alipay {
	tx := Conn.Begin()
	var Pay *Alipay
	if err := tx.Table("alipay").Where("user_id", userId).Find(Pay); err != nil {
		errors.New("查询订单号失败")
	}
	return Pay
}
