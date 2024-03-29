// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameTrainStation = "train_station"

// TrainStation mapped from table <train_station>
type TrainStation struct {
	ID         int64     `gorm:"column:id;primaryKey;comment:id" json:"id"`                   // id
	TrainCode  string    `gorm:"column:train_code;not null;comment:车次编号" json:"train_code"`   // 车次编号
	Index      int32     `gorm:"column:index;not null;comment:站序" json:"index"`               // 站序
	Name       string    `gorm:"column:name;not null;comment:站名" json:"name"`                 // 站名
	NamePinyin string    `gorm:"column:name_pinyin;not null;comment:站名拼音" json:"name_pinyin"` // 站名拼音
	InTime     time.Time `gorm:"column:in_time;comment:进站时间" json:"in_time"`                  // 进站时间
	OutTime    time.Time `gorm:"column:out_time;comment:出站时间" json:"out_time"`                // 出站时间
	StopTime   time.Time `gorm:"column:stop_time;comment:停站时长" json:"stop_time"`              // 停站时长
	Km         float64   `gorm:"column:km;not null;comment:里程（公里）|从上一站到本站的距离" json:"km"`      // 里程（公里）|从上一站到本站的距离
	CreateTime time.Time `gorm:"column:create_time;comment:新增时间" json:"create_time"`          // 新增时间
	UpdateTime time.Time `gorm:"column:update_time;comment:修改时间" json:"update_time"`          // 修改时间
}

// TableName TrainStation's table name
func (*TrainStation) TableName() string {
	return TableNameTrainStation
}
