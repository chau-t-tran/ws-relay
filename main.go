package main

import (
	"github.com/chau-t-tran/ws-relay/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	routes.SetupRootRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
