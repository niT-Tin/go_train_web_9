package global

import (
	"gotrains/userpassenger_web/user-web/config"
	"gotrains/userpassenger_web/user-web/proto"

	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
)

var (
	Trans           ut.Translator
	ServerConfig    *config.ServerConfig = &config.ServerConfig{}
	Rdb             *redis.Client
	UserClient      proto.UserClient
	NewServerConfig *config.NewServerConfig = &config.NewServerConfig{}
)
