package initialize

func Init() {
	InitLogger()
	InitConfig()
	InitDB()
	InitRedis()
	//InitEs()
}
