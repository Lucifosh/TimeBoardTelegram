package main

type Board struct {
	DepartureTimeFromStation   []string //отправление из станции
	DepartureTimeAtStation     []string //отправление на станцию
	ArrivalTimeFromDestination []string //прибытие из кон. пункта
	ArrivalTimeAtDestination   []string //прибытие на кон. пункт
	AdditionalInformation      string   // разная доп инфа
}

func main() {
	// TestParse()
	// q := 0
	// fmt.Scan(&q)
	//Temp()
	TelegramBot()

}
