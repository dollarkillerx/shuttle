package utils

import (
	"github.com/bwmarrin/snowflake"

	"math/rand"
	"sync"
	"time"
)

var sf *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	sf = node
}

func GenerateID() string {
	return sf.Generate().String()
}

func XRand(max int, min int) int {
	rand.Seed(time.Now().UnixNano())
	p := rand.Intn(max)
	if p > min {
		return p
	}

	return p + min
}

const randKey = "QWERTYUIOPASDFGHJKLZXCVBNM1234567890QWERTYUIOPASDFGHJKLZXCVBNM1234567890QWERTYUIOPASDFGHJKLZXCVBNM1234567890"

var gMut sync.Mutex

func RandKey(ln int) string {
	gMut.Lock()
	defer gMut.Unlock()

	if ln == 0 {
		ln = 9
	}
	var result string
	pln := len(randKey)
	for i := 0; i < ln; i++ {
		time.Sleep(time.Nanosecond * 250)
		rand.Seed(time.Now().UnixNano())
		result += string(randKey[rand.Intn(pln)])
	}

	return result
}

const randSms = "ASDFGHJKLQWERTYUOPZXCVBNM0123456798ASDFGHJKLQWERTYUOPZXCVBNM0123456798"

func RandSMSCode(ln int) string {
	gMut.Lock()
	defer gMut.Unlock()

	if ln == 0 {
		ln = 9
	}
	var result string
	pln := len(randSms)
	for i := 0; i < ln; i++ {
		time.Sleep(time.Nanosecond * 250)
		rand.Seed(time.Now().UnixNano())
		result += string(randSms[rand.Intn(pln)])
	}

	return result
}
