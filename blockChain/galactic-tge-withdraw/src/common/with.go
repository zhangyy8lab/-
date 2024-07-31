package common

import (
	"context"
	"encoding/json"
	"fmt"
	"galactic-tge-withdraw/src/config"
	"galactic-tge-withdraw/src/models"
	"galactic-tge-withdraw/src/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"strings"
)

// 合约ABI，确保替换为你实际的合约ABI
const tokenABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}]`

// const BscApiKey = "UIYZSHG5G8FSTXB7FFED9YFB6H3DBU71AN"
const BscApiKey = "DRPTYBTH8A3YPPRVTV1DK3GY3F228KZE83"

type Coin struct {
	ApiToken     string
	UserAddress  string
	TokenAddress string
	Precision    int
	Price        float64
	BlockHeight  int64
	ChainName    string
	CoinKind     string
	CoinCount    float64
	MoneyValue   float64
	Account      models.Account
}

func (c *Coin) getClient() *ethclient.Client {
	client, err := ethclient.Dial(c.ApiToken)
	if err != nil {
		log.Fatalf("Failed to connect to the Arbitrum network: %v", err)
	}
	return client
}

func (c *Coin) TokenBalance() error {
	client := c.getClient()

	// 合约地址和用户地址
	userAddress := common.HexToAddress(c.UserAddress)
	tokenAddress := common.HexToAddress(c.TokenAddress)

	// 查询指定块高
	blockNumber := big.NewInt(c.BlockHeight)

	// 合约ABI解析
	parsedABI, err := abi.JSON(strings.NewReader(tokenABI))
	if err != nil {
		log.Fatalf("Failed to parse token ABI: %v", err)
	}

	data, err := parsedABI.Pack("balanceOf", userAddress)
	if err != nil {
		return err
	}

	callMsg := ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), callMsg, blockNumber)
	if err != nil {
		return err
	}

	var balance = new(big.Int)
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		fmt.Println("parsedABI.UnpackIntoInterface err:", err.Error())
	}

	// 计算用户代币的数量
	value := new(big.Float).Mul(new(big.Float).SetInt(balance), big.NewFloat(1))
	value.Quo(value, new(big.Float).SetInt(big.NewInt(int64(math.Pow10(c.Precision)))))
	c.CoinCount, _ = value.Float64()

	if c.Price == 1 {
		c.MoneyValue, _ = value.Float64()
		return nil
	} else {
		// 创建一个 big.Float 类型用于存储结果
		moneyValue := new(big.Float)
		// 进行乘法运算
		moneyValue.Mul(big.NewFloat(c.CoinCount), big.NewFloat(c.Price))
		c.MoneyValue, _ = moneyValue.Float64()
		return nil
	}
}

func (c *Coin) Balance() error {

	client := c.getClient()
	// 合约地址和用户地址
	userAddress := common.HexToAddress(c.UserAddress)

	// 查询指定块高
	blockNumber := big.NewInt(c.BlockHeight)

	// 获取用户在指定区块高度的以太币余额
	balance, err := client.BalanceAt(context.Background(), userAddress, blockNumber)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	// 计算用户代币的数量
	value := new(big.Float).Mul(new(big.Float).SetInt(balance), big.NewFloat(1))
	value.Quo(value, new(big.Float).SetInt(big.NewInt(int64(math.Pow10(c.Precision)))))
	c.CoinCount, _ = value.Float64()

	if c.Price == 1 {
		c.MoneyValue, _ = value.Float64()
		return nil
	} else {
		// 创建一个 big.Float 类型用于存储结果
		moneyValue := new(big.Float)
		// 进行乘法运算
		moneyValue.Mul(big.NewFloat(c.CoinCount), big.NewFloat(c.Price))

		c.MoneyValue, _ = moneyValue.Float64()
		return nil
	}
}

func (c *Coin) BNBBalance() error {
	baseURL := "https://api.bscscan.com/api"
	endpoint := fmt.Sprintf("%s?module=account&action=balance&address=%s&apikey=%s",
		baseURL, c.UserAddress, BscApiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalf("Failed to get request: %v", err)
	}

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Status != "1" {
		return err
	}

	balance := new(big.Int)
	balance.SetString(result.Result, 10)
	value := new(big.Float).Mul(new(big.Float).SetInt(balance), big.NewFloat(1))
	value.Quo(value, new(big.Float).SetInt(big.NewInt(int64(math.Pow10(18)))))

	c.CoinCount, _ = value.Float64()
	if c.Price == 1 {
		c.MoneyValue, _ = value.Float64()
		return nil
	} else {
		// 创建一个 big.Float 类型用于存储结果
		moneyValue := new(big.Float)
		// 进行乘法运算
		moneyValue.Mul(big.NewFloat(c.CoinCount), big.NewFloat(c.Price))

		c.MoneyValue, _ = moneyValue.Float64()
		return nil
	}
}

func (c *Coin) BnbTokenBalance() error {
	baseURL := "https://api.bscscan.com/api"
	endpoint := fmt.Sprintf("%s?module=account&action=tokenbalance&contractaddress=%s&address=%s&apikey=%s",
		baseURL, c.TokenAddress, c.UserAddress, BscApiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Status != "1" {
		return err
	}

	// 计算用户代币的数量
	balance := new(big.Int)
	balance.SetString(result.Result, 10)
	value := new(big.Float).Mul(new(big.Float).SetInt(balance), big.NewFloat(1))
	value.Quo(value, new(big.Float).SetInt(big.NewInt(int64(math.Pow10(c.Precision)))))

	c.CoinCount, _ = value.Float64()
	if c.Price == 1 {
		c.MoneyValue, _ = value.Float64()
		return nil
	} else {
		// 创建一个 big.Float 类型用于存储结果
		moneyValue := new(big.Float)
		// 进行乘法运算
		moneyValue.Mul(big.NewFloat(c.CoinCount), big.NewFloat(c.Price))

		c.MoneyValue, _ = moneyValue.Float64()
		return nil
	}
}

func (c *Coin) BnbUsdtBalance() error {
	rpcURL := "https://bsc-dataseed.binance.org/"
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// 用户地址
	userAddress := common.HexToAddress(c.UserAddress)
	//blockNumber := big.NewInt(40338106) // 指定的区块高度

	// USDT合约ABI
	usdtABI := `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`
	parsedABI, err := abi.JSON(strings.NewReader(usdtABI))
	if err != nil {
		log.Fatalf("Failed to parse USDT ABI: %v", err)
	}

	data, err := parsedABI.Pack("balanceOf", userAddress)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	fmt.Println("tokenAddress:", c.TokenAddress)
	// 构造调用消息
	contractAddress := common.HexToAddress(c.TokenAddress)
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	// 调用合约
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	var balance = new(big.Int)
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		log.Fatalf("Failed to unpack result: %v", err)
	}

	// 计算用户代币的数量
	//balance := new(big.Int)
	balance.SetString(balance.String(), 10)
	value := new(big.Float).Mul(new(big.Float).SetInt(balance), big.NewFloat(1))
	value.Quo(value, new(big.Float).SetInt(big.NewInt(int64(math.Pow10(c.Precision)))))

	// 创建一个 big.Float 类型用于存储结果
	moneyValue := new(big.Float)
	// 进行乘法运算
	moneyValue.Mul(big.NewFloat(c.CoinCount), big.NewFloat(c.Price))

	c.MoneyValue, _ = moneyValue.Float64()
	return nil

	//fmt.Printf("USDT Balance at block %d: %s\n", "111111", balance.String())
}

func (c *Coin) InsertData() {
	var withdraw models.WithDraw
	withdraw.Stars = c.Account.Star
	withdraw.Height = c.BlockHeight
	withdraw.CoinKind = c.CoinKind
	withdraw.Account = c.UserAddress
	withdraw.Value = c.MoneyValue
	withdraw.Chain = c.ChainName
	withdraw.CoinCount = c.CoinCount
	withdraw.CoinPrice = c.Price

	//// 检查是否已经存在
	//if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
	//	return
	//}

	if err := utils.DB.Debug().Create(&withdraw).Error; err != nil {
		log.Printf("insert data error: %v", err.Error())
		os.Exit(1)
	}
}

func (c *Coin) InsertWithCount() {
	var count models.WithCount
	count.Account = c.Account.Account
	count.Stars = c.Account.Star
	count.CountValue = c.MoneyValue
	if err := utils.DB.Debug().Create(&count).Error; err != nil {
		log.Printf("insert data error: %v", err.Error())
		os.Exit(1)
	}
}

func (c *Coin) CheckWithCountExist() bool {
	var count models.WithCount
	err := utils.DB.Debug().Where("account = ?", c.Account.Account).First(&count).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

// ====================== ETH ======================

func (c *Coin) EthMain() {
	c.ChainName = "Eth"
	c.CoinKind = "mainCoin"
	c.ApiToken = config.Eth.ApiKey
	c.UserAddress = c.Account.Account
	c.TokenAddress = config.Eth.Eth.Address
	c.Precision = config.Eth.Eth.Precision
	c.Price = config.Eth.Eth.Price
	c.BlockHeight = config.Eth.BlockHeight
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.Balance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) EthUsdt() {
	c.CoinKind = "USDT"
	c.TokenAddress = config.Eth.Usdt.Address
	c.Precision = config.Eth.Usdt.Precision
	c.Price = config.Eth.Usdt.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) EthUsdc() {
	c.CoinKind = "USDC"
	c.TokenAddress = config.Eth.Usdc.Address
	c.Precision = config.Eth.Usdc.Precision
	c.Price = config.Eth.Usdc.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) EthBnb() {
	c.CoinKind = "BNB"
	c.TokenAddress = config.Eth.Bnb.Address
	c.Precision = config.Eth.Bnb.Precision
	c.Price = config.Eth.Bnb.Price
	c.CoinCount = 0
	c.MoneyValue = 0
	//c.UserAddress = "0xDD98b5cD53144373AABf9e230C3Bee12dF58C438"
	c.BlockHeight = 20283292

	fmt.Println("------------------------------------")

	fmt.Println("c.Useraddress", c.UserAddress)
	fmt.Println("c.TokenAddress:", c.TokenAddress)
	fmt.Println("------------------------------------")

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	fmt.Println("c.CoinCount:", c.CoinCount)
	fmt.Println("c.Value:", c.MoneyValue)
	//os.Exit(0)
	c.InsertData()
}

func (c *Coin) EthTrias() {
	c.CoinKind = "Trias"
	c.TokenAddress = config.Eth.Trias.Address
	c.Precision = config.Eth.Trias.Precision
	c.Price = config.Eth.Trias.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

// ====================== ARB ======================

func (c *Coin) ArbMain() {
	c.ChainName = "Arbitrum"
	c.CoinKind = "mainCoin"
	c.ApiToken = config.Arb.ApiKey
	c.UserAddress = c.Account.Account
	c.TokenAddress = config.Arb.Eth.Address
	c.Precision = config.Arb.Eth.Precision
	c.Price = config.Arb.Eth.Price
	c.BlockHeight = config.Arb.BlockHeight
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.Balance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) ArbUsdt() {
	c.CoinKind = "USDT"
	c.TokenAddress = config.Arb.Usdt.Address
	c.Precision = config.Arb.Usdt.Precision
	c.Price = config.Arb.Usdt.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) ArbUsdc() {
	c.CoinKind = "USDC"
	c.TokenAddress = config.Arb.Usdc.Address
	c.Precision = config.Arb.Usdc.Precision
	c.Price = config.Arb.Usdc.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

// ====================== BNB ======================

func (c *Coin) BnbMain() {
	c.ChainName = "BNB"
	c.CoinKind = "mainCoin"
	c.ApiToken = config.Bnb.ApiKey
	//c.UserAddress = "0x9B50892B8CbEb07eC9812621Fc088E91f4f84D3f"
	c.UserAddress = c.Account.Account
	c.TokenAddress = config.Bnb.Bnb.Address
	c.Precision = config.Bnb.Bnb.Precision
	c.Price = config.Bnb.Bnb.Price
	c.BlockHeight = config.Bnb.BlockHeight
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.BNBBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	//fmt.Println("c.CoinCount:", c.CoinCount)
	c.InsertData()
}

func (c *Coin) BnbEth() {
	c.CoinKind = "ETH"
	c.TokenAddress = config.Bnb.Eth.Address
	c.Precision = config.Bnb.Eth.Precision
	c.Price = config.Bnb.Eth.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.BnbTokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}
	fmt.Printf("bnb eth: %f\n", c.CoinCount)
	c.InsertData()
}

func (c *Coin) BnbUsdt() {
	c.ChainName = "BNB"
	c.CoinKind = "USDT"
	c.UserAddress = c.Account.Account
	c.TokenAddress = config.Bnb.Usdt.Address
	c.Precision = config.Bnb.Usdt.Precision
	c.Price = config.Bnb.Usdt.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.BnbTokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}
	c.InsertData()
}

func (c *Coin) BnbUsdc() {
	c.CoinKind = "USDC"
	c.TokenAddress = config.Bnb.Usdc.Address
	c.Precision = config.Bnb.Usdc.Precision
	c.Price = config.Bnb.Usdc.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.BnbTokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) BnbTrias() {
	c.CoinKind = "Trias"
	c.TokenAddress = config.Bnb.Trias.Address
	c.Precision = config.Bnb.Trias.Precision
	c.Price = config.Bnb.Trias.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.BnbTokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

// ====================== Polygon ======================

func (c *Coin) PolygonEth() {
	c.ChainName = "Polygon"
	c.CoinKind = "wETH"
	c.ApiToken = config.Polygon.ApiKey
	c.UserAddress = c.Account.Account
	c.TokenAddress = config.Polygon.Eth.Address
	c.Precision = config.Polygon.Eth.Precision
	c.Price = config.Polygon.Eth.Price
	c.BlockHeight = config.Polygon.BlockHeight
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) PolygonUsdt() {
	c.CoinKind = "USDT"
	c.TokenAddress = config.Polygon.Usdt.Address
	c.Precision = config.Polygon.Usdt.Precision
	c.Price = config.Polygon.Usdt.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) PolygonUsdc() {
	c.ChainName = "Polygon"
	c.CoinKind = "USDC"
	c.BlockHeight = 59253900
	//c.UserAddress = "0xA4D8c89f0c20efbe54cBa9e7e7a7E509056228D9"
	c.UserAddress = c.Account.Account
	c.ApiToken = config.Polygon.ApiKey
	c.TokenAddress = config.Polygon.Usdc.Address
	c.Precision = config.Polygon.Usdc.Precision
	c.Price = config.Polygon.Usdc.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	//// 检查是否已经存在
	//if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
	//	return
	//}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}
	fmt.Println("Polygon Precision", c.Precision)
	fmt.Println("Polygon TokenAddress", c.TokenAddress)
	fmt.Println("Polygon UserAddress", c.UserAddress)
	fmt.Println("c.CoinCount:", c.CoinCount)
	fmt.Println("c.MoneyValue:", c.MoneyValue)

	c.InsertData()
}

func (c *Coin) PolygonBnb() {
	c.ChainName = "Polygon"
	c.CoinKind = "BNB"
	c.ApiToken = config.Polygon.ApiKey
	//c.UserAddress = "0xbCF547870155E73400B2Ba7056D4357F0550Dd8C"
	c.UserAddress = c.Account.Account
	c.BlockHeight = 58562175
	c.TokenAddress = config.Polygon.Bnb.Address
	c.Precision = config.Polygon.Bnb.Precision
	c.Price = config.Polygon.Bnb.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Polygon bnb: %f\n", c.CoinCount)
	c.InsertData()
}

func (c *Coin) PolygonTrias() {
	c.CoinKind = "Trias"
	c.TokenAddress = config.Polygon.Trias.Address
	c.Precision = config.Polygon.Trias.Precision
	c.Price = config.Polygon.Trias.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

// ====================== Optimism ======================

func (c *Coin) OptimismMain() {
	c.ChainName = "Optimism"
	c.CoinKind = "mainCoin"
	c.ApiToken = config.Optimism.ApiKey
	c.UserAddress = c.Account.Account
	c.TokenAddress = config.Optimism.Eth.Address
	c.Precision = config.Optimism.Eth.Precision
	c.Price = config.Optimism.Eth.Price
	c.BlockHeight = config.Optimism.BlockHeight
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.Balance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}
	fmt.Println("Optimism mainCoin c.CoinCount: ", c.CoinCount)
	c.InsertData()
}

func (c *Coin) OptimismUsdt() {
	c.CoinKind = "USDT"
	c.UserAddress = c.Account.Account
	c.TokenAddress = config.Optimism.Usdt.Address
	c.Precision = config.Optimism.Usdt.Precision
	c.Price = config.Optimism.Usdt.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

func (c *Coin) OptimismUsdc() {
	c.CoinKind = "USDC"
	c.UserAddress = c.Account.Account
	c.TokenAddress = config.Optimism.Usdc.Address
	c.Precision = config.Optimism.Usdc.Precision
	c.Price = config.Optimism.Usdc.Price
	c.CoinCount = 0
	c.MoneyValue = 0

	// 检查是否已经存在
	if utils.CheckExistWithDraw(c.UserAddress, c.ChainName, c.CoinKind) {
		return
	}

	err := c.TokenBalance()
	if err != nil {
		log.Printf("get user `%s` `%s` `%s` error: %v", c.UserAddress, c.ChainName, c.CoinKind, err.Error())
		os.Exit(1)
	}

	c.InsertData()
}

// WithDrawCount  统计账号信息
func WithDrawCount() {
	for _, account := range utils.GetAll() {
		//fmt.Printf("account:%s star:%d\n", account.Account, account.Star)
		//account.Account = "0xBAFD4E08103faE3ABA7C64c6eB5533304B930d34"
		coin := &Coin{Account: account}

		//// Ethereum
		//coin.EthMain()
		//coin.EthUsdt()
		//coin.EthUsdc()
		//coin.EthBnb()
		//coin.EthTrias()
		//
		//// Arbitrum
		//coin.ArbMain()
		//coin.ArbUsdt()
		//coin.ArbUsdc()
		//
		//// BNB Chain
		coin.BnbMain()
		//coin.BnbUsdt()
		//coin.BnbEth()

		//coin.BnbUsdc()
		//coin.BnbTrias()
		//
		// Polygon
		//coin.PolygonEth()
		//coin.PolygonUsdt()
		//coin.PolygonUsdc()
		//coin.PolygonTrias()
		//coin.PolygonBnb()

		// Optimism
		//coin.OptimismMain()
		//coin.OptimismUsdt()
		//coin.OptimismUsdc()
	}

}

// Census 结算每个用户的资产
func Census() {
	var accounts []models.WithDraw

	utils.DB.Debug().Model(&accounts).Distinct("account").Pluck("account", &accounts)

	for _, account := range accounts {
		var withAccounts []models.WithDraw
		var countValues float64
		utils.DB.Debug().Model(&withAccounts).Where("account = ?", account.Account).Find(&withAccounts)
		for _, withAccount := range withAccounts {
			if withAccount.Value > 0 {
				countValues += withAccount.Value
			}

		}

		utils.WithDrawCount(withAccounts[0], countValues)
		//if countValues >= 1 {
		//	utils.WithDrawCount(withAccounts[0], countValues)
		//
		//	fmt.Println("countValues:", countValues)
		//}

	}
}
