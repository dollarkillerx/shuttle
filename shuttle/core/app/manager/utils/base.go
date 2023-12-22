package utils

import (
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GetDefaultHttpClient(timeout int) *http.Client {
	if timeout == 0 {
		timeout = 3
	}
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func GenerateRandNum(digits int) string {
	rand.Seed(time.Now().UnixNano())                // 用当前时间设置随机数种子
	min := int(math.Pow10(digits - 1))              // 最小值
	max := int(math.Pow10(digits)) - 1              // 最大值
	return strconv.Itoa(rand.Intn(max-min+1) + min) // 生成随机数并转成字符串返回
}
