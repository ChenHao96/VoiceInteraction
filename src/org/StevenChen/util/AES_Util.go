package util

import (
	"bytes"
	"crypto/cipher"
	"crypto/aes"
	"fmt"
	"crypto/md5"
)

func main() {

	secret := "37a0aa9075ec500d"
	api_key := "89b298ba45ec4a479dd9f20076d82b81"
	cmd := "你叫什么"
	data := "{\"key\":\"" + api_key + "\",\"info\":\"" + cmd + "\"}"
	timestamp := "1473140577124"

	keyParam := secret + timestamp + api_key;
	fmt.Println("keyParam:\n", keyParam)

	key := md5.Sum([]byte(keyParam))
	fmt.Println("MD5-Key:\n", fmt.Sprintf("%x", key))
	key2 := make([]byte, len(key))
	for index, value := range key {
		key2[index] = value
	}

	data2, _ := AesEncrypt([]byte(data), key2)
	fmt.Println("AES-MD5-data:\n", data2)
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length - 1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length - 1])
	return origData[:(length - unpadding)]
}