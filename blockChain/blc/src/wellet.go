package src

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	//私钥
	PriKey *ecdsa.PrivateKey
	PubKey []byte
}

// NewWallet 密钥对
func NewWallet() *Wallet {
	// 曲线
	curve := elliptic.P256()

	// 私钥
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Printf("ecdsa.GenerateKey err:%s", err.Error())
		return nil
	}

	// 公钥
	pubKeyRaw := priKey.PublicKey

	// 公钥 x, y 拼接一起
	pubKey := append(pubKeyRaw.X.Bytes(), pubKeyRaw.Y.Bytes()...)

	// 创建 wallet
	wallet := Wallet{
		PriKey: priKey,
		PubKey: pubKey,
	}
	return &wallet
}

// 给定公钥，得到公钥哈希值
func getPubKeyHashFromPubKey(pubKey []byte) []byte {
	hash1 := sha256.Sum256(pubKey)
	//hash160处理
	hasher := ripemd160.New()
	hasher.Write(hash1[:])

	// 公钥哈希，锁定output时就是使用这值
	pubKeyHash := hasher.Sum(nil)

	return pubKeyHash
}

// 得到4字节的校验码
func checkSum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	//4字节checksum
	checksum := second[0:4]
	return checksum
}

// 生成地址
func (w *Wallet) getAddress() string {
	//公钥
	// pubKey := w.PubKey
	pubKeyHash := getPubKeyHashFromPubKey(w.PubKey)

	fmt.Println("pubKeyHashLen", len(pubKeyHash))

	//拼接version和公钥哈希，得到21字节的数据
	payload := append([]byte{byte(0x00)}, pubKeyHash...)

	//生成4字节的校验码
	checksum := checkSum(payload)

	//25字节数据
	payload = append(payload, checksum...)
	address := base58.Encode(payload)
	return address
}

func GetPubKeyHashFromAddress(address string) []byte {
	//base58解码
	decodeInfo := base58.Decode(address)
	if len(decodeInfo) != 25 {
		fmt.Println("getPubKeyHashFromAddress, 传入地址无效")
		return nil
	}
	//需要校验一下地址

	//截取
	pubKeyHash := decodeInfo[1 : len(decodeInfo)-4]
	return pubKeyHash
}