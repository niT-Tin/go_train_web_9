package initialize

import (
	"gotrains/ticketorder_web/ticketorder-web/global"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func InitRedisDb() {
	global.Rdb = redis.NewClient(&redis.Options{
		Addr: global.ServerConfig.RedisConfig.Host + ":" + strconv.Itoa(global.ServerConfig.RedisConfig.Port),
	})
}
