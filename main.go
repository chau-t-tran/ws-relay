package main

import "github.com/chau-t-tran/ws-relay/server"

func main() {
	e := server.GetServer()
	e.Logger.Fatal(e.Start(":5000"))
}
