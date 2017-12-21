package handlers

import (
	"time"

	"github.com/KitlerUA/WeatherForecastBot/weather"
	"github.com/yanzay/tbot"
)

func Start(message *tbot.Message) {
	message.ReplyKeyboard("Welcome!\nPlease choose your location", [][]string{{"Lviv"}})
}

func WeatherToday(message *tbot.Message) {
	message.Reply(weather.Get(time.Now(), time.Now().Add(24*time.Hour)))
}
