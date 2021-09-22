package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetKeyboard(k int) tgbotapi.ReplyKeyboardMarkup {

	m := map[int]tgbotapi.ReplyKeyboardMarkup{
		0: tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Автобус"),
				//tgbotapi.NewKeyboardButton("Поезд"),
				//tgbotapi.NewKeyboardButton("Электричка"),
			),
		),
		1: tgbotapi.NewOneTimeReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("В сторону Калининграда"),
				tgbotapi.NewKeyboardButton("Из Калининграда"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		),
		2: tgbotapi.NewOneTimeReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Назад"),
			),
		),
	}
	n := k
	if n > 2 {
		n = 2
	}
	return m[n]
}
