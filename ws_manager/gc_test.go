package ws_manager

import (
	"context"
	"fmt"
	"log"
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
	suite.timeFormat = "02 Jan 06 15:04:05 MST"
	suite.sessionKey = "abcdefgh"

	suite.wsUrl = fmt.Sprintf("ws://localhost:%d/%s", suite.port, suite.sessionKey)
	suite.port = 4000
	suite.e = echo.New()
	suite.e.GET("/:sessionKey", suite.manager.EchoHandler)
	go func() {
		suite.e.Start(fmt.Sprintf(":%d", suite.port))
	}()
	time.Sleep(2 * time.Second)
}

func (suite *GCTestSuite) TearDownSuite() {
	suite.manager.cronScheduler.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := suite.e.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func (suite *GCTestSuite) SetupTest() {
	log.Println("NEW SESSION")
	suite.manager = CreateSessionManager([]string{})
}

func (suite *GCTestSuite) TearDownTest() {
	log.Println("DESTROYING OLD SESSION")
	suite.manager.cronScheduler.Stop()
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
	registrationTime, err := suite.manager.GetLastUsedTime(suite.sessionKey)
	registrationTimeString := registrationTime.Format(suite.timeFormat)

	time.Sleep(3 * time.Second)

	suite.manager.Broadcast(suite.sessionKey, "", []byte{})
	currentTimeString := time.Now().Format(suite.timeFormat)

	lastUsedTime, err := suite.manager.GetLastUsedTime(suite.sessionKey)
	lastUsedTimeString := lastUsedTime.Format(suite.timeFormat)
	assert.NoError(suite.T(), err)

	assert.NotEqual(suite.T(), lastUsedTimeString, registrationTimeString)
	assert.Equal(suite.T(), currentTimeString, lastUsedTimeString)
}

func (suite *GCTestSuite) TestSessionGetsDeletedAfterOneDay() {
	suite.manager.cronScheduler.Stop()
	suite.manager.cronScheduler.Clear()
	_, _ = suite.manager.cronScheduler.
		// Every(1).
		// Day().
		Every(2).
		Seconds().
		Do(suite.manager.GarbageCollectDaily)
	suite.manager.cronScheduler.StartAsync()

	suite.manager.RegisterSession(suite.sessionKey)

	_, err := suite.manager.GetSession(suite.sessionKey)
	assert.NoError(suite.T(), err)

	_, err = suite.manager.GetLastUsedTime(suite.sessionKey)
	assert.NoError(suite.T(), err)

	time.Sleep(5 * time.Second) // simulate some amount of time

	_, err = suite.manager.GetSession(suite.sessionKey)
	expectedSessionErrorMessage := fmt.Sprintf("Session %s not found", suite.sessionKey)
	assert.Error(suite.T(), err, expectedSessionErrorMessage)

	_, err = suite.manager.GetLastUsedTime(suite.sessionKey)
	expectedTimeErrorMessage := fmt.Sprintf("Session %s last used time not found", suite.sessionKey)
	assert.Error(suite.T(), err, expectedTimeErrorMessage)

	// clean-up
	// suite.manager.cronScheduler.Stop()
}

/*-------------------Test Runner------------------------*/

func TestGCTestSuite(t *testing.T) {
	log.Println("GC TEST")
	suite.Run(t, new(GCTestSuite))
}
