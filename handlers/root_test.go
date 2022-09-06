package handlers

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chau-t-tran/ws-relay/constants"
	"github.com/chau-t-tran/ws-relay/templates"
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
	suite.seed = constants.KEY_TEST_SEED
	suite.e.Renderer = templates.Renderer
}

/*-------------------Tests------------------------------*/

func (suite *RootTestSuite) TestRootRedirectsToRoom() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	rand.Seed(int64(suite.seed))
	key := "/" + utils.RandomKey()
	rand.Seed(int64(suite.seed))

	if assert.NoError(suite.T(), RootHandler(c)) {
		assert.Equal(suite.T(), key, rec.HeaderMap.Get("Location"))
	}
}

func (suite *RootTestSuite) TestSessionHandler() {
	rand.Seed(int64(suite.seed))
	key := "/" + utils.RandomKey()

	req := httptest.NewRequest(http.MethodGet, "/"+key, nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	assert.NoError(suite.T(), SessionHandler(c))
}

/*-------------------Test Runner------------------------*/

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
