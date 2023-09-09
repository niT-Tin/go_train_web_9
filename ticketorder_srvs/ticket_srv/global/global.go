package global

import (
	"gotrains/ticketorder_srvs/ticket_srv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	Config       = &config.ServerConfig{}
	ServerConfig = &config.NewServerConfig{}
)
