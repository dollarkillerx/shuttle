package models

import "gorm.io/gorm"

// Server 服务套餐
type Server struct {
	gorm.Model
	ComboId  string `gorm:"type:varchar(300)" json:"combo_id"` // 套餐id
	Describe string `gorm:"type:varchar(300)" json:"describe"` // 套餐描述
	Traffic  int64  `gorm:"type:varchar(300)" json:"traffic"`  // 套餐流量
}
