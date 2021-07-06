package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// message show realtime fake count user in SiteName
type message struct {
	Count    int    `json:"Count"`
	SiteName string `json:"SiteName"`
}

// hub keep open connection
type hub struct {
	Conn *websocket.Conn
}

// upgrader set buffer R/W size websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// newHub constructor connection instance
func newHub(conn *websocket.Conn) *hub {
	return &hub{
		Conn: conn,
	}
}

func (h *hub) run() {
	rand.Seed(time.Now().UnixNano())

	for {
		_, readMessage, err := h.receive()
		if err != nil {
			log.Println("Could not received message from connection, error : ", err)
			continue
		}

		siteName := string(readMessage)

		for {
			time.Sleep(time.Second)
			message := message{
				Count:    rand.Intn(2001),
				SiteName: siteName,
			}
			err = h.send(message)
			if err != nil {
				log.Println("Could not send message to connection, error : ", err)
				continue
			}
		}
	}
}

// receive message from client
func (h *hub) receive() (messageType int, message []byte, err error) {
	messageType, message, err = h.Conn.ReadMessage()
	if err != nil {
		return 0, []byte{}, err
	}
	return
}

// send message from server
func (h *hub) send(message message) (err error) {
	err = h.Conn.WriteJSON(message)
	return
}
