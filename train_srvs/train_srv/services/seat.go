package services

import (
	"fmt"
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SeatService struct{}

func (s *SeatService) CountSeat(date string, trainCode string, typ global.SeatType) int64 {
	var count int64
	dt, _ := time.Parse("2006-01-02", date)
	global.DB.Where(&model.DailyTrainSeat{Date: dt, TrainCode: trainCode, SeatType: string(typ)}).Count(&count)
	return count
}

// 根据车次生成座位信息
func (s *SeatService) GenTrainSeat(trainCode string) error {
	// 根据trainCode删除原先的座位信息
	res := global.DB.Delete(&model.DailyTrainSeat{}, "train_code = ?", trainCode)
	if res.Error != nil {
		zap.S().Errorw("删除座位信息失败", "msg", res.Error.Error())
		return res.Error
	}
	// 查找该车次的所有车厢
	var carriages []model.TrainCarriage
	result := global.DB.Where(&model.TrainCarriage{TrainCode: trainCode}).Find(&carriages)
	if result.Error != nil {
		zap.S().Errorw("查询车厢信息失败", "msg", result.Error.Error())
		return result.Error
	}

	// 循环生成每个车厢的座位信息
	for _, carriage := range carriages {
		var seats []model.TrainSeat
		var seatIndex int32 = 1
		cols := getColsBySeatType(global.SeatType(carriage.SeatType))
		for i := 1; i <= int(carriage.RowCount); i++ {
			for _, col := range cols {
				seat := model.TrainSeat{
					TrainCode:         trainCode,
					SeatType:          carriage.SeatType,
					CarriageIndex:     carriage.Index,
					Col:               col,
					CarriageSeatIndex: seatIndex,
					CreateTime:        time.Now(),
					UpdateTime:        time.Now(),
					Row_:              fmt.Sprintf("0%d", i),
				}
				seatIndex++
				seats = append(seats, seat)
			}
		}
		// 批量插入
		res := global.DB.CreateInBatches(seats, len(seats))
		if res.Error != nil {
			zap.S().Errorw("生成座位信息失败", "msg", res.Error.Error())
			return res.Error
		}
	}
	return nil
}

func (s *SeatService) GenerateDailySeat(trainCode, date string) error {
	// 删除某日某趟车次的所有座位信息
	result := global.DB.Delete(&model.DailyTrainSeat{}, "train_code = ? and date = ?", trainCode, date)
	if result.Error != nil {
		zap.S().Errorw("删除座位信息失败", "msg", result.Error.Error())
		return result.Error
	}
	// ss := &CarriageService{}
	// carriages, err := ss.GetCarriageByTrainCode(trainCode)
	// stations, err := ss.GetStationsByTrainCode(trainCode)
	// if err != nil {
	// 	zap.S().Errorw("查询车箱信息失败", "msg", err.Error())
	// 	return err
	// }
	// sell := fmt.Sprintf("%0*s", len(stations)-1, "")
	sell := int64(0)
	// 查询某车次的所有座位信息
	var seats []model.TrainSeat
	result = global.DB.Where(&model.TrainSeat{TrainCode: trainCode}).Find(&seats)
	if result.Error != nil {
		zap.S().Errorw("查询座位信息失败", "msg", result.Error.Error())
		return result.Error
	}
	if len(seats) == 0 {
		return status.Errorf(codes.NotFound, "座位信息不存在")
	}
	for _, seat := range seats {
		var seatDaily model.DailyTrainSeat
		// 多此一举
		tp, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		seatDaily.TrainCode = trainCode
		seatDaily.CreateTime = time.Now()
		seatDaily.UpdateTime = time.Now()
		seatDaily.Date = tp
		seatDaily.CarriageIndex = seat.CarriageIndex
		seatDaily.CarriageSeatIndex = seat.CarriageSeatIndex
		seatDaily.SeatType = seat.SeatType
		seatDaily.Row_ = seat.Row_
		seatDaily.Col = seat.Col
		seatDaily.Sell = sell
		result := global.DB.Create(&seatDaily)
		if result.Error != nil {
			zap.S().Errorw("生成座位信息失败", "msg", result.Error.Error())
			return result.Error
		}
	}
	zap.S().Infof("座位信息生成成功 date: %s, train_code: %s", date, trainCode)
	return nil
}

func getColsBySeatType(typ global.SeatType) []string {
	cols := []string{}
	switch typ {
	case global.SeatTypeFirstClass:
		cols = []string{"A", "C", "D", "F"}
	default:
		cols = []string{"A", "B", "C", "D", "F"}
	}
	return cols
}
