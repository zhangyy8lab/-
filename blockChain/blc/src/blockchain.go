package v1

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const genesisData = "Genesis Block"
const dbFile = "blockchain.db"
const blockBucket = "blockBucket"
const lastBlockHashKey = "lastBlockHashKey"

type Blockchain struct {
	db       *bolt.DB
	lastHash []byte
}

type Iterator struct {
	db          *bolt.DB
	currentHash []byte
}

func (bc *Blockchain) AddBlock(data []byte) error {

	prevBlock := bc.lastHash
	newBlock := NewBlock(data, prevBlock)

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

// NewGenesisBlock 创世块
func NewGenesisBlock() *Block {
	genesisBlock := NewBlock([]byte(genesisData), nil)

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
			genesisBlock := NewGenesisBlock()
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

	if !dbExist(dbFile) {
		fmt.Printf("block chain db not exist")
		return NewBlockchain()

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

// 判断db是否存在
func dbExist(dbPath string) bool {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}

	return true
}
