package initialize

import (
	"fmt"
	"gotrains/ticketorder_web/ticketorder-web/global"
	"gotrains/ticketorder_web/ticketorder-web/proto"
	"gotrains/ticketorder_web/ticketorder-web/utils/otgrpc"

	"github.com/opentracing/opentracing-go"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GetServiceByName(name string) *api.AgentService {
	// 创建consul客户端
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	m, err2 := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`, name))
	if err2 != nil {
		panic(err2)
	}
	// 暂时不考虑多个实例的情况
	for _, v := range m {
		return v
	}
	return nil
}

func InitSrvConn() {
	as := GetServiceByName(global.ServerConfig.UserSrvConfig.Name)
	if as == nil || as.Address == "" {
		zap.S().Fatal("获取用户服务失败")
		return
	}
	ip := as.Address
	port := as.Port
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("GetUserList 连接用户服务失败", "msg", err.Error())
		return
	}
	global.UserClient = proto.NewUserClient(conn)
}

func initSrvGeneriticConn(host string, port int, name string, tag string, msg string) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=%s",
		host,
		port,
		name,
		tag,
	)
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		// insecure.NewCredentials()
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// 使用opentracing
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatalf("%s 连接%s服务失败", msg, name)
		return nil, err
	}
	return conn, nil
	// global.UserClient = proto.NewUserClient(conn)
}

func InitSrvConnWithLB() {
	userconn, err := initSrvGeneriticConn(
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port,
		global.ServerConfig.UserSrvConfig.Name,
		global.ServerConfig.UserSrvConfig.Name,
		"用户服务",
	)
	if err != nil {
		return
	}

	trainconn, err := initSrvGeneriticConn(
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port,
		global.ServerConfig.TrainSrvConfig.Name,
		global.ServerConfig.TrainSrvConfig.Name,
		"车次服务",
	)
	if err != nil {
		return
	}

	// ticketconn, err := initSrvGeneriticConn(
	// 	global.ServerConfig.ConsulInfo.Host,
	// 	global.ServerConfig.ConsulInfo.Port,
	// 	global.ServerConfig.TicketSrvConfig.Name,
	// 	global.ServerConfig.TicketSrvConfig.Name,
	// 	"车票服务",
	// )
	// if err != nil {
	// 	return
	// }

	orderconn, err := initSrvGeneriticConn(
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port,
		global.ServerConfig.OrderSrvConfig.Name,
		global.ServerConfig.OrderSrvConfig.Name,
		"订单服务",
	)
	if err != nil {
		return
	}
	global.UserClient = proto.NewUserClient(userconn)
	global.TrainClient = proto.NewTrainClient(trainconn)
	global.OrderClient = proto.NewOrderClient(orderconn)
}
