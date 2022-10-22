package ws_manager

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WSManagerTestSuite struct {
	suite.Suite
	wsUrl       string
	port        int
	sessionKey  string
	manager     SessionManager
	session     []*websocket.Conn
	e           *echo.Echo
	testMessage string
}

type wsResponseAggregator struct {
	mu   sync.Mutex
	data map[string]string
}

func (w *wsResponseAggregator) GetData() map[string]string {
	return w.data
}

func listen(conn *websocket.Conn, agg wsResponseAggregator) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error: %s", err)
			return
		}
		agg.mu.Lock()
		defer agg.mu.Unlock()
		agg.data[conn.LocalAddr().String()] = string(message)
	}
}

/*-------------------Setups/Teardowns-------------------*/

func (suite *WSManagerTestSuite) SetupSuite() {
	suite.port = 4000
	suite.sessionKey = "abcdefgh"
	suite.wsUrl = fmt.Sprintf("ws://localhost:%d/%s", suite.port, suite.sessionKey)
	suite.testMessage = "hello world"
	suite.manager = CreateSessionManager([]string{})

	suite.e = echo.New()
	suite.e.GET("/:sessionKey", suite.manager.EchoHandler)

	go func() {
		suite.e.Logger.Fatal(suite.e.Start(fmt.Sprintf(":%d", suite.port)))
	}()

	time.Sleep(2 * time.Second)
}

func (suite *WSManagerTestSuite) TearDownTest() {
	suite.manager = CreateSessionManager([]string{})
}

/*-------------------Tests------------------------------*/

func (suite *WSManagerTestSuite) TestSessionManagerConstructor() {
	key1 := "abcdefgh"
	key2 := "ijklmnop"
	keys := []string{key1, key2}
	sm := CreateSessionManager(keys)
	_, err1 := sm.GetSession(key1)
	_, err2 := sm.GetSession(key2)
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
}

func (suite *WSManagerTestSuite) TestClientsGetAdded() {
	// baseUrl := fmt.Sprintf("http://localhost:%d", suite.port)
	session, err := suite.manager.GetSession(suite.sessionKey)
	originalSize := len(session)

	dialer := websocket.Dialer{}
	_, _, err = dialer.Dial(suite.wsUrl, nil)
	if err != nil {
		panic(err)
	}

	updatedSession, err := suite.manager.GetSession(suite.sessionKey)
	updatedSize := len(updatedSession)
	assert.Equal(suite.T(), updatedSize, originalSize+1)
}

func (suite *WSManagerTestSuite) TestDoesNotBroadcastToSelf() {
	// add some connections
	dialer1 := websocket.Dialer{}
	conn1, _, err := dialer1.Dial(suite.wsUrl, nil)
	if err != nil {
		panic(err)
	}

	responseData := wsResponseAggregator{
		data: map[string]string{},
	}

	go listen(conn1, responseData)

	conn1.WriteMessage(1, []byte(suite.testMessage))
	// suite.manager.Broadcast(suite.sessionKey, conn1.LocalAddr().String(), suite.testMessage)
	time.Sleep(2 * time.Second)

	assert.Equal(suite.T(), "", responseData.GetData()[conn1.LocalAddr().String()])
}

func (suite *WSManagerTestSuite) TestOneToManyBroadcast() {
	// add multiple connections
	dialer1 := websocket.Dialer{}
	conn1, _, err := dialer1.Dial(suite.wsUrl, nil)
	if err != nil {
		panic(err)
	}

	dialer2 := websocket.Dialer{}
	conn2, _, err := dialer2.Dial(suite.wsUrl, nil)
	if err != nil {
		panic(err)
	}

	dialer3 := websocket.Dialer{}
	conn3, _, err := dialer3.Dial(suite.wsUrl, nil)
	if err != nil {
		panic(err)
	}

	// test aggregator
	responseData := wsResponseAggregator{
		data: map[string]string{},
	}

	// listen on all three connections
	go listen(conn1, responseData)
	go listen(conn2, responseData)
	go listen(conn3, responseData)

	// test broadcast
	conn1.WriteMessage(1, []byte(suite.testMessage))
	//suite.manager.Broadcast(suite.sessionKey, conn1.LocalAddr().String(), suite.testMessage)
	time.Sleep(2 * time.Second)

	log.Println("RESPONSES:", responseData.GetData())

	assert.Equal(suite.T(), suite.testMessage, responseData.GetData()[conn2.LocalAddr().String()])
	assert.Equal(suite.T(), suite.testMessage, responseData.GetData()[conn3.LocalAddr().String()])
}

/*-------------------Test Runner------------------------*/

func TestWSManagerTestSuite(t *testing.T) {
	suite.Run(t, new(WSManagerTestSuite))
}
