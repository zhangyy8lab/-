package utils

import (
	"fmt"
	"galactic-tge-withdraw/src/config"
	"galactic-tge-withdraw/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbInit() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d, sslmode=disable",
		config.PgHost, config.PgDbUser, config.PgDbPass, config.PgDb, config.PgPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	fmt.Println("Successfully connected to database")
}

func GetAll() []models.Account {
	// 查询数据
	var accounts []models.Account
	//result := DB.Debug().Where("star >= 5000").Find(&accounts)
	result := DB.Debug().Where("star >= 5000").Find(&accounts)
	if result.Error != nil {
		panic(result.Error)
	}
	return accounts
}

// CheckExistWithDraw 检查是否已经提现
func CheckExistWithDraw(address, chain, coinKind string) bool {
	var withDraw models.WithDraw
	result := DB.Where(
		"account = ? and chain = ? and coin_kind = ?", address, chain, coinKind).First(&withDraw)
	if result.Error == gorm.ErrRecordNotFound {
		fmt.Printf("address not exist: %s chain: %s coinKind: %s\n", address, chain, coinKind)
		return false
	} else if result.Error != nil {
		panic(result.Error)
	}
	fmt.Printf("address exist: %s chain: %s coinKind: %s\n", address, chain, coinKind)
	return true
}

func WithDrawCount(accountObj models.WithDraw, valuesCount float64) {
	var withCount models.WithCount
	withCount.CountValue = valuesCount
	withCount.Account = accountObj.Account
	withCount.Stars = accountObj.Stars
	DB.Debug().Create(&withCount)

}
