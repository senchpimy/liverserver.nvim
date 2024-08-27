package main

import (
	"C"
	"fmt"
	_ "io"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
)
import "time"

var messages chan bool

const injected_js = `
  <script>
  var i = 0
  var ws = new WebSocket("ws://localhost%s/ws");
   ws.onmessage = function(event) {
    if (event.data === "reload") {
      //window.location.reload();
  console.log("Actualizacion"+i)
  i+=1
    }
  };
  </script>
`

const html = ""

func dummy(w http.ResponseWriter, r *http.Request) {
	js := fmt.Sprintf(injected_js, ":2324")
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, html, js, path, rand.IntN(100))

}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	json := fmt.Sprintf(injected_js, ":2324")
	local_path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	patterns := []string{
		`(?i)(<head\b[^>]*>)`,
		`(?i)(<body\b[^>]*>)`,
		`(?i)(<html\b[^>]*>)`,
	}

	var local_html string
	if path == "/" || strings.HasSuffix(path, ".html") {
		if path == "/" {
			local_html = request("index.html")
			local_html = json + local_html
		} else if strings.HasSuffix(path, ".html") {
			local_html = request(local_path + path)
			local_html = json + local_html
			for _, pattern := range patterns {
				re := regexp.MustCompile(pattern)
				if re.MatchString(local_html) {
					local_html = re.ReplaceAllString(local_html, json+"$1")
					break
				}
			}
		}
	} else {
		local_html = request(local_path + path)
	}
	fmt.Fprint(w, local_html)
}

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	defer conn.Close()
	time.Sleep(1 * time.Second)
	for {
		update := <-messages
		if update {
			err := conn.WriteMessage(websocket.TextMessage, []byte("reload"))
			if err != nil {
				fmt.Println(err)
			}
		}
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
		//fmt.Printf("Starting server on %s\n", port)
		if err := http.ListenAndServe(port, nil); err != nil {
			panic(err)
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

func request(url string) (foo string) {
	resp, err := os.ReadFile(url)
	if err != nil {
		return ""
	}
	return string(resp)
}
