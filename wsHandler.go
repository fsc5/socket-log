package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type wsHandler struct {
	upgrader   websocket.Upgrader
	sendStamps chan<- KeyStamp
	keyStrokes []string
}

func newWsHandler(sendStamps chan<- KeyStamp) *wsHandler {
	return &wsHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		sendStamps: sendStamps,
	}
}

func (ws *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	con, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	stopChan := make(chan bool)
	mutex := &sync.Mutex{}

	go ws.clearList(mutex, stopChan)
	for {
		_, msg, err := con.ReadMessage()
		mutex.Lock()
		if err != nil {
			stopChan <- true
			fmt.Println(err)
			break
		}
		ws.keyStrokes = append(ws.keyStrokes, string(msg))
		stopChan <- false
		mutex.Unlock()
	}
}

func (ws *wsHandler) clearList(mutex *sync.Mutex, interruptChan chan bool) {
	for {
		select {
		case stopped := <-interruptChan:
			if stopped {
				ws.transferStrokes()
				return
			}
		case <-time.After(time.Second * 5):
			mutex.Lock()
			ws.transferStrokes()
			mutex.Unlock()
		}
	}
}

func (ws *wsHandler) transferStrokes() {
	if len(ws.keyStrokes) > 0 {
		ws.sendStamps <- KeyStamp{
			Strokes: ws.keyStrokes,
			Time:    time.Now(),
		}
		ws.keyStrokes = make([]string, 0)
	}
}
