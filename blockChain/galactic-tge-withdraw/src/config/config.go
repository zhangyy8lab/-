package config

import (
	"fmt"
	ini "gopkg.in/ini.v1"
	"os"
)

type Base struct {
	Address   string
	Precision int
	Price     float64
}

type EthereumStruct struct {
	Eth         Base
	Usdt        Base
	Usdc        Base
	Bnb         Base
	Trias       Base
	ApiKey      string
	BlockHeight int64
}

type ArbitrumStruct struct {
	Eth         Base
	Usdt        Base
	Usdc        Base
	ApiKey      string
	BlockHeight int64
}

type BNBStruct struct {
	Eth         Base
	Bnb         Base
	Usdt        Base
	Usdc        Base
	Trias       Base
	ApiKey      string
	BlockHeight int64
}

type PolygonStruct struct {
	Eth         Base
	Bnb         Base
	Usdt        Base
	Usdc        Base
	Trias       Base
	ApiKey      string
	BlockHeight int64
}

type OptimismStruct struct {
	Eth         Base
	Usdt        Base
	Usdc        Base
	ApiKey      string
	BlockHeight int64
}

var (
	PgDb     string
	PgDbUser string
	PgDbPass string
	PgHost   string
	PgPort   int

	Eth      EthereumStruct
	Arb      ArbitrumStruct
	Bnb      BNBStruct
	Polygon  PolygonStruct
	Optimism OptimismStruct
)

func Init() {
	dir, _ := os.Getwd()
	filePath := fmt.Sprintf("%v/src/config/config.ini", dir)
	file, err := ini.Load(filePath)
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadPgDb(file)
	ethLoad(file)
	arbLoad(file)
	bnbLoad(file)
	polygonLoad(file)
	optimismLoad(file)
}

func LoadPgDb(file *ini.File) {
	PgHost = file.Section("pgDb").Key("DbAddress").String()
	PgPort, _ = file.Section("pgDb").Key("DbPort").Int()
	PgDb = file.Section("pgDb").Key("DbName").String()
	PgDbUser = file.Section("pgDb").Key("DbUser").String()
	PgDbPass = file.Section("pgDb").Key("DbPassWord").String()
}

func ethLoad(file *ini.File) {
	Eth.ApiKey = file.Section("eth").Key("ApiKey").String()
	Eth.BlockHeight, _ = file.Section("eth").Key("BlockHeight").Int64()

	Eth.Eth.Address = file.Section("eth.eth").Key("Address").String()
	Eth.Eth.Precision, _ = file.Section("eth.eth").Key("Precision").Int()
	Eth.Eth.Price, _ = file.Section("eth.eth").Key("Price").Float64()

	Eth.Usdt.Address = file.Section("eth.usdt").Key("Address").String()
	Eth.Usdt.Precision, _ = file.Section("eth.usdt").Key("Precision").Int()
	Eth.Usdt.Price, _ = file.Section("eth.usdt").Key("Price").Float64()

	Eth.Usdc.Address = file.Section("eth.usdc").Key("Address").String()
	Eth.Usdc.Precision, _ = file.Section("eth.usdc").Key("Precision").Int()
	Eth.Usdc.Price, _ = file.Section("eth.usdc").Key("Price").Float64()

	Eth.Bnb.Address = file.Section("eth.bnb").Key("Address").String()
	Eth.Bnb.Precision, _ = file.Section("eth.bnb").Key("Precision").Int()
	Eth.Bnb.Price, _ = file.Section("eth.bnb").Key("Price").Float64()

	Eth.Trias.Address = file.Section("eth.trias").Key("Address").String()
	Eth.Trias.Precision, _ = file.Section("eth.trias").Key("Precision").Int()
	Eth.Trias.Price, _ = file.Section("eth.trias").Key("Price").Float64()

}

func arbLoad(file *ini.File) {
	Arb.ApiKey = file.Section("arb").Key("ApiKey").String()
	Arb.BlockHeight, _ = file.Section("arb").Key("BlockHeight").Int64()

	Arb.Eth.Address = file.Section("arb.eth").Key("Address").String()
	Arb.Eth.Precision, _ = file.Section("arb.eth").Key("Precision").Int()
	Arb.Eth.Price, _ = file.Section("arb.eth").Key("Price").Float64()

	Arb.Usdt.Address = file.Section("arb.usdt").Key("Address").String()
	Arb.Usdt.Precision, _ = file.Section("arb.usdt").Key("Precision").Int()
	Arb.Usdt.Price, _ = file.Section("arb.usdt").Key("Price").Float64()

	Arb.Usdc.Address = file.Section("arb.usdc").Key("Address").String()
	Arb.Usdc.Precision, _ = file.Section("arb.usdc").Key("Precision").Int()
	Arb.Usdc.Price, _ = file.Section("arb.usdc").Key("Price").Float64()
}

func bnbLoad(file *ini.File) {
	Bnb.ApiKey = file.Section("bnb").Key("ApiKey").String()
	Bnb.BlockHeight, _ = file.Section("bnb").Key("BlockHeight").Int64()

	Bnb.Eth.Address = file.Section("bnb.eth").Key("Address").String()
	Bnb.Eth.Precision, _ = file.Section("bnb.eth").Key("Precision").Int()
	Bnb.Eth.Price, _ = file.Section("bnb.eth").Key("Price").Float64()

	Bnb.Bnb.Address = file.Section("bnb.bnb").Key("Address").String()
	Bnb.Bnb.Precision, _ = file.Section("bnb.bnb").Key("Precision").Int()
	Bnb.Bnb.Price, _ = file.Section("bnb.bnb").Key("Price").Float64()

	Bnb.Usdt.Address = file.Section("bnb.usdt").Key("Address").String()
	Bnb.Usdt.Precision, _ = file.Section("bnb.usdt").Key("Precision").Int()
	Bnb.Usdt.Price, _ = file.Section("bnb.usdt").Key("Price").Float64()

	Bnb.Usdc.Address = file.Section("bnb.usdc").Key("Address").String()
	Bnb.Usdc.Precision, _ = file.Section("bnb.usdc").Key("Precision").Int()
	Bnb.Usdc.Price, _ = file.Section("bnb.usdc").Key("Price").Float64()

	Bnb.Trias.Address = file.Section("bnb.trias").Key("Address").String()
	Bnb.Trias.Precision, _ = file.Section("bnb.trias").Key("Precision").Int()
	Bnb.Trias.Price, _ = file.Section("bnb.trias").Key("Price").Float64()
}

func polygonLoad(file *ini.File) {
	Polygon.ApiKey = file.Section("polygon").Key("ApiKey").String()
	Polygon.BlockHeight, _ = file.Section("polygon").Key("BlockHeight").Int64()

	Polygon.Eth.Address = file.Section("polygon.wEth").Key("Address").String()
	Polygon.Eth.Precision, _ = file.Section("polygon.wEth").Key("Precision").Int()
	Polygon.Eth.Price, _ = file.Section("polygon.wEth").Key("Price").Float64()

	Polygon.Bnb.Address = file.Section("polygon.bnb").Key("Address").String()
	Polygon.Bnb.Precision, _ = file.Section("polygon.bnb").Key("Precision").Int()
	Polygon.Bnb.Price, _ = file.Section("polygon.bnb").Key("Price").Float64()

	Polygon.Usdt.Address = file.Section("polygon.usdt").Key("Address").String()
	Polygon.Usdt.Precision, _ = file.Section("polygon.usdt").Key("Precision").Int()
	Polygon.Usdt.Price, _ = file.Section("polygon.usdt").Key("Price").Float64()

	Polygon.Usdc.Address = file.Section("polygon.usdc").Key("Address").String()
	Polygon.Usdc.Precision, _ = file.Section("polygon.usdc").Key("Precision").Int()
	Polygon.Usdc.Price, _ = file.Section("polygon.usdc").Key("Price").Float64()

	Polygon.Trias.Address = file.Section("polygon.trias").Key("Address").String()
	Polygon.Trias.Precision, _ = file.Section("polygon.trias").Key("Precision").Int()
	Polygon.Trias.Price, _ = file.Section("polygon.trias").Key("Price").Float64()

}

func optimismLoad(file *ini.File) {
	Optimism.ApiKey = file.Section("optimism").Key("ApiKey").String()
	Optimism.BlockHeight, _ = file.Section("optimism").Key("BlockHeight").Int64()

	Optimism.Eth.Address = file.Section("optimism.eth").Key("Address").String()
	Optimism.Eth.Precision, _ = file.Section("optimism.eth").Key("Precision").Int()
	Optimism.Eth.Price, _ = file.Section("optimism.eth").Key("Price").Float64()

	Optimism.Usdt.Address = file.Section("optimism.usdt").Key("Address").String()
	Optimism.Usdt.Precision, _ = file.Section("optimism.usdt").Key("Precision").Int()
	Optimism.Usdt.Price, _ = file.Section("optimism.usdt").Key("Price").Float64()

	Optimism.Usdc.Address = file.Section("optimism.usdc").Key("Address").String()
	Optimism.Usdc.Precision, _ = file.Section("optimism.usdc").Key("Precision").Int()
	Optimism.Usdc.Price, _ = file.Section("optimism.usdc").Key("Price").Float64()

}
