package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

var reward = 12.5

type TXInputs struct {
	TxId  []byte
	Index int64 // 引用
	//ScriptSig string // 付款人对当前交易的签名
	ScriptSig []byte // 付款人对当前交易的签名
	PubKey    []byte // 付款人的公钥
}

type TXOutput struct {
	ScriptPubKeyHash []byte // 收款人的公钥哈希，
	//ScriptPubKey string  // 收款人的公钥哈希， 可理解为收款人地址
	Value float64 // 转赂金额
}

type Transaction struct {
	TXId      []byte     // 交易id
	TXInputs  []TXInputs // 输入
	TXOutputs []TXOutput // 输出
	TimeStamp uint64
}

// 没有办法直接将地址赋值给TXOutput, 需要提供一下output方法
func newTxOutput(address string, amount float64) TXOutput {
	return TXOutput{
		ScriptPubKeyHash: getPubKeyHashFromPubKey([]byte(address)),
		Value:            amount,
	}
}

// 给一个交易设置hash
func (tx *Transaction) setHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	_ = encoder.Encode(tx)
	hash := sha256.Sum256(buffer.Bytes())

	tx.TXId = hash[:]
	return
}

// NewCoinbaseTx 创建一个coinbase交易
func NewCoinbaseTx(miner /*挖矿人*/ string, data string) *Transaction {
	timeStamp := uint64(time.Now().Unix())
	inputs := TXInputs{
		TxId:      nil,
		Index:     -1,
		ScriptSig: nil,
		PubKey:    []byte(data),
	}
	//output := TXOutput{
	//	Value:        reward,
	//	ScriptPubKey: newTxOutput(miner),
	//}
	output := newTxOutput(miner, reward)

	tx := Transaction{
		TXId:      nil,
		TXInputs:  []TXInputs{inputs},
		TXOutputs: []TXOutput{output},
		TimeStamp: timeStamp,
	}
	tx.setHash()
	return &tx
}

func NewTransaction(from, to string, amount float64, bc *Blockchain) *Transaction {
	var spendUtxo = make(map[string][]int64)
	var retValue float64

	// 遍历账本 找到from能够使用的utxo， 以及utxo需要的钱
	spendUtxo, retValue = bc.findNeedUtxo(from, amount)

	// 判断余额是否足够
	if retValue < amount {
		fmt.Println("余额不足， 创建交易失败")
		return nil
	}

	// 够
	var inputs []TXInputs
	var outputs []TXOutput
	for txid, indexArray := range spendUtxo {

		// 遍历下标，
		for _, i := range indexArray {
			input := TXInputs{
				[]byte(txid),
				int64(i),
				[]byte(from),
				nil,
			}
			inputs = append(inputs, input)
		}

	}

	// output
	//output1 := TXOutput{to, amount}
	output1 := newTxOutput(to, amount)

	outputs = append(outputs, output1)

	// 多余额的找零， 找给自己
	if retValue > amount {
		output2 := TXOutput{[]byte(from), retValue - amount}
		outputs = append(outputs, output2)
	}

	tx := Transaction{
		TXId:      nil,
		TXInputs:  inputs,
		TXOutputs: outputs,
		TimeStamp: uint64(time.Now().Unix()),
	}
	tx.setHash()

	// 需要对输入进行签名
	// todo
	return &tx
}
