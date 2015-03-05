package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
)

func (manager managerStruct) MessagesHandler(mws *websocket.Conn) {
	// 	manager.Mws = mws

	log.Println("Start messagesHandler loop")
	var (
		text        string
		receiveChan = make(chan string, 1)
	)

	go func(receiveChan chan string, mws *websocket.Conn) {
		var (
			err          error
			n            int
			receivedtext = make([]byte, 100)
		)

		for {
			n, err = mws.Read(receivedtext)
			if err != nil {
				log.Println("Received error: ", err)
				continue
			}
			receiveChan <- "Received " + string(receivedtext[:n])
		}
	}(receiveChan, mws)

	for {
		select {
		case text = <-manager.Phone:
			mws.Write([]byte(text))
		case text = <-receiveChan:
			log.Println(text)
		}

	}
}

func (manager managerStruct) DataHandler(dws *websocket.Conn) {
	// 	manager.Dws = dws
	log.Println("Start dataHandler loop")
	var (
		err error
		// 			n   int
		outFile *os.File
		text    string
		data    ParamStruct
	)

	for {

		websocket.JSON.Receive(dws, &data)
		log.Printf("Received message: %+v\n", data)

		if outFile, err = os.Create("params.conf"); err != nil {
			text = "Can't create file with err: " + err.Error()
		}

		outString := "param1 = " + data.Param1 + "\n" +
			"param2 = " + data.Param2 + "\n" +
			"param2 = " + data.Param3 + "\n"

		if _, err = outFile.WriteString(outString); err != nil {
			text = "Can't write to file with err: " + err.Error()
		}
		outFile.Close()

		text = `<a href="params.conf" download="params.conf" >Click to download configuration file</a>`

		manager.Phone <- text

		log.Printf("Sent: %s\n", text)
	}
}

type ParamStruct struct {
	Param1 string
	Param2 string
	Param3 string
}

func main() {
	var manager = managerStruct{
		Phone: make(chan string),
		Done:  make(chan struct{}),
	}
	http.Handle("/messages", websocket.Handler(manager.MessagesHandler))
	http.Handle("/data", websocket.Handler(manager.DataHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}

type managerStruct struct {
	// 	Mws *websocket.Conn
	// 	Dws *websocket.Conn
	Phone chan string
	Done  chan struct{}
}
