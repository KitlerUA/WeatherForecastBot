package handlers

import (
	"time"

	"log"

	"github.com/KitlerUA/WeatherForecastBot/chatslocation"
	"github.com/KitlerUA/WeatherForecastBot/weather"
	"github.com/yanzay/tbot"
)

func WeatherToday(message *tbot.Message) {
	location, err := chatslocation.Get(message.ChatID)
	if err != nil {
		log.Printf("Cannot get location: %s", err)
		message.Reply("We have some problems. Please, try again later")
		return
	}
	replyString := weather.Get(time.Now(), time.Now().Add(24*time.Hour), location)
	message.Reply(replyString)
}

func WeatherTomorrow(message *tbot.Message) {
	location, err := chatslocation.Get(message.ChatID)
	if err != nil {
		log.Printf("Cannot get location: %s", err)
		message.Reply("We have some problems. Please, try again later")
		return
	}
	replyString := weather.Get(time.Now().Add(24*time.Hour), time.Now().Add(48*time.Hour), location)
	message.Reply(replyString)
}

func Weather3Days(message *tbot.Message) {
	location, err := chatslocation.Get(message.ChatID)
	if err != nil {
		log.Printf("Cannot get location: %s", err)
		message.Reply("We have some problems. Please, try again later")
		return
	}
	replyString := weather.Get(time.Now(), time.Now().Add(72*time.Hour), location)
	message.Reply(replyString)
}
