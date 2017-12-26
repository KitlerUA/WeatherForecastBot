package chatslocation

import (
	"context"
	"log"

	"io/ioutil"

	"encoding/json"

	"errors"
	"fmt"

	"reflect"

	"github.com/KitlerUA/WeatherForecastBot/config"
	"github.com/KitlerUA/WeatherForecastBot/db"
	"github.com/olivere/elastic"
)

type LocationElastic struct {
	ChatID   int64
	Location int
}

type CityElastic struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

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

const cityListMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"city":{
			"properties":{
				"id":{
					"type":"int"
				},
				"name":{
					"type":"text"
				},
				"country":{
					"type":"keyword"
				}
			}
		}
	}
}`

func AddOrUpdate(chatID int64, location int) error {
	_, err := db.Get().ElasticsearchVersion(config.Get().ElasticAddress)
	if err != nil {
		log.Fatalf("Cannot ping elastic %s", err)
		return err
	}
	ctx := context.Background()
	query := elastic.NewTermQuery("ChatID", chatID)
	searchResult, err := db.Get().Search("locbot").Type("location").Query(query).Do(ctx)
	if err != nil {
		log.Printf("Cannot search location: %s", err)
		return err
	}
	if searchResult.TotalHits() > 0 {
		hitID := searchResult.Hits.Hits[0].Id
		_, err = db.Get().Update().Index("locbot").Type("location").Id(hitID).
			Script(elastic.NewScriptInline("ctx._source.Location = params.newl").Lang("painless").Param("newl", location)).
			Do(ctx)
		if err != nil {
			log.Printf("Cannot update doc %s : %s", hitID, err)
			return err
		}
	} else {
		log.Printf("Location for %d not found, creating new", chatID)
		loc := LocationElastic{
			ChatID:   chatID,
			Location: location,
		}
		_, err := db.Get().Index().Index("locbot").Type("location").BodyJson(loc).Do(ctx)
		if err != nil {
			log.Fatalf("Cannot put %v: %s", loc, err)
			return err
		}
	}
	return nil
}

func LoadCityList() {
	ctx := context.Background()
	exists, err := db.Get().TypeExists().Index("citbot").Type("city").Do(ctx)
	if err != nil {
		log.Fatalf("Cannot check City type exists: %s", err)
	}
	if exists {
		return
	}
	data, err := ioutil.ReadFile(config.Get().LocationsCodeFileName)
	if err != nil {
		log.Fatalf("Cannot read city list from file: %s", err)
	}
	var cityList []CityElastic
	if err = json.Unmarshal(data, &cityList); err != nil {
		log.Fatalf("Corrupted data in citylist file: %s", err)
	}
	//db.Get().PutMapping().BodyString(cityListMapping)
	for _, c := range cityList {
		_, err := db.Get().Index().Index("citbot").Type("city").BodyJson(c).Do(ctx)
		if err != nil {
			log.Printf("Cannot put %v: %s", c, err)
		}
	}
}

func Get(chatID int64) (int, error) {
	query := elastic.NewTermQuery("ChatID", chatID)
	ctx := context.Background()
	searchResult, err := db.Get().Search("locbot").Type("location").Query(query).Do(ctx)
	if err != nil {
		log.Printf("Cannot search location for chat %d : %s", chatID, err)
		return 0, err
	}
	if searchResult.TotalHits() < 1 {
		err = errors.New(fmt.Sprintf("find 0 matches for cahtID %d", chatID))
		return 0, err
	}
	var eL LocationElastic
	return searchResult.Each(reflect.TypeOf(eL))[0].(LocationElastic).Location, nil
}
