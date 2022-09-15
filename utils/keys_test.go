package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type KeyTestSuite struct {
	suite.Suite
}

/*-------------------Key Tests--------------------------*/

func (suite *KeyTestSuite) TestRandomKeysAreDifferent() {
	m := make(map[string]bool)
	for i := 1; i < 5000; i++ {
		key := RandomKey()
		_, found := m[key]
		assert.True(suite.T(), !found, "Key collision!")

		if found {
			break
		}

		m[key] = true
	}
}

/*-------------------Test Runner------------------------*/

func TestKeyTestSuite(t *testing.T) {
	suite.Run(t, new(KeyTestSuite))
}
