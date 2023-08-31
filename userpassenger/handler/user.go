package handler

import (
	"context"
	"gotrains/userpassenger/global"
	"gotrains/userpassenger/model"
	"gotrains/userpassenger/proto"
	"gotrains/userpassenger/utils"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	proto.UnimplementedUserServer
}

func userModel2Response(user model.User) *proto.UserInfoResponse {
	return &proto.UserInfoResponse{
		Id:     user.ID,
		Mobile: user.Mobile,
		Role:   uint32(user.Role),
	}
}

func (u *UserHandler) GetUserList(ctx context.Context, pageInfo *proto.UserPageInfo) (*proto.UserListResponse, error) {
	zap.S().Infof("Received: %v", "获取用户列表")
	var userList []model.User
	result := global.DB.Find(&userList)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	global.DB.Scopes(utils.Paginate(int(pageInfo.Pn), int(pageInfo.Ps))).Find(&userList)
	for _, user := range userList {
		rsp.Data = append(rsp.Data, userModel2Response(user))
	}
	return rsp, nil
}

func (u *UserHandler) GetUserByMobile(ctx context.Context, mobileRequest *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("mobile = ?", mobileRequest.Mobile).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	return userModel2Response(user), nil
}

func (u *UserHandler) GetUserById(ctx context.Context, idRequest *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, idRequest.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	return userModel2Response(user), nil
}

func (u *UserHandler) CreateUser(ctx context.Context, createUserInfo *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	result := global.DB.Where(&model.User{Mobile: createUserInfo.Mobile}).First(&model.User{})
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user := model.User{
		Mobile: createUserInfo.Mobile,
		Role:   int32(createUserInfo.Role),
		Passwd: utils.Encrypt(createUserInfo.Password),
	}
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return userModel2Response(user), nil
}

func (u *UserHandler) UpdateUser(ctx context.Context, updateUserInfo *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	result := global.DB.First(&model.User{}, updateUserInfo.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	user := model.User{
		Passwd: utils.Encrypt(updateUserInfo.Password),
	}
	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (u *UserHandler) CheckPassWord(ctx context.Context, passwordCheckInfo *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	_, salt, encodedPwd := utils.GetEncoded(passwordCheckInfo.EncryptedPassword)
	result := utils.Verify(passwordCheckInfo.Password, salt, encodedPwd)
	return &proto.CheckResponse{Success: result}, nil
}

func (u *UserHandler) AddPassenger(ctx context.Context, passengerInfo *proto.PassengerInfo) (*emptypb.Empty, error) {
	panic("not implemented") // TODO: Implement
}

func (u *UserHandler) UpdatePassenger(ctx context.Context, passengerInfo *proto.PassengerInfo) (*emptypb.Empty, error) {
	panic("not implemented") // TODO: Implement
}

func (u *UserHandler) DeletePassenger(ctx context.Context, idRequest *proto.IdRequest) (*emptypb.Empty, error) {
	panic("not implemented") // TODO: Implement
}

func (u *UserHandler) GetPassengerList(ctx context.Context, passengerPageInfo *proto.PassengerPageInfo) (*proto.PassengerListResponse, error) {
	panic("not implemented") // TODO: Implement
}
