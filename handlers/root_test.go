package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
	e *echo.Echo
}

/*-------------------Setups/Teardowns-------------------*/

func (suite *RootTestSuite) SetupTest() {
	suite.e = echo.New()
}

/*-------------------Tests------------------------------*/

func (suite *RootTestSuite) TestRootRedirectsToRoom() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	roomId := RootHandler(c)
	if assert.NoError(suite.T(), roomId) {
		assert.Equal(suite.T(), roomId, rec.HeaderMap.Get("Location"))
	}
}

/*-------------------Test Runner------------------------*/

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
