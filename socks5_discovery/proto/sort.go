package proto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type Socks5List struct {
	Socks5 []*Socks5
}

func (x *Socks5List) Len() int {
	return len(x.Socks5)
}

func (x *Socks5List) Less(i, j int) bool {
	return x.Socks5[i].Delay < x.Socks5[j].Delay
}

func (x *Socks5List) Swap(i, j int) {
	x.Socks5[i], x.Socks5[j] = x.Socks5[j], x.Socks5[i]
}

func (s *Socks5) ToToken(aesKey string) string {
	marshal, _ := json.Marshal(s)

	encrypt, err := AESEncrypt(marshal, []byte(aesKey))
	if err != nil {
		panic(err)
	}
	return encrypt
}

func (s *Socks5) FromToken(token string, aesKey string) {
	decrypt, err := AESDecrypt(token, []byte(aesKey))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(decrypt, s)
	if err != nil {
		panic(err)
	}
}

func AESEncrypt(src, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	src = PKCS5Padding(src, block.BlockSize())

	cipherText := make([]byte, len(src))
	blockMode.CryptBlocks(cipherText, src)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func AESDecrypt(src string, key []byte) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])

	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = PKCS5UnPadding(origData)

	return origData, nil
}

func PadKeyString(s string, targetLen int) string {
	if len(s) < targetLen {
		// 字符串长度不足，需要填充
		return s + strings.Repeat(" ", targetLen-len(s))
	} else if len(s) > targetLen {
		// 字符串长度超过，需要截取
		return s[:targetLen]
	} else {
		// 字符串长度正好
		return s
	}
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
