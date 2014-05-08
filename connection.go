package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/rickihastings/cssmate/websocket"
)

type connection struct {
	websocket	*websocket.Conn
	files		[]string
	send 		chan []byte
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func (c *connection) reader() {
	for {
		_, in, err := c.websocket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		var m message
		e := json.Unmarshal(in, &m)

		if e != nil {
			fmt.Println(e)
			break
		}

		c.handle(&m)
	}

	c.websocket.Close()
}

func (c *connection) writer() {
	for in := range c.send {
		err := c.websocket.WriteMessage(websocket.TextMessage, in)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	c.websocket.Close()
}

func (c *connection) handle(m *message) {
	if m.Command == "watching" {
		c.files = strings.Split(m.Data, ",")

		var out message = message{"watching", "true"}
		output, jerr := json.Marshal(out)

		if jerr != nil {
			fmt.Println(jerr)
			return
		}

		c.send <- output
		// create a response and send it back
	}
}