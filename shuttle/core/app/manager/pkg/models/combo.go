package models

import "gorm.io/gorm"

// Combo 套餐
type Combo struct {
	gorm.Model
	ComboID  string  `gorm:"type:varchar(300);uniqueIndex" json:"combo_id"` // 套餐id
	AppID    string  `gorm:"type:varchar(300)" json:"app_id"`               // 应用id
	Describe string  `gorm:"type:varchar(300)" json:"describe"`             // 描述
	Traffic  int64   // 套餐流量 G
	Sort     float32 // 排序
	Day      int64   // 套餐天数
	Amount   float32 // 套餐金额
}
