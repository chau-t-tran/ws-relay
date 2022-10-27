package main

import "github.com/chau-t-tran/ws-to-me/server"

func main() {
	e := server.GetServer()
	e.Logger.Fatal(e.Start(":5000"))
}
