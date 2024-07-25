package blc_demo

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// 每个区块的奖励金额
const subsidy = 10

// TXInput 		表示交易输入
type TXInput struct {
	Txid      []byte // 引用的交易ID，表示交易hash
	Vout      int    // 引用的输出索引
	Scriptsig string // 解锁脚本
}

// TXOutput 		表示交易输出
type TXOutput struct {
	Value        int    // 输出金额
	ScriptPubKey string // 解定脚本
}

// IsCoinbase 将检查交易是否为coinbase交易
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID 设置交易ID
func (tx *Transaction) SetID() {
	// 创建一个缓冲区
	var encoded bytes.Buffer

	// 创建一个hash
	var hash [32]byte

	// 创建一个编码器
	enc := gob.NewEncoder(&encoded)
	if err := enc.Encode(tx); err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Transaction 	表示一笔交易
type Transaction struct {
	ID   []byte     //交易ID 表示交易的hash
	Vin  []TXInput  // 交易输入
	Vout []TXOutput // 交易输出
}

// CanUnlockOutputWith 将检查输出是否可以使用提供的数据解锁
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.Scriptsig == unlockingData
}

// CanBeUnlockedWith 将检查输出是否可以使用提供的数据解锁
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX 创建一个coinbase交易
func NewCoinbaseTX(to, data string) *Transaction {
	// 如果没有数据， 则使用黑夜数据
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	// coinbase交易没有输入，所以Txid为空, Vout为-1
	txin := TXInput{[]byte{}, -1, data}

	// 创建一个输出
	txout := TXOutput{subsidy, to}

	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

// NewUTXOTransaction 创建一个新的UTXO交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	// 获取未花费输出
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	// 判断是否有足够的金额
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// 构建一个输入列表
	for txid, outs := range validOutputs {
		// 交易id转为字节数组
		txId, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		// 遍历输出
		for _, out := range outs {
			// 构建一个输入
			input := TXInput{txId, out, from}

			// 添加到输入列表
			inputs = append(inputs, input)
		}
	}

	// 构建一个输出
	outputs = append(outputs, TXOutput{amount, to})

	// 如果还有剩余金额，则添加一个输出（找零）
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	// 创建一个交易
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}

// IntToHex 		将整数转换为16进制
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
