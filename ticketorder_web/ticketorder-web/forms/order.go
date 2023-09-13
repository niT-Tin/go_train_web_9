package forms

type OrderForm struct {
	DailyTrainTicketId int64  `json:"dailyTrainTicketId"`
	Date               string `json:"date"`
	ImageCode          string `json:"imageCode"`
	ImageCodeId        string `json:"imageCodeId"`
	LineNumber         int32  `json:"lineNumber"`
	Start              string `json:"start"`
	End                string `json:"end"`
	Tickets            []struct {
		PassengerId     int32  `json:"passengerId"`
		PassengerName   string `json:"passengerName"`
		PassengerType   int32  `json:"passengerType"`
		PassengerIdCard string `json:"passengerIdCard"`
		Seat            string `json:"seat"`
		SeatTypeCode    string `json:"seatTypeCode"`
	} `json:"tickets"`
}
