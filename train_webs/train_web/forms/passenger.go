package forms

type PassengerForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=2,max=20"`
	IdCard   string `form:"id_card" json:"id_card" binding:"required,len=18"`
	Type     int64  `form:"type" json:"type" binding:"required,oneof=adult child"`
	UserID   int32  `form:"user_id" json:"user_id" binding:"required"`
	Seat     string `form:"seat" json:"seat" binding:"required"`
	SeatType string `form:"seat_type" json:"seat_type" binding:"required"`
}
