package handlers

import (
	"context"
	"fmt"
	"reflect"

	"strings"

	"github.com/KitlerUA/WeatherForecastBot/chatslocation"
	"github.com/KitlerUA/WeatherForecastBot/db"
	"github.com/olivere/elastic"
	"github.com/yanzay/tbot"
)

func CustomInput(message *tbot.Message) {
	locationFound := false
	query := elastic.NewTermQuery("name", strings.ToLower(message.Text()))
	ctx := context.Background()
	searchResult, err := db.Get().Search().Index("citbot").Type("city").Query(query).Do(ctx)
	if err != nil {
		fmt.Printf("Cannot search for city: %s", err)
	}
	if searchResult.TotalHits() > 0 {
		var cE chatslocation.CityElastic
		id := searchResult.Each(reflect.TypeOf(cE))[0].(chatslocation.CityElastic).Id
		if err := chatslocation.AddOrUpdate(message.ChatID, id); err == nil {
			locationFound = true
		}
	}

	if locationFound {
		message.ReplyKeyboard("Choose period of time", [][]string{{"Today", "Tomorrow"}, {"3 days"}})
		return
	}
	message.Reply("Cannot find your location")
}
