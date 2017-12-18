package weather

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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
	weatherClient := http.Client{
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

	return replyString
}
