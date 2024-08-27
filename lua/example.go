package main

import (
	"C"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

var messages chan bool

const injected_js = `
  <script>
  var ws = new WebSocket("ws://localhost%s/ws");
   ws.onmessage = function(event) {
    if (event.data === "reload") {
      window.location.reload();
    }
  };
  </script>
`

const html = `
 <!DOCTYPE html>
  %s
<html>
<body>

<h1>My First Heading %d</h1>
<p>My first paragraph.</p>

</body>
</html> 
`

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	js := fmt.Sprintf(injected_js, ":2324")
	fmt.Fprintf(w, html, js, rand.IntN(100))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()
	for {
		update := <-messages
		if update {
			err := conn.WriteMessage(websocket.TextMessage, []byte("reload"))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
		time.Sleep(5 * time.Second)
		update = !update
	}
}

//export StartServer
func StartServer() {
	messages = make(chan bool)
	go func() {
		port_num := 2324 // rand.IntN(9999))
		//port_num := rand.IntN(9999-1000) + 1000
		port := fmt.Sprintf(":%d", port_num)
		server := fmt.Sprintf("http://localhost%s/", port)
		http.HandleFunc("/", handler)
		http.HandleFunc("/ws", wsHandler)
		fmt.Printf("Starting server on %s\n", port)
		if err := http.ListenAndServe(port, nil); err != nil {
			fmt.Println("Failed to start server:", err)
		}
		openBrowser(server)
	}()
}

func main() {}

//export SendUpdate
func SendUpdate() {
	messages <- true
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		panic(err)
	}

}
