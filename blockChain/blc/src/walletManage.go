package src

import (
	"encoding/json"
	"fmt"
	"github.com/zhangyy27/docs/blockChain/blc/utils"
	"io/ioutil"
	"os"
)

const WalletFile = "./db/wallet.dat"

type WalletManage struct {

	// key address
	// value wallet 公/私钥
	Wallets map[string]*Wallet
}

func NewWalletManage() *WalletManage {
	//创建一个, Wallets map[string]*wallet
	var wm WalletManage

	//分配空间，一定要分配，否则没有空间
	//wm.Wallets = make(map[string]*Wallet)
	wm.Wallets = make(map[string]*Wallet)
	// 加载再有的钱包
	if !wm.loadFile() {
		return nil
	}
	return &wm
}

// CreateWallet 创建钱包
func (w *WalletManage) CreateWallet() string {
	// 创建密钥对
	wallet := NewWallet()
	if wallet == nil {
		return ""
	}

	// 获取地址
	address := wallet.getAddress()
	w.Wallets[address] = wallet

	// 将秘钥对写入磁盘
	if !w.saveWalletFile() {
		return ""
	}

	// 返回给cli新地址
	return address
}

func (w *WalletManage) saveWalletFile() bool {

	jsonData, err := json.Marshal(w)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return false
	}

	err = ioutil.WriteFile(WalletFile, jsonData, 0600)
	if err != nil {
		fmt.Println("ioutil.WriteFile err", err)
		return false
	}

	return true
}

// 读取wallet.dat文件，加载wm中
func (w *WalletManage) loadFile() bool {
	// 判断文件是否存在
	if !utils.FileIsExist(WalletFile) {
		fmt.Println("文件不存在,无需加载!")
		return true
	}

	// 读取文件
	content, err := ioutil.ReadFile(WalletFile)
	if err != nil {
		fmt.Println("ioutil.ReadFile err:", err)
		os.Exit(1)
		return false
	}
	//
	//fmt.Println("content:", content)
	//
	//gob.Register(elliptic.P256())
	//decoder := gob.NewDecoder(bytes.NewReader(content))
	//err = decoder.Decode(w)
	//if err != nil {
	//	fmt.Println("decoder.Decode err:", err)
	//	os.Exit(1)
	//}

	// 反序列化
	err = json.Unmarshal(content, w)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		os.Exit(1)
		return false
	}

	return true
}
