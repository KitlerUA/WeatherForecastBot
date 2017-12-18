package handlers

import "github.com/yanzay/tbot"

var Locations = map[string]int{
	"L`viv, Ukraine": 702550,
}

var Periods = map[string]tbot.HandlerFunction{
	"Today": WeatherToday,
}
