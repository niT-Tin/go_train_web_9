package models

import (
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

type SeatInfo struct {
	Id            int64  `json:"id"`
	TrainCode     string `json:"train_code"`
	CarriageIndex int32  `json:"carriage_index"`
	SeatType      string `json:"seat_type"`
	SeatIndex     int32  `json:"seat_index"`
	Row           string `json:"row"`
	Column        string `json:"column"`
}

type TrainInfo struct {
	Id           int64  `json:"id"`
	Code         string `json:"code"`
	Type         string `json:"type"`
	StartStation string `json:"start_station"`
	EndStation   string `json:"end_station"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	StartPinyin  string `json:"start_pinyin"`
	EndPinyin    string `json:"end_pinyin"`
}

type TicketInfo struct {
	ID               int64     `gorm:"column:id;primaryKey;comment:id" json:"id"`                                  // id
	Date             time.Time `gorm:"column:date;not null;comment:日期" json:"date"`                                // 日期
	TrainCode        string    `gorm:"column:train_code;not null;comment:车次编号" json:"train_code"`                  // 车次编号
	Start            string    `gorm:"column:start;not null;comment:出发站" json:"start"`                             // 出发站
	StartPinyin      string    `gorm:"column:start_pinyin;not null;comment:出发站拼音" json:"start_pinyin"`             // 出发站拼音
	StartTime        time.Time `gorm:"column:start_time;not null;comment:出发时间" json:"start_time"`                  // 出发时间
	StartIndex       int32     `gorm:"column:start_index;not null;comment:出发站序|本站是整个车次的第几站" json:"start_index"`    // 出发站序|本站是整个车次的第几站
	End              string    `gorm:"column:end;not null;comment:到达站" json:"end"`                                 // 到达站
	EndPinyin        string    `gorm:"column:end_pinyin;not null;comment:到达站拼音" json:"end_pinyin"`                 // 到达站拼音
	EndTime          time.Time `gorm:"column:end_time;not null;comment:到站时间" json:"end_time"`                      // 到站时间
	EndIndex         int32     `gorm:"column:end_index;not null;comment:到站站序|本站是整个车次的第几站" json:"end_index"`        // 到站站序|本站是整个车次的第几站
	FirstClassLast   int32     `gorm:"column:first_class_last;not null;comment:一等座余票" json:"first_class_last"`     // 一等座余票
	FirstClassPrice  float64   `gorm:"column:first_class_price;not null;comment:一等座票价" json:"first_class_price"`   // 一等座票价
	SecondClassLast  int32     `gorm:"column:second_class_last;not null;comment:二等座余票" json:"second_class_last"`   // 二等座余票
	SecondClassPrice float64   `gorm:"column:second_class_price;not null;comment:二等座票价" json:"second_class_price"` // 二等座票价
	SoftBerthLast    int32     `gorm:"column:soft_berth_last;not null;comment:软卧余票" json:"soft_berth_last"`        // 软卧余票
	SoftBerthPrice   float64   `gorm:"column:soft_berth_price;not null;comment:软卧票价" json:"soft_berth_price"`      // 软卧票价
	HardBerthLast    int32     `gorm:"column:hard_berth_last;not null;comment:硬卧余票" json:"hard_berth_last"`        // 硬卧余票
	HardBerthPrice   float64   `gorm:"column:hard_berth_price;not null;comment:硬卧票价" json:"hard_berth_price"`      // 硬卧票价
	CreateTime       time.Time `gorm:"column:create_time;comment:新增时间" json:"create_time"`                         // 新增时间
	UpdateTime       time.Time `gorm:"column:update_time;comment:修改时间" json:"update_time"`                         // 修改时间
}
