package v1

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/zhangyy27/docs/blockChain/blc/src/blockchain/utils"
	"log"
	"os"
	"time"
)

type Block struct {
	Version       uint64
	PrevBlockHash []byte
	Hash          []byte
	Data          []byte
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
		b.Data,
		b.MerkleRoot,
		utils.Uint64ToHex(b.Nonce),
		utils.Uint64ToHex(b.Timestamp),
		utils.Uint64ToHex(b.Bits),
	}

	blockBytes := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}

func NewBlock(data []byte, prevBlockHash []byte) *Block {
	block := &Block{
		Version:       1,
		Data:          data,
		PrevBlockHash: prevBlockHash,
		Nonce:         0,
		Timestamp:     uint64(time.Now().Unix()),
		MerkleRoot:    nil,
		Bits:          0,
		Hash:          nil,
	}
	//block.SetHash()

	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
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
