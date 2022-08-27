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

func (suite *RootTestSuite) SetupTest() {
	suite.e = echo.New()
}

func (suite *RootTestSuite) TestRootHandler() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)
	if assert.NoError(suite.T(), RootHandler(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
	}
}

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
