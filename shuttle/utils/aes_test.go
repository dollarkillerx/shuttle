package utils

import (
	"fmt"
	"testing"

	"github.com/rs/xid"
)

func TestAESDecrypt(t *testing.T) {

	key := PadKeyString(xid.New().String(), 32)
	fmt.Println(key)
	src := []byte("Hello world! 你好呀 我爱你")

	s1, err := AESEncrypt(src, []byte(key))
	if err != nil {
		panic(err)
	}
	fmt.Println("cipher text:", s1)

	s2, err := AESDecrypt(s1, []byte(key))
	if err != nil {
		panic(err)
	}
	fmt.Println("original text:", string(s2))
}
