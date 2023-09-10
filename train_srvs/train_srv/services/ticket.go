package services

import (
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TicketService struct{}

func (t *TicketService) GenerateDailyTicket(trainCode, date string) ([]model.DailyTrainTicket, error) {

	var tickets []model.DailyTrainTicket
	ss := &StationService{}
	dts, err := ss.GetStationsByTrainCodeDaily(trainCode, date)
	ts := &SeatService{}
	tts := &TrainService{}
	t2, err2 := tts.GetTrainByCode(trainCode)
	dt, _ := time.Parse("2006-01-02", date)
	if err2 != nil {
		return nil, err2
	}
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(dts); i++ {
		var kms float64
		for j := i + 1; j < len(dts); j++ {
			kms = dts[j].Km + dts[i].Km
			fc_count := ts.CountSeat(date, trainCode, global.SeatTypeFirstClass)

			sc_count := ts.CountSeat(date, trainCode, global.SeatTypeSecondClass)
			sf_count := ts.CountSeat(date, trainCode, global.SeatTypeSoftBerth)
			hb_count := ts.CountSeat(date, trainCode, global.SeatTypeHardBerth)

			fcp := kms * float64(global.SeatFarePerKmFirstClass) * float64(tts.GetRate(global.TrainType(t2.Type)))
			scp := kms * float64(global.SeatFarePerKmSecondClass) * float64(tts.GetRate(global.TrainType(t2.Type)))
			sfp := kms * float64(global.SeatFarePerKmSoftberth) * float64(tts.GetRate(global.TrainType(t2.Type)))
			hbp := kms * float64(global.SeatFarePerKmHardberth) * float64(tts.GetRate(global.TrainType(t2.Type)))

			ticket := model.DailyTrainTicket{
				TrainCode:        trainCode,
				Start:            dts[i].Name,
				End:              dts[j].Name,
				StartPinyin:      dts[i].NamePinyin,
				EndPinyin:        dts[j].NamePinyin,
				StartTime:        dts[i].OutTime,
				EndTime:          dts[j].InTime,
				Date:             dt,
				FirstClassLast:   int32(fc_count),
				SecondClassLast:  int32(sc_count),
				SoftBerthLast:    int32(sf_count),
				HardBerthLast:    int32(hb_count),
				FirstClassPrice:  fcp,
				SecondClassPrice: scp,
				SoftBerthPrice:   sfp,
				HardBerthPrice:   hbp,
				StartIndex:       dts[i].Index,
				EndIndex:         dts[j].Index,
			}
			tickets = append(tickets, ticket)
		}
	}
	result := global.DB.CreateInBatches(tickets, len(tickets))
	if result.Error != nil {
		zap.S().Errorf("生成车票失败")
		return nil, result.Error
	}
	if result.RowsAffected == 0 && len(tickets) > 0 {
		zap.S().Errorf("生成车票失败")
		return nil, status.Errorf(codes.Internal, "生成车票失败")
	}
	zap.S().Infof("生成车票成功 车票数量: %d", result.RowsAffected)
	return tickets, nil
}

// TODO: 记得添加redis分布式锁
func (t *TicketService) DeducInventory(date, start_s, end_s string, start_time string, train_seat []*model.TrainSeat) error {
	dt, _ := time.Parse("2006-01-02", date)
	st, _ := time.Parse("2006-01-02 15:04:05", start_time)
	// 根据日期、车次、起始站、终点站、发车时间查询车票信息(两车站之间的站数，用于计算票价, 填充sell字段)
	var dailyTrainTicket model.DailyTrainTicket
	var dailyTrainSeat model.DailyTrainSeat
	tx := global.DB.Begin()
	for _, seat := range train_seat {
		res := global.DB.Where(&model.DailyTrainTicket{Date: dt, TrainCode: seat.TrainCode, Start: start_s, End: end_s, StartTime: st}).First(&dailyTrainTicket)
		if res.Error != nil {
			tx.Rollback()
			zap.S().Errorf("查询车票信息失败")
			return res.Error
		}
		// 查询对应的位置情况
		d := global.DB.Where(&model.DailyTrainSeat{Date: dt, TrainCode: seat.TrainCode, SeatType: seat.SeatType, CarriageIndex: seat.CarriageIndex, Col: seat.Col, Row_: seat.Row_}).First(&dailyTrainSeat)
		if d.Error != nil {
			tx.Rollback()
			return d.Error
		}
		// 最低位表示火车开始站，最高位表示火车终点站(这一站到下一站之间)
		// 从0开始
		start := dailyTrainSeat.Sell >> dailyTrainTicket.StartIndex
		mask_bit_count := dailyTrainTicket.EndIndex - dailyTrainTicket.StartIndex
		// 生成最后为数为mask_bit_count位并且每一位都为1的掩码
		mask := int64((1 << mask_bit_count) - 1)
		// 获取该位置相关站的信息
		seats := start & mask
		if seats == 0 {
			// 表示位置没有被订购可以出售
			// 更新位置信息，将Sell的startIndex到endIndex之间的位置置为1
			dailyTrainSeat.Sell = dailyTrainSeat.Sell | (mask << dailyTrainTicket.StartIndex)
			// 更新售卖情况
			result := tx.Save(&dailyTrainSeat)
			typ := global.SeatType(dailyTrainSeat.SeatType)
			switch typ {
			case global.SeatTypeFirstClass:
				dailyTrainTicket.FirstClassLast--
			case global.SeatTypeSecondClass:
				dailyTrainTicket.SecondClassLast--
			case global.SeatTypeSoftBerth:
				dailyTrainTicket.SoftBerthLast--
			case global.SeatTypeHardBerth:
				dailyTrainTicket.HardBerthLast--
			}
			if result.Error != nil {
				return result.Error
			}
			// 更新车票信息
			result = tx.Save(&dailyTrainTicket)
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
			tx.Commit()
		}
		return status.Errorf(codes.Internal, "位置已经被订购，不能出售")
		// 位置已经被订购，不能出售
	}
	return nil
}

// TODO: 记得添加redis分布式锁
func (t *TicketService) ReBackInventory(orderId, date, start_s, end_s string, start_time string, train_seat []*model.TrainSeat) error {
	dt, _ := time.Parse("2006-01-02", date)
	st, _ := time.Parse("2006-01-02 15:04:05", start_time)
	tx := global.DB.Begin()
	// 根据日期、车次、起始站、终点站、发车时间查询车票信息(两车站之间的站数，用于计算票价, 填充sell字段)
	var dailyTrainTicket model.DailyTrainTicket
	var dailyTrainSeat model.DailyTrainSeat
	for _, seat := range train_seat {
		res := global.DB.Where(&model.DailyTrainTicket{Date: dt, TrainCode: seat.TrainCode, Start: start_s, End: end_s, StartTime: st}).First(&dailyTrainTicket)
		if res.Error != nil {
			tx.Rollback()
			zap.S().Errorf("查询车票信息失败")
			return res.Error
		}
		// 查询对应的位置情况
		d := global.DB.Where(&model.DailyTrainSeat{Date: dt, TrainCode: seat.TrainCode, SeatType: seat.SeatType, CarriageIndex: seat.CarriageIndex, Col: seat.Col, Row_: seat.Row_}).First(&dailyTrainSeat)
		if d.Error != nil {
			tx.Rollback()
			return d.Error
		}
		start := dailyTrainSeat.Sell >> dailyTrainTicket.StartIndex
		mask_bit_count := dailyTrainTicket.EndIndex - dailyTrainTicket.StartIndex
		mask := int64((1 << mask_bit_count) - 1)
		// 如果mask位所有位都为1,则将其对应位(也就是已经出售的位)设置为0
		if start&mask == mask {
			// 将Sell的startIndex到endIndex之间的位置置为0
			dailyTrainSeat.Sell = dailyTrainSeat.Sell & (^(mask << dailyTrainTicket.StartIndex))
			result := tx.Save(&dailyTrainSeat)
			typ := global.SeatType(dailyTrainSeat.SeatType)
			switch typ {
			case global.SeatTypeFirstClass:
				dailyTrainTicket.FirstClassLast++
			case global.SeatTypeSecondClass:
				dailyTrainTicket.SecondClassLast++
			case global.SeatTypeSoftBerth:
				dailyTrainTicket.SoftBerthLast++
			case global.SeatTypeHardBerth:
				dailyTrainTicket.HardBerthLast++
			}
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
			// 更新车票信息
			result = tx.Save(&dailyTrainTicket)
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
			tx.Commit()
		}
		return status.Errorf(codes.Internal, "归还错误")
	}
	return nil
}
