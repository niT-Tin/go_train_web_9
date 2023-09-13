package main

import (
	"fmt"
	"gotrains/userpassenger_web/user-web/global"
	"gotrains/userpassenger_web/user-web/initialize"
	"gotrains/userpassenger_web/user-web/utils"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
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
		int(global.ServerConfig.Port),
		global.ServerConfig.Name,
		[]string{"user_web"},
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
