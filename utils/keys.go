package utils

import (
	"math/rand"

	"github.com/chau-t-tran/ws-relay/constants"
)

func RandomKey() string {
	var key []rune

	for i := 0; i < constants.KEY_SIZE; i++ {
		key_index := rand.Intn(constants.KEY_SPACE_SIZE)
		key = append(key, rune(constants.KEY_SPACE[key_index]))
	}

	return string(key)
}
