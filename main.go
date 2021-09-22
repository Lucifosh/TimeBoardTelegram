package main

import (
	"io/ioutil"
	"strings"
)

type Board struct {
	DepartureTimeFromStation   []string //отправление из станции
	DepartureTimeAtStation     []string //отправление на станцию
	ArrivalTimeFromDestination []string //прибытие из кон. пункта
	ArrivalTimeAtDestination   []string //прибытие на кон. пункт
	AdditionalInformation      string   // разная доп инфа
}

func addMap() map[string]string {

	m := make(map[string]string)

	byteString, err := ioutil.ReadFile("links.txt")
	if err != nil {
		panic(err)
	}
	text := string(byteString)
	if len(text) == 0 {
		return m
	}
	text = strings.TrimSpace(text)
	data := strings.Split(text, "\n")
	for _, v := range data {
		temp := strings.Split(v, " ")
		if len(temp) == 2 {
			cut := temp[1][:len(temp[1])]
			for _, c := range cut {
				if c == '\r' {
					cut = cut[:len(cut)-1]
				}
			}
			m[temp[0]] = cut
		} else {
			key := ""
			for i := 0; i < len(temp)-1; i++ {
				key += temp[i] + " "
			}
			cut := temp[len(temp)-1][:len(temp[len(temp)-1])]
			for _, c := range cut {
				if c == '\r' {
					cut = cut[:len(cut)-1]
				}
			}
			m[key] = cut
		}
	}

	return m
}
func main() {
	// TestParse()
	// q := 0
	// fmt.Scan(&q)
	TelegramBot()

}
