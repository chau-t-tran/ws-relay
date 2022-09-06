package main

import (
	"math/rand"
	"time"

	"github.com/chau-t-tran/ws-relay/routes"
	"github.com/chau-t-tran/ws-relay/templates"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	/* ------Pseudorand Seed----- */
	rand.Seed(int64(time.Now().UTC().UnixNano()))

	/* ----------Registry-------- */
	e.Renderer = templates.Renderer

	/* ----------Routes---------- */
	routes.SetupRootRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
