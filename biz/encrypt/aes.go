package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
)

// @brief:填充明文
func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// @brief:去除填充数据
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// @brief:AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	b64 := base64.StdEncoding.EncodeToString(origData)
	origData = []byte(b64)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// @brief:AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	b64, err := base64.StdEncoding.DecodeString(string(origData))
	if err != nil {
		return nil, err
	}
	return b64, nil
}

func EncrptAESHEX(origData, key []byte) (string, error) {
	s, err := AesEncrypt(origData, key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(s), nil
}

func DecrptAESHEX(origData string, key []byte) (string, error) {
	des, err := hex.DecodeString(origData)
	if err != nil {
		return "", err
	}
	s, err := AesDecrypt(des, key)
	if err != nil {
		return "", err
	}
	return string(s), nil
}
