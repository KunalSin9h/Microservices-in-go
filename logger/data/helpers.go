package data

import (
	"crypto/rand"
	"math/big"
)

const charSet = "abcdefghijklmnopqrstuvwzyxABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func GetRandomSequence(l int) string {
	var res string
	for i := 0; i < l; i++ {
		x := big.NewInt(int64(len(charSet)))
		n, _ := rand.Int(rand.Reader, x)
		c := charSet[n.Int64()]
		res += string(c)
	}
	return res
}
