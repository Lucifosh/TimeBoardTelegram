package main

import (
	"strings"
)

func SearchKey(s string, sub string) bool {
	temp := strings.Split(s, " ")
	for _, v := range temp {
		if v == sub {
			return true
		}
	}
	return false
}
func SearchBus(m map[string]string, number string) (map[string]Board, []string) {
	num := strings.ToUpper(number)
	found := make(map[string]Board)
	urls := make([]string, 0)
	for _, v := range m {
		pu := ParseUrl(v)
		for bus, val := range pu {
			if strings.Contains(bus, num) {
				found[bus] = val
				urls = append(urls, v)
			} else if len(found) > 0 {
				return found, urls
			}
		}
	}
	return found, urls
}
