package handlers

import (
	"time"

	"github.com/KitlerUA/WeatherForecastBot/chatslocation"
	"github.com/KitlerUA/WeatherForecastBot/weather"
	"github.com/yanzay/log"
	"github.com/yanzay/tbot"
)

func Start(message *tbot.Message) {
	message.ReplyKeyboard("Welcome!\nPlease choose your location", [][]string{{"Lviv"}})
}

func WeatherToday(message *tbot.Message) {
	location, err := chatslocation.Get(message.ChatID)
	if err != nil {
		log.Printf("Cannot get location: %s", err)
		message.Reply("We have some problems. Please, try again later")
		return
	}
	replyString := weather.Get(time.Now(), time.Now().Add(24*time.Hour), location)
	message.Reply(replyString)
	log.Printf("Reply %s", replyString)
}
