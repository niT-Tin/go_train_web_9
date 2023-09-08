package global

import (
	"gotrains/userpassenger_srvs/user_srv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	Config       = &config.ServerConfig{}
	ServerConfig = &config.NewServerConfig{}
)
