package main

import (
	"golang.org/x/net/websocket"
	//    "io"
	"fmt"
	"net/http"
)

func echoHandler(ws *websocket.Conn) {

	for {
		var (
			err error
			n   int
		)

		receivedtext := make([]byte, 100)

		n, err = ws.Read(receivedtext)

		if err != nil {
			fmt.Printf("Received: %d bytes\n", n)
		}

		s := string(receivedtext[:n])
		t := []byte("Back: "+s)
		fmt.Printf("Received: %d bytes: %s\n", n, s)
		n, err = ws.Write(t)
		if err != nil {
			fmt.Printf("Sent: %d bytes\n", n)
		}
		// 	io.Copy(ws, ws)
		fmt.Printf("Sent: %s\n", t)
	}
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
