package initialize

func Init() {
	InitLogger()
	InitConfig()
	InitRedis()
	InitSrvConn()

	if err := InitTrans("zh"); err != nil {
		panic(err)
	}
}
