package initialize

import (
	"gotrains/userpassenger_web/user-web/global"
	"gotrains/userpassenger_web/user-web/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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
	zap.S().Infof("config info: %#v \n", global.ServerConfig)
	// 如果不是开发环境，就获取一个空闲端口
	if !viper.GetBool("MY_DEV_VAL") {
		zap.S().Info("using config-dev.yaml")
		port, err := utils.GetFreePort()
		if err != nil {
			panic(err)
		}
		global.ServerConfig.Port = int32(port)
	}
	zap.S().Infof("config info: %#v \n", global.ServerConfig)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化config file changed: %s \n", e.Name)
		v.ReadInConfig()
		if err := v.Unmarshal(global.ServerConfig); err != nil {
			zap.S().Errorf("Fatal error Unmarshal config file: %s \n", err)
			panic(err)
		}
		zap.S().Infof("config info: %#v \n", global.ServerConfig)
	})
}

func InitConfigWithNacos() {
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
	if err := v.Unmarshal(global.NewServerConfig); err != nil {
		zap.S().Errorf("Fatal error Unmarshal config file: %s \n", err)
		panic(err)
	}
	zap.S().Infof("Nacos config info: %#v \n", global.NewServerConfig)
	nacos := global.NewServerConfig.NacosInfo

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

	// global.ServerConfig = &config.ServerConfig{}
	yaml.Unmarshal([]byte(content), global.ServerConfig)
	zap.S().Infof("config info: %#v \n", global.ServerConfig)
}
