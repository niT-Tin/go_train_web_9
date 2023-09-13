package initialize

import (
	"gotrains/userpassenger_srvs/user_srv/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func InitConfig() {
	viper.AutomaticEnv()
	v := viper.New()
	if viper.GetBool("MY_DEV_VAL") {
		v.SetConfigFile("./config-dev.yaml")
	} else {
		v.SetConfigFile("./config-pro.yaml")
	}
	if err := v.ReadInConfig(); err != nil {
		zap.S().Errorf("Fatal error ReadInConfig config file: %s \n", err)
		panic(err)
	}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		zap.S().Errorf("Fatal error Unmarshal config file: %s \n", err)
		panic(err)
	}
	zap.S().Infof("Nacos config info: %#v \n", global.ServerConfig)

	nacos := global.ServerConfig.NacosInfo

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(nacos.Host, nacos.Port, constant.WithContextPath("/nacos")),
	}
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(nacos.Namespace),
		constant.WithTimeoutMs(nacos.TimeOutMS),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithCacheDir(nacos.CacheDir),
		constant.WithLogDir(nacos.LogDir),
		constant.WithLogLevel(nacos.LogLevel),
	)

	cfg, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}
	content, err := cfg.GetConfig(vo.ConfigParam{
		DataId: nacos.DataId,
		Group:  nacos.Group,
	})
	if err != nil {
		panic(err)
	}
	// cfg.ListenConfig(vo.ConfigParam{
	// 	onChange: func(namespace, group, dataId, data string) {
	// 		zap.S().Infof("config changed group:%s, dataId:%s, data:%s \n", group, dataId, data)
	// 		err := yaml.Unmarshal([]byte(data), global.ServerConfig)
	// 		if err != nil {
	// 			zap.S().Errorf("Fatal error Unmarshal config file: %s \n", err)
	// 			panic(err)
	// 		}
	// 		zap.S().Infof("Nacos config info: %#v \n", global.ServerConfig)
	// 	},
	// })

	// global.ServerConfig = &config.ServerConfig{}
	yaml.Unmarshal([]byte(content), global.Config)
	zap.S().Infof("config info: %#v \n", global.Config)
	// v.WatchConfig()
	// v.OnConfigChange(func(e fsnotify.Event) {
	// 	zap.S().Infof("配置文件发生变化config file changed: %s \n", e.Name)
	// 	v.ReadInConfig()
	// 	if err := v.Unmarshal(global.Config); err != nil {
	// 		zap.S().Errorf("Fatal error Unmarshal config file: %s \n", err)
	// 		panic(err)
	// 	}
	// 	zap.S().Infof("config info: %#v \n", global.Config)
	// })
}
