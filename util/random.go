package util

import (
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// init initializes the random generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns a random integer in the range [min, max].
func RandomInt(min, max int) int64 {
	return int64(min + rand.Intn(max-min+1))
}

// RandomString returns a random string of the given length.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(letters)

	for i := 0; i < n; i++ {
		c := letters[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner returns a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney returns a random amount of money.
func RandomMoney() int64 {
	return int64(RandomInt(0, 1000))
}

// RandomCurrency returns a random currency from a list of currencies.
func RandomCurrency() string {
	currencies := []string{EUR, USD, IDR}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
