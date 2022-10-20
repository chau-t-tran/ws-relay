package server

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/chau-t-tran/ws-relay/routes"
	"github.com/chau-t-tran/ws-relay/templates"
	"github.com/labstack/echo/v4"
)

func GetServer() *echo.Echo {
	e := echo.New()

	/* ------Pseudorand Seed----- */
	rand.Seed(int64(time.Now().UTC().UnixNano()))

	/* ----------Registry-------- */
	e.Renderer = templates.Renderer
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

	/* ----------Routes---------- */
	routes.SetupRootRoutes(e)
	routes.SetupWSRoutes(e)

	return e
}
