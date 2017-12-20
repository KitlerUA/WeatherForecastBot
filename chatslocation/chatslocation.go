package chatslocation

import (
	"context"
	"log"

	"github.com/KitlerUA/WeatherForecastBot/config"
	"github.com/KitlerUA/WeatherForecastBot/db"
)

var DefaultLocationByChatID map[int64]int

const indexMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"location":{
			"properties":{
				"chatid":{
					"type":"long"
				},
				"location":{
					"type":"int"
				}
			}
		}
	}
}`

func AddOrUpdate(chatID int64) error {
	_, err := db.Get().ElasticsearchVersion(config.Get().ElasticAddress)
	if err != nil {
		log.Fatalf("Cannot ping elastic %s", err)
	}
	ctx := context.Background()
	exists, err := db.Get().IndexExists("wetbot").Do(ctx)
	if err != nil {
		log.Fatalf("Cannot check index: %s", err)
	}
	if !exists {
		createIndex, err := db.Get().CreateIndex("wetbot").BodyString(indexMapping).Do(ctx)
		if err != nil {
			log.Fatalf("Cannot create index: %s ", err)
		}
		if !createIndex.Acknowledged {
			log.Fatal("Not acknowledge")
		}
	}

}
