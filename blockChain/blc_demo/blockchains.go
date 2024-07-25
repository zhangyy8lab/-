package blc_demo

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = "blockchain.db"                                                                      // 数据库文件
const blockBucket = "blocks"                                                                        // 存储区块的桶
const genesisCoinbaseData = "The Times 03/jan/2009 Chancellor on brink of second bailout for banks" // 创世块的交易数据

// Blockchain 包含一个区块链
type Blockchain struct {
	tip []byte   // 最新区块的hash
	db  *bolt.DB // 数据库指针
}

// BlockchainIterator 将用于区块链迭代
type BlockchainIterator struct {
	currentHash []byte   // 当前区块的hash
	db          *bolt.DB // 数据库指针
}

// MineBlock 将用于挖掘⛏新块️
func (bc *Blockchain) MineBlock(transaction []*Transaction) {
	var lastHash []byte

	// 获取最新区块的hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// 创建新区块(包含验证)
	newBlock := NewBlock(transaction, lastHash)

	// 将新区块存储到数据库中
	err = bc.db.Update(func(tx *bolt.Tx) error {
		// 获取区块桶
		b := tx.Bucket([]byte(blockBucket))

		// 将区块存储到数据库中
		err = b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		// 将最新区块的hash存储到数据库中
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		// 更新区块链的tip
		bc.tip = newBlock.Hash
		return nil
	})
}

// FindUnspentTransactions 返回未花费的交易
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	// 未花费的交易
	var unspentTXs []Transaction

	// 已花费的交易
	spentTXOs := make(map[string][]int)

	// 迭代区块链
	bci := bc.Iterator()
	for {
		//获取下一个区块链
		block := bci.Next()

		// 遍历当前区块中的交易
		for _, tx := range block.Transactions {
			// 将交易id转为str
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			// 遍历交易中的输出
			for outIdx, out := range tx.Vout {
				// 检查输出是否已经被花费
				if spentTXOs[txID] != nil {
					// 遍历已花费的输出
					for _, spentOut := range spentTXOs[txID] {
						// 如果输出已经被花费， 则跳过
						if spentOut == outIdx {
							continue Outputs
						}

					}

				}

				// 如果输出可以被解锁， 则将交易添加到未花费的交易中
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// 如果交易不是coinbase交易，则遍历交易的输出
			if tx.IsCoinbase() == false {
				// 遍历交易的输入
				for _, in := range tx.Vin {
					// 如果输入可以解锁， 则将输出添加到已花费的输出中
					if in.CanUnlockOutputWith(address) {
						// 将交易ID转换为字符串
						inTxID := hex.EncodeToString(in.Txid)
						// 将输出添加到已民花费的输出中
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}

		}

		// 如果区块的前一个区块hash为空， 则停止迭代
		if len(block.PrevBlockHash) == 0 {
			break
		}

	}

	// 返回未花费的交易
	return unspentTXs

}

// FindUTXO 返回未花费的输出
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput

	// 未花费的交易
	unspentTransactions := bc.FindUnspentTransactions(address)

	// 遍历交易的输出
	for _, tx := range unspentTransactions {

		// 遍历交易的输出
		for _, out := range tx.Vout {
			// 如果输出可以被解锁， 则将输出添加到未花费的输出中
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs 返回足够未花费输出以满足要求的金额
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	// 未花费的输出
	unspentOutputs := make(map[string][]int)

	// 未花费的交易
	unspentTXs := bc.FindUnspentTransactions(address)

	// 累计金额
	accumulated := 0

Work:
	// 遍历未花费的交易
	for _, tx := range unspentTXs {
		// 将交易ID转换为字符串
		txID := hex.EncodeToString(tx.ID)

		// 遍历交易的输出
		for outIdx, out := range tx.Vout {
			// 如果输出可以被解锁且累计金额小于要求的金额， 则将输出添加到未花费的输出中
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				// 输出的索引添加到未花费的输出中
				accumulated += out.Value

				// 输出的索引添加到未花费的输出中
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				// 如果累计金额大于等于要求的金额， 则停止遍历
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs

}

// Iterator 返回一个迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator {
	// 迭代器对象是一个指向区块链的指针和一个指向数据库的指针
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

// Next 返回区块链的一个区块
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	// 获取当前区块
	err := i.db.View(func(tx *bolt.Tx) error {
		// 获取区块桶
		b := tx.Bucket([]byte(blockBucket))

		// 获取当前区块
		encodedBlock := b.Get(i.currentHash)

		// 返序列化区块
		block = DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// 更新当前区块
	i.currentHash = block.PrevBlockHash
	return block
}

// dbExists 检查数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// NewBlockchain 创建一个新的区块
func NewBlockchain(address string) *Blockchain {
	// 判断数据库是否存在
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	// 存储最新区块的hash
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// 更新数据库
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// 创建一个区块链对象
	bc := Blockchain{tip, db}
	return &bc
}

// CreateBlockchain 创建一个新的区块链
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// 更新数据库
	err = db.Update(func(tx *bolt.Tx) error {
		// 创建一个新的coinbase交易
		cbTx := NewCoinbaseTX(address, genesisCoinbaseData)

		// 创建一个新的区块
		genesis := NewGenesisBlack(cbTx)

		// 创建一个新的桶
		b, err := tx.CreateBucket([]byte(blockBucket))
		if err != nil {
			log.Panic(err)
		}

		// 将区块存储到数据库中
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{
		tip,
		db,
	}
	return &bc
}

// AddBlock 将用于添加区块到区块链
func AddBlock(address string) {
	// 创建一个新的区块

	bc := NewBlockchain(address)
	defer bc.db.Close()

	// 创建一个新的coinbase交易
	//cbTx := NewCoinbaseTX(address, "")

	// 挖掘新块并存放到数据库中
	bc.MineBlock([]*Transaction{NewCoinbaseTX(address, "")})
	fmt.Println("Coinbase交易完成，以奖励的方式将硬币发送到地址：", address)
	fmt.Println("Success!")
}
