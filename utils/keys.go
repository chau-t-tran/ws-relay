package utils

import "crypto/rand"

func randomKey(n int) string {
	key := make([]byte, n)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}

	return string(key)
}
