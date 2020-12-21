package main

import (
	"golangdynamocrud/server"
	"github.com/subosito/gotenv"
)

func main(){
	gotenv.Load()
	server.LoadRoutes()
}