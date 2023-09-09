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

func (t *TicketServer) GenerateTicket(ctx context.Context, ticketreq *proto.TicketRequest) (*proto.TicketListResponse, error) {
	var tickets []model.DailyTrainTicket
	var ticketresp proto.TicketListResponse
	ss := &services.StationService{}
	dts, err := ss.GetStationsByTrainCodeDaily(ticketreq.TrainCode, ticketreq.Date)
	ts := &services.SeatService{}
	tts := &services.TrainService{}
	t2, err2 := tts.GetTrainByCode(ticketreq.TrainCode)
	dt, _ := time.Parse("2006-01-02", ticketreq.Date)
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
			fc_count := ts.CountSeat(ticketreq.Date, ticketreq.TrainCode, global.SeatTypeFirstClass)

			sc_count := ts.CountSeat(ticketreq.Date, ticketreq.TrainCode, global.SeatTypeSecondClass)
			sf_count := ts.CountSeat(ticketreq.Date, ticketreq.TrainCode, global.SeatTypeSoftBerth)
			hb_count := ts.CountSeat(ticketreq.Date, ticketreq.TrainCode, global.SeatTypeHardBerth)

			fcp := kms * float64(global.SeatFarePerKmFirstClass) * float64(tts.GetRate(global.TrainType(t2.Type)))
			scp := kms * float64(global.SeatFarePerKmSecondClass) * float64(tts.GetRate(global.TrainType(t2.Type)))
			sfp := kms * float64(global.SeatFarePerKmSoftberth) * float64(tts.GetRate(global.TrainType(t2.Type)))
			hbp := kms * float64(global.SeatFarePerKmHardberth) * float64(tts.GetRate(global.TrainType(t2.Type)))

			ticket := model.DailyTrainTicket{
				TrainCode:        ticketreq.TrainCode,
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
			}
			tickets = append(tickets, ticket)
			ticketrsp := TicketModel2Response(ticket)
			ticketresp.Data = append(ticketresp.Data, ticketrsp)
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
	return &ticketresp, nil
}
