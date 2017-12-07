package handlers

import "github.com/yanzay/tbot"

func Start(message *tbot.Message){
	message.ReplyKeyboard("Second welcome!", [][]string{{"1", "2"}, {"second row 1", "4"}})
	message
}
