package client

import (
	"fmt"

	"google.dev/google/common/pkg/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresClient(conf conf.PostgresConfiguration, gormConfig *gorm.Config) (*gorm.DB, error) {
	if conf.TimeZone == "" {
		conf.TimeZone = "Asia/Shanghai"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=%s", conf.Host, conf.User, conf.Password, conf.DBName, conf.Port, conf.TimeZone)
	if !conf.SSLMode {
		dsn += " sslmode=disable"
	}

	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}

	return gorm.Open(postgres.Open(dsn), gormConfig)
}
