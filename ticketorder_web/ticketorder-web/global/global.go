package global

import (
	"gotrains/ticketorder_web/ticketorder-web/config"
	"gotrains/ticketorder_web/ticketorder-web/proto"

	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
)

var (
	Trans           ut.Translator
	ServerConfig    *config.ServerConfig = &config.ServerConfig{}
	Rdb             *redis.Client
	UserClient      proto.UserClient
	TrainClient     proto.TrainClient
	TicketClient    proto.TicketClient
	OrderClient     proto.OrderClient
	NewServerConfig *config.NewServerConfig = &config.NewServerConfig{}
)
