package main

import (
	"github.com/yanzay/tbot"
	"github.com/KitlerUA/WeatherForecastBot/config"
	"log"
	"github.com/KitlerUA/WeatherForecastBot/handlers"
)

func main(){
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

	log.Print("Start listen and serve...")
	err = bot.ListenAndServe()
	log.Fatal(err)
}

func logMid(f tbot.HandlerFunction) tbot.HandlerFunction {
	return func(m *tbot.Message){
		log.Print("From ", m.From.UserName, " ", m.From.FirstName, " ", m.From.LastName, " : ", m.Text())
		f(m)
	}
}