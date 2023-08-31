package initialize

func Init() {
	InitLogger()
	InitConfig()
	InitSrvConn()
	InitSentinel()
	//4. 初始化翻译
	if err := InitTrans("zh"); err != nil {
		panic(err)
	}
}
