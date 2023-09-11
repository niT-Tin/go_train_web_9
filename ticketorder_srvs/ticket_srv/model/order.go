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
