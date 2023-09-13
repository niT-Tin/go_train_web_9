package forms

type OrderForm struct {
	DailyTrainTicketId int64  `json:"dailyTrainTicketId"`
	Date               string `json:"date"`
	ImageCode          string `json:"imageCode"`
	ImageCodeId        string `json:"imageCodeId"`
	UserID             int64  `json:"userId"`
	LineNumber         int32  `json:"lineNumber"`
	Start              string `json:"start"`
	TrainCode          string `json:"trainCode"`
	StartTime          string `json:"startTime"`
	EndTime            string `json:"endTime"`
	End                string `json:"end"`
	Tickets            []struct {
		PassengerId     int32  `json:"passengerId"`
		PassengerName   string `json:"passengerName"`
		PassengerType   int32  `json:"passengerType"`
		UserID          int64  `json:"userId"`
		PassengerIdCard string `json:"passengerIdCard"`
		Seat            string `json:"seat"`
		SeatTypeCode    string `json:"seatTypeCode"`
	} `json:"tickets"`
}
