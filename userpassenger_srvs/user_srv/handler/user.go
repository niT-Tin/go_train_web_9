package handler

import (
	"context"
	"gotrains/userpassenger_srvs/user_srv/global"
	"gotrains/userpassenger_srvs/user_srv/model"
	"gotrains/userpassenger_srvs/user_srv/proto"
	"gotrains/userpassenger_srvs/user_srv/utils"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Model2Response(user model.User) *proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.Nickname,
		Gender:   user.Gender,
		Mobile:   user.Mobile,
		Role:     uint32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return &userInfoRsp
}

func (u *UserServer) GetUserList(ctx context.Context, pageInfo *proto.PageInfo) (*proto.UserListResponse, error) {
	zap.S().Infof("获取用户列表")
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = uint32(result.RowsAffected)

	global.DB.Scopes(utils.Paginate(int(pageInfo.Pn), int(pageInfo.Ps))).Find(&users)
	for _, user := range users {
		userInfoRsp := Model2Response(user)
		rsp.Data = append(rsp.Data, userInfoRsp)
	}
	return rsp, nil
}

func (u *UserServer) GetUserByMobile(ctx context.Context, mobileRequest *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: mobileRequest.Mobile}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	userInfoRsp := Model2Response(user)
	return userInfoRsp, nil
}

func (u *UserServer) GetUserById(ctx context.Context, idRequest *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, idRequest.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	userInfoRsp := Model2Response(user)
	return userInfoRsp, nil
}

func (u *UserServer) CreateUser(ctx context.Context, createUserInfo *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: createUserInfo.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user = model.User{
		Mobile:   createUserInfo.Mobile,
		Nickname: createUserInfo.NickName,
	}
	user.Password = utils.Encrypt(createUserInfo.Password)
	result = global.DB.Create(&user)
	if result.Error != nil {
		// 这里注意错误是否为内部错误
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoRsp := Model2Response(user)
	return userInfoRsp, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, updateUserInfo *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, updateUserInfo.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	bt := time.Unix(int64(updateUserInfo.BirthDay), 0)
	user = model.User{
		// BaseModel: model.BaseModel{ID: updateUserInfo.Id},
		Nickname: updateUserInfo.NickName,
		Gender:   updateUserInfo.Gender,
		Birthday: &bt,
	}
	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u *UserServer) CheckPassWord(ctx context.Context, passwordCheckInfo *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	_, salt, encodedPwd := utils.GetEncoded(passwordCheckInfo.EncryptedPassword)
	result := utils.Verify(passwordCheckInfo.Password, salt, encodedPwd)
	return &proto.CheckResponse{Success: result}, nil
}

func (u *UserServer) AddPassenger(ctx context.Context, passengerInfo *proto.PassengerInfo) (*proto.PassengerInfo, error) {
	result := global.DB.First(&model.User{}, passengerInfo.UserId)
	if result.Error != nil {
		return nil, result.Error
	}
	passenger := model.Passenger{
		Name:   passengerInfo.Name,
		UserID: passengerInfo.UserId,
		IdCard: passengerInfo.IdCard,
		Type:   model.PassengerType(passengerInfo.Type),
	}
	result = global.DB.Create(&passenger)
	if result.Error != nil {
		return nil, result.Error
	}
	return &proto.PassengerInfo{
		Id:     passenger.ID,
		Name:   passenger.Name,
		UserId: passenger.UserID,
		IdCard: passenger.IdCard,
		Type:   int64(passenger.Type),
	}, nil
}

func (u *UserServer) UpdatePassenger(ctx context.Context, passengerInfo *proto.PassengerInfo) (*emptypb.Empty, error) {
	var passenger model.Passenger
	result := global.DB.First(&passenger, passengerInfo.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "乘客不存在")
	}
	result = global.DB.First(&model.User{}, passengerInfo.UserId)
	if result.RowsAffected == 0 {
		passenger.IsDelete = true
		global.DB.Save(passenger)
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	passenger = model.Passenger{
		Name:   passengerInfo.Name,
		UserID: passengerInfo.UserId,
		IdCard: passengerInfo.IdCard,
		Type:   model.PassengerType(passengerInfo.Type),
	}
	result = global.DB.Save(passenger)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u *UserServer) DeletePassenger(ctx context.Context, idRequest *proto.IdRequest) (*emptypb.Empty, error) {
	var passenger model.Passenger
	result := global.DB.First(&passenger, idRequest.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "乘客不存在")
	}
	passenger.IsDelete = true
	result = global.DB.Save(passenger)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u *UserServer) GetPassengerList(ctx context.Context, passengerPageInfo *proto.PassengerPageInfo) (*proto.PassengerListResponse, error) {
	var passengerList []model.Passenger
	user := global.DB.First(&model.User{}, passengerPageInfo.UserId)
	if user.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	rsp := &proto.PassengerListResponse{}
	where := global.DB.Where(&model.Passenger{UserID: passengerPageInfo.UserId, BaseModel: model.BaseModel{IsDelete: false}})
	result := where.Scopes(utils.Paginate(int(passengerPageInfo.Pn), int(passengerPageInfo.Ps))).Find(&passengerList)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	rsp.Total = uint32(result.RowsAffected)
	for _, passenger := range passengerList {
		if passenger.IsDelete {
			continue
		}
		rsp.Data = append(rsp.Data, &proto.PassengerInfo{
			Id:     passenger.ID,
			Name:   passenger.Name,
			UserId: passenger.UserID,
			IdCard: passenger.IdCard,
			Type:   int64(passenger.Type),
		})
	}
	return rsp, nil
}
