package main

import (
	"golang.org/x/net/websocket"
	"fmt"
	"net/http"
	"os"
)

func echoHandler(ws *websocket.Conn) {

	for {
		var (
			err error
			n   int
			outFile *os.File
			t []byte
			msg ParamStruct
		)

		receivedtext := make([]byte, 100)

		n, err = ws.Read(receivedtext)
		if err != nil {
			fmt.Printf("Received: %d bytes\n", n)
		}

		s := string(receivedtext[:n])
		
		fmt.Println("Received ", s)
		
		ws.Write([]byte(s))
		
		websocket.JSON.Receive(ws, &msg)
		fmt.Printf("Received message: %+v\n", msg)
		
		if outFile, err = os.Create("params.conf"); err != nil {
			t = []byte("Can't create file with err: " + err.Error())
		}
		
		outString := "param1 = " + msg.Param1 + "\n" +
					"param2 = " + msg.Param2 + "\n" + 
					"param2 = " + msg.Param3 + "\n"

		if _, err = outFile.WriteString(outString); err != nil {
			t = []byte("Can't write to file with err: " + err.Error())
		}
		outFile.Close()
		

		t = []byte(`<a href="params.conf" download="params.conf" >Click to download configuration file</a>`)
		
		n, err = ws.Write(t)
		if err != nil {
			fmt.Printf("Sent: %d bytes\n", n)
		}
		// 	io.Copy(ws, ws)
		fmt.Printf("Sent: %s\n", t)
	}
}

type ParamStruct struct {
	Param1 string
	Param2 string
	Param3 string
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
