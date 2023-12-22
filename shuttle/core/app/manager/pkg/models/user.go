package models

import "gorm.io/gorm"

type User struct {
	gorm.Model         // 继承 BasicModel 结构体
	OS          string `gorm:"type:varchar(300)" json:"os"`                // android, ios, mac, windows
	Token       string `gorm:"type:varchar(300);uniqueIndex" json:"token"` // 用户唯一标识 免密用户
	DeviceID    string `gorm:"type:varchar(300)" json:"device_id"`         // 设备id
	DeviceName  string `gorm:"type:varchar(300)" json:"device_name"`       // 设备名称
	AppID       string `gorm:"type:varchar(300)" json:"app_id"`            // app项目的id ， 免费卖广告， 收费卖套餐
	LastLoginIP string `gorm:"type:varchar(300)" json:"last_login_ip"`     // 上次登陆ip
	LoginIP     string `gorm:"type:varchar(300)" json:"login_ip"`          // 本次登陆ip
	RegIP       string `gorm:"type:varchar(300)" json:"reg_ip"`            // 注册ip
}
