package ws_manager

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chau-t-tran/ws-to-me/constants"
	"github.com/chau-t-tran/ws-to-me/templates"
	"github.com/chau-t-tran/ws-to-me/utils"
	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandlersTestSuite struct {
	suite.Suite
	seed           int
	sessionKeyJSON string
	manager        SessionManager
	e              *echo.Echo
}

/*-------------------Setups/Teardowns-------------------*/

func (suite *HandlersTestSuite) SetupTest() {
	suite.seed = constants.KEY_TEST_SEED
	rand.Seed(int64(suite.seed))
	key := utils.RandomKey()
	rand.Seed(int64(suite.seed))
	suite.sessionKeyJSON = fmt.Sprintf(`{"sessionKey":"%s"}`, key)

	suite.manager = CreateSessionManager([]string{})

	suite.e = echo.New()
	suite.e.Renderer = templates.Renderer
}

/*-------------------Tests------------------------------*/

func (suite *HandlersTestSuite) TestRootHandler() {
	rand.Seed(int64(suite.seed))
	key := "/" + utils.RandomKey()

	req := httptest.NewRequest(http.MethodGet, "/"+key, nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	assert.NoError(suite.T(), suite.manager.RootHandler(c))
}

func (suite *HandlersTestSuite) TestRegisterHandler() {
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	rec := httptest.NewRecorder()
	c := suite.e.NewContext(req, rec)

	if assert.NoError(suite.T(), suite.manager.RegisterHandler(c)) {
		assert.Equal(suite.T(), suite.sessionKeyJSON+"\n", rec.Body.String())
	}
}

/*-------------------Test Runner------------------------*/

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
