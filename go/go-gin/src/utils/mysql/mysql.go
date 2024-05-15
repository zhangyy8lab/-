package mysql

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type Client struct {
	DbName     string `json:"db_name"`
	DbAddress  string `json:"db_address"`
	DbPort     int    `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbPassWord string `json:"db_pass_word"`
	Mode       string `json:"mode"`
}

func (c *Client) Instance() {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True",
		c.DbUser, c.DbPassWord, c.DbAddress, c.DbPort, c.DbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("content mysql DB failed")
	}

	if c.Mode != "release" {
		DB = DB.Debug()
	}
}

type Db struct {
	Client Client
	DB     *gorm.DB
	GCtx   *gin.Context
	Error  error
}
