package services

import (
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *StationService) GetStationsByTrainCode(trainCode string) ([]model.TrainStation, error) {
	var stations []model.TrainStation
	result := global.DB.Where(&model.TrainStation{TrainCode: trainCode}).Find(&stations)
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

func (s *StationService) GenerateStationDaily(stationreq *model.DailyTrainStation) (*model.DailyTrainStation, error) {

	// 删除某日某车次的车站信息
	result := global.DB.Delete(&model.DailyTrainStation{}, "train_code = ? and date = ?", stationreq.TrainCode, stationreq.Date)
	if result.Error != nil {
		return nil, result.Error
	}
	// 获取某日某车次的车站信息
	// dt, _ := time.Parse("2006-01-02", stationreq.Date)
	var stations []model.TrainStation
	zap.S().Infof("trainCode: %s", stationreq.TrainCode)
	res := global.DB.Where(&model.TrainStation{TrainCode: stationreq.TrainCode}).Find(&stations)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(stations) == 0 {
		return nil, status.Errorf(codes.NotFound, "车站信息不存在")
	}
	for _, station := range stations {
		var stationDaily model.DailyTrainStation
		// 多此一举
		tp, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		stationDaily.TrainCode = stationreq.TrainCode
		stationDaily.CreateTime = time.Now()
		stationDaily.UpdateTime = time.Now()
		stationDaily.Date = tp
		stationDaily.Index = station.Index
		stationDaily.Name = station.Name
		stationDaily.NamePinyin = station.NamePinyin
		stationDaily.InTime = station.InTime
		stationDaily.OutTime = station.OutTime
		stationDaily.StopTime = station.StopTime
		stationDaily.Km = station.Km
		result := global.DB.Create(&station)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	zap.S().Infof("车站信息生成成功 date: %s, train_code: %s", stationreq.Date, stationreq.TrainCode)
	return &model.DailyTrainStation{}, nil
}
