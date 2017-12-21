package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var config *Config
var once sync.Once

type Config struct {
	BotToken              string `json:"bot_token"`
	ChatDefaultLocation   string `json:"chat_default_location"`
	ElasticAddress        string `json:"elastic_address"`
	LocationsCodeFileName string `json:"locations_code_file_name"`
}

func Get() Config {
	once.Do(loadConfigFile)
	return *config
}

func loadConfigFile() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Panicf("Cannot load configuration from file: %v", err)
	}
	config = &Config{}
	if err = json.Unmarshal(data, &config); err != nil {
		log.Panicf("Corrupted data in config.json : %v", err)
	}
}
