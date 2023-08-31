package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func main() {
	host := "http://192.168.32.192:9200"
	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	q := elastic.NewMatchQuery("address", "street")
	res, err := client.Search().Index("user").Query(q).Do(context.Background())
	if err != nil {
		panic(err)
	}
	total := res.Hits.TotalHits.Value
	fmt.Printf("%d\n", total)
	for _, v := range res.Hits.Hits {
		if data, err := v.Source.MarshalJSON(); err == nil {
			fmt.Println(data)
		} else {
			panic(err)
		}
	}

}
