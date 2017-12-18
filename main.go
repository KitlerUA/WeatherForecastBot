package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/KitlerUA/WeatherForecastBot/chatslocation"
	"github.com/KitlerUA/WeatherForecastBot/config"
	"github.com/KitlerUA/WeatherForecastBot/handlers"
	"github.com/yanzay/tbot"
)

func main() {
	log.Print("Loading data...")
	chatslocation.DefaultLocationByChatID = make(map[int64]int)
	data, err := ioutil.ReadFile(config.Get().ChatDefaultLocation)
	if err != nil {
		log.Panicf("Cannot load locations from file: %v", err)
	}
	if err = json.Unmarshal(data, &chatslocation.DefaultLocationByChatID); err != nil {
		log.Print("Corrupted data in locations file. Data was`n load", err)
	}
	//save file before close
	defer func() {
		data, err = json.Marshal(chatslocation.DefaultLocationByChatID)
		if err != nil {
			if err = ioutil.WriteFile(config.Get().ChatDefaultLocation, data, 0666); err != nil {
				log.Printf("Cannot save chats location in %s", config.Get().ChatDefaultLocation)
			}
		}
	}()
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
	bot.HandleFunc("{custom}", handlers.CustomInput)

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
