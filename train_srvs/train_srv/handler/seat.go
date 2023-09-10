package handler

import (
	"context"
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"gotrains/train_srvs/train_srv/proto"
	"gotrains/train_srvs/train_srv/services"
	"gotrains/train_srvs/train_srv/utils"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SeatServer) GetSeatDailyByTrainCode(ctx context.Context, dailyreq *proto.SeatDailyRequest) (*proto.SeatDailyResponse, error) {
	var seatDaily model.DailyTrainSeat
	bt, _ := time.Parse("2006-01-02 15:04:05", dailyreq.Date)
	result := global.DB.Where(&model.DailyTrainSeat{
		TrainCode: dailyreq.TrainCode,
		// CarriageIndex: dailyreq.CarriageIndex,
		Date: bt,
	}).First(&seatDaily, dailyreq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "座位不存在")
	}
	seatDailyRsp := DailySeatModel2Response(seatDaily)
	return seatDailyRsp, nil
}

// Seat Server Implement
// func (s *SeatServer) GetSeatDailyList(_ context.Context, _ *proto.SeatDailyPageInfo) (*proto.SeatDailyListResponse, error) {
// 	panic("not implemented")
// }

// 限定为traincode和date以及carriageindex
func (s *SeatServer) GetSeatList(ctx context.Context, pageinfo *proto.SeatPageInfo) (*proto.SeatListResponse, error) {
	var seatList []model.TrainSeat
	bt, _ := time.Parse("2006-01-02", pageinfo.Date)
	result := global.DB.Where(&model.DailyTrainSeat{
		TrainCode:     pageinfo.Seat.TrainCode,
		CarriageIndex: pageinfo.Seat.CarriageIndex,
		Date:          bt,
	}).
		Find(&seatList)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.SeatListResponse{}
	rsp.Total = uint32(result.RowsAffected)

	global.DB.Scopes(utils.Paginate(int(pageinfo.Pn), int(pageinfo.Ps))).Find(&seatList)
	for _, seat := range seatList {
		seatInfoRsp := SeatModel2Response(seat)
		rsp.Data = append(rsp.Data, seatInfoRsp)
	}
	return rsp, nil
}

// 通过id获取座位信息
func (s *SeatServer) GetSeat(ctx context.Context, seatreq *proto.SeatRequest) (*proto.SeatResponse, error) {
	var seat model.TrainSeat
	result := global.DB.First(&seat, seatreq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "座位不存在")
	}
	seatInfoRsp := SeatModel2Response(seat)
	return seatInfoRsp, nil
}

// func (s *SeatServer) CreateSeat(_ context.Context, _ *proto.SeatRequest) (*proto.SeatResponse, error) {
// }

func (s *SeatServer) UpdateSeat(ctx context.Context, seatreq *proto.SeatRequest) (*proto.SeatResponse, error) {
	var seat model.TrainSeat
	result := global.DB.First(&seat, seatreq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "座位不存在")
	}
	seat = model.TrainSeat{
		TrainCode:         seatreq.TrainCode,
		CarriageIndex:     seatreq.CarriageIndex,
		SeatType:          seatreq.SeatType,
		Row_:              seatreq.Row,
		Col:               seatreq.Column,
		CarriageSeatIndex: seatreq.SeatIndex,
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
	}
	result = global.DB.Save(&seat)
	if result.Error != nil {
		return nil, result.Error
	}
	seatInfoRsp := SeatModel2Response(seat)
	return seatInfoRsp, nil
}

func (s *SeatServer) DeleteSeat(ctx context.Context, seatreq *proto.SeatRequest) (*proto.SeatResponse, error) {
	result := global.DB.Delete(&model.TrainSeat{}, seatreq.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "删除座位失败")
	}
	return &proto.SeatResponse{}, nil
}

// func (s *SeatServer) GetSeatDailyListByDate(_ context.Context, _ *proto.SeatDailyPageInfo) (*proto.SeatDailyListResponse, error) {
// }

// func (s *SeatServer) GetAllSeat(_ context.Context, _ *proto.SeatRequest) (*proto.SeatListResponse, error) {

// }

func (s *SeatServer) GenerateSeatDaily(ctx context.Context, seatreq *proto.SeatDailyRequest) (*proto.SeatListResponse, error) {
	ss := &services.SeatService{}
	err := ss.GenerateDailySeat(seatreq.TrainCode, seatreq.Date)
	if err != nil {
		zap.S().Errorw("生成座位信息失败", "msg", err.Error())
		return nil, err
	}
	return &proto.SeatListResponse{}, nil
}
