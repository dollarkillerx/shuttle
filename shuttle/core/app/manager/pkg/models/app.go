package models

import (
	"google.dev/google/shuttle/core/app/manager/generated"
	"gorm.io/gorm"
)

// App App
type App struct {
	gorm.Model
	AppID                    string             `gorm:"type:varchar(300);uniqueIndex" json:"app_id"`  // app id
	AppName                  string             `gorm:"type:varchar(300)" json:"app_name"`            // app name
	Describe                 string             `gorm:"type:varchar(300)" json:"describe"`            // 描述
	AppVersion               float32            `json:"app_version"`                                  // app 版本
	MinimumVersion           float32            `json:"minimum_version"`                              // 最低运行版本
	NoAuthenticationRequired bool               `json:"no_authentication_required"`                   // 是否需要注册验证
	State                    generated.AppState `gorm:"type:varchar(300)" json:"state"`               // app 状态  （enable，disabled）
	ErrorNotification        string             `gorm:"type:varchar(500)" json:"error_notification"`  // 错误通知
	NormalNotification       string             `gorm:"type:varchar(500)" json:"normal_notification"` // 正常通知
}
