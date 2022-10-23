package ws_manager

import (
	"fmt"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GCTestSuite struct {
	suite.Suite
	manager    SessionManager
	sessionKey string
	timeFormat string

	wsUrl string
	port  int
	e     *echo.Echo
}

/*-------------------Setups/Teardowns-------------------*/

func (suite *GCTestSuite) SetupSuite() {
	suite.timeFormat = "02 Jan 06 15:04 MST"
	suite.sessionKey = "abcdefgh"

	suite.port = 9000
	suite.e = echo.New()
	suite.e.GET("/:sessionKey", suite.manager.EchoHandler)
	go func() {
		suite.e.Logger.Fatal(
			suite.e.Start(fmt.Sprintf(":%d", suite.port)),
		)
	}()
	time.Sleep(2 * time.Second)

	suite.wsUrl = fmt.Sprintf("ws://localhost:%d/%s", suite.port, suite.sessionKey)
}

func (suite *GCTestSuite) SetupTest() {
	suite.manager = CreateSessionManager([]string{})
}

/*-------------------Tests------------------------------*/

func (suite *GCTestSuite) TestLastUsedInitiates() {
	suite.manager.RegisterSession(suite.sessionKey)

	currentTimeString := time.Now().Format(suite.timeFormat)

	lastUsedTime, err := suite.manager.GetLastUsedTime(suite.sessionKey)
	lastUsedTimeString := lastUsedTime.Format(suite.timeFormat)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), currentTimeString, lastUsedTimeString)
}

func (suite *GCTestSuite) TestLastUsedUpdatesOnBroadcast() {
	suite.manager.RegisterSession(suite.sessionKey)
	time.Sleep(2)

	suite.manager.Broadcast(suite.sessionKey, "", []byte{})
	currentTimeString := time.Now().Format(suite.timeFormat)

	lastUsedTime, err := suite.manager.GetLastUsedTime(suite.sessionKey)
	lastUsedTimeString := lastUsedTime.Format(suite.timeFormat)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), currentTimeString, lastUsedTimeString)
}

/*-------------------Test Runner------------------------*/

func TestGCTestSuite(t *testing.T) {
	suite.Run(t, new(GCTestSuite))
}
