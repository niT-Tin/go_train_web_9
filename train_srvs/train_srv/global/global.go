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

type SeatType string

type TrainType string

type SeatFarePerKm float32

type TrainTypeRate float32

type SeatTypeCode int32

type TrainTypeCode int32

var (
	SeatTypeFirstClass  SeatType = "一等座"
	SeatTypeSecondClass SeatType = "二等座"
	SeatTypeSoftBerth   SeatType = "软卧"
	SeatTypeHardBerth   SeatType = "硬卧"
)

var (
	SeatTypeFirstClassCode  SeatTypeCode = 1
	SeatTypeSecondClassCode SeatTypeCode = 2
	SeatTypeSoftBerthCode   SeatTypeCode = 3
	SeatTypeHardBerthCode   SeatTypeCode = 4
)

var (
	TrainTypeGTrain TrainType = "G字头列车"
	TrainTypeDTrain TrainType = "D字头列车"
	TrainTypeKTrain TrainType = "K字头列车"
)

var (
	TrainTypeGTrainCode TrainTypeCode = 1
	TrainTypeDTrainCode TrainTypeCode = 2
	TrainTypeKTrainCode TrainTypeCode = 3
)

var (
	TrainTypeGTrainRate TrainTypeRate = 1.3
	TrainTypeDTrainRate TrainTypeRate = 1
	TrainTypeKTrainRate TrainTypeRate = 0.8
)

// 不同座位每公里票价

var (
	SeatFarePerKmFirstClass  SeatFarePerKm = 0.5
	SeatFarePerKmSecondClass SeatFarePerKm = 0.4
	SeatFarePerKmSoftberth   SeatFarePerKm = 0.7
	SeatFarePerKmHardberth   SeatFarePerKm = 0.8
)
