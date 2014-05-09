package main

import (
	"flag"
	"github.com/rickihastings/cssmate/websocket"
)

type message struct {
	Command		string	`json:"command"`
	Data		string	`json:"data"`
}

var (
	bind		string
	folder		string
	host		string
	port		int
	upgrader 	= websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

func main() {
	portPtr := flag.Int("port", 58900, "the default port to run on")
	hostPtr := flag.String("host", "0.0.0.0", "the host to bind to")
	pathPtr := flag.String("path", "", "the folder to monitor, can be local or full path")

	flag.Parse()
	
	setupCssMate(portPtr, hostPtr, pathPtr)
}

