package main

import (
	"fmt"
	"log"
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TelegramBot() {
	bot, err := tgbotapi.NewBotAPI(GetToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	state := 0
	answers := Answer{}

	updates, _ := bot.GetUpdatesChan(u)
	for {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
				text := update.Message.Text
				msg := Say(&state, text, update.Message.Chat.ID, &answers)
				bot.Send(msg)
				fmt.Printf("%v %v %v\n%v\n%v\n", answers.Transport, answers.Direction, answers.searchBy, answers.Links, answers.Needs)
			}
		}
	}
}
