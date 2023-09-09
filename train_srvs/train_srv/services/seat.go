package services

import (
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"time"
)

type SeatService struct{}

func (s *SeatService) CountSeat(date string, trainCode string, typ global.SeatType) int64 {
	var count int64
	dt, _ := time.Parse("2006-01-02", date)
	global.DB.Where(&model.DailyTrainSeat{Date: dt, TrainCode: trainCode, SeatType: string(typ)}).Count(&count)
	return count
}
