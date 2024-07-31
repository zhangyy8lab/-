package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
)

const (
	//usdtContractAddress = "0x55d398326f99059fF775485246999027B3197955" // USDT合约地址
	rpcURL = "https://bsc-dataseed.binance.org/" // BSC节点URL
)

func main() {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// 用户地址
	userAddress := common.HexToAddress("0xF977814e90dA44bFA03b6295A0616a897441aceC")
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

	// 构造调用消息
	contractAddress := common.HexToAddress("0x55d398326f99059fF775485246999027B3197955")
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	// 调用合约
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	// 解析返回值
	var balance = new(big.Int)
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		log.Fatalf("Failed to unpack result: %v", err)
	}

	fmt.Printf("USDT Balance at block %d: %s\n", "111111", balance.String())
}
