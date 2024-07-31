package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

func EccTest() {
	// 生成曲线
	curve := elliptic.P256()
	// 生成私钥
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	// 私钥得到公钥
	pubKey := &priKey.PublicKey

	// 私钥签名
	hash := []byte("hello world")
	r, s, err := ecdsa.Sign(rand.Reader, priKey, hash)
	if err != nil {
		panic(err)
	}

	// 将r, s 转为64进制
	signaTure := append(r.Bytes(), s.Bytes()...)
	fmt.Println("signaTure:", signaTure)

	// 验证
	if !ecdsa.Verify(pubKey, hash, r, s) {
		panic("verify failed")
	}
	fmt.Println("verify success")

}
