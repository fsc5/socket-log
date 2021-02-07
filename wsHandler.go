package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type wsHandler struct {
	upgrader    websocket.Upgrader
	stampDirect chan KeyStamp
}

func newWsHandler(stampDirect chan KeyStamp) *wsHandler {
	return &wsHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		stampDirect: stampDirect,
	}
}

func (ws *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	con, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	keyStrokes := make([]string, 0)
	con.WriteMessage(websocket.TextMessage, []byte("Hello :)"))
	stopChan := make(chan bool)
	mutex := &sync.Mutex{}

	go ws.clearList(&keyStrokes, mutex, stopChan)

	for {
		_, msg, err := con.ReadMessage()
		mutex.Lock()
		if err != nil {
			stopChan <- true
			fmt.Println(err)
			break
		}
		keyStrokes = append(keyStrokes, string(msg))
		stopChan <- false
		mutex.Unlock()
	}
}

func (ws *wsHandler) clearList(keystrokes *[]string, mutex *sync.Mutex, interruptChan chan bool) {
	for {
		select {
		case stopped := <-interruptChan:
			if stopped {
				ws.transferStrokes(keystrokes)
				return
			}
		case <-time.After(time.Second * 5):
			mutex.Lock()
			ws.transferStrokes(keystrokes)
			mutex.Unlock()
		}
	}
}

func (ws *wsHandler) transferStrokes(keystrokes *[]string) {
	if len(*keystrokes) > 0 {
		ws.stampDirect <- KeyStamp{
			Strokes: *keystrokes,
			Time:    time.Now(),
		}
		*keystrokes = make([]string, 0)
	}
}
