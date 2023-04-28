package models

import (
	"fmt"

	"github.com/Web-Engineering-XDU/Project-Backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(config config.MySQL) error {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/huggo?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	return err
}

func DB() *gorm.DB{
	return db
}