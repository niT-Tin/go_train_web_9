package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gotrains/ticketorder_srvs/ticket_srv/global"
	"gotrains/ticketorder_srvs/ticket_srv/model"
	"gotrains/ticketorder_srvs/ticket_srv/proto"
	"gotrains/userpassenger_srvs/user_srv/utils"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderListener struct{}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.Order
	_ = json.Unmarshal(msg.Body, &orderInfo)
	global.DB.Begin()

	// TODO: 链路追踪
	resp, err := global.UserClient.GetPassengerList(context.Background(), &proto.PassengerPageInfo{
		UserId: orderInfo.UserID,
	})

	if err != nil {
		global.DB.Rollback()
		zap.S().Errorf("获取乘客列表失败: %s", err.Error())
		return primitive.RollbackMessageState
	}
	seats := []*proto.SeatInfo{}
	for _, ps := range resp.Data {
		seat := proto.SeatInfo{
			TrainCode: orderInfo.TrainCode,
			// 默认后端计算第几节车厢
			CarriageIndex: 0,
			SeatType:      ps.SeatType,
			// 时间紧，先默认第一个座位
			SeatIndex: 0,
			// 后端计算
			Row:        orderInfo.Row,
			Column:     orderInfo.Colum,
			CreateTime: time.Now().Unix(),
			UpdateTime: time.Now().Unix(),
		}
		seats = append(seats, &seat)
	}

	// TODO: 链路追踪
	if _, err := global.TicketInventoryClient.ReductTicket(context.Background(), &proto.BusinessRequest{
		Seats:        seats,
		OrderId:      orderInfo.OrderSn,
		StartStation: orderInfo.StartStation,
		EndStation:   orderInfo.EndStation,
		StartTime:    orderInfo.StartTime.Format("2006-01-02 15:04:05"),
	}); err != nil {
		global.DB.Rollback()
		return primitive.RollbackMessageState
	}

	result := global.DB.Create(&orderInfo)
	if result.Error != nil {
		global.DB.Rollback()
		zap.S().Errorf("创建订单失败: %s", result.Error.Error())
		return primitive.CommitMessageState
	}
	global.DB.Commit()
	return primitive.RollbackMessageState
}

// 消息回查
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.Order
	_ = json.Unmarshal(msg.Body, &orderInfo)
	var count int64
	global.DB.Model(&model.Order{}).Where("order_sn = ? AND is_delete = ?", orderInfo.OrderSn, false).Count(&count)
	if count > 0 {
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func OrderModel2Info(order model.Order) *proto.OrderInfo {
	return &proto.OrderInfo{
		Id:           order.ID,
		UserId:       order.UserID,
		TrainId:      order.TrainID,
		StartStation: order.StartStation,
		EndStation:   order.EndStation,
		StartTime:    order.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:      order.EndTime.Format("2006-01-02 15:04:05"),
		OrderSn:      order.OrderSn,
		SeatType:     order.SeatType,
		Price:        order.Pirce,
		SeatNumber:   order.SeatNumber,
	}
}

func OrderModel2Response(order model.Order) *proto.OrderInfoResponse {
	return &proto.OrderInfoResponse{
		OrderSn: order.OrderSn,
		UserId:  fmt.Sprintf("%d", order.UserID),
	}
}

func (o *OrderServer) CreateOrder(ctx context.Context, createinfo *proto.CreateOrderInfo) (*proto.OrderInfoResponse, error) {
	var order model.Order
	orderListener := &OrderListener{}
	p, err := rocketmq.NewTransactionProducer(
		orderListener,
		producer.WithNameServer([]string{"localhost:9876"}),
	)
	if err != nil {
		zap.S().Errorf("生成producer失败: %s", err.Error())
		return nil, err
	}
	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败: %s", err.Error())
		return nil, err
	}
	result := global.DB.First(&order, createinfo.OrderInfo.OrderSn)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "订单已存在")
	}
	st, _ := time.Parse("2006-01-02 15:04:05", createinfo.OrderInfo.StartTime)
	et, _ := time.Parse("2006-01-02 15:04:05", createinfo.OrderInfo.EndTime)
	var passengerIds model.PassengerIds
	for _, ps := range createinfo.OrderInfo.Passengers {
		passengerIds = append(passengerIds, fmt.Sprintf("%d", ps.Id))
	}
	order = model.Order{
		UserID:       createinfo.OrderInfo.UserId,
		TrainID:      createinfo.OrderInfo.TrainId,
		StartStation: createinfo.OrderInfo.StartStation,
		EndStation:   createinfo.OrderInfo.EndStation,
		StartTime:    st,
		EndTime:      et,
		SeatType:     createinfo.OrderInfo.SeatType,
		SeatNumber:   createinfo.OrderInfo.SeatNumber,
		Pirce:        createinfo.OrderInfo.Price,
		OrderSn:      uuid.NewV4().String(),
		PassengerIds: passengerIds,
	}
	jsonString, _ := json.Marshal(order)

	_, err = p.SendMessageInTransaction(context.Background(),
		primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}

	return &proto.OrderInfoResponse{
		OrderSn: order.OrderSn,
		UserId:  fmt.Sprintf("%d", order.UserID),
	}, nil
}

func (o *OrderServer) GetOrderList(ctx context.Context, orderpage *proto.OrderPageInfo) (*proto.OrderListResponse, error) {
	var orders []model.Order
	res := global.DB.Where(&model.Order{UserID: orderpage.UserId}).Find(&orders)
	if res.Error != nil {
		zap.S().Errorf("查询订单失败: %s", res.Error.Error())
		return nil, status.Error(codes.Internal, "查询订单失败")
	}
	orderList := []*proto.OrderInfoResponse{}
	for _, order := range orders {
		orderList = append(orderList, OrderModel2Response(order))
	}
	return &proto.OrderListResponse{
		Data: orderList,
	}, nil
}

func (o *OrderServer) GetOrderById(ctx context.Context, oidreq *proto.OIdRequest) (*proto.OrderInfoResponse, error) {
	var order model.Order
	res := global.DB.First(&order, oidreq.UserId)
	if res.Error != nil {
		zap.S().Errorf("查询订单失败：用户id: %d", oidreq.UserId)
		return &proto.OrderInfoResponse{}, res.Error
	}
	if res.RowsAffected == 0 {
		zap.S().Errorf("查询订单失败：用户id: %d", oidreq.UserId)
		return &proto.OrderInfoResponse{}, status.Errorf(codes.NotFound, "该订单不存在")
	}
	return OrderModel2Response(order), nil

}

func (o *OrderServer) UpdateOrder(ctx context.Context, updateorder *proto.UpdateOrderInfo) (*emptypb.Empty, error) {
	var order model.Order
	result := global.DB.Model(&order).Where("order_sn = ?", updateorder.OrderSn)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "该订单不存在")
	}
	for _, ps := range updateorder.Passengers {
		order.PassengerIds = append(order.PassengerIds, fmt.Sprintf("%d", ps.Id))
	}
	result = global.DB.Save(order)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (o *OrderServer) DeleteOrder(ctx context.Context, oidreq *proto.OIdRequest) (*emptypb.Empty, error) {
	var order model.Order
	result := global.DB.Model(&order).Where("order_sn = ?", oidreq.UserId)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "该订单不存在")
	}
	result = global.DB.Delete(&order)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (o *OrderServer) GetOrderListByUserId(ctx context.Context, orderpage *proto.OrderPageInfo) (*proto.OrderListResponse, error) {
	var order []model.Order
	res := global.DB.Scopes(utils.Paginate(int(orderpage.Pn), int(orderpage.Ps))).Find(&order, orderpage.UserId)
	if res.Error != nil {
		zap.S().Errorf("查询订单失败：用户id: %d", orderpage.UserId)
		return &proto.OrderListResponse{}, res.Error
	}
	if res.RowsAffected == 0 {
		zap.S().Errorf("查询订单失败：用户id: %d", orderpage.UserId)
		return &proto.OrderListResponse{}, status.Errorf(codes.NotFound, "该订单不存在")
	}
	var orderList []*proto.OrderInfoResponse
	for _, o := range order {
		orderList = append(orderList, OrderModel2Response(o))
	}
	return &proto.OrderListResponse{Data: orderList}, nil
}
