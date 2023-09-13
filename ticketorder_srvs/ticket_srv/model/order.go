package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type PassengerIds []string

type BaseModel struct {
	ID        int32          `gorm:"primary_key;"`
	CreatedAt time.Time      `gorm:"column:add_time;"`
	UpdatedAt time.Time      `gorm:"column:update_time;"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_time;"`
	IsDelete  bool           `gorm:"column:is_delete;"`
}

type Order struct {
	BaseModel
	UserID       int32     `gorm:"column:user_id;type:int;not null;" json:"user_id"`
	TrainID      int32     `gorm:"column:train_id;type:int;not null;" json:"train_id"`
	TrainCode    string    `gorm:"column:train_code;type:varchar(20);not null;" json:"train_code"`
	StartStation string    `gorm:"column:start_station;type:varchar(20);not null;" json:"start_station"`
	EndStation   string    `gorm:"column:end_station;type:varchar(20);not null;" json:"end_station"`
	StartTime    time.Time `gorm:"column:start_time;type:datetime;not null;" json:"start_time"`
	EndTime      time.Time `gorm:"column:end_time;type:datetime;not null;" json:"end_time"`
	SeatType     string    `gorm:"column:seat_type;type:varchar(20);not null;" json:"seat_type"`
	SeatNumber   string    `gorm:"column:seat_number;type:varchar(20);not null;" json:"seat_number"`
	Price        float32   `gorm:"column:price;type:float;not null;" json:"price"`
	OrderSn      string    `gorm:"column:order_sn;type:varchar(20);not null;unique_index;" json:"order_sn"`
	Row          string    `gorm:"column:row;type:varchar(20);not null;" json:"row"`
	Colum        string    `gorm:"column:colum;type:varchar(20);not null;" json:"colum"`
	PassengerIds PassengerIds
}

func (p *PassengerIds) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &p)
}

func (p PassengerIds) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// marshal passenger ids in json
func (p PassengerIds) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(p))
}

// unmarshal passenger ids in json
func (p *PassengerIds) UnmarshalJSON(data []byte) error {
	var passengerIds []string
	if err := json.Unmarshal(data, &passengerIds); err != nil {
		return err
	}
	*p = PassengerIds(passengerIds)
	return nil
}
