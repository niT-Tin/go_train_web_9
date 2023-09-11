package initialize

import (
	"fmt"

	"gotrains/ticketorder_srvs/ticket_srv/global"

	"gotrains/ticketorder_srvs/ticket_srv/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	consulInfo := global.Config.ConsulInfo
	// 初始化票务服务连接
	ticketConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.Config.TrainSrvName),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		//grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[TicketInventoryClient] 连接 【票务服务连接失败】")
	}

	global.TicketInventoryClient = proto.NewTicketClient(ticketConn)

	// 初始化用户服务连接
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.Config.UserSrvName),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		//grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[UserClient] 连接 【用户服务连接失败】")
	}

	global.UserClient = proto.NewUserClient(userConn)
}
