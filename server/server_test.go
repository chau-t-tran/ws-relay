package server

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
	e                     *echo.Echo
	PORT                  int
	RequestTimeoutSeconds int
}

/*-------------------Setups/Teardowns-------------------*/

func (suite *RootTestSuite) SetupSuite() {
	suite.PORT = 9000
	suite.RequestTimeoutSeconds = 5
}

func (suite *RootTestSuite) TeardownSuite() {
	if suite.e != nil {
		suite.e.Close()
	}
}

/*-------------------Tests------------------------------*/

func (suite *RootTestSuite) TestServerIsOnline() {
	suite.e = GetServer()
	go func() {
		suite.e.Logger.Fatal(suite.e.Start(fmt.Sprintf(":%d", suite.PORT)))
	}()

	got200 := false

	for i := 0; i < suite.RequestTimeoutSeconds; i++ {
		resp, _ := http.Get(fmt.Sprintf("http://localhost:%d/health", suite.PORT))
		if resp != nil {
			got200 = resp.StatusCode == 200
			break
		}
		time.Sleep(1 * time.Second)
	}

	assert.True(suite.T(), got200)
}

/*-------------------Runner-----------------------------*/

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, new(RootTestSuite))
}
