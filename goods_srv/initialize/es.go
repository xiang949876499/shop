package initialize

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"

	"goods_srv/global"
	"goods_srv/model"
)

func InitEs() {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.Es.Host, global.ServerConfig.Es.Port)
	logger := log.New(os.Stdout, "shop", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	//新建mapping和index
	exists, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		_, err = global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
