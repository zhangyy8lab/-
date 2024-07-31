package main

import (
	"galactic-tge-withdraw/src/common"
	"galactic-tge-withdraw/src/config"
	"galactic-tge-withdraw/src/utils"
)

func main() {
	config.Init()
	utils.DbInit()

	//common.WithDrawCount()
	common.Census()
}
