package handler

import (
	"gotrains/train_srvs/train_srv/model"
	"gotrains/train_srvs/train_srv/proto"
)

type TrainServer struct {
	proto.UnimplementedTrainServer
}

type StationServer struct {
	proto.UnimplementedStationServer
}

type TicketServer struct {
	proto.UnimplementedTicketServer
}

type CarriageServer struct {
	proto.UnimplementedCarriageServer
}

type SeatServer struct {
	proto.UnimplementedSeatServer
}

func StationModel2Response(station model.Station) *proto.StationResponse {
	stationInfoRsp := proto.StationResponse{
		Data: &proto.StationInfo{
			Id:          station.ID,
			Name:        station.Name,
			Pinyin:      station.NamePinyin,
			FirstLetter: station.NamePy,
		},
	}
	return &stationInfoRsp
}

func DailyStationModel2Response(station model.DailyTrainStation) *proto.StationDailyResponse {
	stationInfoRsp := proto.StationDailyResponse{
		Data: &proto.StationDailyInfo{
			Date: station.Date.String(),
			Station: &proto.StationInfo{
				Id:          station.ID,
				Name:        station.Name,
				Pinyin:      station.NamePinyin,
				FirstLetter: station.NamePinyin,
			},
		},
	}
	return &stationInfoRsp
}

func TrainModel2Response(train model.Train) *proto.TrainResponse {
	trainInfoRsp := proto.TrainResponse{
		Data: &proto.TrainInfo{
			Id:           train.ID,
			Code:         train.Code,
			Type:         train.Type,
			StartStation: train.Start,
			EndStation:   train.End,
			StartTime:    train.StartTime.Unix(),
			EndTime:      train.EndTime.Unix(),
			StartPinyin:  train.StartPinyin,
			EndPinyin:    train.EndPinyin,
		},
	}
	return &trainInfoRsp
}

func DailyTrainModel2Response(train model.DailyTrain) *proto.TrainDailyResponse {
	trainInfoRsp := proto.TrainDailyResponse{
		Data: &proto.TrainDailyInfo{
			Date: train.Date.String(),
			Train: &proto.TrainInfo{
				Id:           train.ID,
				Code:         train.Code,
				Type:         train.Type,
				StartStation: train.Start,
				EndStation:   train.End,
				StartTime:    train.StartTime.Unix(),
				EndTime:      train.EndTime.Unix(),
				StartPinyin:  train.StartPinyin,
				EndPinyin:    train.EndPinyin,
			},
		},
	}
	return &trainInfoRsp
}

func CarriageModel2Response(carriage model.TrainCarriage) *proto.CarriageResponse {
	carriageInfoRsp := proto.CarriageResponse{
		Data: &proto.CarriageInfo{
			Id:            carriage.ID,
			TrainCode:     carriage.TrainCode,
			CarriageIndex: int32(carriage.Index),
			SeatType:      carriage.SeatType,
			SeatCount:     int32(carriage.SeatCount),
			Row:           int32(carriage.RowCount),
			Column:        int32(carriage.ColCount),
			// CreateTime:    carriage.CreatedAt.Unix(),
		},
	}
	return &carriageInfoRsp
}

func DailyCarriageModel2Response(carriage model.DailyTrainCarriage) *proto.CarriageDailyResponse {
	carriageInfoRsp := proto.CarriageDailyResponse{
		Data: &proto.CarriageDailyInfo{
			Date: carriage.Date.String(),
			Carriage: &proto.CarriageInfo{
				Id:            carriage.ID,
				TrainCode:     carriage.TrainCode,
				CarriageIndex: int32(carriage.Index),
				SeatType:      carriage.SeatType,
				SeatCount:     int32(carriage.SeatCount),
				Row:           int32(carriage.RowCount),
				Column:        int32(carriage.ColCount),
				// CreateTime:    carriage.CreatedAt.Unix(),
			},
		},
	}
	return &carriageInfoRsp
}

func SeatModel2Response(seat model.TrainSeat) *proto.SeatResponse {
	seatInfoRsp := proto.SeatResponse{
		Data: &proto.SeatInfo{
			Id:            seat.ID,
			TrainCode:     seat.TrainCode,
			CarriageIndex: int32(seat.CarriageIndex),
			SeatIndex:     int32(seat.CarriageSeatIndex),
			SeatType:      seat.SeatType,
			Row:           seat.Row_,
			Column:        seat.Col,
			// CreateTime:    seat.CreatedAt.Unix(),
		},
	}
	return &seatInfoRsp
}

func DailySeatModel2Response(seat model.DailyTrainSeat) *proto.SeatDailyResponse {
	seatInfoRsp := proto.SeatDailyResponse{
		Data: &proto.SeatDailyInfo{
			Date: seat.Date.String(),
			Seat: &proto.SeatInfo{
				Id:            seat.ID,
				TrainCode:     seat.TrainCode,
				CarriageIndex: int32(seat.CarriageIndex),
				SeatIndex:     int32(seat.CarriageSeatIndex),
				SeatType:      seat.SeatType,
				Row:           seat.Row_,
				Column:        seat.Col,
				// CreateTime:    seat.CreatedAt.Unix(),
			},
		},
	}
	return &seatInfoRsp
}

func TicketModel2Response(ticket model.DailyTrainTicket) *proto.TicketResponse {
	ticketInfoRsp := proto.TicketResponse{
		Data: &proto.TicketInfo{
			Id:               ticket.ID,
			Date:             ticket.Date.Format("2006-01-02"),
			TrainCode:        ticket.TrainCode,
			StartStation:     ticket.Start,
			StartPinyin:      ticket.StartPinyin,
			EndStation:       ticket.End,
			EndPinyin:        ticket.EndPinyin,
			StartTime:        ticket.StartTime.Unix(),
			StartIndex:       int64(ticket.StartIndex),
			EndIndex:         int64(ticket.EndIndex),
			EndTime:          ticket.EndTime.Unix(),
			FirstClassLast:   int64(ticket.FirstClassLast),
			FirstClassPrice:  float32(ticket.FirstClassPrice),
			SecondClassLast:  int64(ticket.SecondClassLast),
			SecondClassPrice: float32(ticket.SecondClassPrice),
			SoftberthLast:    int64(ticket.SoftBerthLast),
			SoftberthPrice:   float32(ticket.SoftBerthPrice),
			HardberthLast:    int64(ticket.HardBerthLast),
			HardberthPrice:   float32(ticket.HardBerthPrice),
			// StartStation:  ticket.Start,
			// CreateTime:    ticket.CreatedAt.Unix(),
		},
	}
	return &ticketInfoRsp
}
