package global

import (
	"gotrains/train_webs/train_web/config"
	"gotrains/train_webs/train_web/proto"

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
	SeatClient      proto.SeatClient
	StationClient   proto.StationClient
	NewServerConfig *config.NewServerConfig = &config.NewServerConfig{}
)
