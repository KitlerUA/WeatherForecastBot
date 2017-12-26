package handlers

import (
	"github.com/yanzay/tbot"
)

func Start(message *tbot.Message) {
	message.ReplyKeyboard("Welcome!\nPlease choose your location", [][]string{{"Lviv"}})
}
