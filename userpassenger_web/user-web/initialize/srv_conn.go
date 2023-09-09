package initialize

import (
	"fmt"
	"gotrains/userpassenger_web/user-web/global"
	"gotrains/userpassenger_web/user-web/proto"

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

func InitSrvConnWithLB() {
	addr := fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=%s",
		global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port,
		global.ServerConfig.UserSrvConfig.Name,
		global.ServerConfig.UserSrvConfig.Name,
	)
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		// insecure.NewCredentials()
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalw("GetUserList 连接用户服务失败", "msg", err.Error())
		return
	}
	global.UserClient = proto.NewUserClient(conn)
}
