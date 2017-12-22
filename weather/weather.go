package weather

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"context"

	"reflect"

	"strconv"

	"fmt"

	"github.com/KitlerUA/WeatherForecastBot/config"
	"github.com/KitlerUA/WeatherForecastBot/db"
	"github.com/olivere/elastic"
	"github.com/yanzay/log"
)

const indexMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"info":{
			"properties":{
				"location":{
					"type":"long"
				},
				"temp":{
					"type":"float"
				},
				"humidity":{
					"type":"float"
				},
				"description":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"dttxt":{
					"type":"long"
				}
			}
		}
	}
}`

type InfoElastic struct {
	Location    int
	Temp        float32
	Humidity    float32
	Description string
	DtTxt       int64
}

type Info struct {
	Main struct {
		Temp     float32 `json:"temp"`
		Humidity float32 `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	DtTxt string `json:"dt_txt"`
}

type InfoList struct {
	List []Info `json:"list"`
}

func Get(startDate, endDate time.Time, location int) string {

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
		//createIndex, err := db.Get().CreateIndex("wetbot").BodyString(indexMapping).Do(ctx)
		createIndex, err := db.Get().CreateIndex("wetbot").Do(ctx)
		db.Get().PutMapping()
		if err != nil {
			log.Fatalf("Cannot create index: %s ", err)
		}
		if !createIndex.Acknowledged {
			log.Fatal("Not acknowledge")
		}
	}
	query := elastic.NewRangeQuery("DtTxt").Gte(startDate.Unix()).Lte(endDate.Unix())
	queryLocation := elastic.NewTermQuery("Location", location)
	searchResult, err := db.Get().Search().Index("wetbot").Type("info").Query(query).Query(queryLocation).Sort("DtTxt", true).Do(ctx)
	if err != nil {
		log.Fatalf("Cannot search: %s", err)
	}
	if searchResult.Hits.TotalHits == 0 {
		weather, err := getWeatherFromOpenMap(location)
		if err != nil {
			log.Fatalf("Cannot get weather from OpenMap: %s", err)
		}
		for _, info := range weather.List {

			eInfo := infoToElasticInfo(info)
			_, err := db.Get().Index().Index("wetbot").Type("info").BodyJson(eInfo).Do(ctx)
			if err != nil {
				log.Fatalf("Cannot put %v: %s", eInfo, err)
			}
		}
		_, err = db.Get().Flush().Index("wetbot").Do(ctx)
		if err != nil {
			log.Printf("Cannot flush index: %s")
		}
	}
	//log.Printf("Query %s", query)
	searchResult, err = db.Get().Search().Index("wetbot").Type("info").Query(query).Query(queryLocation).Sort("DtTxt", true).Do(ctx)
	if err != nil {
		log.Fatalf("Cannot search: %s", err)
	}

	replyString := ""
	var info InfoElastic
	for _, item := range searchResult.Each(reflect.TypeOf(info)) {
		if t, ok := item.(InfoElastic); ok {
			replyString += time.Unix(t.DtTxt, 0).Format("2006-01-02 15:04:05") + " " + t.Description + " " + strconv.Itoa(int(t.Temp)) + "°C " + strconv.Itoa(int(t.Humidity)) + "%\n"
		}
	}
	return replyString
}

func getWeatherFromOpenMap(location int) (InfoList, error) {
	weatherClient := http.Client{
		Timeout: 5 * time.Second,
	}
	weather := InfoList{}
	requestString := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?id=%s&units=metric&appid=9ebbdc484f058b6e91cba224d761fea2", strconv.FormatInt(int64(location), 10))
	req, err := http.NewRequest(http.MethodGet, requestString, nil)
	if err != nil {
		return weather, err
	}
	res, err := weatherClient.Do(req)
	if err != nil {
		return weather, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return weather, err
	}
	err = json.Unmarshal(body, &weather)
	return weather, err

}

func infoToElasticInfo(info Info) InfoElastic {
	dt, _ := time.Parse("2006-01-02 15:04:05", info.DtTxt)
	res := InfoElastic{
		Location:    702550,
		Temp:        info.Main.Temp,
		Humidity:    info.Main.Humidity,
		Description: info.Weather[0].Description,
		DtTxt:       dt.Unix(),
	}
	return res
}
