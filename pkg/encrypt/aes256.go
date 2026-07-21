package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

// AES-256-CBC 加密
func AES256Encrypt(plaintext, key string) (string, error) {

	// 根据 key，初始化 block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 计算填充字节
	padding := len(key) - len(plaintext)%len(key)
	if padding == 0 {
		padding = len(key)
	}

	// 填充到原文
	plaintext += string(bytes.Repeat([]byte{byte(padding)}, padding))

	// 创建随机 iv 向量
	ciphercode := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphercode[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// 计算密文，密文赋值到 ciphercode 的 iv 值后
	cipher.NewCBCEncrypter(block, iv).
		CryptBlocks(ciphercode[aes.BlockSize:], []byte(plaintext))

	// 密文 base64 编码
	return base64.StdEncoding.EncodeToString(ciphercode), nil
}

// AES-256-CBC 解密
func AES256Decrypt(ciphertext, key string) (string, error) {

	// 根据 key，初始化 block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 密文 base64 解码
	ciphercode, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 获取 iv，密文
	iv := ciphercode[:aes.BlockSize]
	ciphercode = ciphercode[aes.BlockSize:]

	cipher.NewCBCDecrypter(block, iv).
		CryptBlocks(ciphercode, ciphercode)

	plaintext := string(ciphercode)

	return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])], nil
}
