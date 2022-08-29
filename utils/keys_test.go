package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
}

/*-------------------Tests------------------------------*/

func (suite *RootTestSuite) TestRandomKeysAreDifferent() {
	n := 7
	m := make(map[string]bool)
	for i := 1; i < 5000; i++ {
		key := randomKey(n)
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
