package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	flow "github.com/thinksystemio/package-flow"
	"github.com/thinksystemio/package-flow/command"
	response "github.com/thinksystemio/package-response"
	"nhooyr.io/websocket"
)

var tree = flow.NewTree()

var ws = WebSocket{
	Subscribers: map[*websocket.Conn]struct{}{},
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", NotFound)
	r.Get("/echo", Echo)

	r.Get("/api/echo", Echo)
	r.Get("/api/tree", HandleTree)
	r.Post("/api", HandleAPI)
	r.HandleFunc("/api/{node}", HandleWebsocketAPI)

	http.ListenAndServe(":80", r)
}

// Echo allows pinging of this service
func Echo(w http.ResponseWriter, r *http.Request) {
	response := response.CreateResponse()
	response.SendDataWithStatusCode(w, "echo", http.StatusOK)
}

// NotFound redirects to the not found page
func NotFound(w http.ResponseWriter, r *http.Request) {
	response := response.CreateResponse()
	response.SendDataWithStatusCode(w, "not found", http.StatusOK)
}

func HandleTree(w http.ResponseWriter, r *http.Request) {
	// TODO - REMOVE THIS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// TODO - REMOVE THIS
	ws.AddSubscriber(w, r)
}

func HandleWebsocketAPI(w http.ResponseWriter, r *http.Request) {
	// TODO - REMOVE THIS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// TODO - REMOVE THIS

	node := chi.URLParam(r, "node")
	cmd := &command.AddSubscriber{Node: node, W: w, R: r}
	flow.Dispatch(tree, cmd)
}

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	// TODO - REMOVE THIS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// TODO - REMOVE THIS

	response := response.CreateResponse()

	JSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.SendErrorWithStatusCode(w, err, http.StatusBadRequest)
	}

	data := flow.DispatchFromJSON(tree, JSON)
	ws.WriteSubscribers()

	response.SendDataWithStatusCode(w, data, http.StatusOK)
}
