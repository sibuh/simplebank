package util

import (
	"math/rand"
	"strings"
	"time"
)

// random functions created here
var alphabet = "qwertyuiopasdfghjklzxcvbnm"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(max, min int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]

		sb.WriteByte(c)
	}
	return sb.String()
}
func RandomOwner() string {
	return RandString(6)
}
func RandomMoney() int64 {
	return RandInt(1000, 1)
}
func RandomCurrency() string {
	currency := []string{usd, euro, birr, shilng}
	n := len(currency)
	return currency[rand.Intn(n)]
}
func RandomEmail() string {
	return RandString(6) + "@gmail.com"
}
