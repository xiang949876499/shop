package config

//type UserSrvConfig struct {
//	Host string `mapstructure:"host" json:"host"`
//	Port int    `mapstructure:"port" json:"port"`
//	Name string `mapstructure:"name" json:"name"`
//}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

//type AliSmsConfig struct {
//	ApiKey     string `mapstructure:"key" json:"key"`
//	ApiSecrect string `mapstructure:"secrect" json:"secrect"`
//}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
	Expire   int    `mapstructure:"expire" json:"expire"`
}

type ServerConfig struct {
	Name    string    `mapstructure:"name" json:"name"`
	Host    string    `mapstructure:"host" json:"host"`
	Tags    []string  `mapstructure:"tags" json:"tags"`
	Port    int       `mapstructure:"port" json:"port"`
	SrvName string    `mapstructure:"srv_name" json:"srv_name"`
	JWTInfo JWTConfig `mapstructure:"jwt" json:"jwt"`
	//AliSmsInfo  AliSmsConfig  `mapstructure:"sms" json:"sms"`
	RedisInfo  RedisConfig  `mapstructure:"redis" json:"redis"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
}

type LocalConfig struct {
	Nacos NacosConfig `mapstructure:"nacos" json:"nacos"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
