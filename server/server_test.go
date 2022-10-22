package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/chau-t-tran/ws-relay/constants"
	"github.com/chau-t-tran/ws-relay/utils"
	"github.com/chau-t-tran/ws-relay/ws_manager"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
	seed       int
	e          *echo.Echo
	PORT       int
	httpClient http.Client
}

/*-------------------Setups/Teardowns-------------------*/

func (suite *RootTestSuite) SetupSuite() {
	suite.PORT = 9000
	suite.seed = constants.KEY_TEST_SEED
	suite.e = GetServer()
	suite.httpClient = http.Client{
		Timeout: 5 * time.Second,
	}
	go func() {
		suite.e.Logger.Fatal(suite.e.Start(fmt.Sprintf(":%d", suite.PORT)))
	}()
	time.Sleep(5 * time.Second)
}

func (suite *RootTestSuite) TeardownSuite() {
	if suite.e != nil {
		suite.e.Close()
	}
}

/*-------------------Tests------------------------------*/

func (suite *RootTestSuite) TestServerIsOnline() {
	url := fmt.Sprintf("http://localhost:%d/health", suite.PORT)
	resp, err := suite.httpClient.Get(url)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *RootTestSuite) TestLandingPage() {
	url := fmt.Sprintf("http://localhost:%d/", suite.PORT)
	resp, err := suite.httpClient.Get(url)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	contentTypeArray, ok := resp.Header["Content-Type"]
	assert.True(suite.T(), ok)

	contentType := contentTypeArray[0]
	assert.Equal(suite.T(), "text/html; charset=UTF-8", contentType)
}

func (suite *RootTestSuite) TestRegistrationCall() {
	rand.Seed(int64(suite.seed))
	sessionKey := utils.RandomKey()
	rand.Seed(int64(suite.seed))

	respTarget := ws_manager.SessionRegResponse{}
	url := fmt.Sprintf("http://localhost:%d/register", suite.PORT)
	contentType := "application/json"
	jsonBody := []byte(`{}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := suite.httpClient.Post(url, contentType, bodyReader)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&respTarget)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), respTarget.SessionKey, sessionKey)
}

/*-------------------Runner-----------------------------*/

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
