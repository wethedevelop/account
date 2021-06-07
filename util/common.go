package util

import (
	"crypto/rand"
	"math/big"
)

// 随机函数封装
func Rand(max int) int64 {
	randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return randNum.Int64()
}

// RandStringRunes 返回随机字符串
func RandStringRunes(n int64) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[Rand(len(letterRunes))]
	}
	return string(b)
}
