package config

import (
	"sync"
	"io/ioutil"
	"log"
	"encoding/json"
)

var config *Config
var once sync.Once

type Config struct {
	BotToken string `json:"bot_token"`
}

func Get() Config{
	once.Do(loadConfigFile)
	return *config
}

func loadConfigFile (){
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panicf("Cannot load configuration from file: %v", err)
	}
	config = &Config{}
	if err = json.Unmarshal(data, &config); err!= nil {
		log.Panicf("Corrupted data in config.json : %v", err)
	}
}