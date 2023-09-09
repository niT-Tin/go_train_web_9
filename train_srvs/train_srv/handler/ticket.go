package handler

import (
	"context"
	"gotrains/train_srvs/train_srv/proto"
)

// Ticket Server Implement
func (t *TicketServer) GetTicketList(_ context.Context, _ *proto.TicketPageInfo) (*proto.TicketListResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (t *TicketServer) GetTicket(_ context.Context, _ *proto.TicketRequest) (*proto.TicketResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (t *TicketServer) CreateTicket(_ context.Context, _ *proto.TicketRequest) (*proto.TicketResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (t *TicketServer) UpdateTicket(_ context.Context, _ *proto.TicketRequest) (*proto.TicketResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (t *TicketServer) DeleteTicket(_ context.Context, _ *proto.TicketRequest) (*proto.TicketResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (t *TicketServer) GetTicketListByDate(_ context.Context, _ *proto.TicketPageInfo) (*proto.TicketListResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (t *TicketServer) GetAllTicket(_ context.Context, _ *proto.TicketRequest) (*proto.TicketListResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (t *TicketServer) GenerateTicket(_ context.Context, _ *proto.TicketRequest) (*proto.TicketListResponse, error) {
	panic("not implemented") // TODO: Implement
}

// Seat Server Implement
// TODO: 之后删除
func (s *SeatServer) GetSeatDailyList(_ context.Context, _ *proto.SeatDailyPageInfo) (*proto.SeatDailyListResponse, error) {
	panic("not implemented") // TODO: Implement
}
