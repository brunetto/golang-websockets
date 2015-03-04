package main
 
import (
    "golang.org/x/net/websocket"
//    "code.google.com/p/go.net/websocket"
//    "io"
    "net/http"
    "fmt"
)
 
func echoHandler(ws *websocket.Conn) {
 
  for {
    receivedtext := make([]byte, 100)
 
    n,err := ws.Read(receivedtext)
 
    if err != nil {
      fmt.Printf("Received: %d bytes\n",n)
    }
 
    s := string(receivedtext[:n])
    fmt.Printf("Received: %d bytes: %s\n",n,s)
    //io.Copy(ws, ws)
    //fmt.Printf("Sent: %s\n",s)
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


var homepage string = `
<html>
<head>
<meta charset="UTF-8" />
<script>
 
  var serversocket = new WebSocket("ws://localhost:8085/echo");
 
  serversocket.onopen = function() {
    serversocket.send("Connection init");
  }
 
  // Write message on receive
  serversocket.onmessage = function(e) {
    document.getElementById('comms').innerHTML += "Received: " + e.data + "<br>";
  };
 
  function senddata() {
     var data = document.getElementById('sendtext').value;
     serversocket.send(data);
     document.getElementById('comms').innerHTML += "Sent: " + data + "<br>";
  }
 
</script>
 
</head>
<body>
  <input id="sendtext" type="text" />
  <input type="button" id="sendBtn" value="write" onclick="senddata()"></input>
 
  <div id='comms'></div>
 
</body>
</html>
`
