package main

import (
	"log"

	"github.com/KitlerUA/WeatherForecastBot/chatslocation"
	"github.com/KitlerUA/WeatherForecastBot/config"
	"github.com/KitlerUA/WeatherForecastBot/handlers"
	"github.com/KitlerUA/WeatherForecastBot/indexbuilder"
	"github.com/KitlerUA/WeatherForecastBot/weather"
	"github.com/yanzay/tbot"
)

func main() {
	log.Println("Creating indices")
	err := indexbuilder.BuildIndices("wetbot", "locbot", "citbot")
	if err != nil {
		log.Panicf("Cannot build indices: %s", err)
	}
	log.Println("Loading cities")
	chatslocation.LoadCityList()
	log.Print("Creating server...")
	bot, err := tbot.NewServer(config.Get().BotToken)
	if err != nil {
		log.Fatalf("Cannot create bot with given token: %v", err)
	}
	log.Print("Server created")
	log.Print("Adding middleware...")
	bot.AddMiddleware(logMid)
	log.Print("Middleware added")

	bot.HandleFunc("/start", handlers.Start)
	bot.HandleFunc("Today", handlers.WeatherToday)
	bot.HandleFunc("Tomorrow", handlers.WeatherTomorrow)
	bot.HandleFunc("3 days", handlers.Weather3Days)
	bot.HandleFunc("{custom}", handlers.CustomInput)
	log.Printf("Start weather updater...")
	go weather.Update()
	log.Print("Start listen and serve...")
	err = bot.ListenAndServe()
	log.Fatal(err)
}

func logMid(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		log.Print("From ", m.From.UserName, " ", m.From.FirstName, " ", m.From.LastName, " : ", m.Text())
		f(m)
	}
}
