package utils

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
	wsUrl      string
	port       int
	sessionKey string
	manager    SessionManager
	session    Session
	e          *echo.Echo
}

type wsResponseAggregator struct {
	mu   sync.Mutex
	data map[*websocket.Conn]string
}

func (w *wsResponseAggregator) GetData() map[*websocket.Conn]string {
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
		agg.data[conn] = string(message)
	}
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
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return err
		}
		defer conn.Close()
		suite.manager.AddConnection(suite.sessionKey, conn)
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			formattedMessage := fmt.Sprintf("message: %s", message)
			err = conn.WriteMessage(mt, []byte(formattedMessage))
			if err != nil {
				return err
			}
		}
		return nil
	})

	go func() {
		suite.e.Logger.Fatal(suite.e.Start(fmt.Sprintf(":%d", suite.port)))
	}()

	time.Sleep(2 * time.Second)
}

func (suite *WSManagerTestSuite) TearDownTest() {
	suite.manager = SessionManager{
		Sessions: make(map[string]Session),
	}
}

/*-------------------Tests------------------------------*/

func (suite *WSManagerTestSuite) TestClientsGetAdded() {
	// baseUrl := fmt.Sprintf("http://localhost:%d", suite.port)
	originalSize := len(suite.manager.GetConnections(suite.sessionKey))
	log.Println("original:", originalSize)

	dialer := websocket.Dialer{}
	_, _, err := dialer.Dial(suite.wsUrl, nil)
	if err != nil {
		panic(err)
	}

	log.Println("after dial:", len(suite.manager.GetConnections(suite.sessionKey)))
	assert.Equal(suite.T(), len(suite.manager.GetConnections(suite.sessionKey)), originalSize+1)
}

func (suite *WSManagerTestSuite) TestOneToOneBroadcast() {
	// add some connections
	dialer1 := websocket.Dialer{}
	conn1, _, err := dialer1.Dial(suite.wsUrl, nil)
	if err != nil {
		panic(err)
	}

	responseData := wsResponseAggregator{
		data: map[*websocket.Conn]string{conn1: ""},
	}

	go listen(conn1, responseData)

	suite.manager.Broadcast(suite.sessionKey, conn1, "hello world")

	time.Sleep(2 * time.Second)

	assert.Equal(suite.T(), responseData.GetData()[conn1], "hello world")
}

/*-------------------Test Runner------------------------*/

func TestWSManagerTestSuite(t *testing.T) {
	suite.Run(t, new(WSManagerTestSuite))
}
