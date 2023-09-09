package handler

import (
	"context"
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"gotrains/train_srvs/train_srv/proto"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Carrier Server Implement
// TODO: 之后删除
func (c *CarriageServer) GetCarriageDailyList(_ context.Context, _ *proto.CarriageDailyPageInfo) (*proto.CarriageDailyListResponse, error) {
	return nil, nil
}

func (c *CarriageServer) GetCarriageDaily(ctx context.Context, dailyreq *proto.CarriageDailyRequest) (*proto.CarriageDailyResponse, error) {
	var carriageDaily model.DailyTrainCarriage
	result := global.DB.First(&carriageDaily, dailyreq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车厢不存在")
	}
	carriageDailyRsp := DailyCarriageModel2Response(carriageDaily)
	return carriageDailyRsp, nil
}

func (c *CarriageServer) GetCarriageList(ctx context.Context, carriagepage *proto.CarriagePageInfo) (*proto.CarriageListResponse, error) {
	var carriageList []model.TrainCarriage
	result := global.DB.Where(&model.TrainCarriage{TrainCode: carriagepage.TrainCode}).Find(&carriageList)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.CarriageListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	for _, carriage := range carriageList {
		carriageInfoRsp := CarriageModel2Response(carriage)
		rsp.Data = append(rsp.Data, carriageInfoRsp)
	}
	return rsp, nil
}

func (c *CarriageServer) GetCarriage(ctx context.Context, carriagereq *proto.CarriageRequest) (*proto.CarriageResponse, error) {
	var carriage model.TrainCarriage
	result := global.DB.Where(&model.TrainCarriage{ID: carriagereq.Id, TrainCode: carriage.TrainCode}).First(&carriage)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车厢不存在")
	}
	carriageInfoRsp := CarriageModel2Response(carriage)
	return carriageInfoRsp, nil
}

func (c *CarriageServer) CreateCarriage(ctx context.Context, carriagereq *proto.CarriageRequest) (*proto.CarriageResponse, error) {
	var train model.Train
	result := global.DB.Where(&model.Train{Code: carriagereq.TrainCode}).First(&train)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "列车不存在")
	}
	carriage := model.TrainCarriage{
		TrainCode:  carriagereq.TrainCode,
		Index:      carriagereq.CarriageIndex,
		SeatType:   carriagereq.SeatType,
		SeatCount:  carriagereq.SeatCount,
		RowCount:   carriagereq.Row,
		ColCount:   carriagereq.Column,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	result = global.DB.Create(&carriage)
	if result.Error != nil {
		return nil, result.Error
	}
	carriageInfoRsp := CarriageModel2Response(carriage)
	return carriageInfoRsp, nil
}

func (c *CarriageServer) UpdateCarriage(ctx context.Context, carriagereq *proto.CarriageRequest) (*proto.CarriageResponse, error) {
	var carriage model.TrainCarriage
	result := global.DB.First(&carriage, carriagereq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车厢不存在")
	}
	carriage = model.TrainCarriage{
		TrainCode:  carriagereq.TrainCode,
		Index:      carriagereq.CarriageIndex,
		SeatType:   carriagereq.SeatType,
		SeatCount:  carriagereq.SeatCount,
		RowCount:   carriagereq.Row,
		ColCount:   carriagereq.Column,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	result = global.DB.Save(&carriage)
	if result.Error != nil {
		return nil, result.Error
	}
	carriageInfoRsp := CarriageModel2Response(carriage)
	return carriageInfoRsp, nil
}

// should not be invoked
func (c *CarriageServer) DeleteCarriage(ctx context.Context, carriagereq *proto.CarriageRequest) (*proto.CarriageResponse, error) {
	return nil, nil
}

func (c *CarriageServer) GetCarriageDailyListByDate(ctx context.Context, carriagedaily *proto.CarriageDailyPageInfo) (*proto.CarriageDailyListResponse, error) {
	var carriageDailyList []model.DailyTrainCarriage
	bt, _ := time.Parse("2006-01-02", carriagedaily.Date)
	result := global.DB.Where(&model.DailyTrainCarriage{Date: bt}).Find(&carriageDailyList)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.CarriageDailyListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	for _, carriageDaily := range carriageDailyList {
		carriageDailyRsp := DailyCarriageModel2Response(carriageDaily)
		rsp.Data = append(rsp.Data, carriageDailyRsp)
	}
	return rsp, nil
}

func (c *CarriageServer) GetAllCarriage(ctx context.Context, carriagereq *proto.CarriageRequest) (*proto.CarriageListResponse, error) {
	var carriageList []model.TrainCarriage
	result := global.DB.Where(&model.TrainCarriage{TrainCode: carriagereq.TrainCode}).Find(&carriageList)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.CarriageListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	for _, carriage := range carriageList {
		carriageInfoRsp := CarriageModel2Response(carriage)
		rsp.Data = append(rsp.Data, carriageInfoRsp)
	}
	return rsp, nil
}

func (c *CarriageServer) GenerateCarriageDaily(ctx context.Context, carriagereq *proto.CarriageRequest) (*proto.CarriageListResponse, error) {
	panic("not implemented") // TODO: Implement
}
