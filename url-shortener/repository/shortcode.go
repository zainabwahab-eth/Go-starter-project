package repository

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateShortCode() string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, 6)
	max := big.NewInt(int64(len(charset)))

	for i := 0; i < 6; i++ {
		nBig, err := rand.Int(rand.Reader, max)

		if err != nil {
			fmt.Println("Error generating random number:", err)
			return string(result)
		}

		result[i] = charset[nBig.Int64()]
	}

	return string(result)
}
