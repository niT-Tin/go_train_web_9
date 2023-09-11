package global

import (
	"gotrains/ticketorder_srvs/ticket_srv/config"
	"gotrains/ticketorder_srvs/ticket_srv/proto"

	"gorm.io/gorm"
)

var (
	DB                    *gorm.DB
	Config                = &config.ServerConfig{}
	ServerConfig          = &config.NewServerConfig{}
	TicketInventoryClient proto.TicketClient
	UserClient            proto.UserClient
)
