package handler

import (
	"context"
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/model"
	"gotrains/train_srvs/train_srv/proto"
	"gotrains/train_srvs/train_srv/utils"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Station Server Implement
// 多于的方法之后删除
func (s *StationServer) GetStationDailyList(ctx context.Context, dailypage *proto.StationDailyPageInfo) (*proto.StationDailyListResponse, error) {
	return nil, nil
}

func (s *StationServer) GetStationDaily(ctx context.Context, dailyreq *proto.StationDailyRequest) (*proto.StationDailyResponse, error) {
	var station model.DailyTrainStation
	result := global.DB.First(&station, dailyreq.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车站不存在")
	}
	stationdaily := DailyStationModel2Response(station)
	return stationdaily, nil
}

func (s *StationServer) GetStationList(ctx context.Context, stationpage *proto.StationPageInfo) (*proto.StationListResponse, error) {
	zap.S().Infof("获取车站列表")
	var stations []model.Station
	result := global.DB.Find(&stations)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.StationListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	global.DB.Scopes(utils.Paginate(int(stationpage.Pn), int(stationpage.Ps))).Find(&stations)
	for _, station := range stations {
		stationinfo := StationModel2Response(station)
		rsp.Data = append(rsp.Data, stationinfo)
	}
	return rsp, nil
}

func (s *StationServer) GetStation(ctx context.Context, stationreq *proto.StationRequest) (*proto.StationResponse, error) {
	var station model.Station
	result := global.DB.First(&station, stationreq.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车站不存在")
	}
	stationinfo := StationModel2Response(station)
	return stationinfo, nil
}

func (s *StationServer) CreateStation(ctx context.Context, stationreq *proto.StationRequest) (*proto.StationResponse, error) {
	var station model.Station
	result := global.DB.Where(&model.Station{Name: stationreq.Name}).First(&station)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "车站已存在")
	}
	station = model.Station{
		Name:       stationreq.Name,
		NamePinyin: stationreq.Pinyin,
		NamePy:     stationreq.FirstLetter,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	result = global.DB.Create(&station)
	if result.Error != nil {
		// 这里注意错误是否为内部错误
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	stationinfo := StationModel2Response(station)
	return stationinfo, nil
}

func (s *StationServer) UpdateStation(ctx context.Context, stationreq *proto.StationRequest) (*proto.StationResponse, error) {
	var station model.Station
	result := global.DB.First(&station, stationreq.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "车站不存在")
	}
	station = model.Station{
		Name:       stationreq.Name,
		NamePinyin: stationreq.Pinyin,
		NamePy:     stationreq.FirstLetter,
		UpdateTime: time.Now(),
	}
	result = global.DB.Save(station)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	stationinfo := StationModel2Response(station)
	return stationinfo, nil
}

// Should not be invoked
func (s *StationServer) DeleteStation(ctx context.Context, stationreq *proto.StationRequest) (*proto.StationResponse, error) {
	return nil, nil
}

func (s *StationServer) GetStationDailyListByDate(ctx context.Context, dailypage *proto.StationDailyPageInfo) (*proto.StationDailyListResponse, error) {
	var stations []model.DailyTrainStation
	bt, _ := time.Parse("2006-01-02", dailypage.Date)
	result := global.DB.Where(&model.DailyTrainStation{Date: bt}).Find(&stations)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.StationDailyListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	global.DB.Scopes(utils.Paginate(int(dailypage.Pn), int(dailypage.Ps))).Find(&stations)
	for _, station := range stations {
		stationinfo := DailyStationModel2Response(station)
		rsp.Data = append(rsp.Data, stationinfo)
	}
	return rsp, nil
}

func (s *StationServer) GetAllStation(ctx context.Context, stationreq *proto.StationRequest) (*proto.StationListResponse, error) {
	var stations []model.Station
	result := global.DB.Find(&stations)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.StationListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	for _, station := range stations {
		stationinfo := StationModel2Response(station)
		rsp.Data = append(rsp.Data, stationinfo)
	}
	return rsp, nil
}

func (s *StationServer) GenerateStationDaily(_ context.Context, _ *proto.StationRequest) (*proto.StationListResponse, error) {
	panic("not implemented") // TODO: Implement
}
