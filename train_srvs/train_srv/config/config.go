package config

type MysqlConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	User string `mapstructure:"user" json:"user" yaml:"user"`
	Pass string `mapstructure:"pass" json:"pass" yaml:"pass"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type RocketMQConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type ServerConfig struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	// Host       string       `mapstructure:"host" json:"host" yaml:"host"`
	// Port       int          `mapstructure:"port" json:"port" yaml:"port"`
	MySqlInfo      MysqlConfig    `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	ConsulInfo     ConsulConfig   `mapstructure:"consul" json:"consul" yaml:"consul"`
	RocketMQConfig RocketMQConfig `mapstructure:"rocketmq" json:"rocketmq" yaml:"rocketmq"`
}

type NewServerConfig struct {
	NacosInfo NacosConfig `mapstructure:"nacos" yaml:"nacos"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" yaml:"host"`
	Port      uint64 `mapstructure:"port" yaml:"port"`
	Namespace string `mapstructure:"namespace" yaml:"namespace"`
	User      string `mapstructure:"user" yaml:"user"`
	Password  string `mapstructure:"password" yaml:"password"`
	DataId    string `mapstructure:"dataid" yaml:"dataid"`
	Group     string `mapstructure:"group" yaml:"group"`
	TimeOutMS uint64 `mapstructure:"timeoutms" yaml:"timeoutms"`
	CacheDir  string `mapstructure:"cachedir" yaml:"cachedir"`
	LogDir    string `mapstructure:"logdir" yaml:"logdir"`
	LogLevel  string `mapstructure:"loglevel" yaml:"loglevel"`
}
