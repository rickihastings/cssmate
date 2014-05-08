package main

import (
	"flag"
	"github.com/gorilla/websocket"
)

type message struct {
	Command		string	`json:"command"`
	Data		string	`json:"data"`
}

var (
	port		string
	folder		string
	upgrader 	= &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

func main() {
	portPtr := flag.Int("port", 58900, "the default port to run on")
	pathPtr := flag.String("path", "public", "the folder to monitor, can be local or full path")

	flag.Parse()
	
	setupCssMate(portPtr, pathPtr)
}

