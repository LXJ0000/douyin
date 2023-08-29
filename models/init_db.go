package models

import (
	"douyin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	DB, err = gorm.Open(mysql.Open(config.DBConnectString()), &gorm.Config{
		PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction: true, //禁用默认事务操作
	})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&UserInfo{}, &UserLogin{}, &Video{})
	if err != nil {
		panic(err)
	}
}
