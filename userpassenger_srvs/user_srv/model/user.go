package model

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

type PassengerType int64

const (
	// 乘客类型
	PassengerTypeAdult PassengerType = 0 // 成人
	PassengerTypeChild               = 1 // 儿童
)

type BaseModel struct {
	ID        int32          `gorm:"primary_key;"`
	CreatedAt time.Time      `gorm:"column:add_time;"`
	UpdatedAt time.Time      `gorm:"column:update_time;"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_time;"`
	IsDelete  bool           `gorm:"column:is_delete;"`
}

type User struct {
	BaseModel
	Mobile   string    `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string    `gorm:"type:varchar(100);not null"`
	Nickname string    `gorm:"type:varchar(20);"`
	Birthday time.Time `gorm:"type:datetime;"`
	Gender   string    `gorm:"type:varchar(6) comment 'female表示女性，male表示男性';column:gender;default:male;"`
	Role     int32     `gorm:"type:int comment '1 表示普通用户 2 表示管理员';column:role;default:1;"`
}

type Passenger struct {
	BaseModel
	Name     string        `gorm:"column:name;type:varchar(20);not null;"`
	UserID   int32         `gorm:"column:user_id;type:int;not null;"`
	IdCard   string        `gorm:"column:id_card;type:varchar(18);not null;unique_index;"`
	Type     PassengerType `gorm:"column:type;type:int;not null;"`
	SeatType string        `gorm:"column:seat_type;type:varchar(10);not null;"`
	Seat     string        `gorm:"column:seat;type:varchar(10);not null;"`
}

func (p *PassengerType) Scan(value interface{}) error {
	*p = PassengerType(value.(int64))
	return nil
}

func (p *PassengerType) Value() (driver.Value, error) {
	return int64(*p), nil
}
