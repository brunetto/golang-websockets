package main

import (
	"flag"
	"net/http"
	"io"
	"time"
	"code.google.com/p/go.net/websocket"
	"log"
)

// Echo everything back.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

// Send clock ticks.
func TimeServer(ws *websocket.Conn) {
	for {
		if _, err := ws.Write([]byte(time.Now().String())); err != nil {
			log.Println("Write: ", err)
			return
		}
		time.Sleep(1e+9)
	}
}

func greeter(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(homepage))
}

func main() {
	var port string
	flag.StringVar(&port, "http", ":12345", "[-http xxxx] # default 12345")
	flag.Parse()
	http.HandleFunc("/",  greeter)
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./"))))
	http.Handle("/echo", websocket.Handler(EchoServer))
	http.Handle("/time", websocket.Handler(TimeServer))
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: " + err.Error())
	}
}

var homepage string = `
<!DOCTYPE html>
<meta charset="utf-8" />
<title>Clock Server via WebSocket</title>
<script language="JavaScript">
wstest = (function () {
	var hostport = window.location.host
	var timeuri = "ws://"+hostport+"/time";
	var echouri = "ws://"+hostport+"/echo";
	var timesock;
	var echosock;
	var stats;
	var clock;
	var echoback;
	function updatestat(msg, e) {
		stats.innerHTML = '<span style="color: red;">'+msg+'</span> '+e.data;
	}

	function updateclock(e) {
		clock.innerHTML = '<span style="color: blue;">'+e.data+'</span>';
		echosock.send('<span style="color:green;">Sent by client: '+e.data+'</span>');
	}

	return	function () {
		MakeSock = (typeof WebSocket !== "undefined")?WebSocket:MozWebSocket;
		echosock = new MakeSock(echouri);
		echoback = document.getElementById("echoback");
		echosock.onmessage = function(e) { echoback.innerHTML=e.data };

		timesock = new MakeSock(timeuri);
		stats = document.getElementById("status");
		clock = document.getElementById("clock");
		timesock.onopen = function(e) { updatestat("Open:", e) };
		timesock.onclose = function(e) { updatestat("Close:", e) };
		timesock.onmessage = function(e) { updateclock(e) };
		timesock.onerror = function(e) { updatestat("Error:", e) };
	}
})();

window.addEventListener("load", wstest, false);
</script>
<h2>Time Reporting From Server</h2>
<p id="status"></p>
<p id="clock"></p>
<p id="echoback"></p>
</html>
`