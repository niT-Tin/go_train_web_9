package global

import (
	"gotrains/train_srvs/train_srv/config"
	"gotrains/train_srvs/train_srv/query"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	Config       = &config.ServerConfig{}
	ServerConfig = &config.NewServerConfig{}
	Query        = &query.Query{}
)
