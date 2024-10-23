package utils

import "math/rand"

func GenRandCode() string {
	letters := "1234567890"
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
