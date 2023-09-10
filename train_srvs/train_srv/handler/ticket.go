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

// Ticket Server Implement
// func (t *TicketServer) CreateTicket(ctx context.Context, ticketreq *proto.TicketRequest) (*proto.TicketResponse, error) {
// }

func (t *TicketServer) DeleteTicket(ctx context.Context, ticketreq *proto.TicketRequest) (*proto.TicketResponse, error) {
	result := global.DB.Delete(&model.DailyTrainTicket{}, "id = ?", ticketreq.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "删除车票失败")
	}
	return &proto.TicketResponse{}, nil
}

func (t *TicketServer) GetAllTicket(ctx context.Context, ticketreq *proto.TicketPageInfo) (*proto.TicketListResponse, error) {
	var tickets []model.DailyTrainTicket
	ticketresp := &proto.TicketListResponse{}
	result := global.DB.Scopes(utils.Paginate(int(ticketreq.Pn), int(ticketreq.Ps))).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, ticket := range tickets {
		ticketrsp := TicketModel2Response(ticket)
		ticketresp.Data = append(ticketresp.Data, ticketrsp)
	}
	return ticketresp, nil
}

func (t *TicketServer) GetTicketList(ctx context.Context, ticketreq *proto.TicketRequest) (*proto.TicketListResponse, error) {
	var tickets []model.DailyTrainTicket
	ticketresp := &proto.TicketListResponse{}
	dt, _ := time.Parse("2006-01-02", ticketreq.Date)
	result := global.DB.Where(&model.DailyTrainTicket{TrainCode: ticketreq.TrainCode, Date: dt}).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, ticket := range tickets {
		ticketrsp := TicketModel2Response(ticket)
		ticketresp.Data = append(ticketresp.Data, ticketrsp)
	}
	return ticketresp, nil
}

// 每日每趟车次的车票生成
func (t *TicketServer) GenerateTicket(ctx context.Context, ticketreq *proto.TicketRequest) (*proto.TicketListResponse, error) {
	ts := &services.TicketService{}
	dtt, err := ts.GenerateDailyTicket(ticketreq.TrainCode, ticketreq.Date)
	if err != nil {
		return nil, err
	}
	ticketresp := &proto.TicketListResponse{}
	for _, ticket := range dtt {
		ticketrsp := TicketModel2Response(ticket)
		ticketresp.Data = append(ticketresp.Data, ticketrsp)
	}
	return ticketresp, nil
}

// 扣减车票信息
func (t *TicketServer) ReductTicket(ctx context.Context, ticketreq *proto.BusinessRequest) (*proto.TicketResponse, error) {
	ts := &services.TicketService{}
	var tseats []*model.TrainSeat
	for _, seat := range ticketreq.Seats {
		tseat := &model.TrainSeat{
			TrainCode:         seat.TrainCode,
			CarriageIndex:     seat.CarriageIndex,
			SeatType:          seat.SeatType,
			Row_:              seat.Row,
			Col:               seat.Column,
			CarriageSeatIndex: seat.SeatIndex,
		}
		tseats = append(tseats, tseat)
	}
	err := ts.DeducInventory(ticketreq.Date, ticketreq.StartStation, ticketreq.EndStation, ticketreq.StartTime, tseats)
	if err != nil {
		zap.S().Errorf("扣减车票失败")
		return nil, err
	}
	return &proto.TicketResponse{}, nil
}

// 归还车票信息
func (t *TicketServer) RebackTicket(ctx context.Context, ticketreq *proto.BusinessRequest) (*proto.TicketResponse, error) {
	ts := &services.TicketService{}
	var tseats []*model.TrainSeat
	for _, seat := range ticketreq.Seats {
		tseat := &model.TrainSeat{
			TrainCode:         seat.TrainCode,
			CarriageIndex:     seat.CarriageIndex,
			SeatType:          seat.SeatType,
			Row_:              seat.Row,
			Col:               seat.Column,
			CarriageSeatIndex: seat.SeatIndex,
		}
		tseats = append(tseats, tseat)
	}
	err := ts.ReBackInventory(ticketreq.OrderId, ticketreq.Date, ticketreq.StartStation, ticketreq.EndStation, ticketreq.StartTime, tseats)
	if err != nil {
		zap.S().Errorf("归还车票失败")
		return nil, err
	}
	return &proto.TicketResponse{}, nil
}
