package client

import (
	"fmt"

	"google.dev/google/common/pkg/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MySQLClient(conf conf.MySQLConfiguration, gormConfig *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}

	return gorm.Open(mysql.Open(dsn), gormConfig)
}
