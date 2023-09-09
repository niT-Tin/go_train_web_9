package initialize

import (
	"fmt"
	"gotrains/ticketorder_srvs/ticket_srv/global"
	"gotrains/ticketorder_srvs/ticket_srv/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() {
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	// dsn := "root:lzh@tcp(127.0.0.1:3306)/lzhshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.Config.MySqlInfo.User,
		global.Config.MySqlInfo.Pass,
		global.Config.MySqlInfo.Host,
		global.Config.MySqlInfo.Port,
		global.Config.MySqlInfo.Name,
	)
	// dsn := "root:lzh@tcp(127.0.0.1:3306)/lzhshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if nil != err {
		panic(err)
	}
	err = global.DB.AutoMigrate(&model.User{}, &model.Passenger{})
	if err != nil {
		panic(err)
	}
}
