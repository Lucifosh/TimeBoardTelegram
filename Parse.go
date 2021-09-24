package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseUrl(url string) map[string]Board {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatal("status code != 200")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	row := make([]string, 0)
	var rows [][]string
	body := doc.Find("table").Find("tbody").Find("tr")
	body.Find("tr").Each(func(indextr int, tr *goquery.Selection) {
		tr.Find("td").Each(func(indextd int, td *goquery.Selection) {
			row = append(row, td.Text())
		})
		rows = append(rows, row)
		row = nil
	})

	state := "DepFromKgd"
	m := make(map[string]Board)
	b := Board{}
	BusNumber := ""
	space := false
	lastTime := false
	for r, val := range rows {
		if r == 0 {
			continue
		}
		space = false
		for _, v := range val {
			data := strings.Split(v, "\n")
			for _, d := range data {
				s := strings.TrimSpace(d)
				//fmt.Println(len(s), " ", s)
				if len(s) == 0 {
					space = true
					continue
				}
				if space && IsTime(s) && lastTime {
					space = false
					ChangeState(&state)
				}
				if IsNumber(s) {
					tempBus := strings.Split(BusNumber, " ")
					if len(tempBus) == 2 && IsNumber(tempBus[0]) && IsNumber(tempBus[1]) {
						BusNumber = fmt.Sprintf("%v %v", tempBus[0], tempBus[1])
					} else if tempBus[0] == s {
						if len(tempBus) != 1 {
							n, err := strconv.Atoi(tempBus[1])
							if err != nil {
								log.Fatalln(err)
							}
							BusNumber = fmt.Sprintf("%v %v", tempBus[0], n+1)
						} else {
							BusNumber = fmt.Sprintf("%v 1", tempBus[0])
						}
					} else {
						BusNumber = s
					}
					//fmt.Println("Number ", BusNumber)
					state = "DepFromKgd"
					lastTime = false
					space = false
				} else if IsTime(s) {
					temp := strings.Split(s, " ")
					if len(temp) > 1 { // если время в одной строке
						if IsTime(temp[0]) && IsTime(temp[1]) { // проверка что оба время
							Adding(&b, state, temp[0])
							Adding(&b, state, temp[1])
						} else { // если была доп строка, а не время
							Adding(&b, state, s)
						}
					} else {
						Adding(&b, state, s)
					}
					lastTime = true
					space = false
				} else {
					lastTime = false
					ss := ""
					listS := strings.Split(s, " ")
					for _, v := range listS {
						ss += strings.TrimSpace(v) + " "
					}
					s = ss
					ss = ""
					for idx, v := range s {
						if v == '*' {
							if s[idx+1] == '*' {
								ss += "\n" + string(v)
							} else if idx-1 >= 0 && s[idx-1] != '*' || idx == 0 {
								ss += "\n" + string(v)
							} else {
								ss += string(v)
							}
						} else {
							ss += string(v)
						}
					}
					s = ss
					ss = ""

					index := strings.Index(s, "ПО")
					if index != -1 {
						tempS := s[:index] + "\n" + s[index:]
						b.AdditionalInformation += tempS + "\n"
					} else {
						b.AdditionalInformation += s + "\n"
					}

					state = "DepAtKgd"
					space = false
				}
			}
		}
		m[BusNumber] = b
		b = Board{}
	}
	return m
}
