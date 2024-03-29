package main

import (
	"flag"
	"fmt"
	"gotrains/train_srvs/train_srv/global"
	"gotrains/train_srvs/train_srv/handler"
	"gotrains/train_srvs/train_srv/initialize"
	"gotrains/train_srvs/train_srv/proto"
	"gotrains/train_srvs/train_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func Register(address string, port int, name string, tags []string, id string) *api.Client {
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
	return client
}

func GetServiceByName(name string) {
	// 创建consul客户端
	cfg := api.DefaultConfig()
	cfg.Address = "10.102.213.148:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	m, err2 := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`, name))
	if err2 != nil {
		panic(err2)
	}
	for k := range m {
		fmt.Println("key: ", k)
	}
}

func main() {

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50052, "端口号")
	// Port := flag.Int("port", 0, "端口号")
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	flag.Parse()
	g := grpc.NewServer()
	if *Port == 0 {
		var err error
		*Port, err = utils.GetFreePort()
		if err != nil {
			panic(err)
		}
	}
	// 使用默认的健康检查
	grpc_health_v1.RegisterHealthServer(g, health.NewServer())
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	client := Register(global.Config.ConsulInfo.Host, *Port, "train_srv", []string{"train_srv"}, serviceId)
	proto.RegisterTrainServer(g, &handler.TrainServer{})
	proto.RegisterStationServer(g, &handler.StationServer{})
	proto.RegisterCarriageServer(g, &handler.CarriageServer{})
	proto.RegisterSeatServer(g, &handler.SeatServer{})
	proto.RegisterTicketServer(g, &handler.TicketServer{})

	lis, err := net.Listen("tcp", *IP+":"+fmt.Sprint(*Port))
	if err != nil {
		panic(err)
	}

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.Config.RocketMQConfig.Host, global.Config.RocketMQConfig.Port)}),
		consumer.WithGroupName("order_group"),
	)

	if err := c.Subscribe("order_reback", consumer.MessageSelector{}, handler.ReBack); err != nil {
		fmt.Println("读取消息失败")
	}
	_ = c.Start()
	go func() {
		err = g.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Errorf("train_srv注销失败:%s", err.Error())
	}
	zap.S().Info("train_srv注销成功")
	zap.S().Info("train_srv服务退出")
}
