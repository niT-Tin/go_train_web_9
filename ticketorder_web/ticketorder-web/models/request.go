package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type PassengerIds []string

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

type PassengerType int64

const (
	// 乘客类型
	PassengerTypeAdult PassengerType = 0 // 成人
	PassengerTypeChild               = 1 // 儿童
)

type Passenger struct {
	Name     string        `json:"name"`
	UserID   int32         `json:"user_id"`
	IdCard   string        `json:"id_card"`
	Type     PassengerType `json:"type"`
	SeatType string        `json:"seat_type"`
	Seat     string        `json:"seat"`
}

type Order struct {
	UserID       int32     `gorm:"column:user_id;type:int;not null;"`
	TrainID      int32     `gorm:"column:train_id;type:int;not null;"`
	TrainCode    string    `gorm:"column:train_code;type:varchar(20);not null;"`
	StartStation string    `gorm:"column:start_station;type:varchar(20);not null;"`
	EndStation   string    `gorm:"column:end_station;type:varchar(20);not null;"`
	StartTime    time.Time `gorm:"column:start_time;type:datetime;not null;"`
	EndTime      time.Time `gorm:"column:end_time;type:datetime;not null;"`
	SeatType     string    `gorm:"column:seat_type;type:varchar(20);not null;"`
	SeatNumber   string    `gorm:"column:seat_number;type:varchar(20);not null;"`
	Pirce        float32   `gorm:"column:price;type:float;not null;"`
	OrderSn      string    `gorm:"column:order_sn;type:varchar(20);not null;unique_index;"`
	Row          string    `gorm:"column:row;type:varchar(20);not null;"`
	Colum        string    `gorm:"column:colum;type:varchar(20);not null;"`
	PassengerIds PassengerIds
}

func (p *PassengerIds) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &p)
}

func (p PassengerIds) Value() (driver.Value, error) {
	return json.Marshal(p)
}
