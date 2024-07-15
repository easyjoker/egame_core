package egame_core

import (
	"fmt"
	"time"

	"egame_core/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize() {
	config := config.GetConfig()

	// 連線資料庫
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/egame?charset=utf8mb4&parseTime=True&loc=Local", config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	sqlDb, err := DB.DB()
	if err != nil {
		fmt.Println(err)
	}

	sqlDb.SetMaxIdleConns(config.Database.MaxIdle)                               // 10個連線
	sqlDb.SetMaxOpenConns(config.Database.MaxLife)                               // 100個連線
	sqlDb.SetConnMaxLifetime(time.Hour * time.Duration(config.Database.MaxLife)) // 連線最大存活時間
}

func Close() {
	sqlDb, err := DB.DB()
	if err != nil {
		fmt.Println(err)
	}

	sqlDb.Close()
}

func GetDB() *gorm.DB {
	if DB == nil {
		Initialize()
	}
	return DB
}
