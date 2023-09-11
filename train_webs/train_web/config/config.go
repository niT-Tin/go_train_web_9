package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port int    `mapstructure:"port" yaml:"port"`
	Name string `mapstructure:"name" yaml:"name"`
}

type TrainSrvConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port int    `mapstructure:"port" yaml:"port"`
	Name string `mapstructure:"name" yaml:"name"`
}

type TicketSrvConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port int    `mapstructure:"port" yaml:"port"`
	Name string `mapstructure:"name" yaml:"name"`
}

// type OrderSrvConfig struct {
// 	Host string `mapstructure:"host" yaml:"host"`
// 	Port int    `mapstructure:"port" yaml:"port"`
// 	Name string `mapstructure:"name" yaml:"name"`
// }

type SeatSrvConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port int    `mapstructure:"port" yaml:"port"`
	Name string `mapstructure:"name" yaml:"name"`
}

type ServerConfig struct {
	Name            string          `mapstructure:"name" yaml:"name"`
	Port            string          `mapstructure:"port" yaml:"port"`
	UserSrvConfig   UserSrvConfig   `mapstructure:"user_srv" yaml:"user_srv"`
	TrainSrvConfig  TrainSrvConfig  `mapstructure:"train_srv" yaml:"train_srv"`
	TicketSrvConfig TicketSrvConfig `mapstructure:"ticket_srv" yaml:"ticket_srv"`
	SeatSrvConfig   SeatSrvConfig   `mapstructure:"seat_srv" yaml:"seat_srv"`
	// JWTInfo       JWTConfig     `mapstructure:"jwt" yaml:"jwt"`
	// AliSmsConfig  AliSmsConfig  `mapstructure:"ali_sms" yaml:"ali_sms"`
	RedisConfig RedisConfig  `mapstructure:"redis" yaml:"redis"`
	ConsulInfo  ConsulConfig `mapstructure:"consul" yaml:"consul"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" yaml:"key"`
}

type AliSmsConfig struct {
	AccessKeyId  string `mapstructure:"access_key_id" yaml:"access_key_id"`
	AccessSecret string `mapstructure:"access_secret" yaml:"access_secret"`
	TemplateCode string `mapstructure:"template_code" yaml:"template_code"`
	SignName     string `mapstructure:"sign_name" yaml:"sign_name"`
	CodeLen      int    `mapstructure:"code_len" yaml:"code_len"`
	Expire       int    `mapstructure:"expire" yaml:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port int    `mapstructure:"port" yaml:"port"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port int    `mapstructure:"port" yaml:"port"`
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
