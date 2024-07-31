package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/zhangyy27/docs/blockChain/blc/utils"
	"log"
	"os"
	"time"
)

type Block struct {
	Version       uint64
	PrevBlockHash []byte
	Hash          []byte
	Data          []byte
	Transactions  []*Transaction
	MerkleRoot    []byte
	Nonce         uint64
	Timestamp     uint64
	Bits          uint64
}

func (b *Block) SetHash() {
	tmp := [][]byte{
		utils.Uint64ToHex(b.Version),
		b.PrevBlockHash,
		b.Hash,
		//b.Transactions.Serialize(),
		b.MerkleRoot,
		utils.Uint64ToHex(b.Nonce),
		utils.Uint64ToHex(b.Timestamp),
		utils.Uint64ToHex(b.Bits),
	}

	blockBytes := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}

func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
		return nil
	}
	return buffer.Bytes()
}

// HashTransactionsMerkleRoot 默克尔树
func (b *Block) HashTransactionsMerkleRoot() {

	var info [][]byte
	for _, tx := range b.Transactions {
		txHashValue := tx.TXId
		info = append(info, txHashValue)
	}

	infoHashValue := sha256.Sum256(bytes.Join(info, []byte{}))
	b.MerkleRoot = infoHashValue[:]
}

func DeSerialize(data []byte) *Block {
	var block Block
	var buffer bytes.Buffer

	_, err := buffer.Write(data)
	if err != nil {
		os.Exit(1)
	}

	decoder := gob.NewDecoder(&buffer)
	if err = decoder.Decode(&block); err != nil {

		os.Exit(1)
	}

	return &block
}

func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Version:       1,
		Transactions:  txs,
		PrevBlockHash: prevBlockHash,
		Nonce:         0,
		Timestamp:     uint64(time.Now().Unix()),
		MerkleRoot:    nil,
		Bits:          0,
		Hash:          nil,
	}
	block.HashTransactionsMerkleRoot() // 设置merkleRoot
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}
