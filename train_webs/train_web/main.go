package main

import (
	"fmt"
	"gotrains/train_webs/train_web/global"
	"gotrains/train_webs/train_web/initialize"
	"gotrains/train_webs/train_web/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
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
	initialize.InitSentinel()
	initialize.InitConfigWithNacos()
	initialize.InitRedisDb()
	initialize.InitTrans("zh")
	initialize.InitValidator()
	// initialize.InitSrvConn()
	initialize.InitSrvConnWithLB()
	e := initialize.Routers()
	register_client := utils.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := register_client.Register(
		global.ServerConfig.Host,
		global.ServerConfig.Port,
		global.ServerConfig.Name,
		[]string{"train_web"},
		serviceId,
	)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
		return
	}
	zap.S().Infof("server run success on %d", global.ServerConfig.Port)
	go func() {
		if err := e.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("server run failed")
			panic(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Errorf("注销失败:%s", err.Error())
	}
	zap.S().Info("注销成功")
	zap.S().Info("服务退出")
}
