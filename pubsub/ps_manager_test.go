package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WSManagerTestSuite struct {
	suite.Suite
	wsUrl      string
	port       int
	sessionKey string
	manager    SessionManager
	session    Session
	e          *echo.Echo
}

/*-------------------Setups/Teardowns-------------------*/

var (
	upgrader = websocket.Upgrader{}
)

func (suite *WSManagerTestSuite) SetupSuite() {
	suite.port = 4000
	suite.wsUrl = fmt.Sprintf("ws://localhost:%d", suite.port)
	suite.sessionKey = "abcdefgh"
	suite.manager = SessionManager{
		Sessions: make(map[string]Session),
	}

	suite.e = echo.New()
	suite.e.GET("/", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		suite.manager.AddConnection(suite.sessionKey, ws)
		return err
	})

	go func() {
		suite.e.Logger.Fatal(suite.e.Start(fmt.Sprintf(":%d", suite.port)))
	}()

	time.Sleep(2 * time.Second)
}

/*-------------------Tests------------------------------*/

func (suite *WSManagerTestSuite) TestClientsGetAdded() {
	// baseUrl := fmt.Sprintf("http://localhost:%d", suite.port)
	assert.Equal(suite.T(), len(suite.manager.GetConnections(suite.sessionKey)), 0)

	dialer := websocket.Dialer{}
	_, _, err := dialer.Dial(suite.wsUrl, nil)
	if err != nil {
		panic(err)
	}

	assert.Equal(suite.T(), len(suite.manager.GetConnections(suite.sessionKey)), 1)
}

/*-------------------Test Runner------------------------*/

func TestWSManagerTestSuite(t *testing.T) {
	suite.Run(t, new(WSManagerTestSuite))
}
