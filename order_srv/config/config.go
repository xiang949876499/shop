package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type EsConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	db   string `mapstructure:"db" json:"db"`
	//User     string `mapstructure:"user" json:"user"`
	//Password string `mapstructure:"password" json:"password"`
}

type ServerConfig struct {
	Name         string             `mapstructure:"name" json:"name"`
	Host         string             `mapstructure:"host" json:"host"`
	Port         int                `mapstructure:"port" json:"port"`
	Tags         []string           `mapstructure:"tags" json:"tags"`
	Consul       ConsulConfig       `mapstructure:"consul" json:"consul"`
	Redis        RedisConfig        `mapstructure:"redis" json:"redis"`
	Es           EsConfig           `mapstructure:"es" json:"es"`
	Rocketmq     RocketmqConfig     `mapstructure:"rocketmq" json:"rocketmq"`
	GoodSrv      GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrv InventorySrvConfig `mapstructure:"Inventory_srv" json:"Inventory_srv"`
}

type RocketmqConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	GroupName string `mapstructure:"group_name" json:"group_name"`
}

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type InventorySrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"data_id"`
	Group     string `mapstructure:"group" json:"group"`
}

type LocalConfig struct {
	Nacos NacosConfig `mapstructure:"nacos" json:"nacos"`
	Mysql MysqlConfig `mapstructure:"mysql" json:"mysql"`
}
