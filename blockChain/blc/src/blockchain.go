package src

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/zhangyy27/docs/blockChain/blc/utils"
)

const genesisData = "The Times 03/jan/2009 Chancellor on brink of second bailout for banks"
const dbFile = "../db/blockchain.db"
const blockBucket = "blockBucket"
const lastBlockHashKey = "lastBlockHashKey"
const minerAddress = "minerAddress"

type Blockchain struct {
	db       *bolt.DB
	lastHash []byte
}

type Iterator struct {
	db          *bolt.DB
	currentHash []byte
}

// UTXOInfo 定义
type UTXOInfo struct {
	// TxId
	TxId []byte
	// index
	Index     int64
	TXOutputs TXOutput
}

func (bc *Blockchain) AddBlock(txs []*Transaction) error {

	prevBlock := bc.lastHash
	newBlock := NewBlock(txs, prevBlock)

	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		_ = bucket.Put(newBlock.Hash, newBlock.Serialize())
		_ = bucket.Put([]byte(lastBlockHashKey), newBlock.Hash)
		bc.lastHash = newBlock.Hash
		return nil
	})

	if err != nil {
		return errors.New(fmt.Sprintf("add new block failed. %v", err.Error()))
	}
	fmt.Printf("NewBlock Hash: %x\n", newBlock.Hash)
	//fmt.Println("NewBlock Hash:", string(newBlock.Hash))
	return err
}

func (b *Block) toBytes() []byte {
	var result []byte

	return result
}

// NewIterator 区块链迭代器
func (bc *Blockchain) NewIterator() *Iterator {
	return &Iterator{
		db:          bc.db,
		currentHash: bc.lastHash, // 当前最新hash
	}
}

// Next 获取下一个区块
func (it *Iterator) Next() *Block {
	var b *Block

	_ = it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket)) // 获取桶
		if bucket == nil {
			return errors.New(fmt.Sprintf("%s not exist", blockBucket))
		}

		currentBlockBytes := bucket.Get(it.currentHash) // 得到当前块hash

		b = DeSerialize(currentBlockBytes) // 反序列化到区块对象

		it.currentHash = b.PrevBlockHash // 迭代器游标向前移动一个
		return nil
	})

	return b
}

// FindUTXO 找到未花费的交易
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []UTXOInfo {
	var utxoInfos []UTXOInfo

	// 定义一个存在已经消耗过的所有的utxos的集合
	spendUtxos := make(map[string][]int)
	it := bc.NewIterator()

	for {
		// 获取区块
		block := it.Next()

		//遍历交易
		for _, tx := range block.Transactions {
		LABEL:
			//1. 遍历output，判断这个output的锁定脚本是否为我们的目标地址
			for outputIndex, output := range tx.TXOutputs {
				// LABEL:
				fmt.Println("outputIndex:", outputIndex)

				//这里对比的是哪一些utxo与付款人有关系
				// if output.ScriptPubKeyHash /*某一个被公钥哈希锁定output*/ == pubKeyHash /*张三的哈希*/ {
				if bytes.Equal(output.ScriptPubKeyHash, pubKeyHash) {

					//开始过滤
					//当前交易id
					currentTxid := string(tx.TXId)
					//去spentUtxos中查看
					indexArray := spendUtxos[currentTxid]

					//如果不为零，说明这个交易id在篮子中有数据，一定有某个output被使用了
					if len(indexArray) != 0 {
						for _, spendIndex /*0, 1*/ := range indexArray {
							//接着判断下标
							if outputIndex /*当前的*/ == spendIndex {
								continue LABEL
							}
						}
					}

					//找到属于目标地址的output
					// utxos = append(utxos, output)
					utxoinfo := UTXOInfo{tx.TXId, int64(outputIndex), output}
					utxoInfos = append(utxoInfos, utxoinfo)
				}

			}

			//++++++++++++++++++++++遍历inputs+++++++++++++++++++++
			if tx.isCoinbaseTx() {
				//如果是挖矿交易，则不需要遍历inputs
				fmt.Println("发现挖矿交易，无需遍历inputs")
				continue
			}

			for _, input := range tx.TXInputs {
				// if input.PubKey /*付款人的公钥*/ == pubKeyHash /*张三的公钥哈希*/ {
				if bytes.Equal(getPubKeyHashFromPubKey(input.PubKey), pubKeyHash) {
					//map[key交易id][]int
					//map[string][]int{
					//	0x333: {0, 1}
					//}
					spentKey := string(input.TxId)

					//向篮子中添加已经消耗的output
					spendUtxos[spentKey] = append(spendUtxos[spentKey], int(input.Index))
					// spentUtxos[0x333] =[]int{0}
					// spentUtxos[0x333] =[]int{0, 1}
					// spentUtxos[0x222] =[]int{0}

					//不要使用这种方式，否则spendUtxos不会被赋值
					// indexArray := spentUtxos[spentKey]
					// indexArray = append(indexArray, int(input.Index))
				}
			}
		}

		// 退出条件
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return utxoInfos
}

// 找到需要的utxo及剩余
func (bc *Blockchain) findNeedUtxo(from string, amount float64) (map[string][]int64, float64) {
	var retMap = make(map[string][]int64)
	var retValue float64

	utxoinfos := bc.FindUTXO(GetPubKeyHashFromAddress(from))
	for _, utxoinfo := range utxoinfos {
		// 统计总额
		retValue += utxoinfo.TXOutputs.Value

		// 统计需要的utxo
		key := string(utxoinfo.TxId)
		fmt.Println("blockchain: key:", key)
		retMap[key] = append(retMap[key], utxoinfo.Index)

		fmt.Println("11111")
		if retValue >= amount {
			break
		}

	}
	return retMap, retValue
}

// NewGenesisBlock 创世块
func NewGenesisBlock(addr string) *Block {
	//
	coinbase := NewCoinbaseTx(addr, genesisData)
	txs := []*Transaction{coinbase}
	genesisBlock := NewBlock(txs, nil)

	return genesisBlock
}

// NewBlockchain 创建区块链
func NewBlockchain() *Blockchain {
	// 1 区块链不存在，创建
	// 2 区块链存在，返回

	// 创建blockChain 同时添加genesisBlock

	// 创建blockChainDB
	// 更新， 找到目录bucket
	// bucket不存在则创建， 写入创世块
	// bucket存在， 返回最后一个块的hash

	var lastHash []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			// 创建bucket
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
				return err
			}

			// 写入创世块
			genesisBlock := NewGenesisBlock(minerAddress)
			_ = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			_ = bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)

			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}
		return nil
	})

	return &Blockchain{
		db:       db,
		lastHash: lastHash,
	}
}

// GetBlockchainInstance 获取区块链实例
func GetBlockchainInstance() *Blockchain {

	if !utils.FileIsExist(dbFile) {
		fmt.Printf("block chain db not exist\n")
		os.Exit(1)
	}

	var lastHash []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			panic("bucket is not exist")
		}

		lastHash = bucket.Get([]byte(lastBlockHashKey))

		return nil
	})

	return &Blockchain{
		db:       db,
		lastHash: lastHash,
	}
}

func (tx *Transaction) isCoinbaseTx() bool {
	inputs := tx.TXInputs
	//input个数为1，id为nil，索引为-1
	if len(inputs) == 1 && inputs[0].TxId == nil && inputs[0].Index == -1 {
		return true
	}
	return false
}
