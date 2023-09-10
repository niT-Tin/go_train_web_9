package services

import (
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TrainService struct{}

func (s *TrainService) GetTrainByCode(code string) (model.Train, error) {
	var train model.Train
	result := global.DB.Where(&model.Train{Code: code}).First(&train)
	if result.Error != nil {
		return train, result.Error
	}
	return train, nil
}

func (s *TrainService) GetRate(typ global.TrainType) float32 {
	rate := float32(1)
	switch typ {
	case global.TrainTypeGTrain:
		// 高铁
		rate = 1.5
	case global.TrainTypeDTrain:
		// 动车
		rate = 1.1
	case global.TrainTypeKTrain:
		// 绿皮车
		rate = 0.8
	}
	return rate
}

func (s *TrainService) GenerateDaily(date string) error {
	return nil
}

func (s *TrainService) GenerateTrainDaily(trainCode, date string) error {
	// 生成车次数据
	var train model.Train
	res := global.DB.Where(&model.Train{Code: trainCode}).First(&train)
	if res.RowsAffected == 0 && len(train.Code) == 0 {
		return status.Errorf(codes.NotFound, "车次不存在")
	}
	var dailytrain model.DailyTrain
	tp, _ := time.Parse("2006-01-02", date)
	dailytrain.Code = train.Code
	dailytrain.Type = train.Type
	dailytrain.Start = train.Start
	dailytrain.End = train.End
	dailytrain.StartTime = train.StartTime
	dailytrain.EndTime = train.EndTime
	dailytrain.StartPinyin = train.StartPinyin
	dailytrain.EndPinyin = train.EndPinyin
	dailytrain.Date = tp
	dailytrain.CreateTime = time.Now()
	dailytrain.UpdateTime = time.Now()
	result := global.DB.Create(&dailytrain)
	if result.Error != nil {
		zap.S().Errorf("车次信息生成失败: %s", result.Error.Error())
		return result.Error
	}
	zap.S().Infof("车次信息生成成功 date: %s, train_code: %s", date, trainCode)
	// 生成车站数据
	ss := &StationService{}
	_, err := ss.GenerateStationDaily(&model.DailyTrainStation{TrainCode: trainCode, Date: tp})
	if err != nil {
		zap.S().Errorf("车站信息生成失败: %s", err.Error())
		return err
	}
	// 生成车厢数据
	cs := &CarriageService{}
	err = cs.GenerateDailyCarriage(trainCode, date)
	if err != nil {
		zap.S().Errorf("车厢信息生成失败: %s", err.Error())
		return err
	}
	// 生成座位数据
	sss := &SeatService{}
	err = sss.GenerateDailySeat(trainCode, date)
	if err != nil {
		zap.S().Errorf("座位信息生成失败: %s", err.Error())
		return err
	}
	// 生成车票数据
	ts := &TicketService{}
	dtt, err := ts.GenerateDailyTicket(trainCode, date)
	if err != nil {
		zap.S().Errorf("车票信息生成失败: %s", err.Error())
		return err
	}
	zap.S().Infof("车票信息生成成功 date: %s, train_code: %s, count: %d", date, trainCode, len(dtt))
	return nil
}

// 生成某日所有车次的数据
func (t *TrainService) GenDaily(date string) error {
	var trains []model.Train
	result := global.DB.Find(&trains)
	if result.Error != nil {
		return result.Error
	}
	for _, train := range trains {
		err := t.GenerateTrainDaily(train.Code, date)
		if err != nil {
			zap.S().Errorf("车次信息生成失败: %s", err.Error())
			return err
		}
	}
	return nil
}
