package model

import (
	"time"

	"gorm.io/gorm"
)

type PassengerType bool

const (
	// 乘客类型
	PassengerTypeAdult PassengerType = true  // 成人
	PassengerTypeChild               = false // 儿童
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
	Mobile string `gorm:"column:mobile;type:varchar(11);not null;unique_index;"`
	Passwd string `gorm:"column:passwd;type:varchar(32);not null;"`
	Role   int32  `gorm:"type:int comment '1 表示普通用户 2 表示管理员';column:role;default:1;"`
}

type Passenger struct {
	BaseModel
	Name   string        `gorm:"column:name;type:varchar(20);not null;"`
	IdCard string        `gorm:"column:id_card;type:varchar(18);not null;unique_index;"`
	Type   PassengerType `gorm:"column:type;type:tinyint(1);not null;"`
}
