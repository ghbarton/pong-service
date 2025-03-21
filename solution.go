package main

import (
	"context"
	"fmt"
	kitlog "github.com/go-kit/log"
	"github.com/philippseith/signalr"
	"github.com/philippseith/signalr/chatsample/public"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	runHTTPServer(context.TODO())
}

type AppHub struct {
	signalr.Hub
}

func (h *AppHub) SendChatMessage(message string) {
	h.Clients().All().Send("chatMessageReceived", message)
}

func runHTTPServer(ctx context.Context) {
	address := "localhost:8080"

	// create an instance of your hub
	hub := AppHub{}

	// build a signalr.Server using your hub
	// and any server options you may need
	server, _ := signalr.NewServer(
		context.TODO(),
		signalr.SimpleHubFactory(&hub),
		signalr.KeepAliveInterval(2*time.Second),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stderr), true),
	)

	// create a new http.ServerMux to handle your app's http requests
	router := http.NewServeMux()

	// ask the signalr server to map it's server
	// api routes to your custom baseurl
	server.MapHTTP(signalr.WithHTTPServeMux(router), "/chat")

	// in addition to mapping the signalr routes
	// your mux will need to serve the static files
	// which make up your client-side app, including
	// the signalr javascript files. here is an example
	// of doing that using a local `public` package
	// which was created with the go:embed directive
	//
	// fmt.Printf("Serving static content from the embedded filesystem\n")
	router.Handle("/", http.FileServer(http.FS(public.FS)))

	// bind your mux to a given address and start handling requests
	fmt.Printf("Listening for websocket connections on http://%s\n", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
