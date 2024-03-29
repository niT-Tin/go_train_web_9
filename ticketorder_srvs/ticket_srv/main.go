package main

import (
	"flag"
	"fmt"
	"gotrains/ticketorder_srvs/ticket_srv/global"
	"gotrains/ticketorder_srvs/ticket_srv/handler"
	"gotrains/ticketorder_srvs/ticket_srv/initialize"
	"gotrains/ticketorder_srvs/ticket_srv/proto"
	"gotrains/ticketorder_srvs/ticket_srv/utils"
	"gotrains/ticketorder_srvs/ticket_srv/utils/otgrpc"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

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
	Port := flag.Int("port", 50053, "端口号")
	// Port := flag.Int("port", 0, "端口号")
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitSrvConn()
	initialize.InitDB()
	flag.Parse()
	// g := grpc.NewServer()
	if *Port == 0 {
		var err error
		*Port, err = utils.GetFreePort()
		if err != nil {
			panic(err)
		}
	}
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	client := Register(global.Config.ConsulInfo.Host, *Port, "ticketorder_srv", []string{"ticketorder_srv"}, serviceId)
	//初始化jaeger
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", global.Config.JaegerInfo.Host, global.Config.JaegerInfo.Port),
		},
		ServiceName: "ticket_order",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	// server := grpc.NewServer()
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	// 使用默认的健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	proto.RegisterOrderServer(server, &handler.OrderServer{})

	lis, err := net.Listen("tcp", *IP+":"+fmt.Sprint(*Port))
	if err != nil {
		panic(err)
	}
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	closer.Close()
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Errorf("ticketorder注销失败:%s", err.Error())
	}
	zap.S().Info("ticketorder注销成功")
	zap.S().Info("ticketorder服务退出")
}
