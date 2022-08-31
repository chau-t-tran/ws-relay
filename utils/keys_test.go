package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
}

/*-------------------Tests------------------------------*/

func (suite *RootTestSuite) TestRandomKeysAreDifferent() {
	m := make(map[string]bool)
	for i := 1; i < 5000; i++ {
		key := RandomKey()
		fmt.Println(key)
		_, found := m[key]
		assert.True(suite.T(), !found, "Key collision!")

		if found {
			break
		}

		m[key] = true
	}
}

/*-------------------Test Runner------------------------*/

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
