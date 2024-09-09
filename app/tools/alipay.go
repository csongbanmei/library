package tools

import (
	"fmt"
	"time"
)

// 用户买书的时候生成为唯一的订单号
// 首先，time.Now().UnixMilli()获取当前时间的毫秒级时间戳，表示从Unix纪元（1970年1月1日UTC）到当前时间的毫秒数，并将其赋值给变量ts。
// 然后，fmt.Sprintf("%d", ts)将变量ts格式化为一个字符串，使用%d作为占位符来表示一个整数，并将格式化后的字符串赋值给变量outTradeNo。
// 通过将当前时间的毫秒级时间戳转换为字符串形式，可以在生成订单号时使用它作为唯一标识，以确保每个订单号的唯一性。在实际应用中，通常会将订单号与其他订单相关的信息进行关联，如商品、用户等，以便进行订单的处理和跟踪。
func GreateoutTradeNo() string {
	ts := time.Now().UnixMilli()
	outTradeNo := fmt.Sprintf("%d", ts)
	return outTradeNo
}
