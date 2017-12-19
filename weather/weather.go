package weather

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"context"

	"github.com/KitlerUA/WeatherForecastBot/config"
	"github.com/KitlerUA/WeatherForecastBot/db"
	"github.com/yanzay/log"
)

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

func Get(startDate, endDate time.Time) string {

	version, err := db.Get().ElasticsearchVersion(config.Get().ElasticAddress)
	if err != nil {
		log.Fatalf("Cannot ping elastic %s", err)
	}

	log.Printf("Elasticsearch version %s", version)
	ctx := context.Background()
	exists, err := db.Get().IndexExists("twitter").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		createIndex, err := db.Get().CreateIndex("wetbot").Do(ctx)
		if err != nil {
			log.Fatalf("Cannot create index: %s ", err)
		}
		if !createIndex.Acknowledged {
			log.Fatal("Not acknowledge")
		}
	}
	if get, err := db.Get().Index("wetbot").Type("info").
	weather, err := getWeatherFromOpenMap()
	if err != nil {
		log.Fatalf("Cannot get weather from OpenMap: %s", err)
	}
	for _, info := range weather.List {
		unmarshaled, err := json.Marshal(&info)
		if err != nil {
			log.Fatalf("Cannot marshal info: %s", err)
		}

		_, err = db.Get().Index().Index("wetbot").Type("wetinfo").BodyJson(unmarshaled).Do(ctx)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	//get forecast from server, if we have no local data
	/*weatherClient := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, "http://api.openweathermap.org/data/2.5/forecast?id=524901&appid=9ebbdc484f058b6e91cba224d761fea2", nil)
	if err != nil {
		return err.Error()
	}
	res, err := weatherClient.Do(req)
	if err != nil {
		return err.Error()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	weather := InfoList{}
	if err = json.Unmarshal(body, &weather); err != nil {
		return err.Error()
	}
	replyString := ""
	for _, i := range weather.List {
		replyString += i.DtTxt + " "
		replyString += i.Weather[0].Description + "; "
		replyString += "Temperature " + strconv.Itoa(int(i.Main.Temp)) + "°C ; Humidity " + strconv.Itoa(int(i.Main.Humidity)) + "%\n"
	}

	return replyString*/
	return version
}

func getWeatherFromOpenMap() (InfoList, error) {
	weatherClient := http.Client{
		Timeout: 5 * time.Second,
	}
	weather := InfoList{}
	req, err := http.NewRequest(http.MethodGet, "http://api.openweathermap.org/data/2.5/forecast?id=524901&appid=9ebbdc484f058b6e91cba224d761fea2", nil)
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
