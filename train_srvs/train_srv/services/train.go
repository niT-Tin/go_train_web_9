package services

import (
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
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
