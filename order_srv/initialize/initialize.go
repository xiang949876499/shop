package initialize

func Init() {
	InitLogger()
	InitConfig()
	InitDB()
	InitSrvConn()
	//InitEs()
}
