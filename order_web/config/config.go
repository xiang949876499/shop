package config

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type AlipayConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Port         int            `mapstructure:"port" json:"port"`
	Tags         []string       `mapstructure:"tags" json:"tags"`
	GoodsSrv     GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	OrderSrv     GoodsSrvConfig `mapstructure:"order_srv" json:"order_srv"`
	InventorySrv GoodsSrvConfig `mapstructure:"order_srv" json:"inventory_srv"`
	JWT          JWTConfig      `mapstructure:"jwt" json:"jwt"`
	Consul       ConsulConfig   `mapstructure:"consul" json:"consul"`
	Jaeger       JaegerConfig   `mapstructure:"consul" json:"jaeger"`
	//AliPayInfo       AlipayConfig   `mapstructure:"alipay" json:"alipay"`
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

type LocalConfig struct {
	Nacos NacosConfig `mapstructure:"nacos" json:"nacos"`
}
