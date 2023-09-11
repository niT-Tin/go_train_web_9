// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: order.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Order_CreateOrder_FullMethodName          = "/proto.Order/CreateOrder"
	Order_GetOrderList_FullMethodName         = "/proto.Order/GetOrderList"
	Order_GetOrderById_FullMethodName         = "/proto.Order/GetOrderById"
	Order_UpdateOrder_FullMethodName          = "/proto.Order/UpdateOrder"
	Order_DeleteOrder_FullMethodName          = "/proto.Order/DeleteOrder"
	Order_GetOrderListByUserId_FullMethodName = "/proto.Order/GetOrderListByUserId"
)

// OrderClient is the client API for Order service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderClient interface {
	CreateOrder(ctx context.Context, in *CreateOrderInfo, opts ...grpc.CallOption) (*OrderInfoResponse, error)
	GetOrderList(ctx context.Context, in *OrderPageInfo, opts ...grpc.CallOption) (*OrderListResponse, error)
	GetOrderById(ctx context.Context, in *OIdRequest, opts ...grpc.CallOption) (*OrderInfoResponse, error)
	UpdateOrder(ctx context.Context, in *UpdateOrderInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteOrder(ctx context.Context, in *OIdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetOrderListByUserId(ctx context.Context, in *OrderPageInfo, opts ...grpc.CallOption) (*OrderListResponse, error)
}

type orderClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderClient(cc grpc.ClientConnInterface) OrderClient {
	return &orderClient{cc}
}

func (c *orderClient) CreateOrder(ctx context.Context, in *CreateOrderInfo, opts ...grpc.CallOption) (*OrderInfoResponse, error) {
	out := new(OrderInfoResponse)
	err := c.cc.Invoke(ctx, Order_CreateOrder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) GetOrderList(ctx context.Context, in *OrderPageInfo, opts ...grpc.CallOption) (*OrderListResponse, error) {
	out := new(OrderListResponse)
	err := c.cc.Invoke(ctx, Order_GetOrderList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) GetOrderById(ctx context.Context, in *OIdRequest, opts ...grpc.CallOption) (*OrderInfoResponse, error) {
	out := new(OrderInfoResponse)
	err := c.cc.Invoke(ctx, Order_GetOrderById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) UpdateOrder(ctx context.Context, in *UpdateOrderInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Order_UpdateOrder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) DeleteOrder(ctx context.Context, in *OIdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Order_DeleteOrder_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderClient) GetOrderListByUserId(ctx context.Context, in *OrderPageInfo, opts ...grpc.CallOption) (*OrderListResponse, error) {
	out := new(OrderListResponse)
	err := c.cc.Invoke(ctx, Order_GetOrderListByUserId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServer is the server API for Order service.
// All implementations must embed UnimplementedOrderServer
// for forward compatibility
type OrderServer interface {
	CreateOrder(context.Context, *CreateOrderInfo) (*OrderInfoResponse, error)
	GetOrderList(context.Context, *OrderPageInfo) (*OrderListResponse, error)
	GetOrderById(context.Context, *OIdRequest) (*OrderInfoResponse, error)
	UpdateOrder(context.Context, *UpdateOrderInfo) (*emptypb.Empty, error)
	DeleteOrder(context.Context, *OIdRequest) (*emptypb.Empty, error)
	GetOrderListByUserId(context.Context, *OrderPageInfo) (*OrderListResponse, error)
	mustEmbedUnimplementedOrderServer()
}

// UnimplementedOrderServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServer struct {
}

func (UnimplementedOrderServer) CreateOrder(context.Context, *CreateOrderInfo) (*OrderInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedOrderServer) GetOrderList(context.Context, *OrderPageInfo) (*OrderListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderList not implemented")
}
func (UnimplementedOrderServer) GetOrderById(context.Context, *OIdRequest) (*OrderInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderById not implemented")
}
func (UnimplementedOrderServer) UpdateOrder(context.Context, *UpdateOrderInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrder not implemented")
}
func (UnimplementedOrderServer) DeleteOrder(context.Context, *OIdRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteOrder not implemented")
}
func (UnimplementedOrderServer) GetOrderListByUserId(context.Context, *OrderPageInfo) (*OrderListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderListByUserId not implemented")
}
func (UnimplementedOrderServer) mustEmbedUnimplementedOrderServer() {}

// UnsafeOrderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServer will
// result in compilation errors.
type UnsafeOrderServer interface {
	mustEmbedUnimplementedOrderServer()
}

func RegisterOrderServer(s grpc.ServiceRegistrar, srv OrderServer) {
	s.RegisterService(&Order_ServiceDesc, srv)
}

func _Order_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_CreateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).CreateOrder(ctx, req.(*CreateOrderInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_GetOrderList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderPageInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).GetOrderList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_GetOrderList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).GetOrderList(ctx, req.(*OrderPageInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_GetOrderById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).GetOrderById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_GetOrderById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).GetOrderById(ctx, req.(*OIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_UpdateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOrderInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).UpdateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_UpdateOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).UpdateOrder(ctx, req.(*UpdateOrderInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_DeleteOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).DeleteOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_DeleteOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).DeleteOrder(ctx, req.(*OIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Order_GetOrderListByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderPageInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).GetOrderListByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_GetOrderListByUserId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).GetOrderListByUserId(ctx, req.(*OrderPageInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Order_ServiceDesc is the grpc.ServiceDesc for Order service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Order_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Order",
	HandlerType: (*OrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _Order_CreateOrder_Handler,
		},
		{
			MethodName: "GetOrderList",
			Handler:    _Order_GetOrderList_Handler,
		},
		{
			MethodName: "GetOrderById",
			Handler:    _Order_GetOrderById_Handler,
		},
		{
			MethodName: "UpdateOrder",
			Handler:    _Order_UpdateOrder_Handler,
		},
		{
			MethodName: "DeleteOrder",
			Handler:    _Order_DeleteOrder_Handler,
		},
		{
			MethodName: "GetOrderListByUserId",
			Handler:    _Order_GetOrderListByUserId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}
