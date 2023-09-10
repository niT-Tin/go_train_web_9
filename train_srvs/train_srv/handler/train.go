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

// Train Server Implement
// TODO: 多于的方法， 未来删除
func (t *TrainServer) GetTrainDailyList(ctx context.Context, dailyinfo *proto.TrainDailyPageInfo) (*proto.TrainDailyListResponse, error) {
	zap.S().Infof("获取日期车次列表")
	var trains []model.DailyTrain
	result := global.DB.Find(&trains)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.TrainDailyListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	// parse time from string
	bt, _ := time.Parse("2006-01-02", dailyinfo.Date)
	global.DB.Scopes(utils.Paginate(int(dailyinfo.Pn), int(dailyinfo.Ps))).Where(&model.DailyTrain{Date: bt}).Find(&trains)
	for _, train := range trains {
		traindaily := DailyTrainModel2Response(train)
		rsp.Data = append(rsp.Data, traindaily)
	}
	return rsp, nil
}

func (t *TrainServer) GetTrainDaily(ctx context.Context, dailyreq *proto.TrainDailyRequest) (*proto.TrainDailyResponse, error) {
	var train model.DailyTrain
	result := global.DB.First(&train, dailyreq.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车次不存在")
	}
	traindaily := DailyTrainModel2Response(train)
	return traindaily, nil
}

func (t *TrainServer) GetTrainList(ctx context.Context, trainpage *proto.TrainPageInfo) (*proto.TrainListResponse, error) {
	zap.S().Infof("获取车次列表")
	var trains []model.Train
	result := global.DB.Find(&trains)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.TrainListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	global.DB.Scopes(utils.Paginate(int(trainpage.Pn), int(trainpage.Ps))).Find(&trains)
	for _, train := range trains {
		traininfo := TrainModel2Response(train)
		rsp.Data = append(rsp.Data, traininfo)
	}
	return rsp, nil
}

func (t *TrainServer) GetTrain(ctx context.Context, trainreq *proto.TrainRequest) (*proto.TrainResponse, error) {
	var train model.Train
	result := global.DB.First(&train, trainreq.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车次不存在")
	}
	traininfo := TrainModel2Response(train)
	return traininfo, nil
}

func (t *TrainServer) CreateTrain(ctx context.Context, trainreq *proto.TrainRequest) (*proto.TrainResponse, error) {
	var train model.Train
	result := global.DB.Where(&model.Train{Code: trainreq.Code}).First(&train)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "车次已存在")
	}
	train = model.Train{
		Code:        trainreq.Code,
		Type:        trainreq.Type,
		Start:       trainreq.StartStation,
		End:         trainreq.EndStation,
		StartTime:   time.Unix(trainreq.StartTime, 0),
		EndTime:     time.Unix(trainreq.EndTime, 0),
		StartPinyin: trainreq.StartPinyin,
		EndPinyin:   trainreq.EndPinyin,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	result = global.DB.Create(&train)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	traininfo := TrainModel2Response(train)
	return traininfo, nil
}

func (t *TrainServer) UpdateTrain(ctx context.Context, trainreq *proto.TrainRequest) (*proto.TrainResponse, error) {
	var train model.Train
	result := global.DB.First(&train, trainreq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车次不存在")
	}
	train = model.Train{
		Code:        trainreq.Code,
		Type:        trainreq.Type,
		Start:       trainreq.StartStation,
		End:         trainreq.EndStation,
		StartTime:   time.Unix(trainreq.StartTime, 0),
		EndTime:     time.Unix(trainreq.EndTime, 0),
		StartPinyin: trainreq.StartPinyin,
		EndPinyin:   trainreq.EndPinyin,
		UpdateTime:  time.Now(),
	}
	result = global.DB.Save(train)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	traininfo := TrainModel2Response(train)
	return traininfo, nil
}

func (t *TrainServer) DeleteTrain(ctx context.Context, trainreq *proto.TrainRequest) (*proto.TrainResponse, error) {
	// should not be invoked
	return nil, nil
}

func (t *TrainServer) GetTrainDailyListByDate(ctx context.Context, trainpage *proto.TrainDailyPageInfo) (*proto.TrainDailyListResponse, error) {
	zap.S().Infof("获取指定日期车次列表")
	var trains []model.DailyTrain
	result := global.DB.Find(&trains)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.TrainDailyListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	// parse time from string
	bt, _ := time.Parse("2006-01-02 15:04:05", trainpage.Date)
	global.DB.Scopes(utils.Paginate(int(trainpage.Pn), int(trainpage.Ps))).Where(&model.DailyTrain{Date: bt}).Find(&trains)
	for _, train := range trains {
		traindaily := DailyTrainModel2Response(train)
		rsp.Data = append(rsp.Data, traindaily)
	}
	return rsp, nil
}

func (t *TrainServer) GetAllTrain(ctx context.Context, trainreq *proto.TrainRequest) (*proto.TrainListResponse, error) {
	zap.S().Infof("获取所有车次列表")
	var trains []model.Train
	result := global.DB.Find(&trains)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.TrainListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	global.DB.Find(&trains)
	for _, train := range trains {
		traininfo := TrainModel2Response(train)
		rsp.Data = append(rsp.Data, traininfo)
	}
	return rsp, nil
}

// TODO: 之后实现
func (t *TrainServer) GenerateTrainDaily(ctx context.Context, trainreq *proto.TrainDailyRequest) (*proto.TrainListResponse, error) {
	ts := &services.TrainService{}
	err := ts.GenerateTrainDaily(trainreq.Code, trainreq.Date)
	if err != nil {
		return nil, err
	}
	return &proto.TrainListResponse{}, nil
}
