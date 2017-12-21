package db

import (
	"sync"

	"github.com/olivere/elastic"
	econf "github.com/olivere/elastic/config"

	"os"

	"time"

	"github.com/KitlerUA/WeatherForecastBot/config"
	"github.com/yanzay/log"
)

var once sync.Once
var client *elastic.Client

func Get() *elastic.Client {
	once.Do(connectToBD)
	return client
}

func connectToBD() {
	for {
		if len(os.Args) > 1 {
			var err error
			client, err = elastic.NewClientFromConfig(&econf.Config{URL: os.Args[1]})
			if err != nil {
				log.Printf("Cannot connect to Elastic %s", err)
				time.Sleep(2000)
				continue
			}
			break
		} else {
			var err error
			client, err = elastic.NewClientFromConfig(&econf.Config{URL: config.Get().ElasticAddress})

			//client, err = elastic.NewClient()
			if err != nil {
				log.Printf("Cannot connect to Elasticsearch  %s", err)
				time.Sleep(2000)
				continue
			}
			break
		}
	}
}
