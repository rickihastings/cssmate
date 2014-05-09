package main

import (
	"log"
	"fmt"
	"regexp"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/howeyc/fsnotify"
)

var m *martini.ClassicMartini

func setupCssMate(p *int, hs *string, path *string) {
	folder = *path
	host = *hs
	port = *p
	bind = fmt.Sprintf("%s:%d", host, *p)
	// assign folder and port

	c := make(chan bool)
	
	go h.run()
	
	go func() {
		setupMartini()
		c <- true
	}()
	// here we setup martini

	setupWatcher()
	// here we can setup fsnotify without worrying about martini blocking us

	<-c
	// wait till martini is done
}

func setupMartini() {
	m = martini.Classic()
	m.Get("/websocket", wsHandler)
	m.Get("/cssmate.js", jsHandler)
	
	fmt.Println("[martini] listening on", bind)
	log.Fatal(http.ListenAndServe(bind, m))
}

func setupWatcher() {
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

	fmt.Println("[watcher] now watching for changes in", folder)

	<-done
	// wait till we're done

	watcher.Close()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// allow cross domain

	ws, err := upgrader.Upgrade(w, r, w.Header())
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

func jsHandler(w http.ResponseWriter, r *http.Request) (int, string) {
	data, err := Asset("embed/cssmate.js")
	if err != nil || len(data) == 0 {
		return 404, "404 page not found"
	}

	output := string(data)
	hostRegex := regexp.MustCompile("\\/\\* inject\\:host \\*\\/")
	portRegex := regexp.MustCompile("\\/\\* inject\\:port \\*\\/")
	
	output = hostRegex.ReplaceAllLiteralString(output, fmt.Sprintf("'%s'", host))
	output = portRegex.ReplaceAllLiteralString(output, fmt.Sprintf("%d", port))
	// inject the config vars

	return 200, output
}

func fileModified(ev *fsnotify.FileEvent) {
	fn := strings.Replace(ev.Name, folder, "", 1)[1:]
	
	for v := range h.connections {
		if stringInSlice(fn, v.files) {
			fmt.Println("[watcher] reloading", fn, "notifying clients")

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