package blc_demo

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// CLI 处理命令行参数
type CLI struct {
}

func (cli *CLI) createBlockchain(address string) {
	// 创建区块链
	bc := CreateBlockchain(address)
	bc.db.Close()

	fmt.Println("done")
}

// 获取账户余额
func (cli *CLI) getBalance(address string) {
	// 创建区块链
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0

	// 获取地址的UTXO
	UTXOs := bc.FindUTXO(address)

	// 计算余额
	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': '%d'\n", address, balance)

}

func (cli *CLI) printUsage() {

	fmt.Println("Usage:")
	fmt.Println(" English： getbalance -address ADDRESS - Get balance of ADDRESS、、、中文：getbalance - address ADDRESS -获取地址的余额")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")

}

// 验证参数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printChain() {
	// 输入用户地址
	bc := NewBlockchain("")
	defer bc.db.Close()

	bci := bc.Iterator()

	// 循环打印区块链中的区块
	for {

		block := bci.Next()
		timeFormat := time.Unix(block.Timestamp, 0)
		fmt.Println("timestamp:", timeFormat)
		fmt.Println("PrevBlockHash:", block.PrevBlockHash)
		fmt.Println("BlockHash:", block.Hash)

		// 显示完整交易信息
		for _, tx := range block.Transactions {
			var str string
			for _, value := range tx.ID {
				str += strconv.Itoa(int(value))
			}

			fmt.Println("TransactionID:", str)
			fmt.Println("TXInput:")
			for _, txin := range tx.Vin {
				fmt.Println(txin.Txid)
				fmt.Println(txin.Vout)
				fmt.Println(txin.Scriptsig)
			}
			fmt.Println("TXOutput:")
			for _, tx2 := range tx.Vout {
				fmt.Println(tx2.Value)
				fmt.Println(tx2.ScriptPubKey)
			}
		}

		// nonce
		fmt.Println("blockNonce:", block.Nonce)

		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %s \n", strconv.FormatBool(pow.Validate()))

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success")
}

func (cli *CLI) addBlock(address string) {
	AddBlock(address)
}

func (cli *CLI) Run() {
	for {
		fmt.Println("1. getbalance -address ADDRESS - 获取地址的余额 ")
		fmt.Println("2. createblockchain -address ADDRESS - 创建区块链并将创世区块奖励发送到ADDRESS")
		fmt.Println("3. printchain - 打印区块链")
		fmt.Println("4. send -from FROM -to TO -amount AMOUNT - 从FROM地址向TO发送AMOUNT硬币")
		fmt.Println("5. mine ADDRESS挖矿创建区块链并将创世区块奖励发送到ADDRESS")
		fmt.Println("6. exit - Exit")
		fmt.Println("Please enter the command:")

		var cmd string
		switch cmd {
		case "1":
			fmt.Println("Please enter the address")
			var address string
			fmt.Scanln(&address)
			cli.getBalance(address)
			fmt.Println()
		case "2":
			fmt.Println("Please enter the address")
			var address string
			fmt.Scanln(&address)
			cli.createBlockchain(address)
			fmt.Println()
		case "3":
			cli.printChain()
		case "4":
			fmt.Println("Please enter the from address")
			var from string
			fmt.Scanln(&from)

			fmt.Println("Please enter the to address")
			var to string
			fmt.Scanln(&to)

			fmt.Println("Please enter the amount")
			var amount int
			fmt.Scanln(&amount)
			cli.send(from, to, amount)
			fmt.Println()
		case "5":
			fmt.Println("Please enter the address")
			var address string
			fmt.Scanln(&address)
			cli.addBlock(address)
		case "6":
			os.Exit(1)
		}
	}
}
