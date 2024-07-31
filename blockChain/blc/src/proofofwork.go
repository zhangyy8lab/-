package src

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/zhangyy27/docs/blockChain/blc/utils"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	pow := ProofOfWork{block: b}
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	bitIntTmp := big.Int{}
	bitIntTmp.SetString(targetStr, 16)
	pow.target = &bitIntTmp
	return &pow
}

func (pow *ProofOfWork) prepareData(num uint64) []byte {
	block := pow.block
	tmp := [][]byte{
		utils.Uint64ToHex(block.Version),
		block.PrevBlockHash,
		//block.Data,
		block.MerkleRoot,
		utils.Uint64ToHex(num),
		utils.Uint64ToHex(block.Timestamp),
		utils.Uint64ToHex(block.Bits),
	}

	data := bytes.Join(tmp, []byte{})
	return data
}

// Run 进行挖矿进行工作量证明
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var Nonce uint64
	var hash [32]byte

	for {
		hash = sha256.Sum256(pow.prepareData(Nonce))
		tTmp := big.Int{}
		tTmp.SetBytes(hash[:])
		fmt.Printf("hash:%x Nonce:%d\n", hash, Nonce)
		if tTmp.Cmp(pow.target) == -1 {
			break
		} else {

			//
			Nonce++
		}
	}
	return hash[:], Nonce
}

func (pow *ProofOfWork) IsValid() bool {
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	tmpInt := new(big.Int)
	tmpInt.SetBytes(hash[:])
	if tmpInt.Cmp(pow.target) == -1 {
		return true
	}
	return false
}
