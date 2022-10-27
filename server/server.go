package server

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/chau-t-tran/ws-to-me/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GetServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	/* ------Pseudorand Seed----- */
	rand.Seed(int64(time.Now().UTC().UnixNano()))

	/* ----------Registry-------- */
	e.Renderer = templates.Renderer
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

	/* ----------Routes---------- */
	SetupWSRoutes(e)

	return e
}
