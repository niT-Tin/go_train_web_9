package services

import (
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"time"
)

type StationService struct{}

func (s *StationService) GetStationsByTrainCodeDaily(trainCode string, date string) ([]model.DailyTrainStation, error) {
	var stations []model.DailyTrainStation
	dt, _ := time.Parse("2006-01-02", date)
	result := global.DB.Where(&model.DailyTrainStation{TrainCode: trainCode, Date: dt}).Find(&stations)
	if result.Error != nil {
		return nil, result.Error
	}
	return stations, nil
}

func (s *StationService) GetRate(typ global.SeatType) float32 {
	rate := float32(1)
	switch typ {
	case global.SeatTypeFirstClass:
		// 一等座
		rate = 1.5
	case global.SeatTypeSecondClass:
		// 二等座
		rate = 1
	case global.SeatTypeSoftBerth:
		// 软卧
		rate = 1.1
	case global.SeatTypeHardBerth:
		// 硬卧
		rate = 1.2
	}
	return rate
}
