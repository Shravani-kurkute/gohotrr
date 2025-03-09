package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{}

func main() {
	port := "8080" // Default port
	if len(os.Args) > 1 {
		if _, err := strconv.Atoi(os.Args[1]); err == nil {
			port = os.Args[1]
		} else {
			log.Fatalf("Invalid port number: %s\n", os.Args[1])
		}
	}

	go watchFiles(".")

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", injectReloadMiddleware(fs))

	http.HandleFunc("/ws", wsHandler)

	log.Printf("Server started at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Watch for file changes
func watchFiles(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op == fsnotify.Write || event.Op == fsnotify.Create {
				log.Println("File changed:", event.Name)
				broadcastReload()
			}
		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}

// WebSocket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}
	}
}

// Broadcast reload message
func broadcastReload() {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func injectReloadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && !endsWith(r.URL.Path, ".html") {
			next.ServeHTTP(w, r)
			return
		}

		// Buffer the response
		buf := &bytes.Buffer{}
		recorder := &responseRecorder{ResponseWriter: w, body: buf}
		next.ServeHTTP(recorder, r)

		// Inject WebSocket reload script
		modifiedContent := buf.String() + `<script>
            let ws = new WebSocket("ws://" + location.host + "/ws");
            ws.onmessage = () => location.reload();
        </script>`

		w.Header().Set("Content-Length", strconv.Itoa(len(modifiedContent)))
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(recorder.statusCode)
		w.Write([]byte(modifiedContent))
	})
}

// Custom response recorder to capture response body
type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.body.Write(b) // Capture response in buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
