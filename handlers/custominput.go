package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/KitlerUA/WeatherForecastBot/chatslocation"
	"github.com/KitlerUA/WeatherForecastBot/config"
	"github.com/yanzay/tbot"
)

func CustomInput(message *tbot.Message) {
	locationFound := false
	if _, ok := chatslocation.DefaultLocationByChatID[message.ChatID]; !ok {
		for i := range Locations {
			if message.Text() == i {
				chatslocation.DefaultLocationByChatID[message.ChatID] = Locations[i]
				message.Replyf("Your location set to %s", i)
				data, err := json.Marshal(chatslocation.DefaultLocationByChatID)
				if err != nil {

				}
				if err = ioutil.WriteFile(config.Get().ChatDefaultLocation, data, 0666); err != nil {
					log.Printf("Cannot save chats location in %s", config.Get().ChatDefaultLocation)
				}
				locationFound = true
			}
		}
	} else {
		if v, ok := Periods[message.Text()]; ok {
			v(message)
		}
		locationFound = true
	}
	if locationFound {
		message.ReplyKeyboard("Choose period of time", [][]string{{"Today", "Tomorrow"}, {"3 days"}})
		return
	}
	message.Reply("Cannot find your location")
}
