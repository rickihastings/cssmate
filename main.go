package main

import (
	"log"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/go-martini/martini"
	"github.com/howeyc/fsnotify"
)

type message struct {
	Command		string	`json:"command"`
	Data		string	`json:"data"`
}

var (
	folder		string = "public"
	upgrader 	= &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

func main() {
	go h.run()
	SetupServices()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	
	c := &connection{send: make(chan []byte, 256), websocket: ws}
	h.register <- c
	
	defer func() {
		h.unregister <- c
	}()
	
	go c.writer()
	c.reader()
}

func fileModified(ev *fsnotify.FileEvent) {
	fn := strings.Replace(ev.Name, folder, "", 1)[1:]
	
	for v := range h.connections {
		if stringInSlice(fn, v.files) {
			fmt.Println("Reloading", fn, "notifying clients")

			var out message = message{"changed", fn}
			output, jerr := json.Marshal(out)

			if jerr != nil {
				fmt.Println(jerr)
				return
			}

			v.send <- output
			// create a message to send to the client to tell them to reload x file
		}
	}
}

func SetupServices() {
	c := make(chan bool)
	
	go func() {
		SetupMartini()
		c <- true
	}()
	// here we setup martini

	SetupWatcher()
	// here we can setup fsnotify without worrying about martini blocking us

	<-c
	// wait till martini is done
}

func SetupMartini() {
	m := martini.Classic()
	m.Get("/websocket", wsHandler)
	m.Run()
}

func SetupWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// if there is any errors back out

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() {
					fileModified(ev)
				}
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()
	// setup the event loop to keep us going

	err = watcher.Watch(folder)
	if err != nil {
		log.Fatal(err)
	}
	// setup the watcher and wait for errors

	<-done
	// wait till we're done

	watcher.Close()
}