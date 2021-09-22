package main

import (
	"strconv"
	"strings"
	"unicode"
)

func TryToFindTime(s string) bool {
	t := strings.Split(s, " ")
	if len(t) == 1 {
		return false
	}
	count := 0
	for _, v := range t {
		temp := strings.TrimSpace(v)
		line := strings.Split(temp, " ")
		for _, val := range line {
			if IsTime(val) {
				//fmt.Println("TryToFind ", s)
				count++
			}
		}
	}
	return count == 2
}
func IsNumber(s string) bool {
	digits := 0
	if strings.Contains(s, "Калининград") || strings.Contains(s, "~") {
		return false
	}
	newS := strings.TrimSpace(s)
	for _, v := range newS {
		if unicode.IsDigit(v) {
			digits++
		}
	}
	return digits == 3 || digits == 6
}
func IsTime(s string) bool {
	digits := 0
	tempS := strings.TrimSpace(s)
	ss := strings.Split(tempS, " ")
	chz := false
	skb := false
	if tempS == "# 06-20" {
		return true
	}
	if len(ss) > 1 {
		for _, v := range ss {
			if strings.Contains(v, "ч/з") {
				chz = true
			}
		}
	} else {
		if strings.ContainsAny(tempS, "(") {
			skb = true
		}
	}
	for _, v := range ss[0] {
		if unicode.IsDigit(v) {
			digits++
		} else if unicode.IsLetter(v) && (!chz && !skb) {
			return false
		}
	}
	if digits == 4 && strings.Contains(s, "-") {
		return true
	}
	return false
}
func TimeCompare(big, small string) bool {
	b, err := strconv.Atoi(big)
	if err != nil {
		panic(err)
	}
	s, err := strconv.Atoi(small)
	if err != nil {
		panic(err)
	}
	if b == 0 && s > 10 {
		return !(b >= s)
	} else if s == 0 {
		return !(b >= s)
	}
	return b >= s
}
func Adding(board *Board, state, v string) {
	//fmt.Println("ADD ", state, " ", v)
	switch state {
	case "DepFromKgd":
		board.DepartureTimeFromStation = append(board.DepartureTimeFromStation, strings.TrimSpace(v))
	case "ArAtKgd":
		board.ArrivalTimeAtDestination = append(board.ArrivalTimeAtDestination, strings.TrimSpace(v))
	case "DepAtKgd":
		board.DepartureTimeAtStation = append(board.DepartureTimeAtStation, strings.TrimSpace(v))
	case "ArFromKgd":
		board.ArrivalTimeFromDestination = append(board.ArrivalTimeFromDestination, strings.TrimSpace(v))
	}
}
func ChangeState(state *string) {
	if *state == "DepFromKgd" {
		*state = "ArAtKgd"
	} else if *state == "ArAtKgd" {
		*state = "DepAtKgd"
	} else if *state == "DepAtKgd" {
		*state = "ArFromKgd"
	} else if *state == "ArFromKgd" {
		*state = "DepFromKgd"
	}
}
