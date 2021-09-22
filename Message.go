package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Answer struct {
	Transport string   //Вид транспорта
	Direction string   //Направление
	searchBy  string   //поиск по номеру или городу
	Links     []string //Ссылки для ответа
	Needs     []string //Города или номера транспорта
}

func Say(state *int, text string, id int64, answer *Answer) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(id, "")
	if text == "Назад" {
		*state = 0
		msg = Say(state, "", id, answer)
	}
	if *state == 0 { // выбор транспорта
		answer.Links = []string{}
		answer.Needs = []string{}
		if text == "Автобус" {
			*state++
			answer.Transport = "а"
			msg = tgbotapi.NewMessage(id, "Выберите направление")
			msg.ReplyMarkup = GetKeyboard(*state)
		} else {
			msg.ReplyMarkup = GetKeyboard(*state)
			msg = tgbotapi.NewMessage(id, "Выберите транспорт")
			msg.ReplyMarkup = GetKeyboard(*state)
		}
	} else if *state == 1 { // выбор направления
		if text == "В сторону Калининграда" {
			msg = tgbotapi.NewMessage(id, "Укажите, откуда едет транспорт или номер транспорта")
			*state++
			msg.ReplyMarkup = GetKeyboard(*state)
			answer.Direction = "в"
		} else if text == "Из Калининграда" {
			msg = tgbotapi.NewMessage(id, "Укажите, куда едет транспорт или номер транспорта")
			*state++
			msg.ReplyMarkup = GetKeyboard(*state)
			answer.Direction = "и"
		} else {
			msg = tgbotapi.NewMessage(id, "Выберите направление")
			msg.ReplyMarkup = GetKeyboard(*state)
		}
	} else if *state == 2 { // поиск города или номера
		m := addMap()
		if IsNumber(text) { // поиск по номеру
			Buss, urls := SearchBus(m, text)
			if len(Buss) == 0 {
				s := fmt.Sprintf("Не нашел номера \"%v\". Попробуйте снова", text)
				msg = tgbotapi.NewMessage(id, s)

			} else if len(Buss) == 1 { // если найден всего 1 номер то сразу вперед
				answer.Links = append(answer.Links, urls[0])
				for bus := range Buss {
					answer.Needs = append(answer.Needs, bus)
				}
				answer.searchBy = "н"

				*state++
				msg = Say(state, "", id, answer)
			} else { // если найдено больше 1
				k := 1
				message := ""
				for bus, t := range Buss {
					b := strings.Split(bus, " ")[0]
					message += fmt.Sprintf("%v. %v\n%v\n", k, b, t.AdditionalInformation)
					k++
				}
				message += "\nУкажите цифру"
				msg = tgbotapi.NewMessage(id, message)

				for bus := range Buss {
					answer.Needs = append(answer.Needs, bus)
				}
				answer.Links = append(answer.Links, urls[len(urls)-1])
				answer.searchBy = "н"

				*state++
			}
		} else { // поиск по городу
			citys := make([]string, 0)
			urls := make([]string, 0)
			for key, url := range m { // поиск городов
				if SearchKey(key, strings.ToLower(text)) {
					citys = append(citys, key)
					urls = append(urls, url)
				}
			}
			if len(citys) == 0 {
				s := fmt.Sprintf("Не нашел города \"%v\". Попробуйте снова", text)
				msg = tgbotapi.NewMessage(id, s)
			} else if len(citys) == 1 { // если найден 1 город то сразу вперед
				answer.Links = append(answer.Links, urls[0])
				answer.Needs = append(answer.Needs, citys...)
				answer.searchBy = "г"

				*state++
				msg = Say(state, "", id, answer)
			} else { // если найдено больше
				message := ""
				for idx, v := range citys {
					message += fmt.Sprintf("%v. %v\n", idx+1, strings.ToUpper(v))
				}
				message += "\nУкажите цифру"
				msg = tgbotapi.NewMessage(id, message)

				kk := -1
				for _, v := range urls {
					answer.Links = append(answer.Links, v)
					kk--
				}
				answer.searchBy = "г"

				*state++
			}
		}
	} else if *state == 3 { // вывод расписания
		n := 0
		var err error = nil
		if len(answer.Needs) == 1 {
			n = 1
		} else {
			n, err = strconv.Atoi(text)
		}

		if err != nil {
			msg = tgbotapi.NewMessage(id, "Неправильный формат. Попробуйте повторить")
		} else {
			n--
			fmt.Printf("***\n%v is correct\n***\n", text)
			if answer.Transport == "а" { //автобус
				if answer.searchBy == "н" { // поиск по номеру
					if len(answer.Needs) <= n || n < 0 {
						msg = tgbotapi.NewMessage(id, "Такой цифры нет. Попробуйте повторить")
					} else {
						message := ""
						p := ParseUrl(answer.Links[0])
						checkBus := strings.Split(answer.Needs[n], " ")
						if len(checkBus) == 1 {
							message += "Рейс №" + checkBus[0] + "\n"
						} else {
							if IsNumber(checkBus[0]) && IsNumber(checkBus[1]) {
								message += "Рейс №" + checkBus[0] + " " + "Рейс №" + checkBus[1]
							} else {
								message += "Рейс №" + checkBus[0] + "\n"
							}
						}
						message += p[answer.Needs[n]].AdditionalInformation + "\n"
						if answer.Direction == "в" { // в город
							l := len(p[answer.Needs[n]].DepartureTimeAtStation)
							for i := 0; i < l; i++ {
								message += p[answer.Needs[n]].DepartureTimeAtStation[i] + "\t\t"
								message += p[answer.Needs[n]].ArrivalTimeFromDestination[i] + "\n"
							}
							message += "\n"
							msg = tgbotapi.NewMessage(id, message)
							*state = 0
						} else if answer.Direction == "и" { //из города
							l := len(p[answer.Needs[n]].DepartureTimeFromStation)
							for i := 0; i < l; i++ {
								message += p[answer.Needs[n]].DepartureTimeFromStation[i] + "\t\t"
								message += p[answer.Needs[n]].ArrivalTimeAtDestination[i] + "\n"
							}
							message += "\n"
							msg = tgbotapi.NewMessage(id, message)
							*state = 0
						}
					}

				} else if answer.searchBy == "г" { // поиск по городу
					if len(answer.Links) <= n || n < 0 {
						msg = tgbotapi.NewMessage(id, "Такой цифры нет. Попробуйте повторить")
					} else {

						message := ""
						p := ParseUrl(answer.Links[n])

						if answer.Direction == "в" { // в город
							for bus, time := range p {
								checkBus := strings.Split(bus, " ")
								if len(checkBus) == 1 {
									message += "Рейс №" + checkBus[0] + "\n"
								} else {
									if IsNumber(checkBus[0]) && IsNumber(checkBus[1]) {
										message += "Рейс №" + checkBus[0] + " " + "Рейс №" + checkBus[1]
									} else {
										message += "Рейс №" + checkBus[0] + "\n"
									}
								}
								message += time.AdditionalInformation + "\n"
								l := len(time.DepartureTimeAtStation)
								for i := 0; i < l; i++ {
									message += time.DepartureTimeAtStation[i] + "\t\t"
									message += time.ArrivalTimeFromDestination[i] + "\n"
								}
								message += "\n\n"
								msg = tgbotapi.NewMessage(id, message)
								*state = 0
								msg.ReplyMarkup = GetKeyboard(*state)
							}
						} else if answer.Direction == "и" { // из города
							for bus, time := range p {
								checkBus := strings.Split(bus, " ")
								if len(checkBus) == 1 {
									message += "Рейс №" + checkBus[0] + "\n"
								} else {
									if IsNumber(checkBus[0]) && IsNumber(checkBus[1]) {
										message += "Рейс №" + checkBus[0] + " " + "Рейс №" + checkBus[1]
									} else {
										message += "Рейс №" + checkBus[0] + "\n"
									}
								}
								message += time.AdditionalInformation + "\n"
								l := len(time.DepartureTimeFromStation)
								for i := 0; i < l; i++ {
									message += time.DepartureTimeFromStation[i] + "\t\t"
									message += time.ArrivalTimeAtDestination[i] + "\n"
								}
								message += "\n\n"
								msg = tgbotapi.NewMessage(id, message)
								*state = 0
								msg.ReplyMarkup = GetKeyboard(*state)
							}
						}
					}
				}
			}
		}
	}
	return msg
}
