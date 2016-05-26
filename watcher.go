// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"path/filepath"
	"time"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 1 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 2 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 300 * time.Millisecond
)

var (
	filename  string
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func readFileIfModified(lastMod time.Time, env *Environment) ([]byte, time.Time, error) {
	filename  = filepath.Join(config.RootFolder, "/repo-" + env.Name, env.Path, "/planOutput")
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return nil, lastMod, nil
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fi.ModTime(), err
	}
	return p, fi.ModTime(), nil
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, lastMod time.Time, env *Environment) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	changesChannel := getChangesChannel()
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case envID := <-changesChannel:
				log.Printf("Channel got something: %v", envID)
				if envID == env.Id {
					ws.SetWriteDeadline(time.Now().Add(writeWait))
					data := []byte(strconv.Itoa(envID))
					if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
						return
					}					
				}
		case <-fileTicker.C:
			var p []byte
			var err error

			p, lastMod, err = readFileIfModified(lastMod, env)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte("error")
				}
			} else {
				lastError = ""
			}

			if p != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n)
	}

	var envID int
	if envID, err = strconv.Atoi(r.FormValue("envID")); err != nil {
		log.Println(err)
	}

	env := RepoFindEnvironment(envID)
	go writer(ws, lastMod, env)
	reader(ws)
}
