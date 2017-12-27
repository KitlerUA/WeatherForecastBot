package weather

import (
	"context"
	"log"
	"time"

	"github.com/KitlerUA/WeatherForecastBot/chatslocation"
	"github.com/KitlerUA/WeatherForecastBot/db"
	"github.com/olivere/elastic"
)

func Update() {
	ticker := time.NewTicker(3 * time.Hour)
	ctx := context.Background()
	for range ticker.C {
		locations := chatslocation.GetUnique()
		log.Printf("Locations %v", locations)
		for i := range locations {
			weather, err := getWeatherFromOpenMap(locations[i])
			if err != nil {
				log.Printf("Cannot get weather from OpenMap: %s", err)
				continue
			}
			for _, info := range weather.List {
				eInfo := infoToElasticInfo(info, locations[i])
				queryDate := elastic.NewTermQuery("Location", locations[i])
				queryLocation := elastic.NewTermQuery("DtTxt", eInfo.DtTxt)
				query := elastic.NewBoolQuery().Must(queryLocation, queryDate)
				searchResult, err := db.Get().Search("wetbot").Type("info").Query(query).Do(ctx)
				if err != nil {
					log.Printf("Cannot search for forecast in Elasticsearch: %s", err)
				}
				if searchResult.TotalHits() > 0 {
					hitID := searchResult.Hits.Hits[0].Id
					_, err = db.Get().Update().Index("wetbot").Type("info").Id(hitID).
						Script(elastic.NewScriptInline("ctx._source.Description = params.description; ctx._source.Temp = params.temp; ctx._source.Humidity = params.humidity; ctx._source.IconID = params.iconID;").Lang("painless").
							Param("description", eInfo.Description).Param("temp", eInfo.Temp).Param("humidity", eInfo.Humidity).Param("iconID", eInfo.IconID)).
						Do(ctx)
					if err != nil {
						log.Printf("Cannot update doc %s : %s", hitID, err)
						continue
					}
				} else {
					_, err := db.Get().Index().Index("wetbot").Type("info").BodyJson(eInfo).Do(ctx)
					if err != nil {
						log.Printf("Cannot add new foewcast: %s", err)
						continue
					}
				}
			}
			_, err = db.Get().Flush().Index("wetbot").Do(ctx)
			if err != nil {
				log.Printf("Cannot flush index: %s", err)
			}
			log.Printf("Forecast for %d updated", locations[i])
		}

	}
}
