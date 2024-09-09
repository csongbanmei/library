package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"library/app/model"
	"time"
)

func CheckXYZ(context *gin.Context) bool {
	//拿到IP和UA
	ip := context.ClientIP()
	ua := context.GetHeader("user-agent")
	fmt.Printf("ip:%s\nua:%s\n", ip, ua)
	//转下MD5
	hash := md5.New()
	hash.Write([]byte(ip + ua))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	fmt.Printf(hashString)
	//校验是否被ban
	flag, _ := model.Rdb.Get(context, "ban-"+hashString).Bool()
	if flag {
		return false
	}

	i, _ := model.Rdb.Get(context, "xyz-"+hashString).Int()
	fmt.Printf("i:%d\n", i)
	if i > 5 {
		model.Rdb.SetEx(context, "ban-"+hashString, true, 30*time.Minute)
		//SetEx 是 Redis 的一个命令，用于设置一个键值对，并指定键的过期时间。
		//context 是上下文对象，可能用于跟踪操作的上下文信息。
		//"ban-"+hashString 是键的名称，使用了 hashString 变量，并在字符串前面添加了 "ban-" 前缀。
		//true 是键的值，这里设置为布尔值 true。
		//30*time.Minute 是键的过期时间，表示 30 分钟。time.Minute 是 Go 语言中的时间单位，表示分钟。
		//综上所述，这段代码的功能是将一个键值对存储到 Redis 数据库中，并设置键的过期时间为 30 分钟。通常情况下，存储键值对并设置过期时间可以用于实现缓存、限流、会话管理等功能。
		return false
	}

	model.Rdb.Incr(context, "xyz-"+hashString)
	//具体来说，"xyz-"+hashString 是要自增的键，它可能代表了一个计数器或某种状态的值。当这行代码执行时，它会将指定键的值增加1。
	//例如，如果原先该键的值是3，执行了model.Rdb.Incr(context, "xyz-"+hashString) 后，该键的值将变为4
	//通常情况下，自增操作用于统计或计数的场景，可以方便地对某个变量或状态进行累加。在这段代码中，该自增操作可能用于记录某个客户端或用户的请求次数、操作次数或其他与计数相关的信息
	fmt.Println(i)
	model.Rdb.Expire(context, "xyz-"+hashString, 50*time.Minute)
	//具体来说，"xyz-"+hashString 是要设置过期时间的键，它可能代表了某个特定的数据或状态。50*time.Minute 表示过期时间为50分钟。
	//当这行代码执行时，它会将指定键的过期时间设置为50分钟。在50分钟之后，如果对该键进行读取操作，数据库可能会返回空值（或其他指示该键已过期的信息），表示该键已经失效。
	//通过设置过期时间，可以实现数据的自动清理和过期管理。这在一些场景下非常有用，例如缓存数据、临时会话管理等。在这段代码中，该设置过期时间的操作可能用于控制某个数据或状态的生命周期，并在一定时间后自动清理或失效。
	return true
}
