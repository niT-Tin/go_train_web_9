package main

import (
	"fmt"
	"gotrains/ticketorder_web/ticketorder-web/global"
	"gotrains/ticketorder_web/ticketorder-web/initialize"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func Register(address string, port int, name string, tags []string, id string) {
	// 创建consul客户端
	cfg := api.DefaultConfig()
	// consul服务地址
	cfg.Address = fmt.Sprintf("%s:%d", address, 8500)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	// 本服务地址
	registration.Address = address

	check := new(api.AgentServiceCheck)
	check.Interval = "5s"
	check.GRPC = fmt.Sprintf("%s:%d", address, port)
	check.Timeout = "3s"
	check.DeregisterCriticalServiceAfter = "30s"

	registration.Check = check

	err2 := client.Agent().ServiceRegister(registration)
	if err2 != nil {
		panic(err2)
	}
}

func main() {
	initialize.InitLogger()
	initialize.InitConfigWithNacos()
	initialize.InitRedisDb()
	initialize.InitTrans("zh")
	initialize.InitValidator()
	// initialize.InitSrvConn()
	initialize.InitSrvConnWithLB()
	e := initialize.Routers()
	zap.S().Infof("server run success on %s", ":8080")
	if err := e.Run(global.ServerConfig.Port); err != nil {
		zap.S().Panic("server run failed")
		panic(err)
	}
}
