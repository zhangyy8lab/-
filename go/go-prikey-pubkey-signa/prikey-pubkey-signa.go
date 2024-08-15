package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func main() {
	// 生成 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("密钥对生成失败:", err)
		return
	}

	// 原始文本
	plaintext := []byte("Hello, World!")

	fmt.Println("plaintext:", plaintext)
	
	// 加密文本
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privateKey.PublicKey, plaintext, nil)
	if err != nil {
		fmt.Println("加密失败:", err)
		return
	}
	fmt.Println("ciphertext:", ciphertext)

	// 解密文本
	decryptedText, err := privateKey.Decrypt(nil, ciphertext, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}

	fmt.Println("解密后的文本:", string(decryptedText))

	// 对文本进行签名
	hashed := sha256.Sum256(plaintext)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Println("签名失败:", err)
		return
	}

	// 验证签名
	err = rsa.VerifyPKCS1v15(&privateKey.PublicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Println("签名验证失败:", err)
		return
	}

	fmt.Println("签名验证成功")
}

