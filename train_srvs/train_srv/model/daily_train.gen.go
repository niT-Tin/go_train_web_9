// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameDailyTrain = "daily_train"

// DailyTrain mapped from table <daily_train>
type DailyTrain struct {
	ID          int64     `gorm:"column:id;primaryKey;comment:id" json:"id"`                       // id
	Date        time.Time `gorm:"column:date;not null;comment:日期" json:"date"`                     // 日期
	Code        string    `gorm:"column:code;not null;comment:车次编号" json:"code"`                   // 车次编号
	Type        string    `gorm:"column:type;not null;comment:车次类型|枚举[TrainTypeEnum]" json:"type"` // 车次类型|枚举[TrainTypeEnum]
	Start       string    `gorm:"column:start;not null;comment:始发站" json:"start"`                  // 始发站
	StartPinyin string    `gorm:"column:start_pinyin;not null;comment:始发站拼音" json:"start_pinyin"`  // 始发站拼音
	StartTime   time.Time `gorm:"column:start_time;not null;comment:出发时间" json:"start_time"`       // 出发时间
	End         string    `gorm:"column:end;not null;comment:终点站" json:"end"`                      // 终点站
	EndPinyin   string    `gorm:"column:end_pinyin;not null;comment:终点站拼音" json:"end_pinyin"`      // 终点站拼音
	EndTime     time.Time `gorm:"column:end_time;not null;comment:到站时间" json:"end_time"`           // 到站时间
	CreateTime  time.Time `gorm:"column:create_time;comment:新增时间" json:"create_time"`              // 新增时间
	UpdateTime  time.Time `gorm:"column:update_time;comment:修改时间" json:"update_time"`              // 修改时间
}

// TableName DailyTrain's table name
func (*DailyTrain) TableName() string {
	return TableNameDailyTrain
}
