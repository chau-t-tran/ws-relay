package handlers

import (
	"fmt"
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

type APITestSuite struct {
	suite.Suite
	seed           int
	e              *echo.Echo
	sessionKeyJSON string
}

/*-------------------Setups/Teardowns-------------------*/

func (suite *APITestSuite) SetupTest() {
	suite.seed = constants.KEY_TEST_SEED
	suite.e = echo.New()
	suite.e.Renderer = templates.Renderer

	rand.Seed(int64(suite.seed))
	key := utils.RandomKey()
	rand.Seed(int64(suite.seed))
	suite.sessionKeyJSON = fmt.Sprintf(`{"sessionKey":"%s"}`, key)
}

/*-------------------Tests------------------------------*/

func (suite *APITestSuite) TestAPIHandler() {
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	if assert.NoError(suite.T(), APIHandler(c)) {
		assert.Equal(suite.T(), suite.sessionKeyJSON+"\n", rec.Body.String())
	}
}

/*-------------------Test Runner------------------------*/

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
