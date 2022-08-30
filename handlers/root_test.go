package handlers

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chau-t-tran/ws-relay/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
	seed      int
	keyLength int
	e         *echo.Echo
}

/*-------------------Setups/Teardowns-------------------*/

func (suite *RootTestSuite) SetupTest() {
	suite.e = echo.New()
	suite.seed = 42
	suite.keyLength = 7
}

/*-------------------Tests------------------------------*/

func (suite *RootTestSuite) TestRootRedirectsToRoom() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	rand.Seed(int64(suite.seed))
	key := utils.RandomKey(suite.keyLength)
	rand.Seed(int64(suite.seed))

	if assert.NoError(suite.T(), RootHandler(c)) {
		assert.Equal(suite.T(), key, rec.HeaderMap.Get("Location"))
	}
}

/*-------------------Test Runner------------------------*/

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
