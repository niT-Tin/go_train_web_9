package handler

import (
	"context"
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"gotrains/train_srvs/train_srv/proto"
	"gotrains/train_srvs/train_srv/utils"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SeatServer) GetSeatDaily(ctx context.Context, dailyreq *proto.SeatDailyRequest) (*proto.SeatDailyResponse, error) {
	var seatDaily model.DailyTrainSeat
	bt, _ := time.Parse("2006-01-02 15:04:05", dailyreq.Date)
	result := global.DB.Where(&model.DailyTrainSeat{
		TrainCode:     dailyreq.TrainCode,
		CarriageIndex: dailyreq.CarriageIndex,
		Date:          bt,
	}).First(&seatDaily, dailyreq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "座位不存在")
	}
	seatDailyRsp := DailySeatModel2Response(seatDaily)
	return seatDailyRsp, nil
}

// Seat Server Implement
// TODO: 之后删除
func (s *SeatServer) GetSeatDailyList(_ context.Context, _ *proto.SeatDailyPageInfo) (*proto.SeatDailyListResponse, error) {
	panic("not implemented") // TODO: Implement
}

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

func (s *SeatServer) GetSeat(ctx context.Context, seatreq *proto.SeatRequest) (*proto.SeatResponse, error) {
	var seat model.TrainSeat
	result := global.DB.First(&seat, seatreq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "座位不存在")
	}
	seatInfoRsp := SeatModel2Response(seat)
	return seatInfoRsp, nil
}

func (s *SeatServer) CreateSeat(_ context.Context, _ *proto.SeatRequest) (*proto.SeatResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SeatServer) UpdateSeat(_ context.Context, _ *proto.SeatRequest) (*proto.SeatResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SeatServer) DeleteSeat(_ context.Context, _ *proto.SeatRequest) (*proto.SeatResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SeatServer) GetSeatDailyListByDate(_ context.Context, _ *proto.SeatDailyPageInfo) (*proto.SeatDailyListResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SeatServer) GetAllSeat(_ context.Context, _ *proto.SeatRequest) (*proto.SeatListResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SeatServer) GenerateSeatDaily(_ context.Context, _ *proto.SeatRequest) (*proto.SeatListResponse, error) {
	panic("not implemented") // TODO: Implement
}
