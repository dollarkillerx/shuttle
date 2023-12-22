package simple

import (
	"fmt"
	"sync"

	"google.dev/google/common/pkg/client"
	"google.dev/google/common/pkg/conf"
	"google.dev/google/shuttle/core/app/manager/generated"
	"google.dev/google/shuttle/core/app/manager/pkg/models"
	"gorm.io/gorm"
)

type Simple struct {
	db *gorm.DB

	inventoryMu sync.Mutex
}

func NewSimple(config conf.PostgresConfiguration) (*Simple, error) {
	postgresClient, err := client.PostgresClient(config, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("AutoMigrate ....")
	postgresClient.AutoMigrate(
	//&models.User{},
	//
	//&models.App{},
	//&models.Combo{},
	//&models.Server{},
	//&models.Node{},
	//&models.AppNodeMapping{},
	)

	s := &Simple{
		db: postgresClient,
	}

	//s.initBaseInfo()

	fmt.Println("AutoMigrate end ....")

	return s, nil
}

func (s *Simple) DB() *gorm.DB {
	return s.db
}

func (s *Simple) initBaseInfo() {
	s.initApp()   // 应用
	s.initCombo() // 套餐
	//s.initNodes() // 初始化节点
}

func (s *Simple) initApp() {
	apps := []models.App{
		models.App{
			AppID:                    "80f8342f0487",
			AppName:                  "Free Komi VPN",
			Describe:                 "Komi VPN 免费版本",
			AppVersion:               0.01,
			MinimumVersion:           0.01,
			State:                    generated.AppStateEnable,
			NoAuthenticationRequired: true,
			ErrorNotification:        "",
			NormalNotification:       "",
		},
	}

	for i := range apps {
		item := apps[i]
		err := s.db.Model(&models.App{}).
			Where(&models.App{AppID: item.AppID}).
			//Attrs(&upData). // 存在就更新attrs 字段中内容
			FirstOrCreate(&item).Error
		if err != nil {
			panic(err)
		}
	}
}

func (s *Simple) initCombo() {
	freeCombo := []models.Combo{
		models.Combo{
			ComboID:  "80f8342f04871",
			AppID:    "80f8342f0487",
			Describe: "Free",
			Traffic:  100, // 100G
			Sort:     0,
			Day:      1000,
			Amount:   -1,
		},
		models.Combo{
			ComboID:  "80f8342f04872",
			AppID:    "80f8342f0487",
			Describe: "1 WEEK",
			Traffic:  100, // 100G
			Sort:     1,
			Day:      7,
			Amount:   0.5,
		},
		models.Combo{
			ComboID:  "80f8342f04873",
			AppID:    "80f8342f0487",
			Describe: "1 MOUNTH",
			Traffic:  100, // 100G
			Sort:     2,
			Day:      30,
			Amount:   0.99,
		},
		models.Combo{
			ComboID:  "80f8342f04874",
			AppID:    "80f8342f0487",
			Describe: "1 YEAR",
			Traffic:  1000, // 100G
			Sort:     3,
			Day:      360,
			Amount:   9.99,
		},
	}

	for i := range freeCombo {
		item := freeCombo[i]
		err := s.db.Model(&models.Combo{}).
			Where(&models.Combo{
				ComboID: item.ComboID,
			}).FirstOrCreate(&item).Error
		if err != nil {
			panic(err)
		}
	}

}

func (s *Simple) initNodes() {
	freeNodes := []models.Node{
		models.Node{
			NodeID:   "80f8342f04872",
			NodeName: "United States Free",
			Country:  "us",
			Describe: "United States Free",
		},
		models.Node{
			NodeID:   "80f8342f04873",
			NodeName: "United States VIP",
			Country:  "us",
			Describe: "United States VIP",
		},
		models.Node{
			NodeID:   "80f8342f04874",
			NodeName: "Japan VIP",
			Country:  "jp",
			Describe: "Japan VIP",
		},
		models.Node{
			NodeID:   "80f8342f04875",
			NodeName: "Russia VIP",
			Country:  "ru",
			Describe: "Russia VIP",
		},
	}

	for i := range freeNodes {
		item := freeNodes[i]
		err := s.db.Model(&models.Node{}).
			Where(&models.Node{
				NodeID: item.NodeID,
			}).FirstOrCreate(&item).Error
		if err != nil {
			panic(err)
		}
	}
}
