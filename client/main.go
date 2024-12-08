package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

type Message struct {
	MessageType int
	Data        []byte
}

func main() {
	// create connection with default dialer
	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/ws"}
	fmt.Printf("Connecting to %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial: ", err)
	}
	defer conn.Close()

	// Channels for managing messages
	send := make(chan Message)
	done := make(chan struct{})

	// Goroutine to listen to incoming messages in background
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read: ", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	// Goroutine for sending messages
	go func() {
		for {
			select {
			case msg := <-send:
				// write to websocket connection
				err := conn.WriteMessage(msg.MessageType, msg.Data)
				if err != nil {
					log.Println("write: ", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// read input from terminal and send to websocket server
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("type something...")
	for scanner.Scan() {
		text := scanner.Text()
		// send text to channel
		send <- Message{websocket.TextMessage, []byte(text)}
	}

	if err := scanner.Err(); err != nil {
		log.Println("scanner err:", err)
	}
}
