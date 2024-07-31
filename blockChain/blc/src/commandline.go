package src

import (
	"fmt"
	"os"

	"github.com/zhangyy27/docs/blockChain/blc/utils"
)

type CLI struct {
	Bc *Blockchain
}

// 创建区块
func (cli *CLI) createBlockChain() {
	if utils.FileIsExist(dbFile) {
		fmt.Println("数据库已经存在")
		return
	}

	NewBlockchain()
}

func (cli *CLI) deleteBlockChainDB() {
	if err := os.RemoveAll(dbFile); err != nil {
		fmt.Printf("remove `%s` failed. %v", dbFile, err)
		os.Exit(1)
	}
	fmt.Printf("remove `%v` success", dbFile)
}

func (cli *CLI) send(from, to string, amount float64) {
	bc := GetBlockchainInstance()
	coinbaseTx := NewCoinbaseTx(minerAddress, "")

	txs := []*Transaction{coinbaseTx}

	tx := NewTransaction(from, to, amount, bc)
	if tx != nil {
		txs = append(txs, tx)
	}
	//

	err := bc.AddBlock(txs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 输出区块链
func (cli *CLI) printChain() {
	cli.Bc = GetBlockchainInstance()
	it := cli.Bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("cur_Hash: %x\n", block.Hash)
		fmt.Println("blockData:", block.Transactions[0].TXInputs[0].ScriptSig)

		if len(it.currentHash) == 0 {
			break
		}
	}
}

// 获取账户余额
func (cli *CLI) getBalance(address string) {
	bc := GetBlockchainInstance()
	utxoInfos := bc.FindUTXO(GetPubKeyHashFromAddress(address))
	total := 0.0

	for _, utxo := range utxoInfos {
		total += utxo.TXOutputs.Value
	}

	fmt.Printf("address: %v balance total: %f", address, total)
}

func (cli *CLI) createWallet() {
	wm := NewWalletManage()
	address := wm.CreateWallet()
	if len(address) == 0 {
		fmt.Println("create wallet failed")
	}

	fmt.Println(address)
}
