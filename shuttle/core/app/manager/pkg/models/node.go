package models

import (
	"google.dev/google/shuttle/proto/manager"
	"gorm.io/gorm"
)

// Node 节点
type Node struct {
	gorm.Model
	NodeID          string `gorm:"type:varchar(300);uniqueIndex" json:"node_id"` // 节点id
	NodeName        string `gorm:"type:varchar(300)" json:"node_name"`           // 节点名称
	Country         string `gorm:"type:varchar(300)" json:"country"`             // 国家
	Describe        string `gorm:"type:varchar(300)" json:"describe"`            // 描述
	IP              string `gorm:"type:varchar(300)" json:"ip"`                  // 注册ip
	InternetAddress string `gorm:"type:varchar(300)" json:"internet_address"`    // 节点链接地址

	// protocol
	Protocol manager.Protocol `json:"protocol"`                          // 协议
	WssPath  string           `gorm:"type:varchar(300)" json:"wss_path"` // wss path

	MountSupport bool // 挂载支持
}

type AppNodeMapping struct {
	gorm.Model
	AppID  string  `gorm:"type:varchar(300)" json:"app_id"`  // 应用id
	NodeID string  `gorm:"type:varchar(300)" json:"node_id"` // 节点id
	Sort   float32 // 排序
	Free   bool    // 免费节点

	Node *Node `json:"-" gorm:"-"`
}
