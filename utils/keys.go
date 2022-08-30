package utils

import "math/rand"

func RandomKey(n int) string {
	s := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var k []rune

	for i := 0; i < n; i++ {
		k = append(k, rune(s[rand.Intn(len(s))]))
	}

	return string(k)
}
