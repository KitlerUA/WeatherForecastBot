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

type InfoElastic struct {
	Location    int
	Temp        float32
	Humidity    float32
	Description string
	DtTxt       int64
	IconID      string
}

type Info struct {
	Main struct {
		Temp     float32 `json:"temp"`
		Humidity float32 `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
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
	queryDate := elastic.NewRangeQuery("DtTxt").Gte(startDate.Unix()).Lte(endDate.Unix())

	queryLocation := elastic.NewTermQuery("Location", location)
	query := elastic.NewBoolQuery().Must(queryDate, queryLocation)
	searchResult, err := db.Get().Search().Index("wetbot").Type("info").From(0).Size(24).Query(query).Sort("DtTxt", true).Do(ctx)
	if err != nil {
		log.Printf("Cannot search forecast (first): %s", err)
	}
	if err != nil || searchResult.Hits.TotalHits == 0 {
		weather, err := getWeatherFromOpenMap(location)
		if err != nil {
			log.Printf("Cannot get weather from OpenMap: %s", err)
			return "Please, try again later"
		}
		for _, info := range weather.List {

			eInfo := infoToElasticInfo(info, location)
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
	searchResult, err = db.Get().Search().Index("wetbot").Type("info").From(0).Size(24).Query(query).Sort("DtTxt", true).Do(ctx)
	if err != nil {
		log.Fatalf("Cannot search forecast (second): %s", err)
	}

	replyString := ""
	var info InfoElastic
	var currentDate time.Time
	for i, item := range searchResult.Each(reflect.TypeOf(info)) {
		if t, ok := item.(InfoElastic); ok {
			icon, ok := icons[t.IconID]
			if !ok {
				log.Printf("Cannot find icon for %s", t.IconID)
			}
			if i == 0 || currentDate.Day() < time.Unix(t.DtTxt, 0).Day() {
				currentDate = time.Unix(t.DtTxt, 0)
				if i != 0 {
					replyString += "\n"
				}
				replyString += currentDate.Weekday().String() + ", " + currentDate.Format("2006-01-02") + "\n"
			}
			replyString += time.Unix(t.DtTxt, 0).Format("2006-01-02 15:04:05")[11:] + " " + icon + " " + strconv.Itoa(int(t.Temp)) + "°C " + strconv.Itoa(int(t.Humidity)) + "%\n"
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
	log.Printf("Received %d , data from server", len(weather.List))
	return weather, err

}

func infoToElasticInfo(info Info, location int) InfoElastic {
	dt, _ := time.Parse("2006-01-02 15:04:05", info.DtTxt)
	res := InfoElastic{
		Location:    location,
		Temp:        info.Main.Temp,
		Humidity:    info.Main.Humidity,
		Description: info.Weather[0].Description,
		DtTxt:       dt.Unix(),
		IconID:      info.Weather[0].Icon,
	}
	return res
}
