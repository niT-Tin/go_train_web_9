package services

import (
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CarriageService struct{}

func (c *CarriageService) GenerateDailyCarriage(trainCode, date string) error {
	// 删除某日某车次的车厢信息
	result := global.DB.Delete(&model.DailyTrainCarriage{}, "train_code = ? and date = ?", trainCode, date)
	if result.Error != nil {
		return result.Error
	}
	// 获取某日某车次的车站信息
	// dt, _ := time.Parse("2006-01-02", carriagereq.Date)
	var carriages []model.TrainCarriage
	res := global.DB.Where(&model.TrainCarriage{TrainCode: trainCode}).Find(&carriages)
	if res.Error != nil {
		return res.Error
	}
	if len(carriages) == 0 {
		return status.Errorf(codes.NotFound, "车站信息不存在")
	}
	for _, carriage := range carriages {
		var carriageDaily model.DailyTrainCarriage
		// 多此一举
		tp, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		carriageDaily.TrainCode = carriage.TrainCode
		carriageDaily.CreateTime = time.Now()
		carriageDaily.UpdateTime = time.Now()
		carriageDaily.Date = tp
		carriageDaily.Index = carriage.Index
		carriageDaily.SeatType = carriage.SeatType
		carriageDaily.SeatCount = carriage.SeatCount
		carriageDaily.RowCount = carriage.RowCount
		carriageDaily.ColCount = carriage.ColCount
		result := global.DB.Create(&carriageDaily)
		if result.Error != nil {
			return result.Error
		}
	}
	zap.S().Infof("车站信息生成成功 date: %s, train_code: %s", date, trainCode)
	return nil
}

func (c *CarriageService) GetCarriageByTrainCode(trainCode string) ([]*model.TrainCarriage, error) {
	var carriages []*model.TrainCarriage
	result := global.DB.Where(&model.TrainCarriage{TrainCode: trainCode}).Find(&carriages)
	if result.Error != nil {
		return nil, result.Error
	}
	return carriages, nil
}
